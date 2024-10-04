package task

import "SomersaultCloud/api/middleware/taskchain"

type AskTask interface {
	InitContextData(args ...any) *taskchain.TaskContext
	// PreCheckDataTask 数据的前置检查 & 组装TaskContextData对象
	PreCheckDataTask(tc *taskchain.TaskContext)
	// GetHistoryTask 从DB or Cache获取历史记录
	GetHistoryTask(tc *taskchain.TaskContext)
	// GetBotTask 获取prompt & model
	GetBotTask(tc *taskchain.TaskContext)
	// TODO 微调 TBD
	AdjustmentTask(tc *taskchain.TaskContext)
	// AssembleReqTask 组装rpc请求体
	AssembleReqTask(tc *taskchain.TaskContext)
	// CallApiTask 调用api
	CallApiTask(tc *taskchain.TaskContext)
	// ParseRespTask 转换rpc后响应数据
	ParseRespTask(tc *taskchain.TaskContext)
}

// TitleTask 更新标题的时候 所调用的责任链方法 较为简单 可以与上方节点混着用
type TitleTask interface {
	InitContextData(args ...any) *taskchain.TaskContext
	// PreTitleTask 这个任务把前置检查 历史记录 获取bot 全部跳过
	PreTitleTask(tc *taskchain.TaskContext)
	// AssembleTitleReqTask 组装请求
	AssembleTitleReqTask(tc *taskchain.TaskContext)
	// Convert2AskTask 其他的逻辑同AskTask
	Convert2AskTask(tc *taskchain.TaskContext)
}

// ConvertTask 责任链中节点在不同方法中使用的时候 需要根据需求进行一定定制修改
// 对于变化较小的改动 直接在此定义节点并使用即可 上方title_task那其实也可以这么干
type ConvertTask interface {
	StreamArgsTask(tc *taskchain.TaskContext)
	StreamPublishTask(tc *taskchain.TaskContext)
}
