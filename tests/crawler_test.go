package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"wget/webcrawler"
)

func TestWebCrawler_Mirror_DownloadsAndSaves(t *testing.T) {
	// Подготовка
	mockDownloader := NewMockDownloader(map[string][]byte{
		"https://example.com": []byte("<html>test</html>"),
	}, nil)

	mockParser := NewMockHTMLParser([]string{}, []string{}, nil)
	mockPathMapper := NewMockPathMapper(map[string]string{
		"https://example.com": "index.html",
	})
	mockSaver := NewMockFileSaver(nil)

	settings := webcrawler.WebCrawlerSettings{
		MaxDepth:   1,
		MaxWorkers: 1,
	}

	crawler := webcrawler.NewWebCrawler(
		mockDownloader,
		mockParser,
		mockPathMapper,
		mockSaver,
		settings,
	)

	// Вызов
	result, err := crawler.Mirror(context.Background(), "https://example.com")

	// Проверки
	if err != nil {
		t.Fatalf("Mirror returned an error: %v", err)
	}

	// Проверяем, что был вызов Download с нужным URL
	if !mockDownloader.WasCalledWith("https://example.com") {
		t.Fatalf("Expected Download to be called with 'https://example.com', but it wasn't")
	}

	// Проверяем, что файл был "сохранён"
	saved := mockSaver.GetSaved()
	if _, ok := saved["index.html"]; !ok {
		t.Fatalf("Expected file 'index.html' to be saved, but it wasn't")
	}

	// Проверяем, что CountSuccess увеличился
	if result.CountSuccess != 1 {
		t.Errorf("Expected CountSuccess = 1, got %d", result.CountSuccess)
	}

	// Проверяем, что CountError = 0
	if result.CountError != 0 {
		t.Errorf("Expected CountError = 0, got %d", result.CountError)
	}
}

func TestWebCrawler_Mirror_DownloadError(t *testing.T) {
	// Подготовка
	mockDownloader := NewMockDownloader(nil, assert.AnError) // возвращает ошибку

	mockParser := NewMockHTMLParser([]string{}, []string{}, nil)
	mockPathMapper := NewMockPathMapper(map[string]string{})
	mockSaver := NewMockFileSaver(nil)

	settings := webcrawler.WebCrawlerSettings{
		MaxDepth:   1,
		MaxWorkers: 1,
	}

	crawler := webcrawler.NewWebCrawler(
		mockDownloader,
		mockParser,
		mockPathMapper,
		mockSaver,
		settings,
	)

	// Вызов
	result, err := crawler.Mirror(context.Background(), "https://example.com/bad")

	// Проверки
	if err != nil {
		t.Fatalf("Mirror should not return error on download failure, got: %v", err)
	}

	// Проверяем, что был вызов Download
	if !mockDownloader.WasCalledWith("https://example.com/bad") {
		t.Fatalf("Expected Download to be called with 'https://example.com/bad', but it wasn't")
	}

	// Проверяем, что файл **не** был сохранён
	saved := mockSaver.GetSaved()
	if len(saved) != 0 {
		t.Fatalf("Expected no files to be saved, but got %d", len(saved))
	}

	// Проверяем, что CountError увеличился
	if result.CountError != 1 {
		t.Errorf("Expected CountError = 1, got %d", result.CountError)
	}

	// Проверяем, что CountSuccess = 0
	if result.CountSuccess != 0 {
		t.Errorf("Expected CountSuccess = 0, got %d", result.CountSuccess)
	}
}

