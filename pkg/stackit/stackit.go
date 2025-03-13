package stackit

import (
	"bytes"
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"github.com/loft-sh/log"
	"text/template"

	"github.com/stackitcloud/stackit-sdk-go/core/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"

	"github.com/loft-sh/devpod/pkg/client"
	"github.com/pkg/errors"
	"github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/stackitcloud/devpod-provider-stackit/pkg/options"
)

//go:embed cloud-config.yaml
var cloudConfigFS embed.FS

const (
	SSHUserName   = "devpod"
	SSHPort       = "22"
	ubuntuImageID = "117e8764-41c2-405f-aece-b53aa08b28cc"
)

type Stackit struct {
	client *iaas.APIClient
}

func New(options *options.ClientOptions) *Stackit {
	apiClient, err := iaas.NewAPIClient(config.WithRegion(options.Region))
	if err != nil {
		panic(err)
	}
	return &Stackit{
		client: apiClient,
	}
}

func (s *Stackit) Init(ctx context.Context, projectID string) error {
	_, err := s.client.ListServers(ctx, projectID).Execute()
	if err != nil {
		return err
	}
	return nil
}

func (s *Stackit) GetPublicIPOfServer(ctx context.Context, projectId, machineName string) (string, error) {
	server, err := s.getServerByName(ctx, projectId, machineName)
	if err != nil {
		return "", err
	}
	if server == nil {
		return "", errors.New("server not found")
	}
	if server.Nics == nil {
		return "", errors.New("no networks found")
	}
	for _, network := range *server.Nics {
		if network.PublicIp != nil {
			return *network.PublicIp, nil
		}
	}
	return "", errors.New("no public IP found")
}

func (s *Stackit) Status(ctx context.Context, projectId, machineName string) (client.Status, error) {
	server, err := s.getServerByName(ctx, projectId, machineName)
	if err != nil {
		return client.StatusNotFound, nil
	}

	if server == nil {
		return client.StatusNotFound, nil
	}
	if server.GetPowerStatus() == nil {
		return client.StatusNotFound, nil
	}

	return statusFromPowerStateString(*server.GetPowerStatus()), nil
}

func (s *Stackit) Start(ctx context.Context, projectId, machineName string) error {
	server, err := s.getServerByName(ctx, projectId, machineName)
	if err != nil {
		return err
	}
	err = s.client.StartServer(ctx, projectId, *server.Id).Execute()
	if err != nil {
		return err
	}
	return nil
}

func (s *Stackit) Stop(ctx context.Context, projectId, machineName string) error {
	server, err := s.getServerByName(ctx, projectId, machineName)
	if err != nil {
		return err
	}
	err = s.client.StopServer(ctx, projectId, *server.Id).Execute()
	if err != nil {
		return err
	}
	return nil
}

func (s *Stackit) Delete(ctx context.Context, projectId, machineName string) error {
	server, err := s.getServerByName(ctx, projectId, machineName)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	gotError := false

	if len(*server.Nics) == 0 {
		log.Default.Errorf("no networks found for server")
		gotError = true
	} else {
		nics := *server.Nics
		networkId := nics[0].NetworkId
		err = s.client.DeleteNetwork(ctx, projectId, *networkId).Execute()
		if err != nil {
			log.Default.Errorf("failed to delete network")
			gotError = true
		}

		publicIPAddress := nics[0].PublicIp
		publicIP, err := s.getPublicIPByIPAddress(ctx, projectId, *publicIPAddress)
		if err != nil {
			log.Default.Errorf("failed to get public IP")
			gotError = true
		}

		err = s.client.DeletePublicIP(ctx, projectId, *publicIP.Id).Execute()
		if err != nil {
			log.Default.Errorf("failed to delete public IP")
			gotError = true
		}

		for _, id := range *nics[0].SecurityGroups {
			name, err := s.getSecurityGroupNameByID(ctx, projectId, id)
			if err != nil {
				log.Default.Errorf("failed to get name of security group by ID")
				gotError = true
			}

			if name != "default" {
				err = s.client.DeleteSecurityGroup(ctx, projectId, id).Execute()
				if err != nil {
					log.Default.Errorf("failed to delete security group %q", name)
					gotError = true
				}
			}
		}
	}

	if gotError {
		return errors.New("failed to delete all components associated to the devpod")
	}

	return nil
}

func statusFromPowerStateString(state string) client.Status {
	switch state {
	case "CRASHED":
		return client.StatusStopped
	case "STOPPED":
		return client.StatusStopped
	case "RUNNING":
		return client.StatusRunning
	case "ERROR":
		return client.StatusStopped
	}
	return client.StatusNotFound
}

