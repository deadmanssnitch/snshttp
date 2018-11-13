package snshttp

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

// Notification events are sent for messages that are published to the SNS
// topic.
type Notification struct {
	MessageID      string    `json:"MessageId"`
	TopicARN       string    `json:"TopicArn"`
	Subject        string    `json:"Subject"`
	Message        string    `json:"Message"`
	Timestamp      time.Time `json:"Timestamp"`
	UnsubscribeURL string    `json:"UnsubscribeURL"`

	// MessageAttributes contain any attributes added to the message when
	// publishing it to SNS. This is most commonly used when transmitting binary
	// date (using raw message delivery).
	MessageAttributes map[string]MessageAttribute `json:"MessageAttributes"`
}

// Unsubscribe will notify Amazon to remove this subscription from the SNS
// topic. It will make a request to the UnsubscribeURL and error if the
// request times out or the response does not indicate success.
func (e *Notification) Unsubscribe(ctx context.Context) error {
	req, err := http.NewRequest("GET", e.UnsubscribeURL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Server is expected to return 200 OK but we can treat any 200 level code as
	// success.
	if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		return fmt.Errorf("server returned error status=%d", resp.StatusCode)
	}

	return nil
}

type MessageAttribute struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

func (m MessageAttribute) StringValue() string {
	return m.Value
}

func (m MessageAttribute) BinaryValue() ([]byte, error) {
	return base64.StdEncoding.DecodeString(m.Value)
}