func TestWebCrawler_Mirror_DownloadsAbsoluteResources(t *testing.T) {
	// Подготовка
	mockDownloader := NewMockDownloader(map[string][]byte{
		"https://example.com":          []byte("<html><img src='https://example.com/logo.png'></html>"),
		"https://example.com/logo.png": []byte("PNG data"),
	}, nil)

	mockParser := NewMockHTMLParser([]string{"https://example.com/logo.png"}, []string{}, nil)
	mockPathMapper := NewMockPathMapper(map[string]string{
		"https://example.com":          "index.html",
		"https://example.com/logo.png": "logo.png",
	})
	mockSaver := NewMockFileSaver(nil)

	settings := webcrawler.WebCrawlerSettings{
		MaxDepth:   1,
		MaxWorkers: 1,
	}

	crawler := webcrawler.NewWebCrawler(
		mockDownloader,
		mockParser,
		mockPathMapper,
		mockSaver,
		settings,
	)

	// Вызов
	result, err := crawler.Mirror(context.Background(), "https://example.com")

	// Проверки
	if err != nil {
		t.Fatalf("Mirror returned an error: %v", err)
	}

	// Проверяем, что были вызовы Download для основной страницы и ресурса
	if !mockDownloader.WasCalledWith("https://example.com") {
		t.Fatalf("Expected Download to be called with 'https://example.com'")
	}
	if !mockDownloader.WasCalledWith("https://example.com/logo.png") {
		t.Fatalf("Expected Download to be called with 'https://example.com/logo.png'")
	}

	// Проверяем, что оба файла были сохранены
	saved := mockSaver.GetSaved()
	if _, ok := saved["index.html"]; !ok {
		t.Fatalf("Expected file 'index.html' to be saved")
	}
	if _, ok := saved["logo.png"]; !ok {
		t.Fatalf("Expected file 'logo.png' to be saved")
	}

	// Проверяем счётчики
	if result.CountSuccess != 2 { // 1 страница + 1 ресурс
		t.Errorf("Expected CountSuccess = 2, got %d", result.CountSuccess)
	}
	if result.CountError != 0 {
		t.Errorf("Expected CountError = 0, got %d", result.CountError)
	}
}

func TestWebCrawler_Mirror_NormalizesRelativeResourceURLs(t *testing.T) {
	// Подготовка
	mockDownloader := NewMockDownloader(map[string][]byte{
		"https://example.com/dir/":     []byte("<html><img src='/logo.png'></html>"),
		"https://example.com/logo.png": []byte("PNG data"),
	}, nil)

	mockParser := NewMockHTMLParser([]string{"/logo.png"}, []string{}, nil)
	mockPathMapper := NewMockPathMapper(map[string]string{
		"https://example.com/dir/":     "dir/index.html",
		"https://example.com/logo.png": "logo.png",
	})
	mockSaver := NewMockFileSaver(nil)

	settings := webcrawler.WebCrawlerSettings{
		MaxDepth:   1,
		MaxWorkers: 1,
	}

	crawler := webcrawler.NewWebCrawler(
		mockDownloader,
		mockParser,
		mockPathMapper,
		mockSaver,
		settings,
	)

	// Вызов
	result, err := crawler.Mirror(context.Background(), "https://example.com/dir/")

	// Проверки
	if err != nil {
		t.Fatalf("Mirror returned an error: %v", err)
	}

	// Проверяем, что были вызовы Download для основной страницы и ресурса
	if !mockDownloader.WasCalledWith("https://example.com/dir/") {
		t.Fatalf("Expected Download to be called with 'https://example.com/dir/'")
	}
	if !mockDownloader.WasCalledWith("https://example.com/logo.png") {
		t.Fatalf("Expected Download to be called with 'https://example.com/logo.png'")
	}

	// Проверяем, что оба файла были сохранены
	saved := mockSaver.GetSaved()
	if _, ok := saved["dir/index.html"]; !ok {
		t.Fatalf("Expected file 'dir/index.html' to be saved")
	}
	if _, ok := saved["logo.png"]; !ok {
		t.Fatalf("Expected file 'logo.png' to be saved")
	}

	// Проверяем счётчики
	if result.CountSuccess != 2 { // 1 страница + 1 ресурс
		t.Errorf("Expected CountSuccess = 2, got %d", result.CountSuccess)
	}
	if result.CountError != 0 {
		t.Errorf("Expected CountError = 0, got %d", result.CountError)
	}
}

