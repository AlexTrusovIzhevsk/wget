package pathmapper

import (
	"net/url"
	"path"
	"strings"
)

type PathMapper interface {
	Map(url string) string
}

type FilePathMapper struct{}

func (m *FilePathMapper) Map(current string) string {
	u, err := url.Parse(current)
	if err != nil {
		return "default.html"
	}

	pth := u.Path
	if pth == "" || pth == "/" {
		return "index.html"
	}

	if pth[0] == '/' {
		pth = pth[1:]
	}

	if strings.HasSuffix(pth, "/") {
		pth = strings.TrimSuffix(pth, "/")
		return path.Join(pth, "index.html")
	}

	if path.Ext(pth) != "" {
		return pth
	}

	return path.Join(pth, "index.html")
}
