package component

import (
	"fmt"
	"html/template"
	"math"
	"slices"
	"sort"
	"strings"
	"time"

	ut "github.com/nervatura/component/pkg/util"
)

// [Table] constants
const (
	ComponentTypeTable = "table"

	TableEventCurrentPage  = "current_page"
	TableEventFilterChange = "filter_change"
	TableEventAddItem      = "add_item"
	TableEventEditCell     = "edit_cell"
	TableEventRowSelected  = "row_selected"
	TableEventSort         = "table_sort"

	TableFieldTypeString   = "string"
	TableFieldTypeInteger  = "integer"
	TableFieldTypeNumber   = "float"
	TableFieldTypeDateTime = "datetime"
	TableFieldTypeDate     = "date"
	TableFieldTypeTime     = "time"
	TableFieldTypeBool     = "bool"
	TableFieldTypeLink     = "link"
	TableFieldTypeMeta     = "meta"
)

// [Table] TableFieldType values
var TableFieldType []string = []string{TableFieldTypeString, TableFieldTypeInteger, TableFieldTypeNumber,
	TableFieldTypeDateTime, TableFieldTypeDate, TableFieldTypeTime, TableFieldTypeBool, TableFieldTypeLink,
	TableFieldTypeMeta}

// [Table] TableMetaType values
var TableMetaType []string = []string{TableFieldTypeString, TableFieldTypeInteger, TableFieldTypeNumber,
	TableFieldTypeDateTime, TableFieldTypeDate, TableFieldTypeTime, TableFieldTypeBool, TableFieldTypeLink}

/*
Creates an interactive and customizable table control

For example:

	&Table{
	  BaseComponent: BaseComponent{
	    Id: "id_table_string_row_selected",
	  },
	  Rows: []ut.IM{
	    {"name": "Fluffy", "age": 9, "breed": "calico", "gender": "male"},
	    {"name": "Luna", "age": 10, "breed": "long hair", "gender": "female"},
	    {"name": "Cracker", "age": 8, "breed": "fat", "gender": "male"},
	    {"name": "Pig", "age": 6, "breed": "calico", "gender": "female"},
	  },
	  Pagination:  PaginationTypeTop,
	  PageSize:    5,
	  RowSelected: true,
	  SortCol:     "name",
	}
*/
type Table struct {
	BaseComponent
	/* The field name containing the row ID of the data source. If not specified,
	the row index will be used
	*/
	RowKey string `json:"row_key"`
	// Data source of the table
	Rows []ut.IM `json:"rows"`
	// Table column definitions
	Fields []TableField `json:"fields"`
	/* [PaginationType] variable constants:
	[PaginationTypeTop], [PaginationTypeBottom], [PaginationTypeAll], [PaginationTypeNone].
	Default value: [PaginationTypeTop] */
	Pagination string `json:"pagination"`
	// Pagination start value
	CurrentPage int64 `json:"current_page"`
	// Pagination component [PageSize] variable constants: 5, 10, 20, 50, 100. Default value: 10
	PageSize int64 `json:"page_size"`
	// [Pagination] component show/hide page size selector
	HidePaginatonSize bool `json:"hide_paginaton_size"`
	// Show/hide table value filter input row
	TableFilter bool `json:"table_filter"`
	// Show/hide table add item button
	AddItem bool `json:"add_item"`
	// Specifies a short hint that describes the expected value of the input element
	FilterPlaceholder string `json:"filter_placeholder"`
	// Filter input value
	FilterValue string `json:"filter_value"`
	// The filter is case sensitive
	CaseSensitive bool `json:"case_sensitive"`
	// A true value caption in the table cell. Default value: YES
	LabelYes string `json:"label_yes"`
	// A false value caption in the table cell. Default value: NO
	LabelNo string `json:"label_no"`
	// Add item button caption Default empty string
	LabelAdd string `json:"label_add"`
	// Valid [Icon] component value. See more [IconKey] variable values.
	AddIcon string `json:"add_icon"`
	// Table cell padding style value. Example: 8px
	TablePadding string `json:"table_padding"`
	// The order of the table is based on the field name
	SortCol string `json:"sort_col"`
	// Sort in ascending or descending order
	SortAsc bool `json:"sort_asc"`
	// Select an entire row or cell
	RowSelected bool `json:"row_selected"`
}

