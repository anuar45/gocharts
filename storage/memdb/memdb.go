package memdb

import (
	"sync"

	"github.com/anuar45/topgomods/model"
)

type GoRepoDB struct {
	mapdb map[string]model.GoRepo
	mutex *sync.RWMutex
}

type GoModuleDB struct {
	mapdb map[string]model.GoModule
	mutex *sync.RWMutex
}

func NewGoModuleDB() *GoModuleDB {
	return &GoModuleDB{
		mapdb: make(map[string]model.GoModule),
		mutex: new(sync.RWMutex),
	}
}

func NewGoRepoDB() *GoRepoDB {
	return &GoRepoDB{
		mapdb: make(map[string]model.GoRepo),
		mutex: new(sync.RWMutex),
	}
}

func (db GoModuleDB) Save(g model.GoModule) {
	db.mutex.Lock()
	db.mapdb[g.URL] = g
	db.mutex.Unlock()
}

func (db GoModuleDB) FindAll() ([]model.GoModule, error) {
	var g []model.GoModule
	for _, v := range db.mapdb {
		g = append(g, v)
	}

	//sort.Slice(g, func(i, j int) bool { return g[i].Count > g[j].Count })
	return g, nil
}

func (db GoRepoDB) Save(g model.GoRepo) {
	db.mutex.Lock()
	db.mapdb[g.URL] = g
	db.mutex.Unlock()
}

func (db GoRepoDB) FindAll() []model.GoRepo {
	var g []model.GoRepo
	for _, v := range db.mapdb {
		g = append(g, v)
	}

	return g
}
