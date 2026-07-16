package barcode

import (
	"fmt"
	"image"
	"math"
	"os"

	_ "image/png"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
	"golang.org/x/image/draw"
)

const (
	minWidth  = 400
	minHeight = 100
)

/**
 * Read the barcode image file
 */
func Read(path string) (string, error) {

	_, err := os.Stat(path)
	if err != nil {
		return "File not found", err
	}

	file, err := os.Open(path)
	if err != nil {
		return "Failed to open file", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "Failed to decode file", err
	}

	scaleX := float64(minWidth) / float64(img.Bounds().Dx())
	scaleY := float64(minHeight) / float64(img.Bounds().Dy())

	scale := math.Max(scaleX, scaleY)

	if scale > 1 {
		newWidth := int(math.Ceil(float64(img.Bounds().Dx()) * scale))
		newHeight := int(math.Ceil(float64(img.Bounds().Dy()) * scale))

		scaled := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

		draw.NearestNeighbor.Scale(
			scaled,
			scaled.Bounds(),
			img,
			img.Bounds(),
			draw.Over,
			nil,
		)
		img = scaled
	}

	bmpSrc, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "Failed to convert barcode image", err
	}
	reader := oned.NewCode128Reader()

	result, err := reader.Decode(bmpSrc, nil)
	if err != nil {
		return "Failed to decode image", err
	}

	fmt.Println(result.GetText())

	return result.GetText(), nil
}
