package bookmarks

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/panz3r/firefox-bookmarks/bookmarks"
)

// TestWithExampleFiles tests the refactored code using the actual example files
func TestWithExampleFiles(t *testing.T) {
	// Get the path to the example file
	exampleFile := filepath.Join("..", "example", "test_bookmarks.json")

	// Load the example bookmark file
	loader := bookmarks.NewBookmarkLoader()
	bookmarkData, err := loader.LoadBookmarksFromFile(exampleFile)
	if err != nil {
		t.Fatalf("Failed to load example file: %v", err)
	}

	// Verify the loaded data structure
	if bookmarkData.Title != "Bookmarks Menu" {
		t.Errorf("Expected title 'Bookmarks Menu', got %q", bookmarkData.Title)
	}

	if len(bookmarkData.Children) == 0 {
		t.Error("Expected bookmark data to have children")
	}

	// Find the Development folder
	var developmentFolder *bookmarks.BookmarkData
	for i := range bookmarkData.Children {
		if bookmarkData.Children[i].Title == "Development" {
			developmentFolder = &bookmarkData.Children[i]
			break
		}
	}

	if developmentFolder == nil {
		t.Error("Expected to find 'Development' folder")
	} else {
		// Check that Development folder has children
		if len(developmentFolder.Children) == 0 {
			t.Error("Expected Development folder to have children")
		}

		// Verify GitHub bookmark exists
		foundGitHub := false
		for _, child := range developmentFolder.Children {
			if child.Title == "GitHub" && child.URI == "https://github.com" {
				foundGitHub = true
				break
			}
		}
		if !foundGitHub {
			t.Error("Expected to find GitHub bookmark in Development folder")
		}
	}

	// Convert to HTML
	var htmlOutput bytes.Buffer
	err = bookmarks.ConvertBookmarksToHTML(&htmlOutput, bookmarkData)
	if err != nil {
		t.Fatalf("Failed to convert bookmarks to HTML: %v", err)
	}

	htmlResult := htmlOutput.String()

	// Verify HTML contains expected bookmarks from example file
	expectedBookmarks := []string{
		"GitHub",
		"Stack Overflow",
		"BBC News",
		"Example Bookmark",
		"https://github.com",
		"https://stackoverflow.com",
		"https://www.bbc.com/news",
		"https://example.com",
	}

	for _, bookmark := range expectedBookmarks {
		if !strings.Contains(htmlResult, bookmark) {
			t.Errorf("Expected HTML to contain %q", bookmark)
		}
	}

	// Verify description is included
	if !strings.Contains(htmlResult, "This is an example bookmark with a description") {
		t.Error("Expected to find bookmark description in HTML output")
	}

	// Verify folder structure
	expectedFolders := []string{
		"Development",
		"News",
		"Repositories",
	}

	for _, folder := range expectedFolders {
		if !strings.Contains(htmlResult, folder) {
			t.Errorf("Expected HTML to contain folder %q", folder)
		}
	}

	// Verify proper HTML structure
	if !strings.HasPrefix(htmlResult, "<!DOCTYPE NETSCAPE-Bookmark-file-1>") {
		t.Error("HTML should start with proper DOCTYPE")
	}

	if !strings.Contains(htmlResult, "<H1>Bookmarks Menu</H1>") {
		t.Error("Expected main title in HTML")
	}

	// Check that the HTML has proper nesting by counting tags
	dlOpenCount := strings.Count(htmlResult, "<DL><p>")
	dlCloseCount := strings.Count(htmlResult, "</DL><p>")
	if dlOpenCount != dlCloseCount {
		t.Errorf("Mismatched DL tags: %d open, %d close", dlOpenCount, dlCloseCount)
	}
}
