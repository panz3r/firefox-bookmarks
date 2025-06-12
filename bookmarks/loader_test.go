package bookmarks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/pierrec/lz4/v4"
)

func TestFileValidator_IsValidJSONLZ4File(t *testing.T) {
	validator := NewFileValidator()

	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Test case 1: Non-existent file
	t.Run("NonExistentFile", func(t *testing.T) {
		if validator.IsValidJSONLZ4File("nonexistent.jsonlz4") {
			t.Error("Expected false for non-existent file")
		}
	})

	// Test case 2: Valid jsonlz4 file
	t.Run("ValidJSONLZ4File", func(t *testing.T) {
		validFile := filepath.Join(tempDir, "valid.jsonlz4")
		f, err := os.Create(validFile)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		defer f.Close()

		// Write Firefox LZ4 signature
		_, err = f.Write([]byte(FirefoxLZ4Signature))
		if err != nil {
			t.Fatalf("Failed to write signature: %v", err)
		}

		if !validator.IsValidJSONLZ4File(validFile) {
			t.Error("Expected true for valid jsonlz4 file")
		}
	})

	// Test case 3: Invalid signature
	t.Run("InvalidSignature", func(t *testing.T) {
		invalidFile := filepath.Join(tempDir, "invalid.jsonlz4")
		f, err := os.Create(invalidFile)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		defer f.Close()

		// Write invalid signature
		_, err = f.Write([]byte("invalid"))
		if err != nil {
			t.Fatalf("Failed to write invalid signature: %v", err)
		}

		if validator.IsValidJSONLZ4File(invalidFile) {
			t.Error("Expected false for invalid signature")
		}
	})

	// Test case 4: File too short
	t.Run("FileTooShort", func(t *testing.T) {
		shortFile := filepath.Join(tempDir, "short.jsonlz4")
		f, err := os.Create(shortFile)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		defer f.Close()

		// Write only partial signature
		_, err = f.Write([]byte("moz"))
		if err != nil {
			t.Fatalf("Failed to write partial signature: %v", err)
		}

		if validator.IsValidJSONLZ4File(shortFile) {
			t.Error("Expected false for file too short")
		}
	})
}

