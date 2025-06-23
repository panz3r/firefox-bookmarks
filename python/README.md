# Firefox Bookmarks Converter - Python Version

This directory contains the original Python implementation of the Firefox bookmarks converter.

## Files

- `ff_bookmarks.py` - Main Python script
- `requirements.txt` - Python dependencies

## Usage

```bash
# Install dependencies
pip install -r requirements.txt

# Convert bookmarks
python ff_bookmarks.py input_file [-o output_file]
```

## Examples

```bash
# Convert with default output filename
python ff_bookmarks.py bookmarks-backup.jsonlz4

# Convert with custom output filename  
python ff_bookmarks.py bookmarks-backup.jsonlz4 -o my_bookmarks.html

# Convert JSON file
python ff_bookmarks.py bookmarks.json -o bookmarks.html
```

## Note

⚠️ **The Go version is recommended** for better performance and easier distribution. See the main README and `README_GO.md` for details.

The Python version is maintained for compatibility and users who prefer Python environments.
