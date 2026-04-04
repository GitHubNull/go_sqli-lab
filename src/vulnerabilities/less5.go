package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less5 Less-5漏洞: Double Query - Single Quotes
type Less5 struct {
	BaseLesson
}

func (l *Less5) ID() string {
	return "less-5"
}

func (l *Less5) Name() string {
	return "Less-5: Double Query - Single Quotes"
}

func (l *Less5) Description() string {
	return "双查询注入 - 单引号字符串类型"
}

func (l *Less5) Category() string {
	return "Double Injection"
}

func (l *Less5) Route(r *gin.RouterGroup) {
	group := r.Group("/less-5")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less5) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-5")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-5", "请输入ID参数", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-5", "", "You are in...........")
		} else {
			l.renderHTML(c, "less-5", "", err.Error())
		}
		return
	}

	l.renderHTML(c, "less-5", "You are in...........", "")
}

func (r *Registry) registerLess5() {
	r.Register(&Less5{BaseLesson{db: r.db, logger: r.logger}})
}
