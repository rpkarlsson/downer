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
	"regexp"
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

func readFeed(source *url.URL) *rss.Feed {
	resp, err := http.Get(source.String())
	if err != nil {
		fmt.Println("Error when fetching feed:", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error when reading feed body:", err)
		os.Exit(1)
	}

	r := &rss.Feed{}
	err = xml.Unmarshal(body, &r)
	if err != nil {
		fmt.Println("Error during RSS parsing:", err)
		os.Exit(1)
	}

	return r
}

func downloadTorrentFile(torrent rss.Item) ([]byte, error) {
	resp, err := http.Get(torrent.Link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getTorrent(torrent rss.Item, options cliOptions) error {
	var body []byte
	var err error
	if magnet.IsMagnetLink(torrent.Link) {
		body = []byte(magnet.MakeTorrentBody(torrent.Link))
	} else {
		body, err = downloadTorrentFile(torrent)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(
		*options.outPath+strings.ReplaceAll(torrent.Title, "/", "-")+".torrent",
		body,
		0644)

	if err != nil {
		return err
	}

	return nil
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

	_, err := regexp.Compile(*options.pattern)

	if err != nil {
		fmt.Println("Unable to compile pattern:", err)
		os.Exit(1)
	}

	download_history := rss.History{}

	for {
		feedURI, err := url.ParseRequestURI(*options.source)
		if err != nil {
			fmt.Printf("Unable to parse the source \"%s \" as URL.\n", *options.source)
			break
		}
		feed := readFeed(feedURI)
		for _, torrent := range feed.Items {
			isMatching, _ := torrent.IsMatch(*options.pattern)
			if isMatching && !download_history.Contains(torrent) {
				err := getTorrent(torrent, options)
				if err != nil {
					fmt.Printf("Unable to get torrent:", err)
				} else {
					download_history.Add(torrent)
					fmt.Printf("Found torrent %s\n", torrent.Title)
				}
			}
			if download_history.Length() == *options.downloadLimit {
				os.Exit(0)
			}
		}
		time.Sleep(time.Duration(*options.wait) * time.Second)
	}
}
