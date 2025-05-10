package storage

import (
	"sync"
	"sync/atomic"
	"time"
)

type VisitRecord struct {
	Count     int64
	LastVisit int64
}
type MemoryStore struct {
	Mu          sync.RWMutex         // 读写锁
	UrlMap      map[string]string    // 存储短链接和原始url的映射关系
	ReverseMap  map[string]string    // 存储原始url和短链接的反向映射（用于幂等性）
	CreatedTime map[string]time.Time //短链接创建时间
	//VisitCount  map[string]int       // 访问计数器
	Visits sync.Map
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		UrlMap:      make(map[string]string),
		ReverseMap:  make(map[string]string),
		CreatedTime: make(map[string]time.Time),
		//VisitCount:  make(map[string]int),
	}
}

// IncrementVist 原子操作递增计数器
func (m *MemoryStore) IncrementVisit(code string) {
	/*m.Mu.Lock()
	defer m.Mu.Unlock()
	//初始化计数器
	if _, exists := m.VisitCount[code]; !exists {
		m.VisitCount[code] = 0
	}
	m.VisitCount[code]++*/
	//原子操作更新计数器
	record, _ := m.Visits.LoadOrStore(code, &VisitRecord{})
	//安全类型断言
	vRecord, ok := record.(*VisitRecord)
	if !ok {
		panic("unexpected type in visits map")
	}
	atomic.AddInt64(&vRecord.Count, 1)
	atomic.StoreInt64(&vRecord.LastVisit, time.Now().UnixNano())
}
