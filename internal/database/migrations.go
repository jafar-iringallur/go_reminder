package database

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	migrator, err := newMigrator(db)
	if err != nil {
		return err
	}

	err = migrator.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

func RollbackLastMigration(db *gorm.DB) error {
	migrator, err := newMigrator(db)
	if err != nil {
		return err
	}

	err = migrator.Steps(-1)
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

func newMigrator(db *gorm.DB) (*migrate.Migrate, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return nil, err
	}

	return migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
}
