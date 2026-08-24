package cmd

import (
	"fmt"

	"github.com/operaodev/umafetch/internal"
	"github.com/spf13/cobra"
)

const logo = `
██╗   ██╗███╗   ███╗ █████╗ ███████╗███████╗████████╗ ██████╗██╗  ██╗
██║   ██║████╗ ████║██╔══██╗██╔════╝██╔════╝╚══██╔══╝██╔════╝██║  ██║
██║   ██║██╔████╔██║███████║█████╗  █████╗     ██║   ██║     ███████║
██║   ██║██║╚██╔╝██║██╔══██║██╔══╝  ██╔══╝     ██║   ██║     ██╔══██║
╚██████╔╝██║ ╚═╝ ██║██║  ██║██║     ███████╗   ██║   ╚██████╗██║  ██║
 ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝     ╚══════╝   ╚═╝    ╚═════╝╚═╝  ╚═╝
                                                        v1.0.0

                                                        
`

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Download and configure umas for fastfetch",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		fmt.Fprint(out, logo)
		if internal.UmasExist() {
			fmt.Fprint(out, "Umas already downloaded. Do you want to update? (y/n): ")
			var answer string
			fmt.Scanln(&answer)
			if answer != "y" && answer != "Y" {
				fmt.Fprintln(out, "Cancelled.")
				return nil
			}
		}

		fmt.Fprintln(out, "Downloading umas...")
		if err := internal.SaveUmas(); err != nil {
			return fmt.Errorf("error downloading umas: %w", err)
		}
		fmt.Fprintln(out, "Umas downloaded successfully.")

		fmt.Fprintln(out, "Generating default configuration...")
		internal.GenerateDefaultConfig()
		internal.GenerateDefaultTemplateLarge()
		internal.GenerateDefaultTemplateSmall()
		fmt.Fprintln(out, "Configuration generated.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
