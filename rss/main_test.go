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

func TestContainsSameTorrentURL(t *testing.T) {
	history := rss.History{}
	first := rss.Item{Title: "a", Link: "a"}
	first_with_different_title := rss.Item{Title: "b", Link: "a"}
	unseen := rss.Item{Title: "b", Link: "b"}
	history.Add(first)
	if !history.Contains(first_with_different_title) {
		t.Error("Should see a item in history")
	}
	if history.Contains(unseen) {
		t.Error("Should not see a item in history")
	}
}

func TestIsMatch (t *testing.T) {
	item := rss.Item{Title: "a long string"}
	if !item.IsMatch("ong") {
		t.Error("Should be a match")
	}
	if item.IsMatch("foo") {
		t.Error("Should not be a match")
	}
}
