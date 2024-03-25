package component

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	ut "github.com/nervatura/component/pkg/util"
)

const (
	TableEventCurrentPage  = "current_page"
	TableEventFilterChange = "filter_change"
	TableEventAddItem      = "add_item"
	TableEventEditCell     = "edit_cell"
	TableEventRowSelected  = "row_selected"

	TableFieldTypeString   = "string"
	TableFieldTypeNumber   = "number"
	TableFieldTypeDateTime = "datetime"
	TableFieldTypeDate     = "date"
	TableFieldTypeTime     = "time"
	TableFieldTypeBool     = "bool"
	TableFieldTypeDeffield = "deffield"
	TableFieldTypeCustom   = "custom"
)

var TableFieldType []string = []string{TableFieldTypeString, TableFieldTypeNumber, TableFieldTypeDateTime,
	TableFieldTypeDate, TableFieldTypeTime, TableFieldTypeBool, TableFieldTypeDeffield, TableFieldTypeCustom}

type Table struct {
	BaseComponent
	RowKey            string       `json:"row_key"`
	Rows              []ut.IM      `json:"rows"`
	Fields            []TableField `json:"fields"`
	Pagination        string       `json:"pagination"`
	CurrentPage       int64        `json:"current_page"`
	PageSize          int64        `json:"page_size"`
	HidePaginatonSize bool         `json:"hide_paginaton_size"`
	TableFilter       bool         `json:"table_filter"`
	AddItem           bool         `json:"add_item"`
	FilterPlaceholder string       `json:"filter_placeholder"`
	FilterValue       string       `json:"filter_value"`
	LabelYes          string       `json:"label_yes"`
	LabelNo           string       `json:"label_no"`
	LabelAdd          string       `json:"label_add"`
	AddIcon           string       `json:"add_icon"`
	TablePadding      string       `json:"table_padding"`
	SortCol           string       `json:"sort_col"`
	SortAsc           bool         `json:"sort_asc"`
	RowSelected       bool         `json:"row_selected"`
}

type TableField struct {
	Name          string       `json:"name"`
	FieldType     string       `json:"field_type"`
	Label         string       `json:"label"`
	TextAlign     string       `json:"text_align"`
	VerticalAlign string       `json:"vertical_align"`
	Format        bool         `json:"format"`
	Column        *TableColumn `json:"column"`
}

type TableColumn struct {
	Id          string                                                     `json:"id"`
	Header      string                                                     `json:"header"`
	HeaderStyle ut.SM                                                      `json:"header_style"`
	CellStyle   ut.SM                                                      `json:"cell_style"`
	Field       TableField                                                 `json:"field"`
	Cell        func(row ut.IM, col TableColumn, value interface{}) string `json:"-"`
}

func (tbl *Table) Properties() ut.IM {
	return ut.MergeIM(
		tbl.BaseComponent.Properties(),
		ut.IM{
			"row_key":             tbl.RowKey,
			"rows":                tbl.Rows,
			"fields":              tbl.Fields,
			"pagination":          tbl.Pagination,
			"current_page":        tbl.CurrentPage,
			"page_size":           tbl.PageSize,
			"hide_paginaton_size": tbl.HidePaginatonSize,
			"table_filter":        tbl.TableFilter,
			"add_item":            tbl.AddItem,
			"filter_placeholder":  tbl.FilterPlaceholder,
			"filter_value":        tbl.FilterValue,
			"label_yes":           tbl.LabelYes,
			"label_no":            tbl.LabelNo,
			"label_add":           tbl.LabelAdd,
			"add_icon":            tbl.AddIcon,
			"table_padding":       tbl.TablePadding,
			"sort_col":            tbl.SortCol,
			"sort_asc":            tbl.SortAsc,
		})
}

func (tbl *Table) GetProperty(propName string) interface{} {
	return tbl.Properties()[propName]
}

