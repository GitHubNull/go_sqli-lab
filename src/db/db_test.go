package db

import (
	"fmt"
	"os"
	"testing"

	"go-sqli-lab/src/config"
	"go-sqli-lab/src/logger"
)

// getTestConfig 返回测试配置
func getTestConfig(dbType string) config.DatabaseConfig {
	switch dbType {
	case "sqlite":
		return config.DatabaseConfig{
			Type: "sqlite",
			DSN:  "./test_data/test_sqli_lab.db",
		}
	case "mysql":
		return config.DatabaseConfig{
			Type:     "mysql",
			Host:     getEnv("TEST_MYSQL_HOST", "localhost"),
			Port:     getEnvInt("TEST_MYSQL_PORT", 3306),
			User:     getEnv("TEST_MYSQL_USER", "root"),
			Password: getEnv("TEST_MYSQL_PASSWORD", ""),
			DBName:   getEnv("TEST_MYSQL_DB", "test_sqli_lab"),
		}
	case "postgres":
		return config.DatabaseConfig{
			Type:     "postgres",
			Host:     getEnv("TEST_POSTGRES_HOST", "localhost"),
			Port:     getEnvInt("TEST_POSTGRES_PORT", 5432),
			User:     getEnv("TEST_POSTGRES_USER", "postgres"),
			Password: getEnv("TEST_POSTGRES_PASSWORD", ""),
			DBName:   getEnv("TEST_POSTGRES_DB", "test_sqli_lab"),
			SSLMode:  "disable",
		}
	default:
		return config.DatabaseConfig{Type: "sqlite", DSN: "./test_data/test_sqli_lab.db"}
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		var result int
		fmt.Sscanf(v, "%d", &result)
		return result
	}
	return defaultValue
}

// TestDatabase_Connect 测试数据库连接
func TestDatabase_Connect(t *testing.T) {
	dbTypes := []string{"sqlite"}

	// 如果环境变量设置了 MySQL/PostgreSQL 测试配置，也测试它们
	if os.Getenv("TEST_MYSQL_HOST") != "" {
		dbTypes = append(dbTypes, "mysql")
	}
	if os.Getenv("TEST_POSTGRES_HOST") != "" {
		dbTypes = append(dbTypes, "postgres")
	}

	for _, dbType := range dbTypes {
		t.Run(fmt.Sprintf("Connect_%s", dbType), func(t *testing.T) {
			cfg := getTestConfig(dbType)
			log, _ := logger.New(config.LogConfig{Level: "error", Output: "stdout"})
			defer log.Close()

			db, err := New(cfg, log)
			if err != nil {
				t.Fatalf("%s: 连接数据库失败: %v", dbType, err)
			}
			defer db.Close()

			// 测试 Ping
			if err := db.GetDB().Ping(); err != nil {
				t.Fatalf("%s: Ping 失败: %v", dbType, err)
			}
		})
	}
}

// TestDatabase_Migrate 测试数据库迁移
func TestDatabase_Migrate(t *testing.T) {
	dbTypes := []string{"sqlite"}

	if os.Getenv("TEST_MYSQL_HOST") != "" {
		dbTypes = append(dbTypes, "mysql")
	}
	if os.Getenv("TEST_POSTGRES_HOST") != "" {
		dbTypes = append(dbTypes, "postgres")
	}

	for _, dbType := range dbTypes {
		t.Run(fmt.Sprintf("Migrate_%s", dbType), func(t *testing.T) {
			cfg := getTestConfig(dbType)
			log, _ := logger.New(config.LogConfig{Level: "error", Output: "stdout"})
			defer log.Close()

			db, err := New(cfg, log)
			if err != nil {
				t.Fatalf("%s: 连接数据库失败: %v", dbType, err)
			}
			defer db.Close()

			// 执行迁移
			if err := db.Migrate(); err != nil {
				t.Fatalf("%s: 迁移失败: %v", dbType, err)
			}

			// 验证表是否创建成功
			tables := []string{"users", "emails", "uagents", "referers"}
			for i := 1; i <= 12; i++ {
				tables = append(tables, fmt.Sprintf("challenge%d", i))
			}

			for _, table := range tables {
				var count int
				var query string
				if dbType == "postgres" {
					query = fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = '%s'", table)
				} else if dbType == "mysql" {
					query = fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = '%s'", table)
				} else {
					// SQLite
					query = fmt.Sprintf("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='%s'", table)
				}

				if err := db.QueryRow(query).Scan(&count); err != nil {
					t.Fatalf("%s: 检查表 %s 失败: %v", dbType, table, err)
				}
				if count == 0 {
					t.Errorf("%s: 表 %s 未创建", dbType, table)
				}
			}
		})
	}
}

