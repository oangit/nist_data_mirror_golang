package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	linksURL := scapePage("https://nvd.nist.gov/vuln/data-feeds")
	for _, URL := range linksURL {
		fileName := filepath.Base(URL)
		fmt.Printf("Download %v from %v\n", fileName, URL)
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			if err := downloadFile(fileName, URL); err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("file %v exist \n", fileName)
		}
	}
}

func scapePage(url string) []string {

	linksURL := []string{}

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		foundURL := e.Request.AbsoluteURL(e.Attr("href"))
		if strings.Contains(foundURL, "https://nvd.nist.gov/feeds/json/cve/1.0") {
			linksURL = append(linksURL, e.Attr("href"))
		}
	})

	c.Visit(url)

	return linksURL
}

func downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}