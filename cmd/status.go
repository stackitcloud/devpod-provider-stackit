package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/devpod-provider-stackit/pkg/options"
	"github.com/stackitcloud/devpod-provider-stackit/pkg/stackit"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Retrieve the status of an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		options, err := options.FromEnv(false)
		if err != nil {
			return err
		}

		status, err := stackit.New(options.ClientOptions).Status(context.Background(), options.ProjectID, options.MachineID)
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(os.Stdout, status)
		return err
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
