package bookmarks

import (
	"testing"
)

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant interface{}
		expected interface{}
	}{
		{"FirefoxLZ4Signature", FirefoxLZ4Signature, "mozLz4"},
		{"FirefoxLZ4HeaderSize", FirefoxLZ4HeaderSize, 12},
		{"DefaultBufferSize", DefaultBufferSize, 10 * 1024 * 1024},
		{"IndentSize", IndentSize, 4},
		{"BookmarkSeparatorType", BookmarkSeparatorType, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("constant %s = %v, want %v", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

func TestBookmarkData_Structure(t *testing.T) {
	// Test that BookmarkData can be properly initialized
	bookmark := BookmarkData{
		Title:        "Test Bookmark",
		URI:          "https://example.com",
		DateAdded:    1639123456789000,
		LastModified: 1639123456789000,
		TypeCode:     1,
		Annotations: []Annotation{
			{Name: "bookmarkProperties/description", Value: "Test description"},
		},
	}

	if bookmark.Title != "Test Bookmark" {
		t.Errorf("Title = %v, want %v", bookmark.Title, "Test Bookmark")
	}
	if bookmark.URI != "https://example.com" {
		t.Errorf("URI = %v, want %v", bookmark.URI, "https://example.com")
	}
	if len(bookmark.Annotations) != 1 {
		t.Errorf("Annotations length = %v, want %v", len(bookmark.Annotations), 1)
	}
}

func TestAnnotation_Structure(t *testing.T) {
	anno := Annotation{
		Name:  "bookmarkProperties/description",
		Value: "Test description",
	}

	if anno.Name != "bookmarkProperties/description" {
		t.Errorf("Name = %v, want %v", anno.Name, "bookmarkProperties/description")
	}
	if anno.Value != "Test description" {
		t.Errorf("Value = %v, want %v", anno.Value, "Test description")
	}
}