// TestDatabase_Seed 测试种子数据插入
func TestDatabase_Seed(t *testing.T) {
	dbTypes := []string{"sqlite"}

	if os.Getenv("TEST_MYSQL_HOST") != "" {
		dbTypes = append(dbTypes, "mysql")
	}
	if os.Getenv("TEST_POSTGRES_HOST") != "" {
		dbTypes = append(dbTypes, "postgres")
	}

	for _, dbType := range dbTypes {
		t.Run(fmt.Sprintf("Seed_%s", dbType), func(t *testing.T) {
			cfg := getTestConfig(dbType)
			log, _ := logger.New(config.LogConfig{Level: "error", Output: "stdout"})
			defer log.Close()

			db, err := New(cfg, log)
			if err != nil {
				t.Fatalf("%s: 连接数据库失败: %v", dbType, err)
			}
			defer db.Close()

			// 先执行迁移
			if err := db.Migrate(); err != nil {
				t.Fatalf("%s: 迁移失败: %v", dbType, err)
			}

			// 执行种子数据插入
			if err := db.Seed(); err != nil {
				t.Fatalf("%s: 种子数据插入失败: %v", dbType, err)
			}

			// 验证 users 表数据
			var count int
			if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
				t.Fatalf("%s: 查询 users 表失败: %v", dbType, err)
			}
			if count != 8 {
				t.Errorf("%s: users 表应有8条记录，实际有 %d 条", dbType, count)
			}

			// 验证 emails 表数据
			if err := db.QueryRow("SELECT COUNT(*) FROM emails").Scan(&count); err != nil {
				t.Fatalf("%s: 查询 emails 表失败: %v", dbType, err)
			}
			if count != 8 {
				t.Errorf("%s: emails 表应有8条记录，实际有 %d 条", dbType, count)
			}

			// 验证特定用户数据
			var username, password string
			if dbType == "postgres" {
				err = db.QueryRow("SELECT username, password FROM users WHERE id = $1", 1).Scan(&username, &password)
			} else {
				err = db.QueryRow("SELECT username, password FROM users WHERE id = ?", 1).Scan(&username, &password)
			}
			if err != nil {
				t.Fatalf("%s: 查询用户失败: %v", dbType, err)
			}
			if username != "Dumb" || password != "Dumb" {
				t.Errorf("%s: 用户数据不匹配: got %s/%s, want Dumb/Dumb", dbType, username, password)
			}
		})
	}
}

// TestDatabase_Reset 测试数据库重置
func TestDatabase_Reset(t *testing.T) {
	dbTypes := []string{"sqlite"}

	if os.Getenv("TEST_MYSQL_HOST") != "" {
		dbTypes = append(dbTypes, "mysql")
	}
	if os.Getenv("TEST_POSTGRES_HOST") != "" {
		dbTypes = append(dbTypes, "postgres")
	}

	for _, dbType := range dbTypes {
		t.Run(fmt.Sprintf("Reset_%s", dbType), func(t *testing.T) {
			cfg := getTestConfig(dbType)
			log, _ := logger.New(config.LogConfig{Level: "error", Output: "stdout"})
			defer log.Close()

			db, err := New(cfg, log)
			if err != nil {
				t.Fatalf("%s: 连接数据库失败: %v", dbType, err)
			}
			defer db.Close()

			// 先初始化
			if err := db.Migrate(); err != nil {
				t.Fatalf("%s: 迁移失败: %v", dbType, err)
			}
			if err := db.Seed(); err != nil {
				t.Fatalf("%s: 种子数据插入失败: %v", dbType, err)
			}

			// 执行重置
			if err := db.Reset(); err != nil {
				t.Fatalf("%s: 重置失败: %v", dbType, err)
			}

			// 验证数据已重置
			var count int
			if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
				t.Fatalf("%s: 查询 users 表失败: %v", dbType, err)
			}
			if count != 8 {
				t.Errorf("%s: 重置后 users 表应有8条记录，实际有 %d 条", dbType, count)
			}
		})
	}
}

// TestDatabase_Query 测试 LIMIT 语法兼容性
func TestDatabase_Query_LIMIT(t *testing.T) {
	dbTypes := []string{"sqlite"}

	if os.Getenv("TEST_MYSQL_HOST") != "" {
		dbTypes = append(dbTypes, "mysql")
	}
	if os.Getenv("TEST_POSTGRES_HOST") != "" {
		dbTypes = append(dbTypes, "postgres")
	}

	for _, dbType := range dbTypes {
		t.Run(fmt.Sprintf("Query_LIMIT_%s", dbType), func(t *testing.T) {
			cfg := getTestConfig(dbType)
			log, _ := logger.New(config.LogConfig{Level: "error", Output: "stdout"})
			defer log.Close()

			db, err := New(cfg, log)
			if err != nil {
				t.Fatalf("%s: 连接数据库失败: %v", dbType, err)
			}
			defer db.Close()

			// 先初始化
			if err := db.Migrate(); err != nil {
				t.Fatalf("%s: 迁移失败: %v", dbType, err)
			}
			if err := db.Seed(); err != nil {
				t.Fatalf("%s: 种子数据插入失败: %v", dbType, err)
			}

			// 测试 LIMIT 1 OFFSET 0 语法（三种数据库都支持）
			var id int
			var username string
			err = db.QueryRow("SELECT id, username FROM users LIMIT 1 OFFSET 0").Scan(&id, &username)
			if err != nil {
				t.Fatalf("%s: LIMIT 1 OFFSET 0 查询失败: %v", dbType, err)
			}
			if id != 1 || username != "Dumb" {
				t.Errorf("%s: 查询结果不匹配: got id=%d, username=%s", dbType, id, username)
			}
		})
	}
}
