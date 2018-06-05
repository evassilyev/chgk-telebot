package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
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

func (qh *QuestionHandler) LoadPacket(limit int) (packet Packet, err error) {
	//"https://db.chgk.info/xml/random/types12/limit%d"
	//types12 - Типы вопросов: Что? Где? Когда? | Брейн-ринг
	data, err := qh.getXML(fmt.Sprintf("https://db.chgk.info/xml/random/types12/limit%d", limit))
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
