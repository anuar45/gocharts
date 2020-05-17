package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/anuar45/topgomods"
)

// GithubRepoSearch is search response from github
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
}

// Github is github source of Go Repos
type Github struct {
	Token     string
	BaseURL   string
	SearchURL string
	RateLimit int
}

// Init registers source plugin
func Init() {
	topgomods.GoRepoSources["github"] = new(Github)
}

// Configure configures source plugin
func (g *Github) Configure(cfg string) error {
	// TODO: Parse config and take what needed, fail on missing
	g.Token = os.Getenv("GITHUB_TOKEN")
	if g.Token == "" {
		return errors.New("no github token configured")
	}
	return nil
}

// Fetch source interface implmentation
func (g *Github) Fetch() (topgomods.GoRepos, error) {

	goRepos, err := g.fetch()

	return goRepos, err
}

func (g *Github) fetch() (topgomods.GoRepos, error) {

	var goRepos topgomods.GoRepos
	nextURL := g.SearchURL

	for {
		var goReposSearch GithubRepoSearch

		content, headers, _ := topgomods.HTTPGet(nextURL, g.Token)

		err := json.Unmarshal(content, &goReposSearch)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling: %w", err)
		}

		repos := goReposSearch.Repos

		for _, repo := range repos {
			if !repo.IsFork && repo.FullName != "golang/go" {
				log.Println("Processing:", repo.RepoURL)

				gomodURL := strings.Replace(repo.ContentsURL, "{+path}", "go.mod", 1)
				log.Println("GO mod file url:", gomodURL)
				gomodContent, _, _ := topgomods.HTTPGet(gomodURL, g.Token)

				var goModules []topgomods.GoModule
				if len(gomodContent) > 0 {
					goModules, _ = topgomods.ParseGomodFile(gomodContent)
				}

				goRepos = append(goRepos, topgomods.GoRepo{repo.Name, repo.RepoURL, goModules})
				log.Println("Modules:\n", goModules)
			}
		}

		lh := headers["Link"]
		links := ParseLinkHeader(lh[0])

		if val, ok := links["next"]; ok {
			nextURL = val
		} else {
			break
		}
	}

	return goRepos, nil
}

// ParseLinkHeader gets reference links from headers
func ParseLinkHeader(s string) map[string]string {
	links := make(map[string]string)

	sl := strings.Split(s, ",")

	urlRe := regexp.MustCompile(`<(.*)>`)
	relRe := regexp.MustCompile(`rel=\"(.*)\"`)
	for _, line := range sl {
		uri := urlRe.FindStringSubmatch(line)
		rel := relRe.FindStringSubmatch(line)
		if len(uri) == 2 && len(rel) == 2 {
			links[rel[1]] = uri[1]
		}
	}
	return links
}
