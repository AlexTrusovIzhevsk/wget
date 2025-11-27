package webcrawler

import (
	"context"
	"net/url"
	"strings"
	"wget/downloader"
	"wget/parser"
	"wget/pathmapper"
	"wget/storage"
)

type WebCrawler struct {
	Downloader downloader.Downloader
	Parser     parser.Parser
	PathMapper pathmapper.PathMapper
	FileSaver  storage.FileSaver
	Settings   WebCrawlerSettings
}

type WebCrawlerSettings struct {
	MaxDepth   int
	MaxWorkers int
}

type WebCrawlerResult struct {
	CountSuccess int
	CountError   int
}

func NewWebCrawler(
	downloader downloader.Downloader,
	parser parser.Parser,
	pathMapper pathmapper.PathMapper,
	fileSaver storage.FileSaver,
	settings WebCrawlerSettings,
) *WebCrawler {
	return &WebCrawler{
		Downloader: downloader,
		Parser:     parser,
		PathMapper: pathMapper,
		FileSaver:  fileSaver,
		Settings:   settings,
	}
}

func (c *WebCrawler) Mirror(ctx context.Context, url string) (*WebCrawlerResult, error) {
	result := &WebCrawlerResult{}
	return result, c.mirror(ctx, url, url, 1, map[string]bool{}, result)
}

func (c *WebCrawler) mirror(ctx context.Context, baseUrl, url string, depth int, processed map[string]bool, result *WebCrawlerResult) error {
	data, err := c.download(ctx, url)
	c.check(result, err)
	processed[url] = true

	resources, links, err := c.Parser.ParseHTML(data)

	for _, resource := range resources {
		currentUrl := c.normalizeUrl(url, resource)
		if !processed[currentUrl] {
			processed[currentUrl] = true
			_, err := c.download(ctx, currentUrl)
			c.check(result, err)
		}
	}

	for _, link := range links {
		currentUrl := c.normalizeUrl(url, link)
		if depth < c.Settings.MaxDepth && !processed[currentUrl] && strings.HasPrefix(currentUrl, baseUrl) {
			err = c.mirror(ctx, baseUrl, currentUrl, depth+1, processed, result)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *WebCrawler) normalizeUrl(baseUrl, currentUrl string) string {
	base, err := url.Parse(baseUrl)
	if err != nil {
		return currentUrl
	}
	current, err := url.Parse(currentUrl)
	if err != nil {
		return currentUrl
	}
	result := base.ResolveReference(current)

	return result.String()
}

func (c *WebCrawler) download(ctx context.Context, url string) ([]byte, error) {
	data, err := c.Downloader.Download(ctx, url)
	if err != nil {
		return nil, err
	}

	err = c.saveData(url, data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (c *WebCrawler) check(result *WebCrawlerResult, err error) {
	if err != nil {
		result.CountError++
	} else {
		result.CountSuccess++
	}
}

func (c *WebCrawler) saveData(url string, data []byte) error {
	path := c.PathMapper.Map(url)
	err := c.FileSaver.Save(path, data)
	return err
}
