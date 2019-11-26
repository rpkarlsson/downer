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
	XMLName  xml.Name `xml:"rss"`
	Channels []item   `xml:"channel>item"`
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

func checkTorrent(pattern string, outPath string, torrent item) {
	matched, err := regexp.MatchString(pattern, torrent.Title)
	check(err)
	if !matched {
		return
	}

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

	fmt.Println("Found", torrent.Title)
}

func main() {
	source := flag.String("s", "", "A HTTP RSS source.")
	pattern := flag.String("p", "", "The pattern to match RSS feed titles against.")
	outPath := flag.String("o", "./", "Output path. Defaults to current dir.")
	wait := flag.Int("t", 60*15, "Time to sleep between checks in seconds. Defaults to 15 minutes")
	flag.Parse()

	if *source == "" || *pattern == "" {
		fmt.Println("A source and a pattern is required see -h for more info.")
		return
	}

	for {
		fmt.Println("Checking")
		feed := readFeed(*source)
		for _, torrent := range feed.Channels {
			checkTorrent(*pattern, *outPath, torrent)
		}

		time.Sleep(time.Duration(*wait) * time.Second)
	}
}
