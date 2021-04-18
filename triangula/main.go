package main

import (
	"log"
	"os"

	"github.com/RH12503/Triangula-CLI/utils"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Triangula",
		Usage: "An iterative algorithm to generate high quality triangulated images.",
		Commands: []*cli.Command{
			{
				Name: "run",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "image",
						Aliases:  []string{"img"},
						Required: true,
						Usage:    "the image file to triangulate",
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"out"},
						Required: true,
						Usage:    "the file to write output to",
					},
					&cli.UintFlag{
						Name:    "points",
						Aliases: []string{"p"},
						Value:   300,
						Usage:   "the number of points to use in the triangulation",
					},
					&cli.UintFlag{
						Name:    "mutations",
						Aliases: []string{"m", "mut"},
						Value:   2,
						Usage:   "the number of mutations to make",
					},
					&cli.Float64Flag{
						Name:    "variation",
						Aliases: []string{"v"},
						Value:   0.3,
						Usage:   "the variation each mutation causes",
					},
					&cli.UintFlag{
						Name:    "population",
						Aliases: []string{"pop", "size"},
						Value:   400,
						Usage:   "the population size in the algorithm",
					},
					&cli.UintFlag{
						Name:    "cache",
						Aliases: []string{"c"},
						Value:   22,
						Usage:   "the cache size as a power of 2",
					},
					&cli.UintFlag{
						Name:    "block",
						Aliases: []string{"b"},
						Value:   5,
						Usage:   "the size of the block used when rendering",
					},
					&cli.UintFlag{
						Name:    "cutoff",
						Aliases: []string{"cut"},
						Value:   5,
						Usage:   "the size of the block used when rendering",
					},
					&cli.UintFlag{
						Name:    "reps",
						Aliases: []string{"r"},
						Value:   500,
						Usage:   "the number of generations before saving to the output file",
					},
					&cli.UintFlag{
						Name:    "threads",
						Aliases: []string{"t"},
						Value:   0,
						Usage:   "the number of threads to use or 0 to use all cores",
					},
				},
				Action: func(c *cli.Context) error {
					utils.RunAlgorithm(c.String("image"), c.String("output"), c.Uint("points"),
						c.Uint("mutations"), c.Float64("variation"), c.Uint("population"), c.Uint("cache"),
						c.Uint("cutoff"), c.Uint("block"), c.Uint("reps"), c.Uint("threads"))
					return nil
				},
			},
			{
				Name:  "render",
				Usage: "renders to an SVG",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Aliases:  []string{"in", "i"},
						Required: true,
						Usage:    "the .json file containing the points to be rendered",
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"out", "o"},
						Required: true,
						Usage:    "the name of the file to be outputted",
					},
					&cli.StringFlag{
						Name:     "image",
						Aliases:  []string{"img"},
						Required: true,
						Usage:    "the image to use",
					},
				},
				Action: func(c *cli.Context) error {
					utils.RenderSVG(c.String("input"), c.String("output"), c.String("image"))
					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name:  "png",
						Usage: "renders to an PNG",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "format",
								Aliases: []string{"f"},
								Value:   "svg",
								Usage:   "either png or svg",
							},
							&cli.StringFlag{
								Name:    "effect",
								Aliases: []string{"e"},
								Value:   "none",
								Usage:   "either none, gradient, or split",
							},
							&cli.Float64Flag{
								Name:    "scale",
								Aliases: []string{"s"},
								Value:   1,
								Usage:   "either none, gradient, or split",
							},
						},
						Action: func(c *cli.Context) error {
							utils.RenderPNG(c.String("input"), c.String("output"), c.String("image"),
								c.String("effect"), c.Float64("scale"))
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
