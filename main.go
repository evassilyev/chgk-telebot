package main

import (
	"flag"
	"io/ioutil"
	"encoding/json"
)

const (
	version = "1.2"
	sourcesUrl = "https://github.com/evassilyev/chgk-telebot"
)

type Config struct {
	Token string `json:"token"`
	Chat int64 `json:"chat"`
}

func main() {
	var (
		configFile string
	)

	flag.StringVar(&configFile, "config", "config.json", "Configuration file path")
	flag.Parse()

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	telebot, err := NewTelebot(config.Token, config.Chat)
	if err != nil {
		panic(err)
	}

	game := NewGame(telebot)

	game.Play()
}
