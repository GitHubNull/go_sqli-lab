package vulnerabilities

import (
	"database/sql"
	"fmt"
	"time"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less9 Less-9漏洞: Blind - Time Based - Single Quotes
type Less9 struct {
	BaseLesson
}

func (l *Less9) ID() string {
	return "less-9"
}

func (l *Less9) Name() string {
	return "Less-9: Blind - Time Based - Single Quotes"
}

func (l *Less9) Description() string {
	return "盲注 - 基于时间的SQL注入 (单引号)"
}

func (l *Less9) Category() string {
	return "Time Based"
}

func (l *Less9) Route(r *gin.RouterGroup) {
	group := r.Group("/less-9")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less9) handleIndex(c *gin.Context) {
	start := time.Now()
	l.logRequest(c, "less-9")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-9", "请输入ID参数", "")
		return
	}

	// 在SQLite中不支持SLEEP函数，我们通过其他方式模拟
	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	elapsed := time.Since(start)

	// 时间盲注特点：无论成功与否都显示相同信息，但通过响应时间判断
	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-9", "You are in...........", "")
		} else {
			l.renderHTML(c, "less-9", "You are in...........", "")
		}
		return
	}

	// 添加延迟信息到响应头（用于时间盲注测试）
	c.Header("X-Response-Time", fmt.Sprintf("%dms", elapsed.Milliseconds()))
	l.renderHTML(c, "less-9", "You are in...........", "")
}

func (r *Registry) registerLess9() {
	r.Register(&Less9{BaseLesson{db: r.db, logger: r.logger}})
}
