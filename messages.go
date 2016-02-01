package xesende

import (
	"encoding/xml"
	"net/url"
	"time"
)

// Paging gives the details of the page that was accessed.
type Paging struct {
	StartIndex int
	Count      int
	TotalCount int
}

// FailureReason gives detailed information for why a message failed to send.
type FailureReason struct {
	Code        int
	Description string
	Permanent   bool
}

type messageWithBody interface {
	getBodyURI() string
}

// SentMessagesResponse is a list of returned messages along with the paging
// information.
type SentMessagesResponse struct {
	Paging
	Messages []SentMessageResponse
}

// SentMessageResponse is a single sent message. It implements messageWithBody.
type SentMessageResponse struct {
	ID            string
	URI           string
	Reference     string
	Status        string
	LastStatusAt  time.Time
	SubmittedAt   time.Time
	Type          MessageType
	To            string
	From          string
	Summary       string
	Direction     string
	Parts         int
	Username      string
	FailureReason *FailureReason

	bodyURI string
}

func (r SentMessageResponse) getBodyURI() string { return r.bodyURI }

// MessageResponse is a single message. It implements messageWithBody.
type MessageResponse struct {
	ID           string
	URI          string
	Reference    string
	Status       string
	LastStatusAt time.Time
	SubmittedAt  time.Time
	ReceivedAt   time.Time
	Type         MessageType
	To           string
	From         string
	Summary      string
	Direction    string
	ReadAt       time.Time
	SentAt       time.Time
	DeliveredAt  time.Time
	ReadBy       string
	Parts        int
	Username     string
	FailureReason *FailureReason

	bodyURI string
}

func (r MessageResponse) getBodyURI() string { return r.bodyURI }

// ReceivedMessagesResponse is a list of received messages along with the paging
// information.
type ReceivedMessagesResponse struct {
	Paging
	Messages []ReceivedMessageResponse
}

// ReceivedMessageResponse is a single received message. It implements messageWithBody.
type ReceivedMessageResponse struct {
	ID         string
	URI        string
	Reference  string
	Status     string
	ReceivedAt time.Time
	Type       MessageType
	To         string
	From       string
	Summary    string
	Direction  string
	Parts      int
	ReadAt     time.Time
	ReadBy     string

	bodyURI string
}

func (r ReceivedMessageResponse) getBodyURI() string { return r.bodyURI }

// MessageBody is the body of a message.
type MessageBody struct {
	Text         string
	CharacterSet string
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
			Type:         MessageType(message.Type),
			To:           message.To,
			From:         message.From,
			Summary:      message.Summary,
			bodyURI:      message.Body.URI,
			Direction:    message.Direction,
			Parts:        message.Parts,
			Username:     message.Username,
		}

		if message.FailureReason != nil {
			response.Messages[i].FailureReason = &FailureReason{
				Code:        message.FailureReason.Code,
				Description: message.FailureReason.Description,
				Permanent:   message.FailureReason.Permanent,
			}
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
			Type:       MessageType(message.Type),
			To:         message.To,
			From:       message.From,
			Summary:    message.Summary,
			bodyURI:    message.Body.URI,
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
		Type:         MessageType(v.Type),
		To:           v.To,
		From:         v.From,
		Summary:      v.Summary,
		bodyURI:      v.Body.URI,
		Direction:    v.Direction,
		ReadAt:       v.ReadAt.Time,
		SentAt:       v.SentAt.Time,
		DeliveredAt:  v.DeliveredAt.Time,
		ReadBy:       v.ReadBy,
		Parts:        v.Parts,
		Username:     v.Username,
	}

	if v.FailureReason != nil {
		response.FailureReason = &FailureReason{
			Code:        v.FailureReason.Code,
			Description: v.FailureReason.Description,
			Permanent:   v.FailureReason.Permanent,
		}
	}

	return response, nil
}

// Body returns the full body of a single message.
func (c *Client) Body(message messageWithBody) (*MessageBody, error) {
	u, err := url.Parse(message.getBodyURI())
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest("GET", u.Path, nil)
	if err != nil {
		return nil, err
	}

	var v messageBodyResponse
	if _, err = c.do(req, &v); err != nil {
		return nil, err
	}

	return &MessageBody{
		Text:         v.BodyText,
		CharacterSet: v.CharacterSet,
	}, nil
}

type messageBodyResponse struct {
	XMLName      xml.Name `xml:"http://api.esendex.com/ns/ messagebody"`
	BodyText     string   `xml:"bodytext"`
	CharacterSet string   `xml:"characterset"`
}

type messageFailureReason struct {
	Code        int    `xml:"code"`
	Description string `xml:"description"`
	Permanent   bool   `xml:"permanentfailure"`
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
	Direction     string                `xml:"direction"`
	ReadAt        messageHeaderTime     `xml:"readat"`
	SentAt        messageHeaderTime     `xml:"sentat"`
	DeliveredAt   messageHeaderTime     `xml:"deliveredat"`
	ReadBy        string                `xml:"readby"`
	Parts         int                   `xml:"parts"`
	Username      string                `xml:"username"`
	FailureReason *messageFailureReason `xml:"failurereason"`
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
