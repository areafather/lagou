package engine

import (
	"log"

	"github.com/shawpo/lagou/spider/request"
	"github.com/shawpo/lagou/spider/scheduler"
	"github.com/shawpo/lagou/spider/types"
)

// Engine : 爬虫引擎
type Engine struct {
	*scheduler.Scheduler
	WorkerCount int
}

// Run : 运行
func (e *Engine) Run(tasks ...types.Task) {
	log.Println("engine is runing")
	e.creatWorker()

	log.Println("created worker")

	e.Scheduler.AddTask(tasks...)
	e.Scheduler.Run()

}

func (e *Engine) creatWorker() {
	for i := 0; i < e.WorkerCount; i++ {
		worker := make(chan types.Task)
		go func() {
			for {
				// 通知调度器已准备好接受任务
				e.Scheduler.WorkerReady(worker)
				// 阻塞等待任务
				task := <-worker
				tasks, err := Do(task)
				if _, ok := err.(types.NoneNewTask); ok {
					// 将不会再有新任务
					e.Scheduler.CloseTaskChan()
				} else if err != nil {
					// 完成任务时出错，将该任务再一次加入调度器
					e.Scheduler.AddTask(task)
				} else {
					e.Scheduler.CompleteTask(task)
					if len(tasks) > 0 {
						// 将返回的新任务加入调度器
						e.Scheduler.AddTask(tasks...)
					}

				}
			}
		}()
	}
}

// Do : 具体完成任务
func Do(t types.Task) ([]types.Task, error) {
	resp, err := request.Fetch(t.Request)
	if err != nil {
		log.Printf("Fetcher: error fetching url %s: %v", t.Request.URL, err)
		return []types.Task{}, err
	}
	return t.Handler(resp.Body)
}
