package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
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
	app.Version = "0.0.1"

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

func maxTitleWidth(entries Entries) int {
	width := 0
	for _, bookmark := range entries.Entries {
		count := runewidth.StringWidth(bookmark.Title)
		if count > width {
			width = count
		}
	}
	return width
}

func maxURLWidth(entries Entries) int {
	width := 0
	for _, bookmark := range entries.Entries {
		count := runewidth.StringWidth(bookmark.URL)
		if count > width {
			width = count
		}
	}
	return width
}

func maxUserWidth(entries Entries) int {
	width := 0
	for _, bookmark := range entries.Entries {
		count := runewidth.StringWidth(bookmark.User)
		if count > width {
			width = count
		}
	}
	return width
}

func maxBookmarkcountWidth(entries Entries) int {
	width := 0
	for _, bookmark := range entries.Entries {
		count := runewidth.StringWidth(string(bookmark.Bookmarkcount))
		if count > width {
			width = count
		}
	}
	return width
}

func showResult(url string, results Entries) {
	printRSSFeedURL(url)

	titleWidth := maxTitleWidth(results)
	titleFmt := fmt.Sprintf("%%-%ds", titleWidth)

	urlWidth := maxURLWidth(results)
	urlFmt := fmt.Sprintf("%%-%ds", urlWidth)

	userWidth := maxUserWidth(results)
	userFmt := fmt.Sprintf("%%-%ds", userWidth)

	bookmarkcountWidth := maxBookmarkcountWidth(results)
	bookmarkcountFmt := fmt.Sprintf("%%-%ds", bookmarkcountWidth)

	fmt.Fprintf(color.Output, " %s | %s | %s | %s | \n",
		color.BlueString(titleFmt, "Title"),
		color.CyanString(urlFmt, "URL"),
		fmt.Sprintf(userFmt, "User"),
		fmt.Sprintf(bookmarkcountFmt, "Bookmark"),
	)

	fmt.Println(strings.Repeat("-", titleWidth+urlWidth+userWidth+bookmarkcountWidth+18))

	for _, e := range results.Entries {
		fmt.Fprintf(color.Output, " %s | %s | %s | %s \n",
			color.GreenString(runewidth.FillRight(e.Title, titleWidth)),
			fmt.Sprintf(urlFmt, e.URL),
			fmt.Sprintf(userFmt, e.User),
			color.CyanString(strconv.Itoa(e.Bookmarkcount)),
		)
	}
}
