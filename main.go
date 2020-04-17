package main

var VERSION string

func main() {

	githubRepoDB := NewGithubRepoDB()

	goImportDB := NewGoImportDB()

	goImportService := NewGoImportService(githubRepoDB, goImportDB)

	apiServer := NewApiServer(goImportService)

	apiServer.Run()

}
