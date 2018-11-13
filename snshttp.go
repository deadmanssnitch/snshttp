// Package snshttp provides an HTTP handler to make it easier to work with
// webhooks from Amazon SNS.
package snshttp

import "context"

// EventHandler methods are called for each event received from Amazon SNS.
type EventHandler interface {
	// SubscriptionConfirmation is the first event received for a new Amazon SNS
	// webhook subscription and contains information on how to confirm the
	// subscription.
	SubscriptionConfirmation(ctx context.Context, event *SubscriptionConfirmation) error

	// Notification events contain the messages published to the SNS topic. This
	// is the most common type of event.
	Notification(ctx context.Context, event *Notification) error

	// UnsubscribeConfirmation events are sent when the subscription has been
	// canceled and gives the consumer a chance to resubscribe. Note that the
	// UnsubscribeConfirmation event is not sent when deleting a subscription
	// from the console.
	UnsubscribeConfirmation(ctx context.Context, event *UnsubscribeConfirmation) error
}

// DefaultHandler is intended to be mixed in to a struct to provide the most
// common behavior for event receivers. Adding the DefaultHandler to a struct
// will automatically confirm subscriptions and ignore any unsubscribe
// confirmations.
type DefaultHandler struct{}

// SubscriptionConfirmation confirms the subscription.
func (h DefaultHandler) SubscriptionConfirmation(ctx context.Context, event *SubscriptionConfirmation) error {
	return event.Confirm(ctx)
}

// UnsubscribeConfirmation does nothing and ignores the event.
func (h DefaultHandler) UnsubscribeConfirmation(_ context.Context, _ *UnsubscribeConfirmation) error {
	return nil
}
