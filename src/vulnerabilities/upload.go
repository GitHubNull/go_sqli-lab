package vulnerabilities

import (
	"bytes"
	_ "embed"
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
	"strings"

	"go-sqli-lab/src/db"

	"github.com/extrame/xls"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// 嵌入预生成的 .xls 上传模板（真正的二进制 .xls 文件，extrame/xls 可解析）。
//
//go:embed assets/upload_template.xls
var uploadTemplateXLS []byte

// ---------------------------------------------------------------------------
// Excel 解析：.xls（旧版二进制格式）与 .xlsx（OOXML 格式）
// ---------------------------------------------------------------------------

// parseXLS 解析旧版 .xls 二进制文件，返回第一个 Sheet 的所有行数据。
// 使用 extrame/xls 纯 Go 库，无需 CGO。
func parseXLS(data []byte) ([][]string, error) {
	wb, err := xls.OpenReader(bytes.NewReader(data), "utf-8")
	if err != nil {
		return nil, fmt.Errorf("解析 .xls 文件失败: %w", err)
	}

	sheet := wb.GetSheet(0)
	if sheet == nil {
		return nil, fmt.Errorf("文件中没有工作表")
	}

	var rows [][]string
	for i := 0; i <= int(sheet.MaxRow); i++ {
		row := sheet.Row(i)
		if row == nil {
			continue
		}
		var cols []string
		for j := row.FirstCol(); j < row.LastCol(); j++ {
			cols = append(cols, row.Col(j))
		}
		rows = append(rows, cols)
	}
	return rows, nil
}

// parseXLSX 解析 .xlsx 文件，返回第一个 Sheet 的所有行数据。
// 使用项目已有依赖 excelize/v2。
func parseXLSX(data []byte) ([][]string, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("解析 .xlsx 文件失败: %w", err)
	}
	defer f.Close()

	sheet := f.GetSheetName(0)
	if sheet == "" {
		return nil, fmt.Errorf("文件中没有工作表")
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("读取工作表数据失败: %w", err)
	}
	return rows, nil
}

// ---------------------------------------------------------------------------
// SQL 执行：批量验证场景（漏洞点：单元格值直接拼接进 SQL）
// ---------------------------------------------------------------------------

// uploadRowResult 保存一行数据的执行结果
type uploadRowResult struct {
	RowNum  int        // 行号（从 1 开始，含表头）
	Keyword string     // 该行的关键字（第一列值）
	RawSQL  string     // 拼接后的 SQL
	Cols    []string   // 查询列名
	Rows    [][]string // 查询结果
	Err     string     // 错误信息
}

// runUploadQuery 对单个关键字执行「批量验证」查询。
// 漏洞点：keyword 被直接拼接进 SQL（与 Less-1 同款写法），可被 SQL 注入。
func runUploadQuery(database db.Database, keyword string) (cols []string, rows [][]string, rawSQL string, err error) {
	rawSQL = fmt.Sprintf("SELECT id, username, password FROM users WHERE username LIKE '%%%s%%'", keyword)

	rs, qerr := database.Query(rawSQL)
	if qerr != nil {
		return nil, nil, rawSQL, qerr
	}
	defer rs.Close()

	cols, cerr := rs.Columns()
	if cerr != nil {
		return nil, nil, rawSQL, cerr
	}

	for rs.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if serr := rs.Scan(ptrs...); serr != nil {
			return nil, nil, rawSQL, serr
		}
		row := make([]string, len(cols))
		for i, v := range vals {
			row[i] = valueToString(v)
		}
		rows = append(rows, row)
	}
	if rerr := rs.Err(); rerr != nil {
		return nil, nil, rawSQL, rerr
	}

	return cols, rows, rawSQL, nil
}

// processUploadData 对解析出的 Excel 数据逐行执行 SQL 查询。
// 跳过第一行（表头），取每行第一列作为 keyword。
func processUploadData(database db.Database, data [][]string) []uploadRowResult {
	var results []uploadRowResult
	for i, row := range data {
		if i == 0 {
			continue // 跳过表头
		}
		if len(row) == 0 || strings.TrimSpace(row[0]) == "" {
			continue // 跳过空行
		}

		keyword := strings.TrimSpace(row[0])
		cols, rows, rawSQL, err := runUploadQuery(database, keyword)

		r := uploadRowResult{
			RowNum:  i + 1,
			Keyword: keyword,
			RawSQL:  rawSQL,
			Cols:    cols,
			Rows:    rows,
		}
		if err != nil {
			r.Err = err.Error()
		}
		results = append(results, r)
	}
	return results
}

