package models

import "time"

type ShortUrl struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	ShortCode   string    `gorm:"type:varchar(10);unique;notNull"`
	OriginalUrl string    `gorm:"type:text;notNull"`
	CreateAt    time.Time `gorm:"notNull;default:CURRENT_TIMESTAMP"`
	ExpiresAt   *time.Time
	DeleteAt    *time.Time
	//IsDeleted   bool `gorm:"notNull;default:false"`
}

func (ShortUrl) TableName() string {
	return "short_urls"
}
