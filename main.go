package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"wget/downloader"
	"wget/parser"
	"wget/pathmapper"
	"wget/storage"
	"wget/webcrawler"
)

func main() {
	var (
		url    = flag.String("url", "", "URL to mirror")
		depth  = flag.Int("depth", 3, "Max depth for recursion")
		output = flag.String("output", "./mirror", "Output directory")
	)
	flag.Parse()

	if *url == "" {
		log.Fatal("URL is required")
	}

	downloader := &downloader.HTTPDownloader{} // реализация будет ниже
	parser := &parser.HtmlParser{}
	pathMapper := &pathmapper.FilePathMapper{}
	saver := &storage.OsFileSaver{OutputDir: *output}

	settings := webcrawler.WebCrawlerSettings{
		MaxDepth:   *depth,
		MaxWorkers: 0, // not implemented
	}

	crawler := webcrawler.NewWebCrawler(downloader, parser, pathMapper, saver, settings)

	result, err := crawler.Mirror(context.Background(), *url)
	if err != nil {
		log.Fatalf("Mirror failed: %v", err)
	}

	fmt.Printf("Success: %d, Errors: %d\n", result.CountSuccess, result.CountError)
}
