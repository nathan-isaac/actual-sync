package main

import (
	"embed"
	"github.com/nathanjisaac/actual-server-go/cmd"
)

//go:embed node_modules/@actual-app/web/build/*
var buildDirectory embed.FS

func main() {
	cmd.BuildDirectory = buildDirectory

	cmd.Execute()
}
