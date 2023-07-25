package main

import "github.com/milito-78/go-notifier-core"

func main() {
	//For connecting app to db you don't need to create env.
	c := go_notifier_core.DbConfig{
		Name:     "connection name",
		Username: "root",
		Password: "secret",
		Driver:   go_notifier_core.MysqlDriver, // Postgres doesn't implemented
		Host:     "127.0.0.1",
		Port:     "3306",
		DB:       "notifier",
	}

	//To migrate and create tables use this function.
	go_notifier_core.Migrate(c)

	//##########################################

	//You need to init and seed your tables.
	// BUT! Before you seed your tables you should initialize your application.
	//So you should use Initialize function first.
	go_notifier_core.Initialize(c)
	go_notifier_core.Seed()

	/*
		 * HINT: If you used Initialize() you can use handlers and Seed function and Migrate() function too.
			But why Migrate() function has a config type input, cause if you want to migrate your tables without any
			initializing handlers and repositories. You can use it easily.
	*/
	//##########################################

	//To rollback your migrations and remove your db use this function
	go_notifier_core.MigrateRollback(c)
}