// ---------------------------------------------------------------------------
// 模板下载：生成预置的 Excel 模板文件
// ---------------------------------------------------------------------------

// serveUploadTemplate 生成并下载上传模板文件。
// format 为 "XLS" 或 "XLSX"，决定生成的文件格式。
// XLS 使用嵌入的预生成二进制文件（extrame/xls 可解析），XLSX 动态生成。
func serveUploadTemplate(c *gin.Context, format string) {
	if format == "XLSX" {
		cols := []string{"username"}
		rows := [][]string{
			{"admin"},
			{"Dumb"},
		}
		data, err := buildXLSX(cols, rows)
		if err != nil {
			c.String(http.StatusInternalServerError, "生成模板失败: %v", err)
			return
		}
		serveExport(c, "upload_template.xlsx", contentTypeXLSX, data)
	} else {
		serveExport(c, "upload_template.xls", contentTypeXLS, uploadTemplateXLS)
	}
}

// ---------------------------------------------------------------------------
// 页面渲染：上传关卡共用
// ---------------------------------------------------------------------------

// renderUploadPage 渲染文件上传关卡的入口页面。
// format 用于按钮文案和 accept 属性（如 "XLS"/"XLSX"），acceptExt 为文件扩展名（如 ".xls"/".xlsx"）。
// uploaded 为 true 时展示批量验证结果。
func renderUploadPage(c *gin.Context, lessonID, title, desc, format, acceptExt string,
	results []uploadRowResult, errMsg string, uploaded bool) {

	templateURL := "/less/" + lessonID + "/template"

	var b strings.Builder
	b.WriteString(`<!DOCTYPE html><html lang="zh-CN"><head><meta charset="UTF-8">`)
	b.WriteString(`<meta name="viewport" content="width=device-width, initial-scale=1.0">`)
	b.WriteString(`<title>` + html.EscapeString(title) + `</title>`)
	b.WriteString(`<link rel="stylesheet" href="/static/css/style.css"></head><body><div class="container">`)
	b.WriteString(`<h1>` + html.EscapeString(title) + `</h1>`)
	b.WriteString(`<p style="text-align:center;color:#888;">` + html.EscapeString(desc) + `</p>`)

	// 漏洞说明与示例 payload
	b.WriteString(`<div class="content" style="max-width:760px;margin:20px auto;">`)
	b.WriteString(`<h2>漏洞点</h2>`)
	b.WriteString(`<p>后端解析上传的 Excel 文件，将单元格内容直接拼接进 SQL 查询（未做任何过滤/参数化），存在 SQL 注入。</p>`)
	b.WriteString(`<pre style="background:rgba(0,0,0,.4);padding:12px;border-radius:6px;overflow:auto;">` +
		`SELECT id, username, password FROM users WHERE username LIKE '%[单元格值]%'</pre>`)
	b.WriteString(`<h2>测试示例</h2><ul style="padding-left:20px;line-height:1.9;">`)
	b.WriteString(`<li>报错注入：<code>'</code></li>`)
	b.WriteString(`<li>UNION 注入（SQLite/MySQL）：<code>%' UNION SELECT 1,group_concat(username),3 FROM users -- </code></li>`)
	b.WriteString(`<li>UNION 注入（PostgreSQL）：<code>%' UNION SELECT 1,string_agg(username,','),3 FROM users -- </code></li>`)
	b.WriteString(`<li>布尔盲注：<code>admin' AND 1=1 -- </code> 对比 <code>admin' AND 1=2 -- </code></li>`)
	b.WriteString(`</ul></div>`)

	// 模板下载
	b.WriteString(`<div class="nav"><a href="` + templateURL + `">📥 下载 ` +
		html.EscapeString(format) + ` 模板文件</a></div>`)

	// 上传表单
	b.WriteString(`<form method="POST" action="/less/` + lessonID + `/upload" enctype="multipart/form-data">`)
	b.WriteString(`<div><label>选择 ` + html.EscapeString(format) + ` 文件：</label>`)
	b.WriteString(`<input type="file" name="file" accept="` + acceptExt + `" required style="width:100%;padding:12px;border:2px solid #667eea;border-radius:5px;background:rgba(255,255,255,0.1);color:#fff;font-size:16px;"></div>`)
	b.WriteString(`<div><input type="submit" value="上传并验证"></div>`)
	b.WriteString(`</form>`)

	// 全局错误（如文件解析失败）
	if errMsg != "" {
		b.WriteString(renderError(errMsg))
	}

	// 批量验证结果
	if uploaded && len(results) > 0 {
		b.WriteString(`<div class="content" style="max-width:760px;margin:20px auto;">`)
		b.WriteString(`<h2>批量验证结果（共 ` + strconv.Itoa(len(results)) + ` 行）</h2>`)

		for _, r := range results {
			b.WriteString(`<div style="margin-bottom:24px;padding:16px;background:rgba(0,0,0,.3);border-radius:8px;border:1px solid rgba(102,126,234,.3);">`)
			b.WriteString(`<div style="margin-bottom:8px;"><strong>第 ` + strconv.Itoa(r.RowNum) + ` 行</strong>　关键字：<code>` + html.EscapeString(r.Keyword) + `</code></div>`)
			b.WriteString(`<div style="margin-bottom:8px;"><strong>执行的 SQL：</strong></div>`)
			b.WriteString(`<pre style="background:rgba(0,0,0,.4);padding:10px;border-radius:6px;overflow:auto;font-size:13px;">` + html.EscapeString(r.RawSQL) + `</pre>`)

			if r.Err != "" {
				b.WriteString(`<div class="error" style="margin:8px 0;max-width:none;">` + html.EscapeString(r.Err) + `</div>`)
			} else {
				b.WriteString(`<div style="margin-bottom:4px;color:#4ecdc4;">查询结果（` + strconv.Itoa(len(r.Rows)) + ` 行）：</div>`)
				b.WriteString(renderExportTable(r.Cols, r.Rows))
			}
			b.WriteString(`</div>`)
		}
		b.WriteString(`</div>`)
	} else if uploaded {
		b.WriteString(`<div class="content" style="max-width:760px;margin:20px auto;">`)
		b.WriteString(`<p style="text-align:center;color:#888;">（文件中没有有效数据行）</p>`)
		b.WriteString(`</div>`)
	}

	b.WriteString(`<div class="nav"><a href="/">返回首页</a></div>`)
	b.WriteString(`</div></body></html>`)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, b.String())
}

