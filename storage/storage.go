package storage

import "github.com/anuar45/topgomods/model"

type GoRepoRepository interface {
	Save(model.GoRepo)
	FindAll() []model.GoRepo
}

type GoModuleRepository interface {
	Save(g model.GoModule)
	FindAll() ([]model.GoModule, error)
}
