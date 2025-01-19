package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/arran4/abc-justin-rss"
	"io"
	"log"
	"os"
)

func main() {
	var out io.Writer = os.Stdout
	setOutputFile := func(s string) error {
		if oldOut, ok := out.(io.Closer); ok {
			err := oldOut.Close()
			if err != nil {
				return err
			}
		}
		var err error
		out, err = os.Create(s)
		if err != nil {
			return err
		}
		return nil
	}
	flag.Func("o", "Output file", setOutputFile)
	flag.Func("output", "Output file", setOutputFile)
	flag.Parse()
	err, rss := abcjustinrss.FetchAndParseNewsToRSS()
	if err != nil {
		log.Fatal("Failed to fetch and parse new rss: ", err)
	}

	// Output RSS feed
	output, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal RSS: %v", err)
	}

	_, err = fmt.Fprintf(out, "%s%s", xml.Header, output)
	if err != nil {
		log.Fatalf("Failed to format RSS: %v", err)
	}

	if oldOut, ok := out.(io.Closer); ok {
		err := oldOut.Close()
		if err != nil {
			log.Fatalf("Failed to close file: %v", err)
		}
	}
}
