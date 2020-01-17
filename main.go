package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// GithubRepos is github repositories root
type GithubRepoSearchResult struct {
	TotalCount int          `json:"total_count"`
	Repos      []GithubRepo `json:"items"`
	Incomplete bool         `json:"incomplete_results"`
}

// GithubRepo is github repository
type GithubRepo struct {
	//ID int `json:"id"`
	Fork bool   `json:"fork"`
	URL  string `json:"url"`
	Desc string `json:"description"`
}

// GithubItem is github content item
type GithubItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func main() {
	result := SearchGithubRepos("go")
	fmt.Println(result)
}

// SearchGithubRepos searches for repos with specific lang
func SearchGithubRepos(lang string) []GithubRepo {
	var g GithubRepoSearchResult
	u := url.URL{}
	u.Scheme = "https"
	u.Host = "api.github.com"
	u.Path = "/search/repositories"
	q := u.Query()
	q.Set("q", "language:"+lang)
	u.RawQuery = q.Encode()
	//fmt.Println(u)

	client := &http.Client{}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatal("Cant intitialize request:", err)
	}
	req.Header.Add("Accept", "application/vnd.github.mercy-preview+json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&g)
	if err != nil {
		log.Fatal("Error unmarshaling", err)
	}
	return g
}
