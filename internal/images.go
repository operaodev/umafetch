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
	"sort"
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

// downloadImage downloads the image of a UMA and saves it to the image directory.
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

// extractColors extracts the primary and secondary colors from a UMA image.
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

	main, sub, ok := dominantColors(img)
	if !ok {
		return fmt.Errorf("no colors found in image: %s", uma.Image)
	}

	uma.MainColor = toHex(boostSaturation(boostBrightness(main)))
	uma.SubColor = toHex(boostSaturation(boostBrightness(sub)))
	return nil
}

// dominantColors gets the 4 distinct dominant colors from the image,
// then selects main (most vivid) and sub (most distant from main).
func dominantColors(img image.Image) (main, sub [3]byte, ok bool) {
	bounds := img.Bounds()
	freq := make(map[[3]byte]int)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if a < 0x8000 {
				continue
			}
			// Coarse quantization in blocks of 8 (5 bits per channel)
			qr := uint8((r >> 8) & 0xF8)
			qg := uint8((g >> 8) & 0xF8)
			qb := uint8((b >> 8) & 0xF8)
			// Discard near-blacks
			bright := int(qr) + int(qg) + int(qb)
			if bright < 60 {
				continue
			}
			freq[[3]byte{qr, qg, qb}]++
		}
	}

	if len(freq) == 0 {
		return
	}

	type bucket struct {
		color [3]byte
		count int
	}
	buckets := make([]bucket, 0, len(freq))
	for c, n := range freq {
		buckets = append(buckets, bucket{c, n})
	}
	// Sort by descending frequency
	sort.Slice(buckets, func(i, j int) bool {
		return buckets[i].count > buckets[j].count
	})

	// Classify buckets: darks and skin/brown are deprioritized as "background"
	const minDist = 60.0

	isDark := func(c [3]byte) bool {
		return int(c[0])+int(c[1])+int(c[2]) < 180
	}

	isWhite := func(c [3]byte) bool {
		return int(c[0])+int(c[1])+int(c[2]) > 700
	}

	// Skin and brown tones: R dominant, green > blue (excludes pure reds),
	// moderate saturation (excludes very vivid outfit colors).
	isSkinOrBrown := func(c [3]byte) bool {
		r, g, b := float64(c[0]), float64(c[1]), float64(c[2])
		if r < g || r < b || g <= b || r-b < 30 {
			return false
		}
		minC := math.Min(g, b)
		delta := r - minC
		sat := delta / r
		if sat < 0.10 || sat > 0.80 {
			return false
		}
		h := (g - b) / delta * 60
		return h >= 0 && h <= 50
	}

	var bgBuckets, colorBuckets, whiteBuckets []bucket
	var darkTotal, colorTotal, whiteTotal int
	for _, b := range buckets {
		switch {
		case isDark(b.color):
			darkTotal += b.count
			bgBuckets = append(bgBuckets, b)
		case isWhite(b.color):
			whiteTotal += b.count
			whiteBuckets = append(whiteBuckets, b)
		case isSkinOrBrown(b.color):
			bgBuckets = append(bgBuckets, b)
		default:
			colorTotal += b.count
			colorBuckets = append(colorBuckets, b)
		}
	}

	// Extremely dominant black: exceeds 3× the total of colorful pixels
	extremelyDark := colorTotal > 0 && float64(darkTotal) > 3.0*float64(colorTotal)
	// Very dominant white: more white pixels than colorful ones
	whiteDominant := whiteTotal > colorTotal

	clusters := make([][3]byte, 0, 4)

	addClusters := func(src []bucket) {
		for _, b := range src {
			if len(clusters) == 4 {
				break
			}
			tooClose := false
			for _, c := range clusters {
				if colorDistance(b.color, c) < minDist {
					tooClose = true
					break
				}
			}
			if !tooClose {
				clusters = append(clusters, b.color)
			}
		}
	}

	// First pass: only representative colors (not darks, not skin/brown, not whites)
	addClusters(colorBuckets)

	// Include whites if they are very dominant
	if whiteDominant {
		addClusters(whiteBuckets)
	}

	// Include background (darks and skin/brown) only if there is little diversity or black dominates
	if extremelyDark || len(clusters) < 2 {
		addClusters(bgBuckets)
	}

	if len(clusters) == 0 {
		return
	}

	// Main: highest vividness (saturation × normalized brightness)
	mainIdx := 0
	bestVivid := -1.0
	for i, c := range clusters {
		s := float64(saturation(c))
		br := (float64(c[0]) + float64(c[1]) + float64(c[2])) / 3.0
		vivid := s * (br / 255.0)
		if vivid > bestVivid {
			bestVivid = vivid
			mainIdx = i
		}
	}
	main = clusters[mainIdx]

	// Sub: most distant from main among the remaining
	subIdx := -1
	bestDist := -1.0
	for i, c := range clusters {
		if i == mainIdx {
			continue
		}
		d := colorDistance(c, main)
		if d > bestDist {
			bestDist = d
			subIdx = i
		}
	}

	if subIdx < 0 {
		// Only one cluster; use main slightly desaturated as sub
		sub = [3]byte{
			uint8(math.Min(255, float64(main[0])*0.7)),
			uint8(math.Min(255, float64(main[1])*0.7)),
			uint8(math.Min(255, float64(main[2])*0.7)),
		}
	} else {
		sub = clusters[subIdx]
	}

	ok = true
	return
}

