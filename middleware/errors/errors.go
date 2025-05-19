package errors

import (
	"runtime"
)

type StackTrace []uintptr

type Error struct {
	Code       int
	Message    string
	Op         string     /// 操作名称，标识错误发生时正在执行的操作
	Stack      StackTrace // 堆栈跟踪信息，记录错误发生时的调用链
	WrappedErr error      // 包装的原始错误，用于链式错误处理
}

func New(code int, message, op string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Op:      op,
		Stack:   captureStackTrace(), // 捕获当前堆栈跟踪信息
	}
}

// Wrap 包装错误
func Wrap(err error, code int, message, op string) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		Op:         op,
		Stack:      captureStackTrace(),
		WrappedErr: err,
	}
}

// captureStackTrace 捕获当前的堆栈跟踪信息
// 该函数用于记录错误发生时的调用链，帮助开发者定位问题发生的位置
// 返回：
//
//	包含错误发生时调用链信息的StackTrace
func captureStackTrace() StackTrace {
	// 定义一个最多包含32个uintptr元素的数组，用于存储调用链信息
	// uintptr是Go语言中的一种无符号整数类型，通常用于表示内存地址
	var pcs [32]uintptr

	// 调用runtime.Callers函数捕获调用链信息
	// 参数3表示从调用链的第4个调用者开始捕获（因为Callers函数本身占用了前几个调用栈位置）
	// 返回实际捕获的调用者数量
	n := runtime.Callers(3, pcs[:])

	return pcs[0:n]
}
