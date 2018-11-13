package snshttp

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// SubscriptionConfirmation is an initial event sent by Amazon SNS as part of a
// handshake before any Notification events can be sent. A call to Confirm or a
// request to SubscribeURL will finish the handshake and enable Amazon to send
// Notifications to the webhook.
type SubscriptionConfirmation struct {
	MessageID string    `json:"MessageId"`
	TopicARN  string    `json:"TopicArn"`
	Timestamp time.Time `json:"Timestamp"`

	Token        string `json:"Token"`
	Message      string `json:"Message"`
	SubscribeURL string `json:"SubscribeURL"`
}

// Confirm finishes the handshake with Amazon, confirming that the subscription
// should start sending notification. A request is made to the SubscribeURL. An
// error will be returned if the subscription has already been confirmed.
func (e *SubscriptionConfirmation) Confirm(ctx context.Context) error {
	req, err := http.NewRequest("GET", e.SubscribeURL, nil)
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
