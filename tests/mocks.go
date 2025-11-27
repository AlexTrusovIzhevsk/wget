package tests

import (
	"context"
	"errors"
)

// MockDownloader — заглушка для скачивания
type MockDownloader struct {
	responses map[string][]byte
	err       error
	CallLog   []string // лог вызовов
}

func NewMockDownloader(responses map[string][]byte, err error) *MockDownloader {
	return &MockDownloader{
		responses: responses,
		err:       err,
		CallLog:   []string{},
	}
}

func (m *MockDownloader) Download(ctx context.Context, url string) ([]byte, error) {
	m.CallLog = append(m.CallLog, url)
	if m.err != nil {
		return nil, m.err
	}
	data, ok := m.responses[url]
	if !ok {
		return nil, errors.New("URL not found in mock responses")
	}
	return data, nil
}

func (m *MockDownloader) WasCalledWith(url string) bool {
	for _, u := range m.CallLog {
		if u == url {
			return true
		}
	}
	return false
}

// MockHTMLParser — заглушка для парсинга HTML
type MockHTMLParser struct {
	resources []string
	links     []string
	err       error
}

func NewMockHTMLParser(resources, links []string, err error) *MockHTMLParser {
	return &MockHTMLParser{
		resources: resources,
		links:     links,
		err:       err,
	}
}

func (m *MockHTMLParser) ParseHTML(data []byte) (resources []string, links []string, err error) {
	return m.resources, m.links, m.err
}

// MockPathMapper — заглушка для генерации локальных путей
type MockPathMapper struct {
	paths map[string]string // URL -> path
}

func NewMockPathMapper(paths map[string]string) *MockPathMapper {
	return &MockPathMapper{paths: paths}
}

func (m *MockPathMapper) Map(url string) string {
	if path, ok := m.paths[url]; ok {
		return path
	}
	return "default.html"
}

// MockFileSaver — заглушка для сохранения файлов
type MockFileSaver struct {
	saved map[string][]byte
	err   error
}

func NewMockFileSaver(err error) *MockFileSaver {
	return &MockFileSaver{
		saved: make(map[string][]byte),
		err:   err,
	}
}

func (m *MockFileSaver) Save(path string, data []byte) error {
	if m.err != nil {
		return m.err
	}
	m.saved[path] = data
	return nil
}

func (m *MockFileSaver) GetSaved() map[string][]byte {
	return m.saved
}

// MockParserWithDynamicLinks — парсер, который возвращает разные ссылки для разных URL
type MockParserWithDynamicLinks struct {
	linksMap map[string][]string
}

func (m *MockParserWithDynamicLinks) ParseHTML(data []byte) (resources []string, links []string, err error) {
	content := string(data)
	for key, links := range m.linksMap {
		if content == key {
			return []string{}, links, nil
		}
	}
	return []string{}, []string{}, nil
}

// MockDownloaderWithSomeErrors — позволяет указать, какие URL возвращают ошибки
type MockDownloaderWithSomeErrors struct {
	responses map[string][]byte
	errors    map[string]error
	CallLog   []string
}

func (m *MockDownloaderWithSomeErrors) Download(ctx context.Context, url string) ([]byte, error) {
	m.CallLog = append(m.CallLog, url)

	if err, hasErr := m.errors[url]; hasErr {
		return nil, err
	}

	data, ok := m.responses[url]
	if !ok {
		return nil, errors.New("URL not found in mock responses")
	}

	return data, nil
}

func (m *MockDownloaderWithSomeErrors) WasCalledWith(url string) bool {
	for _, u := range m.CallLog {
		if u == url {
			return true
		}
	}
	return false
}
