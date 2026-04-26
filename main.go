package main

import (
	"os"

	"github.com/devfacet/gocmd/v3"
	"github.com/loupeznik/ignoreinit/src"
)

var version = "dev"

func main() {
	os.Args = src.NormalizeGenerationPrintArgs(os.Args)

	flags := src.Flags{}

	src.InitHandlers()

	gocmd.New(gocmd.Options{
		Name:        "ignoreinit",
		Description: "Create .gitignore from the command line",
		Version:     version,
		Flags:       &flags,
		ConfigType:  gocmd.ConfigTypeAuto,
	})
}
