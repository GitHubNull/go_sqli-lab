package vulnerabilities

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less60 Less-60: Challenge 7 - JSON Injection
type Less60 struct {
	BaseLesson
}

func (l *Less60) ID() string {
	return "less-60"
}

func (l *Less60) Name() string {
	return "Less-60: Challenge 7 - JSON Injection"
}

func (l *Less60) Description() string {
	return "挑战关卡7：JSON参数注入"
}

func (l *Less60) Category() string {
	return "Challenge"
}

func (l *Less60) Route(r *gin.RouterGroup) {
	group := r.Group("/less-60")
	{
		group.POST("", l.handleIndex)
		group.POST("/", l.handleIndex)
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less60) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-60")

	var data map[string]interface{}

	// Try to get JSON from body or query
	if c.Request.Method == "POST" {
		if err := c.ShouldBindJSON(&data); err != nil {
			// Try from form data
			jsonStr := c.PostForm("data")
			if jsonStr == "" {
				jsonStr = c.Query("data")
			}
			if jsonStr != "" {
				json.Unmarshal([]byte(jsonStr), &data)
			}
		}
	} else {
		jsonStr := c.Query("data")
		if jsonStr != "" {
			json.Unmarshal([]byte(jsonStr), &data)
		}
	}

	if data == nil {
		data = make(map[string]interface{})
	}

	id, ok := data["id"].(string)
	if !ok || id == "" {
		id = "1"
	}

	// Prevent direct injection characters in JSON
	id = strings.ReplaceAll(id, "'", "")

	query := fmt.Sprintf("SELECT * FROM challenge7 WHERE id=%s LIMIT 0,1", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-60", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-60", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-60", result, "")
}

func (r *Registry) registerLess60() {
	r.Register(&Less60{BaseLesson{db: r.db, logger: r.logger}})
}
