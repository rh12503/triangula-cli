package utils

import (
	"fmt"
	image2 "github.com/RH12503/Triangula/image"
	"image"
	"os"
)

// decodeImage reads and decodes an image.
func decodeImage(imageFile string) (image2.Data, error) {
	file, err := os.Open(imageFile)

	if err != nil {
		fmt.Println("\u001b[31merror reading image file\u001b[0m")
		return image2.Data{}, err
	}

	image, _, err := image.Decode(file)

	file.Close()

	if err != nil {
		fmt.Println("\u001b[31merror decoding image\u001b[0m")
		return image2.Data{}, err
	}

	img := image2.ToData(image)
	return img, nil
}