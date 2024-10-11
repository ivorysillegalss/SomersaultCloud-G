package sequencer

import (
	"SomersaultCloud/constant/common"
	"SomersaultCloud/constant/sys"
	"SomersaultCloud/domain"
	"SomersaultCloud/infrastructure/log"
	"fmt"
	"github.com/thoas/go-funk"
	"strconv"
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
	activeChan      chan int                      // 通过select返回运行时状态码
	version         int64                         // 以时间戳作为版本控制 CAS
}

var (
	streams = make(map[int]*StreamData) // 管理所有流的数据
	mu      sync.Mutex                  // 用于并发控制
	tw      = timingwheel.NewTimingWheel(100*time.Millisecond, 10)
)

func init() {
	// 启动时间轮
	tw.Start()
}

type Sequencer struct {
}

func NewSequencer() *Sequencer {
	return &Sequencer{}
}

func (c *Sequencer) GetData(userId int) (chan domain.ParsedResponse, chan int) {
	data := streams[userId]
	if funk.IsEmpty(data) {
		return nil, nil
	}
	return data.sequenceValue, data.activeChan
}

// Setup 将传过来的流式数据进行顺序的判断
// 若没问题则放入管道中 等待客户端请求下发
func (c *Sequencer) Setup(parsedResp domain.ParsedResponse) {
	if funk.IsEmpty(parsedResp) {
		log.GetTextLogger().Error("nil message")
		return
	}
	identity := parsedResp.GetIdentity() // 获取消息的身份标识
	index := parsedResp.GetIndex()       // 获取消息的序号
	finishReason := parsedResp.GetFinishReason()

	mu.Lock()
	defer mu.Unlock()

	stream, exists := streams[identity]

	if funk.NotEmpty(finishReason) {

		log.GetTextLogger().Info("finish receiving message for user : " + strconv.Itoa(parsedResp.GetIdentity()) + ", end for reason :" + parsedResp.GetFinishReason())
		//normallyEndStream(identity, stream.version)

		return
	}

	if !exists {
		//若流不存在
		if index == sys.FirstMessageIndex {
			// 第一种情况：这是第一条消息，创建新的流数据
			stream = &StreamData{
				sequenceIndex:   sys.FirstMessageIndex,
				unSequenceValue: make(map[int]domain.ParsedResponse),
				sequenceValue:   make(chan domain.ParsedResponse, 100), // 可根据需要设置缓冲区大小
				active:          true,
			}
			streams[identity] = stream
			startStreamTimer(identity)

			fmt.Println(parsedResp.GetIndex())

			stream.sequenceValue <- parsedResp
		} else {
			// 收到了非第一条消息，但流并不存在，记录错误
			log.GetTextLogger().Error(fmt.Sprintf("No active stream for identity %d. Discarding message.\n with Index %d", identity, index))
			return
		}
	} else {
		if !stream.active {
			// 流存在但已不活跃，丢弃消息
			log.GetTextLogger().Error(fmt.Sprintf("Stream for identity %d is no longer active. Discarding message.\n", identity))
			return
		}

		if index == sys.FirstMessageIndex {
			// 第二种情况：流正常结束后，收到新的第一条消息，重置流数据
			resetStreamData(stream)
			stream.sequenceIndex = sys.FirstMessageIndex
			stream.active = true
			startStreamTimer(identity)
			stream.sequenceValue <- parsedResp
		} else {
			// 第三种情况：正常接收后续消息
			if stream.sequenceValue == nil {
				// 管道不存在，说明流已经结束或出错 并发问题
				log.GetTextLogger().Error(fmt.Sprintf("No valid channel for identity %d. Discarding message.\n", identity))
				//TODO 这里并发问题,到底要不要加还没测过
				streams[identity].activeChan <- sys.IllegalRequest
				return
			}

			i := stream.sequenceIndex

			if index == i+1 {
				// 按序接收到消息
				fmt.Println(parsedResp.GetIndex())

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
	stream.version = time.Now().UnixNano()
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

	//stream.timer = tw.AfterFunc(sys.StreamTimeout, func() {
	//	handleStreamTimeout(identity, stream.version)
	//})

}

// 超时处理函数
func handleStreamTimeout(identity int, version int64) {
	timeOutText := fmt.Sprintf("Stream for identity %d timed out. Cleaning up...\n", identity)
	garbageCollectStream(timeOutText, common.Error, identity, version)
	streams[identity].activeChan <- sys.Timeout
}

func normallyEndStream(identity int, version int64) {
	normallyEndText := fmt.Sprintf("Stream for identity %d normally end. Cleaning up...\n", identity)
	tw.AfterFunc(sys.NormallyEndExpiration, func() {
		garbageCollectStream(normallyEndText, common.Info, identity, version)
		streams[identity].activeChan <- sys.Finish
	})
}

func garbageCollectStream(logText string, logLevel string, identity int, version int64) {
	mu.Lock()
	defer mu.Unlock()

	stream, exists := streams[identity]
	if !exists || stream.version != version {
		//流已经不存在 或者版本不一致 （已经被其他goroutine回收了等等情况）
		//streams[identity].activeChan <- illegalRequest
		//无需处理
		return
	}

	switch logLevel {
	case common.Error:
		log.GetTextLogger().Error(logText)
	case common.Info:
		log.GetTextLogger().Info(logText)
	}

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
