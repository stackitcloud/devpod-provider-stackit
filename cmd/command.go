package cmd

import (
	"context"
	"net"
	"os"

	"github.com/loft-sh/devpod/pkg/ssh"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/stackitcloud/devpod-provider-stackit/pkg/options"
	"github.com/stackitcloud/devpod-provider-stackit/pkg/stackit"
)

// commandCmd represents the command command
var commandCmd = &cobra.Command{
	Use:   "command",
	Short: "Run a command on the instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		options, err := options.FromEnv(false)
		if err != nil {
			return err
		}

		ctx := context.Background()

		command := os.Getenv("COMMAND")
		if command == "" {
			return errors.New("command environment variable is missing")
		}

		// Get private key
		privateKey, err := ssh.GetPrivateKeyRawBase(options.MachineFolder)
		if err != nil {
			return errors.Wrap(err, "load private key")
		}

		// Create SSH client
		publicIP, err := stackit.New(options.ClientOptions).GetPublicIPOfServer(ctx, options.ProjectID, options.MachineID)
		if err != nil {
			return errors.Wrap(err, "get public IP")
		}

		// Call external address
		sshClient, err := ssh.NewSSHClient(stackit.SSHUserName, net.JoinHostPort(publicIP, stackit.SSHPort), privateKey)
		if err != nil {
			return errors.Wrap(err, "create ssh client")
		}
		defer func() {
			err = sshClient.Close()
			if err != nil {
				err = errors.Wrap(err, "close ssh client")
			}
		}()

		// Run command
		if err := ssh.Run(ctx, sshClient, command, os.Stdin, os.Stdout, os.Stderr, nil); err != nil {
			return errors.Wrap(err, "run ssh command")
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(commandCmd)
}
