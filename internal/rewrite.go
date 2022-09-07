package internal

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func Rewrite(ctx *cli.Context) error {
	path := ctx.String("path")
	if len(path) < 0 {
		return errors.New("markdown path is required")
	}

	kv := make(map[string]string)
	paths := ctx.StringSlice("paths")
	if len(paths) <= 0 || len(paths)%2 != 0 {
		return fmt.Errorf("paths num %d", len(paths))
	}
	for i := 0; i < len(paths); i = i + 2 {
		v, err := url.PathUnescape(paths[i+1])
		if err != nil {
			return err
		}
		kv[paths[i]] = v
	}

	bins, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	content := string(bins)
	for k, v := range kv {
		content = strings.ReplaceAll(content, k, v)
	}
	if !ctx.Bool("i") {
		fmt.Println(content)
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}
