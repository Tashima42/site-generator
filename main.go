package main

import (
	"log"
	"os"

	"github.com/tashima42/site-generator/generator"
	"github.com/urfave/cli/v2"
)

const version = "0.0.1"

func main() {
	app := &cli.App{
		Name:                   "site-generator",
		UseShortOptionHandling: true,
		Version:                version,
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate static websites using templates",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "title",
						Usage:    "Website title",
						Aliases:  []string{"t"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "source",
						Usage:    "Path for the source files",
						Aliases:  []string{"s"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "destination",
						Usage:    "Destination path for the built files",
						Aliases:  []string{"d"},
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					return generator.Generate(generator.Options{
						Name:            ctx.String("title"),
						SourcePath:      ctx.String("source"),
						DestinationPath: ctx.String("destination"),
					})
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
