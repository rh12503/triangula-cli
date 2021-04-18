package utils

import (
	"encoding/json"
	"fmt"
	"github.com/RH12503/Triangula/algorithm"
	"github.com/RH12503/Triangula/algorithm/evaluator"
	"github.com/RH12503/Triangula/generator"
	"github.com/RH12503/Triangula/mutation"
	"github.com/RH12503/Triangula/normgeom"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"runtime"
	"time"
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
	fmt.Println("\u001B[33mReading image file...")

	img, err := decodeImage(imageFile)

	if err != nil {
		return
	}

	fmt.Println("Initializing algorithm...")

	evaluatorFactory := func(n int) evaluator.Evaluator {
		return evaluator.NewParallel(img, int(cache), int(block), n)
	}
	var mutator mutation.Method

	mutator = mutation.NewGaussianMethod(float64(mutations)/float64(numPoints), variation)

	algo := algorithm.NewSimple(pointFactory, int(population), int(cutoff), evaluatorFactory, mutator)

	fmt.Println("Running algorithm...\u001B[0m")
	generateOutput(algo, outputFile, reps)
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
			fmt.Printf("\u001b[33mGeneration\u001b[0m %v\u001b[37m | \u001b[0m\u001b[33mFitness\u001b[0m %.8f\u001b[37m | \u001b[0m\u001b[33mTime\u001b[0m %.2fms\r",
				stats.Generation, stats.BestFitness, delta)

		}

		jsonOut, err := json.Marshal(algo.Best())
		if err != nil {
			fmt.Println("\u001b[31merror encoding json\u001b[0m")
			return err
		}
		err = ioutil.WriteFile(output+".json", jsonOut, 0644)
		if err != nil {
			fmt.Println("\u001b[31merror writing json output\u001b[0m")
			return err
		}
	}
}
