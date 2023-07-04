package repositories

import (
	"errors"
	"go-notifier-core/domains"
	"gorm.io/gorm"
)

type IMobileSubscriberRepository interface {
	IRepository[domains.NotifierMobileSubscriber]
	GetByMobile(mobile string) (*domains.NotifierMobileSubscriber, error)
	AssignTagToUser(userId uint64, tagsId []uint64) error
	RemoveTagsFromUser(id uint64, entity []uint64) error
	GetSubscribersForTag(tagId uint64, data []domains.NotifierMobileSubscriber)
	GetUnSubscribed(data []domains.NotifierMobileSubscriber)
}

type gormMobileSubscriberRepository struct {
	gormRepository[domains.NotifierMobileSubscriber]
	db *gorm.DB
}

func (g gormMobileSubscriberRepository) GetByMobile(mobile string) (*domains.NotifierMobileSubscriber, error) {
	var tmp domains.NotifierMobileSubscriber
	res := g.db.Where("mobile = ?", mobile).First(&tmp)
	if res.Error != nil {
		return nil, res.Error
	}
	return &tmp, nil
}

func (g gormMobileSubscriberRepository) AssignTagToUser(userId uint64, tagsId []uint64) error {
	if len(tagsId) == 0 {
		return errors.New("tags id is empty")
	}
	tmp := make([]domains.NotifierMobileSubTag, len(tagsId))
	for i, tagId := range tagsId {
		t := domains.NewNotifierMobileSubTag(userId, tagId)
		tmp[i] = *t
	}

	res := g.db.Create(tmp)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (g gormMobileSubscriberRepository) RemoveTagsFromUser(userId uint64, tagsId []uint64) error {
	if len(tagsId) == 0 {
		return errors.New("tags id is empty")
	}
	tmp := make([]domains.NotifierMobileSubTag, len(tagsId))
	for i, tagId := range tagsId {
		t := domains.NewNotifierMobileSubTag(userId, tagId)
		tmp[i] = *t
	}

	res := g.db.Delete(domains.NotifierMobileSubTag{MobileSubscriberId: userId}, "tag_id in ?", tagsId)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (g gormMobileSubscriberRepository) GetSubscribersForTag(tagId uint64, data []domains.NotifierMobileSubscriber) {
	_ = g.db.Scopes(exceptUnsubscribedScope, tagIdScope(tagId)).Find(data)
}

func (g gormMobileSubscriberRepository) GetUnSubscribed(data []domains.NotifierMobileSubscriber) {
	_ = g.db.Scopes(unsubscribedScope).Find(data)
}

func NewGormMobileSubscriberRepository(db *gorm.DB) IMobileSubscriberRepository {
	return &gormMobileSubscriberRepository{
		gormRepository: gormRepository[domains.NotifierMobileSubscriber]{
			db: db,
		},
		db: db,
	}
}

type IMobileSubTagRepository interface {
	IRepository[domains.NotifierMobileSubTag]
}

type gormMobileSubTagRepository struct {
	gormRepository[domains.NotifierMobileSubTag]
	db *gorm.DB
}

func NewGormMobileSubTagRepository(db *gorm.DB) IMobileSubTagRepository {
	return &gormMobileSubTagRepository{
		gormRepository: gormRepository[domains.NotifierMobileSubTag]{
			db: db,
		},
		db: db,
	}
}

type IMobileUnSubEventRepository interface {
	IRepository[domains.NotifierMobileUnsubscribeEvent]
}

type gormMobileUnSubEventRepository struct {
	gormRepository[domains.NotifierMobileUnsubscribeEvent]
	db *gorm.DB
}

func NewGormMobileUnSubEventRepository(db *gorm.DB) IMobileUnSubEventRepository {
	return &gormMobileUnSubEventRepository{
		gormRepository: gormRepository[domains.NotifierMobileUnsubscribeEvent]{
			db: db,
		},
		db: db,
	}
}
