package internal

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

type GitCmd struct {
	dir string
	err error
}

func NewGit(dir string) (*GitCmd, error) {
	f, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !f.IsDir() {
		return nil, fmt.Errorf("%s is not a dir", dir)
	}
	return &GitCmd{
		dir: dir,
	}, nil
}

func (g *GitCmd) Error() error {
	return g.err
}

func (g *GitCmd) Do(subcmd string) error {
	if g.err != nil {
		return g.err
	}

	cmd := exec.Command("/bin/sh", "-c", subcmd)
	cmd.Dir = g.dir
	resp, err := cmd.CombinedOutput()
	if err != nil {
		if !strings.Contains(string(resp), "Your branch is up to date with 'origin/master'.") {
			g.err = fmt.Errorf("opt %s failed, err:%v, output:%s", subcmd, err, string(resp))
			return err
		}
	}
	return nil
}

func copy(src, dst string) error {
	sfn, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, sfn, 0644)
}

func Upload(ctx *cli.Context) error {
	repo := ctx.String("repo")
	dir := ctx.String("dir")
	paths := ctx.StringSlice("paths")

	if len(paths) == 0 {
		return nil
	}
	if strings.Contains(dir, "internal") {
		return fmt.Errorf("internal dir must be a source code dir")
	}

	info, err := os.Stat(repo)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("repo %s is not a dir", repo)
	}

	rawPaths := make(map[string]string)
	for i := range paths {
		item, err := url.PathUnescape(paths[i])
		if err != nil {
			return err
		}
		rawPaths[item] = paths[i]
		paths[i] = item
	}

	dst := strings.Trim(dir, "/")
	repo = filepath.Join(repo, dst)
	defer func() {
		if err != nil {
			os.RemoveAll(repo)
		}
	}()

	info, err = os.Stat(repo)
	if os.IsNotExist(err) {
		err = os.MkdirAll(repo, 0755)
		if err != nil {
			return fmt.Errorf("create dir %s failed, err:%v", repo, err)
		}
	} else if err != nil {
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("dst %s is not a dir", repo)
	}

	ret := make(map[string]string)
	for _, path := range paths {
		remoteGitPath := filepath.Join("github.com/krisdiano/mdimages/raw/master", dst, filepath.Base(path))
		remoteGitPath = fmt.Sprintf("https://%s", remoteGitPath)
		localGitPath := filepath.Join(repo, filepath.Base(path))
		err = copy(path, localGitPath)
		if err != nil {
			return err
		}
		ret[path] = remoteGitPath
	}

	g, err := NewGit(dst)
	if err != nil {
		return err
	}

	g.Do("git add .")
	g.Do("git status")
	g.Do(fmt.Sprintf("git commit -m %s-%s", dir, time.Now().Format("20060102")))
	g.Do("git push")
	err = g.Error()
	if err != nil {
		return err
	}

	for k, v := range ret {
		fmt.Println(rawPaths[k])
		fmt.Println(url.PathEscape(v))
	}
	return nil
}
