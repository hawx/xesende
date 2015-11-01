package xesende

import (
	"encoding/xml"
	"time"
)

// Paging gives the details of the page that was accessed.
type Paging struct {
	StartIndex int
	Count      int
	TotalCount int
}

// MessagesResponse is a list of returned messages along with the paging
// information.
type SentMessagesResponse struct {
	Paging
	Messages []SentMessageResponse
}

// MessageResponse is a single sent message.
type SentMessageResponse struct {
	ID           string
	URI          string
	Reference    string
	Status       string
	LastStatusAt time.Time
	SubmittedAt  time.Time
	Type         string
	To           string
	From         string
	Summary      string
	BodyURI      string
	Direction    string
	Parts        int
	Username     string
}

// MessageResponse is a single message.
type MessageResponse struct {
	ID           string
	URI          string
	Reference    string
	Status       string
	LastStatusAt time.Time
	SubmittedAt  time.Time
	ReceivedAt   time.Time
	Type         string
	To           string
	From         string
	Summary      string
	BodyURI      string
	Direction    string
	ReadAt       time.Time
	SentAt       time.Time
	DeliveredAt  time.Time
	ReadBy       string
	Parts        int
	Username     string
}

// MessagesReceivedResponse is a list of received messages along with the paging
// information.
type ReceivedMessagesResponse struct {
	Paging
	Messages []ReceivedMessageResponse
}

// MessageReceivedResponse is a single received message.
type ReceivedMessageResponse struct {
	ID         string
	URI        string
	Reference  string
	Status     string
	ReceivedAt time.Time
	Type       string
	To         string
	From       string
	Summary    string
	BodyURI    string
	Direction  string
	Parts      int
	ReadAt     time.Time
	ReadBy     string
}

// Sent returns a list of messages sent by the user.
func (c *Client) Sent(opts ...Option) (*SentMessagesResponse, error) {
	req, err := c.newRequest("GET", "/v1.0/messageheaders", nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(req)
	}

	var v messageHeadersResponse

	if _, err := c.do(req, &v); err != nil {
		return nil, err
	}

	response := &SentMessagesResponse{
		Paging: Paging{
			StartIndex: v.StartIndex,
			Count:      v.Count,
			TotalCount: v.TotalCount,
		},
		Messages: make([]SentMessageResponse, len(v.Messages)),
	}

	for i, message := range v.Messages {
		response.Messages[i] = SentMessageResponse{
			ID:           message.ID,
			URI:          message.URI,
			Reference:    message.Reference,
			Status:       message.Status,
			LastStatusAt: message.LastStatusAt.Time,
			SubmittedAt:  message.SubmittedAt.Time,
			Type:         message.Type,
			To:           message.To,
			From:         message.From,
			Summary:      message.Summary,
			BodyURI:      message.Body.URI,
			Direction:    message.Direction,
			Parts:        message.Parts,
			Username:     message.Username,
		}
	}

	return response, nil
}

// Received returns the messages sent to the user.
func (c *Client) Received(opts ...Option) (*ReceivedMessagesResponse, error) {
	req, err := c.newRequest("GET", "/v1.0/inbox/messages", nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(req)
	}

	var v inboxResponse
	if _, err := c.do(req, &v); err != nil {
		return nil, err
	}

	response := &ReceivedMessagesResponse{
		Paging: Paging{
			StartIndex: v.StartIndex,
			Count:      v.Count,
			TotalCount: v.TotalCount,
		},
		Messages: make([]ReceivedMessageResponse, len(v.Messages)),
	}

	for i, message := range v.Messages {
		response.Messages[i] = ReceivedMessageResponse{
			ID:         message.ID,
			URI:        message.URI,
			Reference:  message.Reference,
			Status:     message.Status,
			ReceivedAt: message.ReceivedAt.Time,
			Type:       message.Type,
			To:         message.To,
			From:       message.From,
			Summary:    message.Summary,
			BodyURI:    message.Body.URI,
			Direction:  message.Direction,
			Parts:      message.Parts,
			ReadAt:     message.ReadAt.Time,
			ReadBy:     message.ReadBy,
		}
	}

	return response, nil
}

