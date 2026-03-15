package tmpl

import "embed"

//go:embed *.html *.rss *.yaml static/*
var Files embed.FS
