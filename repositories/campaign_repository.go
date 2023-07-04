package repositories

import (
	"errors"
	"go-notifier-core/domains"
	"gorm.io/gorm"
	"time"
)

type IEmailTemplateRepository interface {
	IRepository[domains.NotifierEmailCampaignTemplate]
}

type gormEmailTemplateRepository struct {
	gormRepository[domains.NotifierEmailCampaignTemplate]
	db *gorm.DB
}

func NewGormEmailTemplateRepository(db *gorm.DB) IEmailTemplateRepository {
	return &gormEmailTemplateRepository{
		gormRepository: gormRepository[domains.NotifierEmailCampaignTemplate]{
			db: db,
		},
		db: db,
	}
}

type IEmailServiceRepository interface {
	IRepository[domains.NotifierEmailService]
}

type gormEmailServiceRepository struct {
	gormRepository[domains.NotifierEmailService]
	db *gorm.DB
}

func NewGormEmailServiceRepository(db *gorm.DB) IEmailServiceRepository {
	return &gormEmailServiceRepository{
		gormRepository: gormRepository[domains.NotifierEmailService]{
			db: db,
		},
		db: db,
	}
}

type IEmailStatusRepository interface {
	IRepository[domains.NotifierEmailStatus]
}

type gormEmailStatusRepository struct {
	gormRepository[domains.NotifierEmailStatus]
	db *gorm.DB
}

func NewGormEmailStatusRepository(db *gorm.DB) IEmailStatusRepository {
	return &gormEmailStatusRepository{
		gormRepository: gormRepository[domains.NotifierEmailStatus]{
			db: db,
		},
		db: db,
	}
}

type IEmailCampaignRepository interface {
	IRepository[domains.NotifierEmailCampaign]
	AssignTagsToCampaign(cmpId uint64, tagsId []uint64) error
	DeleteAllTagsForCampaign(cmpId uint64) error
	GetLatestCampaign() (*domains.NotifierEmailCampaign, error)
	GetCampaignTags(cmpId uint64) []domains.NotifierTag
}

type gormEmailCampaignRepository struct {
	gormRepository[domains.NotifierEmailCampaign]
	db *gorm.DB
}

func (g gormEmailCampaignRepository) AssignTagsToCampaign(cmpId uint64, tagsId []uint64) error {
	if len(tagsId) == 0 {
		return errors.New("tags id is empty")
	}
	tmp := make([]domains.NotifierEmailCampaignTag, len(tagsId))
	for i, tagId := range tagsId {
		t := domains.NewNotifierEmailCampaignTag(cmpId, tagId)
		tmp[i] = *t
	}

	res := g.db.Create(tmp)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (g gormEmailCampaignRepository) DeleteAllTagsForCampaign(cmpId uint64) error {
	res := g.db.Delete(domains.NotifierEmailCampaignTag{CampaignId: cmpId})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (g gormEmailCampaignRepository) GetLatestCampaign() (*domains.NotifierEmailCampaign, error) {
	var tmp domains.NotifierEmailCampaign
	res := g.db.Where("status_id = ?").
		Where("scheduled_at <= ?", time.Now()).
		First(&tmp)

	if res.Error != nil {
		return nil, res.Error
	}
	return &tmp, nil
}

func (g gormEmailCampaignRepository) GetCampaignTags(cmpId uint64) []domains.NotifierTag {
	var cmpTags []domains.NotifierEmailCampaignTag
	res := g.db.Where("campaign_id = ?", cmpId).Find(&cmpTags)
	if res.Error != nil || len(cmpTags) == 0 {
		return []domains.NotifierTag{}
	}
	var tmp []uint64
	for _, tag := range cmpTags {
		tmp = append(tmp, tag.TagId)
	}

	var tags []domains.NotifierTag
	res = g.db.Where("id in ?", tmp).
		Find(&tags)

	if res.Error != nil {
		return nil
	}

	return tags
}

func NewGormEmailCampaignRepository(db *gorm.DB) IEmailCampaignRepository {
	return &gormEmailCampaignRepository{
		gormRepository: gormRepository[domains.NotifierEmailCampaign]{
			db: db,
		},
		db: db,
	}
}
