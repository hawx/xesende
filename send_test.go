package xesende

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleAccountClient_Send() {
	accountClient := New("user@example.com", "pass").Account("EX00000")

	accountClient.Send([]Message{
		{To: "00000000", Body: "Hello"},
	})
}

func TestSendSingleMessage(t *testing.T) {
	const (
		batchID          = "batchID"
		messageID        = "messageID"
		messageURI       = "messageURI"
		accountReference = "EXWHATEVS"
		to               = "358973"
		from             = "mehehherrr"
		body             = "HWEYERW"
		messageType      = Voice
		lang             = ".asmd.amd,ma.s,dma"
		validity         = 234
		characterSet     = "alkjhsdajklhsd"
		retries          = 10292
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messageheaders batchid="`+batchID+`" xmlns="http://api.esendex.com/ns/">
  <messageheader uri="`+messageURI+`" id="`+messageID+`" />
</messageheaders>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	account := client.Account(accountReference)

	result, err := account.Send([]Message{
		{
			To:           to,
			Body:         body,
			MessageType:  messageType,
			Lang:         lang,
			Validity:     validity,
			Retries:      retries,
			CharacterSet: characterSet,
		},
	})

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("POST", h.Request.Method)
	assert.Equal("/v1.0/messagedispatcher", h.Request.URL.String())

	var expectedBodyStr = fmt.Sprintf("<messages>"+
		"<accountreference>%s</accountreference>"+
		"<message>"+
		"<to>%s</to>"+
		"<type>%s</type>"+
		"<lang>%s</lang>"+
		"<validity>%d</validity>"+
		"<characterset>%s</characterset>"+
		"<retries>%d</retries>"+
		"<body>%s</body>"+
		"</message>"+
		"</messages>",
		accountReference, to, messageType, lang, validity, characterSet, retries, body)
	assert.Equal(expectedBodyStr, h.RequestBody)

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	assert.Equal(batchID, result.BatchID)

	assert.Equal(1, len(result.Messages))
	assert.Equal(messageID, result.Messages[0].ID)
	assert.Equal(messageURI, result.Messages[0].URI)
}

func TestSendAtSingleMessage(t *testing.T) {
	const (
		batchID          = "batchID"
		messageID        = "messageID"
		messageURI       = "messageURI"
		accountReference = "EXWHATEVS"
		to               = "cvfgdfg"
		from             = "2342"
		body             = "4334ffdgfh"
		messageType      = SMS
		lang             = ".asmd.amd,ma.s,dma"
		validity         = -3
		characterSet     = "alkjSDhsdajklhsd"
		retries          = 2
	)

	var (
		sendAt    = time.Date(2015, 11, 11, 11, 11, 11, 111111111, time.UTC)
		sendAtStr = "2015-11-11T11:11:11.111111111Z"
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messageheaders batchid="`+batchID+`" xmlns="http://api.esendex.com/ns/">
  <messageheader uri="`+messageURI+`" id="`+messageID+`" />
</messageheaders>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	account := client.Account(accountReference)

	result, err := account.SendAt(sendAt, []Message{
		{
			To:           to,
			Body:         body,
			MessageType:  messageType,
			Lang:         lang,
			Validity:     validity,
			Retries:      retries,
			CharacterSet: characterSet,
		},
	})

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("POST", h.Request.Method)
	assert.Equal("/v1.0/messagedispatcher", h.Request.URL.String())

	var expectedBodyStr = fmt.Sprintf("<messages>"+
		"<accountreference>%s</accountreference>"+
		"<sendat>%s</sendat>"+
		"<message>"+
		"<to>%s</to>"+
		"<type>%s</type>"+
		"<lang>%s</lang>"+
		"<validity>%d</validity>"+
		"<characterset>%s</characterset>"+
		"<retries>%d</retries>"+
		"<body>%s</body>"+
		"</message>"+
		"</messages>",
		accountReference, sendAtStr, to, messageType, lang, validity, characterSet, retries, body)
	assert.Equal(expectedBodyStr, h.RequestBody)

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	assert.Equal(batchID, result.BatchID)

	assert.Equal(1, len(result.Messages))
	assert.Equal(messageID, result.Messages[0].ID)
	assert.Equal(messageURI, result.Messages[0].URI)
}
