package xesende

import (
	"fmt"
	"log"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleClient_Sent() {
	client := New("user@example.com", "pass")

	response, err := client.Sent()
	if err != nil {
		log.Fatal(err)
	}

	for _, message := range response.Messages {
		fmt.Printf("%v %s: %s\n", message.SubmittedAt, message.To, message.Summary)
	}
}

func ExampleClient_Received() {
	client := New("user@example.com", "pass")

	now := time.Now()

	response, err := client.Received(Between(now.AddDate(0, -6, 0), now))
	if err != nil {
		log.Fatal(err)
	}

	for _, message := range response.Messages {
		fmt.Printf("%v %s: %s\n", message.ReceivedAt, message.From, message.Summary)
	}
}

func ExampleClient_Body() {
	client := New("user@example.com", "pass")

	message, err := client.Message("messageId")
	if err != nil {
		log.Fatal(err)
	}

	body, err := client.Body(message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s: %s", message.ID, body.Text)
}

func TestMessagesSent(t *testing.T) {
	const (
		startIndex  = 0
		count       = 15
		totalCount  = 200
		id          = "messageheaderid"
		uri         = "http://somemessageheader"
		reference   = "EXETRTRE"
		status      = "STATUS"
		messageType = Voice
		to          = "4538224364236"
		from        = "428377843"
		summary     = "SUM"
		bodyURI     = "http://rrehekr"
		direction   = "OUT"
		parts       = 1
		username    = "user"
	)

	var (
		lastStatusAt    = time.Date(2012, 1, 1, 12, 0, 5, 0, time.UTC)
		lastStatusAtStr = "2012-01-01T12:00:05.000"
		submittedAt     = time.Date(2012, 1, 1, 12, 0, 2, 0, time.UTC)
		submittedAtStr  = "2012-01-01T12:00:02.000"
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messageheaders startindex="`+strconv.Itoa(startIndex)+`" count="`+strconv.Itoa(count)+`" totalcount="`+strconv.Itoa(totalCount)+`" xmlns="http://api.esendex.com/ns/">
 <messageheader id="`+id+`" uri="`+uri+`">
  <reference>`+reference+`</reference>
  <status>`+status+`</status>
  <laststatusat>`+lastStatusAtStr+`</laststatusat>
  <submittedat>`+submittedAtStr+`</submittedat>
  <type>`+string(messageType)+`</type>
  <to>
    <phonenumber>`+to+`</phonenumber>
  </to>
  <from>
   <phonenumber>`+from+`</phonenumber>
  </from>
  <summary>`+summary+`</summary>
  <body uri="`+bodyURI+`"/>
  <direction>`+direction+`</direction>
  <parts>`+strconv.Itoa(parts)+`</parts>
  <username>`+username+`</username>
 </messageheader>
</messageheaders>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	result, err := client.Sent()

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.0/messageheaders", h.Request.URL.String())

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	assert.Equal(startIndex, result.StartIndex)
	assert.Equal(count, result.Count)
	assert.Equal(totalCount, result.TotalCount)

	if assert.Equal(1, len(result.Messages)) {
		message := result.Messages[0]

		assert.Equal(id, message.ID)
		assert.Equal(uri, message.URI)
		assert.Equal(reference, message.Reference)
		assert.Equal(status, message.Status)
		assert.Equal(lastStatusAt, message.LastStatusAt)
		assert.Equal(submittedAt, message.SubmittedAt)
		assert.Equal(messageType, message.Type)
		assert.Equal(to, message.To)
		assert.Equal(from, message.From)
		assert.Equal(summary, message.Summary)
		assert.Equal(bodyURI, message.bodyURI)
		assert.Equal(direction, message.Direction)
		assert.Equal(parts, message.Parts)
		assert.Equal(username, message.Username)
	}
}

func TestMessagesSentWhenNotOkResponse(t *testing.T) {
	h := newRecordingHandler(``, 401, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	result, err := client.Sent()

	assert := assert.New(t)

	assert.Equal(err, ClientError{
		Method: "GET",
		Path:   "/v1.0/messageheaders",
		Code:   401,
	})
	assert.Nil(result)
}

func TestMessagesSentWithPaging(t *testing.T) {
	const (
		startIndex  = 0
		count       = 15
		totalCount  = 200
		id          = "messageheaderid"
		uri         = "http://somemessageheader"
		reference   = "EXETRTRE"
		status      = "STATUS"
		messageType = SMS
		to          = "4538224364236"
		from        = "428377843"
		summary     = "SUM"
		bodyURI     = "http://rrehekr"
		direction   = "OUT"
		parts       = 1
		username    = "user"
	)

	var (
		lastStatusAt    = time.Date(2012, 1, 1, 12, 0, 5, 0, time.UTC)
		lastStatusAtStr = "2012-01-01T12:00:05.000"
		submittedAt     = time.Date(2012, 1, 1, 12, 0, 2, 0, time.UTC)
		submittedAtStr  = "2012-01-01T12:00:02.000"
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messageheaders startindex="`+strconv.Itoa(startIndex)+`" count="`+strconv.Itoa(count)+`" totalcount="`+strconv.Itoa(totalCount)+`" xmlns="http://api.esendex.com/ns/">
 <messageheader id="`+id+`" uri="`+uri+`">
  <reference>`+reference+`</reference>
  <status>`+status+`</status>
  <laststatusat>`+lastStatusAtStr+`</laststatusat>
  <submittedat>`+submittedAtStr+`</submittedat>
  <type>`+string(messageType)+`</type>
  <to>
    <phonenumber>`+to+`</phonenumber>
  </to>
  <from>
   <phonenumber>`+from+`</phonenumber>
  </from>
  <summary>`+summary+`</summary>
  <body uri="`+bodyURI+`"/>
  <direction>`+direction+`</direction>
  <parts>`+strconv.Itoa(parts)+`</parts>
  <username>`+username+`</username>
 </messageheader>
</messageheaders>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	result, err := client.Sent(Page(5, 10))

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.0/messageheaders", h.Request.URL.Path)

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	query := h.Request.URL.Query()
	assert.Equal("5", query.Get("startindex"))
	assert.Equal("10", query.Get("count"))

	assert.Equal(startIndex, result.StartIndex)
	assert.Equal(count, result.Count)
	assert.Equal(totalCount, result.TotalCount)

	if assert.Equal(1, len(result.Messages)) {
		message := result.Messages[0]

		assert.Equal(id, message.ID)
		assert.Equal(uri, message.URI)
		assert.Equal(reference, message.Reference)
		assert.Equal(status, message.Status)
		assert.Equal(lastStatusAt, message.LastStatusAt)
		assert.Equal(submittedAt, message.SubmittedAt)
		assert.Equal(messageType, message.Type)
		assert.Equal(to, message.To)
		assert.Equal(from, message.From)
		assert.Equal(summary, message.Summary)
		assert.Equal(bodyURI, message.bodyURI)
		assert.Equal(direction, message.Direction)
		assert.Equal(parts, message.Parts)
		assert.Equal(username, message.Username)
	}
}

func TestMessagesByID(t *testing.T) {
	const (
		id          = "messageheaderid"
		uri         = "http://somemessageheader"
		reference   = "EXETRTRE"
		status      = "STATUS"
		messageType = Voice
		to          = "4538224364236"
		from        = "428377843"
		summary     = "SUM"
		bodyURI     = "http://rrehekr"
		direction   = "OUT"
		parts       = 1
		username    = "user"
		readBy      = "john.doe@example.com"
	)

	var (
		lastStatusAt    = time.Date(2012, 1, 1, 12, 0, 5, 0, time.UTC)
		lastStatusAtStr = "2012-01-01T12:00:05.000Z"
		submittedAt     = time.Date(2012, 1, 1, 12, 0, 2, 0, time.UTC)
		submittedAtStr  = "2012-01-01T12:00:02.000Z"
		receivedAt      = time.Date(2012, 1, 1, 12, 0, 1, 50000000, time.UTC)
		receivedAtStr   = "2012-01-01T12:00:01.05Z"
		readAt          = time.Date(2013, 1, 1, 12, 0, 1, 50000000, time.UTC)
		readAtStr       = "2013-01-01T12:00:01.05Z"
		sentAt          = time.Date(2013, 2, 1, 12, 0, 1, 50000000, time.UTC)
		sentAtStr       = "2013-02-01T12:00:01.05Z"
		deliveredAt     = time.Date(2013, 2, 2, 12, 0, 1, 50000000, time.UTC)
		deliveredAtStr  = "2013-02-02T12:00:01.05Z"
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messageheader id="`+id+`" uri="`+uri+`" xmlns="http://api.esendex.com/ns/">
 <reference>`+reference+`</reference>
 <status>`+status+`</status>
 <laststatusat>`+lastStatusAtStr+`</laststatusat>
 <submittedat>`+submittedAtStr+`</submittedat>
 <receivedat>`+receivedAtStr+`</receivedat>
 <type>`+string(messageType)+`</type>
 <to>
   <phonenumber>`+to+`</phonenumber>
 </to>
 <from>
  <phonenumber>`+from+`</phonenumber>
 </from>
 <summary>`+summary+`</summary>
 <body uri="`+bodyURI+`"/>
 <direction>`+direction+`</direction>
 <readat>`+readAtStr+`</readat>
 <sentat>`+sentAtStr+`</sentat>
 <deliveredat>`+deliveredAtStr+`</deliveredat>
 <readby>`+readBy+`</readby>
 <parts>`+strconv.Itoa(parts)+`</parts>
 <username>`+username+`</username>
</messageheader>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	result, err := client.Message(id)

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.0/messageheaders/"+id, h.Request.URL.String())

	assert.Equal(id, result.ID)
	assert.Equal(uri, result.URI)
	assert.Equal(reference, result.Reference)
	assert.Equal(status, result.Status)
	assert.Equal(lastStatusAt, result.LastStatusAt)
	assert.Equal(submittedAt, result.SubmittedAt)
	assert.Equal(receivedAt, result.ReceivedAt)
	assert.Equal(messageType, result.Type)
	assert.Equal(to, result.To)
	assert.Equal(from, result.From)
	assert.Equal(summary, result.Summary)
	assert.Equal(bodyURI, result.bodyURI)
	assert.Equal(direction, result.Direction)
	assert.Equal(readAt, result.ReadAt)
	assert.Equal(sentAt, result.SentAt)
	assert.Equal(deliveredAt, result.DeliveredAt)
	assert.Equal(readBy, result.ReadBy)
	assert.Equal(parts, result.Parts)
	assert.Equal(username, result.Username)
}

func TestMessagesReceived(t *testing.T) {
	const (
		startIndex  = 0
		count       = 15
		totalCount  = 200
		id          = "messageheaderid"
		uri         = "http://somemessageheader"
		reference   = "EXETRTRE"
		status      = "STATUS"
		messageType = SMS
		to          = "4538224364236"
		from        = "428377843"
		summary     = "SUM"
		bodyURI     = "http://rrehekr"
		direction   = "OUT"
		parts       = 1
		readBy      = "someone"
	)

	var (
		receivedAt    = time.Date(2012, 1, 1, 12, 0, 5, 0, time.UTC)
		receivedAtStr = "2012-01-01T12:00:05"
		readAt        = time.Date(2012, 1, 1, 12, 0, 2, 0, time.UTC)
		readAtStr     = "2012-01-01T12:00:02"
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messageheaders startindex="`+strconv.Itoa(startIndex)+`" count="`+strconv.Itoa(count)+`" totalcount="`+strconv.Itoa(totalCount)+`" xmlns="http://api.esendex.com/ns/">
 <messageheader id="`+id+`" uri="`+uri+`">
  <reference>`+reference+`</reference>
  <status>`+status+`</status>
  <receivedat>`+receivedAtStr+`</receivedat>
  <type>`+string(messageType)+`</type>
  <to>
    <phonenumber>`+to+`</phonenumber>
  </to>
  <from>
   <phonenumber>`+from+`</phonenumber>
  </from>
  <summary>`+summary+`</summary>
  <body uri="`+bodyURI+`"/>
  <direction>`+direction+`</direction>
  <parts>`+strconv.Itoa(parts)+`</parts>
  <readat>`+readAtStr+`</readat>
  <readby>`+readBy+`</readby>
 </messageheader>
</messageheaders>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	result, err := client.Received()

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.0/inbox/messages", h.Request.URL.String())

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	assert.Equal(startIndex, result.StartIndex)
	assert.Equal(count, result.Count)
	assert.Equal(totalCount, result.TotalCount)

	if assert.Equal(1, len(result.Messages)) {
		message := result.Messages[0]

		assert.Equal(id, message.ID)
		assert.Equal(uri, message.URI)
		assert.Equal(reference, message.Reference)
		assert.Equal(status, message.Status)
		assert.Equal(receivedAt, message.ReceivedAt)
		assert.Equal(messageType, message.Type)
		assert.Equal(to, message.To)
		assert.Equal(from, message.From)
		assert.Equal(summary, message.Summary)
		assert.Equal(bodyURI, message.bodyURI)
		assert.Equal(direction, message.Direction)
		assert.Equal(parts, message.Parts)
		assert.Equal(readAt, message.ReadAt)
		assert.Equal(readBy, message.ReadBy)
	}
}

func TestMessagesReceivedWithPaging(t *testing.T) {
	const (
		startIndex  = 0
		count       = 15
		totalCount  = 200
		id          = "messageheaderid"
		uri         = "http://somemessageheader"
		reference   = "EXETRTRE"
		status      = "STATUS"
		messageType = SMS
		to          = "4538224364236"
		from        = "428377843"
		summary     = "SUM"
		bodyURI     = "http://rrehekr"
		direction   = "OUT"
		parts       = 1
		readBy      = "someone"
	)

	var (
		receivedAt    = time.Date(2012, 1, 1, 12, 0, 5, 0, time.UTC)
		receivedAtStr = "2012-01-01T12:00:05"
		readAt        = time.Date(2012, 1, 1, 12, 0, 2, 0, time.UTC)
		readAtStr     = "2012-01-01T12:00:02"
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messageheaders startindex="`+strconv.Itoa(startIndex)+`" count="`+strconv.Itoa(count)+`" totalcount="`+strconv.Itoa(totalCount)+`" xmlns="http://api.esendex.com/ns/">
 <messageheader id="`+id+`" uri="`+uri+`">
  <reference>`+reference+`</reference>
  <status>`+status+`</status>
  <receivedat>`+receivedAtStr+`</receivedat>
  <type>`+string(messageType)+`</type>
  <to>
    <phonenumber>`+to+`</phonenumber>
  </to>
  <from>
   <phonenumber>`+from+`</phonenumber>
  </from>
  <summary>`+summary+`</summary>
  <body uri="`+bodyURI+`"/>
  <direction>`+direction+`</direction>
  <parts>`+strconv.Itoa(parts)+`</parts>
  <readat>`+readAtStr+`</readat>
  <readby>`+readBy+`</readby>
 </messageheader>
</messageheaders>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	result, err := client.Received(Page(5, 10))

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.0/inbox/messages", h.Request.URL.Path)

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	query := h.Request.URL.Query()
	assert.Equal("5", query.Get("startindex"))
	assert.Equal("10", query.Get("count"))

	assert.Equal(startIndex, result.StartIndex)
	assert.Equal(count, result.Count)
	assert.Equal(totalCount, result.TotalCount)

	if assert.Equal(1, len(result.Messages)) {
		message := result.Messages[0]

		assert.Equal(id, message.ID)
		assert.Equal(uri, message.URI)
		assert.Equal(reference, message.Reference)
		assert.Equal(status, message.Status)
		assert.Equal(receivedAt, message.ReceivedAt)
		assert.Equal(messageType, message.Type)
		assert.Equal(to, message.To)
		assert.Equal(from, message.From)
		assert.Equal(summary, message.Summary)
		assert.Equal(bodyURI, message.bodyURI)
		assert.Equal(direction, message.Direction)
		assert.Equal(parts, message.Parts)
		assert.Equal(readAt, message.ReadAt)
		assert.Equal(readBy, message.ReadBy)
	}
}

func TestMessagesReceivedWithDateRange(t *testing.T) {
	const (
		startIndex  = 0
		count       = 15
		totalCount  = 200
		id          = "messageheaderid"
		uri         = "http://somemessageheader"
		reference   = "EXETRTRE"
		status      = "STATUS"
		messageType = Voice
		to          = "4538224364236"
		from        = "428377843"
		summary     = "SUM"
		bodyURI     = "http://rrehekr"
		direction   = "OUT"
		parts       = 1
		readBy      = "someone"
	)

	var (
		receivedAt      = time.Date(2012, 1, 1, 12, 0, 5, 0, time.UTC)
		receivedAtStr   = "2012-01-01T12:00:05"
		readAt          = time.Date(2012, 1, 1, 12, 0, 2, 0, time.UTC)
		readAtStr       = "2012-01-01T12:00:02"
		messagesFrom    = time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC)
		messagesFromStr = "2012-01-01T00:00:00Z"
		messagesTo      = time.Date(2012, 6, 1, 0, 0, 0, 0, time.UTC)
		messagesToStr   = "2012-06-01T00:00:00Z"
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messageheaders startindex="`+strconv.Itoa(startIndex)+`" count="`+strconv.Itoa(count)+`" totalcount="`+strconv.Itoa(totalCount)+`" xmlns="http://api.esendex.com/ns/">
 <messageheader id="`+id+`" uri="`+uri+`">
  <reference>`+reference+`</reference>
  <status>`+status+`</status>
  <receivedat>`+receivedAtStr+`</receivedat>
  <type>`+string(messageType)+`</type>
  <to>
    <phonenumber>`+to+`</phonenumber>
  </to>
  <from>
   <phonenumber>`+from+`</phonenumber>
  </from>
  <summary>`+summary+`</summary>
  <body uri="`+bodyURI+`"/>
  <direction>`+direction+`</direction>
  <parts>`+strconv.Itoa(parts)+`</parts>
  <readat>`+readAtStr+`</readat>
  <readby>`+readBy+`</readby>
 </messageheader>
</messageheaders>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	result, err := client.Received(Between(messagesFrom, messagesTo))

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.0/inbox/messages", h.Request.URL.Path)

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	query := h.Request.URL.Query()
	assert.Equal(messagesFromStr, query.Get("start"))
	assert.Equal(messagesToStr, query.Get("finish"))

	assert.Equal(startIndex, result.StartIndex)
	assert.Equal(count, result.Count)
	assert.Equal(totalCount, result.TotalCount)

	if assert.Equal(1, len(result.Messages)) {
		message := result.Messages[0]

		assert.Equal(id, message.ID)
		assert.Equal(uri, message.URI)
		assert.Equal(reference, message.Reference)
		assert.Equal(status, message.Status)
		assert.Equal(receivedAt, message.ReceivedAt)
		assert.Equal(messageType, message.Type)
		assert.Equal(to, message.To)
		assert.Equal(from, message.From)
		assert.Equal(summary, message.Summary)
		assert.Equal(bodyURI, message.bodyURI)
		assert.Equal(direction, message.Direction)
		assert.Equal(parts, message.Parts)
		assert.Equal(readAt, message.ReadAt)
		assert.Equal(readBy, message.ReadBy)
	}
}

func TestSentMessageBody(t *testing.T) {
	const (
		bodyText     = "Hey there"
		characterSet = "GSM"
		bodyPath     = "/message/body"
	)

	h := newRecordingHandler(fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<messagebody xmlns="http://api.esendex.com/ns/">
    <bodytext>%s</bodytext>
    <characterset>%s</characterset>
</messagebody>`, bodyText, characterSet), 200, map[string]string{})

	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	message := SentMessageResponse{bodyURI: "https://example.com" + bodyPath}

	body, err := client.Body(message)

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal(bodyPath, h.Request.URL.Path)

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	assert.Equal(bodyText, body.Text)
	assert.Equal(characterSet, body.CharacterSet)
}

func TestReceivedMessageBody(t *testing.T) {
	const (
		bodyText     = "Hey there"
		characterSet = "GSM"
		bodyPath     = "/message/body"
	)

	h := newRecordingHandler(fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<messagebody xmlns="http://api.esendex.com/ns/">
    <bodytext>%s</bodytext>
    <characterset>%s</characterset>
</messagebody>`, bodyText, characterSet), 200, map[string]string{})

	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	message := ReceivedMessageResponse{bodyURI: "https://example.com" + bodyPath}

	body, err := client.Body(message)

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal(bodyPath, h.Request.URL.Path)

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	assert.Equal(bodyText, body.Text)
	assert.Equal(characterSet, body.CharacterSet)
}

func TestMessageBody(t *testing.T) {
	const (
		bodyText     = "Hey there"
		characterSet = "GSM"
		bodyPath     = "/message/body"
	)

	h := newRecordingHandler(fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<messagebody xmlns="http://api.esendex.com/ns/">
    <bodytext>%s</bodytext>
    <characterset>%s</characterset>
</messagebody>`, bodyText, characterSet), 200, map[string]string{})

	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	message := MessageResponse{bodyURI: "https://example.com" + bodyPath}

	body, err := client.Body(message)

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal(bodyPath, h.Request.URL.Path)

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	assert.Equal(bodyText, body.Text)
	assert.Equal(characterSet, body.CharacterSet)
}