// ---------------------------------------------------------------------------
// 文件上传处理：通用逻辑
// ---------------------------------------------------------------------------

// handleUploadFile 处理文件上传的通用逻辑。
// parseFunc 为对应格式的解析函数，接收文件内容的 []byte。
func handleUploadFile(c *gin.Context, l *BaseLesson, lessonID, title, desc, format, acceptExt string,
	parseFunc func(data []byte) ([][]string, error)) {

	l.logRequest(c, lessonID)

	file, err := c.FormFile("file")
	if err != nil {
		renderUploadPage(c, lessonID, title, desc, format, acceptExt, nil, "请选择要上传的文件", false)
		return
	}

	// 限制文件大小（10MB）
	if file.Size > 10<<20 {
		renderUploadPage(c, lessonID, title, desc, format, acceptExt, nil, "文件大小不能超过 10MB", false)
		return
	}

	src, err := file.Open()
	if err != nil {
		renderUploadPage(c, lessonID, title, desc, format, acceptExt, nil, "打开上传文件失败: "+err.Error(), false)
		return
	}
	defer src.Close()

	// 读取到内存
	data, err := io.ReadAll(src)
	if err != nil {
		renderUploadPage(c, lessonID, title, desc, format, acceptExt, nil, "读取上传文件失败: "+err.Error(), false)
		return
	}

	rows, err := parseFunc(data)
	if err != nil {
		renderUploadPage(c, lessonID, title, desc, format, acceptExt, nil, err.Error(), false)
		return
	}

	results := processUploadData(l.db, rows)
	renderUploadPage(c, lessonID, title, desc, format, acceptExt, results, "", true)
}