func TestWebCrawler_Mirror_RecursesLinksUnlimited(t *testing.T) {
	html1 := "<html><a href='/page2'>Page 2</a></html>"
	html2 := "<html><a href='/page3'>Page 3</a></html>"
	html3 := "<html>End</html>"

	// Подготовка
	mockDownloader := NewMockDownloader(map[string][]byte{
		"https://example.com":       []byte(html1),
		"https://example.com/page2": []byte(html2),
		"https://example.com/page3": []byte(html3),
	}, nil)

	mockParser := &MockParserWithDynamicLinks{
		linksMap: map[string][]string{
			html1: {"/page2"},
			html2: {"/page3"},
		},
	}

	mockPathMapper := NewMockPathMapper(map[string]string{
		"https://example.com":       "index.html",
		"https://example.com/page2": "page2/index.html",
		"https://example.com/page3": "page3/index.html",
	})
	mockSaver := NewMockFileSaver(nil)

	settings := webcrawler.WebCrawlerSettings{
		MaxDepth:   10, // Достаточно большой лимит
		MaxWorkers: 1,
	}

	crawler := webcrawler.NewWebCrawler(
		mockDownloader,
		mockParser,
		mockPathMapper,
		mockSaver,
		settings,
	)

	// Вызов
	result, err := crawler.Mirror(context.Background(), "https://example.com")

	// Проверки
	if err != nil {
		t.Fatalf("Mirror returned an error: %v", err)
	}

	// Проверяем, что были вызовы Download для всех страниц
	if !mockDownloader.WasCalledWith("https://example.com") {
		t.Fatalf("Expected Download to be called with 'https://example.com'")
	}
	if !mockDownloader.WasCalledWith("https://example.com/page2") {
		t.Fatalf("Expected Download to be called with 'https://example.com/page2'")
	}
	if !mockDownloader.WasCalledWith("https://example.com/page3") {
		t.Fatalf("Expected Download to be called with 'https://example.com/page3'")
	}

	// Проверяем, что все файлы сохранены
	saved := mockSaver.GetSaved()
	if _, ok := saved["index.html"]; !ok {
		t.Fatalf("Expected file 'index.html' to be saved")
	}
	if _, ok := saved["page2/index.html"]; !ok {
		t.Fatalf("Expected file 'page2/index.html' to be saved")
	}
	if _, ok := saved["page3/index.html"]; !ok {
		t.Fatalf("Expected file 'page3/index.html' to be saved")
	}

	// Проверяем счётчики
	if result.CountSuccess != 3 { // 3 страницы
		t.Errorf("Expected CountSuccess = 3, got %d", result.CountSuccess)
	}
	if result.CountError != 0 {
		t.Errorf("Expected CountError = 0, got %d", result.CountError)
	}
}

func TestWebCrawler_Mirror_RespectsMaxDepth(t *testing.T) {
	html1 := "<html><a href='/page2'>Page 2</a></html>"
	html2 := "<html><a href='/page3'>Page 3</a></html>"
	html3 := "<html>End</html>"

	// Подготовка
	mockDownloader := NewMockDownloader(map[string][]byte{
		"https://example.com":       []byte(html1),
		"https://example.com/page2": []byte(html2),
		"https://example.com/page3": []byte(html3),
	}, nil)

	mockParser := &MockParserWithDynamicLinks{
		linksMap: map[string][]string{
			html1: {"/page2"},
			html2: {"/page3"},
		},
	}

	mockPathMapper := NewMockPathMapper(map[string]string{
		"https://example.com":       "index.html",
		"https://example.com/page2": "page2/index.html",
		"https://example.com/page3": "page3/index.html",
	})
	mockSaver := NewMockFileSaver(nil)

	settings := webcrawler.WebCrawlerSettings{
		MaxDepth:   2, // Должен зайти на /page2, но не на /page3
		MaxWorkers: 1,
	}

	crawler := webcrawler.NewWebCrawler(
		mockDownloader,
		mockParser,
		mockPathMapper,
		mockSaver,
		settings,
	)

	// Вызов
	result, err := crawler.Mirror(context.Background(), "https://example.com")

	// Проверки
	if err != nil {
		t.Fatalf("Mirror returned an error: %v", err)
	}

	// Проверяем, что были вызовы Download только для 2 уровней (MaxDepth = 2)
	if !mockDownloader.WasCalledWith("https://example.com") {
		t.Fatalf("Expected Download to be called with 'https://example.com'")
	}
	if !mockDownloader.WasCalledWith("https://example.com/page2") {
		t.Fatalf("Expected Download to be called with 'https://example.com/page2'")
	}
	// MaxDepth = 2, значит page3 не должен быть скачан
	if mockDownloader.WasCalledWith("https://example.com/page3") {
		t.Fatalf("Download should not be called for 'https://example.com/page3' due to MaxDepth")
	}

	// Проверяем, что сохранены только index.html и page2/index.html
	saved := mockSaver.GetSaved()
	if _, ok := saved["index.html"]; !ok {
		t.Fatalf("Expected file 'index.html' to be saved")
	}
	if _, ok := saved["page2/index.html"]; !ok {
		t.Fatalf("Expected file 'page2/index.html' to be saved")
	}
	if _, ok := saved["page3/index.html"]; ok {
		t.Fatalf("File 'page3/index.html' should not be saved due to MaxDepth")
	}

	// Проверяем счётчики
	if result.CountSuccess != 2 { // 2 страницы
		t.Errorf("Expected CountSuccess = 2, got %d", result.CountSuccess)
	}
	if result.CountError != 0 {
		t.Errorf("Expected CountError = 0, got %d", result.CountError)
	}
}

