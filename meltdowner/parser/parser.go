package parser

import (
    "strings"
    "strconv"
    "regexp"

    "github.com/wassan128/meltdowner/meltdowner/util"
)

const (
    YEAR int = iota
    MONTH
    DATE
)

type CreatedAt struct {
    Year string
    Month string
    Date string
}

type Header struct {
    Title string
    Date CreatedAt
    Tags []string
    Id string
    Public bool
}
type Post struct {
    Header Header
    Body []byte
}

func valiDate(pattern, date string) bool {
    return regexp.MustCompile(pattern).Match([]byte(date))
}

func parseDate(dateStr string) CreatedAt {
    date := strings.Split(dateStr, "-")

    util.ExitIfFalse(valiDate("^\\d{4}$", date[YEAR]) &&
        valiDate("^\\d{2}$", date[MONTH]) &&
        valiDate("^\\d{2}$", date[DATE]))

    var createdAt CreatedAt
    createdAt.Year = date[YEAR]
    createdAt.Month = date[MONTH]
    createdAt.Date = date[DATE]

    return createdAt
}

func newPost() *Post {
    post := new(Post)

    post.Header.Title = "Untitled"
    post.Header.Date = CreatedAt{Year: "1970", Month: "01", Date: "01"}
    post.Header.Public = true

    return post
}

func ParseMarkdown(markdown []byte) *Post {
    lines := strings.Split(string(markdown), "\n")

    post := newPost()

    cur := 0
    for ; lines[cur] != "---"; cur++ {
        switch header := strings.Split(lines[cur], ":"); header[0] {
        case "title":
            post.Header.Title = strings.TrimSpace(header[1])
        case "date":
            post.Header.Date = parseDate(strings.TrimSpace(header[1]))
        case "tags":
            header[1] = strings.Replace(header[1], " ", "", -1)
            if len(header[1]) > 0 {
                post.Header.Tags = strings.Split(header[1], ",")
            }
        case "public":
            post.Header.Public, _ = strconv.ParseBool(strings.TrimSpace(header[1]))
        }
    }
    cur++

    post.Body = []byte(strings.Join(lines[cur:], "\n"))

    return post
}

