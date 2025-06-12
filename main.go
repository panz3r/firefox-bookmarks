package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/pierrec/lz4/v4"
)

// Constants from the Python version
const (
	FirefoxLZ4Signature   = "mozLz4"
	FirefoxLZ4HeaderSize  = 12
	DefaultBufferSize     = 10 * 1024 * 1024 // 10MB
	IndentSize            = 4
	BookmarkSeparatorType = 3
)

// BookmarkData represents the structure of bookmark data
type BookmarkData struct {
	Title        string         `json:"title,omitempty"`
	URI          string         `json:"uri,omitempty"`
	Children     []BookmarkData `json:"children,omitempty"`
	DateAdded    int64          `json:"dateAdded,omitempty"`
	LastModified int64          `json:"lastModified,omitempty"`
	TypeCode     int            `json:"typeCode,omitempty"`
	Annos        []Annotation   `json:"annos,omitempty"`
}

// Annotation represents bookmark annotations (like descriptions)
type Annotation struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// HTMLWriter wraps an io.Writer with indentation functionality
type HTMLWriter struct {
	writer io.Writer
}

// NewHTMLWriter creates a new HTMLWriter
func NewHTMLWriter(w io.Writer) *HTMLWriter {
	return &HTMLWriter{writer: w}
}

// WriteIndented writes indented text to the output
func (hw *HTMLWriter) WriteIndented(indent int, text string) error {
	indentation := strings.Repeat(" ", IndentSize*indent)
	_, err := fmt.Fprintf(hw.writer, "%s%s\n", indentation, text)
	return err
}

// htmlEscape escapes HTML special characters to prevent XSS and display issues
func htmlEscape(text string) string {
	if text == "" {
		return ""
	}
	return html.EscapeString(text)
}

