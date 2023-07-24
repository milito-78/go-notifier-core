package go_notifier_core

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewNotifierEmailCampaignTemplate(t *testing.T) {
	// Test data
	content := "Test email content"
	name := "Test Campaign"

	// Call the function to create a new NotifierEmailCampaignTemplate
	campaign := NewNotifierEmailCampaignTemplate(content, name)

	// Check if the campaign is not nil
	if campaign == nil {
		t.Error("Expected campaign to be created, but got nil")
	}

	// Check if the content and name are set correctly
	if campaign.Content != content {
		t.Errorf("Expected content to be '%s', but got '%s'", content, campaign.Content)
	}

	if campaign.Name != name {
		t.Errorf("Expected name to be '%s', but got '%s'", name, campaign.Name)
	}

	// Check if CreatedAt and UpdatedAt fields are set to the current time
	currentTime := time.Now()
	if campaign.CreatedAt.After(currentTime) || campaign.CreatedAt.Before(currentTime.Add(-1*time.Second)) {
		t.Error("Expected CreatedAt to be close to the current time")
	}

	if campaign.UpdatedAt.After(currentTime) || campaign.UpdatedAt.Before(currentTime.Add(-1*time.Second)) {
		t.Error("Expected UpdatedAt to be close to the current time")
	}
}

func TestNewNotifierNewNotifierEmailService(t *testing.T) {
	config := SmtpConfig{
		Host:       "testHost",
		Port:       "port",
		Username:   "username",
		Password:   "password",
		Encryption: "tls",
	}
	payload, _ := json.Marshal(config)
	sType := NotifierEmailServiceSMTPType
	name := "Test SMTP"

	// Call the function to create a new NewNotifierEmailService
	service := NewNotifierEmailService(string(payload), sType, name)

	// Check if the service is not nil
	if service == nil {
		t.Error("Expected service to be created, but got nil")
	}

	// Check if the name are set correctly
	if service.Name != name {
		t.Errorf("Expected name to be '%s', but got '%s'", name, service.Name)
	}

	if service.Payload != string(payload) {
		t.Errorf("Expected payload to be '%s', but got '%s'", payload, service.Payload)
	}

	if service.Type != NotifierEmailServiceSMTPType {
		t.Errorf("Expected type to be '%s', but got '%s'", sType, service.Type)
	}

	var tmp SmtpConfig
	err := json.Unmarshal([]byte(service.Payload), &tmp)
	if err != nil {
		t.Errorf("Error during unmarshal payload '%s'", err)
	} else {
		if config.Password != tmp.Password ||
			config.Username != tmp.Username ||
			config.Host != tmp.Host ||
			config.Port != tmp.Port {
			t.Errorf("Payload unmarshal isn't equal to config : '%+v' , '%+v'", config, tmp)
		}
	}
}

func TestNewNotifierEmailStatus(t *testing.T) {
	name := "Sending test"
	status := NewNotifierEmailStatus(name, NotifierEmailStatusSending)

	// Check if the status is not nil
	if status == nil {
		t.Error("Expected status to be created, but got nil")
	}

	if status.Name != name {
		t.Errorf("Expected name to be '%s', but got '%s'", name, status.Name)
	}

	if status.ID != NotifierEmailStatusSending {
		t.Errorf("Expected id to be '%d', but got '%d'", NotifierEmailStatusSending, status.ID)
	}
}

