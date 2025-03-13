package cmd

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/loft-sh/devpod/pkg/ssh"
	"github.com/spf13/cobra"

	"github.com/stackitcloud/devpod-provider-stackit/pkg/options"
	"github.com/stackitcloud/devpod-provider-stackit/pkg/stackit"
)

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

		return s.Create(ctx, options, publicKey)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
