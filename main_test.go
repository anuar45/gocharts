package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLinkHeader(t *testing.T) {
	testcase := struct {
		in   string
		want map[string]string
	}{
		"<https://api.github.com/search/repositories?q=language%3Ago&page=2>; rel=\"next\", <https://api.github.com/search/repositories?q=language%3Ago&page=33>; rel=\"last\"",
		map[string]string{
			"next": "https://api.github.com/search/repositories?q=language%3Ago&page=2",
			"last": "https://api.github.com/search/repositories?q=language%3Ago&page=33",
		},
	}
	got := ParseLinkHeader(testcase.in)
	assert.Equal(t, got, testcase.want, "Should be equal")

}
