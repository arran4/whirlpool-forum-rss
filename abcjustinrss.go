package abcjustinrss

import (
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// RSS defines the structure of the RSS feed.
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

// Channel represents the RSS channel.
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

// Item represents an RSS feed item.
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

const BaseURL = "https://www.abc.net.au"
const JustInURL = BaseURL + "/news/justin"

// IsArticleLink filters valid article links.
func IsArticleLink(link string) bool {
	return strings.HasPrefix(link, "/news/") && strings.Contains(link, "/20") // Filter links with "/news/" and a year
}

// DeducePublicationDate converts relative times like "15 minutes ago" into absolute timestamps.
func DeducePublicationDate(relativeTime string) string {
	now := time.Now()
	re := regexp.MustCompile(`(\d+)\s+(minute|hour|day|week|month)s?\s+ago`)
	matches := re.FindStringSubmatch(relativeTime)

	if len(matches) < 3 {
		// If the relative time format is not recognized, fallback to now.
		return now.Format(time.RFC1123)
	}

	quantity := matches[1]
	unit := matches[2]

	// Parse the quantity as an integer
	value, err := time.ParseDuration(fmt.Sprintf("%s%s", quantity, unitToDuration(unit)))
	if err != nil {
		log.Printf("Failed to parse relative time '%s': %v", relativeTime, err)
		return now.Format(time.RFC1123)
	}

	// Deduce the publication time
	pubTime := now.Add(-value)
	return pubTime.Format(time.RFC1123)
}

// unitToDuration maps time units (e.g., "minute", "hour") to time.Duration format.
func unitToDuration(unit string) string {
	switch unit {
	case "minute":
		return "m"
	case "hour":
		return "h"
	case "day":
		return "24h"
	case "week":
		return "168h" // 7 days * 24 hours
	case "month":
		return "720h" // Approximation: 30 days * 24 hours
	default:
		return "0s" // Default to zero seconds for unrecognized units
	}
}

func FetchAndParseNewsToRSS() (error, RSS) {
	resp, err := http.Get(JustInURL)
	if err != nil {
		return fmt.Errorf("fetching news to rss: %v", err), RSS{}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %v", resp.Status), RSS{}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("parsing news to rss: %v", err), RSS{}
	}

	var items []Item

	baseUrl, err := url.Parse(JustInURL)
	if err != nil {
		return fmt.Errorf("parsing news to rss: %v", err), RSS{}
	}

	// Find articles by inspecting the HTML structure
	for i, s := range doc.Find("article").EachIter() {
		// Extract title and link
		title := strings.TrimSpace(s.Find("a").Text())
		description := strings.TrimSpace(s.Find("p").Text())
		link, exists := s.Find("a").Attr("href")
		if !exists || !IsArticleLink(link) {
			continue
		}

		// Extract relative time (e.g., "15 minutes ago")
		relativeTime := strings.TrimSpace(s.Find("time").Text())
		pubDate := DeducePublicationDate(relativeTime)

		linkUrl, err := baseUrl.Parse(link)
		if err != nil {
			return fmt.Errorf("parsing news to rss: %d %s: %v", i, s.Text(), err), RSS{}
		}

		// Add article to RSS items
		items = append(items, Item{
			Title:       title,
			Link:        linkUrl.String(),
			Description: description,
			PubDate:     pubDate,
			GUID:        linkUrl.String(),
		})
	}

	// Create RSS feed
	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       "ABC News - Just In",
			Link:        JustInURL,
			Description: "Latest news from ABC's Just In section",
			Items:       items,
		},
	}
	return err, rss
}
