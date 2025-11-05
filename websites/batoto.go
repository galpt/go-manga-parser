package websites

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"go-manga-parser/core"
)

type batoParser struct{}

func NewBatoToParser() core.Parser { return &batoParser{} }

func (b *batoParser) Name() string { return "Bato.to" }

func (b *batoParser) Parse(outputDir string, pool core.Worker) error {
	results := []core.Manga{}
	// fallback implementation using simple HTTP + goquery when colly shim is not used.
	u := url.URL{Scheme: "https", Host: "bato.to", Path: "/browse"}
	html, err := core.FetchRenderedHTML(u.String(), 15*time.Second)
	if err != nil || strings.TrimSpace(html) == "" {
		// try a quick HTTP GET parse via net/http and goquery.NewDocumentFromReader
		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Get(u.String())
		if err != nil {
			return fmt.Errorf("failed to fetch bato.to: %w", err)
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("unexpected status %d fetching bato.to: %s", resp.StatusCode, string(bodyBytes))
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to parse bato.to HTML: %w", err)
		}
		doc.Find("#series-list > div").Each(func(i int, s *goquery.Selection) {
			a := s.Find("a").First()
			href, _ := a.Attr("href")
			title := strings.TrimSpace(s.Find(".item-title").Text())
			if title == "" {
				title = strings.TrimSpace(a.Text())
			}
			m := core.Manga{
				ID:    path.Base(href),
				Title: title,
				URL:   href,
			}
			results = append(results, m)
		})
	} else {
		// parse HTML we got from chromedp
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			return fmt.Errorf("failed to parse rendered bato.to HTML: %w", err)
		}
		doc.Find("#series-list > div").Each(func(i int, s *goquery.Selection) {
			a := s.Find("a").First()
			href, _ := a.Attr("href")
			title := strings.TrimSpace(s.Find(".item-title").Text())
			if title == "" {
				title = strings.TrimSpace(a.Text())
			}
			m := core.Manga{
				ID:    path.Base(href),
				Title: title,
				URL:   href,
			}
			results = append(results, m)
		})
	}

	if err := core.WriteJSONAtomically(outputDir, "batoto.json", results); err != nil {
		return err
	}
	return nil
}

// (This file intentionally uses goquery for parsing.)
