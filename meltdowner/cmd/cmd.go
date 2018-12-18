package cmd

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wassan128/meltdowner/meltdowner/build"
	"github.com/wassan128/meltdowner/meltdowner/file"
)

type Opts struct {
	optBool bool
}
var o = &Opts{}

var RootCmd = &cobra.Command{
	Use: "melt",
	Short: "CLI tool for MeltDonwer(blog generator)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("melt(MeltDowner for CLI)")
		fmt.Println("Usage: melt COMMAND [OPTION]")
		fmt.Println("       See also `melt help` for more information.")
		fmt.Println("Author: wassan128")
	},
}

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Print the version number of melt",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("melt v0.1")
	},
}

var initCmd = &cobra.Command{
	Use: "init",
	Short: "Initialize MeltDowner directory",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var generateCmd = &cobra.Command{
	Use: "generate",
	Aliases: []string{"g"},
	Short: "Generate blog(htmls, static files)",
	Run: func(cmd *cobra.Command, args []string) {
		build.Run()
	},
}

var serverCmd = &cobra.Command{
	Use: "server",
	Aliases: []string{"s"},
	Short: "Serve public/ on localhost",
	Run: func(cmd *cobra.Command, args []string) {
		if o.optBool {
			fmt.Println("[*] found option generate before serve.")
			build.Run()
		}
		http.Handle("/", http.FileServer(http.Dir("public")))

		fmt.Println("[*] public/ is being served on http://localhost:5000")
		if err := http.ListenAndServe(":5000", nil); err != nil {
			fmt.Println(err)
			return
		}
	},
}

func getNowYMD() (int, int, int) {
	nowTime := time.Now()
	year := int(nowTime.Year())
	month := int(nowTime.Month())
	date := int(nowTime.Day())
	return year, month, date
}

var newCmd = &cobra.Command{
	Use: "new",
	Short: "Create new post",
	Run: func(cmd *cobra.Command, args []string) {
		title := strings.Replace(args[0], " ", "-", -1)
		fmt.Printf("[*] create new post: %s\n", title)

		year, month, date := getNowYMD()
		dateStr := fmt.Sprintf("%d%d%d", year, month, date)

		mdPaths := file.GetMarkdownPaths("source")
		id := 1
		for _, mdPath := range mdPaths {
			if strings.Index(mdPath, dateStr) != -1 {
				id++
			}
		}

		postPath := fmt.Sprintf("%s_%d_%s", dateStr, id, title)
		mdPath := fmt.Sprintf("%s.md", postPath)
		file.CreateDir(postPath)

		md := file.CreateFile(mdPath)
		defer md.Close()

		fmt.Fprintf(md, "title: %s\n", title)
		fmt.Fprintf(md, "date: %d-%d-%d\n", year, month, date)
		fmt.Fprintf(md, "---\n")

		file.MoveFile(postPath, filepath.Join("source", postPath))
		file.MoveFile(mdPath, filepath.Join("source", mdPath))
	},
}

func init() {
	cobra.OnInitialize()

	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(generateCmd)
	RootCmd.AddCommand(serverCmd)
	serverCmd.Flags().BoolVarP(&o.optBool, "generate", "g", false, "generate before serve")
	RootCmd.AddCommand(newCmd)
}

