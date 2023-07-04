package repositories

import (
	"errors"
	"go-notifier-core/domains"
	"gorm.io/gorm"
)

type ITagRepository interface {
	IRepository[domains.NotifierTag]
	GetByName(name string) (*domains.NotifierTag, error)
}

type gormTagRepository struct {
	gormRepository[domains.NotifierTag]
	db *gorm.DB
}

func (g gormTagRepository) GetByName(name string) (*domains.NotifierTag, error) {
	var x domains.NotifierTag
	res := g.db.Where("name = ?", name).First(&x)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, NotFoundError{}
	} else {
		return &x, res.Error
	}
}

func NewGormTagRepository(db *gorm.DB) ITagRepository {
	return &gormTagRepository{
		gormRepository: gormRepository[domains.NotifierTag]{
			db: db,
		},
		db: db,
	}
}

func exceptUnsubscribedScope(db *gorm.DB) *gorm.DB {
	return db.Where("unsubscribed_event_id is null and unsubscribed_at is null")
}
func unsubscribedScope(db *gorm.DB) *gorm.DB {
	return db.Where("unsubscribed_event_id is not null and unsubscribed_at is not null")
}

func tagIdScope(tagId uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("tag_id = ?", tagId)
	}
}
