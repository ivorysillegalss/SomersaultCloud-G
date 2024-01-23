package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mini-gpt/code"
	"mini-gpt/models"
	"mini-gpt/setting"
	"net/http"
	"net/url"
)

//var logger = setting.GetLogger()
//正常来说应该全局变量的 但是由于代码的先后执行问题先放到下面的函数中

func Execute(message models.ApiRequestMessage) (models.CompletionResponse, error) {
	var logger = setting.GetLogger()

	//构造请求头
	client := setProxy()
	req, err := encodeReq(message)
	if err != nil {
		logger.Error(err)
		return models.CompletionResponse{}, err
	}

	//发送请求
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return models.CompletionResponse{}, err
	}
	defer resp.Body.Close()

	//处理响应信息
	completionResponse, err := decodeResp(resp)
	if err != nil {
		logger.Error("Error decoding JSON:", err)
	}
	return completionResponse, nil
}

// 设置网络代理相关
func setProxy() *http.Client {
	// 解析代理服务器的URL
	proxyURL, err := url.Parse("http://localhost:7890")
	if err != nil {
	}

	// 创建一个新的HTTP客户端，配置代理
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}
	return client
}

// 构造请求头
func encodeReq(reqMessage models.ApiRequestMessage) (*http.Request, error) {

	//这个判空的过程可以优化在结构体 models层中
	if reqMessage.MaxToken != 0 {
		reqMessage.MaxToken = code.DefaultMaxToken
	}

	// 构造请求体
	data := models.CompletionRequest{
		Model: reqMessage.Model,
		//Model:     "gpt-3.5-turbo-instruct", // 替换为当前可用的模型
		Prompt:    reqMessage.InputPrompt,
		MaxTokens: reqMessage.MaxToken,
	}
	jsonData, err := json.Marshal(data)

	if err != nil {
		//logger.Error(err)
		return new(http.Request), err
	}

	// 发送请求
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		//logger.Error(err)
		return new(http.Request), err
	}
	secretKey := setting.Conf.ApiConfig.SecretKey
	//organizationID := "org-KQEFla180NCNd60se8ecJNc7"
	//req.Header.Set("OpenAI-Organization", organizationID)
	//organnizationID 为可选项 有时候需要加上才运行成功 截至commit 不加也行
	key := fmt.Sprintf("Bearer %s", secretKey)
	req.Header.Set("Authorization", key) // 请确保使用你自己的API密钥
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// 解码响应信息
func decodeResp(resp *http.Response) (models.CompletionResponse, error) {
	decoder := json.NewDecoder(resp.Body)
	// 创建一个变量，用于存储解码后的数据
	var completionResponse models.CompletionResponse
	for {
		if err := decoder.Decode(&completionResponse); err == io.EOF {
			// io.EOF 表示已经到达输入流的末尾
			break
		} else if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return completionResponse, err
		}
	}
	return completionResponse, nil
}