// [Table] column definition
type TableField struct {
	// The field name of the data source
	Name string `json:"name"`
	/* [TableFieldType] variable constants:
	[TableFieldTypeString], [TableFieldTypeInteger], [TableFieldTypeNumber], [TableFieldTypeDateTime],
	[TableFieldTypeDate], [TableFieldTypeTime], [TableFieldTypeBool], [TableFieldTypeLink],
	[TableFieldTypeMeta].
	Default value: [TableFieldTypeString] */
	FieldType string `json:"field_type"`
	// The label of the column
	Label string `json:"label"`
	/* [TextAlign] variable constants: [TextAlignLeft], [TextAlignCenter], [TextAlignRight].
	Default value: [TextAlignLeft] */
	TextAlign string `json:"text_align"`
	/* [VerticalAlign] variable constants:
	[VerticalAlignTop], [VerticalAlignMiddle], [VerticalAlignBottom].
	Default value: [VerticalAlignMiddle] */
	VerticalAlign string `json:"vertical_align"`
	/* Formatting of negative (red), positive (green) and zero (line-through) values
	in the case of a number field */
	Format bool `json:"format"`
	// Custom column definition
	Column *TableColumn `json:"-"`
}

// [Table] column
type TableColumn struct {
	// The field name of the data source
	Id string `json:"id"`
	// The label of the column
	Header string `json:"header"`
	// Header cell style settings. Example: ut.SM{"padding": "4px"}
	HeaderStyle ut.SM `json:"header_style"`
	// Data cell style settings. Example: ut.SM{"color": "red"}
	CellStyle ut.SM `json:"cell_style"`
	// Original field definition
	Field TableField `json:"field"`
	// The cell generator function of the table
	Cell func(row ut.IM, col TableColumn, value interface{}) template.HTML `json:"-"`
}

/*
Returns all properties of the [Table]
*/
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
			"case_sensitive":      tbl.CaseSensitive,
			"label_yes":           tbl.LabelYes,
			"label_no":            tbl.LabelNo,
			"label_add":           tbl.LabelAdd,
			"add_icon":            tbl.AddIcon,
			"table_padding":       tbl.TablePadding,
			"sort_col":            tbl.SortCol,
			"sort_asc":            tbl.SortAsc,
		})
}

/*
Returns the value of the property of the [Table] with the specified name.
*/
func (tbl *Table) GetProperty(propName string) interface{} {
	return tbl.Properties()[propName]
}

