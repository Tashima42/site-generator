package main

import (
	"log"
	"os"
	"time"

	"github.com/tashima42/site-generator/generator"
	"github.com/tashima42/site-generator/manager"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	date    = "unknown"
)

func main() {
	app := &cli.App{
		Name:                   "site-generator",
		UseShortOptionHandling: true,
		Version:                version,
		Suggest:                true,
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
						Name:     "destination",
						Usage:    "Destination path for the built files",
						Aliases:  []string{"d"},
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					return generator.Generate(generator.Options{
						Name:            ctx.String("title"),
						DestinationPath: ctx.String("destination"),
					})
				},
			},
			{
				Name:  "manager",
				Usage: "Manage the website pages",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "Add a page",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "pages-path",
								Usage:    "Folder with the current pages",
								Aliases:  []string{"p"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "folder-name",
								Usage:    "Name of the folder to create the pages at",
								Aliases:  []string{"f"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "title",
								Usage:    "Title of the created page",
								Aliases:  []string{"t"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "html-path",
								Usage:    "Path of the HTML page",
								Aliases:  []string{"m"},
								Required: true,
							},
						},
						Action: func(ctx *cli.Context) error {
							return manager.AddPage(manager.Options{
								PagesPath:  ctx.String("pages-path"),
								FolderName: ctx.String("folder-name"),
								Title:      ctx.String("title"),
								HTMLPath:   ctx.String("html-path"),
							})
						},
					},
				},
			},
		},
	}
	compiled, err := time.Parse("2006-01-02T15:04:05Z", date)
	if err == nil {
		app.Compiled = compiled
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
