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
