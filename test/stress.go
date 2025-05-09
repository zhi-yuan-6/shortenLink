package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var (
	logSet = struct {
		sync.Mutex
		set map[string]struct{}
	}{
		set: make(map[string]struct{}), //初始化map
	}
)

func main() {
	var (
		totalRequests int64
		successCount  int64
	)

	//解析命令行参数，由控制台输入参数
	workers := flag.Int("workers", 10, "并发worker数量")
	duration := flag.Int("duration", 2, "测试持续时间(秒)")
	targetURL := flag.String("url", "http://127.0.0.1:8080/2KEN2O", "目标URL")
	flag.Parse()

	if _, err := os.Stat("load_test.log"); err == nil {
		logFile, err := os.OpenFile("load_test.log", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("无法打开日志文件: %v", err)
		}
		log.SetOutput(logFile) //设置日志输出到文件
	} else {
		//创建日志文件
		logFile, err := os.Create("load_test.log")
		if err != nil {
			log.Fatalf("无法创建日志文件: %v", err)
		}
		defer logFile.Close()

		log.SetOutput(logFile) //设置日志输出到文件
	}

	done := make(chan bool) //创建一个通道，用于通知协程退出
	start := time.Now()

	//启动统计协程,定期统计协程退出
	go func() {
		ticker := time.NewTicker(1 * time.Second) //创建一个定时器，每秒钟触发一次
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C: //每秒执行以此
				currentReq := atomic.LoadInt64(&totalRequests)
				currentSuccess := atomic.LoadInt64(&successCount)
				//计算并打印每秒查询率和成功率
				log.Printf("QPS:%d/s,Success Rate:%.2f%%\n", currentReq, float64(currentSuccess)/float64(currentReq)*100)

				//充值计数器，为下一秒统计做准备
				atomic.StoreInt64(&totalRequests, 0)
				atomic.StoreInt64(&successCount, 0)
			case <-done: //接收到退出信号，退出循环
				return
			}
		}
	}()

	//启动worker
	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//创建一个客户端，设置超时时间,增加连接池（1.把建立客户端挪出去，每多少个请求建立一个连接池,不过没必要，模仿的是多用户访问）
			client := &http.Client{
				Timeout: 5 * time.Second,
				Transport: &http.Transport{
					MaxIdleConns:        100,              //最大空闲连接数
					MaxIdleConnsPerHost: 100,              //每个主机的最大空闲连接数
					MaxConnsPerHost:     200,              //每个主机的最大连接数
					IdleConnTimeout:     10 * time.Second, //空闲连接的超时时间
				},
			}
			//循环发送请求，直到测试时间结束
			for time.Since(start).Seconds() < float64(*duration) {
				//发送请求
				resp, err := client.Get(*targetURL)
				//原子增肌啊总请求数
				atomic.AddInt64(&totalRequests, 1)
				if err != nil {
					//log.Printf("HTTP request failed: %v", err)
					logOnce("HTTP request failed: %v", err)
				} else {
					resp.Body.Close() //关闭响应体，防止内存泄漏
					if resp.StatusCode == http.StatusFound {
						//原子增加成功请求数
						atomic.AddInt64(&successCount, 1)
					}
				}
			}
		}()
	}
	wg.Wait()
	close(done) //关闭通道，通知统计协程退出
}

func logOnce(format string, args ...interface{}) {
	var logStr string
	if len(args) == 0 {
		logStr = format
	} else {
		logStr = fmt.Sprintf(format, args...)
	}

	logSet.Lock()
	defer logSet.Unlock()

	if _, exists := logSet.set[logStr]; !exists {
		log.Println(logStr)
		logSet.set[logStr] = struct{}{} //将日志添加到集合中
	}
}
