package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/panz3r/firefox-bookmarks/bookmarks"
)

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

	// Load bookmark data
	loader := bookmarks.NewBookmarkLoader()

	fmt.Printf("Processing bookmark file: %s\n", inputFile)
	bookmarkData, err := loader.LoadBookmarksFromFile(inputFile)
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

	writer := bufio.NewWriter(outputFileHandle)
	err = bookmarks.ConvertBookmarksToHTML(writer, bookmarkData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to convert bookmarks: %v\n", err)
		os.Exit(1)
	}

	// Flush the buffered writer
	err = writer.Flush()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to flush output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted bookmarks to: %s\n", outputPath)
}
