package fastlane

import (
	"context"
	"fmt"
	"os"
)

type PushOptions struct {
	Sound    string
	HTML     bool
	URL      string
	URLTitle string
}

type pusher interface {
	Push(ctx context.Context, message string, opts PushOptions) error
}

func PushSystemNotifications(ctx context.Context, p pusher, events <-chan Event) {
	for e := range events {
		if e.Name != SystemNotificationEventName {
			continue
		}
		sn := e.Data.(SystemNotification)
		msg := fmt.Sprintf("%s: %s", sn.Title, sn.Message)
		err := p.Push(ctx, msg, PushOptions{
			Sound: os.Getenv("PUSHOVER_SOUND"),
		})
		if err != nil {
			fmt.Printf("cannot push notification for event %+v: %v", sn, err)
		}
	}
}
