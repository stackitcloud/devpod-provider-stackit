package cmd

import (
	"context"
	"time"

	"github.com/loft-sh/devpod/pkg/client"
	"github.com/loft-sh/log"
	"github.com/spf13/cobra"

	"github.com/stackitcloud/devpod-provider-stackit/pkg/options"
	"github.com/stackitcloud/devpod-provider-stackit/pkg/stackit"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		options, err := options.FromEnv(false)
		if err != nil {
			return err
		}

		ctx := context.Background()

		stackitClient := stackit.New(options.ClientOptions)
		err = stackitClient.Stop(ctx, options.ProjectID, options.MachineID)

		if err != nil {
			return err
		}

		// Wait until it's stopped
		for {
			status, err := stackitClient.Status(ctx, options.ProjectID, options.MachineID)
			if err != nil {
				log.Default.Errorf("Error retrieving server status: %v", err)
				break
			} else if status == client.StatusStopped {
				break
			}

			time.Sleep(time.Second)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
