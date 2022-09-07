package internal

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/russross/blackfriday/v2"
	"github.com/urfave/cli/v2"
)

func Extract(ctx *cli.Context) error {
	path := ctx.String("path")
	if len(path) < 0 {
		return errors.New("markdown path is required")
	}

	bins, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	md := blackfriday.New()
	root := md.Parse(bins)

	paths := make(map[string]struct{})
	fn := func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if node.Type == blackfriday.Image {
			text := string(node.LinkData.Destination)
			if !strings.HasPrefix(text, "https://github.com") {
				paths[text] = struct{}{}
			}
		}
		return blackfriday.GoToNext
	}
	root.Walk(fn)

	for k := range paths {
		fmt.Println(k)
	}
	return nil
}
