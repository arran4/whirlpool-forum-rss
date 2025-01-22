package abckohlerreport

import (
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/url"
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
const KohlerReportURL = BaseURL + "/news/programs/kohler-report"

// IsArticleLink filters valid article links.
func IsArticleLink(link string) bool {
	return strings.HasPrefix(link, "/news/") && strings.Contains(link, "/20") // Filter links with "/news/" and a year
}

func FetchAndParseToRSS() (error, RSS) {
	resp, err := http.Get(KohlerReportURL)
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

	baseUrl, err := url.Parse(KohlerReportURL)
	if err != nil {
		return fmt.Errorf("parsing news to rss: %v", err), RSS{}
	}

	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("a").Text())
		description := strings.TrimSpace(s.Find("p").Text())
		link, exists := s.Find("a").Attr("href")
		if !exists || !IsArticleLink(link) {
			return
		}

		datetime, exists := s.Find("time").Attr("datetime")
		if !exists {
			log.Printf("Missing datetime attribute in article %s", title)
			return
		}

		parsedTime, err := time.Parse(time.RFC3339, datetime)
		if err != nil {
			log.Printf("Failed to parse datetime '%s': %v", datetime, err)
			return
		}

		linkUrl, err := baseUrl.Parse(link)
		if err != nil {
			log.Printf("Failed to parse link '%s': %v", link, err)
			return
		}

		items = append(items, Item{
			Title:       title,
			Link:        linkUrl.String(),
			Description: description,
			PubDate:     parsedTime.Format(time.RFC1123),
			GUID:        linkUrl.String(),
		})
	})

	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       "ABC News - Kohler Report",
			Link:        KohlerReportURL,
			Description: "News from the world of finance and business.",
			Items:       items,
		},
	}
	return nil, rss
}
