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
	"strings"
	"sync"
	"time"
)

var mu = &sync.Mutex{}

// GithubRepo is github repository
type GithubRepo struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	IsFork       bool   `json:"fork"`
	URL          string `json:"url"`
	Desc         string `json:"description"`
	LanguagesURL string `json:"languages_url"`
	ContentsURL  string `json:"contents_url"`
	GoImports    []string
}

func (g *GithubRepo) IsGo() bool {
	var result bool

	var langs map[string]int

	rb, _ := GetRequestWithLimit(g.LanguagesURL)

	err := json.Unmarshal(rb, &langs)
	if err != nil {
		return false
	}

	if _, ok := langs["Go"]; ok {
		result = true
	}

	return result

}

func (g *GithubRepo) GetImports() {

	gomodURL := strings.Replace(g.ContentsURL, "{+path}", "go.mod", 1)

	gomodData, _ := GetRequestWithLimit(gomodURL)

	if len(gomodData) > 0 {
		g.GoImports = ParseGomodFile(gomodData)
	}
}

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

	nextURL := "https://api.github.com/repositories"

	for {
		var repos []GithubRepo

		rb, lm := GetRequestWithLimit(nextURL)

		err := json.Unmarshal(rb, &repos)
		if err != nil {
			log.Fatal("Error unmarshaling", err)
		}

		for _, repo := range repos {
			fmt.Println("Processing:", repo.URL)
			if repo.IsGo() && !repo.IsFork {
				repo.GetImports()
				fmt.Println("Go repo! Imports from gomod:\n", strings.Join(repo.GoImports, "\n"))
				goRepos = append(goRepos, repo)
			}
		}

		if val, ok := lm["since"]; ok {
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

	urlRe := regexp.MustCompile("<(.*)>")
	relRe := regexp.MustCompile("rel=\"(.*)\"")
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

	time.Sleep(5 * time.Second)

	client := &http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Fatal("Cant intitialize request:", err)
	}
	req.Header.Add("Accept", "application/vnd.github.mercy-preview+json")

	resp, err := client.Do(req)
	if err != nil {
		resp.Body.Close()
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	rb, _ := ioutil.ReadAll(resp.Body)
	lh := resp.Header.Get("Link")
	lm := ParseLinkHeader(lh)

	return rb, lm
}

// ParseGomodFile get imports from gomod file
func ParseGomodFile(b []byte) []string {
	var goimports []string

	bs := bytes.Split(b, []byte("\n"))

	reStart := regexp.MustCompile(`^require \($`)
	reEnd := regexp.MustCompile(`^\}$`)
	reImport := regexp.MustCompile(`^(.*)\sv`)
	reIndirect := regexp.MustCompile(`indirect`)

loop:
	for i := 0; i < len(bs); i++ {
		if reStart.Match(bs[i]) {
			for {
				i++
				importMatch := reImport.FindSubmatch(bytes.TrimSpace(bs[i]))
				if len(importMatch) == 2 && !reIndirect.Match(bs[i]) {
					goimports = append(goimports, string(importMatch[1]))
				}

				if reEnd.Match(bs[i]) {
					break loop
				}
			}
		}
	}

	return goimports
}
