package discovery

import (
	log2 "SomersaultCloud/app/somersaultcloud-common/log"
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
	"time"
)

// ServiceDiscovery 服务发现
type ServiceDiscovery interface {
	WatchService(prefix string, set, del func(key, value string)) error
	Close() error
}

type etcdServiceDiscovery struct {
	cli  *clientv3.Client
	lock sync.Mutex
	ctx  context.Context
}

func (e *etcdServiceDiscovery) WatchService(prefix string, set, del func(key, value string)) error {
	resp, err := e.cli.Get(e.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	//初始化服务列表
	for _, ev := range resp.Kvs {
		set(string(ev.Key), string(ev.Value))
	}

	//根据前缀进行监视 假如有server需要变更 马上修改
	e.watcher(prefix, resp.Header.Revision+1, set, del)
	return nil
}

func (e *etcdServiceDiscovery) Close() error {
	return e.cli.Close()
}

func (e *etcdServiceDiscovery) watcher(prefix string, rev int64, set, del func(key, value string)) {
	rch := e.cli.Watch(e.ctx, prefix, clientv3.WithPrefix(), clientv3.WithRev(rev))
	log2.GetTextLogger().Info("watching prefix: %s now", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			//修改或者新增
			case mvccpb.PUT:
				set(string(ev.Kv.Key), string(ev.Kv.Value))
				//删除
			case mvccpb.DELETE:
				del(string(ev.Kv.Key), string(ev.Kv.Value))
			}
		}
	}
}

func NewServiceDiscovery(ctx context.Context, endpoints []string, timeout time.Duration) ServiceDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	})
	defer cli.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	return &etcdServiceDiscovery{
		cli: cli,
		ctx: ctx,
	}
}
