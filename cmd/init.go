package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/devpod-provider-stackit/pkg/options"
	"github.com/stackitcloud/devpod-provider-stackit/pkg/stackit"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		options, err := options.FromEnv(true)
		if err != nil {
			return err
		}

		return stackit.New(options.ClientOptions).Init(context.Background(), options.ProjectID)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