func (tbl *Table) tableFieldsValidation(propValue interface{}) []TableField {
	fields := []TableField{}
	if fd, valid := propValue.([]TableField); valid && (fd != nil) {
		fields = fd
	}
	if tblFields, valid := propValue.([]interface{}); valid {
		for _, tblField := range tblFields {
			if values, valid := tblField.(ut.IM); valid {
				fields = append(fields, TableField{
					Name:          ut.ToString(values["name"], ""),
					FieldType:     tbl.CheckEnumValue(ut.ToString(values["field_type"], ""), TableFieldTypeString, TableFieldType),
					Label:         ut.ToString(values["label"], ""),
					TextAlign:     tbl.CheckEnumValue(ut.ToString(values["text_align"], ""), TextAlignLeft, TextAlign),
					VerticalAlign: tbl.CheckEnumValue(ut.ToString(values["vertical_align"], ""), VerticalAlignMiddle, VerticalAlign),
					Format:        ut.ToBoolean(values["format"], false),
				})
			}
		}
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
}

/*
It checks the value given to the property of the [Table] and always returns a valid value
*/
func (tbl *Table) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"row_key": func() interface{} {
			return ut.ToString(propValue, "id")
		},
		"rows": func() interface{} {
			return ut.ToIMA(propValue, []ut.IM{})
		},
		"fields": func() interface{} {
			return tbl.tableFieldsValidation(propValue)
		},
		"pagination": func() interface{} {
			return tbl.CheckEnumValue(ut.ToString(propValue, ""), PaginationTypeTop, PaginationType)
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
			if !slices.Contains(pageSize, ut.ToString(value, "")) {
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

/*
Setting a property of the [Table] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
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
		"case_sensitive": func() interface{} {
			tbl.CaseSensitive = ut.ToBoolean(propValue, false)
			return tbl.CaseSensitive
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
		TableFieldTypeInteger: func(i, j int) bool {
			a := ut.ToInteger(tbl.Rows[i][fieldName], 0)
			b := ut.ToInteger(tbl.Rows[j][fieldName], 0)
			if sortAsc {
				return a > b
			}
			return a < b
		},
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
		Header: ut.SM{HeaderRetarget: "#" + tbl.Id},
	}
	switch evt.TriggerName {
	case "top_pagination", "bottom_pagination":
		if evt.Name == PaginationEventPageSize {
			tbl.SetProperty("page_size", tblEvt.Value)
			tbl.SetProperty("current_page", 1)
		} else {
			tblEvt.Name = TableEventCurrentPage
			tbl.SetProperty("current_page", tblEvt.Value)
		}

	case "header_sort":
		sortCol := ut.ToString(evt.Trigger.GetProperty("data").(ut.IM)["fieldname"], "")
		fieldType := ut.ToString(evt.Trigger.GetProperty("data").(ut.IM)["fieldtype"], "")
		if tbl.SortCol == sortCol {
			tbl.SetProperty("sort_asc", !tbl.SortAsc)
		}
		tbl.SetProperty("sort_col", sortCol)
		tbl.SortRows(tbl.SortCol, fieldType, tbl.SortAsc)
		tblEvt.Name = TableEventSort
		tblEvt.Value = sortCol

	case "filter", "btn_add", "link_cell", "data_row":
		evtMap := map[string]func(){
			"filter": func() {
				tblEvt.Name = TableEventFilterChange
				tbl.SetProperty("filter_value", tblEvt.Value)
			},
			"btn_add": func() {
				tblEvt.Name = TableEventAddItem
			},
			"link_cell": func() {
				tblEvt.Name = TableEventEditCell
				tblEvt.Value = evt.Trigger.GetProperty("data")
			},
			"data_row": func() {
				tblEvt.Name = TableEventRowSelected
				tblEvt.Value = evt.Trigger.GetProperty("data")
			},
		}
		evtMap[evt.TriggerName]()

	default:
	}
	if tbl.OnResponse != nil {
		return tbl.OnResponse(tblEvt)
	}
	return tblEvt
}

func (tbl *Table) getComponent(name string, pageCount int64, data ut.IM) (html template.HTML, err error) {
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
	ccInp := func(value string) *Input {
		inp := &Input{
			BaseComponent: BaseComponent{
				Id: tbl.Id + "_" + name, Name: name,
				Style:        ut.SM{"border-radius": "0", "margin": "1px 0 2px"},
				EventURL:     tbl.EventURL,
				Target:       tbl.Target,
				OnResponse:   tbl.response,
				RequestValue: tbl.RequestValue,
				RequestMap:   tbl.RequestMap,
			},
			Type:        InputTypeString,
			Label:       tbl.FilterPlaceholder,
			Placeholder: tbl.FilterPlaceholder,
			Value:       tbl.FilterValue,
			Full:        true,
		}
		inp.SetProperty("value", value)
		return inp
	}
	ccMap := map[string]func() ClientComponent{
		"top_pagination": func() ClientComponent {
			return ccPgn()
		},
		"bottom_pagination": func() ClientComponent {
			return ccPgn()
		},
		"filter": func() ClientComponent {
			return ccInp(tbl.FilterValue)
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
				ButtonStyle: ButtonStyleBorder,
				Icon:        tbl.AddIcon, Label: tbl.LabelAdd,
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
	html, err = cc.Render()
	return html, err
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

type cellFormatOptions struct {
	Value       interface{}
	Label       string
	DateType    string
	FieldName   string
	ResultValue interface{}
	Style       ut.SM
	RowData     ut.IM
}

func (tbl *Table) cellFormat(fmtType string, options cellFormatOptions) template.HTML {
	fmtMap := map[string]func() template.HTML{
		"number": func() template.HTML {
			return template.HTML(fmt.Sprintf(
				`<div class="number-cell">
				<span class="cell-label">%s</span>
				<span %s >%s</span>
			</div>`, options.Label, tbl.getStyle(options.Style), ut.ToString(options.Value, "0")))
		},
		"date": func() template.HTML {
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
			switch v := options.Value.(type) {
			case string:
				tmValue, _ := ut.StringToDateTime(v)
				fmtValue = dateFormat[options.DateType](tmValue)
			case time.Time:
				fmtValue = dateFormat[options.DateType](v)
			}
			return template.HTML(fmt.Sprintf(`<span class="cell-label">%s</span><span>%s</span>`, options.Label, fmtValue))
		},
		"bool": func() template.HTML {
			if (options.Value == 1) || (options.Value == "true") || (options.Value == true) {
				return template.HTML(fmt.Sprintf(
					`<span class="cell-label">%s</span>
					<form-icon iconKey="CheckSquare" ></form-icon>
					<span class="middle"> %s</span>`, options.Label, tbl.LabelYes))
			}
			return template.HTML(fmt.Sprintf(
				`<span class="cell-label">%s</span>
				<form-icon iconKey="SquareEmpty" ></form-icon>
				<span class="middle"> %s</span>`, options.Label, tbl.LabelNo))
		},
		"link": func() template.HTML {
			linkLabel := fmt.Sprintf(
				`<span class="cell-label">%s</span>`, options.Label)
			var link template.HTML
			link, _ = tbl.getComponent("link_cell", 0, ut.IM{
				"value": options.Value, "fieldname": options.FieldName, "result": options.ResultValue, "row": options.RowData,
			})
			return template.HTML(linkLabel + string(link))
		},
		"string": func() template.HTML {
			return template.HTML(fmt.Sprintf(
				`<span class="cell-label">%s</span>
				<span %s >%s</span>`, options.Label, tbl.getStyle(options.Style), options.Value))
		},
	}
	return fmtMap[fmtType]()
}

func (tbl *Table) columns() (cols []TableColumn) {
	getFieldType := func(fType string) string {
		bt := ut.SM{
			TableFieldTypeInteger:  TableFieldTypeNumber,
			TableFieldTypeNumber:   TableFieldTypeNumber,
			TableFieldTypeDate:     TableFieldTypeDateTime,
			TableFieldTypeTime:     TableFieldTypeDateTime,
			TableFieldTypeDateTime: TableFieldTypeDateTime,
			TableFieldTypeBool:     TableFieldTypeBool,
			TableFieldTypeLink:     TableFieldTypeLink,
			TableFieldTypeMeta:     TableFieldTypeMeta,
		}
		if fieldType, found := bt[fType]; found {
			return fieldType
		}
		return TableFieldTypeString
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

			setFieldType := map[string]func(){
				TableFieldTypeNumber: func() {
					coldef.HeaderStyle["text-align"] = TextAlignRight
					coldef.Cell = func(row ut.IM, col TableColumn, value interface{}) template.HTML {
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
						return tbl.cellFormat("number", cellFormatOptions{
							Value: ut.ToFloat(value, 0),
							Label: col.Field.Label,
							Style: style,
						})
					}
				},
				TableFieldTypeDateTime: func() {
					coldef.Cell = func(row ut.IM, col TableColumn, value interface{}) template.HTML {
						return tbl.cellFormat("date", cellFormatOptions{
							Value:    value,
							Label:    col.Field.Label,
							DateType: col.Field.FieldType,
						})
					}
				},
				TableFieldTypeBool: func() {
					coldef.Cell = func(row ut.IM, col TableColumn, value interface{}) template.HTML {
						return tbl.cellFormat("bool", cellFormatOptions{
							Value: value,
							Label: col.Field.Label,
						})
					}
				},
				TableFieldTypeLink: func() {
					coldef.Cell = func(row ut.IM, col TableColumn, value interface{}) template.HTML {
						return tbl.cellFormat("link", cellFormatOptions{
							Value:       value,
							Label:       col.Field.Label,
							FieldName:   col.Field.Name,
							ResultValue: row[col.Field.Name],
							RowData:     row,
						})
					}
				},
				TableFieldTypeMeta: func() {
					coldef.Cell = func(row ut.IM, col TableColumn, value interface{}) template.HTML {
						fieldType := tbl.CheckEnumValue(ut.ToString(row[field.Name+"_meta"], ""), TableFieldTypeString, TableMetaType)
						mResult := map[string]func() template.HTML{
							TableFieldTypeBool: func() template.HTML {
								return tbl.cellFormat("bool", cellFormatOptions{
									Value: value,
									Label: col.Field.Label,
								})
							},
							TableFieldTypeInteger: func() template.HTML {
								return tbl.cellFormat("number", cellFormatOptions{
									Value: ut.ToFloat(value, 0),
									Label: col.Field.Label,
									Style: ut.SM{},
								})
							},
							TableFieldTypeNumber: func() template.HTML {
								return tbl.cellFormat("number", cellFormatOptions{
									Value: ut.ToFloat(value, 0),
									Label: col.Field.Label,
									Style: ut.SM{},
								})
							},
							TableFieldTypeLink: func() template.HTML {
								return tbl.cellFormat("link", cellFormatOptions{
									Value:       ut.ToString(value, ""),
									Label:       col.Field.Label,
									FieldName:   field.Name,
									ResultValue: row[field.Name],
									RowData:     row,
								})
							},
						}
						if slices.Contains([]string{TableFieldTypeBool, TableFieldTypeInteger, TableFieldTypeNumber, TableFieldTypeLink}, fieldType) {
							return mResult[fieldType]()
						}
						return tbl.cellFormat("string", cellFormatOptions{
							Value: ut.ToString(value, ""),
							Label: col.Field.Label,
							Style: ut.SM{},
						})
					}
				},
				TableFieldTypeString: func() {
					coldef.Cell = func(row ut.IM, col TableColumn, value interface{}) template.HTML {
						style := ut.SM{}
						if color, found := row[col.Field.Name+"_color"].(string); found {
							style["color"] = color
						}
						for key, ivalue := range row {
							if key == col.Field.Name+"_link" {
								return tbl.cellFormat("link", cellFormatOptions{
									Value:       ut.ToString(ivalue, ""),
									Label:       col.Field.Label,
									FieldName:   col.Field.Name,
									ResultValue: row[col.Field.Name],
									RowData:     row,
								})
							}
						}
						return tbl.cellFormat("string", cellFormatOptions{
							Value: ut.ToString(value, ""),
							Label: col.Field.Label,
							Style: style,
						})
					}
				},
			}
			setFieldType[getFieldType(field.FieldType)]()

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
	caseValue := func(value string) string {
		if !tbl.CaseSensitive {
			return strings.ToLower(value)
		}
		return value
	}
	getValidRow := func(row ut.IM, filter string) bool {
		for field := range row {
			if strings.Contains(caseValue(ut.ToString(row[field], "")), filter) {
				return true
			}
		}
		return false
	}
	if tbl.FilterValue == "" {
		return tbl.Rows
	}
	for _, row := range tbl.Rows {
		if getValidRow(row, caseValue(tbl.FilterValue)) {
			rows = append(rows, row)
		}
	}
	return rows
}

/*
Based on the values, it will generate the html code of the [Table] or return with an error message.
*/
func (tbl *Table) Render() (html template.HTML, err error) {
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
		"tableComponent": func(name string) (template.HTML, error) {
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
		"cellValue": func(row ut.IM, col TableColumn) template.HTML {
			if col.Cell != nil {
				return col.Cell(row, col, row[col.Id])
			}
			return template.HTML(ut.ToString(row[col.Id], ""))
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
	{{ if ne $.Indicator "none" }} hx-indicator="#{{ $.Indicator }}"{{ end }} 
	{{ if cellStyle $col.HeaderStyle }} style="{{ range $key, $value := $col.HeaderStyle }}{{ $key }}:{{ $value }};{{ end }}"{{ end }} 
	>{{ $col.Header }}</th>
	{{ end }}</tr></thead>
	<tbody>{{ range $index, $row := pageRows }}
	<tr id="{{ rowID $row $index }}" class="{{ pointerClass $row }}" 
	{{ if and ($.RowSelected) (ne $.EventURL "") }} hx-post="{{ $.EventURL }}" hx-target="{{ $.Target }}" hx-swap="{{ $.Swap }}"{{ end }}
	{{ if and ($.RowSelected) (ne $.Indicator "none") }} hx-indicator="#{{ $.Indicator }}"{{ end }}
	>{{ range $icol, $col := cols }}<td
	{{ if cellStyle $col.CellStyle }} style="{{ range $key, $value := $col.CellStyle }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	>{{ cellValue $row $col }}</td>{{ end }}</tr>
	{{ end }}</tbody>
	</table></div>
	{{ if bottomPagination }}<div>{{ tableComponent "bottom_pagination" }}</div>{{ end }}
	</div>`

	return ut.TemplateBuilder("table", tpl, funcMap, tbl)
}

var testTableFields []TableField = []TableField{
	{Name: "name", FieldType: TableFieldTypeString, Label: "Name", TextAlign: TextAlignLeft},
	{Name: "valid", FieldType: TableFieldTypeBool, Label: "Valid"},
	{Name: "date", FieldType: TableFieldTypeDate, Label: "From"},
	{Name: "start", FieldType: TableFieldTypeTime},
	{Name: "stamp", FieldType: TableFieldTypeDateTime, Label: "Stamp"},
	{Name: "levels", FieldType: TableFieldTypeNumber, Label: "Levels", Format: true, VerticalAlign: VerticalAlignMiddle},
	{Name: "product", FieldType: TableFieldTypeLink, Label: "Product"},
	{Name: "deffield", FieldType: TableFieldTypeMeta, Label: "Multiple type"},
	{Column: &TableColumn{Id: "editor", Header: "Custom",
		Cell: func(row ut.IM, col TableColumn, value interface{}) template.HTML {
			btn := Button{
				ButtonStyle: ButtonStylePrimary, Label: "Hello", Disabled: ut.ToBoolean(row["disabled"], false), Small: true}
			res, _ := btn.Render()
			return res
		}}},
	{Column: &TableColumn{Id: "id", CellStyle: ut.SM{"color": "red"}}},
}

var testTableRows []ut.IM = []ut.IM{
	{"id": 1, "name": "Name1", "levels": 0, "valid": "true",
		"date": "2000-03-06", "start": "2019-04-23T05:30:00+02:00", "stamp": "2020-04-20T10:30:00+02:00",
		"name_color": "red", "product": "Product1",
		"deffield": "Customer 1", "deffield_meta": TableFieldTypeLink},
	{"id": 2, "name": "Name2", "name_link": "Name Link",
		"levels": 20, "valid": 1,
		"date": "2008-04-07", "start": "2019-04-23T11:30:00+02:00", "stamp": "2020-04-25T10:30:00+02:00",
		"name_color": "red", "edited": true, "product": "Product2",
		"deffield": "true", "deffield_meta": TableFieldTypeBool},
	{"id": 3, "name": "Name3", "levels": 40, "valid": "false",
		"date": "2022-01-01", "start": "2019-04-23T10:27:00+02:00", "stamp": "2020-04-09T10:30:00+02:00",
		"name_color": "orange", "disabled": true, "product": "Product2",
		"deffield": 123, "deffield_meta": TableFieldTypeInteger},
	{"id": 4, "name": "Name4", "levels": 40, "valid": "false",
		"date": "2022-01-01", "start": "2019-04-23T10:27:00+02:00", "stamp": "2020-04-09T10:30:00+02:00",
		"name_color": "orange", "disabled": true, "product": "Product1",
		"deffield": 123.45, "deffield_meta": TableFieldTypeNumber},
	{"id": 5, "name": "Name5", "levels": 401234.345, "valid": 0,
		"date": "2015-07-26", "start": "", "stamp": time.Now(),
		"name_color": "orange", "product": "Product3",
		"deffield": "value Orange", "deffield_meta": TableFieldTypeString},
	{"id": 6, "name": "Name6", "levels": 40, "valid": false,
		"date": "1999-11-07", "start": "2019-04-23T10:30:00+02:00", "stamp": "2020-04-11T10:30:00+02:00",
		"product": "Product1", "deffield": "Customer 2", "deffield_meta": TableFieldTypeLink},
	{"id": 7, "name": "Name7", "levels": 60, "valid": true,
		"date": "2020-06-06", "start": "2019-04-23T04:10:00+02:00", "stamp": "2020-04-18T10:30:00+02:00",
		"name_color": "green", "product": "Product2",
		"deffield": "Customer 7", "deffield_meta": TableFieldTypeLink},
}

var testTableResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
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
				HeaderReswap:   SwapInnerHTML,
			},
		}
		return re
	}
	return evt
}

// [Table] test and demo data
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
					OnResponse:   testTableResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Rows:       testTableRows,
				Fields:     testTableFields,
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
					OnResponse:   testTableResponse,
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
					OnResponse:   testTableResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Rows:              testTableRows,
				Fields:            testTableFields,
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
					OnResponse:   testTableResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Rows:          testTableRows,
				Fields:        testTableFields,
				Pagination:    PaginationTypeAll,
				CurrentPage:   1,
				TableFilter:   true,
				FilterValue:   "123",
				CaseSensitive: false,
				LabelAdd:      "Add new",
				AddItem:       true,
				SortCol:       "name",
				SortAsc:       true,
			}},
		{
			Label:         "Filtered CaseSensitive",
			ComponentType: ComponentTypeTable,
			Component: &Table{
				BaseComponent: BaseComponent{
					Id:           id + "_table_filtered_cs",
					EventURL:     eventURL,
					OnResponse:   testTableResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Rows:          testTableRows,
				Fields:        testTableFields,
				Pagination:    PaginationTypeAll,
				CurrentPage:   1,
				TableFilter:   true,
				FilterValue:   "Orange",
				CaseSensitive: true,
				LabelAdd:      "Add new",
				AddItem:       true,
				SortCol:       "name",
				SortAsc:       true,
			}},
	}
}
