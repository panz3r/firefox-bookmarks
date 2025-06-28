# Firefox Bookmarks Backup to HTML Converter

A fast, dependency-free tool that converts Firefox bookmark backup files directly to HTML format.

## 🚀 Features

- **Direct conversion**: Convert Firefox `.jsonlz4` backup files directly to HTML without intermediate steps
- **JSON support**: Also supports regular `.json` bookmark files
- **Automatic file detection**: Automatically detects the input file format
- **Preserves metadata**: Maintains bookmark timestamps and descriptions
- **Standard format**: Outputs HTML in the standard Netscape bookmark format
- **Zero dependencies**: Single binary, no runtime requirements
- **Cross-platform**: Builds for Windows, macOS, and Linux

## 📥 Installation

### Pre-built Binaries (Recommended)

**[TBD]**

### Build from Source
```bash
# Clone and build
git clone git@github.com:panz3r/firefox-bookmarks.git

cd firefox-bookmarks

go build -o firefox-bookmarks

# Or use the build script for all platforms
./build.sh
```

## 🚀 Usage

### Basic Usage

```bash
# Show help
./firefox-bookmarks -help

# Convert with auto-generated output filename
./firefox-bookmarks backup.jsonlz4

# Convert with custom output filename
./firefox-bookmarks -o my_bookmarks.html backup.jsonlz4

# Convert JSON bookmark file
./firefox-bookmarks -o bookmarks.html bookmarks.json
```

**Note**: Flags must come before the input file.

## 📁 Input File Formats

### Firefox `.jsonlz4` backup files

#### Location
- **Linux/macOS**: `~/.mozilla/firefox/[profile]/bookmarkbackups/`
- **Windows**: `%APPDATA%\Mozilla\Firefox\Profiles\[profile]\bookmarkbackups\`

#### Format

- **Filename:** `bookmarks-YYYY-MM-DD_HHMMSS_randomhash.jsonlz4`
- These are compressed backup files created automatically by Firefox

### JSON bookmark files
- Regular JSON files containing Firefox bookmark data
- Can be created by manually exporting bookmarks

## 📄 Output Format

Generates HTML files in the standard Netscape bookmark format, compatible with:
- Firefox, Chrome, Safari, Edge, Opera
- Most bookmark management tools
- Other bookmark converters

## 🧠 How it works

1. **File Detection**: The tool automatically detects whether the input is a compressed `.jsonlz4` file or a regular JSON file
2. **Decompression** (if needed): For `.jsonlz4` files, it removes the Mozilla LZ4 header and decompresses the content
3. **JSON Parsing**: Parses the bookmark data structure
4. **HTML Generation**: Recursively converts the bookmark tree to HTML format, preserving:
   - Folder hierarchy
   - Bookmark URLs and titles
   - Creation and modification timestamps
   - Bookmark descriptions (if present)

## 🧪 Example

Run the included demonstration:

```bash
./example/example_usage.sh
```

This will:
- Build the binary if needed
- Show usage examples
- Run a test conversion
- Display performance metrics

## 🚀 Performance

- **Execution time**: ~5-10ms startup + processing
- **Memory usage**: ~8-12MB peak
- **Binary size**: ~2-3MB (no runtime dependencies)
- **Cross-platform**: Native binaries for all major platforms

## 🔧 Troubleshooting

### Binary not found or won't run
- Download the correct binary for your platform from `builds/`
- Make sure the binary has execute permissions: `chmod +x firefox-bookmarks`
- Build from source if needed: `go build -o firefox-bookmarks`

### "not a valid Firefox bookmark backup file"
Ensure you're using a valid Firefox bookmark backup file from:
- **Linux**: `~/.mozilla/firefox/[profile]/bookmarkbackups/`
- **Windows**: `%APPDATA%\Mozilla\Firefox\Profiles\[profile]\bookmarkbackups\`
- **macOS**: `~/Library/Application Support/Firefox/Profiles/[profile]/bookmarkbackups/`

### File permission errors
Ensure you have read permissions for the input file and write permissions for the output directory.

## 🐍 Legacy Python Version

The original Python implementation is still available in the `python/` directory. See [Python README](python/README.md) for details. The Python version requires Python 3.6+ and the `lz4` library, but produces identical output to the Go version.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

This project builds upon excellent work from:
- **[json2html-bookmarks](https://github.com/andreax79/json2html-bookmarks)** by [Andrea Bonomi](https://github.com/andreax79)
- **[jsonlz4_to_json](https://github.com/Robotvasya/jsonlz4_to_json)** by [Robotvasya](https://github.com/Robotvasya)
