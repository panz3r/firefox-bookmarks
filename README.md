# Firefox Bookmarks Backup to HTML Converter

A Python script that converts Firefox bookmark backup files directly to HTML format. This script merges the functionality of jsonlz4 decompression and HTML conversion into a single tool.

## Features

- **Direct conversion**: Convert Firefox `.jsonlz4` backup files directly to HTML without intermediate steps
- **JSON support**: Also supports regular `.json` bookmark files
- **Automatic file detection**: Automatically detects the input file format
- **Preserves metadata**: Maintains bookmark timestamps and descriptions
- **Standard format**: Outputs HTML in the standard Netscape bookmark format

## Requirements

- Python 3.6+
- [`lz4`](https://pypi.org/project/lz4/) library for decompressing Firefox backup files

## Installation

1. Install the required Python package:
```bash
pip install -r requirements.txt
```

Or manually:
```bash
pip install lz4
```

2. Make the script executable (optional):
```bash
chmod +x ff_bookmarks.py
```

## Usage

### Basic usage
```bash
python ff_bookmarks.py input_file
```

### Specify output file
```bash
python ff_bookmarks.py input_file -o output_file.html
```

### Examples

Convert a Firefox jsonlz4 backup to HTML:
```bash
python ff_bookmarks.py bookmarks-2025-06-11_123456_randomhash.jsonlz4
```

Convert with custom output filename:
```bash
python ff_bookmarks.py bookmarks-2025-06-11_123456_randomhash.jsonlz4 -o my_bookmarks.html
```

Convert a regular JSON bookmark file:
```bash
python ff_bookmarks.py bookmarks.json -o bookmarks.html
```

## Input File Formats

### Firefox `jsonlz4` backup files
- Location: `~/.mozilla/firefox/[profile]/bookmarkbackups/`
- Format: `bookmarks-YYYY-MM-DD_HHMMSS_randomhash.jsonlz4`
- These are compressed backup files created automatically by Firefox

### JSON bookmark files
- Regular JSON files containing Firefox bookmark data
- Can be created by manually exporting bookmarks or by first decompressing `.jsonlz4` files

## Output Format

The script generates HTML files in the standard Netscape bookmark format, which can be imported into most web browsers including:
- Firefox
- Chrome/Chromium
- Safari
- Edge
- Opera

## Example

For a complete demonstration of the script's capabilities, you can run the included example:

```bash
chmod +x example/example_usage.sh

./example/example_usage.sh
```

This example script will:
- Check for required dependencies and install them if needed
- Show various usage examples with explanations
- Run a test conversion using the included sample bookmark file
- Generate a test HTML output file that you can examine or import into your browser

The example script also serves as a reference for different ways to use the converter in your own workflows.

## How it works

1. **File Detection**: The script automatically detects whether the input is a compressed `.jsonlz4` file or a regular JSON file
2. **Decompression** (if needed): For `.jsonlz4` files, it removes the Mozilla LZ4 header and decompresses the content
3. **JSON Parsing**: Parses the bookmark data structure
4. **HTML Generation**: Recursively converts the bookmark tree to HTML format, preserving:
   - Folder hierarchy
   - Bookmark URLs and titles
   - Creation and modification timestamps
   - Bookmark descriptions (if present)

## Error Handling

The script includes comprehensive error handling for:
- Missing or invalid input files
- Corrupted compression data
- Invalid JSON data
- File permission issues
- Missing dependencies

## Acknowledgments

This project combines and enhances code from two excellent projects:

- **[json2html-bookmarks](https://github.com/andreax79/json2html-bookmarks)** by [Andrea Bonomi](https://github.com/andreax79) - Provides the HTML generation functionality for converting Firefox bookmark JSON to standard HTML format
- **[jsonlz4_to_json](https://github.com/Robotvasya/jsonlz4_to_json)** by [Robotvasya](https://github.com/Robotvasya) - Provides the LZ4 decompression functionality for handling Firefox `.jsonlz4` backup files

All portions are licensed under the MIT License. See the script header for full license text.

## Troubleshooting

### "Please install the required module 'lz4'"
Install the lz4 package:
```bash
pip install lz4
```

### "not a valid Firefox bookmark backup file"
Ensure you're using a valid Firefox bookmark backup file. These are typically found in:
- Linux: `~/.mozilla/firefox/[profile]/bookmarkbackups/`
- Windows: `%APPDATA%\Mozilla\Firefox\Profiles\[profile]\bookmarkbackups\`
- macOS: `~/Library/Application Support/Firefox/Profiles/[profile]/bookmarkbackups/`

### File permission errors
Ensure you have read permissions for the input file and write permissions for the output directory.
