package go_notifier_core

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
)

type DbConfig struct {
	Username string
	Password string
	Driver   int
	Host     string
	Port     string
	Name     string
	DB       string
}

const (
	MysqlDriver    = iota + 1 // Constant representing the MySQL database driver.
	PostgresDriver            // Constant representing the PostgresSQL database driver.
)

var (
	driversMu sync.Mutex
	drivers   map[string]interface{} // Map to store database drivers associated with their configurations.
)

func init() {
	drivers = make(map[string]interface{}, 2)
}

// getDriver returns the database driver instance based on the provided configuration.
// If the driver instance is already created for the given configuration, it returns the existing one.
// If the driver instance does not exist, it calls the setDriver function to create a new one.
// The function ensures that only one driver instance is created and reused per unique configuration.
// If an invalid driver type is provided in the configuration, the function panics.
func getDriver(config DbConfig) interface{} {
	driversMu.Lock()
	defer driversMu.Unlock()

	if val, ok := drivers[config.Name]; ok {
		return val
	}
	return setDriver(config)
}

// setDriver creates and returns a new database driver instance based on the provided configuration.
// The function supports MySQL driver. If the configuration contains an invalid driver type,
// the function panics with an error message.
func setDriver(config DbConfig) interface{} {
	switch config.Driver {
	case MysqlDriver:
		tmp := mysqlDriver(config)
		drivers[config.Name] = tmp
		return tmp
	case PostgresDriver:
		panic("postgres driver doesn't implemented")
	}
	panic("invalid driver for config '" + config.Name + "'.")
}

// mysqlDriver establishes a connection to the MySQL database using the provided configuration.
// It returns a *gorm.DB object, which represents the database connection.
// If any error occurs during the connection, the function logs a fatal error and exits the application.
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
