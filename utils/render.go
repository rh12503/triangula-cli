package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/RH12503/Triangula-CLI/export"
	"github.com/RH12503/Triangula/normgeom"
)

// decodePoints reads and decodes an JSON file containing the data of points.
func decodePoints(inputFile string) (normgeom.NormPointGroup, error) {
	jsonPoints, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Println("\u001b[31merror reading input file\u001B[0m")
		return normgeom.NormPointGroup{}, err
	}
	var points normgeom.NormPointGroup
	err = json.Unmarshal(jsonPoints, &points)
	if err != nil {
		fmt.Println("\u001b[31merror decoding input file\u001B[0m")
		return normgeom.NormPointGroup{}, err
	}
	return points, nil
}

// RenderPNG renders a triangulation to a PNG.
func RenderPNG(inputFile, outputFile, imageFile, effect string, scale float64) {
	fmt.Println("\u001B[33mReading image file...")

	img, err := decodeImage(imageFile)

	if err != nil {
		return
	}

	fmt.Println("Reading input file...")
	points, err := decodePoints(inputFile)

	if err != nil {
		return
	}

	fmt.Println("Generating PNG...\u001b[0m")
	filename := outputFile + ".png"
	switch e := strings.ToLower(effect); e {
	case "none":
		err = export.WritePNG(filename, points, img, scale)
	case "gradient":
		err = export.WriteEffectPNG(filename, points, img, scale, true)
	case "split":
		err = export.WriteEffectPNG(filename, points, img, scale, false)
	default:
		fmt.Println("\u001b[31munknown effect")
		return
	}

	if err != nil {
		log.Fatal(err)
		fmt.Println("\u001b[31merror generating PNG\u001b[0m")
		return
	}

	fmt.Println("\u001b[32mSuccessfully generated PNG at " + filename + "!\u001B[0m")
}

// RenderSVG renders a triangulation to a SVG.
func RenderSVG(inputFile, outputFile, imageFile string) {
	fmt.Println("\u001B[33mReading image file...")

	img, err := decodeImage(imageFile)

	if err != nil {
		return
	}

	fmt.Println("Reading input file...")
	points, err := decodePoints(inputFile)

	if err != nil {
		return
	}

	fmt.Println("Generating SVG...\u001b[0m")
	err = export.WriteSVG(outputFile+".svg", points, img)

	if err != nil {
		fmt.Println("\u001b[31merror generating SVG\u001b[0m")
		return
	}
	fmt.Println("\u001b[32mSuccessfully generated SVG at " + outputFile + ".svg!\u001B[0m")
}
