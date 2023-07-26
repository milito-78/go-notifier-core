package go_notifier_core

import (
	"github.com/milito-78/go-notifier-core/migrations"
	"gorm.io/gorm"
	"log"
)

type migrator interface {
	migrate() error
}

type migratorRollback interface {
	rollback() error
}

type gormMigrator struct {
	db gorm.Migrator
}

func (g gormMigrator) rollback() error {
	mgs := migrations.GetMigrationsList(g.db)
	for i := len(mgs) - 1; i >= 0; i-- {
		migration := mgs[i]
		err := migration.Down()
		if err != nil {
			return err
		}
	}
	return nil
}

func (g gormMigrator) migrate() error {
	mgs := migrations.GetMigrationsList(g.db)
	for _, migration := range mgs {
		err := migration.Up()
		if err != nil {
			return err
		}
	}
	return nil
}

func Migrate(config DbConfig) {
	m := driverFactory(config)
	err := m.migrate()
	if err != nil {
		log.Fatalf("Error during migrate : %s\n", err)
	} else {
		log.Println("Migration runs successfully")
	}
}

func MigrateRollback(config DbConfig) {
	m := rollbackDriverFactory(config)
	err := m.rollback()
	if err != nil {
		log.Fatalf("Error during rollback : %s\n", err)
	} else {
		log.Println("Migration rollback runs successfully")
	}
}

func driverFactory(config DbConfig) migrator {
	db := getDriver(config)
	return gormMigrator{db: db.(*gorm.DB).Migrator()}
}

func rollbackDriverFactory(config DbConfig) migratorRollback {
	db := getDriver(config)
	return gormMigrator{db: db.(*gorm.DB).Migrator()}
}
