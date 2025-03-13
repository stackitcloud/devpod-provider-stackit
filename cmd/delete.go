package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/devpod-provider-stackit/pkg/options"
	"github.com/stackitcloud/devpod-provider-stackit/pkg/stackit"
)

// deleteCmd represents the delete command.
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an instance and volume",
	RunE: func(_ *cobra.Command, args []string) error {
		options, err := options.FromEnv(false)
		if err != nil {
			return err
		}

		ctx := context.Background()
		s := stackit.New(options.ClientOptions)

		return s.Delete(ctx, options.ProjectID, options.MachineID)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