func TestNewNotifierEmailCampaign(t *testing.T) {
	var emailServiceId = uint64(1)
	scheduledAt := time.Now().Add(time.Second * 30)
	var templateId = uint64(2)
	var statusId = uint64(NotifierEmailStatusSent)
	fromEmail := "from@mail.com"
	fromName := "from me"
	subject := "test it"
	content := "test content that's picked from template"
	name := "a name"

	campaign := NewNotifierEmailCampaign(emailServiceId, &scheduledAt, templateId, statusId, fromEmail, fromName, subject, content, name)
	if campaign == nil {
		t.Error("Expected campaign to be created, but got nil")
	}

	if campaign.Name != name {
		t.Errorf("Expected name to be '%s', but got '%s'", name, campaign.Name)
	}
	if campaign.EmailServiceId != emailServiceId {
		t.Errorf("Expected email service to be '%d', but got '%d'", emailServiceId, campaign.EmailServiceId)
	}
	if campaign.TemplateId != templateId {
		t.Errorf("Expected template id to be '%d', but got '%d'", templateId, campaign.TemplateId)
	}
	if campaign.StatusId != statusId {
		t.Errorf("Expected status id to be '%d', but got '%d'", statusId, campaign.StatusId)
	}
	if campaign.FromEmail != fromEmail {
		t.Errorf("Expected from email to be '%s', but got '%s'", fromEmail, campaign.FromEmail)
	}
	if campaign.FromName != fromName {
		t.Errorf("Expected from name to be '%s', but got '%s'", fromName, campaign.FromName)
	}
	if campaign.Subject != subject {
		t.Errorf("Expected subject  to be '%s', but got '%s'", subject, campaign.Subject)
	}
	if campaign.Content != content {
		t.Errorf("Expected name to be '%s', but got '%s'", name, campaign.Name)
	}
	if campaign.ScheduledAt != nil && !campaign.ScheduledAt.Equal(scheduledAt) {
		t.Errorf("Scheduled at is not equals '%+v', but got '%+v'", scheduledAt, campaign.ScheduledAt)
	}
}

func TestNewNotifierEmailCampaignTag(t *testing.T) {
	var campaignId = uint64(1)
	var tagId = uint64(2)

	campaignTag := NewNotifierEmailCampaignTag(campaignId, tagId)
	if campaignTag == nil {
		t.Error("Expected campaignTag to be created, but got nil")
	}

	if campaignTag.CampaignId != campaignId {
		t.Errorf("Expected campaign id to be '%d', but got '%d'", campaignId, campaignTag.CampaignId)
	}

	if campaignTag.TagId != tagId {
		t.Errorf("Expected tag id to be '%d', but got '%d'", tagId, campaignTag.TagId)
	}
}

func TestNewNotifierEmailUnsubscribeEvent(t *testing.T) {
	reason := "test reason"
	var ID = uint64(NotifierEmailUnsubBounce)

	event := NewNotifierEmailUnsubscribeEvent(reason, ID)
	if event == nil {
		t.Error("Expected event to be created, but got nil")
	}

	if event.Reason != reason {
		t.Errorf("Expected reason to be '%s', but got '%s'", reason, event.Reason)
	}

	if event.ID != ID {
		t.Errorf("Expected id to be '%d', but got '%d'", ID, event.ID)
	}
}

func TestNewNotifierEmailSubscriber(t *testing.T) {
	email := "email@test.com"
	first := "first"
	last := "last"

	subscriber := NewNotifierEmailSubscriber(email, first, last)
	if subscriber == nil {
		t.Error("Expected subscriber to be created, but got nil")
	}

	if subscriber.FirstName != first {
		t.Errorf("Expected first to be '%s', but got '%s'", first, subscriber.FirstName)
	}

	if subscriber.LastName != last {
		t.Errorf("Expected last to be '%s', but got '%s'", first, subscriber.FirstName)
	}

	if subscriber.Email != email {
		t.Errorf("Expected email to be '%s', but got '%s'", first, subscriber.FirstName)
	}

	assert.True(t, subscriber.Unsubscribable(), "new subscriber can unsubscribe at the first time, but got false")

	// Check if CreatedAt and UpdatedAt fields are set to the current time
	currentTime := time.Now()
	if subscriber.CreatedAt.After(currentTime) || subscriber.CreatedAt.Before(currentTime.Add(-1*time.Second)) {
		t.Error("Expected CreatedAt to be close to the current time")
	}

	if subscriber.UpdatedAt.After(currentTime) || subscriber.UpdatedAt.Before(currentTime.Add(-1*time.Second)) {
		t.Error("Expected UpdatedAt to be close to the current time")
	}
}

