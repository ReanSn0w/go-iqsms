package sms

import (
	"errors"
	"os"
)

const (
	loginEnv = "SMS_GATEWAY_LOGIN"
	passEnv  = "SMS_GATEWAY_PASSWORD"
)

var (
	ErrLoginMissing    = errors.New("sms gateway login missing")
	ErrPasswordMissing = errors.New("sms gateway password missing")
)

// Функция инициализирует шлюз из окружения операционной системы
func NewSMSGatewayFromEnv() (*Gateway, error) {
	login, ok := os.LookupEnv(loginEnv)
	if !ok || len(login) == 0 {
		return nil, ErrLoginMissing
	}

	password, ok := os.LookupEnv(passEnv)
	if !ok || len(password) == 0 {
		return nil, ErrPasswordMissing
	}

	// Данные попали в приложение и более незачем хранить их в окружении ОС
	os.Setenv(loginEnv, "")
	os.Setenv(passEnv, "")

	gateway := NewSMSGateway(login, password)
	return gateway, nil
}
