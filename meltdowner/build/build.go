package build

import (
	"os"
	"fmt"
	"path/filepath"
	"text/template"
	"strings"

	"github.com/wassan128/meltdowner/meltdowner/config"
	"github.com/wassan128/meltdowner/meltdowner/file"
	"github.com/wassan128/meltdowner/meltdowner/parser"
	"gopkg.in/russross/blackfriday.v2"
)

var Config config.Config = config.GetConfig()

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
	headerHtml := file.CreateFile("theme/template/header.html")
	headerTmpl := template.Must(template.ParseFiles("theme/template/header.tmpl"))
	headerTmpl.Execute(headerHtml, Config.Blog)

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

		title := []byte(fmt.Sprintf("# %s\n", post.Header.Title))
		content := md2HTML(append(title, post.Body...), renderer)

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
	mdTop := "<ul class='top'>\n"
	for _, post := range posts {
		mdTop += fmt.Sprintf("<li><a href='/%s/%s/%s'>%s</a></li>\n", post.Header.Date.Year,
			post.Header.Date.Month, post.Header.Date.Date, post.Header.Title)
	}
	mdTop += "</ul>\n"

	fmt.Println(mdTop)

	content := md2HTML([]byte(mdTop), renderer)
	htmlString := concatTemplates(content)
	htmlFile := file.CreateFile("index.html")
	defer htmlFile.Close()

	file.MoveFile("index.html", "public/index.html")
	fmt.Fprintln(htmlFile, htmlString)
}

func reset() {
	file.RemoveDir("public")
	file.RemoveFile("theme/template/header.html")
}

func Run() {
	// check `source/` and `template/` existence.
	if ok := checkDirectoryExistence(); !ok {
		fmt.Println("exit.")
		return
	}

	reset()

	renderer := getRenderer()
	mds := file.GetMarkdownPaths("source")

	posts := generatePosts(renderer, mds)
	generateTopPage(renderer, posts)

	file.CopyDir("theme/css", "public/css")
}

