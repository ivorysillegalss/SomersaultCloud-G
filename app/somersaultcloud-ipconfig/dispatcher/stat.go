package dispatcher

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
	ConnectNum   float64
	MessageBytes float64
}

// CalculateScore 计算的时候 以GB为单位进行计算 因为空间资源的相对充裕 小差别可以不算在内
func (s *Stat) CalculateScore() float64 {
	return getGB(s.MessageBytes)
}

func getGB(m float64) float64 {
	return decimal(m / (1 << 30))
}
func decimal(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}

func (s *Stat) Avg(num float64) {
	s.ConnectNum /= num
	s.MessageBytes /= num
}

func (s *Stat) Clone() *Stat {
	newStat := &Stat{
		MessageBytes: s.MessageBytes,
		ConnectNum:   s.ConnectNum,
	}
	return newStat
}

func (s *Stat) Add(st *Stat) {
	if st == nil {
		return
	}
	s.ConnectNum += st.ConnectNum
	s.MessageBytes += st.MessageBytes
}

func (s *Stat) Sub(st *Stat) {
	if st == nil {
		return
	}
	s.ConnectNum -= st.ConnectNum
	s.MessageBytes -= st.MessageBytes
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
