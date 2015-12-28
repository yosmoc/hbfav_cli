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
	Bookmarks []Bookmark `xml:"item"`
}

func main() {
	app := cli.NewApp()
	app.Name = "hbfav"
	app.Usage = "hbfav username"
	app.Version = "0.0.1"

	var bookmarkNumber int
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "number, n",
			Value:       -1,
			Usage:       "the number of bookmarks",
			Destination: &bookmarkNumber,
		},
	}

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
		showResult(url, results, bookmarkNumber)
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

func maxTitleWidth(entries Entries, itemNumber int) int {
	width := 0
	for i := 0; i < itemNumber; i++ {
		count := runewidth.StringWidth(entries.Bookmarks[i].Title)
		if count > width {
			width = count
		}
	}
	return width
}

func maxURLWidth(entries Entries, itemNumber int) int {
	width := 0
	for i := 0; i < itemNumber; i++ {
		count := runewidth.StringWidth(entries.Bookmarks[i].URL)
		if count > width {
			width = count
		}
	}
	return width
}

func maxUserWidth(entries Entries, itemNumber int) int {
	width := 0
	for i := 0; i < itemNumber; i++ {
		count := runewidth.StringWidth(entries.Bookmarks[i].User)
		if count > width {
			width = count
		}
	}
	return width
}

func maxBookmarkcountWidth(entries Entries, itemNumber int) int {
	width := 0
	for i := 0; i < itemNumber; i++ {
		count := runewidth.StringWidth(string(entries.Bookmarks[i].Bookmarkcount))
		if count > width {
			width = count
		}
	}
	return width
}

func showResult(url string, results Entries, itemNumber int) {
	printRSSFeedURL(url)

	var n int
	if itemNumber == -1 || itemNumber >= len(results.Bookmarks) {
		n = len(results.Bookmarks)
	} else {
		n = itemNumber
	}

	titleWidth := maxTitleWidth(results, n)
	titleFmt := fmt.Sprintf("%%-%ds", titleWidth)

	urlWidth := maxURLWidth(results, n)
	urlFmt := fmt.Sprintf("%%-%ds", urlWidth)

	userWidth := maxUserWidth(results, n)
	userFmt := fmt.Sprintf("%%-%ds", userWidth)

	bookmarkcountWidth := maxBookmarkcountWidth(results, n)
	bookmarkcountFmt := fmt.Sprintf("%%-%ds", bookmarkcountWidth)

	fmt.Fprintf(color.Output, " %s | %s | %s | %s | \n",
		color.BlueString(titleFmt, "Title"),
		color.CyanString(urlFmt, "URL"),
		fmt.Sprintf(userFmt, "User"),
		fmt.Sprintf(bookmarkcountFmt, "Bookmark"),
	)

	fmt.Println(strings.Repeat("-", titleWidth+urlWidth+userWidth+bookmarkcountWidth+18))

	for i, e := range results.Bookmarks {
		if i > n-1 {
			return
		}
		fmt.Fprintf(color.Output, " %s | %s | %s | %s \n",
			color.GreenString(runewidth.FillRight(e.Title, titleWidth)),
			fmt.Sprintf(urlFmt, e.URL),
			fmt.Sprintf(userFmt, e.User),
			color.CyanString(strconv.Itoa(e.Bookmarkcount)),
		)
	}
}
