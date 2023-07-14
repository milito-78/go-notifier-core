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

	go_notifier_core.Migrate(c)
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

	err = go_notifier_core.AssignTagsToEmail("milad.test@gmail.com", []string{"newsletter"}, false)
	if err != nil {
		fmt.Println(err)
	}

	//go_notifier_core.MigrateRollback(c)

}
