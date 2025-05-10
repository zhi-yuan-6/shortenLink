package shorten

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"shortenLink/storage"
	"strings"
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

	// 发送请求函数
	sendRequest := func() {
		//准备 JSON 数据
		//生成一个伪随机数
		source := rand.NewSource(time.Now().UnixNano()) //若使用全局随机数生成器，需要加锁
		r := rand.New(source)
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

func TestShortenHandler2(t *testing.T) {
	store := storage.NewMemoryStore()
	engine := gin.Default()
	engine.POST("/api/shorten", ShortenHandler(store))
	//测试用例矩阵
	tests := []struct {
		name       string
		payload    string
		wantStatus int
		wantKey    string
	}{
		{"valid_url", `{"url": "http://google.com"}`, 200, "short_link"},
		{"invalid__json", `{invalid}`, 400, "error"},
		{"empty_url", `{"url":""}`, 400, "error"},
		//{"non_http_url", `{"url":"ftp://example.com"}`, 400, "error"},  //暂时未添加url检验
	}

	for _, tt := range tests {
		//发起请求
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(w, req)
			//检查响应状态码
			if w.Code != tt.wantStatus {
				t.Errorf("状态码错误 ，期望%d 实际%d", tt.wantStatus, w.Code)
			}
			//检查响应数据
			var resp map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatal("响应数据解析失败：", err)
			}

			if _, ok := resp[tt.wantKey]; !ok {
				t.Errorf("响应数据中缺少%s字段", tt.wantKey)
			}
		})
	}
}
