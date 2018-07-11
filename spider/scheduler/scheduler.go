package scheduler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/shawpo/lagou/params"

	"github.com/shawpo/lagou/spider/types"
)

// Scheduler : 爬虫调度器
type Scheduler struct {
	taskChan      chan types.Task
	readyWorker   chan chan types.Task
	tasksQ        []types.Task
	readyWorkersQ []chan types.Task
	completedTask map[string]string
}

// New : 返回一个调度器
func New() *Scheduler {
	return &Scheduler{
		taskChan:      make(chan types.Task),
		readyWorker:   make(chan chan types.Task),
		completedTask: make(map[string]string),
	}
}

// AddTask : 加入新任务
func (s *Scheduler) AddTask(tasks ...types.Task) {
	go func() {
		for _, t := range tasks {
			if !s.IsCompleted(t) {
				s.taskChan <- t
			}
		}
	}()
}

// CloseTaskChan : 关闭taskChan表示不会再有新的任务加入
func (s *Scheduler) CloseTaskChan() {
	close(s.taskChan)
}

// WorkerReady : 添加一个已准备好接受任务的worker
func (s *Scheduler) WorkerReady(w chan types.Task) {
	s.readyWorker <- w
}

// CompleteTask : 提交已完成任务
func (s *Scheduler) CompleteTask(task types.Task) {
	str, err := taskString(task)
	if err != nil {
		log.Println(err)
		return
	}
	s.completedTask[str] = ""
}

// IsCompleted : 当前任务是否已经被完成
func (s *Scheduler) IsCompleted(task types.Task) bool {
	str, err := taskString(task)
	if err != nil {
		log.Println(err)
		return false
	}
	if _, exists := s.completedTask[str]; exists {
		log.Printf("task is IsCompleted: %s\n", str)
		return true
	}
	return false
}

func taskString(task types.Task) (string, error) {
	var body []byte
	if task.Request.GetBody != nil {
		reader, _ := task.Request.GetBody()
		body, _ = ioutil.ReadAll(reader)
	}
	t := struct {
		Method string
		URL    *url.URL
		Body   []byte
	}{
		task.Request.Method,
		task.Request.URL,
		body,
	}

	json, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(json), nil
}

// Run : 运行
func (s *Scheduler) Run() {
	log.Println("scheduler is run")
	// 死循环调度任务
scheduler:
	for {
		var activeWork chan types.Task
		var activeTask types.Task
		if len(s.readyWorkersQ) > 0 && len(s.tasksQ) > 0 {
			activeTask = s.tasksQ[0]
			activeWork = s.readyWorkersQ[0]
		}
		select {
		case task, ok := <-s.taskChan:
			if task.Request != nil && task.Handler != nil {
				s.tasksQ = append(s.tasksQ, task)
			} else if !ok && len(s.tasksQ) <= 0 && len(s.readyWorkersQ) == params.WORKERCOUNT {
				// !ok 表示taskChan已经关闭，不会再有新任务加入
				// len(s.tasksQ) <= 0 表示任务队列已经没有待完成的任务
				// len(s.readyWorkersQ) == params.WORKERCOUNT :
				// 表示所有的worker都已经完成自己的任务，
				// 正在准备接受新任务, 即没有正在处理中的任务
				// 三个条件同时满足说明所有的任务已经完成
				log.Println("任务全部完成！")
				break scheduler
			}
		case worker := <-s.readyWorker:
			s.readyWorkersQ = append(s.readyWorkersQ, worker)
		case activeWork <- activeTask:
			s.readyWorkersQ = s.readyWorkersQ[1:]
			s.tasksQ = s.tasksQ[1:]
		}

	}
}
