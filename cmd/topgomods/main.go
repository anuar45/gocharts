package main

import (
	"log"

	"github.com/anuar45/topgomods/api/rest"
	"github.com/anuar45/topgomods/service"
	_ "github.com/anuar45/topgomods/sources/github"
	"github.com/anuar45/topgomods/storage/memdb"
)

// VERSION of app populated on build
var VERSION string

func main() {

	goRepoDB := memdb.NewGoRepoDB()

	goModuleDB := memdb.NewGoModuleDB()

	goModuleService := service.NewGoModuleService(goRepoDB, goModuleDB, VERSION)

	apiServer := rest.NewApiServer(goModuleService)

	log.Println("Starting server...")
	apiServer.Run()

}
