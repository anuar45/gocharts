package service

import (
	"errors"
	"fmt"
	"log"
	"path"
	"sort"
	"strings"

	"github.com/anuar45/topgomods/model"
	"github.com/anuar45/topgomods/sources"
	"github.com/anuar45/topgomods/storage"
	"golang.org/x/mod/modfile"
)

// Service object for GoImport service
type GoModuleServicer interface {
	GetVersion() string
	Fetch() error
	Repos() (model.GoRepos, error)
	TopModules() ([]model.GoModuleRank, error)
}

type GoModuleService struct {
	Version string
	Config  string
	IsBusy  bool
	GmsRepo storage.GoModuleRepository
	GrsRepo storage.GoRepoRepository
}

func NewGoModuleService(gr storage.GoRepoRepository, gm storage.GoModuleRepository, version string) *GoModuleService {
	return &GoModuleService{
		Version: version,
		GrsRepo: gr,
		GmsRepo: gm,
	}
}

func (g *GoModuleService) Fetch() error {

	if g.IsBusy == true {
		return errors.New("fetch already in progress")
	}

	//log.Println("Starting fetch")

	go func() {
		g.IsBusy = true
		for name, source := range sources.GoRepoSources {
			log.Println("Starting fetch source:", name)
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

func (g *GoModuleService) Repos() (model.GoRepos, error) {
	return g.GrsRepo.FindAll(), nil
}

func (g *GoModuleService) GetVersion() string {
	return g.Version
}

func (g *GoModuleService) TopModules() ([]model.GoModuleRank, error) {
	var moduleRanks []model.GoModuleRank

	modulesCount := make(map[string]int)

	goRepos, _ := g.Repos()

	for _, goRepo := range goRepos {
		for _, module := range goRepo.Modules {
			modulesCount[module.URL]++
		}
	}

	for module, n := range modulesCount {
		moduleRanks = append(moduleRanks, model.GoModuleRank{module, n})
	}

	sort.Slice(moduleRanks, func(i, j int) bool { return moduleRanks[i].Count > moduleRanks[j].Count })

	return moduleRanks, nil
}

// ParseGomodFile extract modules from go.mod file
func ParseGomodFile(b []byte) ([]model.GoModule, error) {
	var modules []model.GoModule

	goModFile, err := modfile.Parse("", b, nil)
	if err != nil {
		return nil, err
	}

	for _, req := range goModFile.Require {
		if !strings.Contains(req.Syntax.Token[0], "golang.org") { //filtering golang x std packages
			modules = append(modules, model.GoModule{path.Base(req.Syntax.Token[0]), req.Syntax.Token[0]})
		}
	}
	fmt.Println(goModFile.Require)

	return modules, nil
}
