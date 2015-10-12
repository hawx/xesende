package xesende

import (
	"encoding/xml"
	"errors"
	"time"
)

// Message is a message to send.
type Message struct {
	To           string
	MessageType  string
	Lang         string
	Validity     int
	CharacterSet string
	Retries      int
	Body         string
}

// SendResponse gives the batchid for the sent batch and lists the details of
// each message sent.
type SendResponse struct {
	BatchID  string
	Messages []SendResponseMessage
}

// SendResponseMessage gives the details for a single sent message.
type SendResponseMessage struct {
	URI string
	ID  string
}

// Send dispatches a list of messages.
func (c *AccountClient) Send(messages []Message) (*SendResponse, error) {
	body := messageDispatchRequest{
		AccountReference: c.reference,
		Message:          make([]messageDispatchRequestMessage, len(messages)),
	}

	for i, message := range messages {
		body.Message[i] = messageDispatchRequestMessage{
			To:           message.To,
			MessageType:  message.MessageType,
			Lang:         message.Lang,
			Validity:     message.Validity,
			CharacterSet: message.CharacterSet,
			Retries:      message.Retries,
			Body:         message.Body,
		}
	}

	req, err := c.newRequest("POST", "/v1.0/messagedispatcher", &body)
	if err != nil {
		return nil, err
	}

	var v messageDispatchResponse
	resp, err := c.do(req, &v)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Expected 200")
	}

	response := &SendResponse{
		BatchID:  v.BatchID,
		Messages: make([]SendResponseMessage, len(v.MessageHeader)),
	}

	for i, message := range v.MessageHeader {
		response.Messages[i] = SendResponseMessage{
			URI: message.URI,
			ID:  message.ID,
		}
	}

	return response, nil
}

// Send dispatches a list of messages.
func (c *AccountClient) SendAt(sendAt time.Time, messages []Message) (*SendResponse, error) {
	body := messageDispatchRequest{
		AccountReference: c.reference,
		SendAt:           &sendAt,
		Message:          make([]messageDispatchRequestMessage, len(messages)),
	}

	for i, message := range messages {
		body.Message[i] = messageDispatchRequestMessage{
			To:           message.To,
			MessageType:  message.MessageType,
			Lang:         message.Lang,
			Validity:     message.Validity,
			CharacterSet: message.CharacterSet,
			Retries:      message.Retries,
			Body:         message.Body,
		}
	}

	req, err := c.newRequest("POST", "/v1.0/messagedispatcher", &body)
	if err != nil {
		return nil, err
	}

	var v messageDispatchResponse
	resp, err := c.do(req, &v)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Expected 200")
	}

	response := &SendResponse{
		BatchID:  v.BatchID,
		Messages: make([]SendResponseMessage, len(v.MessageHeader)),
	}

	for i, message := range v.MessageHeader {
		response.Messages[i] = SendResponseMessage{
			URI: message.URI,
			ID:  message.ID,
		}
	}

	return response, nil
}

type messageDispatchRequest struct {
	XMLName          xml.Name                        `xml:"messages"`
	AccountReference string                          `xml:"accountreference"`
	SendAt           *time.Time                      `xml:"sendat"`
	Message          []messageDispatchRequestMessage `xml:"message"`
}

type messageDispatchRequestMessage struct {
	To           string `xml:"to"`
	MessageType  string `xml:"type"`
	Lang         string `xml:"lang"`
	Validity     int    `xml:"validity"`
	CharacterSet string `xml:"characterset"`
	Retries      int    `xml:"retries"`
	Body         string `xml:"body"`
}

type messageDispatchResponse struct {
	XMLName       xml.Name `xml:"http://api.esendex.com/ns/ messageheaders"`
	BatchID       string   `xml:"batchid,attr"`
	MessageHeader []struct {
		URI string `xml:"uri,attr"`
		ID  string `xml:"id,attr"`
	} `xml:"messageheader"`
}
