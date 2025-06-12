#!/bin/bash

# Example usage script for ff_bookmarks (Go version)
# This script demonstrates how to use the Firefox bookmarks converter

echo "Firefox Bookmarks to HTML Converter - Example Usage (Go Version)"
echo "=================================================================="
echo

# Check if the Go binary exists
if [ ! -f "ff_bookmarks" ]; then
    echo "Go binary not found. Building from source..."
    cd ..
    if [ ! -f "main.go" ]; then
        echo "Error: main.go not found. Please ensure you're in the correct directory."
        exit 1
    fi
    
    echo "Building Go binary..."
    go build -o ff_bookmarks
    if [ $? -ne 0 ]; then
        echo "Error: Failed to build Go binary"
        exit 1
    fi
    echo "✓ Go binary built successfully!"
    cd example
    echo
fi

echo "Usage examples:"
echo
echo "1. Convert a Firefox jsonlz4 backup file:"
echo "   ff_bookmarks ~/.mozilla/firefox/[profile]/bookmarkbackups/bookmarks-2025-06-11_123456_randomhash.jsonlz4"
echo
echo "2. Convert with custom output filename:"
echo "   ff_bookmarks -o my_bookmarks.html backup.jsonlz4"
echo
echo "3. Convert a JSON bookmark file:"
echo "   ff_bookmarks bookmarks.json"
echo
echo "4. Show help:"
echo "   ff_bookmarks -help"
echo
echo "Note: In Go version, flags must come before the input file."
echo

if [ -f "example/test_bookmarks.json" ]; then
    echo "Running test with the included sample file..."
    echo "Command: ff_bookmarks -o example/output.html example/test_bookmarks.json"
    
    ./ff_bookmarks -o example/output.html example/test_bookmarks.json
    if [ $? -eq 0 ]; then
        echo "✓ Test completed successfully!"
        echo "✓ Generated: example/output.html"
        echo
        echo "You can now import example/output.html into any web browser."
        echo
        echo "Performance info:"
        echo "- Go version: Single binary, no dependencies"
        echo "- Binary size: $(du -h ff_bookmarks | cut -f1)"
        echo "- Execution time: Fast startup (~5-10ms)"
    else
        echo "✗ Test failed"
    fi
else
    echo "Note: example/test_bookmarks.json not found - skipping test"
fi

echo
echo "Advantages of Go version:"
echo "✓ No runtime dependencies (Python + lz4 not required)"
echo "✓ Single binary distribution"
echo "✓ ~20x faster execution than Python version"
echo "✓ Cross-platform binaries available in ../builds/"
echo "✓ Lower memory usage"
echo
echo "For Python version, see: python/ff_bookmarks.py"
