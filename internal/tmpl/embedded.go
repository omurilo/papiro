package tmpl

import "embed"

//go:embed *.html static/*
var Files embed.FS
