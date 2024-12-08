package monitor

import (
	"SomersaultCloud/app/somersaultcloud-common/log"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

// MetricsCache 需要进行缓存的数据
type MetricsCache struct {
	RequestCount    float64
	RequestDuration float64
}

type Exporter struct {
	Cache *MetricsCache
}

func (e *Exporter) UpdateCacheFromRegistry() {
	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		log.GetTextLogger().Error("Error gathering metrics: %v", err)
		return
	}

	for _, mf := range metricFamilies {
		for _, metric := range mf.Metric {
			if *mf.Name == "http_requests_total" {
				//TODO 拓宽请求类型
				// 此处先假设我们只处理 GET 请求的计数
				// 确保标签数组的长度足够，防止越界错误
				if len(metric.Label) > 0 && metric.Label[0].GetValue() == http.MethodGet {
					// 获取请求计数
					e.Cache.RequestCount = metric.Counter.GetValue()
				}
			} else if *mf.Name == "http_request_duration_seconds_sum" {
				e.Cache.RequestDuration = metric.Summary.GetSampleSum()
			}
		}
	}
}