// Message returns the message with the given id.
func (c *Client) Message(id string) (*MessageResponse, error) {
	req, err := c.newRequest("GET", "/v1.0/messageheaders/"+id, nil)
	if err != nil {
		return nil, err
	}

	var v messageHeadersResponseMessageHeader
	if _, err = c.do(req, &v); err != nil {
		return nil, err
	}

	response := &MessageResponse{
		ID:           v.ID,
		URI:          v.URI,
		Reference:    v.Reference,
		Status:       v.Status,
		LastStatusAt: v.LastStatusAt.Time,
		SubmittedAt:  v.SubmittedAt.Time,
		ReceivedAt:   v.ReceivedAt.Time,
		Type:         v.Type,
		To:           v.To,
		From:         v.From,
		Summary:      v.Summary,
		BodyURI:      v.Body.URI,
		Direction:    v.Direction,
		ReadAt:       v.ReadAt.Time,
		SentAt:       v.SentAt.Time,
		DeliveredAt:  v.DeliveredAt.Time,
		ReadBy:       v.ReadBy,
		Parts:        v.Parts,
		Username:     v.Username,
	}

	return response, nil
}

type messageHeadersResponse struct {
	XMLName    xml.Name                              `xml:"http://api.esendex.com/ns/ messageheaders"`
	StartIndex int                                   `xml:"startindex,attr"`
	Count      int                                   `xml:"count,attr"`
	TotalCount int                                   `xml:"totalcount,attr"`
	Messages   []messageHeadersResponseMessageHeader `xml:"messageheader"`
}

type messageHeadersResponseMessageHeader struct {
	ID           string            `xml:"id,attr"`
	URI          string            `xml:"uri,attr"`
	Reference    string            `xml:"reference"`
	Status       string            `xml:"status"`
	LastStatusAt messageHeaderTime `xml:"laststatusat"`
	SubmittedAt  messageHeaderTime `xml:"submittedat"`
	ReceivedAt   messageHeaderTime `xml:"receivedat"`
	Type         string            `xml:"type"`
	To           string            `xml:"to>phonenumber"`
	From         string            `xml:"from>phonenumber"`
	Summary      string            `xml:"summary"`
	Body         struct {
		URI string `xml:"uri,attr"`
	} `xml:"body"`
	Direction   string            `xml:"direction"`
	ReadAt      messageHeaderTime `xml:"readat"`
	SentAt      messageHeaderTime `xml:"sentat"`
	DeliveredAt messageHeaderTime `xml:"deliveredat"`
	ReadBy      string            `xml:"readby"`
	Parts       int               `xml:"parts"`
	Username    string            `xml:"username"`
}

type inboxResponse struct {
	XMLName    xml.Name                     `xml:"http://api.esendex.com/ns/ messageheaders"`
	StartIndex int                          `xml:"startindex,attr"`
	Count      int                          `xml:"count,attr"`
	TotalCount int                          `xml:"totalcount,attr"`
	Messages   []inboxResponseMessageHeader `xml:"messageheader"`
}

type inboxResponseMessageHeader struct {
	ID         string            `xml:"id,attr"`
	URI        string            `xml:"uri,attr"`
	Reference  string            `xml:"reference"`
	Status     string            `xml:"status"`
	ReceivedAt messageHeaderTime `xml:"receivedat"`
	Type       string            `xml:"type"`
	To         string            `xml:"to>phonenumber"`
	From       string            `xml:"from>phonenumber"`
	Summary    string            `xml:"summary"`
	Body       struct {
		URI string `xml:"uri,attr"`
	} `xml:"body"`
	Direction string            `xml:"direction"`
	Parts     int               `xml:"parts"`
	ReadAt    messageHeaderTime `xml:"readat"`
	ReadBy    string            `xml:"readby"`
}

const messageHeaderTimeFormat = "2006-01-02T15:04:05.999999999"
const messageHeaderTimeFormatZ = "2006-01-02T15:04:05.999999999Z"

type messageHeaderTime struct {
	time.Time
}

func (t messageHeaderTime) MarshalText() ([]byte, error) {
	return []byte(t.Format(messageHeaderTimeFormat)), nil
}

func (t *messageHeaderTime) UnmarshalText(data []byte) error {
	g, err := time.ParseInLocation(messageHeaderTimeFormat, string(data), time.UTC)
	if err != nil {
		g, err = time.ParseInLocation(messageHeaderTimeFormatZ, string(data), time.UTC)
		if err != nil {
			return err
		}
	}
	*t = messageHeaderTime{g}
	return nil
}
