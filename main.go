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
	"regexp"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
}
type feed struct {
	XMLName xml.Name `xml:"rss"`
	Items   []item   `xml:"channel>item"`
}

func readFeed(source string) *feed {
	resp, err := http.Get(source)
	check(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	r := &feed{}
	err = xml.Unmarshal(body, &r)
	check(err)
	return r
}

func CheckFolderForFile(folder string, filename string) bool {
	if folder == "" {
		folder= "."
	}
	files, err := ioutil.ReadDir(folder)
	check(err)
	for _, file := range files {
		if file.Name() == filename {
			fmt.Printf("File %s already exists\n", filename)
			return true
		}
	}
	return false
}

func torrentName(outPath string, torrent item) string {
	return outPath + strings.ReplaceAll(torrent.Title, "/", "-") + ".torrent"
}

func checkTorrent(pattern string, outPath string, torrent item) {
	matched, err := regexp.MatchString(pattern, torrent.Title)
	check(err)
	if !matched {
		return
	}

	if CheckFolderForFile(outPath, torrentName(outPath, torrent)) {
		return
	}

	resp, err := http.Get(torrent.Link)
	check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	err = ioutil.WriteFile(
		torrentName(outPath, torrent),
		body,
		0644)
	check(err)
	fmt.Printf("Found torrent %s\n", torrent.Title)
}

func diff(new, old *feed) []item {
	if old == nil {
		return new.Items
	}
	var newItems []item
	for _, item := range new.Items {
		if item == old.Items[0] {
			break
		}
		newItems = append(newItems, item)
	}
	return newItems
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

	var previousFeed *feed

	for {
		fmt.Println("Checking")
		feed := readFeed(*source)
		items := diff(feed, previousFeed)
		for _, torrent := range items {
			checkTorrent(*pattern, *outPath, torrent)
		}
		previousFeed = feed
		time.Sleep(time.Duration(*wait) * time.Second)
	}
}
