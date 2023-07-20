package go_notifier_core

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

const (
	MysqlDriver = iota + 1
	PostgresDriver
)

var drivers map[string]interface{}

func init() {
	drivers = make(map[string]interface{}, 2)
}

func getDriver(config DbConfig) interface{} {
	if val, ok := drivers[config.Name]; ok {
		return val
	}
	return setDriver(config)
}

func setDriver(config DbConfig) interface{} {
	switch config.Driver {
	case MysqlDriver:
		tmp := mysqlDriver(config)
		drivers[config.Name] = tmp
		return tmp
	}
	panic("invalid driver for config '" + config.Name + "'.")
}

func mysqlDriver(config DbConfig) *gorm.DB {
	if config.Password != "" {
		config.Password = ":" + config.Password
	}
	dsn := config.Username + config.Password +
		"@tcp(" + config.Host + ":" + config.Port + ")/" +
		config.DB + "?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error during connecting db mysql driver : %s", err)
	}
	return db
}
