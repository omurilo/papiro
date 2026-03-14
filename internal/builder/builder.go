package builder

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/omurilo/papiro/internal/parser"
	"github.com/omurilo/papiro/internal/tmpl"
)

func loadThemeTemplates() (*template.Template, *template.Template, error) {
	var postTmpl, indexTmpl *template.Template
	var err error

	if _, errStat := os.Stat("theme/post_template.html"); errStat == nil {
		fmt.Println("Usando templates customizados da pasta /theme...")

		postTmpl, err = template.ParseFiles("theme/post_template.html")
		if err != nil {
			return nil, nil, fmt.Errorf("erro no post_template.html local: %v", err)
		}

		indexTmpl, err = template.ParseFiles("theme/index_template.html")
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

	if _, err := os.Stat("static"); !os.IsNotExist(err) {
		if err := copyDir("static", "public"); err != nil {
			return fmt.Errorf("erro ao copiar diretório static: %v", err)
		}
		fmt.Println("Diretório estático copiado com sucesso!")
	} else {
		if err := copyDirEmbedded(tmpl.Files, "public"); err != nil {
			return err
		}
	}

	postTmpl, indexTmpl, err := loadThemeTemplates()
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

func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(src, path)
		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, data, info.Mode())
	})
}

func copyDirEmbedded(src embed.FS, dst string) error {
	return fs.WalkDir(src, "static", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, path), 0755)
		}
		data, err := fs.ReadFile(src, path)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(dst, path), data, 0644)
	})
}
