package main

import (
	"fmt"
	"log"
	"scraper/browserhandler"
)

func main() {
	var b browserhandler.BrowserDriver = browserhandler.NewBrowser(browserhandler.Playwright)
	b.Start()
	defer b.Close()
	b.NavigateAndWaitFor("https://news.ycombinator.com", "networkIdle")
	nodes, err := b.NodeFinder(".athing")
	if err != nil {
		log.Fatal(err)
	}
	nodes.All()
	nodes, err = nodes.NodeFinder("td.title > span > a")
	if err != nil {
		log.Fatal(err)
	}
	allContent, err := nodes.All()
	if err != nil {
		log.Fatal(err)
	}
	for index, node := range allContent {
		output, err := node.TextContent()
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("%d: %s\n", index+1, output)
		}
	}
	// b.TextList("td.title > span > a", &titles)
}
