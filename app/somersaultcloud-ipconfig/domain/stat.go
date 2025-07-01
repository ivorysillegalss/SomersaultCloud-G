package domain

import "math"

// Stat 这里对应之前可拓展map发过来的状态信息 根据这些信息算出一个值
// 再根据这个值对服务的状态进行评价
//
//	同时 这里所记录的是机器的剩余指标信息 因为由于物理机的配置可能不同
//	所以记录剩余指标是一种最直观的衡量机器剩余能力的方法 当然 这个是需要根据机器的最大配置进行计算得出的
//	也就是 哪台机器的这个值大 他们所富余的资源就更多
//
// TODO 算法修改
type Stat struct {
	ConnectNum      float64 // 连接数
	MessageBytes    float64 // 字节数
	AvailableMem    float64 // 可用内存
	CpuIdleTime     float64 // CPU空闲时间
	RequestCount    float64 // 请求数
	RequestDuration float64 // 请求时长
}

// 负载均衡各指标权重
const (
	wCpuIdle      = 2.0 // CPU 空闲时间的权重
	wMemAvailable = 1.0 // 可用内存的权重
	wConnectNum   = 1.0 // 连接数的权重
	wMsgBytes     = 0.5 // 消息字节数的权重
	wRequestCount = 0.8 // 请求数的权重
	wRequestDur   = 0.9 // 请求时长的权重
)

// CalculateScore 计算的时候 以GB为单位进行计算 因为空间资源的相对充裕 小差别可以不算在内
func (s *Stat) CalculateScore() float64 {
	memGb := getGB(s.AvailableMem)

	msgGb := getGB(s.MessageBytes)

	// 分数越高代表负载越低 / 可用度越高，这里选择对“空闲指标”加分，对“占用指标”减分
	// 比如：CPU 空闲 / 可用内存 => 正向加分；连接数 / 请求数等 => 反向减分
	score := 0.0

	// 正向：CPU 空闲时间越高越好，可用内存越多越好
	score += wCpuIdle * s.CpuIdleTime
	score += wMemAvailable * memGb

	// 负向：连接数越多，消息量越多，请求数和请求时长越多，表示负载越重
	score -= wConnectNum * s.ConnectNum
	score -= wMsgBytes * msgGb
	score -= wRequestCount * s.RequestCount
	score -= wRequestDur * s.RequestDuration

	return decimal(score)
}

// 提高颗粒度 （M级别或以下的数值通常来说 没有那么敏感）
func getGB(m float64) float64 {
	return decimal(m / (1 << 30))
}

// 截断精度
func decimal(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}

// Avg 用于对多个 Stat 求平均（如多次采样后，计算综合平均状态）。
func (s *Stat) Avg(num float64) {
	if num == 0 {
		return
	}
	s.ConnectNum /= num
	s.MessageBytes /= num
	s.AvailableMem /= num
	s.CpuIdleTime /= num
	s.RequestCount /= num
	s.RequestDuration /= num
}

// Clone 深拷贝对象，复制所有字段。
func (s *Stat) Clone() *Stat {
	return &Stat{
		ConnectNum:      s.ConnectNum,
		MessageBytes:    s.MessageBytes,
		AvailableMem:    s.AvailableMem,
		CpuIdleTime:     s.CpuIdleTime,
		RequestCount:    s.RequestCount,
		RequestDuration: s.RequestDuration,
	}
}

// Add 叠加汇总stat。
func (s *Stat) Add(st *Stat) {
	if st == nil {
		return
	}
	s.ConnectNum += st.ConnectNum
	s.MessageBytes += st.MessageBytes
	s.AvailableMem += st.AvailableMem
	s.CpuIdleTime += st.CpuIdleTime
	s.RequestCount += st.RequestCount
	s.RequestDuration += st.RequestDuration
}

// Sub 差值计算stat。
func (s *Stat) Sub(st *Stat) {
	if st == nil {
		return
	}
	s.ConnectNum -= st.ConnectNum
	s.MessageBytes -= st.MessageBytes
	s.AvailableMem -= st.AvailableMem
	s.CpuIdleTime -= st.CpuIdleTime
	s.RequestCount -= st.RequestCount
	s.RequestDuration -= st.RequestDuration
}

func min(a, b, c float64) float64 {
	m := func(k, j float64) float64 {
		if k > j {
			return j
		}
		return k
	}
	return m(a, m(b, c))
}

func (s *Stat) CalculateStaticScore() float64 {
	return s.ConnectNum
}
