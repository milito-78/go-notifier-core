package main

import (
	"fmt"
	"go-notifier-core/src"
	"log"
)

func main() {
	c := go_notifier_core.DbConfig{
		Username: "root",
		Password: "secret",
		Driver:   go_notifier_core.MysqlDriver,
		Host:     "127.0.0.1",
		Port:     "3306",
		DB:       "notifier",
	}
	//go_notifier_core.MigrateRollback(c)
	//go_notifier_core.Migrate(c)
	//return

	go_notifier_core.Initialize(c)
	_, err := go_notifier_core.CreateTag("all")
	if err != nil {
		log.Printf("Error ! %s", err)
	}
	_, _ = go_notifier_core.CreateTag("newsletter")

	tags, _ := go_notifier_core.TagsList()
	for i, tag := range tags {
		fmt.Printf("I : %d , %+v\n", i, tag)
	}

	subs, err := go_notifier_core.SubscribeEmail("milad.test@gmail.com", "milad", "lname", []string{"all", "temp"}, true)
	fmt.Println(subs)

	err = go_notifier_core.AssignTagsToEmail("milad.test@gmail.com", []string{"newsletter", "xxx", "aa"}, false)
	if err != nil {
		fmt.Println(err)
	}

	err = go_notifier_core.RemoveTagsFromEmail("milad.test@gmail.com", []string{"all", "xxx", "aa"})
	if err != nil {
		fmt.Println(err)
	}

	err = go_notifier_core.DeleteTagByName("all")
	if err != nil {
		fmt.Println(err)
	}

	err = go_notifier_core.UnSubscribeEmail("milad.test@gmail.com", 1)
	if err != nil {
		fmt.Println(err)
	}

	events, _ := go_notifier_core.EmailUnsubscribeEventsList()
	for i, tag := range events {
		fmt.Printf("I : %d , %+v\n", i, tag)
	}

	_, _ = go_notifier_core.SubscribeEmail("milad.test22@gmail.com", "milad", "lname", []string{"all", "temp"}, true)

	subscribers, _ := go_notifier_core.GetTagEmailSubscribers("all")
	for i, tag := range subscribers {
		fmt.Printf("I : %d , %+v\n", i, tag)
	}

	unSubscribers, _ := go_notifier_core.GetUnsubscribedEmails()
	for i, tag := range unSubscribers {
		fmt.Printf("I : %d , %+v\n", i, tag)
	}

	template, _ := go_notifier_core.CreateEmailTemplate("new template", "Hello, Good morning")
	if err != nil {
		log.Printf("Error %s", err)
	} else {
		log.Println(template)
	}

	tem, _ := go_notifier_core.UpdateEmailTemplate(template.ID, "xxxx", "sasas")
	if err != nil {
		log.Printf("Error %s", err)
	} else {
		log.Println(tem)
	}

	err = go_notifier_core.DeleteEmailTemplate(1)
	if err != nil {
		log.Printf("Error %s", err)
	}

	templates, err := go_notifier_core.EmailTemplateList()
	if err != nil {
		log.Printf("Error %s", err)
	}
	for i, tag := range templates {
		fmt.Printf("I : %d , %+v\n", i, tag)
	}

	campaign, err := go_notifier_core.AddEmailCampaign(&go_notifier_core.EmailCampaignCreateData{
		EmailServiceId: 1,
		TemplateId:     3,
		StatusId:       1,
		FromEmail:      "hello@mail.com",
		FromName:       "My mail",
		Subject:        "This is a test",
		Name:           "New Campaign common",
		Tags:           []uint64{1, 2, 3},
	})
	if err != nil {
		log.Printf("Error %s", err)
	} else {
		log.Println(campaign)
	}

	err = go_notifier_core.UpdateEmailCampaignWithId(2, &go_notifier_core.EmailCampaignUpdateData{
		EmailServiceId: 1,
		TemplateId:     3,
		StatusId:       2,
		FromEmail:      "hello@mail.com",
		FromName:       "My mail",
		Subject:        "This is a test",
		Name:           "New Campaign common",
		Tags:           []uint64{1, 2, 3},
	})
	if err != nil {
		log.Printf("Error %s", err)
	}

	err = go_notifier_core.DeleteEmailCampaign(1)
	if err != nil {
		log.Printf("Error %s", err)
	}

	services, _ := go_notifier_core.GetEmailServices()
	for i, service := range services {
		fmt.Printf("I : %d , %+v\n", i, service)
	}

	//go_notifier_core.MigrateRollback(c)

}
