package cmd

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/loft-sh/log"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"net"
	"time"

	"github.com/loft-sh/devpod/pkg/ssh"
	"github.com/spf13/cobra"

	"github.com/stackitcloud/devpod-provider-stackit/pkg/options"
	"github.com/stackitcloud/devpod-provider-stackit/pkg/stackit"
)

const maxConnectionAttempts = 42

type cloudInit struct {
	Status string `json:"status"`
}

func checkConnectionStatus(ctx context.Context, publicIP string, privateKey *[]byte) bool {
	// Call external address
	sshClient, err := ssh.NewSSHClient(stackit.SSHUserName, net.JoinHostPort(publicIP, stackit.SSHPort), *privateKey)
	if err != nil {
		log.Default.Debugf("Error creating ssh client: %v", err)
		return false
	}
	defer func() {
		if err := sshClient.Close(); err != nil {
			log.Default.Debugf("Error closing ssh client: %v", err)
		}
	}()

	buf := new(bytes.Buffer)
	if err := ssh.Run(ctx, sshClient, "cloud-init status || true", &bytes.Buffer{}, buf, &bytes.Buffer{}, nil); err != nil {
		log.Default.Errorf("Error retrieving cloud-init status, %v", err)
		return false
	}

	var status cloudInit
	if err := yaml.Unmarshal(buf.Bytes(), &status); err != nil {
		log.Default.Errorf("Unable to parse cloud-init YAML: %v", err)
		return false
	}

	if status.Status != "done" {
		return false
	}

	return true
}

// createCmd represents the create command.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		options, err := options.FromEnv(false)
		if err != nil {
			return err
		}

		ctx := context.Background()
		s := stackit.New(options.ClientOptions)

		publicKeyBase, err := ssh.GetPublicKeyBase(options.MachineFolder)
		if err != nil {
			return fmt.Errorf("can't get public key: %w", err)
		}

		publicKey, err := base64.StdEncoding.DecodeString(publicKeyBase)
		if err != nil {
			return fmt.Errorf("can't b64decode public key: %w", err)
		}

		if err := s.Create(ctx, options, publicKey); err != nil {
			return err
		}

		privateKey, err := ssh.GetPrivateKeyRawBase(options.MachineFolder)
		if err != nil {
			return errors.Wrap(err, "load private key")
		}

		publicIP, err := stackit.New(options.ClientOptions).GetPublicIPOfServer(ctx, options.ProjectID, options.MachineID)
		if err != nil {
			return errors.Wrap(err, "get public IP")
		}

		attempt := 0

		for {
			if attempt >= maxConnectionAttempts {
				return errors.New("max connection attempts reached")
			}
			attempt++

			log.Default.Debugf("Attempt %d of %d to connect to server using ssh", attempt, maxConnectionAttempts)

			time.Sleep(time.Second)

			if checkConnectionStatus(ctx, publicIP, &privateKey) {
				log.Default.Info("Successfully connected to server")
				break
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
