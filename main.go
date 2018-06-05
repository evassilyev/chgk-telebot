package main

import (
	"flag"
	"io/ioutil"
	"encoding/json"
	"fmt"
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

	/*
		telebot, err := tgbotapi.NewBotAPI(config.Token)
		if err != nil {
			panic(err)
		}
		msg := tgbotapi.NewMessage(config.Chat, "Bot started")
		_, err = telebot.Send(msg)
		if err != nil {
			panic(err)
		}
	*/

	qh := NewQuestionHandler()
	questions, err := qh.LoadPacket(4)
	if err != nil {
		panic(err)
	}
	fmt.Println(questions)

	/*Program the telegram bot*/
}
