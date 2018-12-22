package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"github.com/spf13/cobra"
	"github.com/wassan128/meltdowner/meltdowner/build"
	"github.com/wassan128/meltdowner/meltdowner/config"
	"github.com/wassan128/meltdowner/meltdowner/file"
)

var Config config.Config = config.GetConfig()

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
		file.CreateDir("source")
		file.CreateDir("public")

		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}

		repo, err := git.PlainInit(filepath.Join(wd, "public"), false)
		if err != nil {
			fmt.Println(err)
		}

		worktree, err := repo.Worktree()
		if err != nil {
			fmt.Println(err)
		}

		worktree.Commit("initial commit", &git.CommitOptions{
			Author: &object.Signature{
				Name: Config.GitHub.Id,
				Email: Config.GitHub.Email,
				When: time.Now(),
			},
		})

		branch := plumbing.ReferenceName("refs/heads/gh-pages")
		worktree.Checkout(&git.CheckoutOptions{
			Create: true,
			Force: false,
			Branch: branch,
		})
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

func getNowTime() (int, int, int, int, int, int) {
	nowTime := time.Now()
	year := int(nowTime.Year())
	month := int(nowTime.Month())
	date := int(nowTime.Day())
	hour := int(nowTime.Hour())
	minute := int(nowTime.Minute())
	second := int(nowTime.Second())
	return year, month, date, hour, minute, second
}

var newCmd = &cobra.Command{
	Use: "new",
	Short: "Create new post",
	Run: func(cmd *cobra.Command, args []string) {
		if !file.IsExistPath("source") {
			fmt.Println("[Error] source/ not found.")
			return
		}
		title := strings.Replace(args[0], " ", "-", -1)
		fmt.Printf("[*] create new post: %s\n", title)

		year, month, date, _, _, _ := getNowTime()
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

var deployCmd = &cobra.Command{
	Use: "deploy",
	Short: "Deploy blog",
	Run: func(cmd *cobra.Command, args []string) {
		if o.optBool {
			fmt.Println("[*] found option generate before deploy.")
			build.Run()
		} else {
			if !file.IsExistPath("public") {
				fmt.Println("[Error] public/ not found.")
				return
			}
		}

		repo, err := git.PlainOpen("public")
		if err != nil {
			fmt.Println(err)
		}

		worktree, err := repo.Worktree()
		if err != nil {
			fmt.Println(err)
		}

		worktree.Add(".")

		y, m, d, h, mi, s := getNowTime()
		dateStr := fmt.Sprintf("%d/%d/%d %d:%d:%d", y, m, d, h, mi, s)
		cmsg := fmt.Sprintf("[update] %s", dateStr)

		worktree.Commit(cmsg, &git.CommitOptions{
			Author: &object.Signature{
				Name: Config.GitHub.Id,
				Email: Config.GitHub.Email,
				When: time.Now(),
			},
			All: true,
		})
	},
}

func init() {
	cobra.OnInitialize()

	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(generateCmd)
	RootCmd.AddCommand(serverCmd)
	serverCmd.Flags().BoolVarP(&o.optBool, "generate", "g", false, "generate before serve")
	RootCmd.AddCommand(newCmd)
	RootCmd.AddCommand(initCmd)
	RootCmd.AddCommand(deployCmd)
	deployCmd.Flags().BoolVarP(&o.optBool, "generate", "g", false, "generate before deploy")
}

