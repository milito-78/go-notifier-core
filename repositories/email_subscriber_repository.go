package repositories

import (
	"errors"
	"go-notifier-core/domains"
	"gorm.io/gorm"
)

type IEmailSubscriberRepository interface {
	IRepository[domains.NotifierEmailSubscriber]
	GetByEmail(email string) (*domains.NotifierEmailSubscriber, error)
	AssignTagToUser(userId uint64, tagsId []uint64) error
	RemoveTagsFromUser(id uint64, entity []uint64) error
	GetSubscribersForTag(tagId uint64, data []domains.NotifierEmailSubscriber)
	GetUnSubscribed(data []domains.NotifierEmailSubscriber)
	GetUsersByTagId(tags []domains.NotifierTag) []domains.NotifierEmailSubscriber
}

type gormEmailSubscriberRepository struct {
	gormRepository[domains.NotifierEmailSubscriber]
	db *gorm.DB
}

func (g gormEmailSubscriberRepository) GetByEmail(email string) (*domains.NotifierEmailSubscriber, error) {
	var tmp domains.NotifierEmailSubscriber
	res := g.db.Where("email = ?", email).First(&tmp)
	if res.Error != nil {
		return nil, res.Error
	}
	return &tmp, nil
}

func (g gormEmailSubscriberRepository) GetUsersByTagId(tags []domains.NotifierTag) []domains.NotifierEmailSubscriber {

	var data []domains.NotifierEmailSubscriber
	_ = g.db.
		Table("notifier_email_subscribers AS subs").
		Select("DISTINCT subs.*").
		Where("subs.unsubscribed_event_id IS NOT NULL AND subs.unsubscribed_at IS NOT NULL").
		Joins("INNER JOIN notifier_email_sub_tags AS sub_tags ON subs.id = sub_tags.subscriber_id AND sub_tags.tag_id IN ?", tags).
		Find(&data)

	return data
}

func (g gormEmailSubscriberRepository) AssignTagToUser(userId uint64, tagsId []uint64) error {
	if len(tagsId) == 0 {
		return errors.New("tags id is empty")
	}
	tmp := make([]domains.NotifierEmailSubTag, len(tagsId))
	for i, tagId := range tagsId {
		t := domains.NewNotifierEmailSubTag(userId, tagId)
		tmp[i] = *t
	}

	res := g.db.Create(tmp)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (g gormEmailSubscriberRepository) RemoveTagsFromUser(userId uint64, tagsId []uint64) error {
	if len(tagsId) == 0 {
		return errors.New("tags id is empty")
	}
	tmp := make([]domains.NotifierEmailSubTag, len(tagsId))
	for i, tagId := range tagsId {
		t := domains.NewNotifierEmailSubTag(userId, tagId)
		tmp[i] = *t
	}

	res := g.db.Delete(domains.NotifierEmailSubTag{EmailSubscriberId: userId}, "tag_id in ?", tagsId)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (g gormEmailSubscriberRepository) GetSubscribersForTag(tagId uint64, data []domains.NotifierEmailSubscriber) {
	_ = g.db.Scopes(exceptUnsubscribedScope, tagIdScope(tagId)).Find(data)
}

func (g gormEmailSubscriberRepository) GetUnSubscribed(data []domains.NotifierEmailSubscriber) {
	_ = g.db.Scopes(unsubscribedScope).Find(data)
}

func NewGormEmailSubscriberRepository(db *gorm.DB) IEmailSubscriberRepository {
	return &gormEmailSubscriberRepository{
		gormRepository: gormRepository[domains.NotifierEmailSubscriber]{
			db: db,
		},
		db: db,
	}
}

type IEmailSubTagRepository interface {
	IRepository[domains.NotifierEmailSubTag]
}

type gormEmailSubTagRepository struct {
	gormRepository[domains.NotifierEmailSubTag]
	db *gorm.DB
}

func NewGormEmailSubTagRepository(db *gorm.DB) IEmailSubTagRepository {
	return &gormEmailSubTagRepository{
		gormRepository: gormRepository[domains.NotifierEmailSubTag]{
			db: db,
		},
		db: db,
	}
}

type IEmailUnSubEventRepository interface {
	IRepository[domains.NotifierEmailUnsubscribeEvent]
}

type gormEmailUnSubEventRepository struct {
	gormRepository[domains.NotifierEmailUnsubscribeEvent]
	db *gorm.DB
}

func NewGormEmailUnSubEventRepository(db *gorm.DB) IEmailUnSubEventRepository {
	return &gormEmailUnSubEventRepository{
		gormRepository: gormRepository[domains.NotifierEmailUnsubscribeEvent]{
			db: db,
		},
		db: db,
	}
}

type IEmailMessageRepository interface {
	IRepository[domains.NotifierEmailMessage]
	CheckMessageExists(message *domains.NotifierEmailMessage) error
}

type gormEmailMessageRepository struct {
	gormRepository[domains.NotifierEmailMessage]
	db *gorm.DB
}

func (g gormEmailMessageRepository) CheckMessageExists(message *domains.NotifierEmailMessage) error {
	err := g.db.Where("subscriber_id = ? AND source_id = ? AND source_type like ?", message.SubscriberId, message.SourceId, "%"+message.SourceType+"%").First(message)
	return err.Error
}

func NewGormEmailMessageRepository(db *gorm.DB) IEmailMessageRepository {
	return &gormEmailMessageRepository{
		gormRepository: gormRepository[domains.NotifierEmailMessage]{
			db: db,
		},
		db: db,
	}
}
