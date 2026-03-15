# 📜 Papiro

> Um gerador de sites estáticos absurdamente rápido, minimalista e feito em Go.

O **Papiro** transforma seus textos em Markdown em um blog ou site estático completo em milissegundos. Sem dependências complexas, sem banco de dados, apenas um único binário e a beleza do texto puro.

## Funcionalidades

- **Velocidade Extrema:** Compila dezenas de páginas em milissegundos.
- **Markdown & Frontmatter:** Escreva usando a sintaxe padrão de mercado com metadados em YAML.
- **Tema Embutido:** Comece a escrever imediatamente com o design clássico embutido no binário, ou crie seu próprio tema na pasta `theme/`.
- **Rascunhos (Drafts):** Oculte publicações em andamento usando `draft: true`.
- **RSS Feed Automático:** Geração de `feed.xml` nativa para seus leitores acompanharem as novidades.
- **Configuração Centralizada:** Personalize o título, descrição e URL do seu site facilmente via `papiro.yaml`.

---

## Instalação

Certifique-se de ter o [Go](https://go.dev/) instalado na sua máquina e rode:

```bash
go install [github.com/omurilo/papiro@latest](https://github.com/omurilo/papiro@latest)
```
(O executável papiro estará disponível no seu terminal).

## Como Usar
1. Inicie um novo projeto
Crie uma pasta para o seu blog e rode o comando de inicialização:

```bash
mkdir meu-blog
cd meu-blog
papiro init
```
Isso criará a estrutura básica, incluindo o arquivo de configuração e um post de boas-vindas.


2. Escreva
Abra a pasta content/ e crie seus arquivos .md. O Frontmatter padrão suporta:

```yaml
---
title: "Meu novo post"
date: "2026-03-14"
author: "Seu Nome"
draft: false
---
# Hello World
```

3. Construa o Site
Gere os arquivos HTML finais rodando:

```bash
papiro build
```
Seu site estará pronto dentro da pasta public/! Basta abrir o index.html no navegador ou hospedar essa pasta em qualquer servidor web.

## Estrutura de Diretórios
Após o init e o build, seu ambiente de escrita ficará assim:

```plaintxt
meu-blog/
├── content/              # Seus textos em .md
│   └── ola-mundo.md
├── theme/            # (Opcional) Sobrescreva o tema padrão aqui
│   ├── post_template.html
│   ├── index_template.html
│   ├── feed.rss
│   └── static/           # CSS e Imagens do seu tema
├── papiro.yaml           # Configurações globais (Título, URL, Idioma)
└── public/               # O site gerado pronto para deploy! (Gerado pelo build)
    ├── index.html
    ├── ola_mundo.html
    ├── feed.xml
    └── static/
```

## Configuração (papiro.yaml)
Você pode controlar as variáveis globais do seu site editando o arquivo papiro.yaml na raiz do seu projeto:

```
title: Meu Papiro
description: Um registro das minhas ideias.
url: meublog.com.br
language: pt-BR
```
Esses dados são injetados automaticamente nos templates e no feed RSS.

## Contribuindo
Sinta-se à vontade para abrir issues e enviar pull requests. Toda ajuda é bem-vinda para deixar o motor ainda mais robusto!

## Licença
Este projeto está licenciado sob a licença AGPLv3.0 - veja o arquivo LICENSE para mais detalhes.