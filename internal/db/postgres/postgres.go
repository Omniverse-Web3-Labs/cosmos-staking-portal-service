package postgres

import (
	"app/internal/config"
	"app/internal/db/postgres/query"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	"gorm.io/plugin/opentelemetry/tracing"
)

var defaultConnect *gorm.DB
var subqueryConnect *gorm.DB

func Load() {
	defaultConnect = NewDB(config.GetConfig().Databases.Postgres.Default)
	query.SetDefault(defaultConnect)
	subqueryConnect = NewDB(config.GetConfig().Databases.Postgres.Subquery)
}

func Default() *gorm.DB {
	return defaultConnect
}

func Subquery() *gorm.DB {
	return subqueryConnect
}

func Exists(err error) (bool, error) {
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

func NewDB(dbConfig config.Database) *gorm.DB {
	logger := logger.New(
		logrus.NewWriter(),
		logger.Config{
			SlowThreshold: time.Millisecond,
			LogLevel:      logger.Warn,
			Colorful:      false,
		},
	)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		dbConfig.Host,
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.Port,
	)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger: logger,
	})

	if err != nil {
		panic("failed to connect database")
	}
	if err := db.Use(tracing.NewPlugin(
		tracing.WithoutServerAddress(),
	)); err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get db instance")
	}
	sqlDB.SetMaxIdleConns(2000) // 设置最大空闲连接数
	sqlDB.SetMaxOpenConns(2000) // 设置最大打开连接数
	return db
}