func (tbl *Table) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"row_key": func() interface{} {
			return ut.ToString(propValue, "id")
		},
		"rows": func() interface{} {
			return ut.ToIMA(propValue, []ut.IM{})
		},
		"fields": func() interface{} {
			fields := []TableField{}
			if fd, valid := propValue.([]TableField); valid && (fd != nil) {
				fields = fd
			}
			if len(fields) == 0 {
				if len(tbl.Rows) > 0 {
					for field := range tbl.Rows[0] {
						fields = append(fields,
							TableField{Name: field, FieldType: TableFieldTypeString, Label: field})
					}
				}
			}
			return fields
		},
		"pagination": func() interface{} {
			return tbl.CheckEnumValue(tbl.Pagination, PaginationTypeTop, PaginationType)
		},
		"current_page": func() interface{} {
			value := ut.ToInteger(propValue, 1)
			rows := tbl.filterRows()
			pageCount := int64(math.Ceil(float64(len(rows)) / float64(tbl.PageSize)))
			if value > pageCount {
				value = pageCount
			}
			if value < 1 {
				value = 1
			}
			return value
		},
		"page_size": func() interface{} {
			value := ut.ToInteger(propValue, 10)
			pageSize := []string{}
			for _, ps := range ValidPageSize {
				pageSize = append(pageSize, ut.ToString(ps, ""))
			}
			if !ut.Contains(pageSize, ut.ToString(value, "")) {
				value = ValidPageSize[0]
			}
			return value
		},
		"target": func() interface{} {
			tbl.SetProperty("id", tbl.Id)
			value := ut.ToString(propValue, tbl.Id)
			if value != "this" && !strings.HasPrefix(value, "#") {
				value = "#" + value
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if tbl.BaseComponent.GetProperty(propName) != nil {
		return tbl.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

func (tbl *Table) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"row_key": func() interface{} {
			tbl.RowKey = tbl.Validation(propName, propValue).(string)
			return tbl.RowKey
		},
		"rows": func() interface{} {
			tbl.Rows = tbl.Validation(propName, propValue).([]ut.IM)
			return tbl.Rows
		},
		"fields": func() interface{} {
			tbl.Fields = tbl.Validation(propName, propValue).([]TableField)
			return tbl.Fields
		},
		"pagination": func() interface{} {
			tbl.Pagination = tbl.Validation(propName, propValue).(string)
			return tbl.Pagination
		},
		"current_page": func() interface{} {
			tbl.CurrentPage = tbl.Validation(propName, propValue).(int64)
			return tbl.CurrentPage
		},
		"page_size": func() interface{} {
			tbl.PageSize = tbl.Validation(propName, propValue).(int64)
			return tbl.PageSize
		},
		"hide_paginaton_size": func() interface{} {
			tbl.HidePaginatonSize = ut.ToBoolean(propValue, false)
			return tbl.HidePaginatonSize
		},
		"table_filter": func() interface{} {
			tbl.TableFilter = ut.ToBoolean(propValue, false)
			return tbl.TableFilter
		},
		"add_item": func() interface{} {
			tbl.AddItem = ut.ToBoolean(propValue, false)
			return tbl.AddItem
		},
		"filter_placeholder": func() interface{} {
			tbl.FilterPlaceholder = ut.ToString(propValue, "")
			return tbl.FilterPlaceholder
		},
		"filter_value": func() interface{} {
			tbl.FilterValue = ut.ToString(propValue, "")
			return tbl.FilterValue
		},
		"label_yes": func() interface{} {
			tbl.LabelYes = ut.ToString(propValue, "YES")
			return tbl.LabelYes
		},
		"label_no": func() interface{} {
			tbl.LabelNo = ut.ToString(propValue, "NO")
			return tbl.LabelNo
		},
		"label_add": func() interface{} {
			tbl.LabelAdd = ut.ToString(propValue, "")
			return tbl.LabelAdd
		},
		"add_icon": func() interface{} {
			tbl.AddIcon = ut.ToString(propValue, "Plus")
			return tbl.AddIcon
		},
		"table_padding": func() interface{} {
			tbl.TablePadding = ut.ToString(propValue, "")
			return tbl.TablePadding
		},
		"sort_col": func() interface{} {
			tbl.SortCol = ut.ToString(propValue, "")
			return tbl.SortCol
		},
		"sort_asc": func() interface{} {
			tbl.SortAsc = ut.ToBoolean(propValue, false)
			return tbl.SortAsc
		},
		"target": func() interface{} {
			tbl.Target = tbl.Validation(propName, propValue).(string)
			return tbl.Target
		},
	}
	if _, found := pm[propName]; found {
		return tbl.SetRequestValue(propName, pm[propName](), []string{})
	}
	if tbl.BaseComponent.GetProperty(propName) != nil {
		return tbl.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (tbl *Table) SortRows(fieldName, fieldType string, sortAsc bool) {
	lessMap := map[string]func(i, j int) bool{
		TableFieldTypeNumber: func(i, j int) bool {
			a := ut.ToFloat(tbl.Rows[i][fieldName], 0)
			b := ut.ToFloat(tbl.Rows[j][fieldName], 0)
			if sortAsc {
				return a > b
			}
			return a < b
		},
		TableFieldTypeString: func(i, j int) bool {
			a := ut.ToString(tbl.Rows[i][fieldName], "")
			b := ut.ToString(tbl.Rows[j][fieldName], "")
			if sortAsc {
				return a > b
			}
			return a < b
		},
	}
	lessFn := lessMap[TableFieldTypeString]
	if fn, found := lessMap[fieldType]; found {
		lessFn = fn
	}
	sort.Slice(tbl.Rows, lessFn)
}

func (tbl *Table) response(evt ResponseEvent) (re ResponseEvent) {
	tblEvt := ResponseEvent{
		Trigger: tbl, TriggerName: tbl.Name, Value: evt.Value,
	}
	switch evt.TriggerName {
	case "top_pagination", "bottom_pagination":
		if evt.Name == PaginationEventPageSize {
			tbl.SetProperty("page_size", tblEvt.Value)
			tbl.SetProperty("current_page", 1)
			return tblEvt
		}
		tblEvt.Name = TableEventCurrentPage
		tbl.SetProperty("current_page", tblEvt.Value)

	case "filter":
		tblEvt.Name = TableEventFilterChange
		tbl.SetProperty("filter_value", tblEvt.Value)

	case "btn_add":
		tblEvt.Name = TableEventAddItem

	case "link_cell":
		tblEvt.Name = TableEventEditCell
		tblEvt.Value = evt.Trigger.GetProperty("data")

	case "data_row":
		tblEvt.Name = TableEventRowSelected
		tblEvt.Value = evt.Trigger.GetProperty("data")

	case "header_sort":
		sortCol := ut.ToString(evt.Trigger.GetProperty("data").(ut.IM)["fieldname"], "")
		fieldType := ut.ToString(evt.Trigger.GetProperty("data").(ut.IM)["fieldtype"], "")
		if tbl.SortCol == sortCol {
			tbl.SetProperty("sort_asc", !tbl.SortAsc)
		}
		tbl.SetProperty("sort_col", sortCol)
		tbl.SortRows(tbl.SortCol, fieldType, tbl.SortAsc)
		return tblEvt

	default:
	}
	if tbl.OnResponse != nil {
		return tbl.OnResponse(tblEvt)
	}
	return tblEvt
}

func (tbl *Table) getComponent(name string, pageCount int64, data ut.IM) (res string, err error) {
	ccPgn := func() *Pagination {
		return &Pagination{
			BaseComponent: BaseComponent{
				Id: tbl.Id + "_" + name, Name: name,
				EventURL:     tbl.EventURL,
				Target:       tbl.Target,
				OnResponse:   tbl.response,
				RequestValue: tbl.RequestValue,
				RequestMap:   tbl.RequestMap,
			},
			Value: tbl.CurrentPage, PageSize: tbl.PageSize,
			PageCount:    pageCount,
			HidePageSize: tbl.HidePaginatonSize,
		}
	}
	ccMap := map[string]func() ClientComponent{
		"top_pagination": func() ClientComponent {
			return ccPgn()
		},
		"bottom_pagination": func() ClientComponent {
			return ccPgn()
		},
		"filter": func() ClientComponent {
			return &Input{
				BaseComponent: BaseComponent{
					Id: tbl.Id + "_" + name, Name: name,
					Style:        ut.SM{"border-radius": "0", "margin": "1px 0 2px"},
					EventURL:     tbl.EventURL,
					Target:       tbl.Target,
					Swap:         SwapOuterHTML,
					OnResponse:   tbl.response,
					RequestValue: tbl.RequestValue,
					RequestMap:   tbl.RequestMap,
				},
				Type:        InputTypeText,
				Label:       tbl.FilterPlaceholder,
				Placeholder: tbl.FilterPlaceholder,
				Value:       tbl.FilterValue,
				Full:        true,
			}
		},
		"btn_add": func() ClientComponent {
			return &Button{
				BaseComponent: BaseComponent{
					Id: tbl.Id + "_" + name, Name: name,
					Style:        ut.SM{"padding": "8px 16px", "border-radius": "0", "margin": "1px 0 2px 1px"},
					EventURL:     tbl.EventURL,
					Target:       tbl.Target,
					OnResponse:   tbl.response,
					RequestValue: tbl.RequestValue,
					RequestMap:   tbl.RequestMap,
				},
				Type: ButtonTypeBorder,
				Icon: tbl.AddIcon, Label: tbl.LabelAdd,
			}
		},
		"link_cell": func() ClientComponent {
			return &Label{
				BaseComponent: BaseComponent{
					Id:           tbl.Id + "_" + ut.ToString(data["fieldname"], "") + "_" + ut.ToString(data["row"].(ut.IM)[tbl.RowKey], ""),
					Name:         name,
					EventURL:     tbl.EventURL,
					Target:       tbl.Target,
					Data:         data,
					OnResponse:   tbl.response,
					RequestValue: tbl.RequestValue,
					RequestMap:   tbl.RequestMap,
				},
				Value: ut.ToString(data["value"], ""),
			}
		},
	}
	cc := ccMap[name]()
	res, err = cc.Render()
	return res, err
}

func (tbl *Table) getStyle(styleMap ut.SM) string {
	style := []string{}
	for key, value := range styleMap {
		style = append(style, key+":"+value)
	}
	if len(style) > 0 {
		return fmt.Sprintf(` style="%s;"`, strings.Join(style, ";"))
	}
	return ""
}

func (tbl *Table) columns() (cols []TableColumn) {
	numberCell := func(value float64, label string, style ut.SM) string {
		return fmt.Sprintf(
			`<div class="number-cell">
	    <span class="cell-label">%s</span>
	    <span %s >%g</span>
    </div>`, label, tbl.getStyle(style), ut.ToFloat(value, 0))
	}

	dateCell := func(value interface{}, label, dateType string) string {
		var fmtValue string
		dateFormat := map[string]func(tm time.Time) string{
			TableFieldTypeDate: func(tm time.Time) string {
				return tm.Format("2006-01-02")
			},
			TableFieldTypeTime: func(tm time.Time) string {
				return tm.Format("15:04")
			},
			TableFieldTypeDateTime: func(tm time.Time) string {
				return tm.Format("2006-01-02 15:04")
			},
		}
		switch v := value.(type) {
		case string:
			tmValue, _ := ut.StringToDateTime(v)
			fmtValue = dateFormat[dateType](tmValue)
		case time.Time:
			fmtValue = dateFormat[dateType](v)
		}
		return fmt.Sprintf(`<span class="cell-label">%s</span><span>%s</span>`, label, fmtValue)
	}

	boolCell := func(value interface{}, label string) string {
		if (value == 1) || (value == "true") || (value == true) {
			return fmt.Sprintf(
				`<span class="cell-label">%s</span>
			  <form-icon iconKey="CheckSquare" ></form-icon>
			  <span class="middle"> %s</span>`, label, tbl.LabelYes)
		}
		return fmt.Sprintf(
			`<span class="cell-label">%s</span>
			<form-icon iconKey="SquareEmpty" ></form-icon>
			<span class="middle"> %s</span>`, label, tbl.LabelNo)
	}

	linkCell := func(value, label, fieldname string, resultValue interface{}, rowData ut.IM) string {
		linkLabel := fmt.Sprintf(
			`<span class="cell-label">%s</span>`, label)
		var link string
		link, _ = tbl.getComponent("link_cell", 0, ut.IM{
			"value": value, "fieldname": fieldname, "result": resultValue, "row": rowData,
		})
		return linkLabel + link
	}

	stringCell := func(value string, label string, style ut.SM) string {
		return fmt.Sprintf(
			`<span class="cell-label">%s</span>
			<span %s >%s</span>`, label, tbl.getStyle(style), value)
	}

	cols = []TableColumn{}
	for _, field := range tbl.Fields {
		if field.Column != nil {
			field.Column.Field = field
			cols = append(cols, *field.Column)
		} else {
			coldef := TableColumn{
				Id:          field.Name,
				Header:      ut.ToString(field.Label, field.Name),
				HeaderStyle: ut.SM{},
				CellStyle:   ut.SM{},
				Field:       field,
			}
			switch field.FieldType {

			case TableFieldTypeNumber:
				coldef.HeaderStyle["text-align"] = TextAlignRight
				coldef.Cell = func(row ut.IM, col TableColumn, value interface{}) string {
					style := ut.SM{}
					if col.Field.Format {
						style["font-weight"] = "bold"
						if evalue, found := row["edited"].(bool); found && evalue {
							style["text-decoration"] = "line-through"
						} else if ut.ToFloat(value, 0) != 0 {
							style["color"] = "red"
						} else {
							style["color"] = "green"
						}
					}
					return numberCell(ut.ToFloat(value, 0), col.Field.Label, style)
				}

			case TableFieldTypeDate, TableFieldTypeTime, TableFieldTypeDateTime:
				coldef.Cell = func(row ut.IM, col TableColumn, value interface{}) string {
					return dateCell(value, col.Field.Label, col.Field.FieldType)
				}

			case TableFieldTypeBool:
				coldef.Cell = func(row ut.IM, col TableColumn, value interface{}) string {
					return boolCell(value, col.Field.Label)
				}

			case TableFieldTypeDeffield:
				coldef.Cell = func(row ut.IM, col TableColumn, value interface{}) string {
					switch row["fieldtype"] {
					case "bool":
						return boolCell(value, col.Field.Label)

					case "integer", "float":
						return numberCell(ut.ToFloat(value, 0), col.Field.Label, ut.SM{})

					case "customer", "tool", "product", "trans", "transitem", "transmovement",
						"transpayment", "project", "employee", "place", "urlink":
						return linkCell(ut.ToString(row["export_deffield_value"], ""), col.Field.Label,
							ut.ToString(row["fieldtype"], ""), row[field.Name], row,
						)

					default:
						return stringCell(ut.ToString(value, ""), col.Field.Label, ut.SM{})
					}
				}

			default:
				coldef.Cell = func(row ut.IM, col TableColumn, value interface{}) string {
					style := ut.SM{}
					if color, found := row[col.Field.Name+"_color"].(string); found {
						style["color"] = color
					}
					for key, ivalue := range row {
						if key == "export_"+col.Field.Name {
							return linkCell(ut.ToString(ivalue, ""), col.Field.Label,
								col.Field.Name, row[col.Field.Name], row,
							)
						}
					}
					return stringCell(ut.ToString(value, ""), col.Field.Label, style)
				}
			}

			if tbl.TablePadding != "" {
				coldef.HeaderStyle["padding"] = tbl.TablePadding
				coldef.CellStyle["padding"] = tbl.TablePadding
			}
			if field.VerticalAlign != "" {
				coldef.CellStyle["vertical-align"] = tbl.CheckEnumValue(field.VerticalAlign, VerticalAlignMiddle, VerticalAlign)
			}
			if field.TextAlign != "" {
				coldef.CellStyle["text-align"] = tbl.CheckEnumValue(field.TextAlign, TextAlignLeft, TextAlign)
			}
			cols = append(cols, coldef)
		}
	}
	return cols
}

func (tbl *Table) filterRows() (rows []ut.IM) {
	rows = []ut.IM{}
	getValidRow := func(row ut.IM, filter string) bool {
		for field := range row {
			if strings.Contains(ut.ToString(row[field], ""), filter) {
				return true
			}
		}
		return false
	}
	if tbl.FilterValue == "" {
		return tbl.Rows
	}
	for _, row := range tbl.Rows {
		if getValidRow(row, tbl.FilterValue) {
			rows = append(rows, row)
		}
	}
	return rows
}

func (tbl *Table) Render() (res string, err error) {
	tbl.InitProps(tbl)

	cols := tbl.columns()
	rows := tbl.filterRows()
	pageCount := int64(math.Ceil(float64(len(rows)) / float64(tbl.PageSize)))

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(tbl.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(tbl.Class, " ")
		},
		"topPagination": func() bool {
			return ((pageCount > 1) && ((tbl.Pagination == PaginationTypeTop) || tbl.Pagination == PaginationTypeAll))
		},
		"bottomPagination": func() bool {
			return ((pageCount > 1) && ((tbl.Pagination == PaginationTypeBottom) || tbl.Pagination == PaginationTypeAll))
		},
		"tableComponent": func(name string) (string, error) {
			return tbl.getComponent(name, pageCount, ut.IM{})
		},
		"pageRows": func() []ut.IM {
			if tbl.Pagination != PaginationTypeNone {
				currentPage := tbl.Validation("current_page", tbl.CurrentPage).(int64)
				start := (currentPage - 1) * tbl.PageSize
				end := currentPage * tbl.PageSize
				if end > int64(len(rows)) {
					end = int64(len(rows))
				}
				return rows[start:end]
			}
			return rows
		},
		"colID": func(col TableColumn) string {
			colID := tbl.Id + "_header_" + col.Id
			lbl := &Label{BaseComponent: BaseComponent{
				Id: colID, Name: "header_sort",
				Data:         ut.IM{"fieldname": col.Id, "fieldtype": col.Field.FieldType},
				OnResponse:   tbl.response,
				RequestValue: tbl.RequestValue,
				RequestMap:   tbl.RequestMap,
			}}
			lbl.SetProperty("request_map", lbl)
			return colID
		},
		"rowID": func(row ut.IM, index int) string {
			rowID := ""
			if id, found := row[tbl.RowKey]; found {
				rowID = tbl.Id + "_row_" + ut.ToString(id, "")
			} else {
				rowID = tbl.Id + "_row_" + ut.ToString(index, "")
			}
			if tbl.RowSelected {
				lbl := &Label{BaseComponent: BaseComponent{
					Id: rowID, Name: "data_row", Data: ut.IM{
						"row": row, "index": index,
					},
					OnResponse:   tbl.response,
					RequestValue: tbl.RequestValue,
					RequestMap:   tbl.RequestMap,
				}}
				lbl.SetProperty("request_map", lbl)
			}
			return rowID
		},
		"pointerClass": func(row ut.IM) string {
			if disabled, found := row["disabled"].(bool); found && disabled {
				return "cursor-disabled"
			}
			if tbl.RowSelected {
				return "cursor-pointer"
			}
			return ""
		},
		"cols": func() []TableColumn {
			return cols
		},
		"sortClass": func(colID string) string {
			if tbl.SortCol == colID {
				if tbl.SortAsc {
					return "sort-asc"
				}
				return "sort-desc"
			}
			return "sort-none"
		},
		"cellStyle": func(styleMap ut.SM) bool {
			return len(styleMap) > 0
		},
		"cellValue": func(row ut.IM, col TableColumn) string {
			if col.Cell != nil {
				return col.Cell(row, col, row[col.Id])
			}
			return ut.ToString(row[col.Id], "")
		},
	}
	tpl := `<div id="{{ .Id }}" name="{{ .Name }}" class="responsive {{ customClass }}">
	{{ if or .TableFilter topPagination }}<div>
	{{ if topPagination }}<div>{{ tableComponent "top_pagination" }}</div>{{ end }}
	{{ if .TableFilter }}<div class="row full">
	<div class="cell" >{{ tableComponent "filter" }}</div>
	{{ if .AddItem }}<div class="cell" style="width: 20px;" >{{ tableComponent "btn_add" }}</div>{{ end }}
	</div>{{ end }}</div>{{ end }}
	<div class="table-wrap" ><table class="ui-table"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}>
	<thead><tr>{{ range $icol, $col := cols }}
	<th id="{{ colID $col }}" name="header_cell" 
	class="sort {{ sortClass $col.Id }}" 
	{{ if ne $.EventURL "" }} hx-post="{{ $.EventURL }}" hx-target="{{ $.Target }}" hx-swap="{{ $.Swap }}"{{ end }}
	{{ if ne $.Indicator "" }} hx-indicator="#{{ $.Indicator }}"{{ end }} 
	{{ if cellStyle $col.HeaderStyle }} style="{{ range $key, $value := $col.HeaderStyle }}{{ $key }}:{{ $value }};{{ end }}"{{ end }} 
	>{{ $col.Header }}</th>
	{{ end }}</tr></thead>
	<tbody>{{ range $index, $row := pageRows }}
	<tr id="{{ rowID $row $index }}" class="{{ pointerClass $row }}" 
	{{ if and ($.RowSelected) (ne $.EventURL "") }} hx-post="{{ $.EventURL }}" hx-target="{{ $.Target }}" hx-swap="{{ $.Swap }}"{{ end }}
	{{ if and ($.RowSelected) (ne $.Indicator "") }} hx-indicator="#{{ $.Indicator }}"{{ end }}
	>{{ range $icol, $col := cols }}<td
	{{ if cellStyle $col.CellStyle }} style="{{ range $key, $value := $col.CellStyle }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	>{{ cellValue $row $col }}</td>{{ end }}</tr>
	{{ end }}</tbody>
	</table></div>
	{{ if bottomPagination }}<div>{{ tableComponent "bottom_pagination" }}</div>{{ end }}
	</div>`

	return ut.TemplateBuilder("table", tpl, funcMap, tbl)
}

var testFields []TableField = []TableField{
	{Name: "name", FieldType: TableFieldTypeString, Label: "Name", TextAlign: TextAlignLeft},
	{Name: "valid", FieldType: TableFieldTypeBool, Label: "Valid"},
	{Name: "date", FieldType: TableFieldTypeDate, Label: "From"},
	{Name: "start", FieldType: TableFieldTypeTime},
	{Name: "stamp", FieldType: TableFieldTypeDateTime, Label: "Stamp"},
	{Name: "levels", FieldType: TableFieldTypeNumber, Label: "Levels", Format: true, VerticalAlign: VerticalAlignMiddle},
	{Name: "deffield", FieldType: TableFieldTypeDeffield, Label: "Deffield"},
	{Column: &TableColumn{Id: "editor", Cell: func(row ut.IM, col TableColumn, value interface{}) string {
		btn := Button{
			Type: ButtonTypePrimary, Label: "Hello", Disabled: ut.ToBoolean(row["disabled"], false), Small: true}
		res, _ := btn.Render()
		return res
	}}},
	{Column: &TableColumn{Id: "id", CellStyle: ut.SM{"color": "red"}}},
}

var testRows []ut.IM = []ut.IM{
	{"id": 1, "name": "Name1", "levels": 0, "valid": "true",
		"date": "2000-03-06", "start": "2019-04-23T05:30:00+02:00", "stamp": "2020-04-20T10:30:00+02:00",
		"name_color":            "red",
		"export_deffield_value": "Customer 1", "fieldtype": "customer", "deffield": 123},
	{"id": 2, "name": "Name2", "export_name": "Name link",
		"levels": 20, "valid": 1,
		"date": "2008-04-07", "start": "2019-04-23T11:30:00+02:00", "stamp": "2020-04-25T10:30:00+02:00",
		"name_color": "red", "edited": true,
		"fieldtype": "bool", "deffield": "true"},
	{"id": 3, "name": "Name3", "levels": 40, "valid": "false",
		"date": "2022-01-01", "start": "2019-04-23T10:27:00+02:00", "stamp": "2020-04-09T10:30:00+02:00",
		"name_color": "orange", "disabled": true,
		"fieldtype": "integer", "deffield": 123},
	{"id": 4, "name": "Name4", "levels": 401234.345, "valid": 0,
		"date": "2015-07-26", "start": "", "stamp": time.Now(),
		"name_color": "orange",
		"fieldtype":  "string", "deffield": "value"},
	{"id": 5, "name": "Name5", "levels": 40, "valid": false,
		"date": "1999-11-07", "start": "2019-04-23T10:30:00+02:00", "stamp": "2020-04-11T10:30:00+02:00",
		"export_deffield_value": "Customer 2", "fieldtype": "customer", "deffield": 124},
	{"id": 6, "name": "Name6", "levels": 60, "valid": true,
		"date": "2020-06-06", "start": "2019-04-23T04:10:00+02:00", "stamp": "2020-04-18T10:30:00+02:00",
		"name_color":            "green",
		"export_deffield_value": "Customer 7", "fieldtype": "customer", "deffield": 222},
}

var demoTableResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
	switch evt.Name {
	case TableEventAddItem, TableEventEditCell, TableEventRowSelected:
		re = ResponseEvent{
			Trigger: &Toast{
				Type:    ToastTypeInfo,
				Value:   evt.TriggerName,
				Timeout: 4,
			},
			TriggerName: evt.TriggerName,
			Name:        evt.Name,
			Header: ut.SM{
				HeaderRetarget: "#toast-msg",
				HeaderReswap:   "innerHTML",
			},
		}
		return re
	}
	return evt
}

func TestTable(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeTable,
			Component: &Table{
				BaseComponent: BaseComponent{
					Id:           id + "_table_default",
					EventURL:     eventURL,
					OnResponse:   demoTableResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Rows:       testRows,
				Fields:     testFields,
				Pagination: PaginationTypeNone,
				PageSize:   10,
			}},
		{
			Label:         "String table, top pagination, row selected",
			ComponentType: ComponentTypeTable,
			Component: &Table{
				BaseComponent: BaseComponent{
					Id:           id + "_table_string_row_selected",
					EventURL:     eventURL,
					OnResponse:   demoTableResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Rows: []ut.IM{
					{"name": "Fluffy", "age": 9, "breed": "calico", "gender": "male"},
					{"name": "Luna", "age": 10, "breed": "long hair", "gender": "female"},
					{"name": "Cracker", "age": 8, "breed": "fat", "gender": "male"},
					{"name": "Pig", "age": 6, "breed": "calico", "gender": "female"},
					{"name": "Robin", "age": 7, "breed": "long hair", "gender": "male"},
					{"name": "Sammy", "age": 13, "breed": "fat", "gender": "male"},
					{"name": "Aliece", "age": 9, "breed": "long hair", "gender": "female"},
					{"name": "Mehatable", "age": 5, "breed": "calico", "gender": "female"},
					{"name": "Scorpia", "age": 6, "breed": "long hair", "gender": "female"},
					{"name": "Zoomies", "age": 1, "breed": "fat", "gender": "male"},
					{"name": "Zues", "age": 5, "breed": "long hair", "gender": "male"},
					{"name": "Lord Kittybottom", "age": 9, "breed": "calico", "gender": "male"},
					{"name": "Princess Furball", "age": 5, "breed": "calico", "gender": "female"},
					{"name": "Delerium", "age": 4, "breed": "fat", "gender": "female"},
				},
				Pagination:  PaginationTypeTop,
				PageSize:    5,
				RowSelected: true,
				SortCol:     "name",
			},
		},
		{
			Label:         "Bottom pagination",
			ComponentType: ComponentTypeTable,
			Component: &Table{
				BaseComponent: BaseComponent{
					Id:           id + "_table_bottom_pagination",
					EventURL:     eventURL,
					OnResponse:   demoTableResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Rows:              testRows,
				Fields:            testFields,
				Pagination:        PaginationTypeBottom,
				PageSize:          5,
				CurrentPage:       10,
				TableFilter:       true,
				FilterPlaceholder: "Placeholder text",
				AddIcon:           "Check",
				AddItem:           true,
				TablePadding:      "16px",
			}},
		{
			Label:         "Filtered",
			ComponentType: ComponentTypeTable,
			Component: &Table{
				BaseComponent: BaseComponent{
					Id:           id + "_table_filtered",
					EventURL:     eventURL,
					OnResponse:   demoTableResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Rows:        testRows,
				Fields:      testFields,
				Pagination:  PaginationTypeAll,
				CurrentPage: 1,
				TableFilter: true,
				FilterValue: "123",
				LabelAdd:    "Add new",
				AddItem:     true,
				SortCol:     "name",
				SortAsc:     true,
			}},
	}
}
