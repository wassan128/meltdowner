package config

import (
    "github.com/BurntSushi/toml"
    "github.com/wassan128/meltdowner/meltdowner/util"
)

type BlogInfo struct {
    Title string
    SubTitle string
    RootPath string
    IconURL string
    Author string
}
type GitHubInfo struct {
    Id string
    Email string
    Repo string
}
type Config struct {
    Blog BlogInfo
    GitHub GitHubInfo
}

func GetConfig() Config {
    var config Config

    _, err := toml.DecodeFile("bigbang.toml", &config)
    util.ExitIfError(err)

    return config
}

func (c Config) Validate() []string {
    checks := []struct{
        bad bool
        msg string
    } {
        { c.Blog.Title == "", "Title must specifiy." },
        { c.Blog.SubTitle == "", "Subtitle does not specified." },
        { c.Blog.RootPath == "", "RootPath must specify." },
        { c.Blog.IconURL== "", "Icon url does not specified." },
        { c.Blog.Author== "", "Author does not specified." },
        {
            c.GitHub.Repo != "" &&
            (c.GitHub.Id == "" || c.GitHub.Email == ""),
            "GitHub repository specified but GitHub ID or Email does not specified.",
        },
    }

    var msgs []string
    for _, check := range checks {
        if check.bad {
            msgs = append(msgs, check.msg)
        }
    }

    return msgs
}
