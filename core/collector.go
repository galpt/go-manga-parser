package core

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
)

// NewCollector returns a Colly collector pre-configured to behave more human-like
// (custom User-Agent, randomized delays, and simple rate limiting).
func NewCollector() *colly.Collector {
	ua := randomUserAgent()
	c := colly.NewCollector(
		colly.UserAgent(ua),
		colly.AllowURLRevisit(),
		colly.Async(true),
	)

	// set a politeness rule: small random delay and bounded parallelism
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 3 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		// add some common headers
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		// rotate referer randomly between none and the target host
		if rand.Intn(3) == 0 {
			r.Headers.Set("Referer", r.URL.Scheme+"://"+r.URL.Host+"/")
		}
		// fallback to the default transport
		if r.Ctx.GetAny("http_client") == nil {
			r.Ctx.Put("http_client", http.DefaultClient)
		}
		fmt.Printf("Visiting %s\n", r.URL.String())
	})

	// Add a simple delay between requests to make requests look less bot-like.
	c.OnResponse(func(r *colly.Response) {
		// small random sleep to look human-ish
		time.Sleep(time.Duration(200+rand.Intn(400)) * time.Millisecond)
	})

	return c
}

func randomUserAgent() string {
	agents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36",
	}
	rand.Seed(time.Now().UnixNano())
	return agents[rand.Intn(len(agents))]
}
