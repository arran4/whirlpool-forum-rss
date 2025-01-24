package main

import (
	whirlpoolforumrss "github.com/arran4/whirlpool-forum-rss"
	"log"
	"net/http"
	"net/http/cgi"
)

func main() {
	log.Fatal(cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		action := r.URL.Query().Get("action")
		if action == "" {
			action = "newthreads"
		}

		rss, err := whirlpoolforumrss.GenerateRSS(action)
		if err != nil {
			http.Error(w, "Failed to generate RSS feed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/rss+xml")
		_, err = w.Write(rss)
		if err != nil {
			log.Printf("Error writing RSS: %v", err)
		}
	})))
}
