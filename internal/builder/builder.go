package builder

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"

	"github.com/omurilo/papiro/internal/parser"
	"github.com/omurilo/papiro/internal/tmpl"
)

func LoadTemplates() (*template.Template, *template.Template, error) {
	var postTmpl, indexTmpl *template.Template
	var err error

	if _, errStat := os.Stat("templates/template.html"); errStat == nil {
		fmt.Println("Usando templates customizados da pasta /templates...")

		postTmpl, err = template.ParseFiles("templates/post_template.html")
		if err != nil {
			return nil, nil, fmt.Errorf("erro no post_template.html local: %v", err)
		}

		indexTmpl, err = template.ParseFiles("templates/index_template.html")
		if err != nil {
			return nil, nil, fmt.Errorf("erro no index_templat.html local: %v", err)
		}
	} else {
		fmt.Println("Usando tema padrão embutido no Papiro...")

		postTmpl, err = template.ParseFS(tmpl.Files, "post_template.html")
		if err != nil {
			return nil, nil, fmt.Errorf("erro no template embutido: %v", err)
		}

		indexTmpl, err = template.ParseFS(tmpl.Files, "index_template.html")
		if err != nil {
			return nil, nil, fmt.Errorf("erro no template embutido: %v", err)
		}
	}

	return postTmpl, indexTmpl, nil
}

func BuildSite() error {
	os.MkdirAll("public", 0755)

	postTmpl, indexTmpl, err := LoadTemplates()
	if err != nil {
		return err
	}

	fmt.Println("Templates carregados com sucesso!")

	files, err := os.ReadDir("content")
	if err != nil {
		return err
	}

	var allPosts []parser.PostInfo

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {
			info, err := parser.ProcessFile(file.Name(), postTmpl)
			if err != nil {
				return err
			}

			allPosts = append(allPosts, info)
		}
	}

	sort.Slice(allPosts, func(i, j int) bool {
		return allPosts[i].Meta.Date > allPosts[j].Meta.Date
	})

	parser.MakeIndex(allPosts, indexTmpl)

	return nil
}
