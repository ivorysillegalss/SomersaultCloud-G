# 筋斗云 
Go module
## 项目介绍

项目针对具体的应用场景，结合市面权威LLM，经过Prompt预处理后应用于具体行业和活动中。针对不同功能设计了标签系统来量化任务。用户通过选择标签将需求具象化，通过标签设定的提示词而得到回复，提高工作效率

**内测网址：https://www.jdygpt.com/**

## 项目亮点：

- 为加强业务的可拓展性，提高开发效率。项目基于**整洁架构**结合**DDD**构建，通过**wire**进行依赖注入。项目运行链路由责任链模式组装，热插拔式定制需求。并提供集成三方LLM所需**SPI**，将相关逻辑以工厂模式封装为执行器。

- 为加快响应速度，提高吞吐量，组装请求完毕即返回用户ACK。使用**RabbitMQ**异步生产任务至线程池中。消费层后置异步处理generation，通过SSE以流式输出主动send至客户端中。使生产消费模块解耦，接口平均延迟降低60%。

- 为单批流信息及其状态设计完善的**生命周期**。以**时间戳**实现**单批次**为细粒度的**乐观锁（CAS）**，解决不同批次间的竞态问题。处理单批信息时，使用**时间轮**算法高效管理回收复用管道资源。

- 基于**写多读少**，消息价值随时间递减的业务场景，DB层设计以**会话为单位**存储历史记录，显著降低用户信息读取延迟，提供**Protobuf & Gzip+JSON**两种数据前置写入策略。且**写DB操作配合MQ异步消费**，用户感知较小，解耦读写操作。同时通过**redis-zset + 双哈希表 + lua**，原子性实现哈希类型**单Field过期与LRU存储（提取LRU中台代码）**。从而实现**冷热数据分区**，维护近期热数据。最大限度保证请求打在cache上，非DB上。

- 实现部署上线，使用Nginx反向代理域名。配置SSL实现https协议，确保数据的安全传输。接入gRpc保证通讯速率，提高模块间调用效率。


-   TBD
    -   使用**Flink**进行异构数据库间**流式刷新删除**，将mysql中超7天未刷新的chat记录迁移至Hbase中，将超30天记录从Hbase中删除。

## 项目技术栈：

- RabbitMQ 

- jwt 

- Redis 

- Gin 

- Docker 

- SSE 

- Mysql 

- Wire 

- logrus

- ProtoBUf

TBD

- Grpc

- Flink

**项目功能架构图**

![image](https://github.com/user-attachments/assets/60c61a3b-80e3-4224-8c71-e64380855b75)


**项目结构**

```
├── api
│   ├── controller
│   ├── dto
│   ├── middleware
│   │   └── taskchain
│   └── route
├── bootstrap
├── cmd
├── constant
│   ├── cache
│   ├── common
│   ├── dao
│   ├── mq
│   ├── request
│   ├── sys
│   └── task
├── consume
├── cron
├── domain
├── executor
├── handler
│   └── stream
├── infrastructure
│   ├── log
│   ├── lru
│   │   └── lua
│   ├── mongo
│   ├── mysql
│   ├── pool
│   ├── rabbitmq
│   └── redis
├── internal
│   ├── checkutil
│   ├── compressutil
│   ├── ioutil
│   ├── kvutil
│   ├── requtil
│   └── tokenutil
├── proto
├── repository
│   └── lua
├── sql
├── task
└── usecase
    └── lua
```



