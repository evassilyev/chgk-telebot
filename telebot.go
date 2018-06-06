package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Telebot struct {
	token string
	chat  int64
	tbot  *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}

func NewTelebot(token string, chat int64) (*Telebot, error) {

	telebot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := telebot.GetUpdatesChan(u)

	return &Telebot{
		token: token,
		chat:  chat,
		tbot:  telebot,
		updates:updates,
	}, err
}

func (bot *Telebot)SendMessage(message string) error  {
	msg := tgbotapi.NewMessage(bot.chat, message)
	_, err := bot.tbot.Send(msg)
	return err
}

func (bot *Telebot)ReplyToMessage(msg *tgbotapi.Message, text string) error  {
	rm := tgbotapi.NewMessage(msg.Chat.ID, text)
	rm.ReplyToMessageID = msg.MessageID
	_, err := bot.tbot.Send(rm)
	return err
}

