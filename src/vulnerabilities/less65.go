package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less65 Less-65: Challenge 12 - User-Agent Based
type Less65 struct {
	BaseLesson
}

func (l *Less65) ID() string {
	return "less-65"
}

func (l *Less65) Name() string {
	return "Less-65: Challenge 12 - User-Agent Based"
}

func (l *Less65) Description() string {
	return "挑战关卡12：User-Agent注入"
}

func (l *Less65) Category() string {
	return "Challenge"
}

func (l *Less65) Route(r *gin.RouterGroup) {
	group := r.Group("/less-65")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less65) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-65")

	// Get user-agent header
	ua := c.GetHeader("User-Agent")
	if ua == "" {
		ua = "Mozilla/5.0"
	}

	query := fmt.Sprintf("SELECT * FROM challenge12 WHERE ua='%s' LIMIT 0,1", ua)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-65", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-65", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-65", result, "")
}

func (r *Registry) registerLess65() {
	r.Register(&Less65{BaseLesson{db: r.db, logger: r.logger}})
}
