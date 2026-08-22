package internal

import (
	"fmt"
	"image"
	_ "image/png"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
)

const AppDir = "umafetch"

func imageDir() (string, error) {
	localDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(localDir, AppDir, "images")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return dir, nil
}

// downloadImage descarga la imagen de una UMA y la guarda en el directorio de imágenes.
func downloadImage(uma *Uma) error {
	dir, err := imageDir()
	if err != nil {
		return err
	}

	url := uma.ImageUrl()
	path := filepath.Join(dir, filepath.Base(url))

	if _, err := os.Stat(path); err == nil {
		uma.Image = path
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("image not found: %s (status %d)", url, resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "image/png" {
		return fmt.Errorf("unexpected content type %q for %s", ct, url)
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		os.Remove(path)
		return err
	}

	uma.Image = path
	return nil
}

// extractColors extrae los colores principales y secundarios de una imagen de una UMA.
func extractColors(uma *Uma) error {
	if uma.OutfitID == nil {
		uma.MainColor = "#7E69CC"
		uma.SubColor = "#CCCBF9"
		return nil
	}

	file, err := os.Open(uma.Image)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	colors := getColors(img)
	if len(colors) == 0 {
		return fmt.Errorf("no colors found in image: %s", uma.Image)
	}

	uma.MainColor = colors[0]
	if len(colors) > 1 {
		uma.SubColor = colors[1]
	}

	return nil
}

// funcion para extraer colores presentes de una imagen
func getColors(img image.Image) []string {
	bounds := img.Bounds()
	freq := make(map[[3]byte]int)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if a == 0 {
				continue
			}
			qr := uint8(r >> 8)
			qg := uint8(g >> 8)
			qb := uint8(b >> 8)
			if qr >= 250 && qg >= 250 && qb >= 250 {
				continue
			}
			freq[[3]byte{qr, qg, qb}]++
		}
	}

	type cc struct {
		color [3]byte
		count int
	}
	var sorted []cc
	for c, n := range freq {
		sorted = append(sorted, cc{c, n})
	}
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && sorted[j].count > sorted[j-1].count; j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
	}

	var result []string
	for _, c := range sorted {
		hex := fmt.Sprintf("#%02X%02X%02X", c.color[0], c.color[1], c.color[2])
		duplicate := false
		for _, existing := range result {
			if colorDistance(c.color, hexToRGB(existing)) < 60 {
				duplicate = true
				break
			}
		}
		if !duplicate {
			result = append(result, hex)
		}
	}
	return result
}

func colorDistance(a, b [3]byte) float64 {
	dr := float64(a[0]) - float64(b[0])
	dg := float64(a[1]) - float64(b[1])
	db := float64(a[2]) - float64(b[2])
	return math.Sqrt(dr*dr + dg*dg + db*db)
}

func hexToRGB(hex string) [3]byte {
	var r, g, b byte
	fmt.Sscanf(hex, "#%02X%02X%02X", &r, &g, &b)
	return [3]byte{r, g, b}
}
