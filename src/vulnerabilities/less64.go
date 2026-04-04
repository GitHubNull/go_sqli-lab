package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less64 Less-64: Challenge 11 - Cookie Based
type Less64 struct {
	BaseLesson
}

func (l *Less64) ID() string {
	return "less-64"
}

func (l *Less64) Name() string {
	return "Less-64: Challenge 11 - Cookie Based"
}

func (l *Less64) Description() string {
	return "挑战关卡11：Cookie注入"
}

func (l *Less64) Category() string {
	return "Challenge"
}

func (l *Less64) Route(r *gin.RouterGroup) {
	group := r.Group("/less-64")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less64) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-64")

	// Get id from cookie
	id, err := c.Cookie("id")
	if err != nil || id == "" {
		id = "1"
		c.SetCookie("id", id, 3600, "/", "", false, true)
	}

	query := fmt.Sprintf("SELECT * FROM challenge11 WHERE id='%s' LIMIT 1 OFFSET 0", id)

	var user models.User
	err = l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-64", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-64", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-64", result, "")
}

func (r *Registry) registerLess64() {
	r.Register(&Less64{BaseLesson{db: r.db, logger: r.logger}})
}
