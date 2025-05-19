package errors

const (
	// 客户端错误 1000-1999
	ErrCodeInvalidInput = 1001 // 输入参数错误
	ErrCodeUnauthorized = 1002 // 未授权访问
	ErrCodeRateLimit    = 1003 // 请求限流

	// 服务错误 2000-2999
	ErrCodeDatabase        = 2001 // 数据库操作失败
	ErrCodeCache           = 2002 // 缓存服务异常
	ErrCodeExternalService = 2003 // 外部服务故障

	// 业务错误 3000-3999
	ErrCodeShortCodeConflict = 3001 // 短码冲突
	ErrCodeURLExpired        = 3002 // URL已过期
)
