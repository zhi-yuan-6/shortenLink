package shortLink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"
)

// 测试函数
func TestShortenHandler(t *testing.T) {
	//设置测试参数
	const (
		numberRequest = 1000 //总请求数
		concurrency   = 100  //并发数
	)

	//创建一个 WaitGroup 来等待所有请求完成
	var wg sync.WaitGroup
	wg.Add(numberRequest)

	//创建一个错误计数器
	var mu sync.Mutex
	errors := 0

	//创建一个信号量来控制并发数 执行一次sendRequest函数，就会占用一个通道
	sem := make(chan struct{}, concurrency)

	//设置一个伪随机数
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	// 发送请求函数
	sendRequest := func() {
		//准备 JSON 数据
		//生成一个伪随机数
		randNumber := r.Intn(100000)
		data := ShortenRequest{URL: "https://test" + fmt.Sprintf("%d", randNumber) + ".com"}
		jsonData, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}

		defer wg.Done()
		sem <- struct{}{} //获取信号量

		//bytes.Buffer 是一个可读写的缓冲区，实现了 io.Reader 和 io.Writer 接口，可以方便地进行字节切片的读写操作。
		resp, err := http.Post("http://localhost:8080/shorten", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			mu.Lock()
			errors++
			mu.Unlock()
		} else {
			//读取响应体
			//_, err := ioutil.ReadAll(resp.Body)
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

	//启动并发请求
	for i := 0; i < numberRequest; i++ {
		go sendRequest()
	}
	//等待所有请求完成
	wg.Wait()

	//检查错误数
	if errors > 0 {
		t.Errorf("请求失败错误数: %d", errors)
	}

}