func TestNewNotifierEmailSubTag(t *testing.T) {
	var subscriber = uint64(1)
	var tagId = uint64(2)

	subTag := NewNotifierEmailSubTag(subscriber, tagId)
	if subTag == nil {
		t.Error("Expected subTag to be created, but got nil")
	}

	if subTag.EmailSubscriberId != subscriber {
		t.Errorf("Expected subTag id to be '%d', but got '%d'", subscriber, subTag.EmailSubscriberId)
	}

	if subTag.TagId != tagId {
		t.Errorf("Expected tag id to be '%d', but got '%d'", tagId, subTag.TagId)
	}
}

func TestNewNotifierEmailMessage(t *testing.T) {
	recipientEmail := "to@test.co"
	subscriberId := uint64(1)
	sourceType := "campaign"
	fromEmail := "from@test.co"
	sourceId := uint64(2)
	fromName := "from"
	subject := "test subject"
	emailServiceId := uint64(3)
	msg := "Message from content"

	message := NewNotifierEmailMessage(recipientEmail, subscriberId, sourceType, fromEmail, sourceId, fromName, subject, emailServiceId, msg)
	if message == nil {
		t.Error("Expected message to be created, but got nil")
	}

	if message.RecipientEmail != recipientEmail {
		t.Errorf("Expected recipient email to be '%s', but got '%s'", recipientEmail, message.RecipientEmail)
	}

	if message.SubscriberId != subscriberId {
		t.Errorf("Expected subscriber id to be '%d', but got '%d'", subscriberId, message.SubscriberId)
	}

	if message.SourceType != sourceType {
		t.Errorf("Expected source type to be '%s', but got '%s'", sourceType, message.SourceType)
	}

	if message.FromEmail != fromEmail {
		t.Errorf("Expected from email to be '%s', but got '%s'", fromEmail, message.FromEmail)
	}

	if message.SourceId != sourceId {
		t.Errorf("Expected source id to be '%d', but got '%d'", sourceId, message.SourceId)
	}

	if message.FromName != fromName {
		t.Errorf("Expected from name to be '%s', but got '%s'", fromName, message.FromName)
	}

	if message.Subject != subject {
		t.Errorf("Expected subject to be '%s', but got '%s'", subject, message.Subject)
	}

	if message.EmailServiceId != emailServiceId {
		t.Errorf("Expected email service id to be '%d', but got '%d'", emailServiceId, message.EmailServiceId)
	}

	if message.Message != msg {
		t.Errorf("Expected message to be '%s', but got '%s'", msg, message.Message)
	}

	currentTime := time.Now()
	if message.CreatedAt.After(currentTime) || message.CreatedAt.Before(currentTime.Add(-1*time.Second)) {
		t.Error("Expected CreatedAt to be close to the current time")
	}

	if message.UpdatedAt.After(currentTime) || message.UpdatedAt.Before(currentTime.Add(-1*time.Second)) {
		t.Error("Expected UpdatedAt to be close to the current time")
	}

	if message.FailedAt != nil {
		t.Error("Expected FailedAt to be null")
	}

	if message.SentAt != nil {
		t.Error("Expected SentAt to be null")
	}

	if message.QueuedAt.After(currentTime) || message.QueuedAt.Before(currentTime.Add(-1*time.Second)) {
		t.Error("Expected QueuedAt to be close to the current time")
	}
}

//Mobile tests

