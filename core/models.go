package core

// Basic cross-site models used by parsers and output writers.

type Manga struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	URL         string    `json:"url,omitempty"`
	PublicURL   string    `json:"public_url,omitempty"`
	CoverURL    string    `json:"cover_url,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	State       string    `json:"state,omitempty"`
	Authors     []string  `json:"authors,omitempty"`
	Chapters    []Chapter `json:"chapters,omitempty"`
}

type Chapter struct {
	ID     string `json:"id"`
	Title  string `json:"title,omitempty"`
	Number int    `json:"number,omitempty"`
	URL    string `json:"url,omitempty"`
	Pages  []Page `json:"pages,omitempty"`
}

type Page struct {
	Index int    `json:"index"`
	URL   string `json:"url"`
}

// Parser is implemented by each website parser.
type Parser interface {
	Name() string
	// Parse writes results into outputDir. It may use the provided worker pool for concurrency.
	Parse(outputDir string, pool Worker) error
}

// Worker is a small exported interface implemented by the worker pool. It
// intentionally only exposes Submit to limit coupling.
type Worker interface {
	Submit(func())
}
