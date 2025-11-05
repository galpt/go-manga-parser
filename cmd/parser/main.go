package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"go-manga-parser/core"
	"go-manga-parser/websites"
	"go-manga-parser/worker"
)

func main() {
	sites := flag.String("sites", "mangadex,batoto", "comma-separated list of sites to parse")
	concurrency := flag.Int("concurrency", 8, "number of concurrent workers")
	out := flag.String("out", "output", "output directory")
	flag.Parse()

	if err := os.MkdirAll(*out, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create output dir: %v\n", err)
		os.Exit(1)
	}

	pool := worker.NewPool(*concurrency)
	defer pool.Stop()

	var wg sync.WaitGroup
	for _, s := range strings.Split(*sites, ",") {
		s = strings.TrimSpace(s)
		var p core.Parser
		switch strings.ToLower(s) {
		case "mangadex":
			p = websites.NewMangaDexParser()
		case "batoto":
			p = websites.NewBatoToParser()
		default:
			fmt.Fprintf(os.Stderr, "unknown site: %s\n", s)
			continue
		}

		wg.Add(1)
		go func(parser core.Parser) {
			defer wg.Done()
			fmt.Printf("Starting parser: %s\n", parser.Name())
			if err := parser.Parse(*out, pool); err != nil {
				fmt.Fprintf(os.Stderr, "parser %s failed: %v\n", parser.Name(), err)
			} else {
				fmt.Printf("parser %s finished\n", parser.Name())
			}
		}(p)
	}

	wg.Wait()
}
