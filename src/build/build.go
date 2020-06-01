package build

import (
    "os"
    "fmt"
    "html"
    "path/filepath"
    "text/template"
    "strings"

    "github.com/wassan128/meltdowner/meltdowner/config"
    "github.com/wassan128/meltdowner/meltdowner/file"
    "github.com/wassan128/meltdowner/meltdowner/parser"
    "github.com/wassan128/meltdowner/meltdowner/util"
    "github.com/microcosm-cc/bluemonday"
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

func getRenderer() *ChromaRenderer {
    htmlFlags := setHTMLFlags()
    return &ChromaRenderer{
        html: blackfriday.NewHTMLRenderer(
                blackfriday.HTMLRendererParameters{
                    Flags: htmlFlags,
                },
            ),
        theme: "paraiso-dark",
    }
}

func setExtensionFlags() blackfriday.Extensions {
    extFlags := blackfriday.CommonExtensions
    extFlags |= blackfriday.Footnotes
    extFlags |= blackfriday.HeadingIDs
    extFlags |= blackfriday.Titleblock
    extFlags |= blackfriday.DefinitionLists

    return extFlags
}

func md2HTML(md []byte, renderer *ChromaRenderer) string {
    extFlags := setExtensionFlags()
    raw := blackfriday.Run(md,
        blackfriday.WithExtensions(extFlags),
        blackfriday.WithRenderer(renderer))

    sanitizer := bluemonday.UGCPolicy()
    sanitizer.AllowElements("iframe")
    sanitizer.AllowAttrs("class").Matching(bluemonday.Paragraph).OnElements("ul")
    sanitizer.AllowAttrs("style").Globally()
    html := sanitizer.SanitizeBytes(raw)

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

func createPostDir(createdAt parser.CreatedAt, id string) string {
    postPath := "public"

    paths := []string{createdAt.Year, createdAt.Month, createdAt.Date, id}
    for _, path := range paths {
        postPath = filepath.Join(postPath, path)
        if _, err := os.Stat(postPath); err != nil {
            file.CreateDir(postPath)
        }
    }

    return postPath
}

func createTagDir(tagName string) string {
    tagPath := filepath.Join("public/tags", tagName)
    if _, err := os.Stat(tagPath); err != nil {
        file.CreateDir(tagPath)
    }

    return tagPath
}

func concatRootPath(path string) string {
    return filepath.Join(Config.Blog.RootPath, path)
}

func generatePosts(renderer *ChromaRenderer, mds []string) (posts []parser.Post, tagMap map[string][]*parser.Post) {
    tagMap = map[string][]*parser.Post{}
    for _, mdPath := range mds {
        util.Info(fmt.Sprintf("Start: %s", mdPath))
        md := file.LoadFileContents(mdPath)
        if md == nil {
            fmt.Println("markdown load error")
            return nil, nil
        }

        id := strings.Split(mdPath, "_")[1]
        post := parser.ParseMarkdown(md)
        post.Header.Id = id
        posts = append(posts, *post)

        publicState := ""
        if post.Header.Public == false {
            publicState = "<span class='post-public'>URL限定公開記事</span>"
        }

        tags := ""
        if len(post.Header.Tags) > 0 {
            tags = "<p class='post-tags'>"
            for _, tag := range post.Header.Tags {
                tagPath := fmt.Sprintf("tags/%s", html.EscapeString(tag))
                link := concatRootPath(tagPath)
                tags += fmt.Sprintf("<a href='%s'>#%s</a>", link, html.EscapeString(tag))
                tagMap[tag] = append(tagMap[tag], post)
            }
            tags += "</p>"
        }

        title := []byte(fmt.Sprintf("# %s\n", post.Header.Title))
        content := md2HTML(append(title, post.Body...), renderer)

        date := fmt.Sprintf("<p class='post-date'>投稿日: %s/%s/%s</p>\n", post.Header.Date.Year, post.Header.Date.Month, post.Header.Date.Date)

        state := "<aside>" + tags + date + "</aside>"

        htmlString := concatTemplates("<article>" + publicState + content + state + "</article>")
        htmlFile := file.CreateFile("index.html")
        defer htmlFile.Close()

        postPath := createPostDir(post.Header.Date, id)

        file.MoveFile("index.html", filepath.Join(postPath, "index.html"))
        file.CopyDir(strings.Split(mdPath, ".")[0], postPath)
        fmt.Fprintln(htmlFile, htmlString)
        util.Info(fmt.Sprintf("Done: %s", postPath))
    }

    return posts, tagMap
}

func generateTopPage(renderer *ChromaRenderer, posts []parser.Post) {
    mdTop := "<ul class='top'>\n"
    for _, post := range posts {
        if post.Header.Public == false {
            util.Info(fmt.Sprintf("Found hidden flag: %s", post.Header.Title))
            continue
        }
        date := fmt.Sprintf("%s/%s/%s", post.Header.Date.Year, post.Header.Date.Month, post.Header.Date.Date)
        dateSpan := fmt.Sprintf("<span>%s</span>", date)
        link := concatRootPath(filepath.Join(date, post.Header.Id))
        mdTop += fmt.Sprintf("<li><a href='%s'>%s%s</a></li>\n", link, dateSpan, post.Header.Title)
    }
    mdTop += "</ul>\n"

    content := md2HTML([]byte(mdTop), renderer)
    htmlString := concatTemplates(content)
    htmlFile := file.CreateFile("index.html")
    defer htmlFile.Close()
    file.MoveFile("index.html", "public/index.html")
    fmt.Fprintln(htmlFile, htmlString)
}

func generateTagTopPage(renderer *ChromaRenderer, tagMap map[string][]*parser.Post) {
    for tag, posts := range tagMap {
        tagPath := createTagDir(tag)

        mdTagTop := fmt.Sprintf("<h1 style='color:#aaa;font-weight:normal;text-align:center'>#%s</h1>", tag)
        mdTagTop += "<ul class='top'>\n"
        for _, post := range posts {
            if post.Header.Public == false {
                util.Info(fmt.Sprintf("Found hidden flag: %s", post.Header.Title))
                continue
            }
            date := fmt.Sprintf("%s/%s/%s", post.Header.Date.Year, post.Header.Date.Month, post.Header.Date.Date)
            dateSpan := fmt.Sprintf("<span>%s</span>", date)
            link := concatRootPath(filepath.Join(date, post.Header.Id))
            mdTagTop += fmt.Sprintf("<li><a href='%s'>%s%s</a></li>\n", link, dateSpan, post.Header.Title)
        }
        mdTagTop += "</ul>\n"

        content := md2HTML([]byte(mdTagTop), renderer)
        htmlString := concatTemplates(content)
        htmlFile := file.CreateFile("index.html")
        defer htmlFile.Close()

        file.MoveFile("index.html", filepath.Join(tagPath, "index.html"))
        fmt.Fprintln(htmlFile, htmlString)
    }
}

func reset() {
    dirs, _ := filepath.Glob("public/[0-9][0-9][0-9][0-9]")
    for _, dir := range dirs {
        file.RemoveDir(dir)
    }
    file.RemoveDir("public/css")
    file.RemoveDir("public/tags")
    file.CreateDir("public/tags")
    file.RemoveFile("public/index.html")
    file.RemoveFile("theme/template/header.html")
}

func Run() {
    // check `source/` and `theme/` existence.
    if ok := checkDirectoryExistence(); !ok {
        fmt.Println("Quit build.")
        return
    }

    reset()

    renderer := getRenderer()
    mds := file.GetMarkdownPaths("source")

    posts, tagMap := generatePosts(renderer, mds)
    generateTopPage(renderer, posts)
    generateTagTopPage(renderer, tagMap)

    file.CopyDir("theme/css", "public/css")
}

