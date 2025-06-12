#!/usr/bin/env python3
# -*- coding: utf-8 -*-
#
# Convert Firefox bookmarks from jsonlz4 backup format directly to HTML format
#
# This script combines the functionality of:
# - jsonlz4_conv.py (converts .jsonlz4 to JSON) - Copyright (c) 2022 Robotvasya
#   Source: https://github.com/Robotvasya/jsonlz4_to_json
# - json2html.py (converts JSON bookmarks to HTML) - Copyright (c) 2013 Andrea Bonomi
#   Source: https://github.com/andreax79/json2html-bookmarks
#
# Copyright (c) 2013 Andrea Bonomi - andrea.bonomi@gmail.com (json2html portions)
# Copyright (c) 2022 Robotvasya (jsonlz4 conversion portions)
# Copyright (c) 2025 Mattia Panzeri - Enhanced version merging both functionalities
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.
#

import sys
import math
import pathlib
import argparse
import json
from typing import Dict, Any, TextIO

try:
    import lz4.block
except ImportError:
    sys.stderr.write("Error: Please install the required module 'lz4'.\n")
    sys.stderr.write("You can install it with: pip install lz4\n")
    sys.exit(1)

# Constants
FIREFOX_LZ4_SIGNATURE = b"mozLz4"
FIREFOX_LZ4_HEADER_SIZE = 12
DEFAULT_BUFFER_SIZE = 10 * 1024 * 1024  # 10MB
INDENT_SIZE = 4
BOOKMARK_SEPARATOR_TYPE = 3


def err(message: str = None) -> None:
    """Handle errors and exit."""
    if message:
        sys.stderr.write(f"{sys.argv[0]}: {message}\n")
    else:
        e = sys.exc_info()[1]
        sys.stderr.write(f"{sys.argv[0]}: {str(e)}\n")
    sys.exit(1)


def html_escape(text: str) -> str:
    """Escape HTML special characters to prevent XSS and display issues."""
    if not text:
        return ""
    return (text
            .replace("&", "&amp;")
            .replace("<", "&lt;")
            .replace(">", "&gt;")
            .replace('"', "&quot;")
            .replace("'", "&#x27;"))


def is_valid_jsonlz4_file(filename: str) -> bool:
    """
    Check if the file is a valid Firefox jsonlz4 bookmark backup file.

    :param filename: Path to file
    :return: True if format is valid jsonlz4
    """
    if not pathlib.Path(filename).is_file():
        return False

    try:
        with open(filename, 'rb') as inf:
            file_header = inf.read(len(FIREFOX_LZ4_SIGNATURE))
        return file_header == FIREFOX_LZ4_SIGNATURE
    except OSError:
        return False


def is_json_file(filename: str) -> bool:
    """
    Check if the file is a JSON file by trying to parse it.

    :param filename: Path to file
    :return: True if it's a valid JSON file
    """
    if not pathlib.Path(filename).is_file():
        return False

    try:
        with open(filename, 'r', encoding='utf-8') as f:
            json.load(f)
        return True
    except (json.JSONDecodeError, OSError):
        return False


def decompress_jsonlz4(filename: str) -> Dict[str, Any]:
    """
    Decompress a Firefox jsonlz4 bookmark backup file and return the JSON data.

    :param filename: Path to the .jsonlz4 file
    :return: Parsed JSON data as dictionary
    """
    try:
        with open(filename, 'rb') as inFile:
            # Skip the Firefox LZ4 header
            inFile.seek(FIREFOX_LZ4_HEADER_SIZE)
            compressed_data = inFile.read()

            # Decompress the data with larger buffer size
            decompressed_data = lz4.block.decompress(
                compressed_data, uncompressed_size=DEFAULT_BUFFER_SIZE)

            # Parse JSON
            return json.loads(decompressed_data.decode('utf-8'))

    except lz4.block.LZ4BlockError as e:
        err(f"LZ4 decompression error: {e}")
    except json.JSONDecodeError as e:
        err(f"JSON parsing error: {e}")
    except OSError as e:
        err(f"File reading error: {e}")


def load_json_file(filename: str) -> Dict[str, Any]:
    """
    Load and parse a regular JSON file.

    :param filename: Path to the JSON file
    :return: Parsed JSON data as dictionary
    """
    try:
        with open(filename, "r", encoding="utf-8") as f:
            return json.load(f)
    except (json.JSONDecodeError, OSError) as e:
        err(f"Error loading JSON file: {e}")


def write_indented(output: TextIO, indent: int, text: str) -> None:
    """Write indented text to output stream."""
    indentation = " " * (INDENT_SIZE * indent)
    output.write(f"{indentation}{text}\n")


def convert_firefox_timestamp(timestamp: Any) -> str:
    """
    Convert Firefox timestamp to Unix timestamp string.

    :param timestamp: Firefox timestamp (microseconds since epoch)
    :return: Unix timestamp string or empty string if conversion fails
    """
    try:
        if timestamp:
            return str(int(math.floor(int(timestamp) / 1000000)))
    except (ValueError, TypeError):
        pass
    return ""


