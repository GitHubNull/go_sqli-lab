package vulnerabilities

import (
	"database/sql"
	"fmt"
	"strings"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less59 Less-59: Challenge 6 - WAF Bypass
type Less59 struct {
	BaseLesson
}

func (l *Less59) ID() string {
	return "less-59"
}

func (l *Less59) Name() string {
	return "Less-59: Challenge 6 - WAF Bypass"
}

func (l *Less59) Description() string {
	return "挑战关卡6：WAF绕过注入"
}

func (l *Less59) Category() string {
	return "Challenge"
}

func (l *Less59) Route(r *gin.RouterGroup) {
	group := r.Group("/less-59")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less59) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-59")

	id := c.Query("id")
	if id == "" {
		id = "1"
	}

	// Simple WAF: block certain keywords
	blocked := []string{"union", "select", "sleep", "benchmark"}
	lowerID := strings.ToLower(id)
	for _, keyword := range blocked {
		if strings.Contains(lowerID, keyword) {
			l.renderHTML(c, "less-59", "", "WAF blocked suspicious input")
			return
		}
	}

	query := fmt.Sprintf("SELECT * FROM challenge6 WHERE id='%s' LIMIT 0,1", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-59", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-59", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-59", result, "")
}

func (r *Registry) registerLess59() {
	r.Register(&Less59{BaseLesson{db: r.db, logger: r.logger}})
}
