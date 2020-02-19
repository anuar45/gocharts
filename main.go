package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

var mu = &sync.Mutex{}

const SearchURL = "https://api.github.com/search/repositories?q=language:go"

func main() {
	result := GetGithubGoRepos()
	fmt.Println(len(result))
	f, _ := os.Create("data.out")
	bw := bufio.NewWriter(f)
	for _, repo := range result {
		bw.WriteString(strings.Join([]string{repo.Name, strings.Join(repo.GoImports, ":"), "\n"}, "\t"))
	}
	bw.Flush()
}

// GetGithubGoRepos searches for repos with specific lang
func GetGithubGoRepos() []GithubRepo {
	var goRepos []GithubRepo

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
			fmt.Println("Processing:", repo.RepoURL)
			if !repo.IsFork && repo.FullName != "golang/go" {
				gomodURL := strings.Replace(repo.ContentsURL, "{+path}", "go.mod", 1)
				fmt.Println("GO mod file url:", gomodURL)
				gomodData, _ := GetRequestWithLimit(gomodURL)

				if len(gomodData) > 0 {
					repo.GoImports = ParseGomodFile(gomodData)
				}

				fmt.Println("Imports from gomod:\n", strings.Join(repo.GoImports, "\n"))
				goRepos = append(goRepos, repo)
			}
		}

		if val, ok := links["next"]; ok {
			nextURL = val
		} else {
			break
		}
	}
	return goRepos
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

func GetGoImports(gr []GithubRepo) []GoImport {
	var goimports []GoImport
	goimportsMap := make(map[string]int)
	for _, repo := range gr {
		for _, goimport := range repo.GoImports {
			goimportsMap[goimport]++
		}
	}

	for goimport, count := range goimportsMap {
		goimports = append(goimports, GoImport{goimport, count})
	}

	sort.Slice(goimports, func(i, j int) bool {
		return goimports[i].Count > goimports[j].Count
	})

	return goimports
}
