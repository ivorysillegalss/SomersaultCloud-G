package source

import (
	"SomersaultCloud/app/somersaultcloud-common/discovery"
	"fmt"
)

var eventChan chan *Event

func EventChan() <-chan *Event {
	return eventChan
}

type EventType string

const (
	AddNodeEvent EventType = "addNode"
	DelNodeEvent EventType = "delNode"
)

type Event struct {
	Type            EventType
	IP              string
	Port            string
	ConnectNum      float64
	MessageBytes    float64
	AvailableMem    int64
	CpuIdleTime     float64
	RequestCount    float64
	RequestDuration float64
}

// NewEvent 将传过来的服务端点包装成为一个新的Event节点对象
func NewEvent(info *discovery.EndpointInfo) *Event {
	if info == nil || info.MetaData == nil {
		return nil
	}
	var connNum, msgBytes float64
	if data, ok := info.MetaData["connect_num"]; ok {
		connNum = data.(float64)
	}
	if data, ok := info.MetaData["message_bytes"]; ok {
		msgBytes = data.(float64)
	}
	return &Event{
		Type:         AddNodeEvent,
		IP:           info.IP,
		Port:         info.Port,
		ConnectNum:   connNum,
		MessageBytes: msgBytes,
	}
}

func (e *Event) Key() string {
	return fmt.Sprintf("%s:%s", e.IP, e.Port)
}
