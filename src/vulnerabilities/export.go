package vulnerabilities

import (
	"archive/zip"
	"bytes"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-sqli-lab/src/db"

	"github.com/gin-gonic/gin"
)

// 导出文件的 MIME 类型
const (
	contentTypeXLSX = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	contentTypeXLS  = "application/vnd.ms-excel"
)

// runExportQuery 执行「文件导出」场景下的查询。
// 漏洞点：keyword 被直接拼接进 SQL（与 Less-1 同款写法），可被 SQL 注入。
// 返回列名、字符串化的二维结果、原始 SQL（供页面展示）以及执行错误。
func runExportQuery(database db.Database, keyword string) (cols []string, rows [][]string, rawSQL string, err error) {
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

// valueToString 把数据库返回的任意类型统一转成字符串，便于写入表格文件。
func valueToString(v interface{}) string {
	switch val := v.(type) {
	case nil:
		return ""
	case []byte:
		return string(val)
	case string:
		return val
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'g', -1, 64)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case time.Time:
		return val.Format("2006-01-02 15:04:05")
	default:
		return fmt.Sprintf("%v", val)
	}
}

// serveExport 以附件形式下载导出的表格文件。
func serveExport(c *gin.Context, filename, contentType string, data []byte) {
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Data(http.StatusOK, contentType, data)
}

// serveExportError 查询失败时返回纯文本错误信息，方便观察报错注入。
func serveExportError(c *gin.Context, rawSQL string, err error) {
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, "导出失败，SQL 报错：\n%v\n\n执行的 SQL：\n%s\n", err, rawSQL)
}

// ---------------------------------------------------------------------------
// XLSX：使用标准库 archive/zip 手写最小可用的 OOXML 工作簿（真实 .xlsx）
// ---------------------------------------------------------------------------

