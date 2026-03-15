package builder

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	textTmpl "text/template"

	"github.com/omurilo/papiro/internal/config"
	"github.com/omurilo/papiro/internal/parser"
	"github.com/omurilo/papiro/internal/tmpl"
)

func loadThemeTemplates() (*template.Template, *template.Template, *textTmpl.Template, error) {
	var postTmpl, indexTmpl *template.Template
	var feedTmpl *textTmpl.Template
	var err error

	tmplFuncs := template.FuncMap{
		"date": func(layout string, date parser.YamlDate) string {
			return date.Format(layout)
		},
	}

	if _, errStat := os.Stat("theme/post_template.html"); errStat == nil {
		fmt.Println("Usando templates customizados da pasta /theme...")

		postTmpl, err = template.New("post_template.html").Funcs(tmplFuncs).ParseFiles("theme/post_template.html")
		if err != nil {
			return nil, nil, nil, fmt.Errorf("erro no post_template.html local: %v", err)
		}

		indexTmpl, err = template.New("index_template.html").Funcs(tmplFuncs).ParseFiles("theme/index_template.html")
		if err != nil {
			return nil, nil, nil, fmt.Errorf("erro no index_templat.html local: %v", err)
		}

		feedTmpl, err = textTmpl.New("feed.rss").Funcs(tmplFuncs).ParseFiles("theme/feed.rss")
		if err != nil {
			return nil, nil, nil, fmt.Errorf("erro no feed_template.html local: %v", err)
		}
	} else {
		fmt.Println("Usando tema padrão embutido no Papiro...")

		postTmpl, err = template.New("post_template.html").Funcs(tmplFuncs).ParseFS(tmpl.Files, "post_template.html")
		if err != nil {
			return nil, nil, nil, fmt.Errorf("erro no template embutido: %v", err)
		}

		indexTmpl, err = template.New("index_template.html").Funcs(tmplFuncs).ParseFS(tmpl.Files, "index_template.html")
		if err != nil {
			return nil, nil, nil, fmt.Errorf("erro no template embutido: %v", err)
		}

		feedTmpl, err = textTmpl.New("feed.rss").Funcs(tmplFuncs).ParseFS(tmpl.Files, "feed.rss")
		if err != nil {
			return nil, nil, nil, fmt.Errorf("erro no template embutido: %v", err)
		}
	}

	return postTmpl, indexTmpl, feedTmpl, nil
}

func BuildSite() error {
	cfg, err := config.LoadConfig("papiro.yaml")
	if err != nil {
		return fmt.Errorf("erro ao carregar configuração: %v", err)
	}

	files, err := os.ReadDir("content")
	if err != nil {
		return fmt.Errorf("erro ao tentar ler o diretório de conteúdo: \n\t%v", err)
	}

	os.MkdirAll("public", 0755)

	if _, err := os.Stat("theme/static"); !os.IsNotExist(err) {
		if err := copyDir("theme/static", "public/static"); err != nil {
			return fmt.Errorf("erro ao copiar diretório static: %v", err)
		}
		fmt.Println("Diretório estático copiado com sucesso!")
	} else {
		if err := copyDirEmbedded(tmpl.Files, "static", "public"); err != nil {
			return err
		}
	}

	postTmpl, indexTmpl, feedTmpl, err := loadThemeTemplates()
	if err != nil {
		return err
	}

	fmt.Println("Templates carregados com sucesso!")

	var allPosts []parser.PostInfo

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {
			info, err := parser.ProcessFile(file.Name(), postTmpl)
			if err != nil {
				fmt.Printf("Aviso: erro ao processar %s: %v\n", file.Name(), err)
				continue
			}

			if info.Meta.Draft {
				fmt.Printf("Rascunho gerado (oculto da Home): %s\n", info.URL)
				continue
			}

			allPosts = append(allPosts, info)
		}
	}

	sort.Slice(allPosts, func(i, j int) bool {
		return allPosts[i].Meta.Date.Time.Before(allPosts[j].Meta.Date.Time)
	})

	parser.MakeIndex(allPosts, indexTmpl)

	parser.MakeRSS(allPosts, feedTmpl, cfg)

	return nil
}

func InitSite() error {
	os.MkdirAll("content", 0755)
	os.MkdirAll("theme/static", 0755)

	if err := copyDirEmbedded(tmpl.Files, "static", "theme"); err != nil {
		return err
	}

	if err := copyDirEmbedded(tmpl.Files, "post_template.html", "theme"); err != nil {
		return err
	}

	if err := copyDirEmbedded(tmpl.Files, "index_template.html", "theme"); err != nil {
		return err
	}

	if err := copyDirEmbedded(tmpl.Files, "papiro.yaml", "."); err != nil {
		return err
	}

	os.WriteFile("content/hello-world.md", []byte("---\ntitle: \"Olá, Mundo! Bem-vindo ao Papiro.\"\ndate: 2026-03-14\nauthor: \"Murilo\"\n---\n\n# O Início de uma Nova Jornada\n\nSe você está lendo isso, significa que o comando `init` funcionou perfeitamente e o motor do **Papiro** já está rodando! \n\nEste é um gerador de sites estáticos focado em simplicidade, velocidade e na beleza da escrita em texto puro. Tudo o que você precisa fazer é escrever em Markdown e deixar que o Go faça o resto.\n\n## O que você pode fazer aqui?\n\nComo usamos o padrão Markdown, você tem total liberdade para formatar seus textos de forma rápida:\n\n* Criar **textos em negrito** para dar ênfase.\n* Usar *itálico* para pensamentos ou termos estrangeiros.\n* Fazer listas organizadas, como esta.\n\nSe precisar citar alguém importante, o design clássico cuida disso:\n\n> \"A simplicidade é o último grau de sofisticação.\" \n> — Leonardo da Vinci\n\n### Suporte a Código\n\nE como todo bom desenvolvedor, você pode compartilhar seus trechos de código facilmente. O Papiro já deixa tudo bem formatado:\n\n```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"O Papiro é rápido demais!\")\n}\n```"), 0644)
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

func copyDirEmbedded(src embed.FS, srcPath, dst string) error {
	return fs.WalkDir(src, srcPath, func(path string, d fs.DirEntry, err error) error {
		fmt.Println(path)
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
