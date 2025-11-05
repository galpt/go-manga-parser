# Go Manga Parser

A compact, modular Go-based manga parser built for maintainability and polite scraping. It's designed to be easy to extend and to produce clean, portable JSON outputs that other apps can consume.

> [!NOTE]
> This repository is intended for learning, integration, and building reader apps. Do not use it to break site terms of service or violate copyright.

## Key ideas

- Keep site logic modular: each site has its own parser under `websites/` (for example `websites/mangadex.go`).
- Prefer official JSON APIs when available; fall back to HTML scraping only when necessary.
- Bound concurrency with the included worker pool so the parser behaves politely and predictably.
- Offer an optional headless-browser fallback (`chromedp`) for JavaScript-protected pages and challenges.

## Design goals

- App-agnostic output: the parser focuses on producing clean, well-structured JSON files so third-party reader apps (in Python, Rust, Java, Kotlin, etc.) only need to consume the output files — they don't need to know site-specific scraping details.
- Portable & automatable: the tool is lightweight and portable enough to run on a small server or in CI (for example GitHub Actions) for scheduled runs. This makes it practical to automate periodic parsing (even short intervals like every 5 minutes if desired), while keeping concurrency and politeness settings configurable.

## What this produces

The parser writes one JSON file per site into an `output/` folder in the current working directory. Example:

- `output/mangadex.json`
- `output/batoto.json`

Each JSON file is an array of `Manga` objects. Typical fields:

- `id`, `title`, `description`, `url`, `public_url`, `cover_url`, `tags`, `state`, `authors`
- `chapters`: array of `{ id, title, number, url, pages }` where `pages` is an array of image URLs

## Quick start

Open PowerShell in the repository root and run:

```
go mod tidy
go build ./cmd/parser -o manga-parser.exe
```

You can also use the `compile.bat` to quickly compile the source. It will produce `manga-parser.exe` in the repository root.

Run the parser (example):

```
./manga-parser.exe -sites=mangadex,batoto -concurrency=6 -out=output
```

- `-sites` accepts a comma-separated list of site keys. Current keys: `mangadex`, `batoto`.
- `-concurrency` sets the worker pool size.
- `-out` sets the output directory (defaults to `output`).

## Notes about Cloudflare / JS-protected pages

> [!WARNING]
> Some sites present JavaScript anti-bot pages (Cloudflare, heavy JS, CAPTCHAs). The parser will attempt a headless rendering fallback via `chromedp` for pages where simple requests fail.

- Headless rendering is slower and uses more CPU/memory — enable it only when necessary.
- CAPTCHAs requiring human input cannot be solved automatically and will cause the parser to fail on that page.

## Extending the parser

1. Add a new file to `websites/`, e.g. `websites/my_site.go`.
2. Implement the `core.Parser` interface: `Name() string` and `Parse(outputDir string, pool core.Worker) error`.
3. Keep all selectors and site-specific logic inside that file so fixes are localized and easy to recompile.

## Troubleshooting

> [!TIP]
> If you see frequent `403`/`429` responses, reduce concurrency, increase delays, and be sure to respect the site's robots and policies.

## License

MIT — see `LICENSE`.

