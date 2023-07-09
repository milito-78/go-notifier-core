package go_notifier_core

import "time"

//Campaign Models

type NotifierEmailCampaignTemplate struct {
	UpdatedAt time.Time
	CreatedAt time.Time
	Content   string `gorm:"type=longtext"`
	Name      string
	ID        uint64
}

func NewNotifierEmailCampaignTemplate(content string, name string) *NotifierEmailCampaignTemplate {
	return &NotifierEmailCampaignTemplate{
		Content:   content,
		Name:      name,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
}

type NotifierEmailService struct {
	Payload string `gorm:"type=longtext"`
	Type    string
	Name    string
	ID      uint64
}

type NotifierEmailStatus struct {
	Name string
	ID   uint64
}

type NotifierEmailCampaign struct {
	EmailServiceId uint64
	ScheduledAt    *time.Time
	TemplateId     uint64
	UpdatedAt      time.Time
	CreatedAt      time.Time
	StatusId       uint64
	FromEmail      string
	FromName       string
	Subject        string
	Content        string `gorm:"type=longtext"`
	Name           string
	ID             uint64
}

func NewNotifierEmailCampaign(emailServiceId uint64, scheduledAt *time.Time, templateId uint64, statusId uint64, fromEmail string, fromName string, subject string, content string, name string) *NotifierEmailCampaign {
	return &NotifierEmailCampaign{
		EmailServiceId: emailServiceId,
		ScheduledAt:    scheduledAt,
		TemplateId:     templateId,
		StatusId:       statusId,
		FromEmail:      fromEmail,
		FromName:       fromName,
		Subject:        subject,
		Content:        content,
		Name:           name,
		UpdatedAt:      time.Now(),
		CreatedAt:      time.Now(),
	}
}

type NotifierEmailCampaignTag struct {
	CampaignId uint64
	TagId      uint64
	ID         uint64
}

func NewNotifierEmailCampaignTag(campaignId uint64, tagId uint64) *NotifierEmailCampaignTag {
	return &NotifierEmailCampaignTag{CampaignId: campaignId, TagId: tagId}
}

// Email Subscriber models

type NotifierEmailUnsubscribeEvent struct {
	ID     uint64
	Reason string
}

func NewNotifierEmailUnsubscribeEvent(reason string) *NotifierEmailUnsubscribeEvent {
	return &NotifierEmailUnsubscribeEvent{Reason: reason}
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

//Mobile Subscriber models

type NotifierMobileUnsubscribeEvent struct {
	ID     uint64
	Reason string
}

func NewNotifierMobileUnsubscribeEvent(reason string) *NotifierMobileUnsubscribeEvent {
	return &NotifierMobileUnsubscribeEvent{Reason: reason}
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

// Notification subscriber models

type NotifierNotificationDriver struct {
	Name string
	ID   uint64
}

type NotifierNotificationSubscriber struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	FirstName string
	LastName  string
	DriverId  uint64
	Token     string
	ID        uint64
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

//Tag models

type NotifierTag struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	ID        uint64
}

func NewNotifierTag(name string) *NotifierTag {
	return &NotifierTag{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
