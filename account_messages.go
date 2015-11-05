package xesende

import "net/http"

// Sent returns a list of messages sent by the account.
func (c *AccountClient) Sent(opts ...Option) (*SentMessagesResponse, error) {
	accountOption := func(r *http.Request) {
		q := r.URL.Query()

		q.Add("accountReference", c.reference)

		r.URL.RawQuery = q.Encode()
	}

	return c.Client.Sent(append(opts, accountOption)...)
}

// Received returns the messages sent to the account.
func (c *AccountClient) Received(opts ...Option) (*ReceivedMessagesResponse, error) {
	req, err := c.newRequest("GET", "/v1.0/inbox/"+c.reference+"/messages", nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(req)
	}

	var v inboxResponse
	if _, err = c.do(req, &v); err != nil {
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
