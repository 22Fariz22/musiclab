package migrate

import (
	"github.com/22Fariz22/musiclab/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Migrate applies database migrations
func Migrate(dsn string) error {
	// Инициализация GORM с использованием только для миграций
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Выполнение миграций
	return db.AutoMigrate(&models.Song{})
}
