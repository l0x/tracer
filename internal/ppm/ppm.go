package ppm

import (
	"os"
	"fmt"
	"penrodyn.com/tracer/internal/img"
)

func WritePPM(img *img.Img, fn string) error {
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the header
	fmt.Fprintf(f, "P6\n%d %d\n%d\n", img.Width, img.Height, 255)
	// Draw the rest of the Owl
	_, err = f.Write(img.Pixels)
	return err
}

