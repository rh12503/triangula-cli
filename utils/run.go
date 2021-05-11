package utils

import (
	"encoding/json"
	"fmt"
	"github.com/RH12503/Triangula/fitness"
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

var printReps = 20

// RunAlgorithm runs an algorithm, saves the output, and prints statistics.
func RunAlgorithm(imageFile, outputFile string, numPoints uint, mutations uint,
	variation float64, population, cache, cutoff, block, repetitions, threads uint) {

	runtime.GOMAXPROCS(int(threads))

	reps := int(repetitions)

	pointFactory := func() normgeom.NormPointGroup {
		return (generator.RandomGenerator{}).Generate(int(numPoints))
	}
	color.Yellow("Reading image file...")

	img, err := decodeImage(imageFile)

	if err != nil {
		return
	}

	color.Yellow("Initializing algorithm...")

	evaluatorFactory := func(n int) evaluator.Evaluator {
		return evaluator.NewParallel(fitness.TrianglesImageFunctions(img,  int(block), n), int(cache))
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
			for j := 0; j < printReps && i < reps; j++ {
				algo.Step()
				i++
			}
			stats := algo.Stats()

			delta := float64(time.Since(ti).Microseconds()) / (float64(printReps) * 1000.)
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
