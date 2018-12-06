package build

import (
	"fmt"

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

func Run() {
	if ok := checkDirectoryExistence(); !ok {
		fmt.Println("exit.")
		return
	}

	md := file.LoadMarkdown("source/hello-world.md")
	if md == nil {
		fmt.Println("markdown load error")
		return
	}

	post := parser.ParseMarkdown(md)

	renderer := getRenderer()
	html := md2HTML(post.Body, renderer)

	file.CreateDirForPublish()
	index := file.CreateFile("index.html")
	defer index.Close()

	file.MoveFile("index.html", "public/index.html")
	fmt.Fprintln(index, html)
}

