package xesende

import (
	"encoding/xml"
	"time"
)

// Accounts returns a list of accounts the user is on.
func (c *Client) Accounts() (*AccountsResponse, error) {
	req, err := c.newRequest("GET", "/v1.0/accounts", nil)
	if err != nil {
		return nil, err
	}

	var v accountsResponse
	if _, err = c.do(req, &v); err != nil {
		return nil, err
	}

	response := &AccountsResponse{
		Accounts: make([]AccountResponse, len(v.Accounts)),
	}

	for i, account := range v.Accounts {
		response.Accounts[i] = AccountResponse{
			ID:                account.ID,
			URI:               account.URI,
			Reference:         account.Reference,
			Label:             account.Label,
			Address:           account.Address,
			Type:              account.Type,
			MessagesRemaining: account.MessagesRemaining,
			ExpiresOn:         account.ExpiresOn.Time,
			Role:              account.Role,
			SettingsURI:       account.Settings.URI,
		}
	}

	return response, nil
}

// AccountsResponse is a list of accounts.
type AccountsResponse struct {
	Accounts []AccountResponse
}

// AccountResponse is a single account.
type AccountResponse struct {
	ID                string
	URI               string
	Reference         string
	Label             string
	Address           string
	Type              string
	MessagesRemaining int
	ExpiresOn         time.Time
	Role              string
	SettingsURI       string
}

type accountsResponse struct {
	XMLName  xml.Name                  `xml:"http://api.esendex.com/ns/ accounts"`
	Accounts []accountsResponseAccount `xml:"account"`
}

type accountsResponseAccount struct {
	ID                string       `xml:"id,attr"`
	URI               string       `xml:"uri,attr"`
	Reference         string       `xml:"reference"`
	Label             string       `xml:"label"`
	Address           string       `xml:"address"`
	Type              string       `xml:"type"`
	MessagesRemaining int          `xml:"messagesremaining"`
	ExpiresOn         accountsTime `xml:"expireson"`
	Role              string       `xml:"role"`
	Settings          struct {
		URI string `xml:"uri,attr"`
	} `xml:"settings"`
}

const accountsTimeFormat = "2006-01-02T15:04:05"

type accountsTime struct {
	time.Time
}

func (t accountsTime) MarshalText() ([]byte, error) {
	return []byte(t.Format(accountsTimeFormat)), nil
}

func (t *accountsTime) UnmarshalText(data []byte) error {
	g, err := time.ParseInLocation(accountsTimeFormat, string(data), time.UTC)
	if err != nil {
		return err
	}
	*t = accountsTime{g}
	return nil
}
