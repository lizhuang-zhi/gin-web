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

## 场景流程任务7 - 待定
### 任务
1. 新增网关服务区分访问
2. 新增服务

## 场景流程任务8 - 待定
### 任务
1. 等搭建几个服务后，书写部分openAPI接口，再通过make命令执行shell脚本，一键生成本地测试数据
2. 但是需要考虑后续接口改动的成本