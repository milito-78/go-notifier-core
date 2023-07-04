package go_notifier_core

import (
	"go-notifier-core/domains"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

const (
	MysqlDriver = iota + 1
	PostgresDriver
)

type migrator interface {
	migrate() error
}

type gormMigrator struct {
	db gorm.Migrator
}

func (g gormMigrator) migrate() error {
	//TODO change migrations
	return g.db.AutoMigrate(
		domains.NotifierTag{},
		domains.NotifierEmailUnsubscribeEvent{},
		domains.NotifierEmailSubscriber{},
		domains.NotifierEmailSubTag{},
		domains.NotifierMobileUnsubscribeEvent{},
		domains.NotifierMobileSubscriber{},
		domains.NotifierMobileSubTag{},
		domains.NotifierNotificationDriver{},
		domains.NotifierNotificationSubscriber{},
		domains.NotifierNotificationSubTag{},
		domains.NotifierEmailCampaignTemplate{},
		domains.NotifierEmailService{},
		domains.NotifierEmailStatus{},
		domains.NotifierEmailCampaign{},
		domains.NotifierEmailCampaignTag{},
		domains.NotifierEmailMessage{},
	)
}

type DbConfig struct {
	Username string
	Password string
	Driver   int
	Host     string
	Port     string
	DB       string
}

func Migrate(config DbConfig) {
	m := driverFactory(config)
	err := m.migrate()
	if err != nil {
		log.Fatalf("Error during migrate : %s", err)
	}
}

func driverFactory(config DbConfig) migrator {
	switch config.Driver {
	case MysqlDriver:
		return mysqlDriverMigrator(config)
	}
	return mysqlDriverMigrator(config)
}

func mysqlDriverMigrator(config DbConfig) migrator {
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

	return &gormMigrator{db: db.Migrator()}
}
