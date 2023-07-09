package migrations

import (
	"gorm.io/gorm"
	"time"
)

type migrator interface {
	Up() error
	Down() error
}

type notifierTag struct {
	gorm.Model
	EmailSubscribers        []notifierEmailSubscriber        `gorm:"many2many:notifier_email_sub_tags;"`
	MobileSubscribers       []notifierMobileSubscriber       `gorm:"many2many:notifier_mobile_sub_tags;"`
	NotificationSubscribers []notifierNotificationSubscriber `gorm:"many2many:notifier_notification_sub_tags;"`
	EmailCampaigns          []notifierEmailCampaign          `gorm:"many2many:notifier_email_campaign_tags;"`
	Name                    string                           `gorm:"size:255;index:idx_name,unique;"`
	ID                      uint64                           `gorm:"primarykey"`
}

type createTag struct {
	mg gorm.Migrator
}

func (c createTag) Up() error {
	if !c.mg.HasTable(&notifierTag{}) {
		return c.mg.CreateTable(&notifierTag{})
	}
	return nil
}

func (c createTag) Down() error {
	if c.mg.HasTable(&notifierTag{}) {
		return c.mg.DropTable(&notifierTag{})
	}
	return nil
}

type notifierEmailUnsubscribeEvent struct {
	ID     uint64 `gorm:"primarykey"`
	Reason string `gorm:"size:255;"`
}

type createEmailUnsubscribeEvent struct {
	mg gorm.Migrator
}

func (c createEmailUnsubscribeEvent) Up() error {
	if !c.mg.HasTable(&notifierEmailUnsubscribeEvent{}) {
		return c.mg.CreateTable(&notifierEmailUnsubscribeEvent{})
	}
	return nil
}

func (c createEmailUnsubscribeEvent) Down() error {
	if c.mg.HasTable(&notifierEmailUnsubscribeEvent{}) {
		return c.mg.DropTable(&notifierEmailUnsubscribeEvent{})
	}
	return nil
}

type notifierEmailSubscriber struct {
	UnsubscribedEvent   *notifierEmailUnsubscribeEvent `gorm:"foreignKey:UnsubscribedEventId;"`
	UnsubscribedEventId *uint64
	UnsubscribedAt      *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Tags                []notifierTag `gorm:"many2many:notifier_email_sub_tags;"`
	FirstName           string
	LastName            string
	Email               string `gorm:"index:email_index"`
	ID                  uint64 `gorm:"primarykey"`
}

type createEmailSubscriber struct {
	mg gorm.Migrator
}

func (c createEmailSubscriber) Up() error {
	if !c.mg.HasTable(&notifierEmailSubscriber{}) {
		return c.mg.CreateTable(&notifierEmailSubscriber{})
	}
	return nil
}

func (c createEmailSubscriber) Down() error {
	if c.mg.HasTable(&notifierEmailSubscriber{}) {
		return c.mg.DropTable(&notifierEmailSubscriber{})
	}
	return nil
}

type notifierEmailService struct {
	Payload string `gorm:"type=longtext"`
	Type    string
	Name    string
	ID      uint64
}

type createEmailService struct {
	mg gorm.Migrator
}

func (c createEmailService) Up() error {
	if !c.mg.HasTable(&notifierEmailService{}) {
		return c.mg.CreateTable(&notifierEmailService{})
	}
	return nil
}

func (c createEmailService) Down() error {
	if c.mg.HasTable(&notifierEmailService{}) {
		return c.mg.DropTable(&notifierEmailService{})
	}
	return nil
}

type notifierEmailMessage struct {
	RecipientEmail string
	EmailService   notifierEmailService `gorm:"foreignKey:EmailServiceId"`
	EmailServiceId uint64
	Subscriber     notifierEmailSubscriber `gorm:"foreignKey:SubscriberId"`
	SubscriberId   uint64
	SourceType     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	FromEmail      string
	SourceId       *uint64
	FromName       string
	QueuedAt       *time.Time
	FailedAt       *time.Time
	Message        string
	Subject        string
	SentAt         *time.Time
	ID             uint64 `gorm:"primarykey"`
}

type createEmailMessage struct {
	mg gorm.Migrator
}

func (c createEmailMessage) Up() error {
	if !c.mg.HasTable(&notifierEmailMessage{}) {
		return c.mg.CreateTable(&notifierEmailMessage{})
	}
	return nil
}

func (c createEmailMessage) Down() error {
	if c.mg.HasTable(&notifierEmailMessage{}) {
		return c.mg.DropTable(&notifierEmailMessage{})
	}
	return nil
}

type notifierEmailCampaignTemplate struct {
	UpdatedAt time.Time
	CreatedAt time.Time
	Content   string `gorm:"type=longtext"`
	Name      string
	ID        uint64 `gorm:"primarykey"`
}

type createEmailCampaignTemplate struct {
	mg gorm.Migrator
}

func (c createEmailCampaignTemplate) Up() error {
	if !c.mg.HasTable(&notifierEmailCampaignTemplate{}) {
		return c.mg.CreateTable(&notifierEmailCampaignTemplate{})
	}
	return nil
}

func (c createEmailCampaignTemplate) Down() error {
	if c.mg.HasTable(&notifierEmailCampaignTemplate{}) {
		return c.mg.DropTable(&notifierEmailCampaignTemplate{})
	}
	return nil
}

type notifierEmailCampaignStatus struct {
	ID   uint64 `gorm:"primarykey"`
	Name string `gorm:"size:255;"`
}

type createEmailCampaignStatus struct {
	mg gorm.Migrator
}

