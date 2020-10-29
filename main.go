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
	"os"
	"strings"
	"time"

	"github.com/rpkarlsson/downer/magnet"
	"github.com/rpkarlsson/downer/rss"
)

type cliOptions struct {
	downloadLimit *int
	outPath       *string
	pattern       *string
	source        *string
	wait          *int
}

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

func downloadTorrentFile(torrent rss.Item) []byte {
	resp, err := http.Get(torrent.Link)
	check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	return body
}

func getTorrent(torrent rss.Item, options cliOptions) {
	var body []byte
	if magnet.IsMagnetLink(torrent.Link) {
		body = []byte(magnet.MakeTorrentBody(torrent.Link))
	} else {
		body = downloadTorrentFile(torrent)
	}
	err := ioutil.WriteFile(
		*options.outPath+strings.ReplaceAll(torrent.Title, "/", "-")+".torrent",
		body,
		0644)
	check(err)
}

func parseOptions() cliOptions {
	options := cliOptions{
		downloadLimit: flag.Int("l", -1, "A limit to the amount of torrents to download"),
		outPath:       flag.String("o", "", "Output path. Defaults to current dir."),
		pattern:       flag.String("p", "", "The pattern to match RSS feed titles against."),
		source:        flag.String("s", "", "A HTTP RSS source."),
		wait:          flag.Int("t", 60*15, "Time to sleep between checks in seconds. Defaults to 15 minutes"),
	}

	flag.Parse()
	return options
}

func main() {
	options := parseOptions()
	if *options.source == "" || *options.pattern == "" {
		fmt.Println("A source and a pattern is required see -h for more info.")
		return
	}

	download_history := rss.History{}

	for {
		feedURI, err := url.ParseRequestURI(*options.source)
		if err != nil {
			fmt.Printf("Unable to parse %s as URL.\n", *options.source)
			break
		}
		feed := readFeed(feedURI)
		for _, torrent := range feed.Items {
			if torrent.IsMatch(*options.pattern) && !download_history.Contains(torrent) {
				getTorrent(torrent, options)
				download_history.Add(torrent)
				fmt.Printf("Found torrent %s\n", torrent.Title)
			}
			if download_history.Length() == *options.downloadLimit {
				os.Exit(0)
			}
		}
		time.Sleep(time.Duration(*options.wait) * time.Second)
	}
}
