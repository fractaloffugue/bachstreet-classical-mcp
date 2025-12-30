package models

// Composer represents a composer/person in IMSLP
type Composer struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"fullName,omitempty"`
	Nationality string `json:"nationality,omitempty"`
	BirthYear   string `json:"birthYear,omitempty"`
	DeathYear   string `json:"deathYear,omitempty"`
	Period      string `json:"period,omitempty"` // Baroque, Classical, Romantic, etc.
}

// Work represents a musical composition
type Work struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Composer        string   `json:"composer"`
	ComposerID      string   `json:"composerId,omitempty"`
	Instrumentation string   `json:"instrumentation,omitempty"`
	Key             string   `json:"key,omitempty"`
	OpusNumber      string   `json:"opusNumber,omitempty"`
	MovementCount   int      `json:"movementCount,omitempty"`
	YearComposed    string   `json:"yearComposed,omitempty"`
	Genre           string   `json:"genre,omitempty"`
	ScoreCount      int      `json:"scoreCount,omitempty"`
	PageURL         string   `json:"pageUrl,omitempty"`
}

// Score represents an edition/arrangement of a work
type Score struct {
	ID          string `json:"id"`
	WorkID      string `json:"workId"`
	Edition     string `json:"edition,omitempty"`
	Editor      string `json:"editor,omitempty"`
	Publisher   string `json:"publisher,omitempty"`
	Year        string `json:"year,omitempty"`
	Description string `json:"description,omitempty"`
	PageCount   int    `json:"pageCount,omitempty"`
	FileType    string `json:"fileType"` // PDF, MusicXML, etc.
	DownloadURL string `json:"downloadUrl"`
	FileSize    string `json:"fileSize,omitempty"`
}

// SearchResult wraps search results with metadata
type SearchResult struct {
	Query      string  `json:"query"`
	TotalCount int     `json:"totalCount"`
	Works      []Work  `json:"works,omitempty"`
	Composers  []Composer `json:"composers,omitempty"`
}
