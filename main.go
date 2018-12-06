package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/russross/blackfriday.v2"
)

func main() {
	md, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	htmlFlags := blackfriday.CommonHTMLFlags
	htmlFlags |= blackfriday.FootnoteReturnLinks
	htmlFlags |= blackfriday.SmartypantsAngledQuotes
	htmlFlags |= blackfriday.SmartypantsQuotesNBSP
	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{Flags: htmlFlags, Title: "", CSS: ""})

	extFlags := blackfriday.CommonExtensions
	extFlags |= blackfriday.Footnotes
	extFlags |= blackfriday.HeadingIDs
	extFlags |= blackfriday.Titleblock
	extFlags |= blackfriday.DefinitionLists

	html := blackfriday.Run(md, blackfriday.WithExtensions(extFlags), blackfriday.WithRenderer(renderer))
	fmt.Println(string(html))
}

