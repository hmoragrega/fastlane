package fastlane

import (
	"github.com/google/uuid"
)

const (
	ReviewsEventName            = "REVIEWS"
	ReviewMergedEventName       = "REVIEW-MERGED"
	ReviewsMergedEventName      = "REVIEWS-MERGED"
	MergeEventName              = "MERGE"
	NotificationEventName       = "NOTIFICATION"
	SystemNotificationEventName = "SYSTEM-NOTIFICATION"
)

const (
	SuccessType NotificationType = "success"
	InfoType    NotificationType = "info"
	WarningType NotificationType = "warning"
	ErrorType   NotificationType = "error"
)

type Event struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type NotificationType string

type Notification struct {
	ID      uuid.UUID        `json:"id"`
	Message string           `json:"message"`
	Type    NotificationType `json:"type"`
}

type SystemNotification struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ReviewMerged struct {
	Review      Review   `json:"review"`
	Pipeline    Pipeline `json:"pipeline"`
	HasPipeline bool     `json:"has_pipeline"`
}

type ReviewUpdatedData struct {
	Old Review
	New Review
}

func calculateReviewEvents(current, updated []Review) (events []Event) {
	m := make(map[string]Review, len(current))
	for _, r := range current {
		m[r.ID] = r
	}
	for _, r := range updated {
		old, ok := m[r.ID]
		if !ok {
			continue
		}
		if r.MergeEnabled() && !old.MergeEnabled() {
			events = append(events, Event{
				Name: SystemNotificationEventName,
				Data: SystemNotification{
					Title:   r.Title,
					Message: "Review can be merged! click to merge",
				}})
		}
	}

	events = append(events, Event{
		Name: SystemNotificationEventName,
		Data: SystemNotification{
			Title:   "EL TEST",
			Message: "Review can be merged! click to merge",
		}})

	return events
}
