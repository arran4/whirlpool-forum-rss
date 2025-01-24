package main

import (
	"flag"
	whirlpoolforumrss "github.com/arran4/whirlpool-forum-rss"
	"io"
	"log"
	"os"
)

func main() {
	var action string
	var outputFile string
	flag.StringVar(&action, "action", "newthreads", "Action to scrape (newthreads or popular_views)")
	flag.StringVar(&outputFile, "output", "", "Output file for the RSS feed")
	flag.Parse()

	var out io.Writer = os.Stdout
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			log.Fatalf("Error creating output file: %v", err)
		}
		defer file.Close()
		out = file
	}

	rss, err := whirlpoolforumrss.GenerateRSS(action)
	if err != nil {
		log.Fatalf("Error generating RSS feed: %v", err)
	}

	_, err = out.Write(rss)
	if err != nil {
		log.Fatalf("Error writing RSS feed: %v", err)
	}
}
