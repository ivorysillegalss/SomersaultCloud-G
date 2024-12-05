package domain

const (
	windowSize = 5
)

// 记录在滑窗范围(windowSize) 内的总分值 并且根据这个值计算出总分
type stateWindow struct {
	//存储状态信息
	stateQueue []*Stat
	//通过此chan更新信息 当有新的Stat传进来 则通过此方法和appendStat更新状态
	statChan chan *Stat
	sumStat  *Stat // 记录在窗口内的状态值总分
	idx      int64
}

func newStateWindow() *stateWindow {
	return &stateWindow{
		stateQueue: make([]*Stat, windowSize),
		statChan:   make(chan *Stat),
		sumStat:    &Stat{},
	}
}

func (sw *stateWindow) getStat() *Stat {
	res := sw.sumStat.Clone()
	res.Avg(windowSize)
	return res
}

// appendStat 滑动窗口增加减分
//
//	注意仅删除窗口范围内对应的分值 其状态信息仍然存储
func (sw *stateWindow) appendStat(s *Stat) {
	// 减去即将被删除的state
	sw.sumStat.Sub(sw.stateQueue[sw.idx%windowSize])
	// 更新最新的stat
	sw.stateQueue[sw.idx%windowSize] = s
	// 计算最新的窗口和
	sw.sumStat.Add(s)
	sw.idx++
}
