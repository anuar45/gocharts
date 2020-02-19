package main

type Database interface {
	Save(g GithubRepo)
}

type MapDB map[string]GithubRepo

func NewMapDB() MapDB {
	return make(MapDB)
}

func (db MapDB) Save(g GithubRepo) {
	db[g.FullName] = g
}
