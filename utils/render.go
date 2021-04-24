package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/RH12503/Triangula-CLI/export"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/fatih/color"
)

// decodePoints reads and decodes an JSON file containing the data of points.
func decodePoints(inputFile string) (normgeom.NormPointGroup, error) {
	jsonPoints, err := ioutil.ReadFile(inputFile)
	if err != nil {
		color.Red("error reading input file")
		return normgeom.NormPointGroup{}, err
	}
	var points normgeom.NormPointGroup
	err = json.Unmarshal(jsonPoints, &points)
	if err != nil {
		color.Red("error decoding input file")
		return normgeom.NormPointGroup{}, err
	}
	return points, nil
}

// RenderPNG renders a triangulation to a PNG.
func RenderPNG(inputFile, outputFile, imageFile, effect string, scale float64) {
	color.Yellow("Reading image file...")

	img, err := decodeImage(imageFile)

	if err != nil {
		return
	}

	color.Yellow("Reading input file...")
	points, err := decodePoints(inputFile)

	if err != nil {
		return
	}

	color.Yellow("Generating PNG...")
	filename := outputFile + ".png"
	switch e := strings.ToLower(effect); e {
	case "none":
		err = export.WritePNG(filename, points, img, scale)
	case "gradient":
		err = export.WriteEffectPNG(filename, points, img, scale, true)
	case "split":
		err = export.WriteEffectPNG(filename, points, img, scale, false)
	default:
		color.Red("unknown effect")
		return
	}

	if err != nil {
		log.Fatal(err)
		color.Red("error generating PNG")
		return
	}

	color.Green("Successfully generated PNG at %s!", filename)
}

// RenderSVG renders a triangulation to a SVG.
func RenderSVG(inputFile, outputFile, imageFile string) {
	color.Yellow("Reading image file...")

	img, err := decodeImage(imageFile)

	if err != nil {
		return
	}

	color.Yellow("Reading input file...")
	points, err := decodePoints(inputFile)

	if err != nil {
		return
	}

	color.Yellow("Generating SVG...")
	err = export.WriteSVG(outputFile+".svg", points, img)

	if err != nil {
		color.Red("error generating SVG")
		return
	}
	color.Green("Successfully generated SVG at %s.svg!", outputFile)
}
