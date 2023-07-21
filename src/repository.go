package go_notifier_core

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type (
	NotFoundError struct {
	}
)

func (n NotFoundError) Error() string {
	return "record not found"
}

type IRepository[Model interface{}] interface {
	Create(*Model) error
	Update(*Model) error
	Delete(*Model) error
	Get(id uint64) (*Model, error)
	All(data *[]Model)
}

type gormRepository[m interface{}] struct {
	db *gorm.DB
}

func (g gormRepository[m]) Create(model *m) error {
	res := g.db.Create(model)
	return res.Error
}

func (g gormRepository[m]) Update(model *m) error {
	res := g.db.Save(model)
	return res.Error
}

func (g gormRepository[m]) Delete(model *m) error {
	res := g.db.Delete(model)
	return res.Error
}

func (g gormRepository[m]) Get(id uint64) (*m, error) {
	var x m
	res := g.db.First(&x, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, NotFoundError{}
	} else {
		return &x, res.Error
	}
}

func (g gormRepository[m]) All(data *[]m) {
	g.db.Find(data)
}

//Campaign repositories

type IEmailTemplateRepository interface {
	IRepository[NotifierEmailCampaignTemplate]
}

type gormEmailTemplateRepository struct {
	gormRepository[NotifierEmailCampaignTemplate]
	db *gorm.DB
}

func NewGormEmailTemplateRepository(db *gorm.DB) IEmailTemplateRepository {
	return &gormEmailTemplateRepository{
		gormRepository: gormRepository[NotifierEmailCampaignTemplate]{
			db: db,
		},
		db: db,
	}
}

type IEmailServiceRepository interface {
	IRepository[NotifierEmailService]
}

type gormEmailServiceRepository struct {
	gormRepository[NotifierEmailService]
	db *gorm.DB
}

func NewGormEmailServiceRepository(db *gorm.DB) IEmailServiceRepository {
	return &gormEmailServiceRepository{
		gormRepository: gormRepository[NotifierEmailService]{
			db: db,
		},
		db: db,
	}
}

type IEmailStatusRepository interface {
	IRepository[NotifierEmailCampaignStatus]
	FirstOrCreate(status *NotifierEmailCampaignStatus) error
}

type gormEmailStatusRepository struct {
	gormRepository[NotifierEmailCampaignStatus]
	db *gorm.DB
}

func (g gormEmailStatusRepository) FirstOrCreate(status *NotifierEmailCampaignStatus) error {
	return g.db.FirstOrCreate(status).Error
}

func NewGormEmailStatusRepository(db *gorm.DB) IEmailStatusRepository {
	return &gormEmailStatusRepository{
		gormRepository: gormRepository[NotifierEmailCampaignStatus]{
			db: db,
		},
		db: db,
	}
}

type IEmailCampaignRepository interface {
	IRepository[NotifierEmailCampaign]
	AssignTagsToCampaign(cmpId uint64, tagsId []uint64) error
	DeleteAllTagsForCampaign(cmpId uint64) error
	GetLatestCampaign() (*NotifierEmailCampaign, error)
	GetCampaignTags(cmpId uint64) []NotifierTag
}

type gormEmailCampaignRepository struct {
	gormRepository[NotifierEmailCampaign]
	db *gorm.DB
}

