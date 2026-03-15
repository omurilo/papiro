package parser

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	textTmpl "text/template"
	"time"

	"github.com/omurilo/papiro/internal/config"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

type Frontmatter struct {
	Title  string   `yaml:"title"`
	Date   YamlDate `yaml:"date"`
	Author string   `yaml:"author"`
	Draft  bool     `yaml:"draft"`
}

type YamlDate struct {
	time.Time
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
	Posts           []PostInfo
	BlogUrl         string
	FeedTitle       string
	FeedDescription string
	FeedLanguage    string
}

func (d *YamlDate) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		t, err := time.Parse("2006-01-02", value.Value)
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	}
	return fmt.Errorf("invalid date format")
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
	if meta.Date.IsZero() {
		mDate, _ := time.Parse("2006-01-02", "1970-01-01")
		meta.Date = YamlDate{mDate}
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

func MakeRSS(posts []PostInfo, tmpl *textTmpl.Template, cfg *config.Config) error {
	data := IndexData{
		Posts:           posts,
		BlogUrl:         cfg.URL,
		FeedTitle:       cfg.Title,
		FeedDescription: cfg.Description,
		FeedLanguage:    cfg.Language,
	}

	outputPath := filepath.Join("public", "feed.xml")
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
