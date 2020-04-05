package main

type GithubRepoRepository interface {
	Save(g GithubRepo)
	FindAll() []GithubRepo
}

type GoImportRepository interface {
	Save(g GoImport)
	FindAll() ([]GoImport, error)
}
