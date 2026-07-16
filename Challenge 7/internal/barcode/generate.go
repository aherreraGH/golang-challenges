package barcode

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
)

// var barcodeIDs = []string{
// 	"LIB-GO-001-A7",
// 	"LIB-SYS-204-K9",
// 	"LIB-MATH-88-Z3",
// }

func Create(barcodeIDs string, outputDir string) error {
	const (
		width  = 600
		height = 120
	)

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}

	for _, id := range strings.Split(barcodeIDs, ",") {
		if id != "" {
			if err := generateBarcode(outputDir, id, width, height); err != nil {
				return err
			}

			fmt.Printf("Generated %s.png\n", id)
		}
	}
	return nil
}

func generateBarcode(
	outputDir string,
	value string,
	width int,
	height int,
) error {
	code, err := code128.Encode(value)
	if err != nil {
		return fmt.Errorf("encode Code 128 value: %w", err)
	}

	scaled, err := barcode.Scale(code, width, height)
	if err != nil {
		return fmt.Errorf("scale barcode: %w", err)
	}

	path := filepath.Join(outputDir, value+".png")

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, scaled); err != nil {
		return fmt.Errorf("encode PNG: %w", err)
	}

	return nil
}
