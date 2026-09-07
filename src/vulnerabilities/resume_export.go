package vulnerabilities

import (
	"bytes"
	_ "embed"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"

	"go-sqli-lab/src/db"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// 嵌入预排版的求职登记表模板（布局、样式、字段标签、照片位、明细表头全部在模板文件里）。
// 运行时通过 excelize.OpenReader 加载模板，再按查询列名填充值，实现“模板驱动”而非代码硬编码布局。
//
//go:embed assets/resume_template.xlsx
var resumeTemplateXLSX []byte

// ---------------------------------------------------------------------------
// 简历导出注入：查询（漏洞点：keyword 直接拼接进 SQL，可注入）
// ---------------------------------------------------------------------------

// runResumeExportQuery 执行「简历导出」场景下的查询。
// 查询对象为 resumes（应聘人员求职登记表），比 users 表包含更复杂的结构化字段，
// 符合“复杂布局”的导出场景。keyword 直接拼接进 LIKE 子句，可被 SQL 注入。
func runResumeExportQuery(database db.Database, keyword string) (cols []string, rows [][]string, rawSQL string, err error) {
	rawSQL = fmt.Sprintf("SELECT name, sex, birth, nation, native_place, political, health, marital, "+
		"education, school, major, graduation, phone, email, postal_code, "+
		"work_years, job_intention, salary_expect, address, self_eval "+
		"FROM resumes WHERE name LIKE '%%%s%%'", keyword)

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

// ---------------------------------------------------------------------------
// 生成 Excel：.xlsx（加载模板填充）
// ---------------------------------------------------------------------------

// buildResumeXLSX 加载嵌入模板，把查询结果（第一行）填充到“应聘人员登记表”字段区。
// 模板即完整的求职登记表布局，导出文件就是一份模拟的简历表，不再额外追加明细表。
func buildResumeXLSX(cols []string, rows [][]string) ([]byte, error) {
	f, err := excelize.OpenReader(bytes.NewReader(resumeTemplateXLSX))
	if err != nil {
		return nil, fmt.Errorf("加载简历模板失败: %w", err)
	}
	defer f.Close()

	sheet := f.GetSheetName(0)

	// 填充“应聘人员信息”字段区：用查询列名匹配模板值单元格，
	// 标签与值来自同一列，彻底避免字段错位。
	fillByPlaceholder(f, sheet, cols, rows)

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// fillByPlaceholder 把查询结果第 0 行数据按「字段名 -> 模板单元格」映射填入模板。
// 模板布局与标签预先在 assets/resume_template.xlsx 中定义好，这里只负责“按名字喂数据”，
// 字段标签与值天然对位，彻底避免错位。
func fillByPlaceholder(f *excelize.File, sheet string, cols []string, rows [][]string) {
	if len(rows) == 0 {
		return
	}
	// 建立“列名 -> 列索引”映射（大小写不敏感），兼容注入后列顺序变化
	colIdx := map[string]int{}
	for i, c := range cols {
		colIdx[strings.ToLower(strings.TrimSpace(c))] = i
	}
	// 字段 -> 模板值单元格（与 tools/genresume 模板生成逻辑保持一致）
	for field, coord := range resumeFieldCellMap() {
		idx, ok := colIdx[strings.ToLower(strings.TrimSpace(field))]
		if !ok {
			continue
		}
		_ = f.SetCellValue(sheet, coord, getCellRaw(rows, 0, idx))
	}
}

// resumeFieldCellMap 返回“字段名 -> 模板值单元格坐标”映射。
// 该坐标与模板生成工具 tools/genresume 中的排版严格一致，如需调整排版请同步修改。
func resumeFieldCellMap() map[string]string {
	return map[string]string{
		"name":          "C3",
		"sex":           "E3",
		"birth":         "G3",
		"nation":        "I3",
		"native_place":  "C4",
		"political":     "E4",
		"health":        "G4",
		"marital":       "I4",
		"education":     "C5",
		"school":        "E5",
		"major":         "G5",
		"graduation":    "I5",
		"phone":         "C6",
		"email":         "E6",
		"postal_code":   "G6",
		"work_years":    "I6",
		"job_intention": "C7",
		"salary_expect": "C8",
		"address":       "C9",
		"self_eval":     "C10",
	}
}

// ---------------------------------------------------------------------------
// 生成 Excel：.xls（HTML 模板渲染，Excel 可直接打开）
// ---------------------------------------------------------------------------

// buildResumeXLS 生成 .xls —— 使用 HTML 模板（带 Office 命名空间）渲染求职登记表。
// 布局与样式通过内联 CSS 表达，数据用占位符 ${colName} 替换，字段与值同列对位。
func buildResumeXLS(cols []string, rows [][]string) ([]byte, error) {
	// 字段中文名 -> 模板显示的标签顺序（与数据库列顺序一致）
	fields := []string{"name", "sex", "birth", "nation", "native_place", "political", "health",
		"marital", "education", "school", "major", "graduation", "phone", "email", "postal_code",
		"work_years", "job_intention", "salary_expect", "address", "self_eval"}

	// 第一行作为“应聘人员信息”主数据
	md := map[string]string{}
	if len(rows) > 0 {
		for i, c := range cols {
			md[strings.ToLower(strings.TrimSpace(c))] = getCellRaw(rows, 0, i)
		}
	}

	var b strings.Builder
	b.WriteString(resumeXLSHead)
	b.WriteString(`<table class="resume">`)
	b.WriteString(`<tr><td class="title" colspan="8">应聘人员求职登记表</td></tr>`)
	b.WriteString(`<tr><td class="meta" colspan="8">本表供应聘人员如实填写，请确保所填信息真实有效。　　填表日期：　　　　年　　月　　日</td></tr>`)
	b.WriteString(`<tr><td class="photo" rowspan="6">照­片</td>`)

	writeXLSFieldTable(&b, fields, md)

	b.WriteString(`</table>`)
	b.WriteString(`</body></html>`)
	return []byte(b.String()), nil
}

// writeXLSFieldTable 用 HTML 生成字段网格（label | value 两列一组的 4列 x5行 布局）。
func writeXLSFieldTable(b *strings.Builder, fields []string, md map[string]string) {
	// 每行放 4 组字段：label1 value1 label2 value2 label3 value3 label4 value4
	// 表格列数=8（每字段占 label+value 两列）
	// 依据模板视觉习惯：照片位在 template 头照片, 字段区域为普通 8 列。
	labelMap := map[string]string{
		"name": "姓名", "sex": "性别", "birth": "出生年月", "nation": "民族",
		"native_place": "籍贯", "political": "政治面貌", "health": "健康状况", "marital": "婚姻状况",
		"education": "学历", "school": "毕业院校", "major": "所学专业", "graduation": "毕业时间",
		"phone": "联系电话", "email": "电子邮箱", "postal_code": "邮政编码", "work_years": "工作年限",
		"job_intention": "求职意向", "salary_expect": "期望薪资", "address": "通讯地址", "self_eval": "自我评价",
	}

	// 普通 16 字段（除最后 4 个长文本）：每行 4 对
	short := []string{"name", "sex", "birth", "nation", "native_place", "political", "health", "marital",
		"education", "school", "major", "graduation", "phone", "email", "postal_code", "work_years"}
	// 长文本字段：求职意向 / 期望薪资 / 通讯地址 / 自我评价
	long := []string{"job_intention", "salary_expect", "address", "self_eval"}

	for idx := 0; idx < len(short); idx++ {
		col := idx % 4
		if col == 0 {
			b.WriteString(`<tr>`)
		}
		fld := short[idx]
		label := labelMap[fld]
		val := html.EscapeString(md[fld])
		b.WriteString(`<td class="label">` + label + `</td><td class="value">` + val + `</td>`)
		if col == 3 {
			b.WriteString(`</tr>`)
		}
	}
	// 补齐最后一行未满 4 组的空单元格
	if len(short)%4 != 0 {
		remain := 4 - len(short)%4
		for i := 0; i < remain*2; i++ {
			b.WriteString(`<td class="label empty"></td><td class="value empty"></td>`)
		}
		b.WriteString(`</tr>`)
	}

	// 长文本字段：label 一列 + value 跨 7 列
	for _, fld := range long {
		label := labelMap[fld]
		val := html.EscapeString(md[fld])
		b.WriteString(`<tr><td class="label">` + label + `</td><td class="value" colspan="7">` + val + `</td></tr>`)
	}

	// 底部签名行
	b.WriteString(`<tr><td class="meta" colspan="8">应聘人签名：　　　　　　　　　　　　　　　　　　　　　　　审核意见：</td></tr>`)
}

// resumeXLSHead HTML 模板头部（含 Office Excel 命名空间与内联样式）。
const resumeXLSHead = `<!DOCTYPE html><html xmlns:o="urn:schemas-microsoft-com:office:office" ` +
	`xmlns:x="urn:schemas-microsoft-com:office:excel" xmlns="http://www.w3.org/TR/REC-html40">` +
	`<head><meta charset="UTF-8"><!--[if gte mso 9]><xml><x:ExcelWorkbook>` +
	`<x:ExcelWorksheets><x:ExcelWorksheet><x:Name>求职登记表</x:Name>` +
	`<x:WorksheetOptions><x:DisplayGridlines/></x:WorksheetOptions>` +
	`</x:ExcelWorksheet></x:ExcelWorksheets></x:ExcelWorkbook></xml><![endif]--></head>` +
	`<body><style>` +
	 `table.resume{border-collapse:collapse;width:100%;font-family:"宋体";font-size:14px;color:#333;}` +
	 `table.resume td{border:1px solid #8EA9DB;padding:6px 8px;vertical-align:middle;}` +
	 `td.title{background:#DDEBF7;color:#1F4E78;font-family:"黑体";font-size:26px;font-weight:bold;text-align:center;letter-spacing:8px;padding:14px 0;}` +
	 `td.meta{background:#F7F9FC;color:#595959;font-size:12px;}` +
	 `td.label{background:#4472C4;color:#fff;font-weight:bold;text-align:center;width:12%;white-space:nowrap;}` +
	 `td.value{background:#fff;width:13%;text-align:left;}` +
	 `td.value.empty{background:#fff;}` +
	 `td.photo{background:#F2F2F2;color:#7F7F7F;text-align:center;width:10%;vertical-align:middle;writing-mode:tb-rl;font-size:16px;}` +
	 `</style>`

// ---------------------------------------------------------------------------
// 简历预览查询结果——卡片式应聘人员求职登记表
// ---------------------------------------------------------------------------

// resumePreviewField 描述简历卡片上的一个字段（标签 + 数据库列名）。
type resumePreviewField struct {
	label string
	col   string
}

// 简历卡片短字段组：4 行 × 4 对「标签｜值」，顺序与 Excel 模板一致。
var resumePreviewShort = [][]resumePreviewField{
	{{"姓　　名", "name"}, {"性　　别", "sex"}, {"出生年月", "birth"}, {"民　　族", "nation"}},
	{{"籍　　贯", "native_place"}, {"政治面貌", "political"}, {"健康状况", "health"}, {"婚姻状况", "marital"}},
	{{"学　　历", "education"}, {"毕业院校", "school"}, {"所学专业", "major"}, {"毕业时间", "graduation"}},
	{{"联系电话", "phone"}, {"电子邮箱", "email"}, {"邮政编码", "postal_code"}, {"工作年限", "work_years"}},
}

// 简历卡片长文本字段：每行一个「标签｜值」，值与模板一致跨整行。
var resumePreviewLong = []resumePreviewField{
	{"求职意向", "job_intention"},
	{"期望薪资", "salary_expect"},
	{"通讯地址", "address"},
	{"自我评价", "self_eval"},
}

// renderResumePreviewCards 把查询结果渲染成一张张「应聘人员求职登记表」卡片。
//   - 每条数据一张卡片，字段布局 / 配色完全模拟导出 Excel 模板
//   - 左侧照片位，右侧 4×4 短字段对 + 4 行长文本字段，底部签名行
//   - 多张卡片自上而下排列，形成滚动区域
func renderResumePreviewCards(cols []string, rows [][]string) string {
	// 查询列名 -> 列下标
	colIdx := make(map[string]int, len(cols))
	for i, c := range cols {
		colIdx[strings.ToLower(strings.TrimSpace(c))] = i
	}
	val := func(rowIdx int, f resumePreviewField) string {
		if idx, ok := colIdx[f.col]; ok {
			return getCellRaw(rows, rowIdx, idx)
		}
		return ""
	}

	if len(rows) == 0 {
		return `<div class="resume-card-empty">（未查询到匹配的简历记录）</div>`
	}

	var b strings.Builder
	for i := range rows {
		b.WriteString(`<div class="resume-card">`)
		// 标题 + 元信息
		b.WriteString(`<div class="rc-title">应聘人员求职登记表</div>`)
		b.WriteString(`<div class="rc-meta">本表供应聘人员如实填写，请确保所填信息真实有效。　　填表日期：　　　　年　　月　　日</div>`)

		b.WriteString(`<div class="rc-grid">`)
		// 照片位（占左侧所有行）
		b.WriteString(`<div class="rc-photo">照<br>片</div>`)

		// 短字段组：4 行 × 4 对
		for row := range resumePreviewShort {
			for col := range resumePreviewShort[row] {
				f := resumePreviewShort[row][col]
				b.WriteString(`<div class="rc-label">` + html.EscapeString(f.label) + `</div>`)
				b.WriteString(`<div class="rc-val">` + html.EscapeString(val(i, f)) + `</div>`)
			}
		}

		// 长文本字段：每行一个 label + value(跨整行)
		for _, f := range resumePreviewLong {
			b.WriteString(`<div class="rc-label rc-label-wide">` + html.EscapeString(f.label) + `</div>`)
			b.WriteString(`<div class="rc-val rc-val-wide">` + html.EscapeString(val(i, f)) + `</div>`)
		}

		b.WriteString(`</div>`) // rc-grid

		// 底部签名行
		b.WriteString(`<div class="rc-sign">应聘人签名：　　　　　　　　　　　　　　　　　　　　　　　审核意见：</div>`)
		b.WriteString(`</div>`) // resume-card
	}

	return b.String()
}

// ---------------------------------------------------------------------------
// 工具函数
// ---------------------------------------------------------------------------

// resumeFieldLabel 把查询列名(接口返回的英文列名)转换为可读的中文标签。
// 用于明细表表头与 .xls HTML 渲染，避免直接显示英文列名。
func resumeFieldLabel(col string) string {
	c := strings.ToLower(strings.TrimSpace(col))
	switch c {
	case "name":
		return "姓名"
	case "sex":
		return "性别"
	case "birth":
		return "出生年月"
	case "nation":
		return "民族"
	case "native_place":
		return "籍贯"
	case "political":
		return "政治面貌"
	case "health":
		return "健康状况"
	case "marital":
		return "婚姻状况"
	case "education":
		return "学历"
	case "school":
		return "毕业院校"
	case "major":
		return "所学专业"
	case "graduation":
		return "毕业时间"
	case "phone":
		return "联系电话"
	case "email":
		return "电子邮箱"
	case "postal_code":
		return "邮政编码"
	case "work_years":
		return "工作年限"
	case "job_intention":
		return "求职意向"
	case "salary_expect":
		return "期望薪资"
	case "address":
		return "通讯地址"
	case "self_eval":
		return "自我评价"
	default:
		return col
	}
}

// borderThin 生成 4 边细边框，供动态明细表数据行使用。
func borderThin(color string) []excelize.Border {
	return []excelize.Border{
		{Type: "left", Color: color, Style: 1},
		{Type: "right", Color: color, Style: 1},
		{Type: "top", Color: color, Style: 1},
		{Type: "bottom", Color: color, Style: 1},
	}
}

// toAxis 列索引(从1起)+行号(从1起) -> A1 单元格引用。
func toAxis(col, row int) string {
	c, _ := excelize.CoordinatesToCellName(col, row)
	return c
}

// getCellRaw 安全地从二维切片取值，越界返回空字符串。
func getCellRaw(rows [][]string, rowIdx, colIdx int) string {
	if rowIdx < 0 || rowIdx >= len(rows) {
		return ""
	}
	if colIdx < 0 || colIdx >= len(rows[rowIdx]) {
		return ""
	}
	return rows[rowIdx][colIdx]
}

// ---------------------------------------------------------------------------
// 简历导出页面渲染
// ---------------------------------------------------------------------------

// renderResumeExportPage 渲染简历导出关卡的入口页面。
func renderResumeExportPage(c *gin.Context, lessonID, title, desc, exportFormat, keyword string,
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
		`SELECT name, sex, birth, ... FROM resumes WHERE name LIKE '%[keyword]%'</pre>`)
	b.WriteString(`<h2>测试示例</h2><ul style="padding-left:20px;line-height:1.9;">`)
	b.WriteString(`<li>报错注入：<code>'</code></li>`)
	b.WriteString(`<li>UNION 注入（SQLite/MySQL）：<code>%' UNION SELECT 1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20 FROM resumes -- </code></li>`)
	b.WriteString(`<li>布尔盲注：<code>张伟' AND 1=1 -- </code> 对比 <code>张伟' AND 1=2 -- </code></li>`)
	b.WriteString(`</ul></div>`)

	// 查询表单
	b.WriteString(`<form method="GET" action="">`)
	b.WriteString(`<div><label>查询姓名关键字 (keyword)：</label>`)
	b.WriteString(`<input type="text" id="kw" name="keyword" value="` + html.EscapeString(keyword) + `"></div>`)
	b.WriteString(`<div><input type="submit" value="预览查询结果"></div>`)
	b.WriteString(`</form>`)

	// 导出按钮
	b.WriteString(`<div class="nav"><a href="#" onclick="doExport();return false;">导出简历登记表 (` +
		html.EscapeString(exportFormat) + `)</a></div>`)
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
			b.WriteString(renderResumePreviewCards(cols, rows))
		}
		b.WriteString(`</div>`)
	}

	b.WriteString(`<div class="nav"><a href="/">返回首页</a></div>`)
	b.WriteString(`</div></body></html>`)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, b.String())
}
