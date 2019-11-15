package main

import (
	"strings"
	"sync"
)

const (
	// TaskPrefix 任务key前缀
	TaskPrefix string = "task-"
	// CommitTaskPrefix 提交任务key前缀
	CommitTaskPrefix string = "commit-"
	// ClearTaskPrefix 清除任务
	ClearTaskPrefix string = "clear-"
)

// Event 事件类型
type Event struct {
	Key   string
	Name  string
	Value interface{}
}

// EventListener 用于接收消息回调
type EventListener interface {
	onEvent(event *Event)
}

// MemoryQueue 内存消息队列
type MemoryQueue struct {
	done      chan struct{}
	queue     chan Event
	listeners []EventListener
	wg        sync.WaitGroup
}

func NewMemoryQueue(n int) *MemoryQueue{
	return &MemoryQueue{
		done:      make(chan struct{}, 0),
		queue:     make(chan Event, n),
		listeners: nil,
	}
}
// Push 添加数据
func (mq *MemoryQueue) Push(eventType, name string, value interface{}) {
	mq.queue <- Event{Key: eventType + name, Name: name, Value: value}
	mq.wg.Add(1)
}

// AddListener 添加监听器
func (mq *MemoryQueue) AddListener(listener EventListener) bool {
	for _, item := range mq.listeners {
		if item == listener {
			return false
		}
	}
	mq.listeners = append(mq.listeners, listener)
	return true
}

// Notify 分发消息
func (mq *MemoryQueue) Notify(event *Event) {
	defer mq.wg.Done()
	for _, listener := range mq.listeners {
		listener.onEvent(event)
	}
}

func (mq *MemoryQueue) poll() {
	for {
		select {
		case <-mq.done:
			break
		case event := <-mq.queue:
			mq.Notify(&event)
		}
	}
}

// Start 启动内存队列
func (mq *MemoryQueue) Start() {
	go mq.poll()
}

// Stop 停止内存队列
func (mq *MemoryQueue) Stop() {
	mq.wg.Wait()
	close(mq.done)
}

type ConfigUpdateCallback func(map[string]string)

// Worker 工作进程
type Worker struct {
	name                string
	deferredTaskUpdates map[string][]Task
	onCommit            ConfigUpdateCallback
}

func NewWorker(name string, f func(data map[string]string)) *Worker{
	return &Worker{
		name:                name,
		deferredTaskUpdates: make(map[string][]Task),
		onCommit:            nil,
	}
}

func (w *Worker) onEvent(event *Event) {
	switch {
	// 获取任务事件
	case strings.Contains(event.Key, TaskPrefix):
		w.onTaskEvent(event)
		// 清除本地延迟队列里面的任务
	case strings.Contains(event.Key, ClearTaskPrefix):
		w.onTaskClear(event)
		// 获取commit事件
	case strings.Contains(event.Key, CommitTaskPrefix):
		w.onTaskCommit(event)
	}
}

func (w *Worker) onTaskClear(event *Event) {
	task, err := event.Value.(Task)
	if !err {
		// log
		return
	}
	_, found := w.deferredTaskUpdates[task.Group]
	if !found {
		return
	}
	delete(w.deferredTaskUpdates, task.Group)
	// 还可以继续停止本地已经启动的任务
}

// onTaskCommit 接收任务提交, 从延迟队列中取出数据然后进行业务逻辑处理
func (w *Worker) onTaskCommit(event *Event) {
	// 获取之前本地接收的所有任务
	tasks, found := w.deferredTaskUpdates[event.Name]
	if !found {
		return
	}

	// 获取配置
	config := w.getTasksConfig(tasks)
	if w.onCommit != nil {
		w.onCommit(config)
	}
	delete(w.deferredTaskUpdates, event.Name)
}

// onTaskEvent 接收任务数据，此时需要丢到本地暂存不能进行应用
func (w *Worker) onTaskEvent(event *Event) {
	task, err := event.Value.(Task)
	if !err {
		// log
		return
	}

	// 保存任务到延迟更新map
	configs, found := w.deferredTaskUpdates[task.Group]
	if !found {
		configs = make([]Task, 0)
	}
	configs = append(configs, task)
	w.deferredTaskUpdates[task.Group] = configs
}

// getTasksConfig 获取task任务列表
func (w *Worker) getTasksConfig(tasks []Task) map[string]string {
	config := make(map[string]string)
	for _, t := range tasks {
		config = t.updateConfig(config)
	}
	return config
}

type Task struct {
	Name string
	Group string
	Config map[string]string
}

func (t Task) updateConfig(config map[string]string) map[string]string{
	return config
}

func main() {
	// 生成一个内存队列启动
	queue := NewMemoryQueue(10)
	queue.Start()

	// 生成一个worker
	name := "test"
	worker := NewWorker(name, func(data map[string]string) {
		for key, value := range data {
			println("worker get task key: " + key + " value: " + value)
		}
	})
	// 注册到队列中
	queue.AddListener(worker)

	taskName := "test"
	// events 发送的任务事件
	configs := []map[string]string{
		map[string]string{"task1": "SendEmail", "params1": "Hello world"},
		map[string]string{"task2": "SendMQ", "params2": "Hello world"},
	}

	// 分发任务
	queue.Push(ClearTaskPrefix, taskName, nil)
	for _, conf := range configs {
		queue.Push(TaskPrefix, taskName, Task{Name: taskName, Group: taskName, Config: conf})
	}
	queue.Push(CommitTaskPrefix, taskName, nil)
	// 停止队列
	queue.Stop()
}
