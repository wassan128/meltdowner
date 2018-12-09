package parser

import (
	"strings"
)

type Header struct {
	Title string
	CreatedAt string
	Tags []string
}
type Post struct {
	Header Header
	Body []byte
}

func ParseMarkdown(markdown []byte) *Post {
	lines := strings.Split(string(markdown), "\n")

	var post Post
	cur := 0
	for ; lines[cur] != "---"; cur++ {
		switch header := strings.Split(lines[cur], ":"); header[0] {
		case "title":
			post.Header.Title = strings.TrimSpace(header[1])
		case "date":
			post.Header.CreatedAt = strings.TrimSpace(header[1])
		}
	}
	cur++

	post.Body = []byte(strings.Join(lines[cur:], "\n"))

	return &post
}

