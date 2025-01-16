package models

import "time"

// SongRequest для добавления песни
// @Description Request payload for adding a new song
type SongRequest struct {
	// Group name of the song
	// Required: true
	// Min length: 1
	Group string `json:"group" validate:"required,min=1"`
	// Song name
	// Required: true
	// Min length: 1
	Song string `json:"song" validate:"required,min=1"`
}

// SongDetail для ответа
// @Description Response containing details of a song
type SongDetail struct {
	// Release date of the song
	// Required: true
	ReleaseDate string `json:"releaseDate" validate:"required"`

	// Lyrics or text of the song
	// Required: true
	Text string `json:"text" validate:"required"`

	// External link to the song
	// Required: true
	Link string `json:"link" validate:"required"`
}

// UpdateTrackRequest обновление информации
// @Description Request payload for updating song details
type UpdateTrackRequest struct {
	// ID of the track to update
	// Required: true
	ID uint `json:"id" validate:"required"`

	// Group name
	// Required: true
	// Min length: 1
	GroupName *string `json:"group" validate:"required,min=1"`

	// Song name
	// Required: true
	// Min length: 1
	SongName *string `json:"song" validate:"required,min=1"`

	// Release date
	// Required: true
	ReleaseDate *string `json:"release_date" validate:"required,min=1"`

	// Lyrics or text of the song
	Text *string `json:"text,omitempty"`

	// External link to the song
	Link *string `json:"link,omitempty"`
}

// Group модель базы данных
// @Description Database model for a music group
type Group struct {
	// ID of the group
	// Required: true
	ID uint `gorm:"primaryKey" db:"id"`

	// Name of the group
	// Required: true
	Name string `gorm:"type:varchar(255);not null;unique;index" db:"name"`

	// Creation timestamp
	// Required: true
	CreatedAt time.Time `gorm:"index" db:"created_at"`

	// Update timestamp
	UpdatedAt time.Time `db:"updated_at"`
}

// Song модель базы данных
// @Description Database model for a song
type Song struct {
	// ID of the song
	// Required: true
	ID uint `gorm:"primaryKey" db:"id"`

	// ID of the associated group
	// Required: true
	GroupID uint `gorm:"not null;index;uniqueIndex:idx_group_song,priority:1" db:"group_id"`

	// Associated group
	Group Group `gorm:"foreignKey:GroupID"`

	// Group name
	GroupName string `gorm:"-" db:"group_name"`

	// Name of the song
	// Required: true
	SongName string `gorm:"type:varchar(255);not null;uniqueIndex:idx_group_song,priority:2;index" db:"song_name"` // добавляем отдельный индекс, если часто ищем по имени

	// Release date of the song
	ReleaseDate string `db:"release_date"`

	// Lyrics or text of the song
	Text string `gorm:"type:text" db:"text"`

	// External link to the song
	Link *string `gorm:"index" db:"link"`

	// Creation timestamp
	// Required: true
	CreatedAt time.Time `gorm:"index" db:"created_at"`

	// Update timestamp
	UpdatedAt time.Time `db:"updated_at"`
}
