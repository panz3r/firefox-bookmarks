package bookmarks

// Constants from the Python version
const (
	FirefoxLZ4Signature   = "mozLz4"
	FirefoxLZ4HeaderSize  = 12
	DefaultBufferSize     = 10 * 1024 * 1024 // 10MB
	IndentSize            = 4
	BookmarkSeparatorType = 3
)

// BookmarkData represents the structure of bookmark data
type BookmarkData struct {
	Title        string         `json:"title,omitempty"`
	URI          string         `json:"uri,omitempty"`
	Children     []BookmarkData `json:"children,omitempty"`
	DateAdded    int64          `json:"dateAdded,omitempty"`
	LastModified int64          `json:"lastModified,omitempty"`
	TypeCode     int            `json:"typeCode,omitempty"`
	Annotations  []Annotation   `json:"annos,omitempty"`
}

// Annotation represents bookmark annotations (like descriptions)
type Annotation struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
