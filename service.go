package main

// Service object for GoImport service
type GoImportServicer interface {
	Fetch() error
	FindAll() ([]GoImport, error)
}

type GoImportService struct {
	IsBusy  bool
	GisRepo GoImportRepository
	GrsRepo GithubRepoRepository
}

func NewGoImportService(gr GithubRepoRepository, gi GoImportRepository) *GoImportService {
	return &GoImportService{
		GrsRepo: gr,
		GisRepo: gi,
	}
}

func (g *GoImportService) FindAll() ([]GoImport, error) {
	return g.GisRepo.FindAll()
}
