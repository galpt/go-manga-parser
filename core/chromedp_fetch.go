package core

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

// FetchRenderedHTML fetches a page using a headless Chrome instance and returns
// the full HTML after execution. This helps with sites that require JS (e.g.,
// Cloudflare JS challenges). Caller should keep calls limited to avoid high
// resource usage.
func FetchRenderedHTML(url string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cctx, ccancel := chromedp.NewContext(ctx)
	defer ccancel()

	var html string
	if err := chromedp.Run(cctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &html),
	); err != nil {
		return "", err
	}
	return html, nil
}
