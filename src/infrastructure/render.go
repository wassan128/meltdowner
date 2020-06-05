package infrastructure

import "github.com/russross/blackfriday/v2"

func NewRenderer() *ChromaRenderer {
	htmlFlags := blackfriday.CommonHTMLFlags
	htmlFlags |= blackfriday.FootnoteReturnLinks
	htmlFlags |= blackfriday.SmartypantsAngledQuotes
	htmlFlags |= blackfriday.SmartypantsQuotesNBSP

	return &ChromaRenderer{
		html: blackfriday.NewHTMLRenderer(
			blackfriday.HTMLRendererParameters{
				Flags: htmlFlags,
			},
		),
		// TODO settable
		theme: "paraiso-dark",
	}
}

