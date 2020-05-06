package main

var VERSION string

func main() {

	goRepoDB := NewGoRepoDB()

	goModuleDB := NewGoModuleDB()

	goModuleService := NewGoModuleService(goRepoDB, goModuleDB)

	apiServer := NewApiServer(goModuleService)

	apiServer.Run()

}
