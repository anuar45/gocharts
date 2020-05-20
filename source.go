package topgomods

// GoRepoSource interface for plugable sources of GO Repos
type GoRepoSource interface {
	Fetch() (GoRepos, error)
	Configure(string) error
}

// GoRepoSources collection of GO Repos sources
var (
	GoRepoSources = make(map[string]GoRepoSource)
)
