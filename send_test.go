package xesende

import (
	"net/http/httptest"
	"net/url"
	"testing"

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
		batchID    = "batchID"
		messageID  = "messageID"
		messageURI = "messageURI"
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messageheaders batchid="`+batchID+`" xmlns="http://api.esendex.com/ns/">
  <messageheader uri="`+messageURI+`" id="`+messageID+`" />
</messageheaders>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	account := client.Account("EXWHATEVS")

	result, err := account.Send([]Message{
		{To: "358973", Body: "HWEYERW"},
	})

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("POST", h.Request.Method)
	assert.Equal("/v1.0/messagedispatcher", h.Request.URL.String())

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	assert.Equal(batchID, result.BatchID)

	assert.Equal(1, len(result.Messages))
	assert.Equal(messageID, result.Messages[0].ID)
	assert.Equal(messageURI, result.Messages[0].URI)
}
