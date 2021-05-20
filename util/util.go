package util

import (
	"image"
	"math"
	"os"

	image2 "github.com/RH12503/Triangula/image"
	"github.com/fatih/color"
)

const SvgStart = `<?xml version="1.0"?>
<svg width="%v" height="%v"
     xmlns="http://www.w3.org/2000/svg"
     shape-rendering="crispEdges">
`


// decodeImage reads and decodes an image.
func DecodeImage(imageFile string) (image2.Data, error) {
	file, err := os.Open(imageFile)

	if err != nil {
		color.Red("error reading image file")
		return image2.RGBData{}, err
	}

	image, _, err := image.Decode(file)

	file.Close()

	if err != nil {
		color.Red("error decoding image")
		return image2.RGBData{}, err
	}

	img := image2.ToData(image)
	return img, nil
}


func Scale(num float64, d int) int {
	return int(math.Round(num * float64(d)))
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MultAndRound(v int, s float64) int {
	return int(math.Round(float64(v) * s))
}