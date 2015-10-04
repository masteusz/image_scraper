package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"time"
)

func GetPage(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func GetLinks(r io.Reader) []string {
	links := make([]string, 0)
	doc, err := html.Parse(r)
	if err != nil {
		fmt.Println(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					// fmt.Println(a.Val)
					links = append(links, a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links
}

func CreateDir(path string) error {
	fmt.Println("Creating new directory:", path)
	err := os.Mkdir(path, os.ModeDir)
	return err
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func main() {
	fmt.Println("Starting GoScrape")
	startTime := time.Now()

	data, err := GetPage("http://mateusz.at")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(GetLinks(data))

	fmt.Println("Script finished. Time elapsed:", time.Since(startTime))
	fmt.Println("----------------------------------")

}
