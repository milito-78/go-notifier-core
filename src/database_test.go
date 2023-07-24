package go_notifier_core

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestGetDriver(t *testing.T) {
	config := DbConfig{
		Name:     "name",
		Username: "roottest",
		Password: "secret",
		Driver:   MysqlDriver,
		Host:     "localhost",
		Port:     "3306",
		DB:       "notifier_test",
	}

	// Test when the driver is already set
	driver := getDriver(config)
	assert.NotNil(t, driver, "Driver should not be nil")

	// Test invalid driver type
	invalidConfig := DbConfig{
		Name:   "invalidDB",
		Driver: 999, // Invalid driver type
	}
	assert.Panics(t, func() {
		getDriver(invalidConfig)
	}, "Getting driver with an invalid driver type should panic")

	// Test invalid driver type
	postgresConfig := DbConfig{
		Name:   "invalidDB",
		Driver: PostgresDriver,
	}

	assert.Panics(t, func() {
		getDriver(postgresConfig)
	}, "Panic for postgres driver, not implemented")
}

func TestSetDriver(t *testing.T) {
	// Create a sample DbConfig
	config := DbConfig{
		Name:     "name",
		Username: "roottest",
		Password: "secret",
		Driver:   MysqlDriver,
		Host:     "localhost",
		Port:     "3306",
		DB:       "notifier_test",
	}

	// Test setting the driver for MySQL
	driver := setDriver(config)
	assert.NotNil(t, driver, "Driver should not be nil")

	// Test setting the driver for an invalid driver type
	invalidConfig := DbConfig{
		Name:   "invalidDB",
		Driver: 999, // Invalid driver type
	}
	assert.Panics(t, func() {
		setDriver(invalidConfig)
	}, "Setting driver with an invalid driver type should panic")

	// Test invalid driver type
	postgresConfig := DbConfig{
		Name:   "invalidDB",
		Driver: PostgresDriver,
	}

	assert.Panics(t, func() {
		setDriver(postgresConfig)
	}, "Panic for postgres driver, not implemented")
}

func TestMysqlDriver(t *testing.T) {
	// Create a sample DbConfig
	config := DbConfig{
		Name:     "name",
		Username: "roottest",
		Password: "secret",
		Driver:   MysqlDriver,
		Host:     "localhost",
		Port:     "3306",
		DB:       "notifier_test",
	}

	// Test connecting to MySQL
	db := mysqlDriver(config)
	assert.NotNil(t, db, "Database connection should not be nil")
	// Further testing on the db object and interactions with the MySQL database can be done here.
}

func TestConcurrentGetDriver(t *testing.T) {
	// Create a sample DbConfig for MySQL driver
	config1 := DbConfig{
		Name:     "dbConfig1",
		Driver:   MysqlDriver,
		Host:     "localhost",
		Port:     "3306",
		Username: "roottest",
		Password: "secret",
		DB:       "notifier_test",
	}

	// Create another sample DbConfig for MySQL driver
	config2 := DbConfig{
		Name:     "dbConfig2",
		Driver:   MysqlDriver,
		Host:     "localhost",
		Port:     "3306",
		Username: "roottest",
		Password: "secret",
		DB:       "notifier_test",
	}

	// Ensure that the drivers map is initially empty
	assert.Empty(t, drivers, "drivers map should be empty initially")

	// Create a WaitGroup to synchronize concurrent goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Simultaneously get drivers for both configurations using concurrent goroutines
	go func() {
		defer wg.Done()
		_ = getDriver(config1)
	}()

	go func() {
		defer wg.Done()
		_ = getDriver(config2)
	}()

	// Wait for both goroutines to finish
	wg.Wait()

	// Ensure that both drivers are created and stored in the drivers map
	assert.Equal(t, 2, len(drivers), "drivers map should have two entries")

	// Ensure that drivers for both configurations are correctly stored in the drivers map
	assert.NotNil(t, drivers[config1.Name], "Driver for config1 should not be nil")
	assert.NotNil(t, drivers[config2.Name], "Driver for config2 should not be nil")
}
