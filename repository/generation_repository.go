package repository

import (
	"SomersaultCloud/bootstrap"
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/domain"
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
	script, err := ioutil.LoadLuaScript("cron/lua/hash_expired.lua")
	if err != nil {
		//TODO 打日志
	}

	//JSON序列化存储 也许可以改进
	//marshal, _ := json.Marshal(generationResp)
	marshal, _ := jsoniter.Marshal(generationResp)
	//TODO json包有问题？为什么明明不是空的序列化出来是空的。
	//http.response不可以序列化

	err, _ = g.rcl.ExecuteArgsLuaScript(context.Background(), script, []string{cache.ChatGeneration, cache.ChatGenerationExpired}, strconv.Itoa(generationResp.ChatId), marshal, cache.ChatGenerationTTL)
	if err != nil {
		//同上 TODO 打日志
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
