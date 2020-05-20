package topgomods

import (
	"errors"
	"log"
	"path"
	"sort"
	"strings"

	"golang.org/x/mod/modfile"
)

// Service object for GoImport service
type GoModuleServicer interface {
	Fetch() error
	Repos() (GoRepos, error)
	TopModules() ([]GoModuleRank, error)
}

type GoModuleService struct {
	Config  string
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

func (g *GoModuleService) Fetch() error {

	if g.IsBusy == true {
		return errors.New("fetch already in progress")
	}

	log.Println("Starting fetch")

	go func() {
		g.IsBusy = true
		for name, source := range GoRepoSources {
			err := source.Configure(g.Config)
			if err != nil {
				log.Println("error configuring plugin:", name, err)
				continue
			}
			goRepos, err := source.Fetch()
			if err != nil {
				log.Println("error running fetch from source:", name, err)
				continue
			}

			for _, goRepo := range goRepos {
				g.GrsRepo.Save(goRepo)
			}
		}
		g.IsBusy = false
	}()

	return nil
}

func (g *GoModuleService) Repos() (GoRepos, error) {
	return g.GrsRepo.FindAll(), nil
}

func (g *GoModuleService) TopModules() ([]GoModuleRank, error) {
	var moduleRanks []GoModuleRank

	modulesCount := make(map[string]int)

	goRepos, _ := g.Repos()

	for _, goRepo := range goRepos {
		for _, module := range goRepo.Modules {
			modulesCount[module.URL]++
		}
	}

	for module, n := range modulesCount {
		moduleRanks = append(moduleRanks, GoModuleRank{module, n})
	}

	sort.Slice(moduleRanks, func(i, j int) bool { return moduleRanks[i].Count > moduleRanks[j].Count })

	return moduleRanks, nil
}

// ParseGomodFile get packages from gomod file
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
