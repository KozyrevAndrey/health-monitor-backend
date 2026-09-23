package domain

import (
	"errors"
	"time"
)

// SendStatus is the terminal state of a single notification attempt.
type SendStatus string

const (
	// SendStatusSent means the provider accepted the notification.
	SendStatusSent SendStatus = "sent"
	// SendStatusFailed means the provider explicitly rejected it. Nothing was delivered.
	SendStatusFailed SendStatus = "failed"
	// SendStatusUnknown means we never learned the outcome (timeout, unparseable
	// response, crash). It is never retried: the notification may have gone out.
	SendStatusUnknown SendStatus = "unknown"
)

// ErrUnknownResult marks a delivery whose outcome could not be determined.
// Notifiers return it (wrapped) instead of a plain error when a retry could
// duplicate a real-world side effect such as a phone call.
var ErrUnknownResult = errors.New("notification result unknown")

// ErrSendAlreadyClaimed is returned when a notification for the same incident,
// notifier and sequence was already claimed by someone else.
var ErrSendAlreadyClaimed = errors.New("notification send already claimed")

// Receipt carries provider-side identifiers for a delivered notification, so a
// send can be traced back (or its status polled later).
type Receipt struct {
	ExternalID     string
	ExternalStatus string
	Cost           string
}

// NotificationSend is one notification attempt for one incident. The row is
// written before the provider is called, which is what makes "call once" hold
// across restarts and concurrent ticks.
type NotificationSend struct {
	ID             int64      `json:"id"`
	IncidentID     int64      `json:"incident_id"`
	NotifierID     string     `json:"notifier_id"`
	Sequence       int        `json:"sequence"`
	Status         SendStatus `json:"status"`
	ExternalID     string     `json:"external_id,omitempty"`
	ExternalStatus string     `json:"external_status,omitempty"`
	Cost           string     `json:"cost,omitempty"`
	Error          string     `json:"error,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}
