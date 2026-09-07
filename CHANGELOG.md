# Changelog

本项目所有重要变更均记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [1.5.0] - 2026-09-07

### Added
- 新增 Less-68 关卡:Resume Export XLS Injection(简历导出注入 - 按"应聘人员求职登记表"导出 XLS 文件,查询关键字可 SQL 注入)
- 新增 Less-69 关卡:Resume Export XLSX Injection(简历导出注入 - 按"应聘人员求职登记表"导出 XLSX 文件,查询关键字可 SQL 注入)
- 新增 `resumes` 表迁移/种子/重置逻辑,兼容 SQLite/MySQL/PostgreSQL
- 首页"文件导出注入"分区扩展至 Less-66 ~ Less-69,新增简历登记表卡片预览样式
- 新增简历模板生成工具 `tools/genresume`(基于 excelize 生成 XLSX 模板文件)
- 引入 `github.com/xuri/excelize/v2` 依赖

## [1.4.0] - 2026-09-07

### Added
- 新增 Less-66 关卡:Export XLS Injection(文件导出注入 - 导出 XLS 文件,查询关键字可 SQL 注入)
- 新增 Less-67 关卡:Export XLSX Injection(文件导出注入 - 导出 XLSX 文件,查询关键字可 SQL 注入)
- 首页新增"文件导出注入"分区入口,并添加对应的分类样式
- 新增跨平台启动脚本:`scripts/start-linux.sh`、`scripts/start-macos.command`、`scripts/start-windows.bat`

## [1.3.0] - 2026-04-04

### Added
- 为 CI 集成测试添加测试数据库

## [1.2.0] - 2026-04-04

### Added
- 添加多数据库集成测试与 CI 流水线

### Fixed
- 修复 PostgreSQL 下 LIMIT 语法不兼容问题

## [1.1.0] - 2026-04-04

### Added
- 新增数据库重置功能(含重置 UI 页面与 API 接口)

## [1.0.0] - 2026-04-03

### Added
- 首个正式发布版本,将 sqli-labs 前 65 关移植到 Go 实现
