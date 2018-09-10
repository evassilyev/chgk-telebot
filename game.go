package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"regexp"

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
	timer        = "/set_timer"
	setQtypes    = "/set_question_type"
	showSettings = "/show_settings"
	about        = "/about"
)

type Game struct {
	bot *Telebot
	qh  *QuestionHandler

	questions []Question
	qind      int

	packetLoaded   bool
	lastPacketSize int

	timer      *time.Timer
	alarmTimer *time.Timer
	tout       time.Duration

	//Interactive state
	qtWaiting    bool
	timerWaiting bool

	qtypes QuestionTypes
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
		qtWaiting:      false,
		timerWaiting:   false,

		qtypes: QuestionTypes{
			www:  true,
			br:   true,
			intt: false,
			bes:  false,
			myg:  false,
			eru:  false,
		},
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
	var command string

	//For removing the bot name in command
	atpos := strings.Index(words[0], "@")
	if atpos == -1 {
		command = words[0]
	} else {
		command = words[0][:atpos]
	}

	switch command {

	case help:
		g.sendHelpMessage()
	case help2:
		g.sendHelpMessage()

	case packet:
		g.load(msg, words)
	case packet_rus:
		g.load(msg, words)

	case start:
		g.showQuestion(0)
	case start_rus:
		g.showQuestion(0)

	case next:
		g.next(msg, words)
	case next_rus:
		g.next(msg, words)

	case prev:
		g.prev(msg, words)
	case prev_rus:
		g.prev(msg, words)

	case question:
		g.question(msg, words)
	case question_rus:
		g.question(msg, words)

	case answer:
		g.showAnswer()
	case answer_rus:
		g.showAnswer()

	case info:
		g.showInfo()

	case timer:
		g.setTimer(msg)

	case setQtypes:
		g.setQuestionTypes(msg)

	case showSettings:
		g.showSettings()
	case about:
		g.showAbout()

	default:
		g.bot.ReplyToMessage(msg, "Не знаю такую команду!")
	}
}

func (g *Game) handleAnswerWaiting(msgText string) bool {
	//Question types
	if g.qtWaiting {
		r, _ := regexp.Compile("^(/)?[1-6]+$")
		if !r.MatchString(msgText) {
			g.bot.SendMessage("Не могу распознать типы вопросов!\n" + msgText)
			return true
		}

		var qt QuestionTypes

		if strings.Contains(msgText, "1") {
			qt.www = true
		}
		if strings.Contains(msgText, "2") {
			qt.br = true
		}
		if strings.Contains(msgText, "3") {
			qt.intt = true
		}
		if strings.Contains(msgText, "4") {
			qt.bes = true
		}
		if strings.Contains(msgText, "5") {
			qt.myg = true
		}
		if strings.Contains(msgText, "6") {
			qt.eru = true
		}

		g.qtypes = qt

		g.qtWaiting = false

		g.bot.SendMessage(fmt.Sprintf("Установлены следующие типы загружаемых вопросов: %s\n", g.qtypes.EncodeToUserString()))

		return true
	}

	if g.timerWaiting {
		var minutes string
		fs := strings.Index(msgText, "/")
		if fs == -1 {
			minutes = msgText
		} else {
			minutes = msgText[1:]
		}

		t, err := strconv.ParseFloat(minutes, 64)
		if err != nil {
			g.bot.SendMessage(fmt.Sprintf("Не могу распознать количество минут (%s)!", minutes))
			return true
		}
		if t <= 0.25 {
			g.bot.SendMessage("Слишком маленькое значение таймера!")
			return true
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

		g.timerWaiting = false
		return true
	}
	return false
}

func (g *Game) LoadPacket(packetSize int) {
	p, err := g.qh.LoadPacket(packetSize, g.qtypes)
	if err != nil {
		g.bot.SendMessage(fmt.Sprintf("Не могу загрузить пакет вопросов! Ошибка: %v", err))
		return
	}
	g.questions = p.Questions
	g.bot.SendMessage(fmt.Sprintf("Загружено вопросов: %d\n", len(g.questions)))
	g.bot.SendMessage(fmt.Sprintf("Типы загруженных вопросов: %s\n", g.qtypes.EncodeToUserString()))

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
		"%s - установить таймер в минутах\n"+
		"%s - установить типы вопросов для загрузки\n"+
		"%s - показать настройки\n"+
		"%s - информация о боте\n",
		help, help2,
		packet, packet_rus,
		start, start_rus,
		next, next_rus,
		prev, prev_rus,
		question, question_rus,
		answer, answer_rus,
		info,
		timer,
		setQtypes,
		showSettings,
		about)

	g.bot.SendMessage(helpMessage)
}

func (g *Game) setQuestionTypes(msg *tgbotapi.Message) {
	message := "Отправьте сообщение, содержащее в себе цифры, где\n" +
		"1 - Что? Где? Когда?\n" +
		"2 - Брейн-ринг\n" +
		"3 - Интернет\n" +
		"4 - Бескрылка\n" +
		"5 - Своя игра\n" +
		"6 - Эрудитка\n" +
		"Сообщение отправьте с помощью \"Ответить\" на текущее, либо начав его с символа \"/\".\n"
	g.bot.ReplyToMessage(msg, message)
	g.qtWaiting = true
	return
}

func (g *Game) setTimer(msg *tgbotapi.Message) {
	message := "Отправьте количество минут для таймера в следующем сообщении (Допускаются дробные величины).\n" +
		"Сообщение отправьте с помощью \"Ответить\" на текущее, либо начав его с символа \"/\".\n"
	g.bot.ReplyToMessage(msg, message)
	g.timerWaiting = true

	if g.timer != nil {
		g.timer.Stop()
		g.alarmTimer.Stop()
	}
	return
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
		g.qind+1, g.questions[g.qind])
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

func (g *Game) showSettings() {
	message := fmt.Sprintf("Текущие настройки:\n"+
		"Вопросов в пакете загружается: %d\n"+
		"Таймер устанавливается на: %s\n"+
		"Загружаются вопросы типов: %s\n",
		g.lastPacketSize,
		g.tout.String(),
		g.qtypes.EncodeToUserString())
	g.bot.SendMessage(message)
}

func (g *Game) showAbout() {
	message := fmt.Sprintf("Версия: %s\n"+
		"Исходник: %s\n", version, sourcesUrl)
	g.bot.SendMessage(message)
}
