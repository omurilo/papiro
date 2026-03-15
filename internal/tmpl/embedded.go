package tmpl

import "embed"

//go:embed *.html *.rss static/*
var Files embed.FS
