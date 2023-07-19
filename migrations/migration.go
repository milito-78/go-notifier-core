package migrations

import (
	"gorm.io/gorm"
	"time"
)

type migrator interface {
	Up() error
	Down() error
}

type ModelGorm struct {
	CreatedAt time.Time `gorm:"not null;type:timestamp;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"not null;type:timestamp;default:current_timestamp ON update current_timestamp"`
	ID        uint64    `gorm:"primarykey"`
}

type ModelGormSoftDelete struct {
	DeletedAt *time.Time `gorm:"type:timestamp;"`
}

type notifierTag struct {
	ModelGorm
	EmailSubscribers        []notifierEmailSubscriber        `gorm:"many2many:notifier_email_sub_tags;ForeignKey:id;References:id;JoinForeignKey:TagId;joinReferences:EmailSubscriberId"`
	MobileSubscribers       []notifierMobileSubscriber       `gorm:"many2many:notifier_mobile_sub_tags;"`
	NotificationSubscribers []notifierNotificationSubscriber `gorm:"many2many:notifier_notification_sub_tags;"`
	EmailCampaigns          []notifierEmailCampaign          `gorm:"many2many:notifier_email_campaign_tags;ForeignKey:id;References:id;JoinForeignKey:TagId;joinReferences:CampaignId"`
	Name                    string                           `gorm:"size:255;index:idx_name,unique;not null"`
}

type createTag struct {
	mg gorm.Migrator
}

func (c createTag) Up() error {
	if !c.mg.HasTable(&notifierTag{}) {
		return c.mg.AutoMigrate(&notifierTag{})
	}
	return nil
}

func (c createTag) Down() error {
	if c.mg.HasTable(&notifierTag{}) {
		_ = c.mg.DropTable("notifier_email_sub_tags")
		return c.mg.DropTable(&notifierTag{})
	}
	return nil
}

type notifierEmailUnsubscribeEvent struct {
	ID     uint64 `gorm:"primarykey"`
	Reason string `gorm:"size:255;not null"`
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
	ModelGorm
	UnsubscribedEvent   *notifierEmailUnsubscribeEvent `gorm:"foreignKey:UnsubscribedEventId;"`
	UnsubscribedEventId *uint64
	UnsubscribedAt      *time.Time    `gorm:"type:timestamp;"`
	Tags                []notifierTag `gorm:"many2many:notifier_email_sub_tags;ForeignKey:id;References:id;JoinForeignKey:EmailSubscriberId;joinReferences:TagId"` //
	FirstName           string        `gorm:"size:255;not null"`
	LastName            string        `gorm:"size:255;not null"`
	Email               string        `gorm:"size:255;index:idx_email;not null"`
}

type notifierEmailSubTag struct {
	EmailSubscriberId uint64 `gorm:"primarykey"`
	TagId             uint64 `gorm:"primarykey"`
}

type createEmailSubscriber struct {
	mg gorm.Migrator
}

func (c createEmailSubscriber) Up() error {
	if !c.mg.HasTable(&notifierEmailSubscriber{}) {
		return c.mg.AutoMigrate(&notifierEmailSubscriber{})
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
	Payload string `gorm:"type=json"`
	Type    string `gorm:"size:255;index:idx_type;not null"`
	Name    string `gorm:"size:255;not null"`
	ID      uint64 `gorm:"primarykey"`
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
	ModelGorm
	RecipientEmail string                  `gorm:"not null;size:255;"`
	EmailService   notifierEmailService    `gorm:"foreignKey:EmailServiceId"`
	EmailServiceId uint64                  `gorm:"not null;"`
	Subscriber     notifierEmailSubscriber `gorm:"foreignKey:SubscriberId"`
	SubscriberId   uint64                  `gorm:"not null;"`
	SourceType     string                  `gorm:"not null;size:255;"`
	FromEmail      string                  `gorm:"not null;size:255;"`
	SourceId       *uint64
	FromName       string     `gorm:"not null;size:255;"`
	QueuedAt       *time.Time `gorm:"type:timestamp"`
	FailedAt       *time.Time `gorm:"type:timestamp"`
	Message        string     `gorm:"not null;"`
	Subject        string     `gorm:"not null;size:255;"`
	SentAt         *time.Time `gorm:"type:timestamp"`
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
	ModelGorm
	Content string `gorm:"not null;type=longtext"`
	Name    string `gorm:"not null;size:255;"`
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
	Name string `gorm:"not null;size:255;"`
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
	ModelGorm
	Service        notifierEmailService          `gorm:"foreignKey:EmailServiceId"`
	EmailServiceId uint64                        `gorm:"not null"`
	ScheduledAt    *time.Time                    `gorm:"type:timestamp"`
	Template       notifierEmailCampaignTemplate `gorm:"foreignKey:TemplateId"`
	TemplateId     uint64                        `gorm:"not null"`
	Status         notifierEmailCampaignStatus   `gorm:"foreignKey:StatusId"`
	Tags           []notifierTag                 `gorm:"many2many:notifier_email_campaign_tags;ForeignKey:id;References:id;JoinForeignKey:CampaignId;joinReferences:TagId"`
	StatusId       uint64                        `gorm:"not null"`
	FromEmail      string                        `gorm:"not null;size:255;"`
	FromName       string                        `gorm:"not null;size:255;"`
	Subject        string                        `gorm:"not null;size:255;"`
	Content        string                        `gorm:"not null;type=longtext"`
	Name           string                        `gorm:"not null;size:255;"`
}

type createEmailCampaign struct {
	mg gorm.Migrator
}

func (c createEmailCampaign) Up() error {
	if !c.mg.HasTable(&notifierEmailCampaign{}) {
		return c.mg.AutoMigrate(&notifierEmailCampaign{})
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
	Reason string `gorm:"not null;size:255;"`
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
	ModelGorm
	UnsubscribedEvent   *notifierMobileUnsubscribeEvent `gorm:"foreignKey:UnsubscribedEventId;"`
	UnsubscribedEventId *uint64
	UnsubscribedAt      *time.Time    `gorm:"type:timestamp"`
	Tags                []notifierTag `gorm:"many2many:notifier_mobile_sub_tags;"`
	FirstName           string        `gorm:"not null;size:255;"`
	LastName            string        `gorm:"not null;size:255;"`
	CountryCode         string        `gorm:"not null;size:100;index:country_index"`
	Mobile              string        `gorm:"not null;size:100;index:mobile_index"`
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
	Name string `gorm:"not null;size:255;"`
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
	ModelGorm
	Tags      []notifierTag              `gorm:"many2many:notifier_notification_sub_tags;"`
	FirstName string                     `gorm:"not null;size:255;"`
	LastName  string                     `gorm:"not null;size:255;"`
	Driver    notifierNotificationDriver `gorm:"foreignKey:DriverId;"`
	DriverId  uint64                     `gorm:"not null;"`
	Token     string                     `gorm:"not null;size:144;index:mobile_index"`
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
