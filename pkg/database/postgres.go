package database

import (
	"fmt"
	"log"
	"x-clone/internal/config"
	"x-clone/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.DBName,
		cfg.Database.Password,
		cfg.Database.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database")
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Follower{},
		&model.Repost{},
		&model.Like{},
	)
	if err != nil {
		log.Fatalf("failed to apply migrations")
		return nil, err
	}

	return db, nil
}
