package main

import (
	"github.com/anuar45/topgomods"
	_ "github.com/anuar45/topgomods/source/github"
)

func main() {

	goRepoDB := topgomods.NewGoRepoDB()

	goModuleDB := topgomods.NewGoModuleDB()

	goModuleService := topgomods.NewGoModuleService(goRepoDB, goModuleDB)

	apiServer := topgomods.NewApiServer(goModuleService)

	apiServer.Run()

}
