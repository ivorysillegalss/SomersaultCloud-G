package requtil

import (
	"net/http"
	"net/url"
)

// SetProxy 设置网络代理相关
func SetProxy() *http.Client {
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
