package data

import (
	"embed"
)

//go:embed english.json japanese.json
var DataFS embed.FS
