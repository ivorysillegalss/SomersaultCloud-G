package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mini-gpt/models"
	"mini-gpt/setting"
	"net/http"
	"net/url"
)

//var logger = setting.GetLogger()
//正常来说应该全局变量的 但是由于代码的先后执行问题先放到下面的函数中

func Execute(inputPrompt string) (models.CompletionResponse, error) {

	var logger = setting.GetLogger()

	// 解析代理服务器的URL
	proxyURL, err := url.Parse("http://localhost:7890")
	if err != nil {
		logger.Error(err)
		return models.CompletionResponse{}, err
	}

	// 创建一个新的HTTP客户端，配置代理
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	//// 输入提示
	//scanner := bufio.NewScanner(os.Stdin)
	//scanner.Scan()
	//inputPrompt := scanner.Text()

	// 构造请求体
	data := models.CompletionRequest{
		Model:     "gpt-3.5-turbo-instruct", // 替换为当前可用的模型
		Prompt:    inputPrompt,
		MaxTokens: 1000,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error(err)
		return models.CompletionResponse{}, err
	}

	// 发送请求
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error(err)
		return models.CompletionResponse{}, err
	}
	secretKey := setting.Conf.ApiConfig.SecretKey
	key := fmt.Sprintf("Bearer %s", secretKey)
	req.Header.Set("Authorization", key) // 请确保使用你自己的API密钥
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return models.CompletionResponse{}, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	// 创建一个变量，用于存储解码后的数据
	var completionResponse models.CompletionResponse
	// 循环解码每个JSON对象
	for {
		if err := decoder.Decode(&completionResponse); err == io.EOF {
			// io.EOF 表示已经到达输入流的末尾
			break
		} else if err != nil {
			fmt.Println("Error decoding JSON:", err)
			logger.Error("Error decoding JSON:", err)
			return completionResponse, err
		}

	}
	return completionResponse, nil
}
