package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

//Commands
const (
	help         = "/?"
	help2        = "/HELP"
	packet       = "/get_packet"
	packet_rus   = "/получить_пакет"
	start        = "/start"
	start_rus    = "/начать"
	next         = "/next"
	next_rus     = "/след"
	prev         = "/prev"
	prev_rus     = "/пред"
	question     = "/question"
	question_rus = "/вопрос"
	answer       = "/answer"
	answer_rus   = "/ответ"
	info         = "/info"
	packetB      = "/get_packet@CheGKBot"
	startB       = "/start@CheGKBot"
	nextB        = "/next@CheGKBot"
	prevB        = "/prev@CheGKBot"
	questionB    = "/question@CheGKBot"
	answerB      = "/answer@CheGKBot"
	infoB        = "/info@CheGKBot"
	timer        = "/set_timer"
	timerB       = "/set_timer@CheGKBot"
	setQtypes	= "/set_question_types"
	setQtypesB	= "/set_question_types@CheGKBot"
)

type Game struct {
	bot            *Telebot
	qh             *QuestionHandler

	questions      []Question
	qind           int

	packetLoaded   bool

	timer          *time.Timer
	alarmTimer     *time.Timer
	tout           time.Duration

	lastPacketSize int
	//Interactive state
	qtWaiting bool
	timerWaiting bool
}

func NewGame(bot *Telebot) *Game {
	return &Game{
		bot:            bot,
		qh:             NewQuestionHandler(),
		questions:      []Question{},
		qind:           0,
		packetLoaded:   false,
		timer:          nil,
		alarmTimer:     nil,
		tout:           time.Minute,
		lastPacketSize: 0,
		qtWaiting: false,
		timerWaiting:false,
	}
}

