package config

import (
	"github.com/BurntSushi/toml"
	"github.com/wassan128/meltdowner/meltdowner/util"
)

type BlogConfig struct {
	Title    string
	SubTitle string
	RootPath string
	IconURL  string
	Author   string
}

type GitHubConfig struct {
	Id    string
	Email string
	Repo  string
}

type Config struct {
	Blog   BlogConfig
	GitHub GitHubConfig
}

func LoadConfig() Config {
	var config Config

	_, err := toml.DecodeFile("bigbang.toml", &config)
	util.ExitIfError(err)

	return config
}

func (c *Config) Validate() []string {
	const (
		WARNING = iota
		ERROR
	)

	checks := []struct {
		bad  bool
		msg  string
		kind int
	}{
		{c.Blog.Title == "", "Title must specifiy.", ERROR},
		{c.Blog.SubTitle == "", "Subtitle does not specified.", WARNING},
		{c.Blog.RootPath == "", "RootPath must specify.", ERROR},
		{c.Blog.IconURL == "", "Icon url does not specified.", WARNING},
		{c.Blog.Author == "", "Author does not specified.", WARNING},
		{
			c.GitHub.Repo != "" &&
				(c.GitHub.Id == "" || c.GitHub.Email == ""),
			"GitHub repository specified but GitHub ID or Email does not specified.",
			WARNING,
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
