package main

import (
	"sync"
)

type GithubRepoDB struct {
	mapdb map[string]GithubRepo
	mutex *sync.RWMutex
}

type GoImportDB struct {
	mapdb map[string]GoImport
	mutex *sync.RWMutex
}

func NewGoImportDB() *GoImportDB {
	return &GoImportDB{
		mapdb: make(map[string]GoImport),
		mutex: new(sync.RWMutex),
	}
}

func NewGithubRepoDB() *GithubRepoDB {
	return &GithubRepoDB{
		mapdb: make(map[string]GithubRepo),
		mutex: new(sync.RWMutex),
	}
}

func (db GoImportDB) Save(g GoImport) {
	db.mutex.Lock()
	db.mapdb[g.URL] = g
	db.mutex.Unlock()
}

func (db GoImportDB) FindAll() ([]GoImport, error) {
	var g []GoImport
	for _, v := range db.mapdb {
		g = append(g, v)
	}
	return g, nil
}

func (db GithubRepoDB) Save(g GithubRepo) {
	db.mutex.Lock()
	db.mapdb[g.FullName] = g
	db.mutex.Unlock()
}

func (db GithubRepoDB) FindAll() []GithubRepo {
	var g []GithubRepo
	for _, v := range db.mapdb {
		g = append(g, v)
	}
	return g
}
