package handler

import (
	"SomersaultCloud/app/somersaultcloud-chat/bootstrap"
	"SomersaultCloud/app/somersaultcloud-chat/constant/common"
	"SomersaultCloud/app/somersaultcloud-chat/domain"
	"SomersaultCloud/app/somersaultcloud-chat/internal/requtil"
	"net/http"
)

type OpenaiVisionModelExecutor struct {
	res *bootstrap.Channels
	env *bootstrap.Env
}

func (o OpenaiVisionModelExecutor) AssemblePrompt(tc *domain.AskContextData) *domain.Message {
	var msgs []domain.TypeMessage
	var sysContent []domain.TypeInfo
	sysContent = append(sysContent,
		&domain.TextType{
			Type: common.TextType,
			Text: tc.SysPrompt,
		})
	sysMsg := &domain.TypeMessage{
		Role:    common.SystemRole,
		Content: sysContent,
	}

	var userContent []domain.TypeInfo
	userContent = append(userContent,
		&domain.ImageUrlType{
			Type: common.ImageURLType,
			ImageUrl: &domain.ImageUrl{
				Url:    tc.Message,
				Detail: common.HighDetail,
			},
		})
	userMsg := &domain.TypeMessage{
		Role:    common.UserRole,
		Content: userContent,
	}
	msgs = append(msgs, *sysMsg, *userMsg)
	return &domain.Message{TypeMessage: &msgs}
}

//TODO 图片访问TBD

func (o OpenaiVisionModelExecutor) EncodeReq(tc *domain.AskContextData) *http.Request {
	//TODO implement me
	panic("implement me")
}

// {'choices': [{'index': 0, 'message': {'role': 'assistant', 'content': '{"题目解析":{"题意解释":"求解给定函数在x趋近于0时的极限值。","已知条件":"题目中给出函数的表达式是一个复杂分式，涉及指数函数、对数函数和三次方根。","解决目标":"通过分析和计算，求出极限值。"},"划分题型":"极限问题","公式定理":"使用洛必达法则和极限定义","计算和推导":[{"解释":"首先，分析分子：(1+x)^x- e^{\\\\cos x}。我们知道(1+x)^x在x趋近于0时，可以用e^x来逼近，即(1+x)^x \\\\approx e^x。此外，e^{\\\\cos x}在x趋近于0时趋近于e。","输出":"分子趋近于0。"},{"解释":"接下来，分析分母：\\\\sqrt[3]{1+x-1} = \\\\sqrt[3]{x}，在x趋近于0时，也趋近于0。","输出":"分母趋近于0。"},{"解释":"由于分子和分母均趋于0，可以用洛必达法则。首先求导数：","输出":"分子对x求导：(x+1)^{-1+x} + (x+1)^x \\\\ln(x+1) \\\\text{ 和 } -e^{\\\\cos x} \\\\sin x。"},{"解释":"分母对x求导：\\\\frac{1}{3}(1+x-1)^{-2/3}。","输出":"进行洛必达法则，不断求导后得出极限值。"}],"总结与归纳":"应用洛必达法则可以帮助我们在不确定形式下找到极限值，这一过程说明了在复杂函数组合下应用工具的重要性。"}', 'refusal': None}, 'logprobs': None, 'finish_reason': 'stop'}], 'usage': {'prompt_tokens': 1378, 'completion_tokens': 403, 'total_tokens': 1781}, 'system_fingerprint': 'fp_cbac5eb3c0'}
type T struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string      `json:"role"`
			Content string      `json:"content"`
			Refusal interface{} `json:"refusal"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

func (o OpenaiVisionModelExecutor) ConfigureProxy(tc *domain.AskContextData) *http.Client {
	return requtil.SetProxy()
}

func (o OpenaiVisionModelExecutor) Execute(tc *domain.AskContextData) {
	//TODO implement me
	panic("implement me")
}

func (o OpenaiVisionModelExecutor) ParseResp(tc *domain.AskContextData) (domain.ParsedResponse, string) {
	//TODO implement me
	panic("implement me")
}
