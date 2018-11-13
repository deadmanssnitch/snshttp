# Amazon SNS HTTP Handler

The snshttp package provides an http.Handler adapter for receiving messages
from Amazon SNS over HTTP(s) webhooks. The goal is to reduce the boilerplate
necessary to get a new service or endpoint receiving messages from Amazon SNS.

## Usage

```go
type EventHandler struct {
  snshttp.DefaultHandler
}

func (h *EventHandler) Notification(ctx context.Context, event *snshttp.Notification) error {
  // Process the event here
  fmt.Printf("id=%q subject=%q message=%q timestamp=%q\n",
    event.MessageID,
    event.Subject,
    event.Message,
    event.Timestamp,
  )

  return nil
}


http.Handler("/hooks/sns", snshttp.New(&EventHandler{}))
```

## Timeouts

Amazon SNS expects a webhook to return a response within 15 seconds, any longer
and it considers the request failed and it will try again. Because of this, the
context.Context passed to each snshttp.EventHandler receiver has a 15 second
timeout set from when the request is received.
