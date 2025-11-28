package reader

// BookID uniquely identifies a book within the application.
type BookID string

// Position represents a logical location within a book, independent
// of any particular UI representation.
type Position struct {
	ChapterIndex    int
	OffsetInChapter int
}

// Chapter models a logical chapter or section within a book.
type Chapter struct {
	Index  int
	Title  string
	Offset int // Start offset of the chapter within the linearized text stream.
	Length int // Length of the chapter in runes or characters.
}

// Book represents a logical book with metadata and an ordered list
// of chapters or sections.
type Book struct {
	ID       BookID
	Title    string
	Author   string
	Chapters []Chapter

	// TotalCharacters is an optional aggregate aiding in percentage
	// calculations for navigation and progress display.
	TotalCharacters int
}

// Locatable is implemented by types that can expose a Position
// within a book.
type Locatable interface {
	GetPosition() Position
}

// Bookmark represents a named location within a specific book.
type Bookmark struct {
	Name   string
	BookID BookID
	Pos    Position
}

// GetPosition returns the position associated with the bookmark.
func (b Bookmark) GetPosition() Position {
	return b.Pos
}

// TOCEntry represents an entry in the book's table of contents.
type TOCEntry struct {
	Label  string
	BookID BookID
	Pos    Position
}

// GetPosition returns the position associated with the TOC entry.
func (e TOCEntry) GetPosition() Position {
	return e.Pos
}
