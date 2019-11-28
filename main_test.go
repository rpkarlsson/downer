package main_test

import (
	"github.com/rpkarlsson/downer"
	"testing"
)

func TestIsSeen(t *testing.T) {
	if main.CheckFolderForFile(".", "notAFile") {
		t.Error("Sees files that doesn't exist")
	}
	if !main.CheckFolderForFile(".", "go.mod") {
		t.Error("Doesn't see files that exist")
	}
}
