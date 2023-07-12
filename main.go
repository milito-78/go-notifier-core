package main

import (
	"go-notifier-core/src"
)

func main() {
	c := src.DbConfig{
		Username: "root",
		Password: "secret",
		Driver:   src.MysqlDriver,
		Host:     "127.0.0.1",
		Port:     "3306",
		DB:       "notifier",
	}

	src.Migrate(c)
	//src.Initialize(c)
	//_, err := src.CreateTag("all")
	//if err != nil {
	//	log.Printf("Error ! %s", err)
	//}
	//_, _ = src.CreateTag("newsletter")
	//
	//tags, _ := src.TagsList()
	//for i, tag := range tags {
	//	fmt.Printf("I : %d , %+v\n", i, tag)
	//}
	//
	//subs, err := src.SubscribeEmail("milad.test@gmail.com", "milad", "lname", []string{"all", "temp"}, true)
	//fmt.Println(subs)
	//
	//err = src.AssignTagsToEmail("milad.test@gmail.com", []string{"newsletter"}, false)
	//if err != nil {
	//	fmt.Println(err)
	//}
}
