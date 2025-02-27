package main

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v3"
)

type GenerateCommand struct {
	fileMetadata map[string]any
}

func NewGenerateCmd() Command {
	return &GenerateCommand{}
}

func buildProject(src, dst string, vars map[string]any) error {
	expr, _ := regexp.Compile(`^([^\/]+)`)
	markdown := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithExtensions(meta.New(meta.WithStoresInDocument())),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		newPath := expr.ReplaceAllString(path, dst)

		if d.IsDir() {
			return os.MkdirAll(newPath, 0755)
		}

		// if it's not a markdown file then just copy it
		if !strings.HasSuffix(newPath, ".md") {
			srcFile, err := draftFS.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()
			dstFile, err := os.Create(newPath)
			if err != nil {
				return err
			}
			defer dstFile.Close()
			_, err = io.Copy(dstFile, srcFile)
			return err
		}
		newPath = strings.TrimSuffix(newPath, ".md") + ".html"

		// markdown template file
		mdFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer mdFile.Close()

		// file where we are going to store html
		htmlFile, err := os.Create(newPath)
		if err != nil {
			return err
		}
		defer htmlFile.Close()

		mdFileContent, err := io.ReadAll(mdFile)
		document := markdown.Parser().Parse(text.NewReader(mdFileContent))

		//gathering metadata for current file
		metadata := document.OwnerDocument().Meta()

		var renderedMarkdown bytes.Buffer
		if err := markdown.Renderer().Render(
			&renderedMarkdown,
			mdFileContent,
			document,
		); err != nil {
			return err
		}

		var withLayout bytes.Buffer
		templ := template.Must(template.ParseFiles("layout.html"))
		// should also append wild.yml into this template
		templ.Execute(&withLayout, map[string]any{"Main": template.HTML(renderedMarkdown.String())})

		vars["Self"] = metadata
		t, _ := template.New("finalPage").Parse(withLayout.String())
		return t.Execute(htmlFile, vars)
	})
}

// make some gather metadata method
// invoke it in the beginning
// store all of it in hierarchical structure inside map
// before rendering template append also "Self" field

func (g *GenerateCommand) gatherMetadata(src string) error {
	markdown := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithExtensions(meta.New(meta.WithStoresInDocument())),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
	filesMetadata := map[string]any{}
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// if it's not a markdown file then just copy it
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		//parts := strings.Split(path, "/")

		// markdown template file
		mdFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer mdFile.Close()

		mdFileContent, err := io.ReadAll(mdFile)
		document := markdown.Parser().Parse(text.NewReader(mdFileContent))

		metadata := document.OwnerDocument().Meta()

		filesMetadata[path] = metadata
		return nil
	})
}

func (g *GenerateCommand) Run(ctx context.Context, args []string) error {
	wildFile, err := os.Open("wild.yml")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("Couldn't find wild.yml file. Is that a wild project?")
		}
		return err
	}
	defer wildFile.Close()
	var vars map[string]any
	if err := yaml.NewDecoder(wildFile).Decode(&vars); err != nil {
		return err
	}
	return buildProject("content", "build", vars)
}

func (g *GenerateCommand) Long() string {
	return "this command generates a static website"
}
func (g *GenerateCommand) Short() string {
	return "generate static website"
}

func (g *GenerateCommand) Subcommands() map[string]Command {
	return nil
}

func (g *GenerateCommand) Use() string {
	return "generate"
}
