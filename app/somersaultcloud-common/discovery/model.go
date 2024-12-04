package discovery

import jsoniter "github.com/json-iterator/go"

// EndpointInfo 服务包含的基本信息 Ip & port为基础信息 map中为用户根据需求自行拓展的
//
//	初始的参数有连接树 conn_num 和 message_bytes 消息大小
type EndpointInfo struct {
	IP       string                 `json:"ip"`
	Port     string                 `json:"port"`
	MetaData map[string]interface{} `json:"meta"`
}

func UnMarshal(data []byte) (*EndpointInfo, error) {
	ed := &EndpointInfo{}
	err := jsoniter.Unmarshal(data, ed)
	if err != nil {
		return nil, err
	}
	return ed, nil
}
func (edi *EndpointInfo) Marshal() string {
	data, err := jsoniter.Marshal(edi)
	if err != nil {
		panic(err)
	}
	return string(data)
}
