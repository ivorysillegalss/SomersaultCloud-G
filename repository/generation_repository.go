package repository

import (
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/domain"
	"SomersaultCloud/infrastructure/redis"
	"SomersaultCloud/internal/ioutil"
	"context"
	"encoding/json"
	"strconv"
)

type generationRepository struct {
	rcl redis.Client
}

func (g generationRepository) CacheLuaPollHistory(ctx context.Context, generationResp domain.GenerationResponse) {
	script, err := ioutil.LoadLuaScript("lua/hash_expired.lua")
	if err != nil {
		//TODO 打日志
	}

	//JSON序列化存储 也许可以改进
	marshal, _ := json.Marshal(generationResp)
	err = g.rcl.ExecuteArgsLuaScript(context.Background(), script, []string{cache.ChatGeneration, cache.ChatGenerationExpired}, strconv.Itoa(generationResp.ChatId), string(marshal), cache.ChatGenerationTTL)
	if err != nil {
		//同上 TODO 打日志
	}
}

func NewGenerationRepository(client redis.Client) domain.GenerationRepository {
	return &generationRepository{rcl: client}
}
