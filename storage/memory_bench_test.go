package storage

import "testing"

func BenchmarkConcurrentIncrement(b *testing.B) {
	store := NewMemoryStore()
	code := "testCo"
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.IncrementVisit(code)
		}
	})
}
