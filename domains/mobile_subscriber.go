package domains

import "time"

type NotifierMobileUnsubscribeEvent struct {
	ID     uint64
	Type   uint8
	Reason string
}

func NewNotifierMobileUnsubscribeEvent(Type uint8, reason string) *NotifierMobileUnsubscribeEvent {
	return &NotifierMobileUnsubscribeEvent{Type: Type, Reason: reason}
}

type NotifierMobileSubscriber struct {
	UnsubscribedEventId *uint64
	UnsubscribedAt      *time.Time
	CountryCode         string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	FirstName           string
	LastName            string
	Mobile              string
	ID                  uint64
}

func NewNotifierMobileSubscriber(countryCode, mobile, firstName, lastName string) *NotifierMobileSubscriber {
	return &NotifierMobileSubscriber{
		FirstName:   firstName,
		LastName:    lastName,
		CountryCode: countryCode,
		Mobile:      mobile,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

type NotifierMobileSubTag struct {
	MobileSubscriberId uint64
	CreatedAt          time.Time
	UpdatedAt          time.Time
	TagId              uint64
	ID                 uint64
}

func NewNotifierMobileSubTag(mobileSubscriberId uint64, tagId uint64) *NotifierMobileSubTag {
	return &NotifierMobileSubTag{
		MobileSubscriberId: mobileSubscriberId,
		TagId:              tagId,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}
