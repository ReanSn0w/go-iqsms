# go-iqsms
Библиотека для отправки сообщений через сервис СМС Дисконт

## Быстрый старт

```go
gateway := NewSMSGateway("my_login", "my_password")

gateway.Send(sms.Message{
	Phone: "+79261234567",
	Text:  "Test Message",
})
```

[Документация](https://pkg.go.dev/github.com/ReanSn0w/go-iqsms@v1.0.0/pkg/sms)

Лицензия MIT
