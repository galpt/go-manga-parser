package websites

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"go-manga-parser/core"
)

type mangaDexParser struct{}

func NewMangaDexParser() core.Parser { return &mangaDexParser{} }

func (m *mangaDexParser) Name() string { return "MangaDex" }

// Minimal MangaDex parser using the public API. It fetches up to maxItems
// manga entries (100 per request) and writes a json array to outputDir/mangadex.json.
func (m *mangaDexParser) Parse(outputDir string, pool core.Worker) error {
	client := &http.Client{}
	limit := 100
	offset := 0
	all := []core.Manga{}
	for {
		u := url.URL{
			Scheme: "https",
			Host:   "api.mangadex.org",
			Path:   "/manga",
		}
		q := u.Query()
		q.Set("limit", strconv.Itoa(limit))
		q.Set("offset", strconv.Itoa(offset))
		u.RawQuery = q.Encode()

		resp, err := client.Get(u.String())
		if err != nil {
			return err
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}

		var parsed struct {
			Result string            `json:"result"`
			Data   []json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(body, &parsed); err != nil {
			return err
		}
		if len(parsed.Data) == 0 {
			break
		}

		for _, item := range parsed.Data {
			var entry struct {
				ID         string `json:"id"`
				Attributes struct {
					Title       map[string]string `json:"title"`
					Description map[string]string `json:"description"`
				} `json:"attributes"`
				Relationships []struct {
					Type       string                 `json:"type"`
					ID         string                 `json:"id"`
					Attributes map[string]interface{} `json:"attributes"`
				} `json:"relationships"`
			}
			if err := json.Unmarshal(item, &entry); err != nil {
				// skip malformed entries but continue
				continue
			}

			title := entry.Attributes.Title["en"]
			if title == "" {
				// pick any available
				for _, v := range entry.Attributes.Title {
					title = v
					break
				}
			}

			coverURL := ""
			// try to find cover_art relationship
			for _, rel := range entry.Relationships {
				if rel.Type == "cover_art" {
					if file, ok := rel.Attributes["fileName"].(string); ok {
						// mangadex uses covers on uploads.mangadex.org with /cover_id/file
						coverURL = fmt.Sprintf("https://uploads.mangadex.org/covers/%s/%s", entry.ID, path.Base(file))
					}
				}
			}

			mng := core.Manga{
				ID:          entry.ID,
				Title:       title,
				Description: entry.Attributes.Description["en"],
				URL:         entry.ID,
				PublicURL:   "https://mangadex.org/title/" + entry.ID,
				CoverURL:    coverURL,
			}
			all = append(all, mng)
		}

		// next page
		offset += limit
		// prevent runaway in this initial implementation (1000 items cap)
		if offset > 1000 {
			break
		}
	}

	// write
	if err := core.WriteJSONAtomically(outputDir, "mangadex.json", all); err != nil {
		return err
	}
	return nil
}