func buildXLSX(cols []string, rows [][]string) ([]byte, error) {
	parts := []struct{ name, body string }{
		{"[Content_Types].xml", xlsxContentTypes},
		{"_rels/.rels", xlsxRootRels},
		{"xl/workbook.xml", xlsxWorkbook},
		{"xl/_rels/workbook.xml.rels", xlsxWorkbookRels},
		{"xl/worksheets/sheet1.xml", buildSheetXML(cols, rows)},
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, p := range parts {
		w, err := zw.Create(p.name)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write([]byte(p.body)); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func buildSheetXML(cols []string, rows [][]string) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	b.WriteString(`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`)

	rowNum := 1
	if len(cols) > 0 {
		b.WriteString(fmt.Sprintf(`<row r="%d">`, rowNum))
		for i, col := range cols {
			ref := fmt.Sprintf("%s%d", columnName(i), rowNum)
			b.WriteString(fmt.Sprintf(`<c r="%s" t="inlineStr"><is><t xml:space="preserve">%s</t></is></c>`, ref, escapeXML(col)))
		}
		b.WriteString(`</row>`)
		rowNum++
	}

	for _, row := range rows {
		b.WriteString(fmt.Sprintf(`<row r="%d">`, rowNum))
		for i, cell := range row {
			ref := fmt.Sprintf("%s%d", columnName(i), rowNum)
			b.WriteString(fmt.Sprintf(`<c r="%s" t="inlineStr"><is><t xml:space="preserve">%s</t></is></c>`, ref, escapeXML(cell)))
		}
		b.WriteString(`</row>`)
		rowNum++
	}

	b.WriteString(`</sheetData></worksheet>`)
	return b.String()
}

const xlsxContentTypes = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
	`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
	`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
	`<Default Extension="xml" ContentType="application/xml"/>` +
	`<Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>` +
	`<Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>` +
	`</Types>`

const xlsxRootRels = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
	`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
	`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>` +
	`</Relationships>`

const xlsxWorkbook = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
	`<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" ` +
	`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">` +
	`<sheets><sheet name="Sheet1" sheetId="1" r:id="rId1"/></sheets></workbook>`

const xlsxWorkbookRels = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
	`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
	`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/>` +
	`</Relationships>`

// columnName 把 0 基列索引转成 Excel 列名（0->A, 1->B, 26->AA）。
func columnName(idx int) string {
	name := ""
	for n := idx + 1; n > 0; {
		n--
		name = string(rune('A'+n%26)) + name
		n /= 26
	}
	return name
}

// escapeXML 转义 XML 文本内容，并剔除 XML 1.0 不允许的控制字符。
func escapeXML(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&apos;")
		default:
			if r < 0x20 && r != '\t' && r != '\n' && r != '\r' {
				continue
			}
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ---------------------------------------------------------------------------
// XLS：Excel 兼容的 HTML 表格（带 office/excel 命名空间），以 .xls 扩展名下载
// ---------------------------------------------------------------------------

func buildXLS(cols []string, rows [][]string) []byte {
	var b strings.Builder
	b.WriteString(`<html xmlns:o="urn:schemas-microsoft-com:office:office" `)
	b.WriteString(`xmlns:x="urn:schemas-microsoft-com:office:excel" `)
	b.WriteString(`xmlns="http://www.w3.org/TR/REC-html40">`)
	b.WriteString(`<head><meta charset="UTF-8"><!--[if gte mso 9]><xml><x:ExcelWorkbook>` +
		`<x:ExcelWorksheets><x:ExcelWorksheet><x:Name>Sheet1</x:Name>` +
		`<x:WorksheetOptions><x:DisplayGridlines/></x:WorksheetOptions>` +
		`</x:ExcelWorksheet></x:ExcelWorksheets></x:ExcelWorkbook></xml><![endif]--></head>`)
	b.WriteString(`<body><table border="1">`)

	if len(cols) > 0 {
		b.WriteString(`<tr>`)
		for _, col := range cols {
			b.WriteString(`<th>` + html.EscapeString(col) + `</th>`)
		}
		b.WriteString(`</tr>`)
	}

	for _, row := range rows {
		b.WriteString(`<tr>`)
		for _, cell := range row {
			b.WriteString(`<td>` + html.EscapeString(cell) + `</td>`)
		}
		b.WriteString(`</tr>`)
	}

	b.WriteString(`</table></body></html>`)
	return []byte(b.String())
}

// ---------------------------------------------------------------------------
// 入口页面：Less-66 / Less-67 共用（服务端内联渲染，复用 /static/css/style.css）
// ---------------------------------------------------------------------------

// renderExportPage 渲染导出关卡的入口页面。
// exportFormat 用于按钮文案（如 "XLS"/"XLSX"），exportExt 用于下载 URL 与文件名。
// searched 为 true 时展示预览查询的结果表格 / 报错 / 原始 SQL。
func renderExportPage(c *gin.Context, lessonID, title, desc, exportFormat, keyword string,
	cols []string, rows [][]string, rawSQL, errMsg string, searched bool) {

	exportURL := "/less/" + lessonID + "/export"

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
	b.WriteString(`<p>导出的查询语句直接拼接了 <code>keyword</code> 参数（未做任何过滤/参数化），存在 SQL 注入。</p>`)
	b.WriteString(`<pre style="background:rgba(0,0,0,.4);padding:12px;border-radius:6px;overflow:auto;">` +
		`SELECT id, username, password FROM users WHERE username LIKE '%[keyword]%'</pre>`)
	b.WriteString(`<h2>测试示例</h2><ul style="padding-left:20px;line-height:1.9;">`)
	b.WriteString(`<li>报错注入：<code>'</code></li>`)
	b.WriteString(`<li>UNION 注入（SQLite/MySQL）：<code>%' UNION SELECT 1,group_concat(username),3 FROM users -- </code></li>`)
	b.WriteString(`<li>UNION 注入（PostgreSQL）：<code>%' UNION SELECT 1,string_agg(username,','),3 FROM users -- </code></li>`)
	b.WriteString(`<li>布尔盲注：<code>admin' AND 1=1 -- </code> 对比 <code>admin' AND 1=2 -- </code></li>`)
	b.WriteString(`</ul></div>`)

	// 查询表单（预览，GET 提交到本页）
	b.WriteString(`<form method="GET" action="">`)
	b.WriteString(`<div><label>查询关键字 (keyword)：</label>`)
	b.WriteString(`<input type="text" id="kw" name="keyword" value="` + html.EscapeString(keyword) + `"></div>`)
	b.WriteString(`<div><input type="submit" value="预览查询结果"></div>`)
	b.WriteString(`</form>`)

	// 导出按钮
	b.WriteString(`<div class="nav"><a href="#" onclick="doExport();return false;">📥 导出 ` +
		html.EscapeString(exportFormat) + ` 文件</a></div>`)
	b.WriteString(`<script>function doExport(){var kw=document.getElementById('kw').value;` +
		`window.location='` + exportURL + `?keyword='+encodeURIComponent(kw);}</script>`)

	// 预览结果
	if searched {
		b.WriteString(`<div class="content" style="max-width:760px;margin:20px auto;">`)
		b.WriteString(`<h2>执行的 SQL</h2>`)
		b.WriteString(`<pre style="background:rgba(0,0,0,.4);padding:12px;border-radius:6px;overflow:auto;">` +
			html.EscapeString(rawSQL) + `</pre>`)
		if errMsg != "" {
			b.WriteString(renderError(errMsg))
		} else {
			b.WriteString(`<h2>查询结果（` + strconv.Itoa(len(rows)) + ` 行）</h2>`)
			b.WriteString(renderExportTable(cols, rows))
		}
		b.WriteString(`</div>`)
	}

	b.WriteString(`<div class="nav"><a href="/">返回首页</a></div>`)
	b.WriteString(`</div></body></html>`)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, b.String())
}

// renderExportTable 把查询结果渲染成 HTML 表格（复用 style.css 中的 table 样式）。
func renderExportTable(cols []string, rows [][]string) string {
	var b strings.Builder
	b.WriteString(`<table><thead><tr>`)
	for _, col := range cols {
		b.WriteString(`<th>` + html.EscapeString(col) + `</th>`)
	}
	b.WriteString(`</tr></thead><tbody>`)
	if len(rows) == 0 {
		b.WriteString(`<tr><td colspan="` + strconv.Itoa(len(cols)) + `" style="text-align:center;">（无数据）</td></tr>`)
	}
	for _, row := range rows {
		b.WriteString(`<tr>`)
		for _, cell := range row {
			b.WriteString(`<td>` + html.EscapeString(cell) + `</td>`)
		}
		b.WriteString(`</tr>`)
	}
	b.WriteString(`</tbody></table>`)
	return b.String()
}