func (s *Stackit) Create(ctx context.Context, options *options.Options, publicKey []byte) error {
	sshPublicKey := string(publicKey)

	network, err := s.createNetwork(ctx, options.ProjectID, options.MachineID)
	if err != nil {
		return err
	}

	userdata, err := generateUserData(sshPublicKey)
	if err != nil {
		return err
	}

	createServerPayload := iaas.CreateServerPayload{
		Name:             &options.MachineID,
		AvailabilityZone: &options.AvailabilityZone,
		MachineType:      &options.Flavor,
		BootVolume: &iaas.CreateServerPayloadBootVolume{
			DeleteOnTermination: utils.Ptr(true),
			Size:                utils.Ptr(int64(64)),
			Source: &iaas.BootVolumeSource{
				Id:   utils.Ptr(ubuntuImageID),
				Type: utils.Ptr("image"),
			},
		},
		Networking: &iaas.CreateServerPayloadNetworking{
			CreateServerNetworking: &iaas.CreateServerNetworking{
				NetworkId: network.NetworkId,
			},
		},
		UserData: &userdata,
	}

	server, err := s.client.CreateServer(ctx, options.ProjectID).CreateServerPayload(createServerPayload).Execute()
	if err != nil {
		return err
	}

	server, err = wait.CreateServerWaitHandler(ctx, s.client, options.ProjectID, *server.Id).WaitWithContext(ctx)
	if err != nil {
		return err
	}

	publicIP, err := s.client.CreatePublicIP(ctx, options.ProjectID).CreatePublicIPPayload(iaas.CreatePublicIPPayload{}).Execute()
	if err != nil {
		return err
	}

	err = s.client.AddPublicIpToServer(ctx, options.ProjectID, *server.Id, *publicIP.Id).Execute()
	if err != nil {
		return err
	}

	createSecurityGroupPayload := iaas.CreateSecurityGroupPayload{
		Name: &options.MachineID,
	}

	securityGroup, err := s.client.CreateSecurityGroup(ctx, options.ProjectID).CreateSecurityGroupPayload(createSecurityGroupPayload).Execute()
	if err != nil {
		return err
	}

	createSecurityGroupRulePayload := iaas.CreateSecurityGroupRulePayload{
		Description: utils.Ptr("SSH"),
		Direction:   utils.Ptr("ingress"),
		PortRange: &iaas.PortRange{
			Min: utils.Ptr(int64(22)),
			Max: utils.Ptr(int64(22)),
		},
		Protocol: &iaas.CreateProtocol{
			String: utils.Ptr("tcp"),
		},
	}

	securityGroupRule, err := s.client.CreateSecurityGroupRule(ctx, options.ProjectID, *securityGroup.Id).CreateSecurityGroupRulePayload(createSecurityGroupRulePayload).Execute()
	if err != nil {
		return err
	}

	fmt.Println(*securityGroupRule)

	err = s.client.AddSecurityGroupToServer(ctx, options.ProjectID, *server.Id, *securityGroup.Id).Execute()
	if err != nil {
		return err
	}

	return nil
}

func generateUserData(publicKey string) (string, error) {
	t, err := template.New("cloud-config.yaml").ParseFS(cloudConfigFS, "cloud-config.yaml")
	if err != nil {
		return "", err
	}

	output := new(bytes.Buffer)
	if err := t.Execute(output, map[string]string{
		"PublicKey": publicKey,
		"Username":  SSHUserName,
	}); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(output.Bytes()), nil
}

func (s *Stackit) createVolume(ctx context.Context, projectId, volumeName, volumeAvailabilityZone string, volumeSize int64) (string, error) {
	createVolumePayload := iaas.CreateVolumePayload{
		Name:             &volumeName,
		AvailabilityZone: &volumeAvailabilityZone,
		Size:             &volumeSize,
	}

	volume, err := s.client.CreateVolume(ctx, projectId).CreateVolumePayload(createVolumePayload).Execute()
	if err != nil {
		return "", err
	}

	return *volume.Id, nil
}

func (s *Stackit) getServerByName(ctx context.Context, projectId, serverName string) (*iaas.Server, error) {
	servers, err := s.client.ListServers(ctx, projectId).Details(true).Execute()
	if err != nil {
		return nil, err
	}
	if servers == nil {
		return nil, errors.New("server not found")
	}
	if servers.Items == nil {
		return nil, errors.New("servers not found")
	}

	for _, server := range *servers.Items {
		if *server.Name == serverName {
			return &server, nil
		}
	}
	return nil, errors.New("server not found")
}

func (s *Stackit) getPublicIPByIPAddress(ctx context.Context, projectId, ipAddress string) (*iaas.PublicIp, error) {
	publicIPs, err := s.client.ListPublicIPs(ctx, projectId).Execute()
	if err != nil {
		return nil, err
	}

	if len(*publicIPs.Items) == 0 {
		return nil, errors.New("no public IPs found")
	}

	for _, publicIP := range *publicIPs.Items {
		if *publicIP.Ip == ipAddress {
			return &publicIP, nil
		}
	}

	return nil, errors.New("public IP not found")
}

func (s *Stackit) getSecurityGroupNameByID(ctx context.Context, projectId, securityGroupId string) (string, error) {
	securityGroup, err := s.client.GetSecurityGroup(ctx, projectId, securityGroupId).Execute()
	if err != nil {
		return "", errors.New("security group not found")
	}

	return *securityGroup.Name, nil
}
