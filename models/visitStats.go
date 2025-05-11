package models

import "time"

type VisitStats struct {
	ID         int64  `gorm:"primaryKey"`
	ShortCode  string `gorm:"type:varchar(10);notNull"`
	VisitCount int64  `gorm:"notNull;default:0;check:visit_count >= 0"`
	LastVisit  time.Time
	CreatedAt  time.Time `gorm:"notNull;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"notNull;type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (VisitStats) TableName() string {
	return "visit_stats"
}
