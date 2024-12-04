package repository

import (
	"SomersaultCloud/app/bootstrap"
	"SomersaultCloud/app/constant/cache"
	"SomersaultCloud/app/constant/common"
	"SomersaultCloud/app/domain"
	"SomersaultCloud/app/infrastructure/log"
	"SomersaultCloud/app/infrastructure/redis"
	"context"
	_ "embed"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

var chatGenerationMap = make(map[int]*domain.GenerationResponse)

//go:embed lua/hash_expired.lua
var hashExpiredLuaScript string

type generationRepository struct {
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
	chatGenerationMap[response.ChatId] = response
}

func (g generationRepository) ReadyStreamDataStorage(ctx context.Context, ready domain.StreamGenerationReadyStorageData) {
	_ = g.rcl.SetExpire(ctx, cache.StreamStorageReadyData+common.Infix+strconv.Itoa(ready.UserId), ready, cache.StreamStorageReadyDataExpire)
}

func (g generationRepository) GetStreamDataStorage(ctx context.Context, userId int) *domain.AskContextData {
	get, err := g.rcl.Get(ctx, cache.StreamStorageReadyData+common.Infix+strconv.Itoa(userId))
	if g.rcl.IsEmpty(err) {
		log.GetTextLogger().Error("can't find target id stream data cache , with userId: " + strconv.Itoa(userId))
		return nil
	}
	var dataReady domain.StreamGenerationReadyStorageData
	_ = jsoniter.Unmarshal([]byte(get), dataReady)
	return &domain.AskContextData{UserId: userId, ChatId: dataReady.ChatId, BotId: dataReady.BotId, Message: dataReady.UserContent, History: dataReady.Records}
}

func NewGenerationRepository(dbs *bootstrap.Databases) domain.GenerationRepository {
	return &generationRepository{rcl: dbs.Redis}
}
