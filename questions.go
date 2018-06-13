package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type QuestionHandler struct {
	database []Question
	client   *http.Client
}

type Question struct {
	Question     string `xml:"Question"`
	Answer       string `xml:"Answer"`
	PassCriteria string `xml:"PassCriteria"`
	Authors      string `xml:"Authors"`
	Sources      string `xml:"Sources"`
	Comments     string `xml:"Comments"`
	Notices      string `xml:"Notices"`
	Tournament   string `xml:"tournamentTitle"`
}

type Packet struct {
	Questions []Question `xml:"question"`
}

func NewQuestionHandler() *QuestionHandler {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}
	return &QuestionHandler{
		client: client,
	}
}

func (qh *QuestionHandler) LoadPacket(limit int, qt QuestionTypes) (packet Packet, err error) {
	//"https://db.chgk.info/xml/random/types12/limit%d"
	//types12 - Типы вопросов: Что? Где? Когда? | Брейн-ринг
	data, err := qh.getXML(fmt.Sprintf("https://db.chgk.info/xml/random/types%s/limit%d", qt.EncodeToUrlString(), limit))
	if err != nil {
		return Packet{}, err
	}
	err = xml.Unmarshal(data, &packet)
	return
}

func (qh *QuestionHandler) getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body error: %v", err)
	}

	return data, nil
}

type QuestionTypes struct {
	www  bool
	br   bool
	intt bool
	bes  bool
	myg  bool
	eru  bool
}

func (qt QuestionTypes)EncodeToUrlString()string  {
	result := ""
	if qt.www {
		result = result + "1"
	}
	if qt.br {
		result = result + "2"
	}
	if qt.intt {
		result = result + "3"
	}
	if qt.bes {
		result = result + "4"
	}
	if qt.myg {
		result = result + "5"
	}
	if qt.eru {
		result = result + "6"
	}
	return result
}

func (qt QuestionTypes)EncodeToUserString()string  {
	var rarr []string
	if qt.www {
		rarr = append(rarr, "Что? Где? Когда?")
	}
	if qt.br {
		rarr = append(rarr, "Брейн-ринг")
	}
	if qt.intt {
		rarr = append(rarr, "Интернет")
	}
	if qt.bes {
		rarr = append(rarr, "Бескрылка")
	}
	if qt.myg {
		rarr = append(rarr, "Своя игра")
	}
	if qt.eru {
		rarr = append(rarr, "Эрудитка")
	}
	result := strings.Join(rarr, ", ")
	return result
}
