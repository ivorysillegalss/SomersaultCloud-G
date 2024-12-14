# 筋斗云 
Go module
## 项目介绍

项目针对具体的应用场景，结合市面权威LLM，经过Prompt预处理后应用于具体行业和活动中。针对不同功能设计了标签系统来量化任务。用户通过选择标签将需求具象化，通过标签设定的提示词而得到回复，提高工作效率

**内测网址：https://www.jdygpt.com/**

## 项目亮点：

- 为加强业务的可拓展性，提高开发效率。项目基于**整洁架构**结合**DDD**构建，通过**wire**进行依赖注入。项目运行链路由**责任链模式**组装，热插拔式定制需求。并提供集成三方LLM所需**SPI**，将相关逻辑以工厂模式封装为执行器。

- 为加快响应速度，提高吞吐量，组装请求完毕即返回用户ACK。使用**RabbitMQ**异步生产任务至线程池中。消费层后置异步消费，成功则SSE流式输出推送至客户端，失败则转发至**死信队列**重试。解耦生产消费模块，接口平均延迟降低40%。
  
- 为实现多服务间**负载均衡**，基于**gRPC长连接**配合**etcd的watch机制**，自研exporter与ipconfig模块，实现连接轻量化、延迟低、一致性强的**服务发现**。并同时将状态信息推送至**Prometheus&Grafana**中，实时**监控业务运行状态**。

- 为单批流信息及其状态设计完善的**生命周期**。以**时间戳**实现**单批次**为细粒度的**乐观锁（CAS）**，解决不同批次间的竞态问题。处理单批信息时，使用**时间轮**算法高效管理回收复用管道资源。

- 结合**Redis**与**桶限流**，实现**分布式环境**下单用户接口访问次数**限流**。注册**Hertz**中间件调用。

- 基于**写多读少**，消息价值随时间递减的业务场景，DB层设计以**会话为单位**存储历史记录，显著降低用户信息读取延迟，提供**Protobuf & Gzip+JSON**两种数据前置写入策略。且**写DB操作配合MQ异步消费**，用户感知较小，解耦读写操作。同时通过**redis-zset + 双哈希表 + lua**，原子性实现哈希类型**单Field过期与LRU存储（提取LRU中台代码）**。从而实现**冷热数据分区**，维护近期热数据。最大限度保证请求打在cache上，非DB上。

- 部署上线，使用**Nginx**反向代理域名。配置SSL实现https协议，确保数据的安全传输。


-   TBD
    -   使用**Flink**进行异构数据库间**流式刷新删除**，将mysql中超7天未刷新的chat记录迁移至Hbase中，将超30天记录从Hbase中删除。

## 项目技术栈：

| 技术栈                 | 作用说明                                                     |
| ---------------------- | ------------------------------------------------------------ |
| **RabbitMQ**           | 解耦**存储，api调用**中任务的生产与消费。异步化，提高吞吐    |
| **Redis**              | 提供各类数据的缓存存储，LRU淘汰策略，分布式环境下接口**限流** |
| **MySQL**              | 存储DB层会话和历史记录等数据。会话数据兼容json或pb           |
| **GORM**               | 主流ORM框架，提高CRUD开发效率             |
| **Gin&Hertz**          | 主流高效Web框架，提供**API接口**处理请求，**中间件**前置处理 |
| **gRPC**               | 高效通讯协议，低延迟传送**模块状态信息**                     |
| **ProtoBuf**           | 高效序列化协议，结合gRPC使用。**存储**序列化兼容             |
| **etcd**               | 强一致性kv存储，监控模块状态信息。**服务发现数据源**         |
| **Docker**             | 容器化部署，并提供docker-compose简化多服务编排               |
| **Prometheus&Grafana** | 收集监控信息，提供可视化展示，监控**项目实时运行状态**       |
| **Nginx**              | 反向代理&负载均衡，HTTPS保证数据安全传输                     |
| **JWT**                | 认证 & 鉴权                                                  |
| **SSE**                | 轻量化流式数据传输。（轻于WebSocket）                        |
| **Wire**               | 生成式依赖注入组件                                           |
| **Logrus**             | 日志记录组件                                                 |
| Flink                  | 异构数据库间的数据流式同步与迁移，确保系统中不同数据库之间的数据一致性和高效的数据更新。 |

## 项目图相关

**项目架构图**

![architecture](https://github.com/user-attachments/assets/249b6766-fa69-4424-9005-2b0616973253)

**项目各模块流程图**

![image](https://github.com/user-attachments/assets/41867220-fe93-4882-a4da-cd8a5f3ffcea)

**监控模块流程图**

![image](https://github.com/user-attachments/assets/09e87d93-ea5d-4afa-97a4-5b461a683e16)

**业务模块功能流程图**

![image](https://github.com/user-attachments/assets/60c61a3b-80e3-4224-8c71-e64380855b75)


## 项目结构

```

├── Dockerfile
├── README.md
├── app
│   ├── somersaultcloud-chat
│   ├── somersaultcloud-common
│   ├── somersaultcloud-exporter
│   └── somersaultcloud-ipconfig
├── docker-compose.yaml
├── go.mod
├── go.sum
├── prometheus.yml
├── somersaultcloud.yaml
└── sql

```



