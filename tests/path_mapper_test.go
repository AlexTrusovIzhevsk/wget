package tests

import (
	"testing"
	"wget/pathmapper"
)

func TestPathMapper_Map_Root(t *testing.T) {
	p := &pathmapper.FilePathMapper{}

	path := p.Map("https://example.com/")

	if path != "index.html" {
		t.Errorf("Expected 'index.html', got '%s'", path)
	}
}

func TestPathMapper_Map_Page(t *testing.T) {
	p := &pathmapper.FilePathMapper{}

	path := p.Map("https://example.com/page")

	if path != "page/index.html" {
		t.Errorf("Expected 'page/index.html', got '%s'", path)
	}
}

func TestPathMapper_Map_FileWithExtension(t *testing.T) {
	p := &pathmapper.FilePathMapper{}

	path := p.Map("https://example.com/style.css")

	if path != "style.css" {
		t.Errorf("Expected 'style.css', got '%s'", path)
	}
}

func TestPathMapper_Map_PageWithQuery(t *testing.T) {
	p := &pathmapper.FilePathMapper{}

	path := p.Map("https://example.com/page?param=1")

	if path != "page/index.html" {
		t.Errorf("Expected 'page/index.html', got '%s'", path)
	}
}

func TestPathMapper_Map_PageWithFragment(t *testing.T) {
	p := &pathmapper.FilePathMapper{}

	path := p.Map("https://example.com/page#section")

	if path != "page/index.html" {
		t.Errorf("Expected 'page/index.html', got '%s'", path)
	}
}

func TestPathMapper_Map_Directory(t *testing.T) {
	p := &pathmapper.FilePathMapper{}

	path := p.Map("https://example.com/dir/")

	if path != "dir/index.html" {
		t.Errorf("Expected 'dir/index.html', got '%s'", path)
	}
}
