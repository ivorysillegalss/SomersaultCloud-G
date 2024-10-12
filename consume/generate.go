package consume

import (
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/constant/mq"
	"SomersaultCloud/constant/task"
	"SomersaultCloud/domain"
	"SomersaultCloud/handler"
	"SomersaultCloud/infrastructure/log"
	"SomersaultCloud/sequencer"
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"github.com/thoas/go-funk"
	"strconv"
)

type GenerateEvent struct {
	*baseMessageHandler
	env      *bootstrap.Env
	channels *bootstrap.Channels
}

type generateData struct {
	ChatId     int
	ExecutorId int
}

// TODO 屎山，改进赋值新结构体逻辑
func getParseType(parsedResp map[string]any) domain.ParsedResponse {
	userId := int(parsedResp["UserId"].(float64))
	generateText := parsedResp["GenerateText"].(string)
	finishReason := parsedResp["FinishReason"].(string)
	index := int(parsedResp["Index"].(float64))
	executorId := int(parsedResp["ExecutorId"].(float64))
	chatcmplId := parsedResp["ChatcmplId"].(string)
	switch executorId {
	case task.ChatAskExecutorId:
		return &domain.OpenAIParsedResponse{UserId: userId, GenerateText: generateText, FinishReason: finishReason, Index: index, ChatcmplId: chatcmplId}
	default:
		log.GetTextLogger().Error("GET wrong parse type for UserId:" + strconv.Itoa(userId) + "with executorId : " + strconv.Itoa(executorId))
		return nil
	}
}

// GetGeneration 不能直接将数据放到接口类型当中 需提前将接口实例化或者确定她的具体类型
func (g GenerateEvent) GetGeneration(b []byte) error {
	var parsedResp map[string]any
	_ = jsoniter.Unmarshal(b, &parsedResp)
	//此处保证节点的原序,将MQ消费后的信息存到channel中，等待客户端处理并下发
	newSequencer := sequencer.NewSequencer()
	newSequencer.Setup(getParseType(parsedResp))
	return nil
}

// TODO 直接序列化接口可行？
func (g GenerateEvent) PublishGeneration(data domain.ParsedResponse) {
	marshal, _ := json.Marshal(data)
	g.PublishMessage(mq.UserChatGenerationQueue, marshal)
}
func (g GenerateEvent) AsyncConsumeGeneration() {
	g.ConsumeMessage(mq.UserChatGenerationQueue, g.GetGeneration)
}

// ApiCalling TODO REMOVEORUPDATE好像用不了
func (g GenerateEvent) ApiCalling(b []byte) error {
	var gData generateData
	_ = jsoniter.Unmarshal(b, &gData)
	if funk.IsEmpty(gData) {
		log.GetTextLogger().Error("data for api calling is nil")
		return nil
	}
	executor := handler.NewLanguageModelExecutor(g.env, g.channels, gData.ExecutorId)
	panic(executor)
	//executor.Execute(&domain.AskContextData{Executor: executor,ChatId: gData.ChatId,Conn: domain.NewConnection()})
	return nil
}

// TODO remove
func (g GenerateEvent) AsyncConsumeApiCalling() {
	g.ConsumeMessage(mq.UserChatReadyCallingQueue, g.ApiCalling)
}

// TODO remove
func (g GenerateEvent) PublishApiCalling(data *domain.AskContextData) {
	gData := &generateData{
		ChatId:     data.ChatId,
		ExecutorId: data.ExecutorId,
	}
	marshal, _ := jsoniter.Marshal(gData)
	g.PublishMessage(mq.UserChatReadyCallingQueue, marshal)
}

func NewGenerateEvent(h MessageHandler, e *bootstrap.Env, c *bootstrap.Channels) domain.GenerateEvent {
	messageHandler := h.(*baseMessageHandler)
	//chatReadyCalling := &MessageQueueArgs{
	//	ExchangeName: mq.UserChatReadyCallingExchange,
	//	QueueName:    mq.UserChatReadyCallingQueue,
	//	KeyName:      mq.UserChatReadyCallingKey,
	//}
	chatGetGeneration := &MessageQueueArgs{
		ExchangeName: mq.UserChatGenerationExchange,
		QueueName:    mq.UserChatGenerationQueue,
		KeyName:      mq.UserChatGenerationKey,
	}
	messageHandler.InitMessageQueue(chatGetGeneration)
	//messageHandler.InitMessageQueue(chatGetGeneration, chatReadyCalling)
	return &GenerateEvent{
		baseMessageHandler: messageHandler,
		env:                e,
		channels:           c,
	}
}
