package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
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

type Github struct {
	Token     string
	BaseURL   string
	SearchURL string
	RateLimit int
}

func Init() *Github {
	return &Github{
		SearchURL: "https://api.github.com/search/repositories?q=language:go",
	}
}

func (g *Github) Fetch() ([]GoRepo, error) {

	goRepos, err := g.fetch()

	return goRepos, err
}

func (g *Github) fetch() ([]GoRepo, error) {

	var goRepos []GoRepo
	nextURL := g.SearchURL

	for {
		var goReposSearch GithubRepoSearch

		content, headers, _ := HttpGet(nextURL, g.Token)

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
				gomodContent, _, _ := HttpGet(gomodURL, g.Token)

				var goModules []GoModule
				if len(gomodContent) > 0 {
					goModules, _ = ParseGomodFile(gomodContent)
				}

				goRepos = append(goRepos, GoRepo{repo.Name, repo.RepoURL, goModules})
				log.Println("Modules:\n", goModules)
			}
		}

		// TODO:
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

// HttpGet makes retryable http get requests using 3rd lib
func HttpGet(url, token string) ([]byte, map[string][]string, error) {
	headers := make(map[string][]string)

	client := retryablehttp.NewClient()

	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("cant intitialize request: %w", err)
	}

	if token == "" {
		return nil, nil, errors.New("No github token found")
	}

	req.Header.Add("Authorization", "token "+token)

	resp, err := client.Do(req)
	if err != nil {
		resp.Body.Close()
		return nil, nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading response body: %w", err)
	}

	headers = resp.Header

	return body, headers, nil
}
