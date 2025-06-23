package bookmarks

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pierrec/lz4/v4"
)

// FileValidator handles file validation
type FileValidator struct{}

// NewFileValidator creates a new FileValidator
func NewFileValidator() *FileValidator {
	return &FileValidator{}
}

// IsValidJSONLZ4File checks if the file is a valid Firefox jsonlz4 bookmark backup file
func (fv *FileValidator) IsValidJSONLZ4File(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	header := make([]byte, len(FirefoxLZ4Signature))
	n, err := file.Read(header)
	if err != nil || n != len(FirefoxLZ4Signature) {
		return false
	}

	return string(header) == FirefoxLZ4Signature
}

// IsJSONFile checks if the file is a valid JSON file by trying to parse it
func (fv *FileValidator) IsJSONFile(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data interface{}
	return decoder.Decode(&data) == nil
}

// BookmarkLoader handles loading bookmark data from different file formats
type BookmarkLoader struct {
	validator *FileValidator
}

// NewBookmarkLoader creates a new BookmarkLoader
func NewBookmarkLoader() *BookmarkLoader {
	return &BookmarkLoader{
		validator: NewFileValidator(),
	}
}

// DecompressJSONLZ4 decompresses a Firefox jsonlz4 bookmark backup file and returns the JSON data
func (bl *BookmarkLoader) DecompressJSONLZ4(filename string) (*BookmarkData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("file reading error: %w", err)
	}
	defer file.Close()

	// Skip the Firefox LZ4 header
	_, err = file.Seek(FirefoxLZ4HeaderSize, 0)
	if err != nil {
		return nil, fmt.Errorf("error seeking past header: %w", err)
	}

	// Read the compressed data
	compressedData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading compressed data: %w", err)
	}

	// Decompress the data
	decompressedData := make([]byte, DefaultBufferSize)
	n, err := lz4.UncompressBlock(compressedData, decompressedData)
	if err != nil {
		return nil, fmt.Errorf("LZ4 decompression error: %w", err)
	}

	// Parse JSON
	var bookmarkData BookmarkData
	err = json.Unmarshal(decompressedData[:n], &bookmarkData)
	if err != nil {
		return nil, fmt.Errorf("JSON parsing error: %w", err)
	}

	return &bookmarkData, nil
}

// LoadJSONFile loads and parses a regular JSON file
func (bl *BookmarkLoader) LoadJSONFile(filename string) (*BookmarkData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var bookmarkData BookmarkData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&bookmarkData)
	if err != nil {
		return nil, fmt.Errorf("error loading JSON file: %w", err)
	}

	return &bookmarkData, nil
}

// LoadBookmarksFromFile loads bookmarks from a file, auto-detecting the format
func (bl *BookmarkLoader) LoadBookmarksFromFile(filename string) (*BookmarkData, error) {
	if bl.validator.IsValidJSONLZ4File(filename) {
		return bl.DecompressJSONLZ4(filename)
	} else if bl.validator.IsJSONFile(filename) {
		return bl.LoadJSONFile(filename)
	} else {
		return nil, fmt.Errorf("file '%s' is not a valid Firefox bookmark backup file (.jsonlz4) or JSON file", filename)
	}
}
