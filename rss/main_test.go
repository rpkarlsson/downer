package rss_test

import (
	"testing"

	"github.com/rpkarlsson/downer/rss"
)

func TestRead(t *testing.T) {
	history := rss.History{}

	item1 := rss.Item{Title: "Item 1"}
	history.Add(item1)

	item2 := rss.Item{Title: "Item 2"}
	history.Add(item2)

	seen_items := history.Items()

	if !(item1 == seen_items[0]) {
		t.Error("History was not set and read correctly.")
	}

	if !(item2 == seen_items[1]) {
		t.Error("History was not set and read correctly.")
	}

}

func TestContains(t *testing.T) {
	history := rss.History{}
	item := rss.Item{Title: "Item"}
	if history.Contains(item) {
		t.Error("Should not see item in history")
	}

	history.Add(item)

	if !history.Contains(item) {
		t.Error("Should see item in history")
	}
}