func TestFileValidator_IsJSONFile(t *testing.T) {
	validator := NewFileValidator()
	tempDir := t.TempDir()

	// Test case 1: Non-existent file
	t.Run("NonExistentFile", func(t *testing.T) {
		if validator.IsJSONFile("nonexistent.json") {
			t.Error("Expected false for non-existent file")
		}
	})

	// Test case 2: Valid JSON file
	t.Run("ValidJSONFile", func(t *testing.T) {
		validFile := filepath.Join(tempDir, "valid.json")
		testData := map[string]interface{}{
			"title": "Test",
			"uri":   "https://example.com",
		}

		data, err := json.Marshal(testData)
		if err != nil {
			t.Fatalf("Failed to marshal test data: %v", err)
		}

		err = os.WriteFile(validFile, data, 0644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		if !validator.IsJSONFile(validFile) {
			t.Error("Expected true for valid JSON file")
		}
	})

	// Test case 3: Invalid JSON file
	t.Run("InvalidJSONFile", func(t *testing.T) {
		invalidFile := filepath.Join(tempDir, "invalid.json")
		err := os.WriteFile(invalidFile, []byte("invalid json {"), 0644)
		if err != nil {
			t.Fatalf("Failed to write invalid test file: %v", err)
		}

		if validator.IsJSONFile(invalidFile) {
			t.Error("Expected false for invalid JSON file")
		}
	})
}

func TestBookmarkLoader_LoadJSONFile(t *testing.T) {
	loader := NewBookmarkLoader()
	tempDir := t.TempDir()

	// Test case 1: Valid JSON file
	t.Run("ValidJSONFile", func(t *testing.T) {
		testData := BookmarkData{
			Title: "Test Bookmarks",
			Children: []BookmarkData{
				{
					Title: "Test Bookmark",
					URI:   "https://example.com",
				},
			},
		}

		validFile := filepath.Join(tempDir, "test.json")
		data, err := json.Marshal(testData)
		if err != nil {
			t.Fatalf("Failed to marshal test data: %v", err)
		}

		err = os.WriteFile(validFile, data, 0644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		result, err := loader.LoadJSONFile(validFile)
		if err != nil {
			t.Fatalf("LoadJSONFile failed: %v", err)
		}

		if result.Title != testData.Title {
			t.Errorf("Title = %v, want %v", result.Title, testData.Title)
		}
		if len(result.Children) != 1 {
			t.Errorf("Children length = %v, want %v", len(result.Children), 1)
		}
		if result.Children[0].URI != testData.Children[0].URI {
			t.Errorf("Child URI = %v, want %v", result.Children[0].URI, testData.Children[0].URI)
		}
	})

	// Test case 2: Non-existent file
	t.Run("NonExistentFile", func(t *testing.T) {
		_, err := loader.LoadJSONFile("nonexistent.json")
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})

	// Test case 3: Invalid JSON file
	t.Run("InvalidJSONFile", func(t *testing.T) {
		invalidFile := filepath.Join(tempDir, "invalid.json")
		err := os.WriteFile(invalidFile, []byte("invalid json"), 0644)
		if err != nil {
			t.Fatalf("Failed to write invalid test file: %v", err)
		}

		_, err = loader.LoadJSONFile(invalidFile)
		if err == nil {
			t.Error("Expected error for invalid JSON file")
		}
	})
}

func TestBookmarkLoader_DecompressJSONLZ4(t *testing.T) {
	loader := NewBookmarkLoader()
	tempDir := t.TempDir()

	// Test case 1: Valid compressed file
	t.Run("ValidCompressedFile", func(t *testing.T) {
		// Create test bookmark data
		testData := BookmarkData{
			Title: "Test Bookmarks",
			Children: []BookmarkData{
				{
					Title: "Test Bookmark",
					URI:   "https://example.com",
				},
			},
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(testData)
		if err != nil {
			t.Fatalf("Failed to marshal test data: %v", err)
		}

		// Compress with LZ4
		compressedData := make([]byte, lz4.CompressBlockBound(len(jsonData)))
		compressedSize, err := lz4.CompressBlock(jsonData, compressedData, nil)
		if err != nil {
			t.Fatalf("Failed to compress test data: %v", err)
		}

		// Create test file
		testFile := filepath.Join(tempDir, "test.jsonlz4")
		f, err := os.Create(testFile)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		defer f.Close()

		// Write Firefox LZ4 header
		_, err = f.Write([]byte(FirefoxLZ4Signature))
		if err != nil {
			t.Fatalf("Failed to write signature: %v", err)
		}

		// Write padding to reach header size
		padding := make([]byte, FirefoxLZ4HeaderSize-len(FirefoxLZ4Signature))
		_, err = f.Write(padding)
		if err != nil {
			t.Fatalf("Failed to write padding: %v", err)
		}

		// Write compressed data
		_, err = f.Write(compressedData[:compressedSize])
		if err != nil {
			t.Fatalf("Failed to write compressed data: %v", err)
		}

		// Test decompression
		result, err := loader.DecompressJSONLZ4(testFile)
		if err != nil {
			t.Fatalf("DecompressJSONLZ4 failed: %v", err)
		}

		if result.Title != testData.Title {
			t.Errorf("Title = %v, want %v", result.Title, testData.Title)
		}
		if len(result.Children) != 1 {
			t.Errorf("Children length = %v, want %v", len(result.Children), 1)
		}
	})

	// Test case 2: Non-existent file
	t.Run("NonExistentFile", func(t *testing.T) {
		_, err := loader.DecompressJSONLZ4("nonexistent.jsonlz4")
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})
}

func TestBookmarkLoader_LoadBookmarksFromFile(t *testing.T) {
	loader := NewBookmarkLoader()
	tempDir := t.TempDir()

	// Test case 1: Valid JSON file
	t.Run("ValidJSONFile", func(t *testing.T) {
		testData := BookmarkData{
			Title: "Test Bookmarks",
			Children: []BookmarkData{
				{Title: "Test", URI: "https://example.com"},
			},
		}

		validFile := filepath.Join(tempDir, "test.json")
		data, err := json.Marshal(testData)
		if err != nil {
			t.Fatalf("Failed to marshal test data: %v", err)
		}

		err = os.WriteFile(validFile, data, 0644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		result, err := loader.LoadBookmarksFromFile(validFile)
		if err != nil {
			t.Fatalf("LoadBookmarksFromFile failed: %v", err)
		}

		if result.Title != testData.Title {
			t.Errorf("Title = %v, want %v", result.Title, testData.Title)
		}
	})

	// Test case 2: Invalid file format
	t.Run("InvalidFileFormat", func(t *testing.T) {
		invalidFile := filepath.Join(tempDir, "invalid.txt")
		err := os.WriteFile(invalidFile, []byte("not a bookmark file"), 0644)
		if err != nil {
			t.Fatalf("Failed to write invalid test file: %v", err)
		}

		_, err = loader.LoadBookmarksFromFile(invalidFile)
		if err == nil {
			t.Error("Expected error for invalid file format")
		}
	})
}
