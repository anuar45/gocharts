package utils

import (
	"fmt"
	"io/ioutil"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

// HTTPGetWithHeaders makes retryable http get requests using 3rd lib
func HTTPGetWithHeaders(url string, headers map[string]string) ([]byte, map[string][]string, error) {
	respHeaders := make(map[string][]string)

	//log.Println("quering")
	client := retryablehttp.NewClient()
	client.RetryMax = 100
	client.RetryWaitMin = 5 * time.Second
	client.RetryWaitMax = 1 * time.Minute

	//log.Println("quering")
	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("cant intitialize request: %w", err)
	}

	for header, value := range headers {
		req.Header.Add(header, value)
	}

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

	respHeaders = resp.Header

	return body, respHeaders, nil
}

// HTTPGet wraps HTTPGetWithHeaders
func HTTPGet(url string) ([]byte, map[string][]string, error) {
	reqHeaders := make(map[string]string)

	return HTTPGetWithHeaders(url, reqHeaders)
}
