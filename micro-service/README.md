# 微服务

## 目标
1. 涉及网关服务学习
2. 涉及使用grpc学习
3. 涉及使用http学习

## 场景流程任务1 - 活动服务搭建
### 任务
1. 搭建活动服务,提供http接口,用于新增活动数据到内存和MongoDB数据库中
2. 活动服务提供grpc调用,用于客户端获取活动数据
### tips
1. http服务和grpc服务需要开启两个协程启动,如果grpc服务写在http服务后,则会阻塞主线程,导致http服务无法启动

## 场景流程任务2 - 增加组件Options
### 任务
1. 新增Options,包含logger组件
2. 新增MongoDB组件

## 场景流程任务3 - 新增grpc协议，单独封装到manager包下
### 任务
1. 书写更多grpc服务，然后封装管理

## 场景流程任务4 - 新增yaml配置文件与Makefile
### 任务
1. 新增配置文件

## 场景流程任务5 - 整合微服务项目结构
### 任务
1. 提取各个服务到./cluster目录下
- cluster    服务集群：多服务集成
    - activity
    - mail
- core       公共抽象层：公共抽象接口、结构体
- client     客户端
- protobuf   protobuf协议文件（make proto生成）
    - gen-pb
    - proto
- configs    yaml配置文件
- scripts    执行脚本

## 场景流程任务6 - yaml配置新增日志配置
### 任务
1. yaml配置新增日志配置
2. 收集日志到本地可视化工具，并展示[Elasticsearch、Kibana]（后续再尝试）

## 场景流程任务7 - 重新整理微服务项目结构
### 任务
micro-service/
├── client/                   # 客户端，如Go语言的客户端, 用于测试grpc方法
│   ├── rpc
│   │   ├── pbclient.go       # new client
│   │   └── rpcmethod.go      # 调用rpc方法
│   └── main.go               # client启动入口
├── cluster/
│   ├── activity/
│   │   ├── api               # Activity服务的HTTP接口
│   │   ├── manager           # Activity服务的GRPC接口服务实现
│   │   ├── model             # Activity服务的HTTP接口数据持久化
│   │   ├── router            # Activity服务的HTTP接口路由层
│   │   └── activity.go       # Activity服务启动入口
│   ├── mail/
│   │   └── mail.go           # Mail服务启动入口
│   ├── common/
│   │   └── config            # 配置管理，如读取和解析配置文件
│   │   └── logger            # 日志处理
│   │   └── mongodb           # 数据库处理
│   │   └── utils             # 实用工具函数
├── cmd/
│   └── server/
│       └── main.go           # 统一启动所有服务的入口文件
├── configs/
│   └── local           
│       └── config.yaml       # 配置文件，如日志级别、MongoDB连接信息等
├── protobuf/                 
│   └── gen-pb                # 生成pb文件
│   └── proto                 # 定义的协议文件
│       └── notice.proto      # 公告协议文件
├── scripts/                  # 启动脚本和部署脚本
├── logs                      # 日志文件
├── Makefile                  # make命令脚本
└── README.md                 # 项目文档，包括如何运行、部署和贡献代码的说明

## 场景流程任务8 - 广播与客户端代码整理
### 任务
1. grpc广播
2. 整理客户端分类测试

## 场景流程任务9 - 使用cli、metadata初始化
### 任务 
1. cli、metadata

## 场景流程任务10 - 待定
### 任务
1. 新增网关服务区分访问(实现网关映射到不同服务的grpc)
2. 整合网关内容(之前的服务发现、ip限流等)
3. 实现新增端点(endpoints)

## 场景流程任务11 - 待定
### 任务
1. 等搭建几个服务后，书写部分openAPI接口，再通过make命令执行shell脚本，一键生成本地测试数据(待定,先不做)