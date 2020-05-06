package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var gomodFile = `
// This is a generated file. Do not edit directly.
// Run hack/pin-dependency.sh to change pinned dependency versions.
// Run hack/update-vendor.sh to update go.mod files and the vendor directory.

module k8s.io/kubernetes

go 1.13

require (
	bitbucket.org/bertimus9/systemstat v0.0.0-20180207000608-0eeff89b0690
	github.com/Rican7/retry v0.1.0 // indirect
)


replace (
	bitbucket.org/bertimus9/systemstat => bitbucket.org/bertimus9/systemstat v0.0.0-20180207000608-0eeff89b0690
	cloud.google.com/go => cloud.google.com/go v0.38.0
)
`

var goRepo = GithubRepo{
	ID:           20580498,
	FullName:     "kubernetes/kubernetes",
	IsFork:       false,
	RepoURL:      "https://api.github.com/repos/kubernetes/kubernetes",
	Desc:         "Production-Grade Container Scheduling and Management",
	LanguagesURL: "https://api.github.com/repos/kubernetes/kubernetes/languages",
	ContentsURL:  "https://api.github.com/repos/kubernetes/kubernetes/contents/{+path}",
}

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

func TestParseGomod(t *testing.T) {
	want := []GoModule{
		GoModule{"systemstat", "bitbucket.org/bertimus9/systemstat"},
		GoModule{"retry", "github.com/Rican7/retry"},
	}
	got, err := ParseGomodFile([]byte(gomodFile))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, want, got, "Should be equal")

}
