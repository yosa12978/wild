package main

import (
	"context"
	"embed"
	_ "embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"regexp"
)

//go:embed draft
var draftFS embed.FS

type DraftCommand struct {
}

func copyDir(dst string) error {
	expr, _ := regexp.Compile(`^([^/]+)`)
	return fs.WalkDir(draftFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		newPath := expr.ReplaceAllString(path, dst)
		if d.IsDir() {
			return os.MkdirAll(newPath, 0755)
		}
		src, err := draftFS.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()
		dstFile, err := os.Create(newPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()
		_, err = io.Copy(dstFile, src)
		return err
	})
}

func (d *DraftCommand) Run(ctx context.Context, args []string) error {
	fmt.Println("creating default site structure")
	return copyDir(args[0])
}

func (d *DraftCommand) Subcommands() map[string]Command {
	return map[string]Command{}
}

func (d *DraftCommand) Use() string {
	return "draft"
}

func (d *DraftCommand) Short() string {
	return "draft creates something"
}

func (d *DraftCommand) Long() string {
	return "draft creates something long"
}
