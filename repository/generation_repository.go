package repository

import (
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/domain"
	"SomersaultCloud/infrastructure/log"
	"SomersaultCloud/infrastructure/redis"
	"SomersaultCloud/internal/ioutil"
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/thoas/go-funk"
	"strconv"
)

var chatGenerationMap map[int]*domain.GenerationResponse

type generationRepository struct {
	rcl redis.Client
}

func (g generationRepository) CacheLuaPollHistory(ctx context.Context, generationResp domain.GenerationResponse) {
	script, _ := ioutil.LoadLuaScript("cron/lua/hash_expired.lua")

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

func NewGenerationRepository(dbs *bootstrap.Databases) domain.GenerationRepository {
	return &generationRepository{rcl: dbs.Redis}
}
