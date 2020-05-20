package topgomods

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

// HTTPGet makes retryable http get requests using 3rd lib
func HTTPGet(url, token string) ([]byte, map[string][]string, error) {
	headers := make(map[string][]string)

	log.Println("quering")
	client := retryablehttp.NewClient()

	log.Println("quering")
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
