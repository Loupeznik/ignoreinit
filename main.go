package main

import (
	"github.com/devfacet/gocmd/v3"
	"github.com/loupeznik/ignoreinit/src"
)

const version = "1.1.1"

func main() {
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
