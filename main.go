// Package downer solves all your torrent rss feed needs
//
// Will keep track of the files that's already been seen during
// the programs runtime.
// Will not download already existing files.
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"github.com/rpkarlsson/downer/rss"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFeed(source *url.URL) *rss.Feed {
	resp, err := http.Get(source.String())
	check(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	r := &rss.Feed{}
	err = xml.Unmarshal(body, &r)
	check(err)
	return r
}

func downloadTorrent(pattern string, outPath string, torrent rss.Item) {
	resp, err := http.Get(torrent.Link)
	check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	err = ioutil.WriteFile(
		outPath+strings.ReplaceAll(torrent.Title, "/", "-")+".torrent",
		body,
		0644)
	check(err)
}
func main() {
	source := flag.String("s", "", "A HTTP RSS source.")
	pattern := flag.String("p", "", "The pattern to match RSS feed titles against.")
	outPath := flag.String("o", "", "Output path. Defaults to current dir.")
	wait := flag.Int("t", 60*15, "Time to sleep between checks in seconds. Defaults to 15 minutes")
	flag.Parse()

	if *source == "" || *pattern == "" {
		fmt.Println("A source and a pattern is required see -h for more info.")
		return
	}

	download_history := rss.History{}

	for {
		feedURI, err := url.ParseRequestURI(*source)
		if err != nil {
			fmt.Printf("Unable to parse %s as URL.\n", *source)
			break
		}
		feed := readFeed(feedURI)
		for _, torrent := range feed.Items {
			if torrent.IsMatch(*pattern) && !download_history.Contains(torrent) {
				downloadTorrent(*pattern, *outPath, torrent)
				download_history.Add(torrent)
				fmt.Printf("Found torrent %s\n", torrent.Title)
			}
		}
		time.Sleep(time.Duration(*wait) * time.Second)
	}
}
