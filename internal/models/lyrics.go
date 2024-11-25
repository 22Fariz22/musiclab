package models

import "time"

// SongRequest для добавления песни
type SongRequest struct {
	Group string `json:"group" validate:"required,min=1"`
	Song  string `json:"song" validate:"required,min=1"`
}

// SongDetail для ответа
type SongDetail struct {
	ReleaseDate string `json:"releaseDate" validate:"required"`
	Text        string `json:"text" validate:"required"`
	Link        string `json:"link" validate:"required"`
}

// UpdateTrackRequest обновление информации
type UpdateTrackRequest struct {
	ID          uint   	`json:"id" validate:"required"`
	GroupName   string  `json:"group" validate:"required"`       
	SongName    string  `json:"song" validate:"required"`         
	ReleaseDate *string `json:"release_date,omitempty"`          
	Text        *string `json:"text,omitempty"`                  
	Link        *string `json:"link,omitempty"`                 
}

// Song модель базы данных
type Song struct {
	ID          uint   `gorm:"primaryKey"`
	GroupName   string `gorm:"type:varchar(255);not null"`
	SongName    string `gorm:"type:varchar(255);not null"`
	ReleaseDate *time.Time
	Text        string    `gorm:"type:text"` 
	Link        *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}