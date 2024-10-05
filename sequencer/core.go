package sequencer

import (
	"SomersaultCloud/domain"
	"SomersaultCloud/infrastructure/log"
	"fmt"
	"sync"
	"time"

	"github.com/RussellLuo/timingwheel"
)

type StreamData struct {
	sequenceIndex   int                           // 当前已处理的序列号
	unSequenceValue map[int]domain.ParsedResponse // 存储失序的消息
	sequenceValue   chan domain.ParsedResponse    // 存储按序到达的消息
	timer           *timingwheel.Timer            // 管理该流的计时器
	active          bool                          // 标记流的状态，是否仍然活跃
}

var (
	streams = make(map[int]*StreamData) // 管理所有流的数据
	mu      sync.Mutex                  // 用于并发控制
	tw      = timingwheel.NewTimingWheel(100*time.Millisecond, 10)
)

const streamTimeout = 10 * time.Second // 设置整个流的超时时间
// 第一条信息的索引
const firstMessageIndex = 1

func init() {
	// 启动时间轮
	tw.Start()
}

// Setup 将传过来的流式数据进行顺序的判断
// 若没问题则放入管道中 等待客户端请求下发
// TODO 差一个正常发送信息，正常断流的情况
func Setup(parsedResp domain.ParsedResponse) {
	identity := parsedResp.GetIdentity() // 获取消息的身份标识
	index := parsedResp.GetIndex()       // 获取消息的序号

	mu.Lock()
	defer mu.Unlock()

	stream, exists := streams[identity]

	if !exists {
		//若流不存在
		if index == firstMessageIndex {
			// 第一种情况：这是第一条消息，创建新的流数据
			stream = &StreamData{
				sequenceIndex:   firstMessageIndex,
				unSequenceValue: make(map[int]domain.ParsedResponse),
				sequenceValue:   make(chan domain.ParsedResponse, 100), // 可根据需要设置缓冲区大小
				active:          true,
			}
			streams[identity] = stream
			startStreamTimer(identity)
			stream.sequenceValue <- parsedResp
		} else {
			// 收到了非第一条消息，但流并不存在，记录错误
			log.GetTextLogger().Error(fmt.Sprintf("No active stream for identity %d. Discarding message.\n", identity))
			return
		}
	} else {
		if !stream.active {
			// 流存在但已不活跃，丢弃消息
			log.GetTextLogger().Error(fmt.Sprintf("Stream for identity %d is no longer active. Discarding message.\n", identity))
			return
		}

		if index == firstMessageIndex {
			// 第二种情况：流正常结束后，收到新的第一条消息，重置流数据
			resetStreamData(stream)
			stream.sequenceIndex = firstMessageIndex
			stream.active = true
			startStreamTimer(identity)
			stream.sequenceValue <- parsedResp
		} else {
			// 第三种情况：正常接收后续消息
			if stream.sequenceValue == nil {
				// 管道不存在，说明流已经结束或出错
				log.GetTextLogger().Error(fmt.Sprintf("No valid channel for identity %d. Discarding message.\n", identity))
				return
			}

			i := stream.sequenceIndex

			if index == i+1 {
				// 按序接收到消息
				stream.sequenceValue <- parsedResp
				stream.sequenceIndex = i + 1
				checkUnorderedMessages(stream)
			} else if index > i+1 {
				// 发生失序，存储消息
				stream.unSequenceValue[index] = parsedResp
			} else {
				// 重复或过期的消息，忽略或记录
				log.GetTextLogger().Warn(fmt.Sprintf("Received outdated message for identity %d. Index: %d\n", identity, index))
			}
		}
	}
}

// 检查并处理失序的消息
func checkUnorderedMessages(stream *StreamData) {
	i := stream.sequenceIndex

	for {
		nextIndex := i + 1
		if msg, exists := stream.unSequenceValue[nextIndex]; exists {
			// 处理失序的下一条消息
			stream.sequenceValue <- msg
			stream.sequenceIndex = nextIndex
			delete(stream.unSequenceValue, nextIndex)
			i = nextIndex
		} else {
			break
		}
	}
}

// 重置流数据，但保留管道以避免重新分配
func resetStreamData(stream *StreamData) {
	stream.sequenceIndex = 0
	stream.unSequenceValue = make(map[int]domain.ParsedResponse)
	// 清空管道中的残留数据
	for len(stream.sequenceValue) > 0 {
		<-stream.sequenceValue
	}
	if stream.timer != nil {
		stream.timer.Stop()
	}
}

// 使用 timingwheel 启动计时器
func startStreamTimer(identity int) {
	stream := streams[identity]

	// 如果已有计时器，则先停止它
	if stream.timer != nil {
		stream.timer.Stop()
	}

	// 使用时间轮的 AfterFunc 来设定超时
	stream.timer = tw.AfterFunc(streamTimeout, func() {
		handleStreamTimeout(identity)
	})
}

// 超时处理函数
func handleStreamTimeout(identity int) {
	mu.Lock()
	defer mu.Unlock()

	stream, exists := streams[identity]
	if !exists {
		return
	}

	log.GetTextLogger().Error(fmt.Sprintf("Stream for identity %d timed out. Cleaning up...\n", identity))

	// 标记流为非活跃状态
	stream.active = false

	// 停止计时器
	if stream.timer != nil {
		stream.timer.Stop()
		stream.timer = nil
	}

	// 清理失序消息
	stream.unSequenceValue = nil

	// 关闭并清理管道
	if stream.sequenceValue != nil {
		close(stream.sequenceValue)
		stream.sequenceValue = nil
	}

	// 从全局映射中删除该流（可选）
	delete(streams, identity)
}
