package topgomods

import (
	"sync"
)

type GoRepoDB struct {
	mapdb map[string]GoRepo
	mutex *sync.RWMutex
}

type GoModuleDB struct {
	mapdb map[string]GoModule
	mutex *sync.RWMutex
}

func NewGoModuleDB() *GoModuleDB {
	return &GoModuleDB{
		mapdb: make(map[string]GoModule),
		mutex: new(sync.RWMutex),
	}
}

func NewGoRepoDB() *GoRepoDB {
	return &GoRepoDB{
		mapdb: make(map[string]GoRepo),
		mutex: new(sync.RWMutex),
	}
}

func (db GoModuleDB) Save(g GoModule) {
	db.mutex.Lock()
	db.mapdb[g.URL] = g
	db.mutex.Unlock()
}

func (db GoModuleDB) FindAll() ([]GoModule, error) {
	var g []GoModule
	for _, v := range db.mapdb {
		g = append(g, v)
	}

	//sort.Slice(g, func(i, j int) bool { return g[i].Count > g[j].Count })
	return g, nil
}

func (db GoRepoDB) Save(g GoRepo) {
	db.mutex.Lock()
	db.mapdb[g.URL] = g
	db.mutex.Unlock()
}

func (db GoRepoDB) FindAll() []GoRepo {
	var g []GoRepo
	for _, v := range db.mapdb {
		g = append(g, v)
	}

	return g
}
