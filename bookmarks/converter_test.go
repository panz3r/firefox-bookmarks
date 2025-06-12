package bookmarks

import (
	"bytes"
	"strings"
	"testing"
)

func TestHTMLWriter_WriteIndented(t *testing.T) {
	tests := []struct {
		name   string
		indent int
		text   string
		want   string
	}{
		{
			name:   "NoIndent",
			indent: 0,
			text:   "Hello World",
			want:   "Hello World\n",
		},
		{
			name:   "SingleIndent",
			indent: 1,
			text:   "Hello World",
			want:   "    Hello World\n",
		},
		{
			name:   "DoubleIndent",
			indent: 2,
			text:   "Hello World",
			want:   "        Hello World\n",
		},
		{
			name:   "EmptyText",
			indent: 1,
			text:   "",
			want:   "    \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewHTMLWriter(&buf)

			err := writer.WriteIndented(tt.indent, tt.text)
			if err != nil {
				t.Fatalf("WriteIndented failed: %v", err)
			}

			got := buf.String()
			if got != tt.want {
				t.Errorf("WriteIndented() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHTMLConverter_htmlEscape(t *testing.T) {
	converter := NewHTMLConverter()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "NoEscaping",
			input: "Hello World",
			want:  "Hello World",
		},
		{
			name:  "EscapeAmpersand",
			input: "Tom & Jerry",
			want:  "Tom &amp; Jerry",
		},
		{
			name:  "EscapeLessThan",
			input: "a < b",
			want:  "a &lt; b",
		},
		{
			name:  "EscapeGreaterThan",
			input: "a > b",
			want:  "a &gt; b",
		},
		{
			name:  "EscapeQuotes",
			input: `Say "Hello"`,
			want:  "Say &#34;Hello&#34;",
		},
		{
			name:  "EmptyString",
			input: "",
			want:  "",
		},
		{
			name:  "MultipleSpecialChars",
			input: `<script>alert("XSS & stuff");</script>`,
			want:  "&lt;script&gt;alert(&#34;XSS &amp; stuff&#34;);&lt;/script&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := converter.htmlEscape(tt.input)
			if got != tt.want {
				t.Errorf("htmlEscape() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHTMLConverter_convertFirefoxTimestamp(t *testing.T) {
	converter := NewHTMLConverter()

	tests := []struct {
		name      string
		timestamp int64
		want      string
	}{
		{
			name:      "ZeroTimestamp",
			timestamp: 0,
			want:      "",
		},
		{
			name:      "ValidTimestamp",
			timestamp: 1639123456789000,
			want:      "1639123456",
		},
		{
			name:      "AnotherValidTimestamp",
			timestamp: 1234567890123000,
			want:      "1234567890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := converter.convertFirefoxTimestamp(tt.timestamp)
			if got != tt.want {
				t.Errorf("convertFirefoxTimestamp() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHTMLConverter_formatDateAttributes(t *testing.T) {
	converter := NewHTMLConverter()

	tests := []struct {
		name string
		data *BookmarkData
		want string
	}{
		{
			name: "NoTimestamps",
			data: &BookmarkData{
				Title: "Test",
				URI:   "https://example.com",
			},
			want: "",
		},
		{
			name: "OnlyDateAdded",
			data: &BookmarkData{
				Title:     "Test",
				URI:       "https://example.com",
				DateAdded: 1639123456789000,
			},
			want: ` ADD_DATE="1639123456"`,
		},
		{
			name: "OnlyLastModified",
			data: &BookmarkData{
				Title:        "Test",
				URI:          "https://example.com",
				LastModified: 1639123456789000,
			},
			want: ` LAST_MODIFIED="1639123456"`,
		},
		{
			name: "BothTimestamps",
			data: &BookmarkData{
				Title:        "Test",
				URI:          "https://example.com",
				DateAdded:    1639123456789000,
				LastModified: 1639123456789000,
			},
			want: ` ADD_DATE="1639123456" LAST_MODIFIED="1639123456"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := converter.formatDateAttributes(tt.data)
			if got != tt.want {
				t.Errorf("formatDateAttributes() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHTMLConverter_writeHTMLHeader(t *testing.T) {
	converter := NewHTMLConverter()

	tests := []struct {
		name  string
		title string
		want  string
	}{
		{
			name:  "SimpleTitle",
			title: "My Bookmarks",
			want: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
    It will be read and overwritten.
    DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>My Bookmarks</H1>
<DL><p>
`,
		},
		{
			name:  "TitleWithSpecialChars",
			title: "Tom & Jerry's <Bookmarks>",
			want: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
    It will be read and overwritten.
    DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Tom &amp; Jerry&#39;s &lt;Bookmarks&gt;</H1>
<DL><p>
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewHTMLWriter(&buf)

			err := converter.writeHTMLHeader(writer, tt.title)
			if err != nil {
				t.Fatalf("writeHTMLHeader failed: %v", err)
			}

			got := buf.String()
			if got != tt.want {
				t.Errorf("writeHTMLHeader() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHTMLConverter_writeFolder(t *testing.T) {
	converter := NewHTMLConverter()

	tests := []struct {
		name   string
		data   *BookmarkData
		indent int
		want   string
	}{
		{
			name: "SimpleFolderNoTimestamps",
			data: &BookmarkData{
				Title: "Development",
			},
			indent: 1,
			want: `    <DT><H3>Development</H3>
    <DL><p>
`,
		},
		{
			name: "FolderWithTimestamps",
			data: &BookmarkData{
				Title:        "Development",
				DateAdded:    1639123456789000,
				LastModified: 1639123456789000,
			},
			indent: 1,
			want: `    <DT><H3 ADD_DATE="1639123456" LAST_MODIFIED="1639123456">Development</H3>
    <DL><p>
`,
		},
		{
			name: "FolderWithSpecialChars",
			data: &BookmarkData{
				Title: "Tom & Jerry's <Folder>",
			},
			indent: 2,
			want: `        <DT><H3>Tom &amp; Jerry&#39;s &lt;Folder&gt;</H3>
        <DL><p>
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewHTMLWriter(&buf)

			err := converter.writeFolder(writer, tt.data, tt.indent)
			if err != nil {
				t.Fatalf("writeFolder failed: %v", err)
			}

			got := buf.String()
			if got != tt.want {
				t.Errorf("writeFolder() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHTMLConverter_writeBookmark(t *testing.T) {
	converter := NewHTMLConverter()

	tests := []struct {
		name   string
		data   *BookmarkData
		indent int
		want   string
	}{
		{
			name: "SimpleBookmark",
			data: &BookmarkData{
				Title: "GitHub",
				URI:   "https://github.com",
			},
			indent: 1,
			want: `    <DT><A HREF="https://github.com">GitHub</A>
`,
		},
		{
			name: "BookmarkWithTimestamps",
			data: &BookmarkData{
				Title:        "GitHub",
				URI:          "https://github.com",
				DateAdded:    1639123456789000,
				LastModified: 1639123456789000,
			},
			indent: 1,
			want: `    <DT><A HREF="https://github.com" ADD_DATE="1639123456" LAST_MODIFIED="1639123456">GitHub</A>
`,
		},
		{
			name: "BookmarkWithDescription",
			data: &BookmarkData{
				Title: "Example",
				URI:   "https://example.com",
				Annotations: []Annotation{
					{Name: "bookmarkProperties/description", Value: "This is an example"},
				},
			},
			indent: 1,
			want: `    <DT><A HREF="https://example.com">Example</A>
    <DD>This is an example
`,
		},
		{
			name: "BookmarkNoTitle",
			data: &BookmarkData{
				URI: "https://example.com",
			},
			indent: 1,
			want: `    <DT><A HREF="https://example.com">https://example.com</A>
`,
		},
		{
			name: "BookmarkWithSpecialChars",
			data: &BookmarkData{
				Title: "Tom & Jerry's <Site>",
				URI:   "https://example.com?q=test&value=1",
			},
			indent: 1,
			want: `    <DT><A HREF="https://example.com?q=test&amp;value=1">Tom &amp; Jerry&#39;s &lt;Site&gt;</A>
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewHTMLWriter(&buf)

			err := converter.writeBookmark(writer, tt.data, tt.indent)
			if err != nil {
				t.Fatalf("writeBookmark failed: %v", err)
			}

			got := buf.String()
			if got != tt.want {
				t.Errorf("writeBookmark() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHTMLConverter_ConvertToHTML(t *testing.T) {
	converter := NewHTMLConverter()

	// Test with sample bookmark structure
	testData := &BookmarkData{
		Title: "Bookmarks Menu",
		Children: []BookmarkData{
			{
				Title: "Development",
				Children: []BookmarkData{
					{
						Title: "GitHub",
						URI:   "https://github.com",
					},
					{
						Title: "Stack Overflow",
						URI:   "https://stackoverflow.com",
					},
				},
			},
			{
				Title: "Example Bookmark",
				URI:   "https://example.com",
				Annotations: []Annotation{
					{Name: "bookmarkProperties/description", Value: "This is an example"},
				},
			},
			{
				TypeCode: BookmarkSeparatorType, // Should be skipped
			},
		},
	}

	var buf bytes.Buffer
	writer := NewHTMLWriter(&buf)

	err := converter.ConvertToHTML(writer, testData, 0)
	if err != nil {
		t.Fatalf("ConvertToHTML failed: %v", err)
	}

	result := buf.String()

	// Check that the output contains expected elements
	expectedParts := []string{
		"<!DOCTYPE NETSCAPE-Bookmark-file-1>",
		"<H1>Bookmarks Menu</H1>",
		"<DT><H3>Development</H3>",
		`<DT><A HREF="https://github.com">GitHub</A>`,
		`<DT><A HREF="https://stackoverflow.com">Stack Overflow</A>`,
		`<DT><A HREF="https://example.com">Example Bookmark</A>`,
		"<DD>This is an example",
		"</DL><p>",
	}

	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("Expected output to contain %q, but it didn't.\nFull output:\n%s", part, result)
		}
	}

	// Check that separator is not in output
	if strings.Contains(result, "separator") {
		t.Error("Output should not contain separator elements")
	}
}

func TestConvertBookmarksToHTML(t *testing.T) {
	// Test the convenience function
	testData := &BookmarkData{
		Title: "Test Bookmarks",
		Children: []BookmarkData{
			{
				Title: "Test Bookmark",
				URI:   "https://example.com",
			},
		},
	}

	var buf bytes.Buffer
	err := ConvertBookmarksToHTML(&buf, testData)
	if err != nil {
		t.Fatalf("ConvertBookmarksToHTML failed: %v", err)
	}

	result := buf.String()

	// Basic check that it produced HTML output
	if !strings.Contains(result, "<!DOCTYPE NETSCAPE-Bookmark-file-1>") {
		t.Error("Expected HTML DOCTYPE in output")
	}
	if !strings.Contains(result, "Test Bookmarks") {
		t.Error("Expected title in output")
	}
	if !strings.Contains(result, "https://example.com") {
		t.Error("Expected bookmark URL in output")
	}
}
