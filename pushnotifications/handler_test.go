package pushnotifications

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReceived(t *testing.T) {
	const (
		id          = "ghree"
		messageId   = "gkjuhkserir"
		accountId   = "khraeherre"
		messageText = "rejreijroire"
		from        = "5u25425425"
		to          = "49034814jn,f"
	)

	ch := make(chan ReceivedNotification, 1)

	s := httptest.NewServer(Received(func(message ReceivedNotification) {
		ch <- message
	}))
	defer s.Close()

	body := `<InboundMessage>
 <Id>` + id + `</Id>
 <MessageId>` + messageId + `</MessageId>
 <AccountId>` + accountId + `</AccountId>
 <MessageText>` + messageText + `</MessageText>
 <From>` + from + `</From>
 <To>` + to + `</To>
</InboundMessage>`
	resp, err := http.Post(s.URL, "application/xml", strings.NewReader(body))

	assert := assert.New(t)

	assert.Nil(err)
	assert.Equal(resp.StatusCode, 200)

	select {
	case message := <-ch:
		assert.Equal(message.Id, id)
		assert.Equal(message.MessageId, messageId)
		assert.Equal(message.AccountId, accountId)
		assert.Equal(message.MessageText, messageText)
		assert.Equal(message.From, from)
		assert.Equal(message.To, to)

	case <-time.After(time.Second):
		assert.Fail("timeout")
	}
}

func TestDelivered(t *testing.T) {
	const (
		id        = "ghree"
		messageId = "gkjuhkserir"
		accountId = "khraeherre"
	)

	var (
		occurredAt    = time.Date(2012, 1, 1, 12, 0, 5, 0, time.UTC)
		occurredAtStr = "2012-01-01T12:00:05"
	)

	ch := make(chan DeliveredNotification, 1)

	s := httptest.NewServer(Delivered(func(message DeliveredNotification) {
		ch <- message
	}))
	defer s.Close()

	body := `<MessageDelivered>
 <Id>` + id + `</Id>
 <MessageId>` + messageId + `</MessageId>
 <AccountId>` + accountId + `</AccountId>
 <OccurredAt>` + occurredAtStr + `</OccurredAt>
</MessageDelivered>`
	resp, err := http.Post(s.URL, "application/xml", strings.NewReader(body))

	assert := assert.New(t)

	assert.Nil(err)
	assert.Equal(resp.StatusCode, 200)

	select {
	case message := <-ch:
		assert.Equal(message.Id, id)
		assert.Equal(message.MessageId, messageId)
		assert.Equal(message.AccountId, accountId)
		assert.Equal(message.OccurredAt, occurredAt)

	case <-time.After(time.Second):
		assert.Fail("timeout")
	}
}

func TestFailed(t *testing.T) {
	const (
		id        = "ghree"
		messageId = "gkjuhkserir"
		accountId = "khraeherre"
	)

	var (
		occurredAt    = time.Date(2012, 1, 1, 12, 0, 5, 0, time.UTC)
		occurredAtStr = "2012-01-01T12:00:05"
	)

	ch := make(chan FailedNotification, 1)

	s := httptest.NewServer(Failed(func(message FailedNotification) {
		ch <- message
	}))
	defer s.Close()

	body := `<MessageFailed>
 <Id>` + id + `</Id>
 <MessageId>` + messageId + `</MessageId>
 <AccountId>` + accountId + `</AccountId>
 <OccurredAt>` + occurredAtStr + `</OccurredAt>
</MessageFailed>`
	resp, err := http.Post(s.URL, "application/xml", strings.NewReader(body))

	assert := assert.New(t)

	assert.Nil(err)
	assert.Equal(resp.StatusCode, 200)

	select {
	case message := <-ch:
		assert.Equal(message.Id, id)
		assert.Equal(message.MessageId, messageId)
		assert.Equal(message.AccountId, accountId)
		assert.Equal(message.OccurredAt, occurredAt)

	case <-time.After(time.Second):
		assert.Fail("timeout")
	}
}
