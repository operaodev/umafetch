package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "umafetch",
	Short:         "Umamusume - Fastfetch theme",
	Version:       "1.0.2",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.SetVersionTemplate(`v{{.Version}}`)
}

func Execute() error {
	return rootCmd.Execute()
}
