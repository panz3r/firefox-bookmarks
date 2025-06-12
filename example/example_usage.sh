#!/bin/bash

# Example usage script for ff_bookmarks.py
# This script demonstrates how to use the Firefox bookmarks converter

echo "Firefox Bookmarks to HTML Converter - Example Usage"
echo "==================================================="
echo

# Check if the script exists
if [ ! -f "ff_bookmarks.py" ]; then
    echo "Error: ff_bookmarks.py not found in current directory"
    exit 1
fi

# Check if lz4 is installed
python3 -c "import lz4" 2>/dev/null
if [ $? -ne 0 ]; then
    echo "Installing required dependencies..."
    pip install lz4
fi

echo "Usage examples:"
echo
echo "1. Convert a Firefox jsonlz4 backup file:"
echo "   python ff_bookmarks.py ~/.mozilla/firefox/[profile]/bookmarkbackups/bookmarks-2025-06-11_123456_randomhash.jsonlz4"
echo
echo "2. Convert with custom output filename:"
echo "   python ff_bookmarks.py backup.jsonlz4 -o my_bookmarks.html"
echo
echo "3. Convert a JSON bookmark file:"
echo "   python ff_bookmarks.py bookmarks.json"
echo
echo "4. Show help:"
echo "   python ff_bookmarks.py --help"
echo

if [ -f "example/test_bookmarks.json" ]; then
    echo "Running test with the included sample file..."
    python ff_bookmarks.py example/test_bookmarks.json -o example/output.html
    if [ $? -eq 0 ]; then
        echo "✓ Test completed successfully!"
        echo "✓ Generated: example/output.html"
        echo
        echo "You can now import example/output.html into any web browser."
    else
        echo "✗ Test failed"
    fi
else
    echo "Note: example/test_bookmarks.json not found - skipping test"
fi
