package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var STARTLINK string = "http://mateusz.at"

func GetPage(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func GetLinks(r io.Reader) ([]string, []string) {
	links := make([]string, 0)
	image_srcs := make([]string, 0)
	doc, err := html.Parse(r)
	if err != nil {
		fmt.Println(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
					break
				}
			}
		}
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key == "src" {
					image_srcs = append(image_srcs, a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links, image_srcs
}

func CreateDir(path string) error {
	fmt.Println("Creating new directory:", path)
	err := os.Mkdir(path, 0711)
	return err
}

func IsSameDomain(url string) bool {
	return strings.Contains(url, STARTLINK)
}

func RelativeToAbsoluteLink(url string) string {
	if len(url) == 0 {
		return url
	}
	if strings.Contains(url, "http") {
		return url
	}
	if url[0] == '/' {
		return STARTLINK + url
	} else {
		return STARTLINK + "/" + url
	}
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

func Crawl(startpage string) {
	fringe := make([]string, 0)
	fringe = append(fringe, startpage)
	visited := make(map[string]bool)
	downloaded := make(map[string]bool)

	// fmt.Println(fringe)
	for {
		if len(fringe) == 0 {
			break
		}
		nextElem := fringe[0]
		fringe = fringe[1:]
		if visited[nextElem] {
			continue
		}
		fmt.Println(len(fringe), nextElem)
		data, err := GetPage(nextElem)
		if err != nil {
			fmt.Println(err)
		}
		links, images := GetLinks(data)

		for _, image := range images {
			fullImage := RelativeToAbsoluteLink(image)
			if downloaded[fullImage] {
				continue
			}
			DownloadFromUrl(fullImage, "out")
			downloaded[fullImage] = true
		}

		for _, link := range links {
			fullLink := RelativeToAbsoluteLink(link)
			if IsSameDomain(fullLink) {
				fringe = append(fringe, fullLink)
			}
		}

		visited[nextElem] = true

	}
}

func DownloadFromUrl(url string, path string) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(path + "/" + fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()
	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	// fmt.Println(n, "bytes downloaded.")
}

func main() {
	fmt.Println("Starting GoScrape")
	startTime := time.Now()

	Crawl(STARTLINK)

	fmt.Println("Script finished. Time elapsed:", time.Since(startTime))
	fmt.Println("----------------------------------")

}
