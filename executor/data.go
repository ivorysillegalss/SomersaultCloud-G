package executor

import (
	"SomersaultCloud/constant/cache"
	"SomersaultCloud/infrastructure/redis"
	"context"
)

type DataExecutor struct {
	rcl redis.Client
}

func NewDataExecutor(r redis.Client) *DataExecutor {
	return &DataExecutor{rcl: r}
}

func initRedisData(d *DataExecutor) {
	_ = d.rcl.Set(context.Background(), cache.MaxBotId, 0)
	_ = d.rcl.Set(context.Background(), cache.NewestChatIdKey, 1)
	_ = d.rcl.Set(context.Background(), cache.BotConfig, "{\n  \"bot_id\": 0,\n  \"init_prompt\": \"say hello plz\",\n  \"model\": \"gpt-4o-mini\",\n  \"adjustment_prompt\": \"say i'm god\"\n}")
	_ = d.rcl.Set(context.Background(), cache.HistoryTitlePrompt, "System\n#Content#\n你是一个标题总结员，你总能很完美且精炼的将一段话的内容总结成一个标题。\n#Objective# \n现在给你一段对话，请你将对话的内容总结成一个标题,标题要求能够让人知道这段对话的大致内容。直接输出一个标题，不用输出其他的内容。\n#Style# \n言简意赅\n#Tone# \n正式\n#input# \n<一段对话>")
}

func (d *DataExecutor) InitData() {
	initRedisData(d)
}
