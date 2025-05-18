package models

import (
	"gorm.io/gorm"
	"shortenLink/dto"
	"time"
)

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

// 创建记录
func CreateVisitStats(shortCode string) error {
	return dto.DB.Create(&VisitStats{ShortCode: shortCode}).Error
}

func GetVisitCount(code string) (int64, time.Time, error) {
	var visitStats VisitStats
	err := dto.DB.Where("short_code=?", code).First(&visitStats).Error
	return visitStats.VisitCount, visitStats.LastVisit, err
}

func IncrementVisit(code string) error {
	/*
		var visitStats VisitStats
		if err := dto.DB.Where("short_code = ?", code).First(&visitStats).Error; err != nil {
			return err
		}

		visitStats.VisitCount++
		visitStats.LastVisit = time.Now()

		return dto.DB.Model(&visitStats).Updates(visitStats).Error*/
	result := dto.DB.Model(&VisitStats{}).
		Where("short_code = ?", code).
		Updates(map[string]interface{}{
			"visit_count": gorm.Expr("visit_count + 1"),
			"last_visit":  time.Now()})
	return result.Error
}

/*
// IncrementVist 原子操作递增计数器
func (m *MemoryStore) IncrementVisit(code string) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	//初始化计数器
	if _, exists := m.VisitCount[code]; !exists {
		m.VisitCount[code] = 0
	}
		m.VisitCount[code]++
	//原子操作更新计数器
	record, _ := m.Visits.LoadOrStore(code, &VisitRecord{})
	//安全类型断言
	vRecord, ok := record.(*VisitRecord)
	if !ok {
		panic("unexpected type in visits map")
	}
	atomic.AddInt64(&vRecord.Count, 1)
	atomic.StoreInt64(&vRecord.LastVisit, time.Now().UnixNano())
	//}
*/
