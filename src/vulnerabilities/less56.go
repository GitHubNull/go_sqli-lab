package vulnerabilities

import (
	"database/sql"
	"fmt"
	"time"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less56 Less-56: Challenge 3 - Time Based Blind
type Less56 struct {
	BaseLesson
}

func (l *Less56) ID() string {
	return "less-56"
}

func (l *Less56) Name() string {
	return "Less-56: Challenge 3 - Time Blind"
}

func (l *Less56) Description() string {
	return "挑战关卡3：时间延迟盲注"
}

func (l *Less56) Category() string {
	return "Challenge"
}

func (l *Less56) Route(r *gin.RouterGroup) {
	group := r.Group("/less-56")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less56) handleIndex(c *gin.Context) {
	start := time.Now()
	l.logRequest(c, "less-56")

	id := c.Query("id")
	if id == "" {
		id = "1"
	}

	query := fmt.Sprintf("SELECT * FROM challenge3 WHERE id='%s' LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-56", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-56", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s<br>Query Time: %v",
		user.Username, user.Password, time.Since(start))
	l.renderHTML(c, "less-56", result, "")
}

func (r *Registry) registerLess56() {
	r.Register(&Less56{BaseLesson{db: r.db, logger: r.logger}})
}
