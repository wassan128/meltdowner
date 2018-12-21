package config

import (
	"github.com/BurntSushi/toml"
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
	if err != nil {
		panic(err)
	}

	return config
}
