package parser

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

type Frontmatter struct {
	Title  string `yaml:"title"`
	Date   string `yaml:"date"`
	Author string `yaml:"author"`
}

type PageData struct {
	Meta    Frontmatter
	Content template.HTML
}

type PostInfo struct {
	Meta Frontmatter
	URL  string
}

type IndexData struct {
	Posts []PostInfo
}

func parseMarkdown(content []byte) (Frontmatter, []byte, error) {
	var meta Frontmatter
	body := content

	if bytes.HasPrefix(content, []byte("---\n")) || bytes.HasPrefix(content, []byte("---\r\n")) {
		parts := bytes.SplitN(content, []byte("---"), 3)

		if len(parts) == 3 {
			err := yaml.Unmarshal(parts[1], &meta)
			if err != nil {
				return meta, nil, err
			}

			body = bytes.TrimSpace(parts[2])
		}
	}

	if meta.Title == "" {
		meta.Title = "Post sem título"
	}
	if meta.Date == "" {
		meta.Date = "1970-01-01"
	}

	return meta, body, nil
}

func MakeIndex(posts []PostInfo, tmpl *template.Template) error {
	data := IndexData{
		Posts: posts,
	}

	outputPath := filepath.Join("public", "index.html")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer outputFile.Close()

	if err := tmpl.Execute(outputFile, data); err != nil {
		return err
	}

	return nil
}

func ProcessFile(filename string, tmpl *template.Template) (PostInfo, error) {
	var info PostInfo

	mdPath := filepath.Join("content", filename)
	mdContent, err := os.ReadFile(mdPath)
	if err != nil {
		return info, err
	}

	meta, mdBody, err := parseMarkdown(mdContent)

	var buf bytes.Buffer
	if err := goldmark.Convert(mdBody, &buf); err != nil {
		return info, err
	}

	data := PageData{
		Meta:    meta,
		Content: template.HTML(buf.String()),
	}

	nameWithoutExtension := strings.TrimSuffix(filename, ".md")
	finalName := nameWithoutExtension + ".html"
	outputPath := filepath.Join("public", finalName)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return info, err
	}

	defer outputFile.Close()

	if err := tmpl.Execute(outputFile, data); err != nil {
		return info, err
	}

	info.Meta = meta
	info.URL = finalName

	return info, nil
}
