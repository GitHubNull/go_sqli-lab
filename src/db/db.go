package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"go-sqli-lab/src/config"
	"go-sqli-lab/src/logger"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

// Database 数据库接口
type Database interface {
	GetDB() *sql.DB
	Migrate() error
	Seed() error
	Close() error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// defaultDatabase 默认数据库实现
type defaultDatabase struct {
	db     *sql.DB
	config config.DatabaseConfig
	logger logger.Logger
}

// New 创建数据库实例
func New(cfg config.DatabaseConfig, log logger.Logger) (Database, error) {
	var dsn string
	var driver string

	switch cfg.Type {
	case "sqlite":
		driver = "sqlite"
		dsn = cfg.DSN
		// 确保目录存在
		dir := filepath.Dir(dsn)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("创建数据库目录失败: %w", err)
		}
	case "mysql":
		driver = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	case "postgres":
		driver = "postgres"
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", cfg.Type)
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	log.Info("数据库连接成功", "type", cfg.Type)

	return &defaultDatabase{
		db:     db,
		config: cfg,
		logger: log,
	}, nil
}

func (d *defaultDatabase) GetDB() *sql.DB {
	return d.db
}

func (d *defaultDatabase) Close() error {
	return d.db.Close()
}

func (d *defaultDatabase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	d.logger.Debug("执行查询", "query", query, "args", args)
	return d.db.Query(query, args...)
}

func (d *defaultDatabase) QueryRow(query string, args ...interface{}) *sql.Row {
	d.logger.Debug("执行单行查询", "query", query, "args", args)
	return d.db.QueryRow(query, args...)
}

func (d *defaultDatabase) Exec(query string, args ...interface{}) (sql.Result, error) {
	d.logger.Debug("执行语句", "query", query, "args", args)
	return d.db.Exec(query, args...)
}

// Migrate 执行数据库迁移
func (d *defaultDatabase) Migrate() error {
	d.logger.Info("开始数据库迁移")

	// 创建users表
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(20) NOT NULL,
		password VARCHAR(20) NOT NULL
	);`
	if d.config.Type == "postgres" {
		usersTable = `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(20) NOT NULL,
			password VARCHAR(20) NOT NULL
		);`
	} else if d.config.Type == "mysql" {
		usersTable = `
		CREATE TABLE IF NOT EXISTS users (
			id INT NOT NULL AUTO_INCREMENT,
			username VARCHAR(20) NOT NULL,
			password VARCHAR(20) NOT NULL,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	}

	if _, err := d.db.Exec(usersTable); err != nil {
		return fmt.Errorf("创建users表失败: %w", err)
	}

	// 创建emails表
	emailsTable := `
	CREATE TABLE IF NOT EXISTS emails (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email_id VARCHAR(30) NOT NULL
	);`
	if d.config.Type == "postgres" {
		emailsTable = `
		CREATE TABLE IF NOT EXISTS emails (
			id SERIAL PRIMARY KEY,
			email_id VARCHAR(30) NOT NULL
		);`
	} else if d.config.Type == "mysql" {
		emailsTable = `
		CREATE TABLE IF NOT EXISTS emails (
			id INT NOT NULL AUTO_INCREMENT,
			email_id VARCHAR(30) NOT NULL,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	}

	if _, err := d.db.Exec(emailsTable); err != nil {
		return fmt.Errorf("创建emails表失败: %w", err)
	}

	// 创建uagents表
	uagentsTable := `
	CREATE TABLE IF NOT EXISTS uagents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uagent VARCHAR(256) NOT NULL,
		ip_address VARCHAR(35) NOT NULL,
		username VARCHAR(20) NOT NULL
	);`
	if d.config.Type == "postgres" {
		uagentsTable = `
		CREATE TABLE IF NOT EXISTS uagents (
			id SERIAL PRIMARY KEY,
			uagent VARCHAR(256) NOT NULL,
			ip_address VARCHAR(35) NOT NULL,
			username VARCHAR(20) NOT NULL
		);`
	} else if d.config.Type == "mysql" {
		uagentsTable = `
		CREATE TABLE IF NOT EXISTS uagents (
			id INT NOT NULL AUTO_INCREMENT,
			uagent VARCHAR(256) NOT NULL,
			ip_address VARCHAR(35) NOT NULL,
			username VARCHAR(20) NOT NULL,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	}

	if _, err := d.db.Exec(uagentsTable); err != nil {
		return fmt.Errorf("创建uagents表失败: %w", err)
	}

	// 创建referers表
	referersTable := `
	CREATE TABLE IF NOT EXISTS referers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		referer VARCHAR(256) NOT NULL,
		ip_address VARCHAR(35) NOT NULL
	);`
	if d.config.Type == "postgres" {
		referersTable = `
		CREATE TABLE IF NOT EXISTS referers (
			id SERIAL PRIMARY KEY,
			referer VARCHAR(256) NOT NULL,
			ip_address VARCHAR(35) NOT NULL
		);`
	} else if d.config.Type == "mysql" {
		referersTable = `
		CREATE TABLE IF NOT EXISTS referers (
			id INT NOT NULL AUTO_INCREMENT,
			referer VARCHAR(256) NOT NULL,
			ip_address VARCHAR(35) NOT NULL,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	}

	if _, err := d.db.Exec(referersTable); err != nil {
		return fmt.Errorf("创建referers表失败: %w", err)
	}

	// 创建挑战关卡表 (Less-54 ~ Less-65)
	for i := 1; i <= 12; i++ {
		tableName := fmt.Sprintf("challenge%d", i)
		challengeTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(20) NOT NULL,
			password VARCHAR(20) NOT NULL
		);`, tableName)

		if d.config.Type == "postgres" {
			challengeTable = fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				id SERIAL PRIMARY KEY,
				username VARCHAR(20) NOT NULL,
				password VARCHAR(20) NOT NULL
			);`, tableName)
		} else if d.config.Type == "mysql" {
			challengeTable = fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				id INT NOT NULL AUTO_INCREMENT,
				username VARCHAR(20) NOT NULL,
				password VARCHAR(20) NOT NULL,
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`, tableName)
		}

		if _, err := d.db.Exec(challengeTable); err != nil {
			return fmt.Errorf("创建%s表失败: %w", tableName, err)
		}
	}

	d.logger.Info("数据库迁移完成")
	return nil
}

