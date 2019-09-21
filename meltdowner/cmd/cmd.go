package cmd

import (
    "fmt"
    "net/http"
    "os"
    "os/exec"
    "strings"
    "time"
    "path/filepath"

    "gopkg.in/src-d/go-git.v4"
    gitConfig "gopkg.in/src-d/go-git.v4/config"
    "gopkg.in/src-d/go-git.v4/plumbing"
    "gopkg.in/src-d/go-git.v4/plumbing/object"
    "github.com/spf13/cobra"
    "github.com/wassan128/meltdowner/meltdowner/build"
    "github.com/wassan128/meltdowner/meltdowner/config"
    "github.com/wassan128/meltdowner/meltdowner/file"
    "github.com/wassan128/meltdowner/meltdowner/util"
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
        if o.optBool {
            util.Info("-r option found: reset before init.")
            file.RemoveDir("source")
            file.RemoveDir("public")
            util.Info("Deleted source/ and public/ because reset option found.")
        }
        file.CreateDir("source")
        file.CreateDir("public")
        file.CreateDir("public/tags")
        util.Info("Created source/ and public/")

        wd, err := os.Getwd()
        util.ExitIfError(err)

        repo, err := git.PlainInit(filepath.Join(wd, "public"), false)
        util.ExitIfError(err)

        worktree, err := repo.Worktree()
        util.ExitIfError(err)

        _, err = repo.CreateRemote(&gitConfig.RemoteConfig{
            Name: "origin",
            URLs: []string{Config.GitHub.Repo},
        })
        util.ExitIfError(err)
        if Config.GitHub.Repo != "" {
            util.Info(fmt.Sprintf("Created remote(origin->%s)", Config.GitHub.Repo))
        } else {
            util.Info("Remote repository does not registered")
        }

        branch := plumbing.ReferenceName("refs/heads/gh-pages")
        err = worktree.Pull(&git.PullOptions{
            RemoteName: "origin",
            ReferenceName: branch,
        })
        if err != nil {
            _, err = worktree.Commit("initial commit", &git.CommitOptions{
                Author: &object.Signature{
                    Name: Config.GitHub.Id,
                    Email: Config.GitHub.Email,
                    When: time.Now(),
                },
            })
            util.ExitIfError(err)
            util.Info("Created initial commit")
        } else {
            util.Info("Pulled from remote")
        }

        err = worktree.Checkout(&git.CheckoutOptions{
            Create: true,
            Branch: branch,
        })
        util.ExitIfError(err)
        util.Info(fmt.Sprintf("Checked out to %s", branch))

        util.Info("Done")
    },
}

var generateCmd = &cobra.Command{
    Use: "generate",
    Aliases: []string{"g"},
    Short: "Generate blog(htmls, static files)",
    Run: func(cmd *cobra.Command, args []string) {
        build.Run()

        util.Info("Done")
    },
}

var serverCmd = &cobra.Command{
    Use: "server",
    Aliases: []string{"s"},
    Short: "Serve public/ on localhost",
    Run: func(cmd *cobra.Command, args []string) {
        if o.optBool {
            util.Info("-g option found: generate before serve.")
            build.Run()
        }
        http.Handle(Config.Blog.RootPath,
            http.StripPrefix(Config.Blog.RootPath,
                http.FileServer(http.Dir("public"))))

        util.Info(fmt.Sprintf("public/ is being served on http://localhost:5000%s", Config.Blog.RootPath))
        err := http.ListenAndServe(":5000", nil)
        util.ExitIfError(err)
    },
}

func getNowTime() (year, month, date, hour, minute, second int) {
    nowTime := time.Now()
    year = int(nowTime.Year())
    month = int(nowTime.Month())
    date = int(nowTime.Day())
    hour = int(nowTime.Hour())
    minute = int(nowTime.Minute())
    second = int(nowTime.Second())
    return
}

var newCmd = &cobra.Command{
    Use: "new",
    Short: "Create new post",
    Run: func(cmd *cobra.Command, args []string) {
        title := args[0]
        titleForFileName := strings.Replace(title, " ", "-", -1)
        util.Info(fmt.Sprintf("Create new post: %s", title))

        year, month, date, _, _, _ := getNowTime()
        dateStr := fmt.Sprintf("%d%02d%02d", year, month, date)

        mdPaths := file.GetMarkdownPaths("source")
        id := 1
        for _, mdPath := range mdPaths {
            if strings.Index(mdPath, dateStr) != -1 {
                id++
            }
        }

        postPath := fmt.Sprintf("%s_%d_%s", dateStr, id, titleForFileName)
        mdPath := fmt.Sprintf("%s.md", postPath)
        file.CreateDir(postPath)

        md := file.CreateFile(mdPath)
        defer md.Close()

        fmt.Fprintf(md, "title: %s\n", title)
        fmt.Fprintf(md, "date: %d-%02d-%02d\n", year, month, date)
        fmt.Fprintf(md, "tags:\n")
        fmt.Fprintf(md, "---\n")

        file.MoveFile(postPath, filepath.Join("source", postPath))
        file.MoveFile(mdPath, filepath.Join("source", mdPath))

        util.Info("Done")
    },
}

var updateCmd = &cobra.Command{
    Use: "update",
    Short: "Update(pull & build) MeltDowner tool",
    Run: func(cmd *cobra.Command, args []string) {
        repo, err := git.PlainOpen(".")
        util.ExitIfError(err)

        worktree, err := repo.Worktree()
        util.ExitIfError(err)

        err = worktree.Pull(&git.PullOptions{
            RemoteName: "origin",
        })
        util.WarningIfError(err)

        err = exec.Command("go", "build", "meltdowner/main.go").Run()
        util.ExitIfError(err)

        err = exec.Command("mv", "-f", "main", "melt").Run()
        util.ExitIfError(err)
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
    initCmd.Flags().BoolVarP(&o.optBool, "reset", "r", false, "to reset to delete public/ and source/")
    RootCmd.AddCommand(updateCmd)
}

