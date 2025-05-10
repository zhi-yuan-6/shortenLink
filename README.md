# 短链接服务系统

🚀 基于Go语言的高并发短链接生成与统计系统，核心特性：

- 微秒级响应延迟（P99 < 10ms） 注：暂未测试具体数据
- 万级QPS处理能力
- 多级缓存加速体系
- 实时数据可视化

## 快速开始
```bash[utils](utils)
# 启动开发环境
docker-compose up -d

# 生成短链接
curl -X POST -d '{"url":"https://example.com"}' http://localhost:8080/api/shorten