// Seed 插入种子数据
func (d *defaultDatabase) Seed() error {
	d.logger.Info("开始插入种子数据")

	// 检查是否已有数据
	var count int
	if err := d.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return fmt.Errorf("检查users表数据失败: %w", err)
	}
	if count > 0 {
		d.logger.Info("数据已存在，跳过种子数据插入")
		return nil
	}

	// 插入users数据
	users := []struct {
		id       int
		username string
		password string
	}{
		{1, "Dumb", "Dumb"},
		{2, "Angelina", "I-kill-you"},
		{3, "Dummy", "p@ssword"},
		{4, "secure", "crappy"},
		{5, "stupid", "stupidity"},
		{6, "superman", "genious"},
		{7, "batman", "mob!le"},
		{8, "admin", "admin"},
	}

	for _, u := range users {
		_, err := d.db.Exec(
			"INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
			u.id, u.username, u.password,
		)
		if err != nil {
			// PostgreSQL使用$1, $2等参数占位符
			if d.config.Type == "postgres" {
				_, err = d.db.Exec(
					"INSERT INTO users (id, username, password) VALUES ($1, $2, $3)",
					u.id, u.username, u.password,
				)
			}
			if err != nil {
				return fmt.Errorf("插入用户数据失败: %w", err)
			}
		}
	}

	// 插入emails数据
	emails := []struct {
		id      int
		emailID string
	}{
		{1, "Dumb@dhakkan.com"},
		{2, "Angel@iloveu.com"},
		{3, "Dummy@dhakkan.local"},
		{4, "secure@dhakkan.local"},
		{5, "stupid@dhakkan.local"},
		{6, "superman@dhakkan.local"},
		{7, "batman@dhakkan.local"},
		{8, "admin@dhakkan.com"},
	}

	for _, e := range emails {
		_, err := d.db.Exec(
			"INSERT INTO emails (id, email_id) VALUES (?, ?)",
			e.id, e.emailID,
		)
		if err != nil {
			if d.config.Type == "postgres" {
				_, err = d.db.Exec(
					"INSERT INTO emails (id, email_id) VALUES ($1, $2)",
					e.id, e.emailID,
				)
			}
			if err != nil {
				return fmt.Errorf("插入email数据失败: %w", err)
			}
		}
	}

	d.logger.Info("种子数据插入完成")
	return nil
}
