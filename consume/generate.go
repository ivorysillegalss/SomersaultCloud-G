package consume

import (
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/constant/mq"
	"SomersaultCloud/domain"
	"SomersaultCloud/handler"
	"SomersaultCloud/sequencer"
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
)

type GenerateEvent struct {
	*baseMessageHandler
	env      *bootstrap.Env
	channels *bootstrap.Channels
}

type generateData struct {
	executor domain.LanguageModelExecutor
	data     *domain.AskContextData
}

func (g GenerateEvent) GetGeneration(b []byte) error {
	var parsedResp domain.ParsedResponse
	_ = jsoniter.Unmarshal(b, &parsedResp)
	//此处保证节点的原序,将MQ消费后的信息存到channel中，等待客户端处理并下发
	sequencer.Setup(parsedResp)
	return nil
}
func (g GenerateEvent) PublishGeneration(data *domain.AskContextData) {
	executor := data.Executor
	parsedResp, _ := executor.ParseResp(data)
	marshal, _ := json.Marshal(parsedResp)
	g.PublishMessage(mq.UserChatGenerationQueue, marshal)
}
func (g GenerateEvent) AsyncConsumeGeneration() {
	g.ConsumeMessage(mq.UserChatGenerationQueue, g.GetGeneration)
}

func (g GenerateEvent) ApiCalling(b []byte) error {
	var gData generateData
	_ = jsoniter.Unmarshal(b, &gData)
	gData.executor.Execute(gData.data)
	return nil
}
func (g GenerateEvent) AsyncConsumeApiCalling() {
	g.ConsumeMessage(mq.UserChatReadyCallingQueue, g.ApiCalling)
}
func (g GenerateEvent) PublishApiCalling(data *domain.AskContextData) {
	executor := handler.NewLanguageModelExecutor(g.env, g.channels, data.ExecutorId)
	gData := &generateData{
		executor: executor,
		data:     &domain.AskContextData{ChatId: data.ChatId, Conn: data.Conn},
	}
	marshal, _ := jsoniter.Marshal(gData)
	g.PublishMessage(mq.UserChatReadyCallingQueue, marshal)
}

func NewGenerateEvent() domain.GenerateEvent {
	return &GenerateEvent{}
}