// isValidJSONLZ4File checks if the file is a valid Firefox jsonlz4 bookmark backup file
func isValidJSONLZ4File(filename string) bool {
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

// isJSONFile checks if the file is a valid JSON file by trying to parse it
func isJSONFile(filename string) bool {
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

// decompressJSONLZ4 decompresses a Firefox jsonlz4 bookmark backup file and returns the JSON data
func decompressJSONLZ4(filename string) (*BookmarkData, error) {
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

// loadJSONFile loads and parses a regular JSON file
func loadJSONFile(filename string) (*BookmarkData, error) {
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

// convertFirefoxTimestamp converts Firefox timestamp to Unix timestamp string
func convertFirefoxTimestamp(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	return fmt.Sprintf("%d", int64(math.Floor(float64(timestamp)/1000000)))
}

// formatDateAttributes formats date attributes for HTML bookmark tags
func formatDateAttributes(data *BookmarkData) string {
	var attributes []string

	if data.DateAdded != 0 {
		dateAdded := convertFirefoxTimestamp(data.DateAdded)
		if dateAdded != "" {
			attributes = append(attributes, fmt.Sprintf(` ADD_DATE="%s"`, dateAdded))
		}
	}

	if data.LastModified != 0 {
		lastModified := convertFirefoxTimestamp(data.LastModified)
		if lastModified != "" {
			attributes = append(attributes, fmt.Sprintf(` LAST_MODIFIED="%s"`, lastModified))
		}
	}

	return strings.Join(attributes, "")
}

// writeHTMLHeader writes the HTML document header
func writeHTMLHeader(writer *HTMLWriter, title string) error {
	header := fmt.Sprintf(`<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
    It will be read and overwritten.
    DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>%s</H1>
<DL><p>`, htmlEscape(title))

	return writer.WriteIndented(0, header)
}

// writeFolder writes a bookmark folder to HTML
func writeFolder(writer *HTMLWriter, data *BookmarkData, indent int) error {
	title := htmlEscape(data.Title)
	dateAttrs := formatDateAttributes(data)

	err := writer.WriteIndented(indent, fmt.Sprintf(`<DT><H3%s>%s</H3>`, dateAttrs, title))
	if err != nil {
		return err
	}
	return writer.WriteIndented(indent, `<DL><p>`)
}

// writeBookmark writes a single bookmark to HTML
func writeBookmark(writer *HTMLWriter, data *BookmarkData, indent int) error {
	uri := data.URI
	title := data.Title
	if title == "" {
		title = uri
	}
	title = htmlEscape(title)
	dateAttrs := formatDateAttributes(data)

	err := writer.WriteIndented(indent,
		fmt.Sprintf(`<DT><A HREF="%s"%s>%s</A>`, htmlEscape(uri), dateAttrs, title))
	if err != nil {
		return err
	}

	// Handle bookmark descriptions
	for _, anno := range data.Annos {
		if anno.Name == "bookmarkProperties/description" {
			description := htmlEscape(anno.Value)
			err = writer.WriteIndented(indent, fmt.Sprintf(`<DD>%s`, description))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// convertBookmarksToHTML converts bookmark data to HTML format recursively
func convertBookmarksToHTML(writer *HTMLWriter, data *BookmarkData, indent int) error {
	// Handle containers (folders) with children
	if data.Children != nil && len(data.Children) > 0 {
		if indent == 0 {
			// Output the main header
			title := data.Title
			if title == "" {
				title = "Bookmarks Menu"
			}
			err := writeHTMLHeader(writer, title)
			if err != nil {
				return err
			}
		} else {
			// Output a folder
			err := writeFolder(writer, data, indent)
			if err != nil {
				return err
			}
		}

		// Process children
		for _, child := range data.Children {
			// Skip separators (typeCode 3)
			if child.TypeCode == BookmarkSeparatorType {
				continue
			}
			err := convertBookmarksToHTML(writer, &child, indent+1)
			if err != nil {
				return err
			}
		}

		return writer.WriteIndented(indent, `</DL><p>`)
	} else if data.URI != "" {
		// Output a bookmark
		return writeBookmark(writer, data, indent)
	}

	return nil
}

// printUsage prints the usage information
func printUsage() {
	fmt.Printf(`
ff_bookmarks [-o OUTPUT_FILE] input_file

Converts Firefox bookmark backup files to HTML format.
Supports both .jsonlz4 (compressed backup) and .json (uncompressed) input files.

Examples:
    ff_bookmarks bookmarks-2025-06-11.jsonlz4
    ff_bookmarks -o my_bookmarks.html bookmarks-2025-06-11.jsonlz4  
    ff_bookmarks -o bookmarks.html bookmarks.json

Options:
`)
	flag.PrintDefaults()
}

func main() {
	var outputFile string
	var showHelp bool

	flag.StringVar(&outputFile, "o", "", "Path to output HTML file. If omitted, uses input filename with .html extension")
	flag.BoolVar(&showHelp, "help", false, "Show this help message")
	flag.Usage = printUsage
	flag.Parse()

	if showHelp {
		printUsage()
		return
	}

	// Check if input file is provided
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Error: Input file is required.\n")
		printUsage()
		os.Exit(1)
	}

	inputFile := flag.Arg(0)

	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Input file '%s' does not exist.\n", inputFile)
		os.Exit(1)
	}

	// Determine output filename
	var outputPath string
	if outputFile != "" {
		outputPath = outputFile
	} else {
		ext := filepath.Ext(inputFile)
		outputPath = strings.TrimSuffix(inputFile, ext) + ".html"
	}

	// Determine file type and load data accordingly
	var bookmarkData *BookmarkData
	var err error

	if isValidJSONLZ4File(inputFile) {
		fmt.Printf("Processing Firefox jsonlz4 bookmark backup: %s\n", inputFile)
		bookmarkData, err = decompressJSONLZ4(inputFile)
	} else if isJSONFile(inputFile) {
		fmt.Printf("Processing JSON bookmark file: %s\n", inputFile)
		bookmarkData, err = loadJSONFile(inputFile)
	} else {
		fmt.Fprintf(os.Stderr, "Error: '%s' is not a valid Firefox bookmark backup file (.jsonlz4) or JSON file.\n", inputFile)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Conversion failed: %v\n", err)
		os.Exit(1)
	}

	// Convert to HTML
	fmt.Println("Converting bookmarks to HTML format...")
	outputFileHandle, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFileHandle.Close()

	writer := NewHTMLWriter(bufio.NewWriter(outputFileHandle))
	err = convertBookmarksToHTML(writer, bookmarkData, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to convert bookmarks: %v\n", err)
		os.Exit(1)
	}

	// Flush the buffered writer
	if bufferedWriter, ok := writer.writer.(*bufio.Writer); ok {
		err = bufferedWriter.Flush()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to flush output: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Successfully converted bookmarks to: %s\n", outputPath)
}
