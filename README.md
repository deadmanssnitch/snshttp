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

## Double Requests

When using authentication Amazon SNS will make an initial request without
authentication information to determine which scheme (Basic or Digest) the
endpoint is using. Amazon will then make a request using the correct
authentication scheme. These double requests will happen for every webhook from
SNS but only one will be received by the EventHandler.

## Timeouts

Amazon SNS expects a webhook to return a response within 15 seconds, any longer
and it considers the request failed and it will try again. Because of this, the
context.Context passed to each snshttp.EventHandler receiver has a 15 second
timeout set from when the request is received.

## Thanks

Continued development is sponsored by [Dead Man's Snitch](https://deadmanssnitch.com).

Ever been surprised that a critical scheduled task was silently failing to
run? Whether it's sending invoices, cache clearing, or backups; Dead Man's
Snitch makes it easy to [monitor cron jobs](https://deadmanssnitch.com/docs/cron-job-monitoring)
and [Amazon SNS](https://deadmanssnitch.com/docs/amazon-sns) subscriptions.

Get started with [Dead Man's Snitch](https://deadmanssnitch.com/plans) for free
