package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var mu = &sync.Mutex{}

const SearchURL = "https://api.github.com/search/repositories?q=language:go"

func (g *GoImportService) Fetch() error {
	if g.IsBusy {
		return errors.New("Fetch already in progress...")
	}
	go func() {
		g.IsBusy = true
		g.fetch()
		g.IsBusy = false
	}()
	return nil
}

// Fetch calculates go mod imports statistics
func (g *GoImportService) fetch() error {

	countImports := make(map[string]int)
	nextURL := SearchURL

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
				fmt.Println("Processing:", repo.RepoURL)
				g.GrsRepo.Save(repo)

				gomodURL := strings.Replace(repo.ContentsURL, "{+path}", "go.mod", 1)
				fmt.Println("GO mod file url:", gomodURL)
				gomodData, _ := GetRequestWithLimit(gomodURL)

				if len(gomodData) > 0 {

					goImports := ParseGomodFile(gomodData)
					repo.GoImports = goImports

					for _, goImport := range goImports {
						countImports[goImport]++
					}
				}

				fmt.Println("Imports from gomod:\n", strings.Join(repo.GoImports, "\n"))
				for k, v := range countImports {
					g.GisRepo.Save(GoImport{k, v})
				}
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
		log.Println("No token in GITHUB_TOKEN env")
		return nil, nil
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

// ParseGomodFile get imports from gomod file
func ParseGomodFile(b []byte) []string {
	var goimports []string

	bs := bytes.Split(b, []byte("\n"))

	reStart := regexp.MustCompile(`^require \($`)
	reEnd := regexp.MustCompile(`\}`)
	reImport := regexp.MustCompile(`^(.*)\sv`)
	reIndirect := regexp.MustCompile(`indirect`)

	var requireBlock bool

	for i := 0; i < len(bs); i++ {
		if reStart.Match(bs[i]) {
			requireBlock = true
		}
		if requireBlock {
			importMatch := reImport.FindSubmatch(bytes.TrimSpace(bs[i]))
			if len(importMatch) == 2 && !reIndirect.Match(bs[i]) {
				goimports = append(goimports, string(importMatch[1]))
			}

			if reEnd.Match(bs[i]) {
				break
			}
		}
	}

	return goimports
}
