package rss

import (
	"encoding/xml"
)

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Items   []Item   `xml:"channel>item"`
}

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type History struct {
	seen_items []Item
}

func (h *History) Add(item Item) {
	h.seen_items = append(h.seen_items, item)
}

func (h *History) Items() []Item {
	return h.seen_items
}

func (h *History) Contains(item Item) bool {
	for _, seen_item := range h.Items() {
		if seen_item == item {
			return true
		}
	}
	return false
}
