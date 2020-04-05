package main

// Just for parsing search response
type GithubRepoSearch struct {
	Repos []GithubRepo `json:"items"`
}

// GithubRepo is github repository
type GithubRepo struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	FullName     string `json:"full_name"`
	IsFork       bool   `json:"fork"`
	RepoURL      string `json:"url"`
	Desc         string `json:"description"`
	LanguagesURL string `json:"languages_url"`
	ContentsURL  string `json:"contents_url"`
	GoImports    []string
}

// Count package imports
type GoImport struct {
	URL   string
	Count int
}
