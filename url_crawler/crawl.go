package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Enter a URL")
		os.Exit(1)
	}
	retrieve(args[0])
}

func retrieve(uri string) {
	resp, err := http.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	links := collectLinks(resp.Body)
	fmt.Println(links)
}

func collectLinks(httpBody io.Reader) []string {
	links := []string{}
	collection := map[string]bool{}
	page := html.NewTokenizer(httpBody)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			return links
		}
		token := page.Token()
		if tokenType == html.StartTagToken && token.DataAtom.String() == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					tl := trimRightFromHash(attr.Val)
					if collection[tl] {
						continue
					}
					links = append(links, tl)
				}
			}
		}
	}
}

func trimRightFromHash(uri string) string {
	if id := strings.Index(uri, "#"); id != -1 {
		return uri[:id]
	}
	return uri
}
