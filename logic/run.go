package logic

import (
	"encoding/json"
	"fmt"
	"github.com/RH12503/Triangula-CLI/util"
	"github.com/RH12503/Triangula/fitness"
	"github.com/RH12503/Triangula/image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"runtime"
	"strings"
	"time"

	"github.com/RH12503/Triangula/algorithm"
	"github.com/RH12503/Triangula/algorithm/evaluator"
	"github.com/RH12503/Triangula/generator"
	"github.com/RH12503/Triangula/mutation"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/fatih/color"
)

// RunAlgorithm runs an algorithm, saves the output, and prints statistics.
func RunAlgorithm(imageFile, outputFile string, numPoints uint, shape string, mutations uint,
	variation float64, population, cache, cutoff, block, repetitions, threads uint) {

	runtime.GOMAXPROCS(int(threads))

	reps := int(repetitions)

	pointFactory := func() normgeom.NormPointGroup {
		return (generator.RandomGenerator{}).Generate(int(numPoints))
	}
	color.Yellow("Reading image file...")

	img, err := util.DecodeImage(imageFile)

	if err != nil {
		return
	}

	color.Yellow("Initializing algorithm...")

	var fitnessFuncsFactory func(target image.Data, blockSize, n int) []fitness.CacheFunction

	switch shape {
	case "triangles":
		fitnessFuncsFactory = fitness.TrianglesImageFunctions
		break
	case "polygons":
		fitnessFuncsFactory = fitness.PolygonsImageFunctions
		break
	default:
		color.Red("invalid shape type")
		return
	}

	evaluatorFactory := func(n int) evaluator.Evaluator {
		return evaluator.NewParallel(fitnessFuncsFactory(img, int(block), n), int(cache))
	}
	var mutator mutation.Method

	mutator = mutation.NewGaussianMethod(float64(mutations)/float64(numPoints), variation)

	algo := algorithm.NewModifiedGenetic(pointFactory, int(population), int(cutoff), evaluatorFactory, mutator)

	color.Yellow("Running algorithm...")
	filename := outputFile

	if !strings.HasSuffix(filename, ".json") {
		filename += ".json"
	}

	generateOutput(algo, filename, reps)
}

// generateOutput is an utility function for running an algorithm.
func generateOutput(algo algorithm.Algorithm, output string, reps int) error {
	for {

		for i := 0; i < reps; {
			ti := time.Now()
			p := i
			for time.Since(ti).Milliseconds() < 500 && i < reps {
				algo.Step()
				i++
			}
			stats := algo.Stats()

			delta := float64(time.Since(ti).Microseconds()) / (float64(i-p) * 1000.)
			fmt.Fprintf(color.Output, "Generation %v"+color.WhiteString(" | ")+color.YellowString("Fitness")+" %.8f"+color.WhiteString(" | ")+color.YellowString("Time")+" %.2fms\r",
				stats.Generation, stats.BestFitness, delta)
		}

		jsonOut, err := json.Marshal(algo.Best())
		if err != nil {
			color.Red("error encoding json")
			return err
		}
		err = ioutil.WriteFile(output, jsonOut, 0644)
		if err != nil {
			color.Red("error writing json output")
			return err
		}
	}
}
