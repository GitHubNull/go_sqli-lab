# Go SQLi Lab Makefile

# 变量定义
BINARY_NAME=go-sqli-lab
BUILD_DIR=bin
SRC_DIR=src
MAIN_FILE=$(SRC_DIR)/main.go
GO=go

# 颜色定义
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: all build clean run test setup-db deps help

# 默认目标
all: build

# 构建项目
build:
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# 清理构建文件
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -rf logs/
	@rm -rf tmp/
	@echo "$(GREEN)Clean complete$(NC)"

# 运行项目
run:
	@echo "$(GREEN)Starting $(BINARY_NAME)...$(NC)"
	$(GO) run $(MAIN_FILE)

# 运行项目（带数据库初始化）
run-setup:
	@echo "$(GREEN)Starting $(BINARY_NAME) with setup...$(NC)"
	$(GO) run $(MAIN_FILE) --setup-db

# 初始化数据库
setup-db:
	@echo "$(YELLOW)Setting up database...$(NC)"
	$(GO) run $(MAIN_FILE) --setup-db
	@echo "$(GREEN)Database setup complete$(NC)"

# 安装依赖
deps:
	@echo "$(YELLOW)Installing dependencies...$(NC)"
	$(GO) mod tidy
	@echo "$(GREEN)Dependencies installed$(NC)"

# 运行测试
test:
	@echo "$(YELLOW)Running tests...$(NC)"
	$(GO) test -v ./...

# 运行测试（带覆盖率）
test-coverage:
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	$(GO) test -cover ./...

# 代码格式化
fmt:
	@echo "$(YELLOW)Formatting code...$(NC)"
	$(GO) fmt ./...
	@echo "$(GREEN)Formatting complete$(NC)"

# 代码检查
vet:
	@echo "$(YELLOW)Running go vet...$(NC)"
	$(GO) vet ./...

# 检查日志
logs:
	@if [ -f logs/app.log ]; then \
		tail -f logs/app.log; \
	else \
		echo "$(RED)Log file not found. Start the server first.$(NC)"; \
	fi

# 查看帮助信息
help:
	@echo "$(GREEN)Go SQLi Lab - Available Commands:$(NC)"
	@echo ""
	@echo "  $(YELLOW)make build$(NC)        - Build the binary"
	@echo "  $(YELLOW)make run$(NC)          - Run the application"
	@echo "  $(YELLOW)make run-setup$(NC)    - Run with database setup"
	@echo "  $(YELLOW)make setup-db$(NC)     - Initialize database only"
	@echo "  $(YELLOW)make clean$(NC)        - Clean build files"
	@echo "  $(YELLOW)make deps$(NC)         - Install dependencies"
	@echo "  $(YELLOW)make test$(NC)         - Run tests"
	@echo "  $(YELLOW)make test-coverage$(NC)- Run tests with coverage"
	@echo "  $(YELLOW)make fmt$(NC)          - Format code"
	@echo "  $(YELLOW)make vet$(NC)          - Run go vet"
	@echo "  $(YELLOW)make logs$(NC)         - View logs"
	@echo "  $(YELLOW)make help$(NC)         - Show this help"
	@echo ""

# Docker 相关命令（可选）
docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -p 8080:8080 $(BINARY_NAME)
