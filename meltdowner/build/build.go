package build
import (
	"os"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/wassan128/meltdowner/meltdowner/file"
	"github.com/wassan128/meltdowner/meltdowner/parser"
	"gopkg.in/russross/blackfriday.v2"
)

func setHTMLFlags() blackfriday.HTMLFlags {
	htmlFlags := blackfriday.CommonHTMLFlags
	htmlFlags |= blackfriday.FootnoteReturnLinks
	htmlFlags |= blackfriday.SmartypantsAngledQuotes
	htmlFlags |= blackfriday.SmartypantsQuotesNBSP

	return htmlFlags
}

func getRenderer() *blackfriday.HTMLRenderer {
	htmlFlags := setHTMLFlags()
	return blackfriday.NewHTMLRenderer(
		blackfriday.HTMLRendererParameters{
			Flags: htmlFlags,
			Title: "",
			CSS: "",
		},
	)
}

func setExtensionFlags() blackfriday.Extensions {
	extFlags := blackfriday.CommonExtensions
	extFlags |= blackfriday.Footnotes
	extFlags |= blackfriday.HeadingIDs
	extFlags |= blackfriday.Titleblock
	extFlags |= blackfriday.DefinitionLists

	return extFlags
}

func md2HTML(md []byte, renderer *blackfriday.HTMLRenderer) string {
	extFlags := setExtensionFlags()
	html := blackfriday.Run(md,
		blackfriday.WithExtensions(extFlags),
		blackfriday.WithRenderer(renderer))
	return string(html)
}

func checkDirectoryExistence() bool {
	if !file.IsExistPath("source") {
		fmt.Println("[Error] source/ does not exists.")
		return false
	}
	if !file.IsExistPath("theme") {
		fmt.Println("[Error] theme/ does not exists.")
		return false
	}
	return true
}

func concatTemplates(content string) string {
	header := string(file.LoadFileContents("theme/template/header.html"))
	footer := string(file.LoadFileContents("theme/template/footer.html"))

	html := strings.Join([]string{header, content, footer}, "\n")

	return html
}

func createPostDir(publicDir string, createdAt parser.CreatedAt) string {
	var postPath string

	paths := []string{publicDir, createdAt.Year, createdAt.Month, createdAt.Date}
	for _, path := range paths {
		postPath = filepath.Join(postPath, path)
		if _, err := os.Stat(postPath); err != nil {
			file.CreateDir(postPath)
		}
	}

	return postPath
}

func generatePosts(renderer *blackfriday.HTMLRenderer, mds []string) []parser.Post {
	var posts []parser.Post

	for _, mdPath := range mds {
		fmt.Println("[*] Start: ", mdPath)
		md := file.LoadFileContents(mdPath)
		if md == nil {
			fmt.Println("markdown load error")
			return nil
		}

		post := parser.ParseMarkdown(md)
		posts = append(posts, *post)

		content := md2HTML(post.Body, renderer)

		htmlString := concatTemplates(content)
		htmlFile := file.CreateFile("index.html")
		defer htmlFile.Close()

		postPath := createPostDir("public", post.Header.Date)

		file.MoveFile("index.html", filepath.Join(postPath, "index.html"))
		fmt.Fprintln(htmlFile, htmlString)
		fmt.Println("[*] Done: ", postPath)
	}

	return posts
}

func generateTopPage(renderer *blackfriday.HTMLRenderer, posts []parser.Post) {
	var mdTop string
	for _, post := range posts {
		mdTop += fmt.Sprintf("* [%s](/%s/%s/%s)\n", post.Header.Title,
			post.Header.Date.Year, post.Header.Date.Month, post.Header.Date.Date)
	}

	content := md2HTML([]byte(mdTop), renderer)
	htmlString := concatTemplates(content)
	htmlFile := file.CreateFile("index.html")
	defer htmlFile.Close()

	file.MoveFile("index.html", filepath.Join("public", "index.html"))
	fmt.Fprintln(htmlFile, htmlString)
}

func Run() {
	// check `source/` and `template/` existence.
	if ok := checkDirectoryExistence(); !ok {
		fmt.Println("exit.")
		return
	}

	renderer := getRenderer()
	mds := file.GetMarkdownPaths("source")

	posts := generatePosts(renderer, mds)
	generateTopPage(renderer, posts)

	file.CopyDir(filepath.Join("theme", "css"), filepath.Join("public", "css"))
}

