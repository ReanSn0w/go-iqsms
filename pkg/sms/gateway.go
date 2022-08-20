package sms

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL = "https://api.iqsms.ru/messages/v2/"
)

var (
	ErrWrongResponseParts = errors.New("wrong parts in response")
)

type Gateway struct {
	client *http.Client
	auth   string
}

// Инициализация нового шлюза для отправки сообщений
func NewSMSGateway(login string, password string) *Gateway {
	return &Gateway{
		auth: base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", login, password))),
		client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

// Передача сообщения
func (g *Gateway) Send(msg Message) (*Result, error) {
	body, err := g.request(fmt.Sprintf("send/?%s", msg.Encode()))
	if err != nil {
		return nil, err
	}

	return g.getMessageStatus(body)
}

// Проверка состояния отправленного сообщения (до 200 id в запросе)
func (g *Gateway) CheckMessages(ids ...string) ([]*Result, error) {
	values := url.Values{}
	for _, id := range ids {
		values.Add("id", id)
	}

	body, err := g.request(fmt.Sprintf("status/?%s", values.Encode()))
	if err != nil {
		return []*Result{}, err
	}

	return g.getMessagesStatus(body)
}

// Проверка очереди статусов отправленных сообщений
func (g *Gateway) CheckQuery(name string, limit int) ([]*Result, error) {
	values := url.Values{}
	values.Add("statusQueueName", name)
	values.Add("limit", fmt.Sprint(limit))

	body, err := g.request(fmt.Sprintf("statusQueue/?%s", values.Encode()))
	if err != nil {
		return []*Result{}, err
	}

	return g.getMessagesStatus(body)
}

// Метод для проверки баланса на аккаунте
func (g *Gateway) Balance() (float64, error) {
	body, err := g.request("balance/")
	if err != nil {
		return 0, err
	}

	parts := strings.Split(body, ";")
	if len(parts) != 2 {
		return 0, ErrWrongResponseParts
	}

	return strconv.ParseFloat(parts[1], 64)
}

// Метод для получения списка доступных подписей отправителя
func (g *Gateway) Senders() ([]string, error) {
	body, err := g.request("senders/")
	if err != nil {
		return []string{}, err
	}

	return strings.Split(body, "\n"), nil
}

func (g *Gateway) getMessagesStatus(body string) ([]*Result, error) {
	resultStrings := strings.Split(body, "\n")
	results := []*Result{}

	for _, item := range resultStrings {
		val, err := g.getMessageStatus(item)
		if err != nil {
			return results, err
		}

		results = append(results, val)
	}

	return results, nil
}

func (g *Gateway) getMessageStatus(value string) (*Result, error) {
	if !strings.Contains(value, ";") {
		return nil, errors.New(value)
	}

	values := strings.Split(value, ";")
	if len(values) != 2 {
		return nil, errors.New("incorect values count")
	}

	return &Result{ID: values[1], Status: status(values[0])}, nil
}

func (g *Gateway) request(url string) (string, error) {
	req, err := http.NewRequest("GET", baseURL+url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Basic "+g.auth)

	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(resp.Body)
	return buffer.String(), err
}
