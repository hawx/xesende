package pushnotifications

import (
	"encoding/xml"
	"net/http"
	"time"
)

type ReceivedNotification struct {
	Id          string
	MessageId   string
	AccountId   string
	MessageText string
	From        string
	To          string
}

type receivedNotificationResponse struct {
	Id          string `xml:"Id"`
	MessageId   string `xml:"MessageId"`
	AccountId   string `xml:"AccountId"`
	MessageText string `xml:"MessageText"`
	From        string `xml:"From"`
	To          string `xml:"To"`
}

func Received(f func(ReceivedNotification)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var v receivedNotificationResponse

		if err := xml.NewDecoder(r.Body).Decode(&v); err != nil {
			w.WriteHeader(500)
			return
		}

		f(ReceivedNotification{
			Id:          v.Id,
			MessageId:   v.MessageId,
			AccountId:   v.AccountId,
			MessageText: v.MessageText,
			From:        v.From,
			To:          v.To,
		})
	})
}

type DeliveredNotification struct {
	Id         string
	MessageId  string
	AccountId  string
	OccurredAt time.Time
}

type deliveredNotificationResponse struct {
	Id         string           `xml:"Id"`
	MessageId  string           `xml:"MessageId"`
	AccountId  string           `xml:"AccountId"`
	OccurredAt notificationTime `xml:"OccurredAt"`
}

func Delivered(f func(DeliveredNotification)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var v deliveredNotificationResponse

		if err := xml.NewDecoder(r.Body).Decode(&v); err != nil {
			w.WriteHeader(500)
			return
		}

		f(DeliveredNotification{
			Id:         v.Id,
			MessageId:  v.MessageId,
			AccountId:  v.AccountId,
			OccurredAt: v.OccurredAt.Time,
		})
	})
}

type FailedNotification struct {
	Id         string
	MessageId  string
	AccountId  string
	OccurredAt time.Time
}

type failedNotificationResponse struct {
	Id         string           `xml:"Id"`
	MessageId  string           `xml:"MessageId"`
	AccountId  string           `xml:"AccountId"`
	OccurredAt notificationTime `xml:"OccurredAt"`
}

func Failed(f func(FailedNotification)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var v failedNotificationResponse

		if err := xml.NewDecoder(r.Body).Decode(&v); err != nil {
			w.WriteHeader(500)
			return
		}

		f(FailedNotification{
			Id:         v.Id,
			MessageId:  v.MessageId,
			AccountId:  v.AccountId,
			OccurredAt: v.OccurredAt.Time,
		})
	})
}

const notificationTimeFormat = "2006-01-02T15:04:05"

type notificationTime struct {
	time.Time
}

func (t notificationTime) MarshalText() ([]byte, error) {
	return []byte(t.Format(notificationTimeFormat)), nil
}

func (t *notificationTime) UnmarshalText(data []byte) error {
	g, err := time.ParseInLocation(notificationTimeFormat, string(data), time.UTC)
	if err != nil {
		return err
	}
	*t = notificationTime{g}
	return nil
}
