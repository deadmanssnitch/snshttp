package snshttp

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// requestTimeout is how long Amazon will wait from a response from the
	// server before considering the request failed. This does not appear to be
	// configurable, see the documentation below.
	//
	// https://docs.aws.amazon.com/sns/latest/dg/DeliveryPolicies.html#delivery-policy-maximum-receive-rate
	requestTimeout = 15 * time.Second
)

type handler struct {
	handler     EventHandler
	credentials *authOption
}

// New creates a http.Handler for receiving webhooks from an Amazon SNS
// subscription and dispatching them to the EventHandler. Options are applied
// in the order they're provided and may clobber previous options.
func New(eventHandler EventHandler, opts ...Option) http.Handler {
	handler := &handler{
		handler: eventHandler,
	}

	for _, opt := range opts {
		opt.apply(handler)
	}

	return handler
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var err error

	if !h.credentials.Check(req) {
		resp.Header().Set("WWW-Authenticate", `Basic realm="ses"`)
		http.Error(resp, "Unauthorized", http.StatusUnauthorized)

		return
	}

	// Amazon will consider a request failed if it takes longer than 15 seconds
	// to execute. This does not appear to be configurable.
	ctx, cancel := context.WithTimeout(req.Context(), requestTimeout)
	defer cancel()

	// Wrap the Request with the timeout so it's enforced when reading the body.
	req = req.WithContext(ctx)

	// Use the Type header so we can avoid parsing the body unless we know it's
	// an event we support.
	switch req.Header.Get("X-Amz-Sns-Message-Type") {

	// Notifications should be the most common case and switch statements are
	// checked in definition order.
	case "Notification":
		event := &Notification{}

		err = readEvent(req.Body, event)
		if err != nil {
			break
		}

		err = h.handler.Notification(ctx, event)

	case "SubscriptionConfirmation":
		event := &SubscriptionConfirmation{}

		err = readEvent(req.Body, event)
		if err != nil {
			break
		}

		err = h.handler.SubscriptionConfirmation(ctx, event)

	case "UnsubscribeConfirmation":
		event := &UnsubscribeConfirmation{}

		err = readEvent(req.Body, event)
		if err != nil {
			break
		}

		err = h.handler.UnsubscribeConfirmation(ctx, event)

	// Amazon (or someone else?) sent an unknown type
	default:
		http.NotFound(resp, req)
		return
	}

	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	// Success! Signals Amazon to mark message as received.
	resp.WriteHeader(http.StatusNoContent)
}

// readEvent reads and parses
func readEvent(reader io.Reader, event interface{}) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, event)
	if err != nil {
		return err
	}

	return nil
}
