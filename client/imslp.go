package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"bachstreet-classical-mcp/models"
)

const (
	imslpWikiAPIBase = "https://imslp.org/api.php"
	imslpBaseURL     = "https://imslp.org"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchWorks searches for musical works by title or composer
func (c *Client) SearchWorks(query string, limit int) ([]models.Work, error) {
	// Use MediaWiki API to search
	params := url.Values{}
	params.Set("action", "query")
	params.Set("list", "search")
	params.Set("srsearch", query)
	params.Set("srnamespace", "0") // Main namespace
	params.Set("srlimit", fmt.Sprintf("%d", limit))
	params.Set("format", "json")

	apiURL := fmt.Sprintf("%s?%s", imslpWikiAPIBase, params.Encode())

	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search works: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Query struct {
			Search []struct {
				Title   string `json:"title"`
				PageID  int    `json:"pageid"`
				Snippet string `json:"snippet"`
			} `json:"search"`
		} `json:"query"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	works := make([]models.Work, 0, len(result.Query.Search))
	for _, item := range result.Query.Search {
		// Parse composer and work title from IMSLP format
		// IMSLP pages are typically "Work Title (Composer, Name)"
		composer, title := parseIMSLPTitle(item.Title)

		works = append(works, models.Work{
			ID:       fmt.Sprintf("%d", item.PageID),
			Title:    title,
			Composer: composer,
			PageURL:  fmt.Sprintf("%s/wiki/%s", imslpBaseURL, url.PathEscape(item.Title)),
		})
	}

	return works, nil
}

// GetWorkDetails retrieves detailed information about a specific work
func (c *Client) GetWorkDetails(pageTitle string) (*models.Work, error) {
	// Use MediaWiki API to get page content
	params := url.Values{}
	params.Set("action", "parse")
	params.Set("page", pageTitle)
	params.Set("prop", "text|categories")
	params.Set("format", "json")

	apiURL := fmt.Sprintf("%s?%s", imslpWikiAPIBase, params.Encode())

	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get work details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Parse struct {
			Title      string `json:"title"`
			PageID     int    `json:"pageid"`
			Text       struct {
				Content string `json:"*"`
			} `json:"text"`
			Categories []struct {
				Category string `json:"*"`
			} `json:"categories"`
		} `json:"parse"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	composer, title := parseIMSLPTitle(result.Parse.Title)

	work := &models.Work{
		ID:       fmt.Sprintf("%d", result.Parse.PageID),
		Title:    title,
		Composer: composer,
		PageURL:  fmt.Sprintf("%s/wiki/%s", imslpBaseURL, url.PathEscape(pageTitle)),
	}

	// Extract additional metadata from categories and content
	categories := make([]string, len(result.Parse.Categories))
	for i, cat := range result.Parse.Categories {
		categories[i] = cat.Category
	}
	work.Instrumentation = extractInstrumentation(categories)
	work.Key = extractKey(result.Parse.Text.Content)
	work.OpusNumber = extractOpusNumber(title)

	return work, nil
}

// GetScoreLinks retrieves download links for scores of a work
func (c *Client) GetScoreLinks(pageTitle string) ([]models.Score, error) {
	// Use MediaWiki API to get images (PDF files) from the page
	params := url.Values{}
	params.Set("action", "query")
	params.Set("titles", pageTitle)
	params.Set("prop", "images")
	params.Set("imlimit", "50")
	params.Set("format", "json")

	apiURL := fmt.Sprintf("%s?%s", imslpWikiAPIBase, params.Encode())

	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get score links: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Query struct {
			Pages map[string]struct {
				Images []struct {
					Title string `json:"title"`
				} `json:"images"`
			} `json:"pages"`
		} `json:"query"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	scores := []models.Score{}

	// Get file details for each image
	for _, page := range result.Query.Pages {
		for _, img := range page.Images {
			// Filter for PDF files
			if strings.HasSuffix(strings.ToLower(img.Title), ".pdf") {
				fileURL, fileSize, err := c.getFileURL(img.Title)
				if err != nil {
					continue // Skip files we can't access
				}

				scores = append(scores, models.Score{
					ID:          img.Title,
					Description: img.Title,
					FileType:    "PDF",
					DownloadURL: fileURL,
					FileSize:    fileSize,
				})
			}
		}
	}

	return scores, nil
}

// getFileURL retrieves the actual download URL for a file
func (c *Client) getFileURL(fileName string) (string, string, error) {
	params := url.Values{}
	params.Set("action", "query")
	params.Set("titles", fileName)
	params.Set("prop", "imageinfo")
	params.Set("iiprop", "url|size")
	params.Set("format", "json")

	apiURL := fmt.Sprintf("%s?%s", imslpWikiAPIBase, params.Encode())

	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var result struct {
		Query struct {
			Pages map[string]struct {
				ImageInfo []struct {
					URL  string `json:"url"`
					Size int    `json:"size"`
				} `json:"imageinfo"`
			} `json:"pages"`
		} `json:"query"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", err
	}

	for _, page := range result.Query.Pages {
		if len(page.ImageInfo) > 0 {
			info := page.ImageInfo[0]
			sizeStr := fmt.Sprintf("%.2f MB", float64(info.Size)/(1024*1024))
			return info.URL, sizeStr, nil
		}
	}

	return "", "", fmt.Errorf("no file info found")
}

// Helper functions

func parseIMSLPTitle(title string) (composer, workTitle string) {
	// IMSLP format: "Work Title (Composer, Name)" or just "Work Title"
	if idx := strings.Index(title, " ("); idx != -1 {
		workTitle = strings.TrimSpace(title[:idx])
		composerPart := strings.TrimSuffix(title[idx+2:], ")")
		composer = strings.TrimSpace(composerPart)
	} else {
		workTitle = title
		composer = "Unknown"
	}
	return
}

func extractInstrumentation(categories []string) string {
	instruments := []string{}
	for _, cat := range categories {
		// Look for instrument categories
		if strings.Contains(cat, "For ") {
			instruments = append(instruments, cat)
		}
	}
	return strings.Join(instruments, ", ")
}

func extractKey(content string) string {
	keys := []string{"C major", "C minor", "D major", "D minor", "E major", "E minor",
		"F major", "F minor", "G major", "G minor", "A major", "A minor", "B major", "B minor"}
	for _, key := range keys {
		if strings.Contains(content, key) {
			return key
		}
	}
	return ""
}

func extractOpusNumber(title string) string {
	if strings.Contains(strings.ToLower(title), "op.") {
		parts := strings.Split(title, "Op.")
		if len(parts) > 1 {
			opNum := strings.TrimSpace(strings.Split(parts[1], ",")[0])
			return "Op." + opNum
		}
	}
	return ""
}