func TestNewNotifierMobileUnsubscribeEvent(t *testing.T) {
	reason := "test reason"
	var ID = uint64(1)

	event := NewNotifierMobileUnsubscribeEvent(reason, ID)
	if event == nil {
		t.Error("Expected event to be created, but got nil")
	}

	if event.Reason != reason {
		t.Errorf("Expected reason to be '%s', but got '%s'", reason, event.Reason)
	}

	if event.ID != ID {
		t.Errorf("Expected id to be '%d', but got '%d'", ID, event.ID)
	}
}

func TestNewNotifierNewNotifierMobileService(t *testing.T) {
	config := map[string]interface{}{
		"Url":   "testHost",
		"Token": "username",
	}
	payload, _ := json.Marshal(config)
	sType := NotifierMobileServiceKavehNegarType
	name := "Test SMTP"

	// Call the function to create a new NewNotifierEmailService
	service := NewNotifierMobileDriver(string(payload), sType, name)

	// Check if the service is not nil
	if service == nil {
		t.Error("Expected service to be created, but got nil")
	}

	// Check if the name are set correctly
	if service.Name != name {
		t.Errorf("Expected name to be '%s', but got '%s'", name, service.Name)
	}

	if service.Payload != string(payload) {
		t.Errorf("Expected payload to be '%s', but got '%s'", payload, service.Payload)
	}

	if service.Type != NotifierMobileServiceKavehNegarType {
		t.Errorf("Expected type to be '%s', but got '%s'", sType, service.Type)
	}

	var tmp SmtpConfig
	err := json.Unmarshal([]byte(service.Payload), &tmp)
	if err != nil {
		t.Errorf("Error during unmarshal payload '%s'", err)
	}
}

func TestNewNotifierMobileSubscriber(t *testing.T) {
	mobile := "9194567890"
	countryCode := "+98"
	first := "first"
	last := "last"

	subscriber := NewNotifierMobileSubscriber(countryCode, mobile, first, last)
	if subscriber == nil {
		t.Error("Expected subscriber to be created, but got nil")
	}

	if subscriber.FirstName != first {
		t.Errorf("Expected first to be '%s', but got '%s'", first, subscriber.FirstName)
	}

	if subscriber.LastName != last {
		t.Errorf("Expected last to be '%s', but got '%s'", last, subscriber.LastName)
	}

	if subscriber.Mobile != mobile {
		t.Errorf("Expected mobile to be '%s', but got '%s'", mobile, subscriber.Mobile)
	}

	if subscriber.CountryCode != countryCode {
		t.Errorf("Expected country code to be '%s', but got '%s'", mobile, subscriber.Mobile)
	}

	assert.True(t, subscriber.Unsubscribable(), "new subscriber can unsubscribe at the first time, but got false")

	// Check if CreatedAt and UpdatedAt fields are set to the current time
	currentTime := time.Now()
	if subscriber.CreatedAt.After(currentTime) || subscriber.CreatedAt.Before(currentTime.Add(-1*time.Second)) {
		t.Error("Expected CreatedAt to be close to the current time")
	}

	if subscriber.UpdatedAt.After(currentTime) || subscriber.UpdatedAt.Before(currentTime.Add(-1*time.Second)) {
		t.Error("Expected UpdatedAt to be close to the current time")
	}
}

func TestNewNotifierMobileSubTag(t *testing.T) {
	var subscriber = uint64(1)
	var tagId = uint64(2)

	subTag := NewNotifierMobileSubTag(subscriber, tagId)
	if subTag == nil {
		t.Error("Expected subTag to be created, but got nil")
	}

	if subTag.MobileSubscriberId != subscriber {
		t.Errorf("Expected subTag id to be '%d', but got '%d'", subscriber, subTag.MobileSubscriberId)
	}

	if subTag.TagId != tagId {
		t.Errorf("Expected tag id to be '%d', but got '%d'", tagId, subTag.TagId)
	}
}

//Notification tests

