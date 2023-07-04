package go_notifier_core

import (
	"errors"
	"github.com/golobby/container/v3"
	"go-notifier-core/domains"
	"go-notifier-core/repositories"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

func dbFactory(config DbConfig) *gorm.DB {
	switch config.Driver {
	case MysqlDriver:
		return mysqlDriverDb(config)
	}
	return mysqlDriverDb(config)
}

func mysqlDriverDb(config DbConfig) *gorm.DB {
	if config.Password != "" {
		config.Password = ":" + config.Password
	}
	dsn := config.Username + config.Password +
		"@tcp(" + config.Host + ":" + config.Port + ")/" +
		config.DB + "?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error during connecting db mysql driver : %s", err)
	}

	return db
}

func initRepositories() {
	_ = container.Singleton(func(db *gorm.DB) repositories.ITagRepository {
		return repositories.NewGormTagRepository(db)
	})

	//Email repositories #start
	_ = container.Singleton(func(db *gorm.DB) repositories.IEmailUnSubEventRepository {
		return repositories.NewGormEmailUnSubEventRepository(db)
	})

	_ = container.Singleton(func(db *gorm.DB) repositories.IEmailSubTagRepository {
		return repositories.NewGormEmailSubTagRepository(db)
	})

	_ = container.Singleton(func(db *gorm.DB) repositories.IEmailSubscriberRepository {
		return repositories.NewGormEmailSubscriberRepository(db)
	})
	//Email repositories #end

	//Mobile repositories #start
	_ = container.Singleton(func(db *gorm.DB) repositories.IMobileUnSubEventRepository {
		return repositories.NewGormMobileUnSubEventRepository(db)
	})

	_ = container.Singleton(func(db *gorm.DB) repositories.IMobileSubTagRepository {
		return repositories.NewGormMobileSubTagRepository(db)
	})

	_ = container.Singleton(func(db *gorm.DB) repositories.IMobileSubscriberRepository {
		return repositories.NewGormMobileSubscriberRepository(db)
	})
	//Mobile repositories #end

	//Notification repositories #start
	_ = container.Singleton(func(db *gorm.DB) repositories.INotifierNotificationDriverRepository {
		return repositories.NewGormNotifierNotificationDriverRepository(db)
	})

	_ = container.Singleton(func(db *gorm.DB) repositories.INotificationSubTagRepository {
		return repositories.NewGormNotificationSubTagRepository(db)
	})

	_ = container.Singleton(func(db *gorm.DB) repositories.INotificationSubscriberRepository {
		return repositories.NewGormNotificationSubscriberRepository(db)
	})
	//Notification repositories #end

	//Campaign repositories #start
	_ = container.Singleton(func(db *gorm.DB) repositories.IEmailTemplateRepository {
		return repositories.NewGormEmailTemplateRepository(db)
	})

	_ = container.Singleton(func(db *gorm.DB) repositories.IEmailServiceRepository {
		return repositories.NewGormEmailServiceRepository(db)
	})

	_ = container.Singleton(func(db *gorm.DB) repositories.IEmailStatusRepository {
		return repositories.NewGormEmailStatusRepository(db)
	})

	_ = container.Singleton(func(db *gorm.DB) repositories.IEmailCampaignRepository {
		return repositories.NewGormEmailCampaignRepository(db)
	})
	//Campaign repositories #end
}

func initMailers() {
	_ = container.NamedSingleton("smtp", func() Mailer {
		return new(SmtpMailer)
	})
}

func Initialize(config DbConfig) {
	err := container.NamedSingleton("db", func() *gorm.DB {
		return dbFactory(config)
	})
	if err != nil {
		log.Fatalf("Error during generate singleton : %s", err)
	}

	initRepositories()
	initMailers()
}

// Tag functions #start

// CreateTag
/**
for creating tag use this function. If tag exists you will get error.
*Hint* Tags stored as `lower-cases` in db
*/
func CreateTag(name string) (*domains.NotifierTag, error) {
	var tgRepo repositories.ITagRepository
	err := container.Resolve(&tgRepo)
	if err != nil {
		return nil, err
	}
	nLower := strings.ToLower(name)

	res, err := tgRepo.GetByName(nLower)
	if res.ID != 0 && err == nil {
		return res, errors.New("tag is exists")
	}

	tmp := domains.NewNotifierTag(nLower)
	err = tgRepo.Create(tmp)
	if err != nil {
		return nil, err
	}
	return tmp, nil
}

