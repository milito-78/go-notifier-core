package domains

import "time"

type NotifierNotificationDriver struct {
	Name string
	ID   uint64
}

type NotifierNotificationSubscriber struct {
	UnsubscribedEventId *uint64
	UnsubscribedAt      *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
	FirstName           string
	LastName            string
	DriverId            uint64
	Token               string
	ID                  uint64
}

func NewNotifierNotificationSubscriber(token, firstName, lastName string, driverId uint64) *NotifierNotificationSubscriber {
	return &NotifierNotificationSubscriber{
		FirstName: firstName,
		LastName:  lastName,
		DriverId:  driverId,
		Token:     token,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type NotifierNotificationSubTag struct {
	NotificationSubscriberId uint64
	CreatedAt                time.Time
	UpdatedAt                time.Time
	TagId                    uint64
	ID                       uint64
}

func NewNotifierNotificationSubTag(notificationSubscriberId uint64, tagId uint64) *NotifierNotificationSubTag {
	return &NotifierNotificationSubTag{
		NotificationSubscriberId: notificationSubscriberId,
		TagId:                    tagId,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}
}
