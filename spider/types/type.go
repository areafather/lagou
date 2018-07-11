package types

import (
	"io"
	"net/http"
)

// Task : 存储请求以及请求响应需要的处理器函数
type Task struct {
	Request *http.Request
	Handler func(io.Reader) ([]Task, error)
}

// NoneNewTask : 表示将不会再有新的任务加入
type NoneNewTask struct {
}

// 实现 error interface
func (e NoneNewTask) Error() string {
	return "Will none task, shoule close taskChan"
}
