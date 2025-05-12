package stackit

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/loft-sh/devpod/pkg/ssh"
	"github.com/pkg/errors"
	"github.com/stackitcloud/stackit-sdk-go/core/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"

	"github.com/stackitcloud/devpod-provider-stackit/pkg/options"
)

func (s *Stackit) createNetwork(ctx context.Context, projectId string, networkName string) (*iaas.Network, error) {
	createNetworkPayload := iaas.CreateNetworkPayload{
		Name: utils.Ptr(networkName),
		AddressFamily: &iaas.CreateNetworkAddressFamily{
			Ipv4: &iaas.CreateNetworkIPv4Body{
				PrefixLength: utils.Ptr(int64(29)), // 29 is the largest possible prefix length
			},
		},
	}

	network, err := s.client.CreateNetwork(ctx, projectId).CreateNetworkPayload(createNetworkPayload).Execute()
	if err != nil {
		return nil, err
	}

	network, err = wait.CreateNetworkWaitHandler(ctx, s.client, projectId, *network.NetworkId).WaitWithContext(ctx)
	if err != nil {
		return nil, err
	}

	return network, nil
}

func (s *Stackit) deleteNetwork(ctx context.Context, projectId string, networkId string) error {
	err := s.client.DeleteNetworkExecute(ctx, projectId, networkId)
	if err != nil {
		return err
	}
	_, err = wait.DeleteNetworkWaitHandler(ctx, s.client, projectId, networkId).WaitWithContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Stackit) waitForSSHToBeReady(options *options.Options, publicIP, sshUserName, sshPort string) error {
	// Get private key
	privateKey, err := ssh.GetPrivateKeyRawBase(options.MachineFolder)
	if err != nil {
		return errors.Wrap(err, "load private key")
	}

	sshSuccess := false
	for !sshSuccess {
		sshClient, err := ssh.NewSSHClient(sshUserName, net.JoinHostPort(publicIP, sshPort), privateKey)
		if err == nil {
			sshSuccess = true
			err := sshClient.Close()
			if err != nil {
				return err
			}
		} else {
			fmt.Println("SSH not ready...")
			time.Sleep(5 * time.Second)
		}
	}
	return nil
}
