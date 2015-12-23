package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/codegangsta/cli"
)

type Bookmark struct {
	Title         string `xml:"title"`
	URL           string `xml:"link"`
	User          string `xml:"creator"`
	Bookmarkcount int    `xml:"bookmarkcount"`
}

type Entries struct {
	Entries []Bookmark `xml:"item"`
}

func main() {
	app := cli.NewApp()
	app.Name = "hbfav"
	app.Usage = "hbfav username"
	app.Action = func(c *cli.Context) {
		username := ""
		if len(c.Args()) > 0 {
			username = c.Args()[0]
		} else {
			fmt.Println(app.Usage)
			return
		}
		url := buildURL(username)
		feed := getRSSFeed(url)
		results := parseItems(feed)
		showResult(url, results)
	}
	app.Run(os.Args)
}

func buildURL(user string) string {
	return fmt.Sprintf("http://b.hatena.ne.jp/%s/favorite.rss", url.QueryEscape(user))
}

func getRSSFeed(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	feed, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	return feed
}

func parseItems(feed []byte) Entries {
	entries := Entries{}

	err := xml.Unmarshal(feed, &entries)
	if err != nil {
		fmt.Println(err)
	}

	return entries
}

func printRSSFeedURL(url string) {
	fmt.Printf("Feed URL: %s\n\n", url)
}

func showResult(url string, results Entries) {
	printRSSFeedURL(url)
	for _, entry := range results.Entries {
		fmt.Println(entry)
	}
}
