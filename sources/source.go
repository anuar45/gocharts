package sources

import "github.com/anuar45/topgomods/model"

// GoRepoSource interface for plugable sources of GO Repos
type GoRepoSource interface {
	Fetch() (model.GoRepos, error)
	Configure(string) error
}

// GoRepoSources collection of GO Repos sources
var (
	GoRepoSources = make(map[string]GoRepoSource)
)
