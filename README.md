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
