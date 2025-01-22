package main

import (
	"encoding/xml"
	"fmt"
	"github.com/arran4/abc-kohler-report-rss"
	"log"
	"net/http"
	"net/http/cgi"
)

func main() {
	log.Fatal(cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err, rss := abckohlerreport.FetchAndParseToRSS()
		if err != nil {
			http.Error(w, "Failed to fetch and parse RSS", http.StatusInternalServerError)
			return
		}

		output, err := xml.MarshalIndent(rss, "", "  ")
		if err != nil {
			http.Error(w, "Failed to marshal RSS", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = fmt.Fprintf(w, "%s%s", xml.Header, output)
	})))
}