func (c createEmailCampaignStatus) Up() error {
	if !c.mg.HasTable(&notifierEmailCampaignStatus{}) {
		return c.mg.CreateTable(&notifierEmailCampaignStatus{})
	}
	return nil
}

func (c createEmailCampaignStatus) Down() error {
	if c.mg.HasTable(&notifierEmailCampaignStatus{}) {
		return c.mg.DropTable(&notifierEmailCampaignStatus{})
	}
	return nil
}

type notifierEmailCampaign struct {
	EmailServiceId uint64
	ScheduledAt    *time.Time
	Template       notifierEmailCampaignTemplate `gorm:"foreignKey:TemplateId"`
	TemplateId     uint64
	UpdatedAt      time.Time
	CreatedAt      time.Time
	Status         notifierEmailCampaignStatus `gorm:"foreignKey:TemplateId"`
	Tags           []notifierTag               `gorm:"many2many:notifier_email_campaign_tags;"`
	StatusId       uint64
	FromEmail      string
	FromName       string
	Subject        string
	Content        string `gorm:"type=longtext"`
	Name           string
	ID             uint64
}

type createEmailCampaign struct {
	mg gorm.Migrator
}

func (c createEmailCampaign) Up() error {
	if !c.mg.HasTable(&notifierEmailCampaign{}) {
		return c.mg.CreateTable(&notifierEmailCampaign{})
	}
	return nil
}

func (c createEmailCampaign) Down() error {
	if c.mg.HasTable(&notifierEmailCampaign{}) {
		return c.mg.DropTable(&notifierEmailCampaign{})
	}
	return nil
}

//Mobile notification

type notifierMobileUnsubscribeEvent struct {
	ID     uint64 `gorm:"primarykey"`
	Reason string `gorm:"size:255;"`
}

type createMobileUnsubscribeEvent struct {
	mg gorm.Migrator
}

func (c createMobileUnsubscribeEvent) Up() error {
	if !c.mg.HasTable(&notifierMobileUnsubscribeEvent{}) {
		return c.mg.CreateTable(&notifierMobileUnsubscribeEvent{})
	}
	return nil
}

func (c createMobileUnsubscribeEvent) Down() error {
	if c.mg.HasTable(&notifierMobileUnsubscribeEvent{}) {
		return c.mg.DropTable(&notifierMobileUnsubscribeEvent{})
	}
	return nil
}

type notifierMobileSubscriber struct {
	UnsubscribedEvent   *notifierMobileUnsubscribeEvent `gorm:"foreignKey:UnsubscribedEventId;"`
	UnsubscribedEventId *uint64
	UnsubscribedAt      *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Tags                []notifierTag `gorm:"many2many:notifier_mobile_sub_tags;"`
	FirstName           string
	LastName            string
	CountryCode         string `gorm:"index:country_index"`
	Mobile              string `gorm:"index:mobile_index"`
	ID                  uint64 `gorm:"primarykey"`
}

type createMobileSubscriber struct {
	mg gorm.Migrator
}

func (c createMobileSubscriber) Up() error {
	if !c.mg.HasTable(&notifierMobileSubscriber{}) {
		return c.mg.CreateTable(&notifierMobileSubscriber{})
	}
	return nil
}

func (c createMobileSubscriber) Down() error {
	if c.mg.HasTable(&notifierMobileSubscriber{}) {
		return c.mg.DropTable(&notifierMobileSubscriber{})
	}
	return nil
}

//Notification

type notifierNotificationDriver struct {
	ID   uint64 `gorm:"primarykey"`
	Name string `gorm:"size:255;"`
}

type createNotificationDriver struct {
	mg gorm.Migrator
}

func (c createNotificationDriver) Up() error {
	if !c.mg.HasTable(&notifierNotificationDriver{}) {
		return c.mg.CreateTable(&notifierNotificationDriver{})
	}
	return nil
}

func (c createNotificationDriver) Down() error {
	if c.mg.HasTable(&notifierNotificationDriver{}) {
		return c.mg.DropTable(&notifierNotificationDriver{})
	}
	return nil
}

type notifierNotificationSubscriber struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Tags      []notifierTag `gorm:"many2many:notifier_notification_sub_tags;"`
	FirstName string
	LastName  string
	Driver    notifierNotificationDriver `gorm:"foreignKey:DriverId;"`
	DriverId  uint64
	Token     string `gorm:"index:mobile_index"`
	ID        uint64 `gorm:"primarykey"`
}

type createNotificationSubscriber struct {
	mg gorm.Migrator
}

func (c createNotificationSubscriber) Up() error {
	if !c.mg.HasTable(&notifierNotificationSubscriber{}) {
		return c.mg.CreateTable(&notifierNotificationSubscriber{})
	}
	return nil
}

func (c createNotificationSubscriber) Down() error {
	if c.mg.HasTable(&notifierNotificationSubscriber{}) {
		return c.mg.DropTable(&notifierNotificationSubscriber{})
	}
	return nil
}

func GetMigrationsList(migr gorm.Migrator) []migrator {
	var migrations = make([]migrator, 12)
	migrations[0] = createTag{migr}
	migrations[1] = createEmailUnsubscribeEvent{migr}
	migrations[2] = createEmailSubscriber{migr}
	migrations[3] = createEmailService{migr}
	migrations[4] = createEmailMessage{migr}
	migrations[5] = createEmailCampaignTemplate{migr}
	migrations[6] = createEmailCampaignStatus{migr}
	migrations[7] = createEmailCampaign{migr}
	migrations[8] = createMobileUnsubscribeEvent{migr}
	migrations[9] = createMobileSubscriber{migr}
	migrations[10] = createNotificationDriver{migr}
	migrations[11] = createNotificationSubscriber{migr}

	return migrations
}
