package cmd

import (
	"fmt"

	"github.com/operaodev/umafetch/internal"
	"github.com/spf13/cobra"
)

var (
	flagGenTemplate string
	flagSwitch      bool
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage fastfetch templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		cfg, err := internal.ConfigLoad()
		if err != nil {
			return err
		}

		if flagSwitch {
			current := cfg.Template
			var next internal.Template
			if current == internal.TemplateLarge {
				next = internal.TemplateSmall
			} else {
				next = internal.TemplateLarge
			}
			cfg.Template = next
			cfg.ConfigSave()
			fmt.Fprintf(out, "Template switched: %s -> %s\n", current, next)
			return nil
		}

		if flagGenTemplate != "" {
			switch flagGenTemplate {
			case "small":
				fmt.Fprint(out, "Generate small template? (y/n): ")
			case "large":
				fmt.Fprint(out, "Generate large template? (y/n): ")
			default:
				return fmt.Errorf("unknown template: %s (use 'small' or 'large')", flagGenTemplate)
			}
			var answer string
			fmt.Scanln(&answer)
			if answer != "y" && answer != "Y" {
				fmt.Fprintln(out, "Cancelled.")
				return nil
			}

			switch flagGenTemplate {
			case "small":
				if err := internal.GenerateDefaultTemplateSmall(); err != nil {
					return err
				}
				fmt.Fprintln(out, "Small template generated.")
			case "large":
				if err := internal.GenerateDefaultTemplateLarge(); err != nil {
					return err
				}
				fmt.Fprintln(out, "Large template generated.")
			}
			return nil
		}

		fmt.Fprint(out, "Generate both templates? (y/n): ")
		var answer string
		fmt.Scanln(&answer)
		if answer != "y" && answer != "Y" {
			fmt.Fprintln(out, "Cancelled.")
			return nil
		}

		if err := internal.GenerateDefaultTemplateLarge(); err != nil {
			return err
		}
		if err := internal.GenerateDefaultTemplateSmall(); err != nil {
			return err
		}
		fmt.Fprintln(out, "Both templates generated.")
		return nil
	},
}

func init() {
	templateCmd.Flags().StringVar(&flagGenTemplate, "gen-template", "", "Generate template (small or large)")
	templateCmd.Flags().BoolVar(&flagSwitch, "switch", false, "Switch between small and large template")
	rootCmd.AddCommand(templateCmd)
}
