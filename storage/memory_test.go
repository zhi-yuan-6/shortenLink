package storage

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestConcurrentIncrement(t *testing.T) {
	store := NewMemoryStore()
	code := "2KEN2O"

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			store.IncrementVisit(code)
		}()
	}
	wg.Wait()

	/*if store.VisitCount[code] != 100 {
		t.Errorf("并发访问计数错误，期望100，实际：%d", store.VisitCount[code])
	}*/
	record, _ := store.Visits.Load(code)
	if atomic.LoadInt64(&record.(*VisitRecord).Count) != 1000 {
		t.Errorf("并发访问计数错误，期望1000，实际：%d", atomic.LoadInt64(&record.(*VisitRecord).Count))
	}
}
