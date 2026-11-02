package template

import "embed"

//go:embed layout/* shared/* page/*
var FS embed.FS
