package testutils

/*
import (
	"errors"
	"math/rand"
	"shortenLink/storage"
	"testing"
)

type FaultStorage struct {
	storage.Storage
	failureRate float64 // 0.0 to 1.0
}

func (f *FaultStorage) Save(url string) (string, error) {
	if rand.Float64() < f.failureRate {
		return "", errors.New("injected storage failure")
	}
	return f.Storage.Save(url)
}

func TestSaveWithFailure(t *testing.T) {
	store := &FaultStorage{
		Storage:     storage.NewMemoryStore(),
		failureRate: 0.3,
	}
	//验证服务层能否正确处理存储层错误
}
*/
