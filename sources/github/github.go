package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/anuar45/topgomods/model"
	"github.com/anuar45/topgomods/sources"
	"github.com/anuar45/topgomods/utils"
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
func init() {
	sources.GoRepoSources["github"] = new(Github)
}

// Configure configures source plugin
func (g *Github) Configure(cfg string) error {
	// TODO: Parse config and take what needed, fail on missing
	g.Token = os.Getenv("GITHUB_TOKEN")
	if g.Token == "" {
		return errors.New("no github token configured")
	}

	g.SearchURL = "https://api.github.com/search/repositories?q=language:go"
	return nil
}

// Fetch source interface implmentation
func (g *Github) Fetch() (model.GoRepos, error) {

	goRepos, err := g.fetch()

	return goRepos, err
}

func (g *Github) fetch() (model.GoRepos, error) {

	var goRepos model.GoRepos
	nextURL := g.SearchURL

	//l og.Println("starting fetch")

	for {
		var goReposSearch GithubRepoSearch

		content, headers, _ := utils.HTTPGet(nextURL)
		//l og.Println("starting fetch")
		err := json.Unmarshal(content, &goReposSearch)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling: %w", err)
		}

		// log.Println("starting fetch 2")
		repos := goReposSearch.Repos

		for _, repo := range repos {
			if !repo.IsFork && repo.FullName != "golang/go" {
				log.Println("Processing:", repo.RepoURL)

				gomodURL := strings.Replace(repo.ContentsURL, "{+path}", "go.mod", 1)

				log.Println("Go mod file url:", gomodURL)
				gomodContent, _, _ := utils.HTTPGetWithHeaders(gomodURL,
					map[string]string{
						"Authorization": "Bearer " + g.Token,
						"Accept":        "application/vnd.github.VERSION.raw",
					})

				var goModules []model.GoModule
				log.Println(string(gomodContent))
				if len(gomodContent) > 0 {
					modules, _ := utils.ParseGomodFile(gomodContent)

					for _, module := range modules {
						goModules = append(goModules, model.GoModule{path.Base(module), module})
					}
				}

				goRepos = append(goRepos, model.GoRepo{repo.Name, repo.RepoURL, goModules})
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
