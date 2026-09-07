#!/usr/bin/env bash
# Go SQLi Lab - macOS 启动脚本（Finder 双击即运行）
# 首次运行会自动编译并初始化数据库，然后启动服务并打开浏览器。
# 设置环境变量 SQLI_LAB_NO_BROWSER=1 可禁止自动打开浏览器。
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

echo "=============================================="
echo "   Go SQLi Lab 快速启动"
echo "=============================================="

# ---- 检测 Go ----
if ! command -v go >/dev/null 2>&1; then
    echo "[错误] 未找到 Go 环境。" >&2
    echo "请先安装 Go 1.21 或更高版本，并确保 go 命令已加入 PATH。" >&2
    echo "下载地址: https://go.dev/dl/" >&2
    exit 1
fi

# ---- 构建 ----
mkdir -p bin
echo "[1/3] 编译中..."
go build -o bin/go-sqli-lab src/main.go

# ---- 初始化数据库（幂等）----
echo "[2/3] 初始化数据库..."
./bin/go-sqli-lab --setup-db

# ---- 读取服务端口（默认 8080）----
PORT="$(awk '/^server:/{s=1} s&&/port:/{print $2; exit}' config/app.yaml)"
PORT="${PORT:-8080}"

echo "[3/3] 启动服务: http://localhost:$PORT"
echo "按 Ctrl+C 停止服务。"

if [ -z "$SQLI_LAB_NO_BROWSER" ]; then
    ( sleep 2; open "http://localhost:$PORT" ) &
fi

exec ./bin/go-sqli-lab
