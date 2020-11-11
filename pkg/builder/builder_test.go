package builder

import (
	"testing"
)

func TestBuildTokenFlags(t *testing.T) {
	tokenParams := buildTokenFlags("token", "host", true)
	if tokenParams != " --auth-token token --server host --insecure --prune" {
		t.Errorf("Wrong \"%s\" token command ", tokenParams)
	}

}

func TestBuildTokenFlagsWithoutPrune(t *testing.T) {
	tokenParams := buildTokenFlags("token", "host", false)
	if tokenParams != " --auth-token token --server host --insecure" {
		t.Errorf("Wrong \"%s\" token command ", tokenParams)
	}

}
