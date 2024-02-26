package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"io"
	"log"
	"mini-gpt/constant"
	exception "mini-gpt/error"
	"mini-gpt/models"
	"mini-gpt/setting"
	"mini-gpt/utils/redisUtils"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

//var logger = setting.GetLogger()
//正常来说应该全局变量的 但是由于代码的先后执行问题先放到下面的函数中

// 优化处理prompt 是否合法
func optimizationPrompt(originalPrompt string) (reply string) {
	originalPrompt = strings.TrimSpace(originalPrompt)
	length := len([]rune(originalPrompt))
	reply = originalPrompt
	if length <= 1 {
		reply = "请说详细些..."
		return
	}
	if length > constant.DefaultMaxToken {
		reply = "问题字数超出设定限制，请精简问题"
		return
	}

	return
}

func completions(userChat *models.UserChat, respBody io.ReadCloser) {
	// 异步读取数据
	//TODO 暂且改回同步 运行有未知问题
	//go func(u *models.UserChat, respBody io.ReadCloser) {
	u := userChat
	defer func() {
		u.Answer.Counter++
		//更新回答的数量 保持问问题的数量与回答的数量一致

		_ = respBody.Close()
		if err := recover(); err != nil {
			_, file, line, _ := runtime.Caller(3)
			log.Println("ERROR:", err, file, line)
		}
	}()

	u.Answer.Buffer.Reset()
	//清空缓冲区

	u.Answer.Mu.Lock()
	//加锁防止竞争
	//if u.Question.Counter-u.Answer.Counter == 1 {
	//如果问题的数量多于答案 则等回答
	_, err3 := io.Copy(&u.Answer.Buffer, respBody)
	//返回的第一个为复制的字节总数
	if err3 != nil {
		return
	}
	//} else {
	//这里有异常可以进一步通过日志等方式封装 TODO
	//return
	//}
	u.Answer.Mu.Unlock()

	//}(userChat, respBody)
}

func getUser(userId string) *models.UserChat {

	/*注意！
	这里给redis传递的 是对应结构体UserChat的指针！ 不是整个结构体
	因为在UserChat结构体中
	Mu      sync.Mutex
	是互斥锁 而复制一个互斥锁是不推荐的 于是直接传递锁的指针
	*/
	userChat, err2 := redisUtils.GetStruct[*models.UserChat](constant.UserCachePrefix + userId)
	//go redis中 如果get不到想要的值 是会返回一个如下的错 redis.Nil 通过进行此异常处理 来判断用户此前是否有数据缓存在redis中
	//如果没有 则创建一个新的
	if errors.Is(redis.Nil, err2) {
		newUserChat := models.NewUserChat(userId)
		_ = redisUtils.SetStruct(constant.UserCachePrefix+userId, newUserChat)
		return newUserChat
		//这里应该也不会有错了 直接不处理了
	}
	//但是问题来了 这种做法是正确的吗？并发的效率性和安全性可能都不如sync.Map
	return userChat
}

// 优化 异步读取数据流
func Execute(uid string, message *models.ApiRequestMessage) (*models.CompletionResponse, error) {
	//应该在外部传入调用者的标识逻辑 UserChat等
	var logger = setting.GetLogger()

	//ctx, cancel := context.WithTimeout(context.Background(), constant.DefaultMaxLimitedTime)
	//设置超时时间限制 若一次请求超出此事件则驳回时间 取消掉
	//defer cancel()

	message.InputPrompt = optimizationPrompt(message.InputPrompt)
	//优化prompt 这里这么写可能不太好 后期可能需要优化 （根据不同的情况决定是否需要对prompt进行优化）

	userChat := getUser(uid)

	if userChat.Question.Doing {
		return models.ExceptionCompletionResponse("上个问题正在处理中，请稍等..."), nil
		//这里需要再次丰富一下对应的异常类型 太多了太繁杂了
	}

	//更新一下用户的状态 相当于上锁
	userChat.Question.Doing = true
	defer func() {
		userChat.Question.Doing = false
	}()

	var done bool
	//done表示用户此次 执行的结果是否完成
	for !done {
		select {
		// 超时结束
		//case <-ctx.Done():
		//	done = true
		// 每秒检测结果是否完全返回，用于提前结束
		default:
			if userChat.Question.Counter == userChat.Answer.Counter {
				done = true
			}
		}
	}

	//构造请求头
	client := setProxy()
	req, err := encodeReq(message)
	if err != nil {
		logger.Error(err)
		return models.ErrorCompletionResponse(), err
	}

	//发送请求
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return models.ErrorCompletionResponse(), err
	}

	//判断状态是否正常运行
	executeStatus := statusClarify(resp)
	if executeStatus.StatusCode != constant.APIExecuteSuccessStatus {
		return models.ErrorCompletionResponse(), &exception.ExecuteError{
			ExecuteTime: time.Now(),
			Status:      executeStatus.Status,
			StatusCode:  executeStatus.StatusCode,
		}
	}

	// 异步读取数据 同时将数据传输到对应用户的缓冲区当中
	// 先改回同步
	completions(userChat, resp.Body)

	userChat.Answer.Mu.Lock()
	defer userChat.Answer.Mu.Unlock()

	//读取之前异步写下的结果 （如果已经处理完了的话）此时的数据是一个json字符串
	chatJson := userChat.Answer.Buffer.String()
	userChat.Answer.Buffer.Reset()

	completionResponse, err := decodeResp(chatJson)
	//反序列化
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
func encodeReq(reqMessage *models.ApiRequestMessage) (*http.Request, error) {

	//这个判空的过程可以优化在结构体 models层中
	if reqMessage.MaxToken == 0 {
		reqMessage.MaxToken = constant.DefaultMaxToken
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
	req, err := http.NewRequest("POST", constant.ApiServerOpenAI, bytes.NewBuffer(jsonData))
	if err != nil {
		//logger.Error(err)
		return new(http.Request), err
	}
	secretKey := setting.Conf.ApiConfig.SecretKey
	//organizationID := "org-KQEFla180NCNd60se8ecJNc7"
	//req.Header.Set("OpenAI-Organization", organizationID)
	//organizationID 为可选项 有时候需要加上才运行成功 截至commit 不加也行
	key := fmt.Sprintf("Bearer %s", secretKey)
	req.Header.Set("Authorization", key) // 请确保使用你自己的API密钥
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// 解码响应信息
func decodeResp(chatJson string) (*models.CompletionResponse, error) {
	reader := strings.NewReader(chatJson)
	// 创建json.Decoder实例
	decoder := json.NewDecoder(reader)
	// 创建一个变量，用于存储解码后的数据
	var completionResponse *models.CompletionResponse
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

// 判断状态
func statusClarify(resp *http.Response) *models.ExecuteStatus {
	return &models.ExecuteStatus{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
	}
}
