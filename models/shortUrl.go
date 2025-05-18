package models

import (
	"shortenLink/dto"
	"time"
)

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

func (model *ShortUrl) CreateShortenUrl() error {
	return dto.DB.Create(model).Error
}

func GetShortenCode(url string) (string, error) {
	var shortUrl ShortUrl
	err := dto.DB.Where("original_url=?", url).First(&shortUrl).Error
	return shortUrl.ShortCode, err
}

func GetOriginalURL(code string) (string, error) {
	var shortURL ShortUrl
	err := dto.DB.Where("short_code=?", code).First(&shortURL).Error
	if err != nil {
		return "", err
	}
	return shortURL.OriginalUrl, nil
}

/*// 删除
func DeleteShortenURL(url string) error {
	err := storage.DB.Where("original_url=?", url).Error
	return err
}*/
