# 筋斗云 
Go module
## 项目介绍

项目针对具体的应用场景，结合市面权威LLM，经过Prompt预处理后应用于具体行业和活动中。针对不同功能设计了标签系统来量化任务。用户通过选择标签将需求具象化，通过标签设定的提示词而得到回复，提高工作效率

**内测网址：http://www.jdygpt.com/**

## 项目亮点：

-   基于**整洁架构**结合**DDD**构建项目，通过**wire**进行依赖注入，二次封装基础设施，轻松拓展新业务。

-   通过**责任链模式**实现普遍req任务组装，通过**统一接口**实现，可排列组合满足不同请求处理需求。**工厂模式**第三方LLM封装为执行器，开发者可通过实现该接口轻松集成不同LLM。

-   请求组装完毕，返回用户**ACK**。使用**RabbitMQ**异步发送任务至线程池中处理。

-   基于**写多读少**，历史记录价值普遍低的业务场景，DB层设计**读放大**存储模式，DB以**用户维度**存储历史记录，将消息**序列化Gzip压缩**后写入。显著降低读取延迟。写DB操作异步前提下，用户感知较小。

-   结合**redis-zset**实现lru存储，**提取lru中台代码**。通过**双哈希表+lua脚本**，哈希类型单field过期功能实现。用于将近5条活跃chat，单chat内近10条消息记录写入cache中。

-   **Nginx**处理域名挂靠，反向代理与路由转发请求，使用**gRpc**进行服务间高效远程调用。

-   **RabbitMQ**重连封装。**docker-compose**动态部署。日志框架**logrus**接入，封装JSON,Text二格式日志。Jwt用户鉴权。



-   TBD
    -   使用**Flink**进行异构数据库间**流式刷新删除**，将mysql中超7天未刷新的chat记录迁移至Hbase中，将超30天记录从Hbase中删除。

    -   消费层将generation处理后，通过**SSE**以流式输出主动send至客户端中。




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

TBD

- Grpc

- Flink



**项目结构**

├─api
│  ├─controller
│  ├─dto
│  ├─middleware
│  │  └─taskchain
│  └─route
├─bootstrap
├─cmd
├─constant
│  ├─cache
│  ├─common
│  ├─dao
│  ├─mq
│  ├─request
│  ├─sys
│  └─task
├─consume
├─cron
│  └─lua
├─domain
├─executor
├─handler
├─infrastructure
│  ├─log
│  ├─lru
│  │  └─lua
│  ├─mongo
│  ├─mysql
│  ├─pool
│  ├─rabbitmq
│  └─redis
├─internal
│  ├─checkutil
│  ├─compressutil
│  ├─ioutil
│  ├─requtil
│  └─tokenutil
├─repository
├─task
└─usecase
    └─lua

