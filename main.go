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

	log "github.com/sirupsen/logrus"

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

func logAndExitOnErr(message string, err error) {
	if err != nil {
		log.Fatal(message, " ", err)
	}
}

func readFeed(source *url.URL) *rss.Feed {
	resp, err := http.Get(source.String())
	logAndExitOnErr("Error when fetching feed:", err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	logAndExitOnErr("Error when reading feed body:", err)

	r := &rss.Feed{}
	err = xml.Unmarshal(body, &r)
	logAndExitOnErr("Error during RSS parsing:", err)

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
	logAndExitOnErr("error in pattern.", err)

	download_history := rss.History{}

	for {
		feedURI, err := url.ParseRequestURI(*options.source)
		if err != nil {
			log.Error("Unable to parse source URI: ", *options.source)
			break
		}
		feed := readFeed(feedURI)
		for _, torrent := range feed.Items {
			isMatching, _ := torrent.IsMatch(*options.pattern)
			if isMatching && !download_history.Contains(torrent) {
				err := getTorrent(torrent, options)
				if err != nil {
					log.Error("Unable to get torrent: ", err)
				} else {
					download_history.Add(torrent)
					log.Info("Found ", torrent.Title)
				}
			}
			if download_history.Length() == *options.downloadLimit {
				os.Exit(0)
			}
		}
		time.Sleep(time.Duration(*options.wait) * time.Second)
	}
}
