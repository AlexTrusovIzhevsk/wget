package tests

import (
	"testing"
	"wget/parser"
)

// Вспомогательная функция для сравнения слайсов
func assertEqualSlices(t *testing.T, actual, expected []string) {
	if len(actual) != len(expected) {
		t.Fatalf("Expected %d items, got %d", len(expected), len(actual))
	}

	for i, v := range actual {
		if v != expected[i] {
			t.Fatalf("Expected %v at index %d, got %v", expected[i], i, v)
		}
	}
}

func TestHTMLParser_ParseHTML_ExtractsResourcesAndLinks(t *testing.T) {
	html := `
<html>
<head>
    <link rel="stylesheet" href="/assets/style.css">
    <script src="/js/app.js"></script>
</head>
<body>
    <img src="/images/logo.png" alt="Logo">
    <a href="/about">О нас</a>
    <a href="https://external.com">Внешний сайт</a>
    <a href="/contact">Контакты</a>
</body>
</html>`

	expectedResources := []string{"/assets/style.css", "/js/app.js", "/images/logo.png"}
	expectedLinks := []string{"/about", "https://external.com", "/contact"}

	p := &parser.HtmlParser{} // или как у тебя реализовано

	resources, links, err := p.ParseHTML([]byte(html))

	if err != nil {
		t.Fatalf("ParseHTML returned an error: %v", err)
	}

	// Проверяем, что ресурсы и ссылки совпадают
	assertEqualSlices(t, resources, expectedResources)
	assertEqualSlices(t, links, expectedLinks)
}

func TestHTMLParser_ParseHTML_EmptyHTML(t *testing.T) {
	p := &parser.HtmlParser{}

	resources, links, err := p.ParseHTML([]byte(""))

	if err != nil {
		t.Fatalf("ParseHTML returned an error: %v", err)
	}

	if len(resources) != 0 {
		t.Errorf("Expected 0 resources, got %d", len(resources))
	}
	if len(links) != 0 {
		t.Errorf("Expected 0 links, got %d", len(links))
	}
}

func TestHTMLParser_ParseHTML_NoLinksOrResources(t *testing.T) {
	html := "<html><body><h1>Hello</h1></body></html>"
	p := &parser.HtmlParser{}

	resources, links, err := p.ParseHTML([]byte(html))

	if err != nil {
		t.Fatalf("ParseHTML returned an error: %v", err)
	}

	if len(resources) != 0 {
		t.Errorf("Expected 0 resources, got %d", len(resources))
	}
	if len(links) != 0 {
		t.Errorf("Expected 0 links, got %d", len(links))
	}
}

func TestHTMLParser_ParseHTML_InvalidHTML(t *testing.T) {
	invalidHTML := "<html><body><h1>Unclosed tag</body></html>" // intentionally broken
	p := &parser.HtmlParser{}

	resources, links, err := p.ParseHTML([]byte(invalidHTML))

	if err != nil {
		t.Fatalf("ParseHTML should not return error on invalid HTML, got: %v", err)
	}

	// Should not panic, and may return some valid tags
	// For this test, we just check it doesn't crash
	_ = resources
	_ = links
}

func TestHTMLParser_ParseHTML_IgnoresIgnoredSchemesEverywhere(t *testing.T) {
	html := `<html>
		<head>
			<link href="javascript:alert(1)">
		</head>
		<body>
			<a href="javascript:alert(1)">JS</a>
			<a href="mailto:test@example.com">Email</a>
			<img src="javascript:alert(1)">
			<iframe src="mailto:test@example.com">
			</iframe>
			<video src="tel:+123456789"></video>
			<a href="/valid">Valid</a>
		</body>
	</html>`
	p := &parser.HtmlParser{}

	resources, links, err := p.ParseHTML([]byte(html))

	if err != nil {
		t.Fatalf("ParseHTML returned an error: %v", err)
	}

	// Должна быть только одна настоящая ссылка
	expectedLinks := []string{"/valid"}
	assertEqualSlices(t, links, expectedLinks)

	// Resources должны быть пустыми (все игнорируемые)
	if len(resources) != 0 {
		t.Errorf("Expected 0 resources, got %d", len(resources))
	}
}
