#!/bin/bash

# Performance comparison between Python and Go versions
# Tests both versions with the same input file

echo "Firefox Bookmarks Converter - Performance Comparison"
echo "==================================================="
echo

INPUT_FILE="example/test_bookmarks.json"

if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: Test file $INPUT_FILE not found"
    exit 1
fi

echo "Testing with file: $INPUT_FILE"
echo

# Test Python version
echo "Python version:"
if [ -f "python/ff_bookmarks.py" ]; then
    time python3 python/ff_bookmarks.py "$INPUT_FILE" -o test_python_perf.html
    PYTHON_SIZE=$(du -h test_python_perf.html | cut -f1)
    echo "Output size: $PYTHON_SIZE"
else
    echo "Python version (python/ff_bookmarks.py) not found"
fi

echo

# Test Go version  
echo "Go version:"
if [ -f "ff_bookmarks" ]; then
    time ./ff_bookmarks -o test_go_perf.html "$INPUT_FILE"
    GO_SIZE=$(du -h test_go_perf.html | cut -f1)
    echo "Output size: $GO_SIZE"
else
    echo "Go version (ff_bookmarks) not found - run 'go build' first"
fi

echo

# Compare outputs
if [ -f "test_python_perf.html" ] && [ -f "test_go_perf.html" ]; then
    echo "Comparing outputs:"
    if diff test_python_perf.html test_go_perf.html > /dev/null; then
        echo "✅ Outputs are identical"
    else
        echo "❌ Outputs differ"
    fi
    
    # Cleanup
    rm -f test_python_perf.html test_go_perf.html
fi

echo
echo "Binary sizes:"
if [ -f "ff_bookmarks" ]; then
    echo "Go binary: $(du -h ff_bookmarks | cut -f1)"
fi

echo "Python script: $(du -h python/ff_bookmarks.py | cut -f1)"
