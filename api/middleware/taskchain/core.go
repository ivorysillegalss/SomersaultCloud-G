package taskchain

import (
	task2 "SomersaultCloud/constant/task"
	"SomersaultCloud/usecase/task"
)

type TaskContext struct {
	BusinessType        string
	BusinessCode        int
	TaskContextData     *TaskContextData
	TData               any // 特定类型的数据
	Exception           bool
	TaskContextResponse *TaskContextResponse
}

// TaskContextData 此处定义责任链节点所使用的共同数据
type TaskContextData struct {
	*task.AskContextData
}

type TaskContextResponse struct {
	Message string
	Code    int
}

type TaskContextFactory struct {
	Nodes       []func(tc *TaskContext)
	TaskContext *TaskContext
}

func NewTaskContextFactory() *TaskContextFactory {
	return &TaskContextFactory{}
}

// InterruptExecute 错误执行包装类
func (t *TaskContext) InterruptExecute(message string) {
	t.Exception = true
	t.TaskContextResponse.Code = task2.FailCode
	t.TaskContextResponse.Message = message
}

// ExecuteChain 执行责任链
func (t *TaskContextFactory) ExecuteChain(handlers ...func(tc *TaskContext)) {
	for _, handler := range handlers {
		//具体的错误处理包装在TaskResponse中 由各节点自行处理
		if t.TaskContext.Exception {
			// 具体的错误原因类型在 实现类中包装
			return
		}
		handler(t.TaskContext)
	}
}

// Puts 加入节点
func (t *TaskContextFactory) Puts(handlers ...func(tc *TaskContext)) {
	nodes := &t.Nodes
	//可变参数 + 解包 加入节点
	*nodes = append(*nodes, handlers...)
}

// List 列举节点
func (t *TaskContextFactory) List() []func(tc *TaskContext) {
	return t.Nodes
}

// TODO 目前是即用即装配链子 初步封装template
// 			kv形式存储		 k for 业务类型 & v for 链子配置 形成工厂类策略化等等乱七八糟的
// 			这个以后完善 没有实际链子太抽象了
