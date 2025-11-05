# Go Manga Parser

A modular, Go-based manga parser focused on maintainability and humane scraping defaults.

> [!NOTE]
> This repository is intended for learning, integration, and building reader apps. Do not use it to break site terms of service or violate copyright.

## Key ideas

- Use site-specific parsers under `websites/` (for example `websites/mangadex.go`).
- Prefer official JSON APIs where available; fall back to HTML scraping when needed.
- Keep concurrency bounded with the included worker pool to stay polite to sites.
- Provide an optional headless-browser fallback (`chromedp`) for JS-protected pages (Cloudflare challenges).

## What this produces

The parser writes one JSON file per site into an `output/` folder in the current working directory. Example:

- `output/mangadex.json`
- `output/batoto.json`

Each JSON file is an array of `Manga` objects. Typical fields:

- `id`, `title`, `description`, `url`, `public_url`, `cover_url`, `tags`, `state`, `authors`
- `chapters`: array of `{ id, title, number, url, pages }` where `pages` is an array of image URLs

## Quick start (PowerShell)

Open PowerShell in the repository root and run:

```powershell
cd "d:\coding\ultimate parser\go-manga-parser"
go mod tidy
go build ./cmd/parser -o manga-parser.exe
```

Run the parser (example):

```powershell
.
\manga-parser.exe -sites=mangadex,batoto -concurrency=6 -out=output
```

- `-sites` accepts a comma-separated list of site keys. Current keys: `mangadex`, `batoto`.
- `-concurrency` sets the worker pool size.
- `-out` sets the output directory (defaults to `output`).

## Notes about Cloudflare / JS-protected pages

> [!WARNING]
> Some sites present JavaScript anti-bot pages (Cloudflare, heavy JS, CAPTCHAs). The parser will attempt a headless rendering fallback via `chromedp` for pages where simple requests fail.

- Headless rendering is slower and uses more CPU/memory — use it sparingly.
- CAPTCHAs requiring human input cannot be solved automatically and will cause the parser to fail on that page.

## Extending the parser

1. Add a new file to `websites/`, e.g. `websites/my_site.go`.
2. Implement the `core.Parser` interface: `Name() string` and `Parse(outputDir string, pool core.Worker) error`.
3. Keep all selectors and site-specific logic inside that file so fixes are local and easy to recompile.

## Troubleshooting

> [!TIP]
> If you see frequent `403`/`429` responses, reduce concurrency, increase delays, and be sure to respect the site's robots and policies.

## License

MIT — see `LICENSE`.

