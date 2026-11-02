package template

import "embed"

//go:embed all:layout all:page all:shared
var FS embed.FS
