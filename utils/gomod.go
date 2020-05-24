package utils

import (
	"strings"

	"golang.org/x/mod/modfile"
)

// ParseGomodFile extract required modules from go.mod file
func ParseGomodFile(b []byte) ([]string, error) {
	var requiredModules []string

	goModFile, err := modfile.Parse("", b, nil)
	if err != nil {
		return nil, err
	}

	for _, req := range goModFile.Require {
		if !strings.Contains(req.Syntax.Token[0], "golang.org") { //filtering golang x std packages
			requiredModules = append(requiredModules, req.Syntax.Token[0])
		}
	}
	// fmt.Println(goModFile.Require)

	return requiredModules, nil
}
