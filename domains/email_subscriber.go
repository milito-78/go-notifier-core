package domains

import "time"

type NotifierEmailUnsubscribeEvent struct {
	ID     uint64
	Type   uint8
	Reason string
}

func NewNotifierEmailUnsubscribeEvent(Type uint8, reason string) *NotifierEmailUnsubscribeEvent {
	return &NotifierEmailUnsubscribeEvent{Type: Type, Reason: reason}
}

type NotifierEmailSubscriber struct {
	UnsubscribedEventId *uint64
	UnsubscribedAt      *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
	FirstName           string
	LastName            string
	Email               string
	ID                  uint64
}

func NewNotifierEmailSubscriber(email, firstName, lastName string) *NotifierEmailSubscriber {
	return &NotifierEmailSubscriber{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type NotifierEmailSubTag struct {
	EmailSubscriberId uint64
	CreatedAt         time.Time
	UpdatedAt         time.Time
	TagId             uint64
	ID                uint64
}

func NewNotifierEmailSubTag(emailSubscriberId uint64, tagId uint64) *NotifierEmailSubTag {
	return &NotifierEmailSubTag{
		EmailSubscriberId: emailSubscriberId,
		TagId:             tagId,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

type NotifierEmailMessage struct {
	RecipientEmail string
	EmailServiceId uint64
	SubscriberId   uint64
	SourceType     string
	Message        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	FromEmail      string
	SourceId       uint64
	FromName       string
	QueuedAt       *time.Time
	FailedAt       *time.Time
	Subject        string
	SentAt         *time.Time
	ID             uint64
}

func NewNotifierEmailMessage(recipientEmail string, subscriberId uint64, sourceType string, fromEmail string, sourceId uint64, fromName string, subject string, emailServiceId uint64, message string) *NotifierEmailMessage {
	queuedAt := time.Now()
	return &NotifierEmailMessage{
		EmailServiceId: emailServiceId,
		RecipientEmail: recipientEmail,
		SubscriberId:   subscriberId,
		SourceType:     sourceType,
		FromEmail:      fromEmail,
		SourceId:       sourceId,
		FromName:       fromName,
		QueuedAt:       &queuedAt,
		Subject:        subject,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Message:        message,
	}
}
