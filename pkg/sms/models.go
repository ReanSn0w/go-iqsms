package sms

import (
	"net/url"
	"time"
)

type status string

const (
	// Статусы сообщений
	StatusAccepted      status = "accepted"       // Сообщение принято в обработку
	StatusQueued        status = "queued"         // Сообщение находится в очереди
	StatusDelivered     status = "delivered"      // Сообщение доставлено
	StatusDeliveryError status = "delivery error" // Ошибка доставки SMS
	StatusSMSCSubmit    status = "smsc submit"    // Сообщение доставлено в SMSC
	StatusSMSCReject    status = "smsc reject"    // Сообщение отвергнуто SMSC
	StatusIncorrectID   status = "incorrect id"   // Неверный идентификатор сообщения
)

// Структура описывает сообщение, которое сервис может отправить
//
// Обязательные поля: Phone, Text
type Message struct {
	Phone    string    // Телефон в формате +71234567890
	Text     string    // Текст сообщения
	WapURL   string    // Wap-push ссылка (прим: wap.yousite.ru)
	Sender   string    // Подпись отправителя
	Flash    bool      // Flash SMS – сообщение, которое сразу отображается на экране и не сохраняется в памяти телефона
	Schedule time.Time // Время для отложенной отправки
	Queue    string    // Название очереди статусов отправленных сообщений
}

func (m *Message) Encode() string {
	values := url.Values{}
	values.Add("phone", m.Phone)
	values.Add("text", m.Text)

	if m.WapURL != "" {
		values.Add("wapurl", m.WapURL)
	}

	if m.Sender != "" {
		values.Add("sender", m.Sender)
	}

	if m.Flash {
		values.Add("flash", "1")
	}

	unix := m.Schedule.Unix()
	if unix > time.Now().Unix() {
		values.Add("scheduleTime", m.Schedule.Format("2006-01-02T15:04:05-07:00"))
	}

	if m.Queue != "" {
		values.Add("statusQueueName", m.Queue)
	}

	return values.Encode()
}

// Структура для описания результата
type Result struct {
	ID     string
	Status status
}
