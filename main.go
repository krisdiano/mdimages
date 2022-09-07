package main

import (
	"log"
	"os"

	"github.com/krisdiano/mdimages/internal"
	"github.com/urfave/cli/v2"
)

func ExtractAction(app *cli.App) {
	cmd := &cli.Command{
		Name:  "extract",
		Usage: "Extract the path of the image",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Required: true,
			},
		},
		Action: internal.Extract,
	}
	app.Commands = append(app.Commands, cmd)
}

func UploadAction(app *cli.App) {
	cmd := &cli.Command{
		Name:  "upload",
		Usage: "Upload local image to github repo",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "repo",
				Value: "/Users/litianxiang/Desktop/workspace/repo/mdimages",
			},
			&cli.StringFlag{
				Name:  "dir",
				Value: "default",
			},
			&cli.StringSliceFlag{
				Name:     "paths",
				Required: true,
			},
		},
		Action: internal.Upload,
	}
	app.Commands = append(app.Commands, cmd)
}

func RewriteAction(app *cli.App) {
	cmd := &cli.Command{
		Name:  "rewrite",
		Usage: "Convert local path to github url",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "i",
				Value: false,
			},
			&cli.StringFlag{
				Name:     "path",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:     "paths",
				Required: true,
			},
		},
		Action: internal.Rewrite,
	}
	app.Commands = append(app.Commands, cmd)
}

func main() {
	app := &cli.App{
		Name:  "mdit",
		Usage: "Markdown images transfer.",
	}

	ExtractAction(app)
	UploadAction(app)
	RewriteAction(app)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
