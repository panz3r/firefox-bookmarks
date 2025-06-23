package bookmarks

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/panz3r/firefox-bookmarks/bookmarks"
)

// TestIntegration_CompleteWorkflow tests the complete workflow from JSON file to HTML output
func TestIntegration_CompleteWorkflow(t *testing.T) {
	// Create test bookmark data that matches the example
	testData := bookmarks.BookmarkData{
		Title: "Bookmarks Menu",
		Children: []bookmarks.BookmarkData{
			{
				Title: "Development",
				Children: []bookmarks.BookmarkData{
					{
						Title:        "GitHub",
						URI:          "https://github.com",
						DateAdded:    1639123456789000,
						LastModified: 1639123456789000,
					},
					{
						Title:        "Stack Overflow",
						URI:          "https://stackoverflow.com",
						DateAdded:    1639123456789000,
						LastModified: 1639123456789000,
					},
				},
				DateAdded:    1639123456789000,
				LastModified: 1639123456789000,
			},
			{
				Title:        "Example Bookmark",
				URI:          "https://example.com",
				DateAdded:    1639123456789000,
				LastModified: 1639123456789000,
				Annotations: []bookmarks.Annotation{
					{
						Name:  "bookmarkProperties/description",
						Value: "This is an example bookmark with a description",
					},
				},
			},
		},
		DateAdded:    1639123456789000,
		LastModified: 1639123456789000,
	}

	// Create temporary directory for test files
	tempDir := t.TempDir()

	// Create JSON test file
	jsonFile := filepath.Join(tempDir, "test_bookmarks.json")
	jsonData, err := json.MarshalIndent(testData, "", "    ")
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	err = os.WriteFile(jsonFile, jsonData, 0644)
	if err != nil {
		t.Fatalf("Failed to write JSON test file: %v", err)
	}

	// Test loading JSON file
	loader := bookmarks.NewBookmarkLoader()
	bookmarkData, err := loader.LoadBookmarksFromFile(jsonFile)
	if err != nil {
		t.Fatalf("Failed to load bookmark data: %v", err)
	}

	// Verify loaded data
	if bookmarkData.Title != testData.Title {
		t.Errorf("Title = %v, want %v", bookmarkData.Title, testData.Title)
	}
	if len(bookmarkData.Children) != len(testData.Children) {
		t.Errorf("Children count = %v, want %v", len(bookmarkData.Children), len(testData.Children))
	}

	// Test HTML conversion
	var htmlOutput bytes.Buffer
	err = bookmarks.ConvertBookmarksToHTML(&htmlOutput, bookmarkData)
	if err != nil {
		t.Fatalf("Failed to convert to HTML: %v", err)
	}

	htmlResult := htmlOutput.String()

	// Verify HTML output contains expected elements
	expectedElements := []string{
		"<!DOCTYPE NETSCAPE-Bookmark-file-1>",
		"<H1>Bookmarks Menu</H1>",
		"<DT><H3",
		"Development</H3>",
		`<DT><A HREF="https://github.com"`,
		"GitHub</A>",
		`<DT><A HREF="https://stackoverflow.com"`,
		"Stack Overflow</A>",
		`<DT><A HREF="https://example.com"`,
		"Example Bookmark</A>",
		"<DD>This is an example bookmark with a description",
		`ADD_DATE="1639123456"`,
		`LAST_MODIFIED="1639123456"`,
		"</DL><p>",
	}

	for _, element := range expectedElements {
		if !strings.Contains(htmlResult, element) {
			t.Errorf("Expected HTML to contain %q, but it didn't.\nHTML output:\n%s", element, htmlResult)
		}
	}

	// Verify HTML structure basics
	if !strings.HasPrefix(htmlResult, "<!DOCTYPE NETSCAPE-Bookmark-file-1>") {
		t.Error("HTML should start with proper DOCTYPE")
	}

	if strings.Count(htmlResult, "<DL><p>") != strings.Count(htmlResult, "</DL><p>") {
		t.Error("Mismatched <DL><p> and </DL><p> tags")
	}
}