func TestWebCrawler_Mirror_AvoidsDuplicates(t *testing.T) {
	html1 := "<html><a href='/page2'>Page 2</a><a href='/page3'>Page 3</a></html>"
	html2 := "<html><a href='/page3'>Page 3 again</a></html>" // page3 встречается снова
	html3 := "<html>End</html>"

	// Подготовка
	mockDownloader := NewMockDownloader(map[string][]byte{
		"https://example.com":       []byte(html1),
		"https://example.com/page2": []byte(html2),
		"https://example.com/page3": []byte(html3),
	}, nil)

	mockParser := &MockParserWithDynamicLinks{
		linksMap: map[string][]string{
			html1: {"/page2", "/page3"},
			html2: {"/page3"}, // page3 встречается снова
		},
	}

	mockPathMapper := NewMockPathMapper(map[string]string{
		"https://example.com":       "index.html",
		"https://example.com/page2": "page2/index.html",
		"https://example.com/page3": "page3/index.html",
	})
	mockSaver := NewMockFileSaver(nil)

	settings := webcrawler.WebCrawlerSettings{
		MaxDepth:   10, // Без ограничения
		MaxWorkers: 1,
	}

	crawler := webcrawler.NewWebCrawler(
		mockDownloader,
		mockParser,
		mockPathMapper,
		mockSaver,
		settings,
	)

	// Вызов
	result, err := crawler.Mirror(context.Background(), "https://example.com")

	// Проверки
	if err != nil {
		t.Fatalf("Mirror returned an error: %v", err)
	}

	// Проверяем, что каждый URL скачан **только один раз**
	if !mockDownloader.WasCalledWith("https://example.com") {
		t.Fatalf("Expected Download to be called with 'https://example.com'")
	}
	if !mockDownloader.WasCalledWith("https://example.com/page2") {
		t.Fatalf("Expected Download to be called with 'https://example.com/page2'")
	}
	if !mockDownloader.WasCalledWith("https://example.com/page3") {
		t.Fatalf("Expected Download to be called with 'https://example.com/page3'")
	}

	// Убедимся, что вызовов не было больше, чем нужно
	// В реальности можно проверить количество вызовов в CallLog
	// Но для простоты проверим, что вызовы были, и файлы сохранены
	saved := mockSaver.GetSaved()
	if _, ok := saved["index.html"]; !ok {
		t.Fatalf("Expected file 'index.html' to be saved")
	}
	if _, ok := saved["page2/index.html"]; !ok {
		t.Fatalf("Expected file 'page2/index.html' to be saved")
	}
	if _, ok := saved["page3/index.html"]; !ok {
		t.Fatalf("Expected file 'page3/index.html' to be saved")
	}

	// Проверяем счётчики
	if result.CountSuccess != 3 { // 3 уникальные страницы
		t.Errorf("Expected CountSuccess = 3, got %d", result.CountSuccess)
	}
	if result.CountError != 0 {
		t.Errorf("Expected CountError = 0, got %d", result.CountError)
	}

	// Проверим, что вызовов Download было ровно 3
	if len(mockDownloader.CallLog) != 3 {
		t.Fatalf("Expected 3 calls to Download, got %d", len(mockDownloader.CallLog))
	}
}