// boostBrightness increases brightness if the color is too dark to display in the terminal.
func boostBrightness(c [3]byte) [3]byte {
	const minBright = 100
	bright := (int(c[0]) + int(c[1]) + int(c[2])) / 3
	if bright >= minBright {
		return c
	}
	factor := float64(minBright) / math.Max(float64(bright), 1)
	clamp := func(v float64) uint8 {
		if v > 255 {
			return 255
		}
		return uint8(v)
	}
	return [3]byte{
		clamp(float64(c[0]) * factor),
		clamp(float64(c[1]) * factor),
		clamp(float64(c[2]) * factor),
	}
}

// boostSaturation increases the color saturation in HSL space while preserving hue and lightness.
func boostSaturation(c [3]byte) [3]byte {
	r, g, b := float64(c[0])/255, float64(c[1])/255, float64(c[2])/255
	maxC := math.Max(r, math.Max(g, b))
	minC := math.Min(r, math.Min(g, b))
	delta := maxC - minC
	l := (maxC + minC) / 2

	if delta == 0 {
		return c // pure gray, no hue to saturate
	}

	// Calculate hue
	var h float64
	switch maxC {
	case r:
		h = math.Mod((g-b)/delta, 6)
	case g:
		h = (b-r)/delta + 2
	default:
		h = (r-g)/delta + 4
	}
	h /= 6
	if h < 0 {
		h += 1
	}

	// Maximum saturation
	s := math.Min(delta/(1-math.Abs(2*l-1)), 1.0)

	// HSL to RGB
	c2 := (1 - math.Abs(2*l-1)) * s
	x := c2 * (1 - math.Abs(math.Mod(h*6, 2)-1))
	m := l - c2/2

	var r1, g1, b1 float64
	switch int(h * 6) {
	case 0:
		r1, g1, b1 = c2, x, 0
	case 1:
		r1, g1, b1 = x, c2, 0
	case 2:
		r1, g1, b1 = 0, c2, x
	case 3:
		r1, g1, b1 = 0, x, c2
	case 4:
		r1, g1, b1 = x, 0, c2
	default:
		r1, g1, b1 = c2, 0, x
	}

	clamp := func(v float64) uint8 {
		v = (v + m) * 255
		if v > 255 {
			return 255
		}
		if v < 0 {
			return 0
		}
		return uint8(v)
	}
	return [3]byte{clamp(r1), clamp(g1), clamp(b1)}
}

func toHex(c [3]byte) string {
	return fmt.Sprintf("#%02X%02X%02X", c[0], c[1], c[2])
}

func saturation(c [3]byte) int {
	return int(max3(c[0], c[1], c[2])) - int(min3(c[0], c[1], c[2]))
}

func max3(a, b, c byte) byte {
	if a > b && a > c {
		return a
	}
	if b > c {
		return b
	}
	return c
}

func min3(a, b, c byte) byte {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}

func colorDistance(a, b [3]byte) float64 {
	dr := float64(a[0]) - float64(b[0])
	dg := float64(a[1]) - float64(b[1])
	db := float64(a[2]) - float64(b[2])
	return math.Sqrt(dr*dr + dg*dg + db*db)
}
