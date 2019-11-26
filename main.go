package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
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

func main() {
	source := flag.String("s", "", "A HTTP RSS source.")
	pattern := flag.String("p", "", "The pattern to match RSS feed titles against.")
	outPath := flag.String("o", "./", "Output path. Defaults to current dir.")
	flag.Parse()

	if *source == "" || *pattern == "" {
		fmt.Println("A source and a pattern is required see -h for more info.")
		return
	}

	// Debug Read from file
	// dat, err := ioutil.ReadFile("sample.xml")
	// check(err)

	resp, err := http.Get(*source)
	check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	r := &feed{}
	err = xml.Unmarshal(body, &r)
	check(err)

	for _, torrent := range r.Channels {
		matched, err := regexp.MatchString(*pattern, torrent.Title)
		check(err)
		if matched {
			resp, err := http.Get(torrent.Link)
			check(err)
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			check(err)
			err = ioutil.WriteFile(
				*outPath+strings.ReplaceAll(torrent.Title, "/", "-")+".torrent",
				body,
				0644)
			check(err)
			break

		}
	}
}
