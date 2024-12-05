package dispatcher

import (
	"SomersaultCloud/app/somersaultcloud-ipconfig/bootstrap"
	"SomersaultCloud/app/somersaultcloud-ipconfig/domain"
	"SomersaultCloud/app/somersaultcloud-ipconfig/source"
	"sort"
	"sync"
)

type Dispatcher struct {
	IpConfigEnv    *bootstrap.IpConfigEnv
	candidateTable map[string]*EndPort
	sync.RWMutex
}

var dp *Dispatcher

func init() {
	dp = &Dispatcher{}
	dp.candidateTable = make(map[string]*EndPort)
}

func (dp *Dispatcher) Handle() {
	go func() {
		for event := range source.EventChan() {
			switch event.Type {
			case source.AddNodeEvent:
				dp.addNode(event)
			case source.DelNodeEvent:
				dp.delNode(event)
			}
		}
	}()
}

func (dp *Dispatcher) Do(ctx *domain.IpConfContext) []*EndPort {
	//获得所有候选的endports
	eds := dp.getCandidateEndPort(ctx)

	//计算每一个的得分并且重赋值
	for _, ed := range eds {
		ed.CalculateScore(ctx)
	}

	//类Lambda 根据得分降序排序
	sort.Slice(eds, func(i, j int) bool {
		return eds[i].Score > eds[j].Score
	})

	return eds
}

// 获取候选EndPort的时候是遍历map 从这个方法获取所有候选列表 再基于他们的状态信息进行计算
func (dp *Dispatcher) getCandidateEndPort(ctx *domain.IpConfContext) []*EndPort {
	dp.RLock()
	defer dp.RUnlock()
	candidateList := CloneEndPort(dp.candidateTable)
	for _, ed := range dp.candidateTable {
		candidateList = append(candidateList, ed)
	}
	return candidateList
}

func (dp *Dispatcher) delNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	delete(dp.candidateTable, event.Key())
}

// addNode 增加节点
func (dp *Dispatcher) addNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	var (
		ed *EndPort
		ok bool
	)

	//塞进去的时候 先判断这个节点是否存在 不存在创建 存在更新状态
	if ed, ok = dp.candidateTable[event.Key()]; !ok {
		ed = NewEndPort(event.IP, event.Port)
		dp.candidateTable[event.Key()] = ed
	}

	ed.UpdateStat(&Stat{
		ConnectNum:   event.ConnectNum,
		MessageBytes: event.MessageBytes,
	})

}