// DeleteTagByName
// for deleting tag by its name use this function.
func DeleteTagByName(name string) error {
	var tgRepo repositories.ITagRepository
	err := container.Resolve(&tgRepo)
	if err != nil {
		return err
	}

	tmp := &domains.NotifierTag{
		Name: strings.ToLower(name),
	}
	err = tgRepo.Delete(tmp)
	if err != nil {
		return err
	}
	return nil
}

// GetTagByName
// for getting tag by its name use this function.
func GetTagByName(name string) (*domains.NotifierTag, error) {
	var tgRepo repositories.ITagRepository
	err := container.Resolve(&tgRepo)
	if err != nil {
		return nil, err
	}

	nLower := strings.ToLower(name)
	res, err := tgRepo.GetByName(nLower)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// TagsList
// for getting tags list use this function.
func TagsList() ([]domains.NotifierTag, error) {
	var tgRepo repositories.ITagRepository
	err := container.Resolve(&tgRepo)
	if err != nil {
		return nil, err
	}
	var data []domains.NotifierTag
	tgRepo.All(data)
	return data, nil
}

func fetchTags(tags []string, createTag bool) ([]uint64, error) {
	var tagsEntity []uint64

	if len(tags) == 0 {
		tmp, err := CreateTag("all")
		if err != nil && tmp.ID == 0 {
			return nil, err
		}
		tagsEntity = append(tagsEntity, tmp.ID)
		return tagsEntity, nil
	}

	if createTag {
		for _, tag := range tags {
			tmp, err := CreateTag(tag)
			if err != nil && tmp.ID == 0 {
				return nil, err
			}
			tagsEntity = append(tagsEntity, tmp.ID)
		}
	} else {
		for _, tag := range tags {
			tmp, err := GetTagByName(tag)
			if err != nil {
				return nil, err
			}
			tagsEntity = append(tagsEntity, tmp.ID)
		}
	}
	return tagsEntity, nil
}

func checkAllTagExists(tags []string) bool {
	for _, tag := range tags {
		if strings.ToLower(tag) == "all" {
			return true
		}
	}
	return false
}

// Tag functions #end

// Email subscribe functions #start

func SubscribeEmail(email, fName, lName string, tags []string, createTag bool) (*domains.NotifierEmailSubscriber, error) {
	tagsEntity, err := fetchTags(tags, createTag)
	if err != nil {
		return nil, err
	}

	var subRepo repositories.IEmailSubscriberRepository
	err = container.Resolve(&subRepo)
	if err != nil {
		return nil, err
	}
	subscriber := domains.NewNotifierEmailSubscriber(email, fName, lName)
	err = subRepo.Create(subscriber)
	if err != nil {
		return nil, err
	}
	err = subRepo.AssignTagToUser(subscriber.ID, tagsEntity)
	if err != nil {
		return subscriber, err
	}

	return subscriber, nil
}

func AssignTagsToEmail(email string, tags []string, createTag bool) error {
	if len(tags) == 0 {
		return errors.New("tags is empty")
	}
	tagsEntity, err := fetchTags(tags, createTag)
	if err != nil {
		return err
	}

	var subRepo repositories.IEmailSubscriberRepository
	err = container.Resolve(&subRepo)
	if err != nil {
		return err
	}
	subscriber, err := subRepo.GetByEmail(email)
	if err != nil {
		return err
	}

	err = subRepo.AssignTagToUser(subscriber.ID, tagsEntity)
	if err != nil {
		return err
	}

	return nil
}

func RemoveTagsFromEmail(email string, tags []string) error {
	if len(tags) == 0 {
		return errors.New("tags is empty")
	}
	if checkAllTagExists(tags) {
		return errors.New("you can't remove user from all tag")
	}

	var tagsEntity []uint64
	for _, tag := range tags {
		tmp, err := GetTagByName(tag)
		if err == nil {
			tagsEntity = append(tagsEntity, tmp.ID)
		}
	}

	var subRepo repositories.IEmailSubscriberRepository
	err := container.Resolve(&subRepo)
	if err != nil {
		return err
	}
	subscriber, err := subRepo.GetByEmail(email)
	if err != nil {
		return err
	}

	err = subRepo.RemoveTagsFromUser(subscriber.ID, tagsEntity)
	if err != nil {
		return err
	}

	return nil
}

func UnSubscribeEmail(email string, unsubId uint64) error {
	var subRepo repositories.IEmailSubscriberRepository
	err := container.Resolve(&subRepo)
	if err != nil {
		return err
	}

	subscriber, err := subRepo.GetByEmail(email)
	if err != nil {
		return err
	}

	var eventRepo repositories.IEmailUnSubEventRepository
	err = container.Resolve(&eventRepo)
	if err != nil {
		return err
	}

	_, err = eventRepo.Get(unsubId)
	if err != nil {
		return err
	}

	subscriber.UnsubscribedEventId = &unsubId
	n := time.Now()
	subscriber.UnsubscribedAt = &n

	err = subRepo.Update(subscriber)
	if err != nil {
		return err
	}

	return nil
}

func EmailUnsubscribeEventsList() ([]domains.NotifierEmailUnsubscribeEvent, error) {
	var eventRepo repositories.IEmailUnSubEventRepository
	err := container.Resolve(&eventRepo)
	if err != nil {
		return nil, err
	}
	var data []domains.NotifierEmailUnsubscribeEvent
	eventRepo.All(data)
	return data, nil
}

func GetTagEmailSubscribers(tag string) ([]domains.NotifierEmailSubscriber, error) {
	tmp, err := GetTagByName(tag)
	if err != nil {
		return nil, err
	}

	var subRepo repositories.IEmailSubscriberRepository
	err = container.Resolve(&subRepo)
	if err != nil {
		return nil, err
	}
	var data []domains.NotifierEmailSubscriber
	subRepo.GetSubscribersForTag(tmp.ID, data)
	return data, nil
}

func GetUnsubscribedEmails() ([]domains.NotifierEmailSubscriber, error) {
	var subRepo repositories.IEmailSubscriberRepository
	err := container.Resolve(&subRepo)
	if err != nil {
		return nil, err
	}

	var data []domains.NotifierEmailSubscriber
	subRepo.GetUnSubscribed(data)
	return data, nil
}

// Email subscribe functions #end

// Mobile subscribe functions #start

func SubscribeMobile(countryCode, mobile, fName, lName string, tags []string, createTag bool) (*domains.NotifierMobileSubscriber, error) {
	tagsEntity, err := fetchTags(tags, createTag)
	if err != nil {
		return nil, err
	}

	var subRepo repositories.IMobileSubscriberRepository
	err = container.Resolve(&subRepo)
	if err != nil {
		return nil, err
	}
	subscriber := domains.NewNotifierMobileSubscriber(countryCode, mobile, fName, lName)
	err = subRepo.Create(subscriber)
	if err != nil {
		return nil, err
	}
	err = subRepo.AssignTagToUser(subscriber.ID, tagsEntity)
	if err != nil {
		return subscriber, err
	}

	return subscriber, nil
}

func AssignTagsToMobile(mobile string, tags []string, createTag bool) error {
	if len(tags) == 0 {
		return errors.New("tags is empty")
	}

	tagsEntity, err := fetchTags(tags, createTag)
	if err != nil {
		return err
	}

	var subRepo repositories.IMobileSubscriberRepository
	err = container.Resolve(&subRepo)
	if err != nil {
		return err
	}
	subscriber, err := subRepo.GetByMobile(mobile)
	if err != nil {
		return err
	}

	err = subRepo.AssignTagToUser(subscriber.ID, tagsEntity)
	if err != nil {
		return err
	}

	return nil
}

func RemoveTagsFromMobile(mobile string, tags []string) error {
	if len(tags) == 0 {
		return errors.New("tags is empty")
	}

	if checkAllTagExists(tags) {
		return errors.New("you can't remove user from all tag")
	}
	var tagsEntity []uint64
	for _, tag := range tags {
		tmp, err := GetTagByName(tag)
		if err == nil {
			tagsEntity = append(tagsEntity, tmp.ID)
		}
	}

	var subRepo repositories.IMobileSubscriberRepository
	err := container.Resolve(&subRepo)
	if err != nil {
		return err
	}
	subscriber, err := subRepo.GetByMobile(mobile)
	if err != nil {
		return err
	}

	err = subRepo.RemoveTagsFromUser(subscriber.ID, tagsEntity)
	if err != nil {
		return err
	}

	return nil
}

func UnSubscribeMobile(mobile string, unsubId uint64) error {
	var subRepo repositories.IMobileSubscriberRepository
	err := container.Resolve(&subRepo)
	if err != nil {
		return err
	}

	subscriber, err := subRepo.GetByMobile(mobile)
	if err != nil {
		return err
	}

	var eventRepo repositories.IMobileUnSubEventRepository
	err = container.Resolve(&eventRepo)
	if err != nil {
		return err
	}

	_, err = eventRepo.Get(unsubId)
	if err != nil {
		return err
	}

	subscriber.UnsubscribedEventId = &unsubId
	n := time.Now()
	subscriber.UnsubscribedAt = &n

	err = subRepo.Update(subscriber)
	if err != nil {
		return err
	}

	return nil
}

func MobileUnsubscribeEventsList() ([]domains.NotifierMobileUnsubscribeEvent, error) {
	var eventRepo repositories.IMobileUnSubEventRepository
	err := container.Resolve(&eventRepo)
	if err != nil {
		return nil, err
	}
	var data []domains.NotifierMobileUnsubscribeEvent
	eventRepo.All(data)
	return data, nil
}

func GetTagMobileSubscribers(tag string) ([]domains.NotifierMobileSubscriber, error) {
	tmp, err := GetTagByName(tag)
	if err != nil {
		return nil, err
	}

	var subRepo repositories.IMobileSubscriberRepository
	err = container.Resolve(&subRepo)
	if err != nil {
		return nil, err
	}
	var data []domains.NotifierMobileSubscriber
	subRepo.GetSubscribersForTag(tmp.ID, data)
	return data, nil
}

func GetUnsubscribedMobiles() ([]domains.NotifierMobileSubscriber, error) {
	var subRepo repositories.IMobileSubscriberRepository
	err := container.Resolve(&subRepo)
	if err != nil {
		return nil, err
	}

	var data []domains.NotifierMobileSubscriber
	subRepo.GetUnSubscribed(data)
	return data, nil
}

// Mobile subscribe functions #end

// Notification subscribe functions #start

func AddNewToken(token, fName, lName string, driverId uint64, tags []string, createTag bool) (*domains.NotifierNotificationSubscriber, error) {
	var tagsEntity []uint64
	if len(tags) == 0 {
		tmp, err := CreateTag("all")
		if err != nil && tmp.ID == 0 {
			return nil, err
		}
		tagsEntity = append(tagsEntity, tmp.ID)
		createTag = false
	} else {
		if createTag {
			for _, tag := range tags {
				tmp, err := CreateTag(tag)
				if err != nil && tmp.ID == 0 {
					return nil, err
				}
				tagsEntity = append(tagsEntity, tmp.ID)
			}
		} else {
			for _, tag := range tags {
				tmp, err := GetTagByName(tag)
				if err != nil {
					return nil, err
				}
				tagsEntity = append(tagsEntity, tmp.ID)
			}
		}
	}

	var driverRepo repositories.INotifierNotificationDriverRepository
	err := container.Resolve(&driverRepo)
	if err != nil {
		return nil, err
	}

	_, err = driverRepo.Get(driverId)
	if err != nil {
		return nil, err
	}

	var subRepo repositories.INotificationSubscriberRepository
	err = container.Resolve(&subRepo)
	if err != nil {
		return nil, err
	}
	subscriber := domains.NewNotifierNotificationSubscriber(token, fName, lName, driverId)
	err = subRepo.Create(subscriber)
	if err != nil {
		return nil, err
	}
	err = subRepo.AssignTagToUser(subscriber.ID, tagsEntity)
	if err != nil {
		return subscriber, err
	}

	return subscriber, nil
}

func AssignTagsToToken(token string, tags []string, createTag bool) error {
	if len(tags) == 0 {
		return errors.New("tags is empty")
	}

	var tagsEntity []uint64
	if createTag {
		for _, tag := range tags {
			tmp, err := CreateTag(tag)
			if err != nil && tmp.ID == 0 {
				return err
			}
			tagsEntity = append(tagsEntity, tmp.ID)
		}
	} else {
		for _, tag := range tags {
			tmp, err := GetTagByName(tag)
			if err != nil {
				return err
			}
			tagsEntity = append(tagsEntity, tmp.ID)
		}
	}

	var subRepo repositories.INotificationSubscriberRepository
	err := container.Resolve(&subRepo)
	if err != nil {
		return err
	}
	subscriber, err := subRepo.GetByNotification(token)
	if err != nil {
		return err
	}

	err = subRepo.AssignTagToUser(subscriber.ID, tagsEntity)
	if err != nil {
		return err
	}

	return nil
}

func RemoveTagsFromToken(token string, tags []string) error {
	if len(tags) == 0 {
		return errors.New("tags is empty")
	}

	if checkAllTagExists(tags) {
		return errors.New("you can't remove user from all tag")
	}

	var tagsEntity []uint64
	for _, tag := range tags {
		tmp, err := GetTagByName(tag)
		if err == nil {
			tagsEntity = append(tagsEntity, tmp.ID)
		}
	}

	var subRepo repositories.INotificationSubscriberRepository
	err := container.Resolve(&subRepo)
	if err != nil {
		return err
	}
	subscriber, err := subRepo.GetByNotification(token)
	if err != nil {
		return err
	}

	err = subRepo.RemoveTagsFromUser(subscriber.ID, tagsEntity)
	if err != nil {
		return err
	}

	return nil
}

func RemoveToken(token string) error {
	var subRepo repositories.INotificationSubscriberRepository
	err := container.Resolve(&subRepo)
	if err != nil {
		return err
	}

	subscriber, err := subRepo.GetByNotification(token)
	if err != nil {
		return err
	}

	err = subRepo.Delete(subscriber)
	if err != nil {
		return err
	}
	return nil
}

func GetTagTokenSubscribers(tag string) ([]domains.NotifierNotificationSubscriber, error) {
	tmp, err := GetTagByName(tag)
	if err != nil {
		return nil, err
	}

	var subRepo repositories.INotificationSubscriberRepository
	err = container.Resolve(&subRepo)
	if err != nil {
		return nil, err
	}
	var data []domains.NotifierNotificationSubscriber
	subRepo.GetSubscribersForTag(tmp.ID, data)
	return data, nil
}

func GetTagAndDriverTokenSubscribers(tag string, driverId uint64) ([]domains.NotifierNotificationSubscriber, error) {
	tmp, err := GetTagByName(tag)
	if err != nil {
		return nil, err
	}

	var subRepo repositories.INotificationSubscriberRepository
	err = container.Resolve(&subRepo)
	if err != nil {
		return nil, err
	}
	var data []domains.NotifierNotificationSubscriber
	subRepo.GetSubscribersForTagAndDriver(tmp.ID, driverId, data)
	return data, nil
}

func NotificationDriversList() ([]domains.NotifierNotificationDriver, error) {
	var tgRepo repositories.INotifierNotificationDriverRepository
	err := container.Resolve(&tgRepo)
	if err != nil {
		return nil, err
	}
	var data []domains.NotifierNotificationDriver
	tgRepo.All(data)
	return data, nil
}

// Notification subscribe functions #end

func CreateEmailTemplate(name, content string) (*domains.NotifierEmailCampaignTemplate, error) {
	var tmRepo repositories.IEmailTemplateRepository
	err := container.Resolve(&tmRepo)
	if err != nil {
		return nil, err
	}
	tmp := domains.NewNotifierEmailCampaignTemplate(content, name)
	err = tmRepo.Create(tmp)
	if err != nil {
		return nil, err
	}
	return tmp, nil
}

func UpdateEmailTemplate(id uint64, name, content string) (*domains.NotifierEmailCampaignTemplate, error) {
	var tmRepo repositories.IEmailTemplateRepository
	err := container.Resolve(&tmRepo)
	if err != nil {
		return nil, err
	}
	tmp, err := tmRepo.Get(id)
	if err != nil {
		return nil, err
	}

	tmp.Name = name
	tmp.Content = content
	err = tmRepo.Update(tmp)
	if err != nil {
		return nil, err
	}
	return tmp, nil
}

func DeleteEmailTemplate(id uint64) error {
	var tmRepo repositories.IEmailTemplateRepository
	err := container.Resolve(&tmRepo)
	if err != nil {
		return err
	}
	err = tmRepo.Delete(&domains.NotifierEmailCampaignTemplate{ID: id})
	return err
}

func EmailTemplateList() ([]domains.NotifierEmailCampaignTemplate, error) {
	var tgRepo repositories.IEmailTemplateRepository
	err := container.Resolve(&tgRepo)
	if err != nil {
		return nil, err
	}
	var data []domains.NotifierEmailCampaignTemplate
	tgRepo.All(data)
	return data, nil
}

type EmailCampaignCreateData struct {
	EmailServiceId uint64
	ScheduledAt    *time.Time
	TemplateId     uint64
	StatusId       uint64
	FromEmail      string
	FromName       string
	Subject        string
	Name           string
	Tags           []uint64
}

func AddEmailCampaign(data *EmailCampaignCreateData) (*domains.NotifierEmailCampaign, error) {
	var tmRepo repositories.IEmailTemplateRepository
	err := container.Resolve(&tmRepo)
	if err != nil {
		return nil, err
	}
	temp, err := tmRepo.Get(data.TemplateId)
	if err != nil {
		return nil, err
	}

	var cmRepo repositories.IEmailCampaignRepository
	err = container.Resolve(&cmRepo)
	if err != nil {
		return nil, err
	}
	tmp := domains.NewNotifierEmailCampaign(
		data.EmailServiceId,
		data.ScheduledAt,
		data.TemplateId,
		data.StatusId,
		data.FromEmail,
		data.FromName,
		data.Subject,
		temp.Content,
		data.Name,
	)
	err = cmRepo.Create(tmp)
	if err != nil {
		return nil, err
	}

	err = cmRepo.AssignTagsToCampaign(tmp.ID, data.Tags)
	if err != nil {
		return nil, err
	}
	return tmp, nil
}

func DeleteEmailCampaign(campaign uint64) error {
	var cmRepo repositories.IEmailCampaignRepository
	err := container.Resolve(&cmRepo)
	if err != nil {
		return err
	}

	tmp, err := cmRepo.Get(campaign)
	if err != nil {
		return err
	}

	err = cmRepo.Delete(tmp)
	if err != nil {
		return err
	}
	return nil
}

type EmailCampaignUpdateData struct {
	EmailServiceId uint64
	ScheduledAt    *time.Time
	TemplateId     uint64
	StatusId       uint64
	FromEmail      string
	FromName       string
	Subject        string
	Name           string
	Tags           []uint64
}

func UpdateEmailCampaign(cmpId uint64, data *EmailCampaignUpdateData) error {
	var tmRepo repositories.IEmailTemplateRepository
	err := container.Resolve(&tmRepo)
	if err != nil {
		return err
	}
	temp, err := tmRepo.Get(data.TemplateId)
	if err != nil {
		return err
	}

	var cmRepo repositories.IEmailCampaignRepository
	err = container.Resolve(&cmRepo)
	if err != nil {
		return err
	}
	campaign, err := cmRepo.Get(cmpId)
	if err != nil {
		return err
	}

	campaign.FromEmail = data.FromEmail
	campaign.FromName = data.FromEmail
	campaign.StatusId = data.StatusId
	campaign.ScheduledAt = data.ScheduledAt
	campaign.TemplateId = temp.ID
	campaign.EmailServiceId = data.EmailServiceId
	campaign.Subject = data.Subject
	campaign.Content = temp.Content
	campaign.Name = data.Name
	campaign.UpdatedAt = time.Now()

	err = cmRepo.DeleteAllTagsForCampaign(cmpId)
	if err != nil {
		return err
	}
	err = cmRepo.AssignTagsToCampaign(campaign.ID, data.Tags)
	if err != nil {
		return err
	}
	return err
}
