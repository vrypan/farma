package main

import (
	"embed"

	"github.com/vrypan/farma/cmd"
)

//go:embed test_frame/* test_frame/.well-known
var staticFiles embed.FS

func main() {
	cmd.StaticFiles = staticFiles
	cmd.Execute()
}
