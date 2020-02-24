package main

type GithubRepoDatabase interface {
	Save(g GithubRepo)
	GetAll() []GithubRepo
}

type GoImportDatabase interface {
	Save(g GoImport)
	GetAll() []GoImport
}

type GithubRepoDB map[string]GithubRepo
type GoImportDB map[string]GoImport

func NewGoImportDB() GoImportDB {
	return make(GoImportDB)
}

func NewGithubRepoDB() GithubRepoDB {
	return make(GithubRepoDB)
}

func (db GoImportDB) Save(g GoImport) {
	db[g.URL] = g
}

func (db GoImportDB) GetAll() []GoImport {
	var g []GoImport
	for _, v := range db {
		g = append(g, v)
	}
	return g
}

func (db GithubRepoDB) Save(g GithubRepo) {
	db[g.FullName] = g
}

func (db GithubRepoDB) GetAll() []GithubRepo {
	var g []GithubRepo
	for _, v := range db {
		g = append(g, v)
	}
	return g
}
