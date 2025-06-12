# Bookmarks Package

This package provides functionality for loading and converting Firefox bookmark backup files to HTML format.

## Features

- Load Firefox `.jsonlz4` compressed bookmark backup files
- Load regular `.json` bookmark files
- Convert bookmark data to standard HTML bookmark format
- Support for bookmark descriptions and timestamps
- Proper HTML escaping for security
- Comprehensive test coverage

## Package Structure

- `types.go` - Core data structures and constants
- `loader.go` - File validation and bookmark loading functionality
- `converter.go` - HTML conversion functionality
- `*_test.go` - Comprehensive test suites

## Usage

```go
import "github.com/panz3r/firefox-bookmarks/bookmarks"

// Load bookmarks from a file (auto-detects format)
loader := bookmarks.NewBookmarkLoader()
bookmarkData, err := loader.LoadBookmarksFromFile("bookmarks.json")
if err != nil {
    log.Fatal(err)
}

// Convert to HTML
var htmlOutput bytes.Buffer
err = bookmarks.ConvertBookmarksToHTML(&htmlOutput, bookmarkData)
if err != nil {
    log.Fatal(err)
}

fmt.Print(htmlOutput.String())
```

## Testing

Run all tests:
```bash
go test ./bookmarks
```

Run tests with coverage:
```bash
go test ./bookmarks -cover
```

Run integration tests:
```bash
go test ./bookmarks -run Integration
```

## Test Coverage

The package includes comprehensive tests covering:

- Unit tests for all public functions
- Integration tests for complete workflows
- Error handling scenarios
- Edge cases (empty bookmarks, special characters, etc.)
- File format validation
- HTML output validation

Current test coverage: ~89% of statements.