func (g gormEmailCampaignRepository) AssignTagsToCampaign(cmpId uint64, tagsId []uint64) error {
	if len(tagsId) == 0 {
		return errors.New("tags id is empty")
	}
	tmp := make([]NotifierEmailCampaignTag, len(tagsId))
	for i, tagId := range tagsId {
		t := NewNotifierEmailCampaignTag(cmpId, tagId)
		tmp[i] = *t
	}

	res := g.db.Create(tmp)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (g gormEmailCampaignRepository) DeleteAllTagsForCampaign(cmpId uint64) error {
	res := g.db.Where("campaign_id = ?", cmpId).Delete(&NotifierEmailCampaignTag{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (g gormEmailCampaignRepository) GetLatestCampaign() (*NotifierEmailCampaign, error) {
	var tmp NotifierEmailCampaign
	res := g.db.Where("status_id = ?", NotifierEmailStatusDraft).
		Where("scheduled_at <= ? or scheduled_at IS NULL", time.Now()).
		Order("ID asc").
		First(&tmp)

	if res.Error != nil {
		return nil, res.Error
	}
	return &tmp, nil
}

func (g gormEmailCampaignRepository) GetCampaignTags(cmpId uint64) []NotifierTag {
	var cmpTags []NotifierEmailCampaignTag
	res := g.db.Where("campaign_id = ?", cmpId).Find(&cmpTags)
	if res.Error != nil || len(cmpTags) == 0 {
		return []NotifierTag{}
	}
	var tmp []uint64
	for _, tag := range cmpTags {
		tmp = append(tmp, tag.TagId)
	}

	var tags []NotifierTag
	res = g.db.Where("id in ?", tmp).
		Find(&tags)

	if res.Error != nil {
		return nil
	}

	return tags
}

func NewGormEmailCampaignRepository(db *gorm.DB) IEmailCampaignRepository {
	return &gormEmailCampaignRepository{
		gormRepository: gormRepository[NotifierEmailCampaign]{
			db: db,
		},
		db: db,
	}
}

// Email repositories

type IEmailSubscriberRepository interface {
	IRepository[NotifierEmailSubscriber]
	GetByEmail(email string) (*NotifierEmailSubscriber, error)
	AssignTagToUser(userId uint64, tagsId []uint64) error
	RemoveTagsFromUser(id uint64, entity []uint64) error
	GetSubscribersForTag(tagId uint64, data *[]NotifierEmailSubscriber)
	GetUnSubscribed(data *[]NotifierEmailSubscriber)
	GetUsersByTagId(tags []NotifierTag, data *[]NotifierEmailSubscriber)
	GetByEmailWithTags(email string) (*NotifierEmailSubscriber, error)
}

type gormEmailSubscriberRepository struct {
	gormRepository[NotifierEmailSubscriber]
	db *gorm.DB
}

func (g gormEmailSubscriberRepository) GetByEmail(email string) (*NotifierEmailSubscriber, error) {
	var tmp NotifierEmailSubscriber
	res := g.db.Where("email = ?", email).First(&tmp)
	if res.Error != nil {
		return nil, res.Error
	}
	return &tmp, nil
}
func (g gormEmailSubscriberRepository) GetByEmailWithTags(email string) (*NotifierEmailSubscriber, error) {
	var tmp NotifierEmailSubscriber
	res := g.db.Preload("Tags").Where("email = ?", email).First(&tmp)
	if res.Error != nil {
		return nil, res.Error
	}
	return &tmp, nil
}

func (g gormEmailSubscriberRepository) GetUsersByTagId(tags []NotifierTag, data *[]NotifierEmailSubscriber) {
	ids := make([]uint64, len(tags))
	for i := 0; i < len(tags); i++ {
		ids[i] = tags[i].ID
	}
	_ = g.db.
		Table("notifier_email_subscribers AS subs").
		Select("DISTINCT subs.*").
		Where("subs.unsubscribed_event_id IS NULL AND subs.unsubscribed_at IS NULL").
		Joins("INNER JOIN notifier_email_sub_tags AS sub_tags ON subs.id = sub_tags.email_subscriber_id AND sub_tags.tag_id IN ?", ids).
		Find(data)
}

func (g gormEmailSubscriberRepository) AssignTagToUser(userId uint64, tagsId []uint64) error {
	if len(tagsId) == 0 {
		return errors.New("tags id is empty")
	}
	tmp := make([]NotifierEmailSubTag, len(tagsId))
	for i, tagId := range tagsId {
		t := NewNotifierEmailSubTag(userId, tagId)
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
	tmp := make([]NotifierEmailSubTag, len(tagsId))
	for i, tagId := range tagsId {
		t := NewNotifierEmailSubTag(userId, tagId)
		tmp[i] = *t
	}

	res := g.db.Delete(NotifierEmailSubTag{EmailSubscriberId: userId}, "tag_id in ?", tagsId)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (g gormEmailSubscriberRepository) GetSubscribersForTag(tagId uint64, data *[]NotifierEmailSubscriber) {
	_ = g.db.Scopes(exceptUnsubscribedScope).
		Where("id IN (SELECT email_subscriber_id FROM notifier_email_sub_tags WHERE tag_id = ?)", tagId).
		Find(data)
}

func (g gormEmailSubscriberRepository) GetUnSubscribed(data *[]NotifierEmailSubscriber) {
	_ = g.db.Scopes(unsubscribedScope).Find(data)
}

func NewGormEmailSubscriberRepository(db *gorm.DB) IEmailSubscriberRepository {
	return &gormEmailSubscriberRepository{
		gormRepository: gormRepository[NotifierEmailSubscriber]{
			db: db,
		},
		db: db,
	}
}

type IEmailSubTagRepository interface {
	IRepository[NotifierEmailSubTag]
}

type gormEmailSubTagRepository struct {
	gormRepository[NotifierEmailSubTag]
	db *gorm.DB
}

func NewGormEmailSubTagRepository(db *gorm.DB) IEmailSubTagRepository {
	return &gormEmailSubTagRepository{
		gormRepository: gormRepository[NotifierEmailSubTag]{
			db: db,
		},
		db: db,
	}
}

type IEmailUnSubEventRepository interface {
	IRepository[NotifierEmailUnsubscribeEvent]
	FirstOrCreate(status *NotifierEmailUnsubscribeEvent) error
}

type gormEmailUnSubEventRepository struct {
	gormRepository[NotifierEmailUnsubscribeEvent]
	db *gorm.DB
}

func (g gormEmailUnSubEventRepository) FirstOrCreate(status *NotifierEmailUnsubscribeEvent) error {
	return g.db.FirstOrCreate(status).Error
}

func NewGormEmailUnSubEventRepository(db *gorm.DB) IEmailUnSubEventRepository {
	return &gormEmailUnSubEventRepository{
		gormRepository: gormRepository[NotifierEmailUnsubscribeEvent]{
			db: db,
		},
		db: db,
	}
}

type IEmailMessageRepository interface {
	IRepository[NotifierEmailMessage]
	CheckMessageExists(message *NotifierEmailMessage) error
}

type gormEmailMessageRepository struct {
	gormRepository[NotifierEmailMessage]
	db *gorm.DB
}

func (g gormEmailMessageRepository) CheckMessageExists(message *NotifierEmailMessage) error {
	err := g.db.Where("subscriber_id = ? AND source_id = ? AND source_type like ?", message.SubscriberId, message.SourceId, "%"+message.SourceType+"%").First(message)
	return err.Error
}

func NewGormEmailMessageRepository(db *gorm.DB) IEmailMessageRepository {
	return &gormEmailMessageRepository{
		gormRepository: gormRepository[NotifierEmailMessage]{
			db: db,
		},
		db: db,
	}
}

// Mobile repositories

type IMobileSubscriberRepository interface {
	IRepository[NotifierMobileSubscriber]
	GetByMobile(mobile string) (*NotifierMobileSubscriber, error)
	AssignTagToUser(userId uint64, tagsId []uint64) error
	RemoveTagsFromUser(id uint64, entity []uint64) error
	GetSubscribersForTag(tagId uint64, data []NotifierMobileSubscriber)
	GetUnSubscribed(data []NotifierMobileSubscriber)
}

type gormMobileSubscriberRepository struct {
	gormRepository[NotifierMobileSubscriber]
	db *gorm.DB
}

func (g gormMobileSubscriberRepository) GetByMobile(mobile string) (*NotifierMobileSubscriber, error) {
	var tmp NotifierMobileSubscriber
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
	tmp := make([]NotifierMobileSubTag, len(tagsId))
	for i, tagId := range tagsId {
		t := NewNotifierMobileSubTag(userId, tagId)
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
	tmp := make([]NotifierMobileSubTag, len(tagsId))
	for i, tagId := range tagsId {
		t := NewNotifierMobileSubTag(userId, tagId)
		tmp[i] = *t
	}

	res := g.db.Delete(NotifierMobileSubTag{MobileSubscriberId: userId}, "tag_id in ?", tagsId)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (g gormMobileSubscriberRepository) GetSubscribersForTag(tagId uint64, data []NotifierMobileSubscriber) {
	_ = g.db.Scopes(exceptUnsubscribedScope, tagIdScope(tagId)).Find(data)
}

func (g gormMobileSubscriberRepository) GetUnSubscribed(data []NotifierMobileSubscriber) {
	_ = g.db.Scopes(unsubscribedScope).Find(data)
}

func NewGormMobileSubscriberRepository(db *gorm.DB) IMobileSubscriberRepository {
	return &gormMobileSubscriberRepository{
		gormRepository: gormRepository[NotifierMobileSubscriber]{
			db: db,
		},
		db: db,
	}
}

type IMobileSubTagRepository interface {
	IRepository[NotifierMobileSubTag]
}

type gormMobileSubTagRepository struct {
	gormRepository[NotifierMobileSubTag]
	db *gorm.DB
}

func NewGormMobileSubTagRepository(db *gorm.DB) IMobileSubTagRepository {
	return &gormMobileSubTagRepository{
		gormRepository: gormRepository[NotifierMobileSubTag]{
			db: db,
		},
		db: db,
	}
}

type IMobileUnSubEventRepository interface {
	IRepository[NotifierMobileUnsubscribeEvent]
}

type gormMobileUnSubEventRepository struct {
	gormRepository[NotifierMobileUnsubscribeEvent]
	db *gorm.DB
}

func NewGormMobileUnSubEventRepository(db *gorm.DB) IMobileUnSubEventRepository {
	return &gormMobileUnSubEventRepository{
		gormRepository: gormRepository[NotifierMobileUnsubscribeEvent]{
			db: db,
		},
		db: db,
	}
}

// Notification repositories

type INotificationSubscriberRepository interface {
	IRepository[NotifierNotificationSubscriber]
	GetByNotification(token string) (*NotifierNotificationSubscriber, error)
	AssignTagToUser(userId uint64, tagsId []uint64) error
	RemoveTagsFromUser(id uint64, entity []uint64) error
	GetSubscribersForTag(tagId uint64, data []NotifierNotificationSubscriber)
	GetSubscribersForTagAndDriver(tagId, driverId uint64, data []NotifierNotificationSubscriber)
}

type gormNotificationSubscriberRepository struct {
	gormRepository[NotifierNotificationSubscriber]
	db *gorm.DB
}

func (g gormNotificationSubscriberRepository) GetByNotification(token string) (*NotifierNotificationSubscriber, error) {
	var tmp NotifierNotificationSubscriber
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
	tmp := make([]NotifierNotificationSubTag, len(tagsId))
	for i, tagId := range tagsId {
		t := NewNotifierNotificationSubTag(userId, tagId)
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
	tmp := make([]NotifierNotificationSubTag, len(tagsId))
	for i, tagId := range tagsId {
		t := NewNotifierNotificationSubTag(userId, tagId)
		tmp[i] = *t
	}

	res := g.db.Delete(NotifierNotificationSubTag{NotificationSubscriberId: userId}, "tag_id in ?", tagsId)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (g gormNotificationSubscriberRepository) GetSubscribersForTag(tagId uint64, data []NotifierNotificationSubscriber) {
	_ = g.db.Scopes(tagIdScope(tagId)).Find(data)
}

func (g gormNotificationSubscriberRepository) GetSubscribersForTagAndDriver(tagId, driverId uint64, data []NotifierNotificationSubscriber) {
	_ = g.db.Scopes(tagIdScope(tagId), driverIdScope(driverId)).Find(data)
}

func driverIdScope(driverId uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("driver_id = ?", driverId)
	}
}

func NewGormNotificationSubscriberRepository(db *gorm.DB) INotificationSubscriberRepository {
	return &gormNotificationSubscriberRepository{
		gormRepository: gormRepository[NotifierNotificationSubscriber]{
			db: db,
		},
		db: db,
	}
}

type INotificationSubTagRepository interface {
	IRepository[NotifierNotificationSubTag]
}

type gormNotificationSubTagRepository struct {
	gormRepository[NotifierNotificationSubTag]
	db *gorm.DB
}

func NewGormNotificationSubTagRepository(db *gorm.DB) INotificationSubTagRepository {
	return &gormNotificationSubTagRepository{
		gormRepository: gormRepository[NotifierNotificationSubTag]{
			db: db,
		},
		db: db,
	}
}

type INotifierNotificationDriverRepository interface {
	IRepository[NotifierNotificationDriver]
}

type gormNotifierNotificationDriverRepository struct {
	gormRepository[NotifierNotificationDriver]
	db *gorm.DB
}

func NewGormNotifierNotificationDriverRepository(db *gorm.DB) INotifierNotificationDriverRepository {
	return &gormNotifierNotificationDriverRepository{
		gormRepository: gormRepository[NotifierNotificationDriver]{
			db: db,
		},
		db: db,
	}
}

//Tag repositories

type ITagRepository interface {
	IRepository[NotifierTag]
	GetByName(name string) (*NotifierTag, error)
}

type gormTagRepository struct {
	gormRepository[NotifierTag]
	db *gorm.DB
}

func (g gormTagRepository) GetByName(name string) (*NotifierTag, error) {
	var x NotifierTag
	res := g.db.Where("name = ?", name).First(&x)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, NotFoundError{}
	} else {
		return &x, res.Error
	}
}

func NewGormTagRepository(db *gorm.DB) ITagRepository {
	return &gormTagRepository{
		gormRepository: gormRepository[NotifierTag]{
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
