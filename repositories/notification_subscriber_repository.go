package repositories

import (
	"errors"
	"go-notifier-core/domains"
	"gorm.io/gorm"
)

type INotificationSubscriberRepository interface {
	IRepository[domains.NotifierNotificationSubscriber]
	GetByNotification(token string) (*domains.NotifierNotificationSubscriber, error)
	AssignTagToUser(userId uint64, tagsId []uint64) error
	RemoveTagsFromUser(id uint64, entity []uint64) error
	GetSubscribersForTag(tagId uint64, data []domains.NotifierNotificationSubscriber)
	GetSubscribersForTagAndDriver(tagId, driverId uint64, data []domains.NotifierNotificationSubscriber)
}

type gormNotificationSubscriberRepository struct {
	gormRepository[domains.NotifierNotificationSubscriber]
	db *gorm.DB
}

func (g gormNotificationSubscriberRepository) GetByNotification(token string) (*domains.NotifierNotificationSubscriber, error) {
	var tmp domains.NotifierNotificationSubscriber
	res := g.db.Where("token = ?", token).First(&tmp)
	if res.Error != nil {
		return nil, res.Error
	}
	return &tmp, nil
}

func (g gormNotificationSubscriberRepository) AssignTagToUser(userId uint64, tagsId []uint64) error {
	if len(tagsId) == 0 {
		return errors.New("tags id is empty")
	}
	tmp := make([]domains.NotifierNotificationSubTag, len(tagsId))
	for i, tagId := range tagsId {
		t := domains.NewNotifierNotificationSubTag(userId, tagId)
		tmp[i] = *t
	}

	res := g.db.Create(tmp)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (g gormNotificationSubscriberRepository) RemoveTagsFromUser(userId uint64, tagsId []uint64) error {
	if len(tagsId) == 0 {
		return errors.New("tags id is empty")
	}
	tmp := make([]domains.NotifierNotificationSubTag, len(tagsId))
	for i, tagId := range tagsId {
		t := domains.NewNotifierNotificationSubTag(userId, tagId)
		tmp[i] = *t
	}

	res := g.db.Delete(domains.NotifierNotificationSubTag{NotificationSubscriberId: userId}, "tag_id in ?", tagsId)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (g gormNotificationSubscriberRepository) GetSubscribersForTag(tagId uint64, data []domains.NotifierNotificationSubscriber) {
	_ = g.db.Scopes(tagIdScope(tagId)).Find(data)
}

func (g gormNotificationSubscriberRepository) GetSubscribersForTagAndDriver(tagId, driverId uint64, data []domains.NotifierNotificationSubscriber) {
	_ = g.db.Scopes(tagIdScope(tagId), driverIdScope(driverId)).Find(data)
}

func driverIdScope(driverId uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("driver_id = ?", driverId)
	}
}

func NewGormNotificationSubscriberRepository(db *gorm.DB) INotificationSubscriberRepository {
	return &gormNotificationSubscriberRepository{
		gormRepository: gormRepository[domains.NotifierNotificationSubscriber]{
			db: db,
		},
		db: db,
	}
}

type INotificationSubTagRepository interface {
	IRepository[domains.NotifierNotificationSubTag]
}

type gormNotificationSubTagRepository struct {
	gormRepository[domains.NotifierNotificationSubTag]
	db *gorm.DB
}

func NewGormNotificationSubTagRepository(db *gorm.DB) INotificationSubTagRepository {
	return &gormNotificationSubTagRepository{
		gormRepository: gormRepository[domains.NotifierNotificationSubTag]{
			db: db,
		},
		db: db,
	}
}

type INotifierNotificationDriverRepository interface {
	IRepository[domains.NotifierNotificationDriver]
}

type gormNotifierNotificationDriverRepository struct {
	gormRepository[domains.NotifierNotificationDriver]
	db *gorm.DB
}

func NewGormNotifierNotificationDriverRepository(db *gorm.DB) INotifierNotificationDriverRepository {
	return &gormNotifierNotificationDriverRepository{
		gormRepository: gormRepository[domains.NotifierNotificationDriver]{
			db: db,
		},
		db: db,
	}
}
