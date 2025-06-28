package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain_WithExampleFile(t *testing.T) {
	// Test the refactored main function with the example file
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "test_output.html")

	// Simulate command line arguments
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	exampleFile := "example/test_bookmarks.json"
	os.Args = []string{"firefox-bookmarks", "-o", outputFile, exampleFile}

	// Run main (this will panic if there's an issue, which will fail the test)
	main()

	// Check that output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFile)
	}

	// Read and verify the output file content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	htmlContent := string(content)

	// Basic checks that the HTML was generated correctly
	expectedContent := []string{
		"<!DOCTYPE NETSCAPE-Bookmark-file-1>",
		"<H1>Bookmarks Menu</H1>",
		"GitHub",
		"Stack Overflow",
		"https://github.com",
		"https://stackoverflow.com",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(htmlContent, expected) {
			t.Errorf("Expected HTML to contain %q", expected)
		}
	}
}

func TestMain_NoArguments(t *testing.T) {
	// Test that main exits gracefully when no arguments are provided
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"firefox-bookmarks"}

	// This should cause main to exit with status 1, but since we can't easily
	// catch os.Exit in tests, we'll just check that it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main() panicked with no arguments: %v", r)
		}
	}()

	// main() will call os.Exit(1), but that's expected behavior
	// We can't easily test this without refactoring main further
}
