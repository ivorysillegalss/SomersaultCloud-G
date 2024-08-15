package sys

const (
	// GoRoutinePoolTypesAmount 线程池种类数量
	GoRoutinePoolTypesAmount = 1
	// ExecuteRpcGoRoutinePool 执行rap任务的线程池的 id
	ExecuteRpcGoRoutinePool = 1
	// GenerationResponseChannelBuffer 新初始化传递rpc数据的channel大小
	GenerationResponseChannelBuffer = 100
	// GenerateQueryRetryLimit 查询rpc返回值最大次数
	GenerateQueryRetryLimit = 10
)
