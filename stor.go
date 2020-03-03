package main

type GithubRepoDatabase interface {
	SaveAll(g GithubRepo)
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

func (db GoImportDB) SaveAll(gis []GoImport) {
	for _, v := range gis {
		db[v.URL] = v
	}
}

func (db GoImportDB) GetAll() []GoImport {
	var g []GoImport
	for _, v := range db {
		g = append(g, v)
	}
	return g
}

func (db GithubRepoDB) SaveAll(grs []GithubRepo) {
	for _, v := range grs {
		db[v.FullName] = v
	}
}

func (db GithubRepoDB) GetAll() []GithubRepo {
	var g []GithubRepo
	for _, v := range db {
		g = append(g, v)
	}
	return g
}
