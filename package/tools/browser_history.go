package tools

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

type BrowserHistory struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type GetRecentInput struct {
	Limit int `json:"limit", description:"The number of history to get"`
}
type GetByKeywordInput struct {
	Keyword string `json:"keyword", description:"The keyword to get history"`
	Limit   int    `json:"limit", description:"The number of history to get"`
}
type GetByDomainInput struct {
	Domain string `json:"domain", description:"The domain to get history"`
	Limit  int    `json:"limit", description:"The number of history to get"`
}
type BrowserHistoryTool interface {
	GetRecent(ctx context.Context, input GetRecentInput) ([]BrowserHistory, error)
	GetByDomain(ctx context.Context, input GetByDomainInput) ([]BrowserHistory, error)
	GetByKeyword(ctx context.Context, input GetByKeywordInput) ([]BrowserHistory, error)
}

type browserHistoryTool struct {
	historyFile string
}

func NewBrowserHistoryTool(chromeProfilePath string) (BrowserHistoryTool, error) {
	return &browserHistoryTool{
		historyFile: filepath.Join(chromeProfilePath, "History"),
	}, nil
}

// ------------ SQL ------------
const getRecentQuery = `
SELECT 
    urls.url,
    urls.title
FROM visits
JOIN urls ON visits.url = urls.id
ORDER BY visits.visit_time DESC
LIMIT ?;
`

const getByDomainQuery = `
SELECT 
    urls.url,
    urls.title
FROM visits
JOIN urls ON visits.url = urls.id
WHERE urls.url LIKE ?
ORDER BY visits.visit_time DESC
LIMIT ?;
`

// just match the keyword in the title, url
const getByKeywordQuery = `
SELECT 
    urls.url,
    urls.title
FROM visits
JOIN urls ON visits.url = urls.id
WHERE urls.url LIKE ? OR urls.title LIKE ?
ORDER BY visits.visit_time DESC
LIMIT ?;
`

// ------------ Core ------------
func copyFile(src string) (string, error) {
	tmp := src + "_tmp_copy"

	in, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("open src: %w", err)
	}
	defer in.Close()

	out, err := os.Create(tmp)
	if err != nil {
		return "", fmt.Errorf("create tmp: %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return "", fmt.Errorf("copy file: %w", err)
	}

	return tmp, nil
}

func openTempDB(historyFile string) (*sql.DB, string, error) {
	tmp, err := copyFile(historyFile)
	if err != nil {
		return nil, "", err
	}

	db, err := sql.Open("sqlite", tmp)
	if err != nil {
		os.Remove(tmp)
		return nil, "", err
	}

	return db, tmp, nil
}

// ------------ API ------------
func (t *browserHistoryTool) GetRecent(ctx context.Context, input GetRecentInput) ([]BrowserHistory, error) {
	db, tmp, err := openTempDB(t.historyFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	defer os.Remove(tmp)

	rows, err := db.QueryContext(ctx, getRecentQuery, input.Limit)
	if err != nil {
		return nil, fmt.Errorf("query recent history: %w", err)
	}
	defer rows.Close()

	var list []BrowserHistory
	for rows.Next() {
		var h BrowserHistory
		if err := rows.Scan(&h.URL, &h.Title); err != nil {
			return nil, err
		}
		list = append(list, h)
	}
	return list, nil
}

func (t *browserHistoryTool) GetByDomain(ctx context.Context, input GetByDomainInput) ([]BrowserHistory, error) {
	db, tmp, err := openTempDB(t.historyFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	defer os.Remove(tmp)

	domain := strings.TrimSpace(input.Domain)
	if !strings.Contains(domain, "%") {
		domain = "%" + domain + "%"
	}

	rows, err := db.QueryContext(ctx, getByDomainQuery, domain, input.Limit)
	if err != nil {
		return nil, fmt.Errorf("query history by domain: %w", err)
	}
	defer rows.Close()

	var list []BrowserHistory
	for rows.Next() {
		var h BrowserHistory
		if err := rows.Scan(&h.URL, &h.Title); err != nil {
			return nil, err
		}
		list = append(list, h)
	}
	return list, nil
}

func (t *browserHistoryTool) GetByKeyword(ctx context.Context, input GetByKeywordInput) ([]BrowserHistory, error) {
	db, tmp, err := openTempDB(t.historyFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	defer os.Remove(tmp)

	keyword := strings.TrimSpace(input.Keyword)
	if !strings.Contains(keyword, "%") {
		keyword = "%" + keyword + "%"
	}

	rows, err := db.QueryContext(ctx, getByKeywordQuery, keyword, input.Limit)
	if err != nil {
		return nil, fmt.Errorf("query history by keyword: %w", err)
	}
	defer rows.Close()

	var list []BrowserHistory
	for rows.Next() {
		var h BrowserHistory
		if err := rows.Scan(&h.URL, &h.Title); err != nil {
			return nil, err
		}
		list = append(list, h)
	}
	return list, nil
}
