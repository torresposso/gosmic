package services

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Chapter struct {
	ID    string // Filename without extension (e.g. "01-pre-flight-check")
	Title string // Readable title
}

type DocService struct {
	chaptersDir string
	mdParser    goldmark.Markdown
}

func NewDocService(chaptersDir string) *DocService {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	return &DocService{
		chaptersDir: chaptersDir,
		mdParser:    md,
	}
}

// GetChapters returns a sorted list of all available chapters
func (s *DocService) GetChapters() ([]Chapter, error) {
	absPath, _ := filepath.Abs(s.chaptersDir)
	fmt.Printf("DEBUG: Reading chapters from: %s (Abs: %s)\n", s.chaptersDir, absPath)

	files, err := os.ReadDir(s.chaptersDir)
	if err != nil {
		fmt.Printf("DEBUG: Error reading directory: %v\n", err)
		return nil, err
	}
	fmt.Printf("DEBUG: Found %d files\n", len(files))

	var chapters []Chapter
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".md" {
			name := strings.TrimSuffix(f.Name(), ".md")
			// Beautify title: "01-pre-flight-check" -> "01 Pre Flight Check"
			// Ideally we would parse the H1 from file, but this is faster for now
			title := strings.ReplaceAll(name, "-", " ")
			title = strings.Title(title)

			chapters = append(chapters, Chapter{
				ID:    name,
				Title: title,
			})
		}
	}

	// Ensure sorted by filename
	sort.Slice(chapters, func(i, j int) bool {
		return chapters[i].ID < chapters[j].ID
	})

	return chapters, nil
}

// GetChapterContent returns the HTML content of a specific chapter
func (s *DocService) GetChapterContent(chapterID string) (string, error) {
	// Security check to prevent directory traversal
	if strings.Contains(chapterID, "..") || strings.Contains(chapterID, "/") {
		return "", fmt.Errorf("invalid chapter ID")
	}

	path := filepath.Join(s.chaptersDir, chapterID+".md")
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := s.mdParser.Convert(content, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
