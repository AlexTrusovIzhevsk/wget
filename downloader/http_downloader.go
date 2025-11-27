package downloader

import (
	"context"
	"io"
	"net/http"
)

type Downloader interface {
	Download(ctx context.Context, url string) (data []byte, err error)
}

type HTTPDownloader struct{}

func (d *HTTPDownloader) Download(ctx context.Context, url string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, &DownloadError{URL: url, StatusCode: response.StatusCode}
	}

	return io.ReadAll(response.Body)
}

type DownloadError struct {
	URL        string
	StatusCode int
}

func (e *DownloadError) Error() string {
	return e.URL + ": " + http.StatusText(e.StatusCode)
}
