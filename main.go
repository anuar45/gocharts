package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	githubAPI     = "api.github.com"
	searchRepoURI = "/search/repositories?q="
)

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

	req, err := http.NewRequest("GET", githubAPI, nil)
	if err != nil {
		log.Fatal("Cant intitialize request:", err)
	}
	req.Header.Add("Accept", "application/vnd.github.mercy-preview+json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making sending request:", err)
	}

	fmt.Println(resp)
}
