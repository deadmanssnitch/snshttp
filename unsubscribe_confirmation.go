package snshttp

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// UnsubscribeConfirmation events are received when a subscription is canceled
// via the API. No unsubscribe event is fired when deleting a subscription
// through the AWS web console.
type UnsubscribeConfirmation struct {
	MessageID string    `json:"MessageId"`
	TopicARN  string    `json:"TopicArn"`
	Timestamp time.Time `json:"Timestamp"`

	Token        string `json:"Token"`
	Message      string `json:"Message"`
	SubscribeURL string `json:"SubscribeURL"`
}

// Resubscribe notifies Amazon to reinstate the subscription. A request is made
// to the SubscribeURL.
func (e *UnsubscribeConfirmation) Resubscribe(ctx context.Context) error {
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