func TestNewNotifierNewNotifierNotificationService(t *testing.T) {
	config := map[string]interface{}{
		"Url":           "testHost",
		"ApplicationId": "username",
	}
	payload, _ := json.Marshal(config)
	sType := NotifierNotificationServiceFirebaseType
	name := "Test SMTP"

	// Call the function to create a new NewNotifierEmailService
	service := NewNotifierNotificationService(string(payload), sType, name)

	// Check if the service is not nil
	if service == nil {
		t.Error("Expected service to be created, but got nil")
	}

	// Check if the name are set correctly
	if service.Name != name {
		t.Errorf("Expected name to be '%s', but got '%s'", name, service.Name)
	}

	if service.Payload != string(payload) {
		t.Errorf("Expected payload to be '%s', but got '%s'", payload, service.Payload)
	}

	if service.Type != NotifierNotificationServiceFirebaseType {
		t.Errorf("Expected type to be '%s', but got '%s'", sType, service.Type)
	}

	var tmp SmtpConfig
	err := json.Unmarshal([]byte(service.Payload), &tmp)
	if err != nil {
		t.Errorf("Error during unmarshal payload '%s'", err)
	}
}

func TestNewNotifierNotificationSubscriber(t *testing.T) {
	token := "token fcm or gcm"
	first := "first"
	last := "last"
	driverId := uint64(1)

	subscriber := NewNotifierNotificationSubscriber(token, first, last, driverId)
	if subscriber == nil {
		t.Error("Expected subscriber to be created, but got nil")
	}

	if subscriber.FirstName != first {
		t.Errorf("Expected first to be '%s', but got '%s'", first, subscriber.FirstName)
	}

	if subscriber.LastName != last {
		t.Errorf("Expected last to be '%s', but got '%s'", last, subscriber.LastName)
	}

	if subscriber.Token != token {
		t.Errorf("Expected mobile to be '%s', but got '%s'", token, subscriber.Token)
	}

	if subscriber.DriverId != driverId {
		t.Errorf("Expected country code to be '%d', but got '%d'", driverId, subscriber.DriverId)
	}

	// Check if CreatedAt and UpdatedAt fields are set to the current time
	currentTime := time.Now()
	if subscriber.CreatedAt.After(currentTime) || subscriber.CreatedAt.Before(currentTime.Add(-1*time.Second)) {
		t.Error("Expected CreatedAt to be close to the current time")
	}

	if subscriber.UpdatedAt.After(currentTime) || subscriber.UpdatedAt.Before(currentTime.Add(-1*time.Second)) {
		t.Error("Expected UpdatedAt to be close to the current time")
	}
}

func TestNewNotifierNotificationSubTag(t *testing.T) {
	var subscriber = uint64(1)
	var tagId = uint64(2)

	subTag := NewNotifierNotificationSubTag(subscriber, tagId)
	if subTag == nil {
		t.Error("Expected subTag to be created, but got nil")
	}

	if subTag.NotificationSubscriberId != subscriber {
		t.Errorf("Expected sub id to be '%d', but got '%d'", subscriber, subTag.NotificationSubscriberId)
	}

	if subTag.TagId != tagId {
		t.Errorf("Expected tag id to be '%d', but got '%d'", tagId, subTag.TagId)
	}
}

//Tag test

func TestNewNotifierTag(t *testing.T) {
	name := "all"

	tag := NewNotifierTag(name)
	if tag == nil {
		t.Error("Expected subTag to be created, but got nil")
	}

	if tag.Name != name {
		t.Errorf("Expected name to be '%s', but got '%s'", name, tag.Name)
	}

	// Check if CreatedAt and UpdatedAt fields are set to the current time
	currentTime := time.Now()
	if tag.CreatedAt.After(currentTime) || tag.CreatedAt.Before(currentTime.Add(-1*time.Second)) {
		t.Error("Expected CreatedAt to be close to the current time")
	}

	if tag.UpdatedAt.After(currentTime) || tag.UpdatedAt.Before(currentTime.Add(-1*time.Second)) {
		t.Error("Expected UpdatedAt to be close to the current time")
	}
}
