package utils

import (
	"image"
	"os"

	image2 "github.com/RH12503/Triangula/image"
	"github.com/fatih/color"
)

// decodeImage reads and decodes an image.
func decodeImage(imageFile string) (image2.Data, error) {
	file, err := os.Open(imageFile)

	if err != nil {
		color.Red("error reading image file")
		return image2.Data{}, err
	}

	image, _, err := image.Decode(file)

	file.Close()

	if err != nil {
		color.Red("error decoding image")
		return image2.Data{}, err
	}

	img := image2.ToData(image)
	return img, nil
}
