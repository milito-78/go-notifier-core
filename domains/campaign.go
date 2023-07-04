package domains

import "time"

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
