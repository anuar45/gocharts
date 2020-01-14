package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	githubAPI     = "api.github.com"
	searchRepoURI = "/search/repositories?q="
)

type GithubRepos struct {
	TotalCount int          `json:"total_count"`
	Repos      []GithubRepo `json:"items"`
	Incomplete bool         `json:"incomplete_results"`
}

type GithubRepo struct {
	//ID int `json:"id"`
	Fork bool   `json:"fork"`
	Url  string `json:"url"`
	Desc string `json:"description"`
}

func main() {
	u := url.URL{}
	u.Scheme = "https"
	u.Host = "api.github.com"
	u.Path = "/search/repositories"
	q := u.Query()
	q.Set("q", "language:go")
	u.RawQuery = q.Encode()
	fmt.Println(u)

	client := &http.Client{}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatal("Cant intitialize request:", err)
	}
	req.Header.Add("Accept", "application/vnd.github.mercy-preview+json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making sending request:", err)
	}

	var g GithubRepos

	err = json.NewDecoder(resp.Body).Decode(&g)
	if err != nil {
		log.Fatal("Error unmarshaling", err)
	}
	fmt.Println(resp)
	fmt.Println(g)
}