// TestIntegration_ErrorHandling tests error handling throughout the workflow
func TestIntegration_ErrorHandling(t *testing.T) {
	loader := bookmarks.NewBookmarkLoader()
	tempDir := t.TempDir()

	// Test 1: File doesn't exist
	t.Run("FileNotExist", func(t *testing.T) {
		_, err := loader.LoadBookmarksFromFile("nonexistent.json")
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})

	// Test 2: Invalid file format
	t.Run("InvalidFileFormat", func(t *testing.T) {
		invalidFile := filepath.Join(tempDir, "invalid.txt")
		err := os.WriteFile(invalidFile, []byte("not a bookmark file"), 0644)
		if err != nil {
			t.Fatalf("Failed to create invalid file: %v", err)
		}

		_, err = loader.LoadBookmarksFromFile(invalidFile)
		if err == nil {
			t.Error("Expected error for invalid file format")
		}
		if !strings.Contains(err.Error(), "not a valid Firefox bookmark backup file") {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	// Test 3: Invalid JSON content
	t.Run("InvalidJSON", func(t *testing.T) {
		invalidJsonFile := filepath.Join(tempDir, "invalid.json")
		err := os.WriteFile(invalidJsonFile, []byte("{ invalid json"), 0644)
		if err != nil {
			t.Fatalf("Failed to create invalid JSON file: %v", err)
		}

		_, err = loader.LoadBookmarksFromFile(invalidJsonFile)
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}
	})
}

// TestIntegration_EmptyBookmarks tests handling of edge cases like empty bookmark structures
func TestIntegration_EmptyBookmarks(t *testing.T) {
	tests := []struct {
		name string
		data bookmarks.BookmarkData
	}{
		{
			name: "EmptyBookmarks",
			data: bookmarks.BookmarkData{
				Title:    "Empty Bookmarks",
				Children: []bookmarks.BookmarkData{},
			},
		},
		{
			name: "BookmarksWithNoTitle",
			data: bookmarks.BookmarkData{
				Children: []bookmarks.BookmarkData{
					{
						URI: "https://example.com",
					},
				},
			},
		},
		{
			name: "OnlyFolders",
			data: bookmarks.BookmarkData{
				Title: "Only Folders",
				Children: []bookmarks.BookmarkData{
					{
						Title:    "Empty Folder",
						Children: []bookmarks.BookmarkData{},
					},
				},
			},
		},
		{
			name: "OnlyBookmarks",
			data: bookmarks.BookmarkData{
				Title: "Only Bookmarks",
				Children: []bookmarks.BookmarkData{
					{
						Title: "Bookmark 1",
						URI:   "https://example1.com",
					},
					{
						Title: "Bookmark 2",
						URI:   "https://example2.com",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var htmlOutput bytes.Buffer
			err := bookmarks.ConvertBookmarksToHTML(&htmlOutput, &tt.data)
			if err != nil {
				t.Fatalf("ConvertBookmarksToHTML failed: %v", err)
			}

			htmlResult := htmlOutput.String()

			// Basic validation that HTML was generated
			if !strings.Contains(htmlResult, "<!DOCTYPE NETSCAPE-Bookmark-file-1>") {
				t.Error("Expected valid HTML output")
			}

			// Should contain a title (either from data or default)
			if !strings.Contains(htmlResult, "<H1>") {
				t.Error("Expected H1 tag in output")
			}
		})
	}
}

// TestIntegration_SpecialCharacters tests handling of special characters throughout the pipeline
func TestIntegration_SpecialCharacters(t *testing.T) {
	testData := bookmarks.BookmarkData{
		Title: "Test & Special <Characters>",
		Children: []bookmarks.BookmarkData{
			{
				Title: `Quotes "test" & symbols <script>`,
				URI:   "https://example.com?param=value&other=test",
				Annotations: []bookmarks.Annotation{
					{
						Name:  "bookmarkProperties/description",
						Value: `Description with "quotes" & <tags>`,
					},
				},
			},
		},
	}

	var htmlOutput bytes.Buffer
	err := bookmarks.ConvertBookmarksToHTML(&htmlOutput, &testData)
	if err != nil {
		t.Fatalf("ConvertBookmarksToHTML failed: %v", err)
	}

	htmlResult := htmlOutput.String()

	// Verify special characters are properly escaped
	expectedEscapedContent := []string{
		"Test &amp; Special &lt;Characters&gt;",
		"Quotes &#34;test&#34; &amp; symbols &lt;script&gt;",
		"https://example.com?param=value&amp;other=test",
		"Description with &#34;quotes&#34; &amp; &lt;tags&gt;",
	}

	for _, content := range expectedEscapedContent {
		if !strings.Contains(htmlResult, content) {
			t.Errorf("Expected escaped content %q in HTML output", content)
		}
	}

	// Verify no unescaped dangerous content
	dangerousContent := []string{
		"<script>",
		"& symbols", // unescaped ampersand followed by text
	}

	for _, danger := range dangerousContent {
		if strings.Contains(htmlResult, danger) {
			t.Errorf("Found unescaped dangerous content %q in HTML output", danger)
		}
	}
}
