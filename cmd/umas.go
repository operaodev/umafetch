package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/operaodev/umafetch/internal"
	"github.com/spf13/cobra"
)

var (
	flagTarget string
	flagFull   bool
	flagSearch string
)

var umasCmd = &cobra.Command{
	Use:   "umas",
	Short: "List and select umas",
	RunE: func(cmd *cobra.Command, args []string) error {
		umas, err := internal.FindUmas()
		if err != nil {
			return err
		}

		seen := make(map[int]bool)
		for _, u := range umas {
			seen[u.ID] = true
		}

		if flagSearch != "" {
			return listByName(cmd, umas, flagSearch)
		}

		if flagTarget != "" {
			return selectUma(cmd, umas)
		}

		if flagFull {
			return listFull(cmd, umas, seen)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Total umas: %d, outfits: %d\n", len(seen), len(umas))
		return nil
	},
}

func listFull(cmd *cobra.Command, umas []internal.Uma, seen map[int]bool) error {
	type charInfo struct {
		ID      int
		Name    string
		Outfits int
	}

	chars := make(map[int]*charInfo)
	for _, u := range umas {
		if c, ok := chars[u.ID]; ok {
			c.Outfits++
		} else {
			chars[u.ID] = &charInfo{ID: u.ID, Name: u.Name, Outfits: 1}
		}
	}

	var list []*charInfo
	for _, c := range chars {
		list = append(list, c)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})

	out := cmd.OutOrStdout()
	for i, c := range list {
		fmt.Fprintf(out, "%d. %s - %d Outfits\n", i+1, c.Name, c.Outfits)
	}
	return nil
}

func listByName(cmd *cobra.Command, umas []internal.Uma, name string) error {
	target := strings.ToLower(name)

	var matches []internal.Uma
	for _, u := range umas {
		if strings.Contains(strings.ToLower(u.Name), target) {
			matches = append(matches, u)
		}
	}

	if len(matches) == 0 {
		return fmt.Errorf("no uma found with name '%s'", name)
	}

	type charGroup struct {
		Name    string
		Outfits []internal.Uma
	}

	seen := make(map[int]*charGroup)
	var order []int
	for _, u := range matches {
		if g, ok := seen[u.ID]; ok {
			g.Outfits = append(g.Outfits, u)
		} else {
			seen[u.ID] = &charGroup{Name: u.Name, Outfits: []internal.Uma{u}}
			order = append(order, u.ID)
		}
	}

	out := cmd.OutOrStdout()
	idx := 1
	for _, id := range order {
		g := seen[id]
		sort.Slice(g.Outfits, func(i, j int) bool {
			return g.Outfits[i].Order < g.Outfits[j].Order
		})
		fmt.Fprintf(out, "%d. %s (%d)\n", idx, g.Name, len(g.Outfits))
		for _, u := range g.Outfits {
			if flagFull {
				main := paintBG(u.MainColor, fmt.Sprintf(" %s ", u.MainColor))
				sub := paintBG(u.SubColor, fmt.Sprintf(" %s ", u.SubColor))
				fmt.Fprintf(out, "   (%d). %s - %s %s\n", u.Order, u.Title, main, sub)
			} else {
				fmt.Fprintf(out, "   (%d). %s\n", u.Order, u.Title)
			}
		}
		idx++
	}
	return nil
}

func selectUma(cmd *cobra.Command, umas []internal.Uma) error {
	parts := strings.SplitN(flagTarget, ",", 2)
	name := strings.TrimSpace(parts[0])

	var order *int
	if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
		n, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return fmt.Errorf("invalid order number: %s", parts[1])
		}
		order = &n
	}

	if order != nil {
		var matches []internal.Uma
		for _, u := range umas {
			if strings.EqualFold(u.Name, name) {
				matches = append(matches, u)
			}
		}
		if len(matches) == 0 {
			return fmt.Errorf("no uma found with name '%s'", name)
		}
		found := false
		for _, u := range matches {
			if u.Order == *order {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("outfit %d not found for '%s'", *order, name)
		}
	}

	cfg, err := internal.ConfigLoad()
	if err != nil {
		return err
	}

	cfg.Theme.Name = &name
	cfg.Theme.Outfit = order
	cfg.ConfigSave()

	out := cmd.OutOrStdout()
	if order != nil {
		fmt.Fprintf(out, "Selected: %s (outfit=%d)\n", name, *order)
	} else {
		fmt.Fprintf(out, "Selected: %s (random outfit)\n", name)
	}
	return nil
}

func paintBG(hex, text string) string {
	if hex == "" || len(hex) != 7 {
		return text
	}
	var r, g, b byte
	fmt.Sscanf(hex, "#%02X%02X%02X", &r, &g, &b)
	return fmt.Sprintf("\033[48;2;%d;%d;%dm%s\033[0m", r, g, b, text)
}

func init() {
	umasCmd.Flags().BoolVarP(&flagFull, "full", "f", false, "Show full list with outfit details")
	umasCmd.Flags().StringVar(&flagSearch, "search", "", "Search uma by name")
	umasCmd.Flags().StringVar(&flagTarget, "target", "", "Select uma by name and optional order (name or name,order)")
	rootCmd.AddCommand(umasCmd)
}
