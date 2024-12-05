package dispatcher

import (
	"sync/atomic"
	"unsafe"
)

type EndPort struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
	//根据这个分数进行排序 并返回给客户端可用性最高的选项
	Score  float64      `json:"-"`
	Stats  *Stat        `json:"-"`
	Window *stateWindow `json:"-"`
}

func NewEndPort(ip, port string) *EndPort {
	ed := &EndPort{
		IP:     ip,
		Port:   port,
		Window: newStateWindow(),
	}
	ed.Stats = ed.Window.getStat()
	go func() {
		for stat := range ed.Window.statChan {
			ed.Window.appendStat(stat)
			newStat := ed.Window.getStat()
			atomic.SwapPointer((*unsafe.Pointer)((unsafe.Pointer)(ed.Stats)), unsafe.Pointer(newStat))
		}
	}()
	return ed
}

func CloneEndPort(eps map[string]*EndPort) []*EndPort {
	endports := make([]*EndPort, 0, len(eps))
	for _, ep := range eps {
		endports = append(endports, ep)
	}
	return endports
}

func (ed *EndPort) UpdateStat(s *Stat) {
	ed.Window.statChan <- s
}