func TestWebCrawler_Mirror_RecoversFromLinkError(t *testing.T) {
	html1 := "<html><a href='/page2'>Page 2</a><a href='/page3'>Page 3</a></html>"
	html2 := "<html>OK</html>"

	// page3 возвращает ошибку
	mockDownloaderWithErrors := &MockDownloaderWithSomeErrors{
		responses: map[string][]byte{
			"https://example.com":       []byte(html1),
			"https://example.com/page2": []byte(html2),
		},
		errors: map[string]error{
			"https://example.com/page3": assert.AnError,
		},
	}

	mockParser := &MockParserWithDynamicLinks{
		linksMap: map[string][]string{
			html1: {"/page2", "/page3"},
			html2: {},
		},
	}

	mockPathMapper := NewMockPathMapper(map[string]string{
		"https://example.com":       "index.html",
		"https://example.com/page2": "page2/index.html",
		"https://example.com/page3": "page3/index.html",
	})
	mockSaver := NewMockFileSaver(nil)

	settings := webcrawler.WebCrawlerSettings{
		MaxDepth:   10,
		MaxWorkers: 1,
	}

	crawler := webcrawler.NewWebCrawler(
		mockDownloaderWithErrors,
		mockParser,
		mockPathMapper,
		mockSaver,
		settings,
	)

	// Вызов
	result, err := crawler.Mirror(context.Background(), "https://example.com")

	// Проверки
	if err != nil {
		t.Fatalf("Mirror should not return error on link failure, got: %v", err)
	}

	// Проверяем, что были вызовы для всех трёх URL
	if !mockDownloaderWithErrors.WasCalledWith("https://example.com") {
		t.Fatalf("Expected Download to be called with 'https://example.com'")
	}
	if !mockDownloaderWithErrors.WasCalledWith("https://example.com/page2") {
		t.Fatalf("Expected Download to be called with 'https://example.com/page2'")
	}
	if !mockDownloaderWithErrors.WasCalledWith("https://example.com/page3") {
		t.Fatalf("Expected Download to be called with 'https://example.com/page3'")
	}

	// Проверяем, что успешно сохранены только 2 файла (page3 не скачался)
	saved := mockSaver.GetSaved()
	if _, ok := saved["index.html"]; !ok {
		t.Fatalf("Expected file 'index.html' to be saved")
	}
	if _, ok := saved["page2/index.html"]; !ok {
		t.Fatalf("Expected file 'page2/index.html' to be saved")
	}
	if _, ok := saved["page3/index.html"]; ok {
		t.Fatalf("File 'page3/index.html' should not be saved due to error")
	}

	// Проверяем счётчики
	if result.CountSuccess != 2 { // 2 успешных (index + page2)
		t.Errorf("Expected CountSuccess = 2, got %d", result.CountSuccess)
	}
	if result.CountError != 1 { // 1 ошибка (page3)
		t.Errorf("Expected CountError = 1, got %d", result.CountError)
	}
}

func TestWebCrawler_Mirror_IgnoresExternalLinks(t *testing.T) {
	html1 := "<html><a href='https://external.com'>External</a><a href='/internal'>Internal</a></html>"
	html2 := "<html>Internal page</html>"

	// Подготовка
	mockDownloader := NewMockDownloader(map[string][]byte{
		"https://example.com":          []byte(html1),
		"https://example.com/internal": []byte(html2),
	}, nil)

	mockParser := &MockParserWithDynamicLinks{
		linksMap: map[string][]string{
			string([]byte(html1)): {"/internal", "https://external.com"},
		},
	}

	mockPathMapper := NewMockPathMapper(map[string]string{
		"https://example.com":          "index.html",
		"https://example.com/internal": "internal/index.html",
	})
	mockSaver := NewMockFileSaver(nil)

	settings := webcrawler.WebCrawlerSettings{
		MaxDepth:   10,
		MaxWorkers: 1,
	}

	crawler := webcrawler.NewWebCrawler(
		mockDownloader,
		mockParser,
		mockPathMapper,
		mockSaver,
		settings,
	)

	// Вызов
	result, err := crawler.Mirror(context.Background(), "https://example.com")

	// Проверки
	if err != nil {
		t.Fatalf("Mirror returned an error: %v", err)
	}

	// Проверяем, что был вызов для внутренней ссылки
	if !mockDownloader.WasCalledWith("https://example.com") {
		t.Fatalf("Expected Download to be called with 'https://example.com'")
	}
	if !mockDownloader.WasCalledWith("https://example.com/internal") {
		t.Fatalf("Expected Download to be called with 'https://example.com/internal'")
	}

	// Проверяем, что внешняя ссылка НЕ вызывалась
	if mockDownloader.WasCalledWith("https://external.com") {
		t.Fatalf("Download should not be called for 'https://external.com' (external link)")
	}

	// Проверяем, что сохранены только внутренние файлы
	saved := mockSaver.GetSaved()
	if _, ok := saved["index.html"]; !ok {
		t.Fatalf("Expected file 'index.html' to be saved")
	}
	if _, ok := saved["internal/index.html"]; !ok {
		t.Fatalf("Expected file 'internal/index.html' to be saved")
	}
	if _, ok := saved["external.com"]; ok {
		t.Fatalf("File for external link should not be saved")
	}

	// Проверяем счётчики
	if result.CountSuccess != 2 { // 2 внутренние страницы
		t.Errorf("Expected CountSuccess = 2, got %d", result.CountSuccess)
	}
	if result.CountError != 0 {
		t.Errorf("Expected CountError = 0, got %d", result.CountError)
	}
}

