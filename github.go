package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
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

var mu = &sync.Mutex{}

const searchURL = "https://api.github.com/search/repositories?q=language:go"

func (g *GoModuleService) Fetch() error {
	if g.IsBusy {
		return errors.New("fetch already in progress")
	}
	go func() {
		g.IsBusy = true
		g.fetch()
		g.IsBusy = false
	}()
	return nil
}

func (g *GoModuleService) fetch() error {
	nextURL := searchURL

	for {
		var goReposSearch GithubRepoSearch

		rb, links := GetRequestWithLimit(nextURL)

		err := json.Unmarshal(rb, &goReposSearch)
		if err != nil {
			log.Fatal("Error unmarshaling", err)
		}

		repos := goReposSearch.Repos

		for _, repo := range repos {
			if !repo.IsFork && repo.FullName != "golang/go" {
				log.Println("Processing:", repo.RepoURL)

				gomodURL := strings.Replace(repo.ContentsURL, "{+path}", "go.mod", 1)
				log.Println("GO mod file url:", gomodURL)
				gomodData, _ := GetRequestWithLimit(gomodURL)

				var goModules []GoModule
				if len(gomodData) > 0 {
					goModules, _ = ParseGomodFile(gomodData)
				}

				g.GrsRepo.Save(GoRepo{repo.Name, repo.RepoURL, goModules})
				log.Println("Modules:\n", goModules)
			}
		}

		if val, ok := links["next"]; ok {
			nextURL = val
		} else {
			break
		}
	}

	return nil
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

// GetRequestWithLimit is attemp to make simple rate limiter using sleep and mutex
func GetRequestWithLimit(u string) ([]byte, map[string]string) {
	mu.Lock()
	defer mu.Unlock()

	time.Sleep(2 * time.Second)

	client := &http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Fatal("Cant intitialize request:", err)
	}
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("No token in GITHUB_TOKEN env")
	}

	req.Header.Add("Authorization", "token "+token)
	req.Header.Add("Accept", "application/vnd.github.VERSION.raw")

	resp, err := client.Do(req)
	if err != nil {
		resp.Body.Close()
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	rb, _ := ioutil.ReadAll(resp.Body)
	lh := resp.Header.Get("Link")
	links := ParseLinkHeader(lh)

	return rb, links
}
