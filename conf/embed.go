package conf

import (
	"embed"
)

//go:embed *.ini
var Files embed.FS
