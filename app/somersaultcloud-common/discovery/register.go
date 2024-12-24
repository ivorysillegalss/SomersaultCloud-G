package discovery

import (
	"SomersaultCloud/app/somersaultcloud-common/log"
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ServiceRegister struct {
	cli           *clientv3.Client //etcd client
	leaseID       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	val           string
	ctx           context.Context
}

// NewServiceRegister 服务注册 设置租约
func NewServiceRegister(cli *clientv3.Client, key string, ctx context.Context, info *EndpointInfo, lease int64) (*ServiceRegister, error) {
	ser := &ServiceRegister{
		cli: cli,
		key: key,
		val: info.Marshal(),
		ctx: ctx,
	}
	err := ser.putKeyWithLease(lease)
	if err != nil {
		return nil, nil
	}

	return ser, nil
}

func (s *ServiceRegister) putKeyWithLease(lease int64) error {
	//新建租约 设置其TTL
	resp, err := s.cli.Grant(s.ctx, lease)
	if err != nil {
		log.GetTextLogger().Error("grant lease error")
		return err
	}

	//注册服务 绑定租约
	_, err = s.cli.Put(s.ctx, s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		log.GetTextLogger().Error("put service or bind lease error with {}", resp.ID)
		return err
	}

	//指定对当前所交流的键值对注册自动续约
	keepAliveChan, err := s.cli.KeepAlive(s.ctx, resp.ID)
	if err != nil {
		log.GetTextLogger().Error("register keep-alive error with {}", resp.ID)
	}

	s.leaseID = resp.ID
	s.keepAliveChan = keepAliveChan
	return nil
}

// ListenLeaseRespChan 为每一个服务都创建一个线程 监听他的续租情况
//
//	通过阻塞 channel 实现
func (s *ServiceRegister) ListenLeaseRespChan() {
	for leaseResp := range s.keepAliveChan {
		log.GetTextLogger().Info("lease success leaseID:%d, Put key:%s,val:%s reps:+%v",
			s.leaseID, s.key, s.val, leaseResp)
	}
	log.GetTextLogger().Info("lease failed !!!  leaseID:%d, Put key:%s,val:%s", s.leaseID, s.key, s.val)
}

// UpdateValue 更新服务发现值
func (s *ServiceRegister) UpdateValue(val *EndpointInfo) error {
	jsonVal := val.Marshal()
	_, err := s.cli.Put(s.ctx, s.key, jsonVal, clientv3.WithLease(s.leaseID))
	if err != nil {
		log.GetTextLogger().Error("update value error leaseID=%d", s.leaseID)
		return err
	}
	s.val = jsonVal
	log.GetTextLogger().Info("ServiceRegister.updateValue leaseID=%d Put key=%s,val=%s, success!", s.leaseID, s.key, s.val)
	return nil
}

// Close 注销服务
func (s *ServiceRegister) Close() error {
	//撤销租约
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	log.GetTextLogger().Info("lease close !!!  leaseID:%d, Put key:%s,val:%s  success!", s.leaseID, s.key, s.val)
	return s.cli.Close()
}
