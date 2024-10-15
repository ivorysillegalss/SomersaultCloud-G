package repository

import (
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/domain"
	"SomersaultCloud/infrastructure/log"
	"SomersaultCloud/infrastructure/redis"
	"context"
	_ "embed"
	jsoniter "github.com/json-iterator/go"
	"github.com/thoas/go-funk"
	"strconv"
)

var chatGenerationMap map[int]*domain.GenerationResponse
var chatStreamValue map[int]chan domain.ParsedResponse

//go:embed lua/hash_expired.lua
var hashExpiredLuaScript string

type generationRepository struct {
	//streamValue chan domain.ParsedResponse
	rcl redis.Client
}

func (g generationRepository) CacheLuaPollHistory(ctx context.Context, generationResp domain.GenerationResponse) {
	//script, _ := ioutil.LoadLuaScript("repository/lua/hash_expired.lua")
	script := hashExpiredLuaScript

	//JSON序列化存储 也许可以改进
	marshal, _ := jsoniter.Marshal(generationResp)

	err, _ := g.rcl.ExecuteArgsLuaScript(context.Background(), script, []string{cache.ChatGeneration, cache.ChatGenerationExpired}, strconv.Itoa(generationResp.ChatId), marshal, cache.ChatGenerationTTL)
	if err != nil {
		log.GetJsonLogger().WithFields("lua", err.Error()).Error("CacheLuaPollHistory Lua executing error")
	}
}

func (g generationRepository) InMemoryPollHistory(ctx context.Context, response *domain.GenerationResponse) {
	if funk.IsEmpty(chatGenerationMap) {
		chatGenerationMap = make(map[int]*domain.GenerationResponse)
	}
	chatGenerationMap[response.ChatId] = response
}

func (g generationRepository) InMemorySetStreamValue(ctx context.Context, response domain.ParsedResponse) {
	identity := response.GetIdentity()

	chatStreamValue[identity] <- response
}

func (g generationRepository) InMemoryGetStreamValue(userId int) chan domain.ParsedResponse {
	responses := chatStreamValue[userId]
	return responses
}

//func (g generationRepository) GetStreamChannel() chan domain.ParsedResponse {
//	return g.streamValue
//}
//
//func (g generationRepository) SendStreamValueChannel(response domain.ParsedResponse) {
//	value := g.streamValue
//	if funk.IsEmpty(value) {
//		//TODO 常量替换
//		value = make(chan domain.ParsedResponse, 100)
//	}
//	value <- response
//}

func NewGenerationRepository(dbs *bootstrap.Databases) domain.GenerationRepository {
	return &generationRepository{rcl: dbs.Redis}
}
