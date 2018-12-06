package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"log"

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

func createDirForPublish() {
	if err := os.Mkdir("public", 0777); err != nil {
		fmt.Println(err)
	}
}
func createFile(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return file
}
func moveFile(dstPath string, srcPath string) {
	if err := os.Rename(dstPath, srcPath); err != nil {
		fmt.Println(err)
	}
}

func main() {
	md, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	renderer := getRenderer()
	html := md2HTML(md, renderer)

	createDirForPublish()
	index := createFile("index.html")
	defer index.Close()

	moveFile("index.html", "public/index.html")

	fmt.Fprintln(index, html)
}

