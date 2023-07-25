package main

import (
	"github.com/milito-78/go-notifier-core"
	"time"
)

func main() {
	c := go_notifier_core.DbConfig{
		Name:     "connection name",
		Username: "root",
		Password: "secret",
		Driver:   go_notifier_core.MysqlDriver, // Postgres doesn't implemented
		Host:     "127.0.0.1",
		Port:     "3306",
		DB:       "notifier",
	}
	//For use worker you need to Initialize it first.
	go_notifier_core.Initialize(c)

	//You need to create a list of workers that you need.
	list := go_notifier_core.WorkersList{
		go_notifier_core.WorkerConfig{
			Duration: time.Second * 10,
			Worker:   go_notifier_core.EmailWorker{}, // You can customize your worker. It must be implemented from IWorker interface.
			Name:     "Email worker",
		},
	}

	//After create list, you should pass list to start it.
	go_notifier_core.WorkerStart(list)

	//To keep your app running you can use chanel
	var infinity chan interface{}
	<-infinity
}
