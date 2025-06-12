package bookmarks

import (
	"fmt"
	"html"
	"io"
	"math"
	"strings"
)

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

// HTMLConverter handles conversion of bookmark data to HTML format
type HTMLConverter struct{}

// NewHTMLConverter creates a new HTMLConverter
func NewHTMLConverter() *HTMLConverter {
	return &HTMLConverter{}
}

// htmlEscape escapes HTML special characters to prevent XSS and display issues
func (hc *HTMLConverter) htmlEscape(text string) string {
	if text == "" {
		return ""
	}
	return html.EscapeString(text)
}

// convertFirefoxTimestamp converts Firefox timestamp to Unix timestamp string
func (hc *HTMLConverter) convertFirefoxTimestamp(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	return fmt.Sprintf("%d", int64(math.Floor(float64(timestamp)/1000000)))
}

// formatDateAttributes formats date attributes for HTML bookmark tags
func (hc *HTMLConverter) formatDateAttributes(data *BookmarkData) string {
	var attributes []string

	if data.DateAdded != 0 {
		dateAdded := hc.convertFirefoxTimestamp(data.DateAdded)
		if dateAdded != "" {
			attributes = append(attributes, fmt.Sprintf(` ADD_DATE="%s"`, dateAdded))
		}
	}

	if data.LastModified != 0 {
		lastModified := hc.convertFirefoxTimestamp(data.LastModified)
		if lastModified != "" {
			attributes = append(attributes, fmt.Sprintf(` LAST_MODIFIED="%s"`, lastModified))
		}
	}

	return strings.Join(attributes, "")
}

// writeHTMLHeader writes the HTML document header
func (hc *HTMLConverter) writeHTMLHeader(writer *HTMLWriter, title string) error {
	header := fmt.Sprintf(`<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
    It will be read and overwritten.
    DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>%s</H1>
<DL><p>`, hc.htmlEscape(title))

	return writer.WriteIndented(0, header)
}

// writeFolder writes a bookmark folder to HTML
func (hc *HTMLConverter) writeFolder(writer *HTMLWriter, data *BookmarkData, indent int) error {
	title := hc.htmlEscape(data.Title)
	dateAttrs := hc.formatDateAttributes(data)

	err := writer.WriteIndented(indent, fmt.Sprintf(`<DT><H3%s>%s</H3>`, dateAttrs, title))
	if err != nil {
		return err
	}
	return writer.WriteIndented(indent, `<DL><p>`)
}

// writeBookmark writes a single bookmark to HTML
func (hc *HTMLConverter) writeBookmark(writer *HTMLWriter, data *BookmarkData, indent int) error {
	uri := data.URI
	title := data.Title
	if title == "" {
		title = uri
	}
	title = hc.htmlEscape(title)
	dateAttrs := hc.formatDateAttributes(data)

	err := writer.WriteIndented(indent,
		fmt.Sprintf(`<DT><A HREF="%s"%s>%s</A>`, hc.htmlEscape(uri), dateAttrs, title))
	if err != nil {
		return err
	}

	// Handle bookmark descriptions
	for _, anno := range data.Annotations {
		if anno.Name == "bookmarkProperties/description" {
			description := hc.htmlEscape(anno.Value)
			err = writer.WriteIndented(indent, fmt.Sprintf(`<DD>%s`, description))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ConvertToHTML converts bookmark data to HTML format recursively
func (hc *HTMLConverter) ConvertToHTML(writer *HTMLWriter, data *BookmarkData, indent int) error {
	// Handle containers (folders) with children
	if data.Children != nil {
		if indent == 0 {
			// Output the main header
			title := data.Title
			if title == "" {
				title = "Bookmarks Menu"
			}
			err := hc.writeHTMLHeader(writer, title)
			if err != nil {
				return err
			}
		} else {
			// Output a folder
			err := hc.writeFolder(writer, data, indent)
			if err != nil {
				return err
			}
		}

		// Process children (if any)
		for _, child := range data.Children {
			// Skip separators (typeCode 3)
			if child.TypeCode == BookmarkSeparatorType {
				continue
			}
			err := hc.ConvertToHTML(writer, &child, indent+1)
			if err != nil {
				return err
			}
		}

		return writer.WriteIndented(indent, `</DL><p>`)
	} else if data.URI != "" {
		// Output a bookmark
		return hc.writeBookmark(writer, data, indent)
	}

	return nil
}

// ConvertBookmarksToHTML is a convenience function that creates a converter and converts bookmarks
func ConvertBookmarksToHTML(writer io.Writer, data *BookmarkData) error {
	converter := NewHTMLConverter()
	htmlWriter := NewHTMLWriter(writer)
	return converter.ConvertToHTML(htmlWriter, data, 0)
}
