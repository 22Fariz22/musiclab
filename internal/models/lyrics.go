package models

import "time"

// SongRequest для добавления песни
type SongRequest struct {
	Group string `json:"group" validate:"required"`
	Song  string `json:"song" validate:"required"`
}

// SongDetail для ответа
type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// Song модель базы данных
type Song struct {
	ID          uint   `gorm:"primaryKey"`
	GroupName   string `gorm:"type:varchar(255);not null"`
	SongName    string `gorm:"type:varchar(255);not null"`
	ReleaseDate *time.Time
	Text        []string `gorm:"type:text[]"`
	Link        *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
