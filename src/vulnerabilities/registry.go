package vulnerabilities

import (
	"fmt"

	"go-sqli-lab/src/db"
	"go-sqli-lab/src/logger"

	"github.com/gin-gonic/gin"
)

// Registry 漏洞注册表
type Registry struct {
	db              db.Database
	logger          logger.Logger
	vulnerabilities map[string]Vulnerability
}

// NewRegistry 创建漏洞注册表
func NewRegistry(database db.Database, log logger.Logger) *Registry {
	return &Registry{
		db:              database,
		logger:          log,
		vulnerabilities: make(map[string]Vulnerability),
	}
}

// Register 注册漏洞
func (r *Registry) Register(v Vulnerability) {
	r.vulnerabilities[v.ID()] = v
	r.logger.Info("注册漏洞", "id", v.ID(), "name", v.Name())
}

// Get 获取漏洞
func (r *Registry) Get(id string) (Vulnerability, bool) {
	v, ok := r.vulnerabilities[id]
	return v, ok
}

// GetAll 获取所有漏洞
func (r *Registry) GetAll() map[string]Vulnerability {
	return r.vulnerabilities
}

// RegisterAll 注册所有漏洞路由
func (r *Registry) RegisterAll(router *gin.Engine) {
	// 注册各个级别的漏洞
	r.registerLess1()
	r.registerLess2()
	r.registerLess3()
	r.registerLess4()
	r.registerLess5()
	r.registerLess6()
	r.registerLess7()
	r.registerLess8()
	r.registerLess9()
	r.registerLess10()
	r.registerLess11()
	r.registerLess12()
	r.registerLess13()
	r.registerLess14()
	r.registerLess15()
	r.registerLess16()
	r.registerLess17()
	r.registerLess18()
	r.registerLess19()
	r.registerLess20()
	r.registerLess21()
	r.registerLess22()
	r.registerLess23()
	r.registerLess24()
	r.registerLess25()
	r.registerLess26()
	r.registerLess27()
	r.registerLess28()
	r.registerLess29()
	r.registerLess30()
	r.registerLess31()
	r.registerLess32()
	r.registerLess33()
	r.registerLess34()
	r.registerLess35()
	r.registerLess36()
	r.registerLess37()
	r.registerLess38()
	r.registerLess39()
	r.registerLess40()
	r.registerLess41()
	r.registerLess42()
	r.registerLess43()
	r.registerLess44()
	r.registerLess45()
	r.registerLess46()
	r.registerLess47()
	r.registerLess48()
	r.registerLess49()
	r.registerLess50()
	r.registerLess51()
	r.registerLess52()
	r.registerLess53()
	r.registerLess54()
	r.registerLess55()
	r.registerLess56()
	r.registerLess57()
	r.registerLess58()
	r.registerLess59()
	r.registerLess60()
	r.registerLess61()
	r.registerLess62()
	r.registerLess63()
	r.registerLess64()
	r.registerLess65()
	r.registerLess66()
	r.registerLess67()

	// 注册漏洞路由组
	vulnGroup := router.Group("/less")
	for _, v := range r.vulnerabilities {
		v.Route(vulnGroup)
	}

	// 漏洞列表API
	router.GET("/api/vulnerabilities", r.listVulnerabilities)
	router.GET("/api/vulnerabilities/:id", r.getVulnerability)
}

// listVulnerabilities 获取漏洞列表
func (r *Registry) listVulnerabilities(c *gin.Context) {
	list := make([]gin.H, 0, len(r.vulnerabilities))
	for _, v := range r.vulnerabilities {
		list = append(list, gin.H{
			"id":          v.ID(),
			"name":        v.Name(),
			"description": v.Description(),
			"category":    v.Category(),
		})
	}
	c.JSON(200, gin.H{
		"data":  list,
		"total": len(list),
	})
}

// getVulnerability 获取单个漏洞信息
func (r *Registry) getVulnerability(c *gin.Context) {
	id := c.Param("id")
	v, ok := r.vulnerabilities[id]
	if !ok {
		c.JSON(404, gin.H{"error": fmt.Sprintf("漏洞 %s 不存在", id)})
		return
	}
	c.JSON(200, gin.H{
		"id":          v.ID(),
		"name":        v.Name(),
		"description": v.Description(),
		"category":    v.Category(),
	})
}
