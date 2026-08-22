package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

const (
	defaultSlogan  = "Umamusume! Pretty Derby - Where dreams are won on the track!"
	defaultProfile = "A horse girl who runs with all her heart. Born from the legend of the racetrack, she strives to become the ultimate champion."
)

// RenderUma renderiza una UMA usando el template y la configuración proporcionados
func (uma Uma) RenderUma(template string, cfg Config) (string, error) {
	slogan := uma.Slogan
	if slogan == "" {
		slogan = defaultSlogan
	}

	profile := uma.Profile
	if profile == "" {
		profile = defaultProfile
	}

	sloganLines := buildTextBlock("Slogan", slogan, uma.SubColor, 35, 45)
	bioLines := buildTextBlock("Profile", profile, uma.SubColor, 35, 45)
	separator := separatorModule(uma.MainColor, cfg)

	result := template
	result = strings.ReplaceAll(result, "{Image}", uma.Image)
	result = strings.ReplaceAll(result, "{PrimaryColor}", uma.MainColor)
	result = strings.ReplaceAll(result, "{SecondaryColor}", uma.SubColor)
	result = strings.ReplaceAll(result, "{Name}", uma.Name)
	result = strings.ReplaceAll(result, "{Title}", uma.Title)
	result = strings.ReplaceAll(result, "{SloganLines}", sloganLines)
	result = strings.ReplaceAll(result, "{BioLines}", bioLines)
	result = strings.ReplaceAll(result, "{separator}", separator)

	dir, err := os.MkdirTemp("", "umafetch_*")
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, "config.jsonc")
	if err := os.WriteFile(path, []byte(result), 0644); err != nil {
		return "", err
	}

	return path, nil
}

func separatorModule(color string, cfg Config) string {
	return fmt.Sprintf(`{"type":"custom","format":"%s","outputColor":"%s"}`, cfg.Separator.Build(), color)
}

func buildTextBlock(key, text, color string, firstLen, restLen int) string {
	lines := wrapLines(text, firstLen, restLen)

	var modules []string
	for i, line := range lines {
		if i == 0 {
			modules = append(modules, fmt.Sprintf(
				`{"type":"custom","key":%s,"keyColor":%s,"format":"%s"}`,
				escJSON(key), escJSON(color), line,
			))
		} else {
			modules = append(modules, fmt.Sprintf(
				`{"type":"custom","format":"%s"}`, line,
			))
		}
	}

	return strings.Join(modules, ",\n        ")
}

func escJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func wrapLines(text string, firstLimit, restLimit int) []string {
	var lines []string
	remaining := text

	for idx := 0; utf8.RuneCountInString(remaining) > 0; idx++ {
		limit := restLimit
		if idx == 0 {
			limit = firstLimit
		}

		runeCount := utf8.RuneCountInString(remaining)
		if runeCount <= limit {
			lines = append(lines, remaining)
			break
		}

		cut := limit
		for i := limit - 1; i >= limit/2; i-- {
			r, _ := utf8.DecodeRuneInString(remaining[i:])
			if r == ' ' {
				cut = i
				break
			}
		}
		lines = append(lines, remaining[:cut])
		remaining = strings.TrimLeft(remaining[cut:], " ")
	}

	return lines
}
