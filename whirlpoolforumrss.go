package whirlpoolforumrss

import (
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

func GenerateRSS(action string) ([]byte, error) {
	baseURL := "https://forums.whirlpool.net.au/forum/"
	url := fmt.Sprintf("%s?action=%s", baseURL, action)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forum page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse page content: %w", err)
	}

	rssTitle := "Whirlpool Forums - " + strings.Split(strings.TrimSpace(doc.Find("title").Text()), " - ")[0]
	items := []Item{}

	descriptionTemplate := `Section: {{.Section}}
Tag: {{.Tag}}
Replies: {{.Replies}}
First Post: {{.FirstPostAuthor}} ({{.FirstPostTime}})
Last Post: {{.LastPostAuthor}} ({{.LastPostTime}})`
	descTmpl, err := template.New("desc").Parse(descriptionTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse description template: %w", err)
	}

	doc.Find("#threads tbody tr.thread").Each(func(i int, s *goquery.Selection) {
		section := strings.TrimSpace(s.Find("tr.section").PrevAll().First().Find("a.title").Text())
		title := strings.TrimSpace(s.Find("a.title").Text())
		topicLink, _ := s.Find("a.title").Attr("href")
		archiveLink := strings.Replace(topicLink, "/thread/", "/archive/", 1)
		firstPostAuthor := strings.TrimSpace(s.Find("td.oldest a").First().Text())
		firstPostTime := strings.TrimSpace(s.Find("td.oldest").Contents().Last().Text())
		lastPostAuthor := strings.TrimSpace(s.Find("td.newest span a").Text())
		lastPostTime := strings.TrimSpace(s.Find("td.newest span").Contents().Last().Text())
		replies := strings.TrimSpace(s.Find("td.reps").Text())
		tag := strings.TrimSpace(s.Find("a.group").Text())

		descData := map[string]string{
			"Section":         section,
			"Tag":             tag,
			"Replies":         replies,
			"FirstPostAuthor": firstPostAuthor,
			"FirstPostTime":   firstPostTime,
			"LastPostAuthor":  lastPostAuthor,
			"LastPostTime":    lastPostTime,
		}
		var descriptionBuilder strings.Builder
		if err := descTmpl.Execute(&descriptionBuilder, descData); err != nil {
			return
		}

		pubDate, _ := time.Parse("2 Jan 2006 3:04 PM", lastPostTime)

		items = append(items, Item{
			Title:       fmt.Sprintf("[%s] %s", section, title),
			Link:        fmt.Sprintf("https://forums.whirlpool.net.au%s", archiveLink),
			Description: descriptionBuilder.String(),
			PubDate:     pubDate.Format(time.RFC1123),
			GUID:        fmt.Sprintf("https://forums.whirlpool.net.au%s", topicLink),
		})
	})

	rss := RSS{
		XMLName: xml.Name{Local: "rss"},
		Version: "2.0",
		Channel: Channel{
			Title:       rssTitle,
			Link:        baseURL,
			Description: fmt.Sprintf("RSS feed for %s", action),
			Items:       items,
		},
	}

	xmlData, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal RSS feed: %w", err)
	}

	return xmlData, nil
}
