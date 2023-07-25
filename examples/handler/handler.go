package main

import (
	"encoding/json"
	"github.com/milito-78/go-notifier-core"
	"log"
	"strconv"
)

func main() {

	c := go_notifier_core.DbConfig{
		Name:     "connection name",
		Username: "root",
		Password: "secret",
		Driver:   go_notifier_core.MysqlDriver,
		Host:     "127.0.0.1",
		Port:     "3306",
		DB:       "notifier",
	}

	//For using package you should first initialize it.
	go_notifier_core.Initialize(c)

	//Create a tag called all. all tag is a default tag for subscribers. All subscribers have this tag.
	tag, err := go_notifier_core.CreateTag("all")
	if err != nil {
		log.Printf("Error during create a tag : %s", err)
	}

	//subscribe 1000 emails
	for i := 0; i < 1000; i++ {
		iString := strconv.Itoa(i)
		// Subscribe email with `all` tag.
		// If you want to assign subscriber a tag that's not created before, and you want to create that tag,
		// you should pass the last param true. This one allows your application that create new tags inside list.
		// For example :
		// _, err = go_notifier_core.SubscribeEmail("fname"+iString + ".lname"+iString + "@test.com"+iString, "fname"+iString, "lname"+iString, []string{"new tag","another tag"}, false)
		// If the last param equals to false, the function returns error when your tags not exists in database.
		// If the last param equals to true, it stores your new tags in db and assign them to subscriber
		_, err = go_notifier_core.SubscribeEmail("fname"+iString+".lname"+iString+"@test.com"+iString, "fname"+iString, "lname"+iString, []string{}, false)
	}

	//Create a mail service (SMTP)
	config := go_notifier_core.SmtpConfig{
		Host:       "sandbox.smtp.mailtrap.io",
		Port:       "2525",
		Username:   "username",
		Password:   "password",
		Encryption: "tls",
	}

	bt, err := json.Marshal(&config)
	if err != nil {
		log.Printf("Error during json config : %s", err)
	}

	emailService, err := go_notifier_core.CreateEmailService("smtp test", go_notifier_core.NotifierEmailServiceSMTPType, bt)

	//Create a email template
	template, error := go_notifier_core.CreateEmailTemplate("Test Template", "<h1> Hello, Good morning</h1></br><p>This is a test</p>")

	//Create a campaign
	//To send a scheduled email for subscribers, you should create a campaign
	//You can schedule your campaign to send in a selected time.
	//You can set tags to send a group of subscribers or send to all of them with `all` tag.
	//For a campaign, You should select email service, template, from email, from name, and tags list.
	campaign, err := go_notifier_core.AddEmailCampaign(&go_notifier_core.EmailCampaignCreateData{
		EmailServiceId: emailService.ID,
		TemplateId:     template.ID,
		StatusId:       go_notifier_core.NotifierEmailStatusDraft,
		FromEmail:      "go.notifier.service@mail.com",
		FromName:       "Go Notifier",
		Subject:        "This is a test",
		Name:           "New Campaign",
		Tags:           []uint64{1},
	})

}
