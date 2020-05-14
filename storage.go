package topgomods

type GoRepoRepository interface {
	Save(g GoRepo)
	FindAll() []GoRepo
}

type GoModuleRepository interface {
	Save(g GoModule)
	FindAll() ([]GoModule, error)
}
