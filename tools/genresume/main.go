// Command genresume 一次性生成「应聘人员求职登记表」Excel 模板文件。
//
// 模板设计思路（模板文件即“样式/布局的单一来源”）：
//   - 所有文字、字段标签、配色、边框、合并单元格、行高列宽、照片区
//     全部固化在这个 .xlsx 模板文件里，可用 Excel/WPS 直接打开微调，运行时只负责填数据。
//   - 「值」单元格预先写入 RESUME:<字段名> 占位标记；程序加载模板后，按查询返回的
//     列名找到匹配字段并替换占位符为真实数据（字段错位问题彻底消除）。
//   - 模板即一张完整的求职登记表，导出文件就是模拟的简历表本身，不再追加任何明细表。
//
// 用法： go run ./tools/genresume   （输出到 src/vulnerabilities/assets/resume_template.xlsx）
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

func main() {
	if err := generate(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("模板已生成: src/vulnerabilities/assets/resume_template.xlsx")
}

func generate() error {
	f := excelize.NewFile()
	const s = "求职登记表"
	if err := f.SetSheetName("Sheet1", s); err != nil {
		return err
	}

	// 配色 / 字体常量
	const (
		fontCN     = "宋体"
		titleFont  = "黑体"
		colorDark  = "1F4E78"
		colorText  = "333333"
		labelFill  = "4472C4"
		valueFill  = "FFFFFF"
		titleFill  = "DDEBF7"
		gridColor  = "8EA9DB"
		photoFill  = "F2F2F2"
	)

	var (
		stTitle, stMeta, stLabel, stValue, stArea, stPhoto, stCell                              int
		err                                                                    error
	)

	border := func(color string) []excelize.Border {
		return []excelize.Border{
			{Type: "left", Color: color, Style: 1},
			{Type: "right", Color: color, Style: 1},
			{Type: "top", Color: color, Style: 1},
			{Type: "bottom", Color: color, Style: 1},
		}
	}
	borderDash := func(color string) []excelize.Border {
		return []excelize.Border{
			{Type: "left", Color: color, Style: 4},
			{Type: "right", Color: color, Style: 4},
			{Type: "top", Color: color, Style: 4},
			{Type: "bottom", Color: color, Style: 4},
		}
	}

	// 大标题
	if stTitle, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 20, Color: colorDark, Family: titleFont},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{titleFill}, Pattern: 1},
	}); err != nil {
		return err
	}
	// 元信息
	if stMeta, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Color: "595959", Family: fontCN},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
		Border:    border(colorDark),
	}); err != nil {
		return err
	}
	// 字段标签
	if stLabel, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF", Family: fontCN},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{labelFill}, Pattern: 1},
		Border:    border(gridColor),
	}); err != nil {
		return err
	}
	// 字段值
	if stValue, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 11, Color: colorText, Family: fontCN},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{valueFill}, Pattern: 1},
		Border:    border(gridColor),
	}); err != nil {
		return err
	}
	// 长文本区
	if stArea, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 11, Color: colorText, Family: fontCN},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{valueFill}, Pattern: 1},
		Border:    border(gridColor),
	}); err != nil {
		return err
	}
	// 照片区
	if stPhoto, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 12, Color: "7F7F7F", Family: fontCN},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{photoFill}, Pattern: 1},
		Border:    borderDash("BFBFBF"),
	}); err != nil {
		return err
	}
	// 普通格子（未用的残留样式单元格兜底）
	if stCell, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Color: colorText, Family: fontCN},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{valueFill}, Pattern: 1},
		Border:    border("BFBFBF"),
	}); err != nil {
		return err
	}
	_ = stCell

	// 列宽：A 照片位，B..I 字段网格（9 列）
	widths := map[string]float64{
		"A": 12, "B": 12, "C": 16, "D": 12, "E": 16, "F": 12, "G": 16, "H": 12, "I": 16,
	}
	for c, w := range widths {
		if err := f.SetColWidth(s, c, c, w); err != nil {
			return err
		}
	}

	// 第1行：大标题（跨 A:I）
	_ = f.SetCellValue(s, "A1", "应聘人员求职登记表")
	_ = f.MergeCell(s, "A1", "I1")
	_ = f.SetCellStyle(s, "A1", "I1", stTitle)
	_ = f.SetRowHeight(s, 1, 42)

	// 第2行：顶部说明（跨 A:I）
	_ = f.SetCellValue(s, "A2", "本表供应聘人员如实填写，请确保所填信息真实有效。　　填表日期：　　　　年　　月　　日")
	_ = f.MergeCell(s, "A2", "I2")
	_ = f.SetCellStyle(s, "A2", "I2", stMeta)
	_ = f.SetRowHeight(s, 2, 22)

	// 照片区：A3:A10 合并（跨普通字段区与长文本区左列）
	_ = f.SetCellValue(s, "A3", "照\n片")
	_ = f.MergeCell(s, "A3", "A10")
	_ = f.SetCellStyle(s, "A3", "A10", stPhoto)

	// 普通字段 16 个：每行 4 对「标签｜值」，共 4 行
	// 标签列 B/D/F/H，值列 C/E/G/I
	type pair struct{ label, field string }
	pairs := []pair{
		{"姓　　名", "name"}, {"性　　别", "sex"}, {"出生年月", "birth"}, {"民　　族", "nation"},
		{"籍　　贯", "native_place"}, {"政治面貌", "political"}, {"健康状况", "health"}, {"婚姻状况", "marital"},
		{"学　　历", "education"}, {"毕业院校", "school"}, {"所学专业", "major"}, {"毕业时间", "graduation"},
		{"联系电话", "phone"}, {"电子邮箱", "email"}, {"邮政编码", "postal_code"}, {"工作年限", "work_years"},
	}
	// 字段标签列与值列（B=2,C=3, D=4,E=5, F=6,G=7, H=8,I=9）
	labelCols := []int{2, 4, 6, 8}
	valueCols := []int{3, 5, 7, 9}

	apply := func(cell, txt string, styleIdx int) {
		_ = f.SetCellValue(s, cell, txt)
		_ = f.SetCellStyle(s, cell, cell, styleIdx)
	}
	applyValue := func(cell, field string) {
		_ = f.SetCellValue(s, cell, "RESUME:"+field)
		_ = f.SetCellStyle(s, cell, cell, stValue)
	}

	startRow := 3
	for i := 0; i < len(pairs); i++ {
		row := startRow + i/4
		lc := labelCols[i%4]
		vc := valueCols[i%4]
		apply(toCell(lc, row), pairs[i].label, stLabel)
		applyValue(toCell(vc, row), pairs[i].field)
		_ = f.SetRowHeight(s, row, 26)
	}

	// 长文本区（行 7-10）：求职意向 / 期望薪资 / 通讯地址 / 自我评价
	// 行7：求职意向（B7 标签，C7:I7 值）
	apply("B7", "求职意向", stLabel)
	_ = f.SetCellValue(s, "C7", "RESUME:job_intention")
	_ = f.MergeCell(s, "C7", "I7")
	_ = f.SetCellStyle(s, "C7", "I7", stArea)
	_ = f.SetRowHeight(s, 7, 26)

	// 行8：期望薪资（B8 标签，C8:I8 值）
	apply("B8", "期望薪资", stLabel)
	_ = f.SetCellValue(s, "C8", "RESUME:salary_expect")
	_ = f.MergeCell(s, "C8", "I8")
	_ = f.SetCellStyle(s, "C8", "I8", stArea)
	_ = f.SetRowHeight(s, 8, 26)

	// 行9：通讯地址（B9 标签，C9:I9 值）
	apply("B9", "通讯地址", stLabel)
	_ = f.SetCellValue(s, "C9", "RESUME:address")
	_ = f.MergeCell(s, "C9", "I9")
	_ = f.SetCellStyle(s, "C9", "I9", stArea)
	_ = f.SetRowHeight(s, 9, 26)

	// 行10-11：自我评价（B10 标签，C10:I11 值，跨两行）
	apply("B10", "自我评价", stLabel)
	_ = f.SetCellValue(s, "C10", "RESUME:self_eval")
	_ = f.MergeCell(s, "C10", "I11")
	_ = f.SetCellStyle(s, "C10", "I11", stArea)
	_ = f.SetRowHeight(s, 10, 44)
	_ = f.SetRowHeight(s, 11, 44)

	// 第12行：底部签名
	_ = f.SetCellValue(s, "A12", "应聘人签名：　　　　　　　　　　　　　　　　　　　　　　　审核意见：")
	_ = f.MergeCell(s, "A12", "I12")
	_ = f.SetCellStyle(s, "A12", "I12", stMeta)
	_ = f.SetRowHeight(s, 12, 26)

	// 冻结窗格（第3行起滚动）
	_ = f.SetPanes(s, &excelize.Panes{
		Freeze:      true,
		YSplit:      2,
		TopLeftCell: "A3",
		ActivePane:  "bottomLeft",
	})

	// 横向 A4 打印
	landscape := "landscape"
	paperSize := 9
	_ = f.SetPageLayout(s, &excelize.PageLayoutOptions{Orientation: &landscape, Size: &paperSize})

	// 删除多余 Sheet
	_ = f.DeleteSheet("Sheet2")
	_ = f.DeleteSheet("Sheet3")

	if err := os.MkdirAll("src/vulnerabilities/assets", 0755); err != nil {
		return err
	}
	return f.SaveAs("src/vulnerabilities/assets/resume_template.xlsx")
}

// toCell 将列索引(int, 从1起)与行号(int, 从1起)转换为 A1 形式的单元格引用。
func toCell(col, row int) string {
	c, _ := excelize.CoordinatesToCellName(col, row)
	return c
}
