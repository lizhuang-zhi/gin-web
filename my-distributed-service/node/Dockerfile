# node/Dockerfile

# 使用Go的官方镜像作为构建环境
FROM golang:1.18-alpine as builder

# 设置工作目录
WORKDIR /app

# 将项目根目录下的go.mod和go.sum复制到容器中
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 将整个项目复制到容器中
COPY . .

# 构建服务节点应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -o /service-node ./node/main.go

# 使用scratch作为最小运行环境
FROM scratch

# 从builder阶段复制编译完成的应用程序
COPY --from=builder /service-node /service-node

# 运行服务节点应用程序
CMD ["/service-node"]