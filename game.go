package main

import (
	"strings"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Game struct {
	bot *Telebot
	qh *QuestionHandler
}

func NewGame(bot *Telebot) *Game {
	return &Game{
		bot:bot,
		qh:NewQuestionHandler(),
	}
}

//TODO change panic to log output
func (g *Game)Play() {
	err := g.bot.SendMessage("Bot started")
	if err != nil {
		panic(err)
	}

	for update := range g.bot.updates {
		if update.Message == nil {
			continue
		}
		go g.parseMessage(update.Message)
	}
}

func (g *Game)parseMessage(msg *tgbotapi.Message)  {
	words := strings.Split(msg.Text, " ")
	if len(words) == 0 {
		g.bot.ReplyToMessage(msg, "Command handling error")
	}
}
/*
qh := NewQuestionHandler()
packet, err := qh.LoadPacket(4)
if err != nil {
panic(err)
}
fmt.Println(len(packet.Questions))
for _, q := range packet.Questions {
fmt.Println(q)
}
*/

/*Program the telegram bot*/
