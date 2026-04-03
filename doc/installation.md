# 安装指南

## 系统要求

- **操作系统**: Windows / Linux / macOS
- **Go 版本**: 1.21 或更高
- **内存**: 至少 512MB
- **磁盘**: 至少 100MB 可用空间

## 快速安装

### 1. 从源码安装

```bash
# 克隆仓库
git clone https://github.com/yourusername/go-sqli-lab.git
cd go-sqli-lab

# 安装依赖
go mod tidy

# 构建
make build
# 或直接: go build -o bin/go-sqli-lab src/main.go
```

### 2. 初始化数据库

```bash
# 自动创建数据库和表
make setup-db
# 或直接: go run src/main.go --setup-db
```

### 3. 启动服务

```bash
# 启动服务器
make run
# 或直接: go run src/main.go
```

访问 http://localhost:8080

## 数据库配置

### SQLite (默认)

无需额外配置，数据存储在 `./data/sqli_lab.db`

### MySQL

```yaml
# config/app.yaml
database:
  type: mysql
  host: localhost
  port: 3306
  user: root
  password: your_password
  dbname: security
```

或设置环境变量：
```bash
export SQLI_LAB_DB_TYPE=mysql
export SQLI_LAB_DB_HOST=localhost
export SQLI_LAB_DB_USER=root
export SQLI_LAB_DB_PASSWORD=your_password
export SQLI_LAB_DB_NAME=security
```

### PostgreSQL

```yaml
# config/app.yaml
database:
  type: postgres
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  dbname: security
  sslmode: disable
```

## Docker 部署

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o go-sqli-lab src/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/go-sqli-lab .
COPY --from=builder /app/config ./config
COPY --from=builder /app/static ./static
EXPOSE 8080
CMD ["./go-sqli-lab"]
```

构建和运行：
```bash
docker build -t go-sqli-lab .
docker run -p 8080:8080 go-sqli-lab
```

## 常见问题

### 1. 端口被占用

修改配置或使用环境变量：
```bash
export SQLI_LAB_PORT=8081
```

### 2. 数据库连接失败

- 检查数据库服务是否运行
- 验证连接配置是否正确
- 查看日志文件 `logs/app.log`

### 3. 权限问题

确保程序有权限创建目录和文件：
```bash
chmod -R 755 .
```

## 卸载

直接删除项目目录即可：
```bash
rm -rf go-sqli-lab
```

如有数据文件，请手动删除：
```bash
rm -rf data/
rm -rf logs/
```
