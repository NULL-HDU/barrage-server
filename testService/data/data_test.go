package data

import (
	"testing"
)

func TestReply(t *testing.T) {
	t.Logf("result: % x", Reply)
}

func TestRandomUserID(t *testing.T) {
	t.Logf("result: % x", RandomUserID())
}
