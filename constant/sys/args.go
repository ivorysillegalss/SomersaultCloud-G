package sys

import "time"

const (
	// GoRoutinePoolTypesAmount 线程池种类数量
	GoRoutinePoolTypesAmount = 1
	// ExecuteRpcGoRoutinePool 执行chat任务的线程池的 id
	ExecuteRpcGoRoutinePool = 1
	// GenerationResponseChannelBuffer 新初始化传递rpc数据的channel大小
	GenerationResponseChannelBuffer = 100
	// GenerateQueryRetryLimit 查询rpc返回值最大次数
	GenerateQueryRetryLimit = 15

	// DefaultPoolGoRoutineAmount 默认的线程池中线程的数量
	DefaultPoolGoRoutineAmount = 20

	// GzipCompress 压缩方式为Gzip 搭配Json序列化
	GzipCompress = "gzip&json"
	// ProtoBufCompress 序列化方式
	ProtoBufCompress = "protobuf"

	StreamOverSignal = "data: [DONE]"
	StreamPrefix     = "data: "
)

// 排序相关
const (
	Finish         = 0
	Timeout        = -1
	IllegalRequest = -2

	NormallyEndExpiration = time.Second      //指单次会话所有流信息存储在channel中的缓存时间
	StreamTimeout         = 10 * time.Second // 设置整个流的超时时间
	FirstMessageIndex     = 1                // 第一条信息的索引

)
