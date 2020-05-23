package utils

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

func TestParseGomod(t *testing.T) {
	want := []string{
		"bitbucket.org/bertimus9/systemstat",
		"github.com/Rican7/retry",
	}
	got, err := ParseGomodFile([]byte(gomodFile))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, want, got, "Should be equal")

}
