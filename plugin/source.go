package plugin

import "github.com/anuar45/topgomods"

// GoRepoSource interface for plugable sources of GO Repos
type GoRepoSource interface {
	Fetch() (topgomods.GoRepos, error)
	Configure(string) error
}

// GoRepoSources collection of GO Repos sources
var GoRepoSources map[string]GoRepoSource
