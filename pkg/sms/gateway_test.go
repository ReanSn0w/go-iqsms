package sms_test

import (
	"testing"
	"time"

	"github.com/ReanSn0w/go-iqsms/pkg/sms"
)

func Test_SendMessage(t *testing.T) {
	gateway, err := sms.NewSMSGatewayFromEnv()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	res, err := gateway.Send(sms.Message{
		Phone: "+79261234567",
		Text:  "Test Message",
	})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	time.Sleep(time.Second * 5)
	items, err := gateway.CheckMessages(res.ID)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, item := range items {
		if item == nil {
			t.Log("empty element in slice")
			t.Fail()
			continue
		}

		t.Logf("MessageID: %s, Status: %s", item.ID, item.Status)
	}
}
