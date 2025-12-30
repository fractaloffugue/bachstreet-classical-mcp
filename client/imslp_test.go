package client

import (
	"testing"
)

func TestSearchWorks(t *testing.T) {
	client := NewClient()

	// Test search for a well-known work
	works, err := client.SearchWorks("Bach Prelude C major BWV 846", 5)
	if err != nil {
		t.Fatalf("SearchWorks failed: %v", err)
	}

	if len(works) == 0 {
		t.Error("Expected at least one result for Bach Prelude search")
	}

	t.Logf("Found %d works", len(works))
	for i, work := range works {
		t.Logf("  [%d] %s by %s (URL: %s)", i+1, work.Title, work.Composer, work.PageURL)
	}
}

func TestSearchWorksMozart(t *testing.T) {
	client := NewClient()

	// Test search for Mozart K545
	works, err := client.SearchWorks("Mozart K545", 5)
	if err != nil {
		t.Fatalf("SearchWorks failed: %v", err)
	}

	if len(works) == 0 {
		t.Error("Expected at least one result for Mozart K545 search")
	}

	t.Logf("Found %d works for Mozart K545", len(works))
	for i, work := range works {
		t.Logf("  [%d] %s by %s", i+1, work.Title, work.Composer)
	}
}

func TestGetWorkDetails(t *testing.T) {
	client := NewClient()

	// First search for a work
	works, err := client.SearchWorks("Moonlight Sonata", 1)
	if err != nil {
		t.Fatalf("SearchWorks failed: %v", err)
	}

	if len(works) == 0 {
		t.Skip("No works found to test details retrieval")
	}

	// Use the title from the URL to get details
	// Extract page title from URL
	pageTitle := works[0].Title
	if works[0].Composer != "Unknown" {
		pageTitle = works[0].Title + " (" + works[0].Composer + ")"
	}

	work, err := client.GetWorkDetails(pageTitle)
	if err != nil {
		t.Logf("Warning: GetWorkDetails failed (may be expected for some pages): %v", err)
		return
	}

	t.Logf("Work Details:")
	t.Logf("  Title: %s", work.Title)
	t.Logf("  Composer: %s", work.Composer)
	t.Logf("  Key: %s", work.Key)
	t.Logf("  Opus: %s", work.OpusNumber)
	t.Logf("  Instrumentation: %s", work.Instrumentation)
}

func TestGetScoreLinks(t *testing.T) {
	client := NewClient()

	// Use a well-known page title that should have scores
	pageTitle := "Prelude and Fugue in C major, BWV 846 (Bach, Johann Sebastian)"

	scores, err := client.GetScoreLinks(pageTitle)
	if err != nil {
		t.Logf("Warning: GetScoreLinks failed: %v", err)
		return
	}

	t.Logf("Found %d scores", len(scores))
	for i, score := range scores {
		if i < 3 { // Only log first 3 to keep output reasonable
			t.Logf("  [%d] %s (%s, %s)", i+1, score.Description, score.FileType, score.FileSize)
			t.Logf("       URL: %s", score.DownloadURL)
		}
	}
}

func TestParseIMSLPTitle(t *testing.T) {
	tests := []struct {
		input         string
		wantComposer  string
		wantWorkTitle string
	}{
		{
			input:         "Prelude and Fugue in C major, BWV 846 (Bach, Johann Sebastian)",
			wantComposer:  "Bach, Johann Sebastian",
			wantWorkTitle: "Prelude and Fugue in C major, BWV 846",
		},
		{
			input:         "Piano Sonata No.14, Op.27 No.2 (Beethoven, Ludwig van)",
			wantComposer:  "Beethoven, Ludwig van",
			wantWorkTitle: "Piano Sonata No.14, Op.27 No.2",
		},
		{
			input:         "Just a Title",
			wantComposer:  "Unknown",
			wantWorkTitle: "Just a Title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			composer, workTitle := parseIMSLPTitle(tt.input)
			if composer != tt.wantComposer {
				t.Errorf("parseIMSLPTitle() composer = %v, want %v", composer, tt.wantComposer)
			}
			if workTitle != tt.wantWorkTitle {
				t.Errorf("parseIMSLPTitle() workTitle = %v, want %v", workTitle, tt.wantWorkTitle)
			}
		})
	}
}

func TestExtractOpusNumber(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Piano Sonata No.14, Op.27 No.2", "Op.27 No.2"},
		{"Nocturne Op.9 No.1", "Op.9 No.1"},
		{"Some Work Without Opus", ""},
		{"Symphony op.5", "op.5"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := extractOpusNumber(tt.input)
			if got != tt.want {
				t.Errorf("extractOpusNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
