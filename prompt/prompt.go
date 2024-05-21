package prompt

import (
	"fmt"
	"mini-gpt/constant"
	"mini-gpt/utils/redisUtils"
	"mini-gpt/utils/stringsUtil"
)

// OpenaiPrompt记录部分需要高频用到的提示词
// eg 总结标题
var OpenaiPrompt map[string]string

// 初始化对应的map
func init() {
	OpenaiPrompt = make(map[string]string)
}

// 加载总结标题的prompt
func loadOpenaiTitlePrompt() {
	TitlePromptValue, _ := redisUtils.Get(stringsUtil.Concat(constant.OpenaiPrompt, constant.Infix, constant.Conclude2TitlePrompt))
	OpenaiPrompt[constant.Conclude2TitlePrompt] = TitlePromptValue
}

// 加载所有prompt的入口方法
func LoadPrompt() {
	loadOpenaiTitlePrompt()
}

// 将对应的提示词存入redis中
func SetOpenaiTitlePrompt() {
	var conclude2TitlePrompt string
	conclude2TitlePrompt = "#Content#\n你是一个标题总结员，你总能很完美且精炼的将一段话的内容总结成一个标题。\n" +
		"#Objective# \n现在给你一段对话，请你将对话的内容总结成一个标题,标题要求能够让人知道这段对话的大致内容。直接输出一个标题，不用输出其他的内容。" +
		"\n#Style# \n言简意赅\n#Tone# \n正式\n#input# \n<一段对话>"
	err := redisUtils.Set(stringsUtil.Concat(constant.OpenaiPrompt, constant.Infix, constant.Conclude2TitlePrompt), conclude2TitlePrompt)
	fmt.Println(err)
}
