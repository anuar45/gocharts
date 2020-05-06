package main

import (
	"path"
	"sort"
	"strings"

	"golang.org/x/mod/modfile"
)

// Service object for GoImport service
type GoModuleServicer interface {
	Fetch() error
	Repos() ([]GoRepo, error)
	TopModules() ([]GoModuleRank, error)
}

type GoModuleService struct {
	IsBusy  bool
	GmsRepo GoModuleRepository
	GrsRepo GoRepoRepository
}

func NewGoModuleService(gr GoRepoRepository, gm GoModuleRepository) *GoModuleService {
	return &GoModuleService{
		GrsRepo: gr,
		GmsRepo: gm,
	}
}

func (g *GoModuleService) Repos() ([]GoRepo, error) {
	return g.GrsRepo.FindAll(), nil
}

func (g *GoModuleService) TopModules() ([]GoModuleRank, error) {
	var moduleRanks []GoModuleRank

	modulesCount := make(map[string]int)

	repos, _ := g.Repos()

	for _, repo := range repos {
		for _, module := range repo.Modules {
			modulesCount[module.URL]++
		}
	}

	for module, n := range modulesCount {
		moduleRanks = append(moduleRanks, GoModuleRank{module, n})
	}

	sort.Slice(moduleRanks, func(i, j int) bool { return moduleRanks[i].Count > moduleRanks[j].Count })

	return moduleRanks, nil
}

// ParseGomodFile get imports from gomod file
func ParseGomodFile(b []byte) ([]GoModule, error) {
	var modules []GoModule

	goModFile, err := modfile.Parse("", b, nil)
	if err != nil {
		return nil, err
	}

	for _, req := range goModFile.Require {
		if !strings.Contains(req.Syntax.Token[0], "golang.org") { //filtering golang std packages
			modules = append(modules, GoModule{path.Base(req.Syntax.Token[0]), req.Syntax.Token[0]})
		}
	}

	return modules, nil
}
