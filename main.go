package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/russross/blackfriday.v2"
)

func setHTMLFlags() blackfriday.HTMLFlags {
	htmlFlags := blackfriday.CommonHTMLFlags
	htmlFlags |= blackfriday.FootnoteReturnLinks
	htmlFlags |= blackfriday.SmartypantsAngledQuotes
	htmlFlags |= blackfriday.SmartypantsQuotesNBSP

	return htmlFlags
}

func setExtensionFlags() blackfriday.Extensions {
	extFlags := blackfriday.CommonExtensions
	extFlags |= blackfriday.Footnotes
	extFlags |= blackfriday.HeadingIDs
	extFlags |= blackfriday.Titleblock
	extFlags |= blackfriday.DefinitionLists

	return extFlags
}

func createDirForPublish() {
	if err := os.Mkdir("public", 0777); err != nil {
		fmt.Println(err)
	}
}

func main() {
	md, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	htmlFlags := setHTMLFlags()
	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{Flags: htmlFlags, Title: "", CSS: ""})


	extFlags := setExtensionFlags()
	html := blackfriday.Run(md, blackfriday.WithExtensions(extFlags), blackfriday.WithRenderer(renderer))

	createDirForPublish()
	fmt.Println(string(html))
}

