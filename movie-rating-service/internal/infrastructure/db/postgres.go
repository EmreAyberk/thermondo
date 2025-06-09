package db

import (
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"movie-rating-service/config"
	"movie-rating-service/internal/domain"
	"time"
)

/*
No interface is needed because:
- GORM already provides all required methods and abstractions.
- Only mock at the service or domain layer for tests.
- Avoids unnecessary interface pollution.
*/

var db *gorm.DB

func BeginTransaction() *gorm.DB {
	return db.Begin()
}

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.Cfg.DbConfig.User,
		config.Cfg.DbConfig.Password,
		config.Cfg.DbConfig.Host,
		config.Cfg.DbConfig.Port,
		config.Cfg.DbConfig.Name,
		config.Cfg.DbConfig.SSLMode,
	)
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
	if err != nil {

		return nil, fmt.Errorf("error occurred while connecting to database: %w", err)
	}

	if config.Cfg.DebugMode {
		dbConn = dbConn.Debug()
	}

	sqlDB, err := dbConn.DB()
	if err != nil {
		return nil, fmt.Errorf("error occurred while connecting to database: %w", err)
	}

	sqlDB.SetConnMaxLifetime(time.Minute * 5)
	sqlDB.SetConnMaxIdleTime(time.Minute * 5)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)

	if cast.ToBool(config.Cfg.Migrate) {
		err = migrate(dbConn)
		if err != nil {
			return nil, fmt.Errorf("error occurred while migrating database: %w", err)
		}
	}

	db = dbConn
	return dbConn, nil
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&domain.User{}, &domain.Movie{}, &domain.Rating{})
}
