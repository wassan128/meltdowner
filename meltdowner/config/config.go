package config

import (
	"github.com/BurntSushi/toml"
)

type BlogInfo struct {
	Title string
	Author string
}
type GitHubInfo struct {
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