//TODO change panic to log output
func (g *Game) Play() {
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


func (g *Game) parseMessage(msg *tgbotapi.Message) {

	if g.handleAnswerWaiting(msg.Text) {
		return
	}

	words := strings.Split(msg.Text, " ")
	if len(words) == 0 {
		g.bot.ReplyToMessage(msg, "Command handling error")
		return
	}

	switch words[0] {

	case help:
		g.sendHelpMessage()
	case help2:
		g.sendHelpMessage()

	case packet:
		g.load(msg, words)
	case packet_rus:
		g.load(msg, words)
	case packetB:
		g.load(msg, words)

	case start:
		g.showQuestion(0)
	case start_rus:
		g.showQuestion(0)
	case startB:
		g.showQuestion(0)

	case next:
		g.next(msg, words)
	case next_rus:
		g.next(msg, words)
	case nextB:
		g.next(msg, words)

	case prev:
		g.prev(msg, words)
	case prev_rus:
		g.prev(msg, words)
	case prevB:
		g.prev(msg, words)

	case question:
		g.question(msg, words)
	case question_rus:
		g.question(msg, words)
	case questionB:
		g.question(msg, words)

	case answer:
		g.showAnswer()
	case answer_rus:
		g.showAnswer()
	case answerB:
		g.showAnswer()

	case info:
		g.showInfo()
	case infoB:
		g.showInfo()

	case timer:
		g.setTimer(msg, words)
	case timerB:
		g.setTimer(msg, words)

	case setQtypes:
		g.setQuestionTypes(msg)
	case setQtypesB:
		g.setQuestionTypes(msg)

	default:
		g.bot.ReplyToMessage(msg, "Не знаю такую команду!")
	}
}

func (g *Game)handleAnswerWaiting(msgText string) bool {
	if g.qtWaiting {
		fmt.Println(msgText)

		g.qtWaiting = false
		return true
	}

	if g.timerWaiting {

		g.timerWaiting = false
		return true
	}
	return false
}

func (g *Game) LoadPacket(packetSize int) {
	p, err := g.qh.LoadPacket(packetSize)
	if err != nil {
		g.bot.SendMessage(fmt.Sprintf("Не могу загрузить пакет вопросов! Ошибка: %#v", err))
		return
	}
	g.questions = p.Questions
	g.bot.SendMessage(fmt.Sprintf("Загружено вопросов: %d ", len(g.questions)))

	g.qind = 0

	g.packetLoaded = true
}

func (g *Game) sendHelpMessage() {
	helpMessage := fmt.Sprintf("Команды для работы с ботом:\n"+
		"%s, %s - Показать справку по командам\n"+
		"%s N, %s N - загрузить пакет из N вопросов\n"+
		"%s, %s - начать игру\n"+
		"%s, %s - следующий вопрос\n"+
		"%s, %s - предыдущий вопрос\n"+
		"%s N, %s N - перейти к вопросу под номером N\n"+
		"%s, %s - показать ответ\n"+
		"%s - показать информацию о вопросе (автор, источники и т.д.)\n"+
		"%s - установить таймер в минутах\n",
		help, help2,
		packet, packet_rus,
		start, start_rus,
		next, next_rus,
		prev, prev_rus,
		question, question_rus,
		answer, answer_rus,
		info,
		timer)

	g.bot.SendMessage(helpMessage)
}

func (g *Game)setQuestionTypes(msg *tgbotapi.Message)  {
	g.qtWaiting = true
	g.bot.ReplyToMessage(msg, "Ответьте на это сообщение")
	return
}

func (g *Game) setTimer(msg *tgbotapi.Message, words []string) {
	if len(words) != 2 {
		g.bot.ReplyToMessage(msg, "Укажите после команды (/set_timer) таймаут для таймера в минутах")
		return
	}
	t, err := strconv.ParseFloat(words[1], 64)
	if err != nil {
		g.bot.ReplyToMessage(msg, fmt.Sprintf("Не могу распознать количество минут (%s)!", words[1]))
		return
	}
	if t <= 0.25 {
		g.bot.ReplyToMessage(msg, fmt.Sprintf("Слишком маленькое значение таймера!"))
		return
	}

	g.tout = time.Duration(int64(float64(time.Minute.Nanoseconds()) * t))
	if g.timer == nil {
		g.timer = time.NewTimer(g.tout)
		g.alarmTimer = time.NewTimer(g.tout - (15 * time.Second))
		g.timer.Stop()
		g.alarmTimer.Stop()
	} else {
		g.timer.Reset(g.tout)
		g.alarmTimer.Reset(g.tout - (15 * time.Second))
	}
	g.bot.SendMessage(fmt.Sprintf("Таймер на %.1f мин. установлен\n", t))
}

func (g *Game) load(msg *tgbotapi.Message, words []string) {
	var (
		packetSize int
		err        error
	)

	if len(words) == 2 {
		packetSize, err = strconv.Atoi(words[1])
		if err != nil {
			g.bot.SendMessage(fmt.Sprintf("(%s) Не могу распознать число!", words[1]))
			return
		}
		g.lastPacketSize = packetSize
	} else if g.lastPacketSize != 0 {
		packetSize = g.lastPacketSize
	} else {
		g.bot.ReplyToMessage(msg, "Укажите после команды (/get_packet) , сколько вопросов нужно загрузить")
		return
	}

	g.LoadPacket(packetSize)
}

func (g *Game) next(msg *tgbotapi.Message, words []string) {
	if g.qind+1 >= len(g.questions) {
		g.bot.ReplyToMessage(msg, "Уже последний вопрос!")
		return
	}
	g.qind = g.qind + 1
	g.showQuestion(g.qind)
}

func (g *Game) prev(msg *tgbotapi.Message, words []string) {
	if g.qind-1 < 0 {
		g.bot.ReplyToMessage(msg, "Уже первый вопрос!")
		return
	}
	g.qind = g.qind - 1
	g.showQuestion(g.qind)
}

func (g *Game) question(msg *tgbotapi.Message, words []string) {
	if !g.packetLoaded {
		g.bot.SendMessage("Пакет не загружен! Загрузите пакет командой /get_packet")
		return
	}
	if len(words) != 2 {
		g.bot.ReplyToMessage(msg, "Укажите после команды (/question) номер вопроса")
		return
	}
	qn, err := strconv.Atoi(words[1])
	if err != nil {
		g.bot.ReplyToMessage(msg, fmt.Sprintf("Не могу распознать номер вопроса (%s)!", words[1]))
		return
	}
	if qn <= 0 || qn > len(g.questions) {
		g.bot.ReplyToMessage(msg, fmt.Sprintf("Неправильный номер вопроса (%d)! В пакете (%d) вопросов!", qn, len(g.questions)))
		return
	}

	g.qind = qn - 1
	g.showQuestion(g.qind)
}

func (g *Game) showQuestion(qi int) {
	if !g.packetLoaded {
		g.bot.SendMessage("Пакет не загружен! Загрузите пакет командой /get_packet")
		return
	}
	g.qind = qi
	questionMsg := fmt.Sprintf("Вопрос №%d\n"+
		"%s\n"+
		"\n",
		g.qind+1, g.questions[g.qind].Question)
	g.bot.SendMessage(questionMsg)

	if g.timer != nil {
		g.timer.Reset(g.tout)
		g.alarmTimer.Reset(g.tout)
		go func() {
			<-g.timer.C
			g.showAnswer()
		}()
		g.alarmTimer.Reset(g.tout - (15 * time.Second))
		go func() {
			<-g.alarmTimer.C
			g.bot.SendMessage("Осталось 15 секунд")
		}()
	}
}

func (g *Game) showAnswer() {
	if !g.packetLoaded {
		g.bot.SendMessage("Пакет не загружен! Загрузите пакет командой /get_packet")
		return
	}
	answerMsg := fmt.Sprintf("Ответ: %s\n"+
		"Зачет: %s\n"+
		"Комментарий: %s\n"+
		"Замечания: %s\n",
		g.questions[g.qind].Answer,
		g.questions[g.qind].PassCriteria,
		g.questions[g.qind].Comments,
		g.questions[g.qind].Notices)

	g.bot.SendMessage(answerMsg)

	if g.timer != nil {
		g.timer.Stop()
		g.alarmTimer.Stop()
	}
}

func (g *Game) showInfo() {
	if !g.packetLoaded {
		g.bot.SendMessage("Пакет не загружен! Загрузите пакет командой /get_packet")
		return
	}
	answerMsg := fmt.Sprintf(
		"Чемпионат: %s\n"+
			"Автор: %s\n"+
			"Источники: %s\n",
		g.questions[g.qind].Tournament,
		g.questions[g.qind].Authors,
		g.questions[g.qind].Sources)
	g.bot.SendMessage(answerMsg)
}