// not implemented
//func TestWebCrawler_Mirror_ReplacesAbsoluteResourceURLsInHTML(t *testing.T) {
//	originalHTML := `<html>
//<head>
//    <link rel="stylesheet" href="https://example.com/style.css">
//    <link rel="icon" href="https://example.com/favicon.ico">
//</head>
//<body>
//    <img src="https://example.com/logo.png" alt="Logo">
//</body>
//</html>`
//
//	// Ожидаемый HTML после замены URL на локальные
//	expectedHTML := `<html>
//<head>
//    <link rel="stylesheet" href="style.css">
//    <link rel="icon" href="favicon.ico">
//</head>
//<body>
//    <img src="logo.png" alt="Logo">
//</body>
//</html>`
//
//	// Подготовка
//	mockDownloader := NewMockDownloader(map[string][]byte{
//		"https://example.com":             []byte(originalHTML),
//		"https://example.com/style.css":   []byte("body { color: red; }"),
//		"https://example.com/favicon.ico": []byte("ico data"),
//		"https://example.com/logo.png":    []byte("png data"),
//	}, nil)
//
//	mockParser := NewMockHTMLParser([]string{"/style.css", "/favicon.ico", "/logo.png"}, []string{}, nil)
//	mockPathMapper := NewMockPathMapper(map[string]string{
//		"https://example.com":             "index.html",
//		"https://example.com/style.css":   "style.css",
//		"https://example.com/favicon.ico": "favicon.ico",
//		"https://example.com/logo.png":    "logo.png",
//	})
//	mockSaver := NewMockFileSaver(nil)
//
//	settings := webcrawler.WebCrawlerSettings{
//		MaxDepth:   1,
//		MaxWorkers: 1,
//	}
//
//	crawler := webcrawler.NewWebCrawler(
//		mockDownloader,
//		mockParser,
//		mockPathMapper,
//		mockSaver,
//		settings,
//	)
//
//	// Вызов
//	result, err := crawler.Mirror(context.Background(), "https://example.com")
//
//	// Проверки
//	if err != nil {
//		t.Fatalf("Mirror returned an error: %v", err)
//	}
//
//	// Проверяем, что файлы сохранены
//	saved := mockSaver.GetSaved()
//	if _, ok := saved["index.html"]; !ok {
//		t.Fatalf("Expected file 'index.html' to be saved")
//	}
//
//	// Проверяем, что в index.html были заменены URL
//	savedHTML := string(saved["index.html"])
//	if savedHTML != expectedHTML {
//		t.Errorf("Expected saved HTML to have local URLs, got:\n%s\nExpected:\n%s", savedHTML, expectedHTML)
//	}
//
//	// Проверяем счётчики
//	if result.CountSuccess != 4 { // 1 HTML + 3 ресурса
//		t.Errorf("Expected CountSuccess = 4, got %d", result.CountSuccess)
//	}
//	if result.CountError != 0 {
//		t.Errorf("Expected CountError = 0, got %d", result.CountError)
//	}
//}
