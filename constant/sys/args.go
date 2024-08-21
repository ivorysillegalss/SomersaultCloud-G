package sys

const (
	// GoRoutinePoolTypesAmount 线程池种类数量
	GoRoutinePoolTypesAmount = 1
	// ExecuteRpcGoRoutinePool 执行chat任务的线程池的 id
	ExecuteRpcGoRoutinePool = 1
	// GenerationResponseChannelBuffer 新初始化传递rpc数据的channel大小
	GenerationResponseChannelBuffer = 100
	// GenerateQueryRetryLimit 查询rpc返回值最大次数
	GenerateQueryRetryLimit = 10

	// DefaultPoolGoRoutineAmount 默认的线程池中线程的数量
	DefaultPoolGoRoutineAmount = 20

	// GzipCompress 压缩方式为Gzip
	GzipCompress = 1
)