def format_date_attributes(data: Dict[str, Any]) -> str:
    """
    Format date attributes for HTML bookmark tags.

    :param data: Bookmark data dictionary
    :return: Formatted date attributes string
    """
    attributes = []

    if 'dateAdded' in data:
        date_added = convert_firefox_timestamp(data['dateAdded'])
        if date_added:
            attributes.append(f' ADD_DATE="{date_added}"')

    if 'lastModified' in data:
        last_modified = convert_firefox_timestamp(data['lastModified'])
        if last_modified:
            attributes.append(f' LAST_MODIFIED="{last_modified}"')

    return ''.join(attributes)


def write_html_header(output: TextIO, title: str) -> None:
    """Write the HTML document header."""
    header = f"""<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
    It will be read and overwritten.
    DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>{html_escape(title)}</H1>
<DL><p>"""
    write_indented(output, 0, header)


def write_folder(output: TextIO, data: Dict[str, Any], indent: int) -> None:
    """Write a bookmark folder to HTML."""
    title = html_escape(data.get('title', ''))
    date_attrs = format_date_attributes(data)

    write_indented(output, indent, f'<DT><H3{date_attrs}>{title}</H3>')
    write_indented(output, indent, '<DL><p>')


def write_bookmark(output: TextIO, data: Dict[str, Any], indent: int) -> None:
    """Write a single bookmark to HTML."""
    uri = data.get('uri', '')
    title = html_escape(data.get('title', uri))
    date_attrs = format_date_attributes(data)

    write_indented(
        output, indent, f'<DT><A HREF="{html_escape(uri)}"{date_attrs}>{title}</A>')

    # Handle bookmark descriptions
    annos = data.get('annos')
    if isinstance(annos, list):
        for anno in annos:
            if (isinstance(anno, dict) and
                    anno.get('name') == 'bookmarkProperties/description'):
                description = html_escape(anno.get('value', ''))
                write_indented(output, indent, f'<DD>{description}')


def convert_bookmarks_to_html(output: TextIO, data: Dict[str, Any], indent: int = 0) -> None:
    """
    Convert bookmark data to HTML format recursively.

    :param output: Output file handle
    :param data: Bookmark data dictionary
    :param indent: Current indentation level
    """
    children = data.get('children')
    uri = data.get('uri')

    # Handle containers (folders) with children
    if children is not None and isinstance(children, list):
        if indent == 0:
            # Output the main header
            title = data.get('title', 'Bookmarks Menu')
            write_html_header(output, title)
        else:
            # Output a folder
            write_folder(output, data, indent)

        # Process children
        for child in children:
            # Skip separators (typeCode 3)
            if child.get('typeCode') == BOOKMARK_SEPARATOR_TYPE:
                continue
            convert_bookmarks_to_html(output, child, indent + 1)

        write_indented(output, indent, '</DL><p>')

    elif uri is not None:
        # Output a bookmark
        write_bookmark(output, data, indent)


def main() -> None:
    """Main function to handle command line arguments and orchestrate the conversion."""
    cli_usage_description = '''
    ff_bookmarks.py input_file [-o OUTPUT_FILE] 
    
    Converts Firefox bookmark backup files to HTML format.
    Supports both .jsonlz4 (compressed backup) and .json (uncompressed) input files.
    
    Examples:
        ff_bookmarks.py bookmarks-2025-06-11.jsonlz4
        ff_bookmarks.py bookmarks-2025-06-11.jsonlz4 -o my_bookmarks.html
        ff_bookmarks.py bookmarks.json -o bookmarks.html
    
or help:  
    ff_bookmarks.py --help
    '''

    parser = argparse.ArgumentParser(
        usage=cli_usage_description,
        description="Converts Firefox bookmark backup files from .jsonlz4 or .json format to HTML format",
        formatter_class=argparse.RawDescriptionHelpFormatter
    )

    parser.add_argument(
        'input_file',
        help="Path to Firefox bookmark backup file (.jsonlz4 or .json)"
    )

    parser.add_argument(
        '-o', '--output',
        dest='output_file',
        help="Path to output HTML file. If omitted, uses input filename with .html extension",
        required=False
    )

    args = parser.parse_args()

    # Check if input file exists
    input_path = pathlib.Path(args.input_file)
    if not input_path.is_file():
        err(f"Input file '{args.input_file}' does not exist.")

    # Determine output filename
    if args.output_file:
        output_path = pathlib.Path(args.output_file)
    else:
        output_path = input_path.with_suffix(".html")

    try:
        # Determine file type and load data accordingly
        if is_valid_jsonlz4_file(args.input_file):
            print(
                f"Processing Firefox jsonlz4 bookmark backup: {args.input_file}")
            bookmark_data = decompress_jsonlz4(args.input_file)
        elif is_json_file(args.input_file):
            print(f"Processing JSON bookmark file: {args.input_file}")
            bookmark_data = load_json_file(args.input_file)
        else:
            err(f"'{args.input_file}' is not a valid Firefox bookmark backup file (.jsonlz4) or JSON file.")

        # Convert to HTML
        print("Converting bookmarks to HTML format...")
        with open(output_path, "w", encoding='utf-8') as output:
            convert_bookmarks_to_html(output, bookmark_data)

        print(f"Successfully converted bookmarks to: {output_path}")

    except Exception as e:
        err(f"Conversion failed: {e}")


if __name__ == '__main__':
    main()
