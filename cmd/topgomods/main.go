package main

import (
	"log"

	"github.com/anuar45/topgomods"
	_ "github.com/anuar45/topgomods/sources/github"
)

// VERSION of app populated on build
var VERSION string

func main() {

	goRepoDB := topgomods.NewGoRepoDB()

	goModuleDB := topgomods.NewGoModuleDB()

	goModuleService := topgomods.NewGoModuleService(goRepoDB, goModuleDB)

	apiServer := topgomods.NewApiServer(goModuleService)

	log.Println("Starting server...")
	apiServer.Run()

}
