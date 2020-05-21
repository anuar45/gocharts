package model

import (
	"path"
	"strings"

	"golang.org/x/mod/modfile"
)

// GoRepo is golang project
type GoRepo struct {
	Name    string
	URL     string
	Modules []GoModule
}

//GoRepos slice of GoRepo
type GoRepos []GoRepo

// GoModule is go pkg/lib
type GoModule struct {
	Name string
	URL  string
}

// GoModuleRank is imports count of go module
type GoModuleRank struct {
	URL   string
	Count int
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
