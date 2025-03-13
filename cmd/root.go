package cmd

import (
	"os"

	"github.com/loft-sh/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devpod-provider-stackit",
	Short: "DevPod on STACKIT",
	PersistentPreRunE: func(cobraCmd *cobra.Command, args []string) error {
		log.Default.MakeRaw()

		logLevel := os.Getenv("DEVPOD_LOG_LEVEL")
		if logLevel != "" {
			if lvl, err := logrus.ParseLevel(logLevel); err != nil {
				log.Default.Error(errors.Wrap(err, "invalid log level provided, continuing"))
			} else {
				log.Default.SetLevel(lvl)
			}
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
