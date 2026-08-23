package cmd

import (
	"fmt"
	"strconv"

	"github.com/operaodev/umafetch/internal"
	"github.com/spf13/cobra"
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render uma theme and print the config path",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !internal.UmasExist() {
			return fmt.Errorf("umas not downloaded, run 'umafetch install' first")
		}

		cfg, err := internal.ConfigLoad()
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}

		flags := cmd.Flags()

		if t := flags.Lookup("template"); t != nil && t.Changed {
			switch t.Value.String() {
			case "large":
				cfg.Template = internal.TemplateLarge
			case "small":
				cfg.Template = internal.TemplateSmall
			default:
				return fmt.Errorf("unknown template: %s (use 'large' or 'small')", t.Value.String())
			}
		}

		if n := flags.Lookup("name"); n != nil && n.Changed {
			name := n.Value.String()
			cfg.Theme.Name = &name
		}

		if o := flags.Lookup("outfit"); o != nil && o.Changed {
			n, err := strconv.Atoi(o.Value.String())
			if err != nil {
				return fmt.Errorf("invalid outfit number: %s", o.Value.String())
			}
			cfg.Theme.Outfit = &n
		}

		uma, err := internal.FindUma(cfg)
		if err != nil {
			return fmt.Errorf("error finding uma: %w", err)
		}

		tmpl, err := internal.FindTemplate(cfg)
		if err != nil {
			return fmt.Errorf("error finding template: %w", err)
		}

		configPath, err := uma.RenderUma(tmpl, cfg)
		if err != nil {
			return fmt.Errorf("error rendering: %w", err)
		}

		// Print only the path to stdout so it can be captured by the shell
		fmt.Println(configPath)
		return nil
	},
}

func init() {
	renderCmd.Flags().String("template", "", "Override template (large or small)")
	renderCmd.Flags().String("name", "", "Override uma name")
	renderCmd.Flags().String("outfit", "", "Override outfit number")
	rootCmd.AddCommand(renderCmd)
}
