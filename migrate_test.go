package go_notifier_core

import (
	"gorm.io/gorm"
	"testing"
)

type mockMigration struct {
	upCalled   bool
	downCalled bool
}

func (m *mockMigration) Up() error {
	m.upCalled = true
	return nil
}

func (m *mockMigration) Down() error {
	m.downCalled = true
	return nil
}

type mockMigrator struct {
	db         gorm.Migrator
	migrations []*mockMigration
}

func (m *mockMigrator) rollback() error {
	mgs := m.migrations
	for i := len(mgs) - 1; i >= 0; i-- {
		migration := mgs[i]
		err := migration.Down()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *mockMigrator) migrate() error {
	mgs := m.migrations
	for _, migration := range mgs {
		err := migration.Up()
		if err != nil {
			return err
		}
	}
	return nil
}

func TestMockMigrate(t *testing.T) {
	mockMigration1 := mockMigration{}
	mockMigration2 := mockMigration{}
	migrations := []*mockMigration{&mockMigration1, &mockMigration2}
	mockMigrator := &mockMigrator{migrations: migrations}
	err := mockMigrator.migrate()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	if !mockMigration1.upCalled {
		t.Error("Expected Up() to be called on the first migration")
	}

	if !mockMigration2.upCalled {
		t.Error("Expected Up() to be called on the second migration")
	}
}

func TestMockRollback(t *testing.T) {
	mockMigration1 := &mockMigration{}
	mockMigration2 := &mockMigration{}
	migrations := []*mockMigration{mockMigration1, mockMigration2}
	mockMigrator := &mockMigrator{migrations: migrations}
	err := mockMigrator.rollback()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	if !mockMigration1.downCalled {
		t.Error("Expected Down() to be called on the first migration")
	}

	if !mockMigration2.downCalled {
		t.Error("Expected Down() to be called on the second migration")
	}
}
