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
	Reset() error
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

	// 创建简历（应聘人员求职登记表——Less-68/69）
	resumesTable := `
	CREATE TABLE IF NOT EXISTS resumes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(50) NOT NULL,
		sex VARCHAR(10) NOT NULL,
		birth VARCHAR(30),
		nation VARCHAR(30),
		native_place VARCHAR(100),
		political VARCHAR(30),
		health VARCHAR(20),
		marital VARCHAR(20),
		education VARCHAR(30),
		school VARCHAR(100),
		major VARCHAR(80),
		graduation VARCHAR(30),
		phone VARCHAR(30),
		email VARCHAR(100),
		postal_code VARCHAR(12),
		work_years VARCHAR(20),
		job_intention VARCHAR(100),
		salary_expect VARCHAR(40),
		address VARCHAR(200),
		self_eval TEXT
	);`
	if d.config.Type == "postgres" {
		resumesTable = `
		CREATE TABLE IF NOT EXISTS resumes (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			sex VARCHAR(10) NOT NULL,
			birth VARCHAR(30),
			nation VARCHAR(30),
			native_place VARCHAR(100),
			political VARCHAR(30),
			health VARCHAR(20),
			marital VARCHAR(20),
			education VARCHAR(30),
			school VARCHAR(100),
			major VARCHAR(80),
			graduation VARCHAR(30),
			phone VARCHAR(30),
			email VARCHAR(100),
			postal_code VARCHAR(12),
			work_years VARCHAR(20),
			job_intention VARCHAR(100),
			salary_expect VARCHAR(40),
			address VARCHAR(200),
			self_eval TEXT
		);`
	} else if d.config.Type == "mysql" {
		resumesTable = `
		CREATE TABLE IF NOT EXISTS resumes (
			id INT NOT NULL AUTO_INCREMENT,
			name VARCHAR(50) NOT NULL,
			sex VARCHAR(10) NOT NULL,
			birth VARCHAR(30),
			nation VARCHAR(30),
			native_place VARCHAR(100),
			political VARCHAR(30),
			health VARCHAR(20),
			marital VARCHAR(20),
			education VARCHAR(30),
			school VARCHAR(100),
			major VARCHAR(80),
			graduation VARCHAR(30),
			phone VARCHAR(30),
			email VARCHAR(100),
			postal_code VARCHAR(12),
			work_years VARCHAR(20),
			job_intention VARCHAR(100),
			salary_expect VARCHAR(40),
			address VARCHAR(200),
			self_eval TEXT,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	}

	if _, err := d.db.Exec(resumesTable); err != nil {
		return fmt.Errorf("创建resumes表失败: %w", err)
	}

	d.logger.Info("数据库迁移完成")
	return nil
}

// Seed 插入种子数据
func (d *defaultDatabase) Seed() error {
	d.logger.Info("开始插入种子数据")

	// 简历表独立于 users 表判断：即便 users 已有数据，也要保证 resumes 已种入
	if err := d.seedResumes(); err != nil {
		return fmt.Errorf("插入简历数据失败: %w", err)
	}

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

// Reset 重置数据库 - 清除所有数据并重新初始化
func (d *defaultDatabase) Reset() error {
	d.logger.Info("开始重置数据库")

	tables := []string{"users", "emails", "uagents", "referers", "resumes"}
	// 添加挑战关卡表
	for i := 1; i <= 12; i++ {
		tables = append(tables, fmt.Sprintf("challenge%d", i))
	}

	// 清除所有表的数据
	for _, table := range tables {
		if err := d.clearTable(table); err != nil {
			d.logger.Error("清除表数据失败", "table", table, "error", err)
			return fmt.Errorf("清除表 %s 失败: %w", table, err)
		}
		d.logger.Info("已清除表数据", "table", table)
	}

	// 重新插入种子数据
	if err := d.forceSeed(); err != nil {
		return fmt.Errorf("重新插入种子数据失败: %w", err)
	}

	// 插入挑战关卡种子数据
	if err := d.seedChallenges(); err != nil {
		return fmt.Errorf("插入挑战关卡数据失败: %w", err)
	}

	d.logger.Info("数据库重置完成")
	return nil
}

// clearTable 清除指定表的数据
func (d *defaultDatabase) clearTable(tableName string) error {
	var query string

	switch d.config.Type {
	case "sqlite":
		// SQLite 支持 DELETE 和 TRUNCATE（通过 DELETE 实现）
		query = fmt.Sprintf("DELETE FROM %s", tableName)
	case "mysql":
		// MySQL 支持 TRUNCATE
		query = fmt.Sprintf("TRUNCATE TABLE %s", tableName)
	case "postgres":
		// PostgreSQL 支持 TRUNCATE
		query = fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY", tableName)
	default:
		query = fmt.Sprintf("DELETE FROM %s", tableName)
	}

	_, err := d.db.Exec(query)
	return err
}

// forceSeed 强制插入种子数据（不检查数据是否存在）
func (d *defaultDatabase) forceSeed() error {
	d.logger.Info("强制插入种子数据")

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
		var err error
		if d.config.Type == "postgres" {
			_, err = d.db.Exec(
				"INSERT INTO users (id, username, password) VALUES ($1, $2, $3)",
				u.id, u.username, u.password,
			)
		} else {
			_, err = d.db.Exec(
				"INSERT INTO users (id, username, password) VALUES (?, ?, ?)",
				u.id, u.username, u.password,
			)
		}
		if err != nil {
			return fmt.Errorf("插入用户数据失败: %w", err)
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
		var err error
		if d.config.Type == "postgres" {
			_, err = d.db.Exec(
				"INSERT INTO emails (id, email_id) VALUES ($1, $2)",
				e.id, e.emailID,
			)
		} else {
			_, err = d.db.Exec(
				"INSERT INTO emails (id, email_id) VALUES (?, ?)",
				e.id, e.emailID,
			)
		}
		if err != nil {
			return fmt.Errorf("插入email数据失败: %w", err)
		}
	}

	if err := d.seedResumes(); err != nil {
		return fmt.Errorf("插入简历数据失败: %w", err)
	}

	d.logger.Info("种子数据插入完成")
	return nil
}

// seedChallenges 为挑战关卡表插入种子数据
func (d *defaultDatabase) seedChallenges() error {
	d.logger.Info("插入挑战关卡种子数据")

	// 为每个挑战表插入一些随机数据
	challengeUsers := []struct {
		username string
		password string
	}{
		{"challenge_user1", "pass1"},
		{"challenge_user2", "pass2"},
		{"challenge_user3", "pass3"},
		{"challenge_user4", "pass4"},
		{"challenge_user5", "pass5"},
	}

	for i := 1; i <= 12; i++ {
		tableName := fmt.Sprintf("challenge%d", i)

		// 为每个挑战表插入1-3个随机用户
		for j, u := range challengeUsers {
			if j >= 3+i%3 { // 每个表插入不同数量的用户
				break
			}

			var err error
			if d.config.Type == "postgres" {
				_, err = d.db.Exec(
					fmt.Sprintf("INSERT INTO %s (id, username, password) VALUES ($1, $2, $3)", tableName),
					j+1, u.username, u.password,
				)
			} else {
				_, err = d.db.Exec(
					fmt.Sprintf("INSERT INTO %s (id, username, password) VALUES (?, ?, ?)", tableName),
					j+1, u.username, u.password,
				)
			}
			if err != nil {
				return fmt.Errorf("插入%s数据失败: %w", tableName, err)
			}
		}
	}

	d.logger.Info("挑战关卡数据插入完成")
	return nil
}

// seedResumes 为简历（应聘人员求职登记表）插入真实的中文测试数据。
// 仅在 resumes 表为空时插入，避免重复。
func (d *defaultDatabase) seedResumes() error {
	var count int
	if err := d.db.QueryRow("SELECT COUNT(*) FROM resumes").Scan(&count); err != nil {
		return fmt.Errorf("检查resumes表数据失败: %w", err)
	}
	if count > 0 {
		d.logger.Info("resumes 表已有数据，跳过种子录入")
		return nil
	}

	d.logger.Info("插入简历种子数据")

	// 字段顺序：name, sex, birth, nation, native_place, political, health, marital,
	//           education, school, major, graduation, phone, email, postal_code,
	//           work_years, job_intention, salary_expect, address, self_eval
	resumes := []struct {
		name, sex, birth, nation, nativePlace, political, health, marital string
		education, school, major, graduation, phone, email, postalCode    string
		workYears, jobIntention, salaryExpect, address, selfEval          string
	}{
		{
			"张伟", "男", "1992-05-18", "汉族", "湖南长沙",
			"中共党员", "健康", "已婚",
			"本科", "中南大学", "计算机科学与技术", "2014-06",
			"138-0731-5566", "zhangwei@example.com", "410000",
			"8年", "Java 高级开发工程师", "面议",
			"湖南省长沙市岳麓区麓山南路 932 号",
			"热爱技术，8 年 Java 后端研发经验，主导过多个高并发、分布式电商系统；熟悉 Spring Cloud、MySQL 调优、Redis 缓存与消息队列；具备良好的系统设计能力与团队协作精神。",
		},
		{
			"李娜", "女", "1995-11-02", "汉族", "四川成都",
			"共青团员", "健康", "未婚",
			"硕士", "四川大学", "软件工程", "2021-06",
			"139-8022-3311", "lina@example.com", "610000",
			"3年", "前端开发工程师", "12-18K",
			"四川省成都市武侯区科华北路 65 号",
			"3 年前端开发经验，精通 Vue3、TypeScript、Vite；对组件化、微前端、性能优化有较深理解；业余喜欢开源与写技术博客，沟通表达能力强。",
		},
		{
			"王强", "男", "1988-09-23", "满族", "辽宁沈阳",
			"群众", "良好", "未婚",
			"本科", "大连理工大学", "自动化", "2011-06",
			"136-0411-8899", "wangqiang@example.com", "110000",
			"12年", "嵌入式软件工程师", "18-25K",
			"辽宁省沈阳市和平区文化路 3 号",
			"长期从事嵌入式 Linux 驱动与 BSP 开发，熟悉 ARM、RTOS、C/C++；参与过车载、工控等多类量产项目，有较强的硬件问题定位能力。",
		},
		{
			"赵敏", "女", "1997-01-15", "汉族", "浙江杭州",
			"中共党员", "健康", "未婚",
			"本科", "浙江大学", "会计学", "2018-06",
			"187-5712-6644", "zhaomin@example.com", "310000",
			"5年", "财务主管", "面议",
			"浙江省杭州市西湖区文三路 478 号",
			"注册会计（CPA），5 年财务经验，熟悉全套账务、税务筹划与财务分析；擅长用数据支撑经营决策，责任心强，抗压能力好。",
		},
		{
			"刘洋", "男", "1999-03-08", "回族", "宁夏银川",
			"共青团员", "健康", "未婚",
			"本科", "西安电子科技大学", "通信工程", "2021-06",
			"131-0951-2233", "liuyang@example.com", "750000",
			"2年", "测试开发工程师", "9-13K",
			"宁夏回族自治区银川市金凤区北京中路 168 号",
			"熟悉 Python、pytest、接口自动化与性能测试工具；对质量保障有热情，善于从测试角度发现设计与实现缺陷，细心耐心。",
		},
		{
			"孙芳", "女", "1990-07-30", "汉族", "湖北武汉",
			"群众", "健康", "已婚",
			"大专", "武汉职业技术学院", "电子商务", "2012-06",
			"150-0712-7788", "sunfang@example.com", "430000",
			"10年", "运营经理", "面议",
			"湖北省武汉市洪山区珞喻路 1037 号",
			"10 年电商运营与团队管理经验，擅长用户增长、活动策划与数据分析；独立操盘过千万级 GMV 店铺，熟悉主流平台规则与打法。",
		},
	}

	// 组装 SQL 占位符：跨库统一使用位置参数 (?, ?...)，PostgreSQL 使用 $1..$21
	var placeholders string
	for i := 1; i <= 20; i++ {
		if i > 1 {
			placeholders += ", "
		}
		if d.config.Type == "postgres" {
			placeholders += fmt.Sprintf("$%d", i)
		} else {
			placeholders += "?"
		}
	}

	insertSQL := fmt.Sprintf(
		"INSERT INTO resumes (name, sex, birth, nation, native_place, political, health, marital, "+
			"education, school, major, graduation, phone, email, postal_code, "+
			"work_years, job_intention, salary_expect, address, self_eval) VALUES (%s)", placeholders)

	for _, r := range resumes {
		if _, err := d.db.Exec(insertSQL,
			r.name, r.sex, r.birth, r.nation, r.nativePlace, r.political, r.health, r.marital,
			r.education, r.school, r.major, r.graduation, r.phone, r.email, r.postalCode,
			r.workYears, r.jobIntention, r.salaryExpect, r.address, r.selfEval,
		); err != nil {
			return fmt.Errorf("插入简历数据失败: %w", err)
		}
	}

	d.logger.Info("简历种子数据插入完成", "count", len(resumes))
	return nil
}
