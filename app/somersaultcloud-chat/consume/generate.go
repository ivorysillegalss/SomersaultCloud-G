package consume

import (
	"SomersaultCloud/app/somersaultcloud-chat/bootstrap"
	"SomersaultCloud/app/somersaultcloud-chat/constant/mq"
	"SomersaultCloud/app/somersaultcloud-chat/domain"
	"SomersaultCloud/app/somersaultcloud-common/log"
	"SomersaultCloud/app/somersaultcloud-common/set"
	"context"
	jsoniter "github.com/json-iterator/go"
)

var retrySet set.HashSet

func init() {
	retrySet = *set.NewHashSet()
}

type GenerateEvent struct {
	*baseMessageHandler
	env                  *bootstrap.Env
	channels             *bootstrap.Channels
	generationRepository domain.GenerationRepository
}

func (g GenerateEvent) StreamDataReady(b []byte) error {
	var dataReady domain.StreamGenerationReadyStorageData
	_ = jsoniter.Unmarshal(b, dataReady)
	g.generationRepository.ReadyStreamDataStorage(context.Background(), dataReady)
	return nil
}

func (g GenerateEvent) AsyncStreamStorageDataReady() {
	g.ConsumeMessage(mq.UserChatReadyCallingQueue, g.StreamDataReady)
}

func (g GenerateEvent) PublishStreamReadyStorageData(data *domain.StreamGenerationReadyStorageData) {
	//此方法应用于流信息 发起调用前 所以此时没generateText
	marshal, _ := jsoniter.Marshal(data)
	g.PublishMessage(mq.UserChatReadyCallingQueue, marshal)
}

func (g GenerateEvent) PublishChatGenerate(data *domain.AskContextData) {
	marshal, _ := jsoniter.Marshal(data)
	g.PublishMessage(mq.UserChatReadyCallingQueue, marshal)
}

// TODO askContextData序列化优化
func (g GenerateEvent) ConsumeChatGenerate() {
	g.ConsumeMessage(mq.UserChatReadyCallingQueue, g.DoGenerate)
}

func (g GenerateEvent) DoGenerate(b []byte) error {
	var data domain.AskContextData
	_ = jsoniter.Unmarshal(b, data)
	execute := data.Executor.Execute(&data)

	//return error的时候 会自动重试
	if !execute {
		//通过一个hashset 判断有没有重试过 如果有的话 删除 （用户颗粒度）
		//TODO bitmap&布隆优化

		//已经重试过的话 丢弃
		if retrySet.Contains(data.UserId) {
			retrySet.Clear()
			return nil
		}
		retrySet.Add(data.UserId)

		c := &consumeErr{err: "message reject"}
		log.GetTextLogger().Error(c.Error()+"%v", data)
		return c
	}
	return nil
}

type consumeErr struct {
	err string
}

func (c *consumeErr) Error() string {
	return c.err
}

func NewGenerateEvent(h MessageHandler, e *bootstrap.Env, c *bootstrap.Channels, g domain.GenerationRepository) domain.GenerateEvent {
	messageHandler := h.(*baseMessageHandler)
	chatReadyCalling := &MessageQueueArgs{
		ExchangeName:         mq.UserChatReadyCallingExchange,
		QueueName:            mq.UserChatReadyCallingQueue,
		KeyName:              mq.UserChatReadyCallingKey,
		ExistDeadLetterQueue: true,
		DeadLetterExchange:   mq.UserChatDeadLetterRetryExchange,
		DeadLetterRoutingKey: mq.UserChatDeadLetterRetryKey,
	}
	//TODO 死信队列设置代码抽离
	chatDeadLetterRetry := &MessageQueueArgs{
		ExchangeName:         mq.UserChatDeadLetterRetryExchange,
		QueueName:            mq.UserChatDeadLetterRetryQueue,
		KeyName:              mq.UserChatDeadLetterRetryKey,
		ExistDeadLetterQueue: true,
		DeadLetterExchange:   mq.UserChatDeadLetterRetryExchange,
		DeadLetterRoutingKey: mq.UserChatDeadLetterRetryKey,
	}
	messageHandler.InitMessageQueue(chatReadyCalling, chatDeadLetterRetry)
	return &GenerateEvent{
		baseMessageHandler:   messageHandler,
		env:                  e,
		channels:             c,
		generationRepository: g,
	}
}
