package redirect

import (
	"io"
	"net/http"
	"sync"
	"testing"
)

func TestRedirectHandler(t *testing.T) {
	const (
		numberRequest = 2000 //总请求数
		concurrency   = 100  //并发数
	)
	var wg sync.WaitGroup
	wg.Add(numberRequest)

	//计数器
	var mu sync.Mutex
	errors := 0

	//创建一个信号量来控制并发数 执行一次sendRequest函数，就会占用一个通道
	sem := make(chan struct{}, concurrency)

	redirect := func() {
		defer wg.Done()
		sem <- struct{}{}
		//发起请求
		resp, err := http.Get("http://localhost:8080/2KEN2O")
		if err != nil {
			mu.Lock()
			errors++
			mu.Unlock()
		} else {
			//读取响应体
			_, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				mu.Lock()
				errors++
				mu.Unlock()
			}
		}
		<-sem //释放信号量
	}

	for i := 0; i < numberRequest; i++ {
		go redirect()
	}
	wg.Wait()

	//检查错误数
	if errors > 0 {
		t.Errorf("有 %d 个请求失败", errors)
	} else {
		t.Logf("所有请求成功")
	}
}
