package models

import "time"

// SongRequest для добавления песни
type SongRequest struct {
	Group string `json:"group" validate:"required,min=1"`
	Song  string `json:"song" validate:"required,min=1"`
}

// SongDetail для ответа
type SongDetail struct {
	ReleaseDate time.Time `json:"releaseDate" validate:"required"`
	Text        string    `json:"text" validate:"required"`
	Link        string    `json:"link" validate:"required"`
}

// UpdateTrackRequest обновление информации
type UpdateTrackRequest struct {
	ID          uint    `json:"id" validate:"required"`
	GroupName   string  `json:"group" validate:"required"`
	SongName    string  `json:"song" validate:"required"`
	ReleaseDate *string `json:"release_date,omitempty"`
	Text        *string `json:"text,omitempty"`
	Link        *string `json:"link,omitempty"`
}

// Group модель базы данных
type Group struct {
	ID        uint      `gorm:"primaryKey" db:"id"`
	Name      string    `gorm:"type:varchar(255);not null;unique;index" db:"name"`
	CreatedAt time.Time `gorm:"index" db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Song модель базы данных
type Song struct {
	ID          uint      `gorm:"primaryKey" db:"id"`
	GroupID     uint      `gorm:"not null;index;uniqueIndex:idx_group_song,priority:1" db:"group_id"`
	Group       Group     `gorm:"foreignKey:GroupID"`
	SongName    string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_group_song,priority:2;index" db:"song_name"` // добавляем отдельный индекс, если часто ищем по имени
	ReleaseDate time.Time `gorm:"index" db:"release_date"`
	Text        string    `gorm:"type:text" db:"text"`
	Link        *string   `gorm:"index" db:"link"`
	CreatedAt   time.Time `gorm:"index" db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
