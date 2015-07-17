package xesende

import (
	"encoding/xml"
	"errors"
)

type Messages []Message
type Message struct {
	To   string
	Body string
}

type SendResponse struct {
	BatchId  string
	Messages []SendResponseMessage
}

type SendResponseMessage struct {
	Uri string
	Id  string
}

// Send dispatches a list of messages.
func (c *AccountClient) Send(messages Messages) (*SendResponse, error) {
	body := messageDispatchRequest{
		AccountReference: c.reference,
		Message:          make([]messageDispatchRequestMessage, len(messages)),
	}

	for i, message := range messages {
		body.Message[i] = messageDispatchRequestMessage{To: message.To, Body: message.Body}
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
		BatchId:  v.BatchId,
		Messages: make([]SendResponseMessage, len(v.MessageHeader)),
	}

	for i, message := range v.MessageHeader {
		response.Messages[i] = SendResponseMessage{
			Uri: message.Uri,
			Id:  message.Id,
		}
	}

	return response, nil
}

type messageDispatchRequest struct {
	XMLName          xml.Name                        `xml:"messages"`
	AccountReference string                          `xml:"accountreference"`
	Message          []messageDispatchRequestMessage `xml:"message"`
}

type messageDispatchRequestMessage struct {
	To   string `xml:"to"`
	Body string `xml:"body"`
}

type messageDispatchResponse struct {
	XMLName       xml.Name `xml:"http://api.esendex.com/ns/ messageheaders"`
	BatchId       string   `xml:"batchid,attr"`
	MessageHeader []struct {
		Uri string `xml:"uri,attr"`
		Id  string `xml:"id,attr"`
	} `xml:"messageheader"`
}
