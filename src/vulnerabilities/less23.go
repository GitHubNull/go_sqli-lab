package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less23 GET - Error Based - Comments
type Less23 struct {
	BaseLesson
}

func (l *Less23) ID() string {
	return "less-23"
}

func (l *Less23) Name() string {
	return "Less-23: GET - Error Based - Comments"
}

func (l *Less23) Description() string {
	return "GET请求 - 基于注释过滤的SQL注入"
}

func (l *Less23) Category() string {
	return "WAF Bypass"
}

func (l *Less23) Route(r *gin.RouterGroup) {
	group := r.Group("/less-23")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less23) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-23")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-23", "请输入ID参数", "")
		return
	}

	// 过滤注释符
	id = filterComments(id)

	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 0,1", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-23", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-23", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-23", result, "")
}

// filterComments 过滤SQL注释符
func filterComments(input string) string {
	result := input
	// 移除各种注释符
	result = replaceAll(result, "--", "")
	result = replaceAll(result, "#", "")
	result = replaceAll(result, "/*", "")
	result = replaceAll(result, "*/", "")
	return result
}

func replaceAll(s, old, new string) string {
	for {
		ns := ""
		found := false
		for i := 0; i < len(s); {
			if i+len(old) <= len(s) && s[i:i+len(old)] == old {
				ns += new
				i += len(old)
				found = true
			} else {
				ns += string(s[i])
				i++
			}
		}
		s = ns
		if !found {
			break
		}
	}
	return s
}

func (r *Registry) registerLess23() {
	r.Register(&Less23{BaseLesson{db: r.db, logger: r.logger}})
}
