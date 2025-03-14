package component

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"html/template"
	"slices"
	"strings"
	"time"

	ut "github.com/nervatura/component/pkg/util"
)

// [Browser] constants
const (
	ComponentTypeBrowser = "browser"

	BrowserEventChange       = "browser_change"
	BrowserEventSearch       = "browser_search"
	BrowserEventBookmark     = "browser_bookmark"
	BrowserEventExport       = "browser_export"
	BrowserEventHelp         = "browser_help"
	BrowserEventView         = "browser_view"
	BrowserEventAddFilter    = "browser_add_filter"
	BrowserEventChangeFilter = "browser_change_filter"
	BrowserEventShowTotal    = "browser_show_total"
	BrowserEventSetColumn    = "browser_set_column"
	BrowserEventEditRow      = "browser_edit_row"
	BrowserExportLimit       = 40000
)

var browserDefaultLabel ut.SM = ut.SM{
	"browser_title":        "Data browser",
	"browser_view":         "Data view",
	"browser_views":        "Views",
	"browser_columns":      "Columns",
	"browser_total":        "Total",
	"browser_filter":       "Filter",
	"browser_comparison":   "?",
	"browser_filters":      "Filter criterias",
	"browser_search":       "Search",
	"browser_bookmark":     "Bookmark",
	"browser_export":       "Export",
	"browser_help":         "Help",
	"browser_result":       "record(s) found",
	"browser_placeholder":  "Filter",
	"browser_label_new":    "NEW",
	"browser_label_delete": "Delete",
	"browser_label_yes":    "YES",
	"browser_label_no":     "NO",
	"browser_label_ok":     "OK",
	"browser_export_error": "There is too much data, it cannot be exported. Use the Filters to limit the displayed result!",
}

var browserFilterComp []SelectOption = []SelectOption{
	{Value: "==", Text: "=="},
	{Value: "!=", Text: "!="},
	{Value: "<", Text: "<"},
	{Value: "<=", Text: "<="},
	{Value: ">", Text: ">"},
	{Value: ">=", Text: ">="},
}

// [Browser] filter query filter type
type BrowserFilter struct {
	Or    bool        `json:"or"`    // and (False) or (True)
	Field string      `json:"field"` // Fieldname and alias
	Comp  string      `json:"comp"`  // ==,!=,<,<=,>,>=
	Value interface{} `json:"value"`
}

// [Browser] meta data definition
type BrowserMetaField struct {
	/* [TableMetaType] variable constants:
	[TableFieldTypeString], [TableFieldTypeInteger], [TableFieldTypeNumber], [TableFieldTypeDateTime],
	[TableFieldTypeDate], [TableFieldTypeTime], [TableFieldTypeBool], [TableFieldTypeLink].
	Default value: [TableFieldTypeString] */
	FieldType string `json:"field_type"`
	// The label of the column
	Label string `json:"label"`
}

// [Browser] field total data definition
type BrowserTotalField struct {
	// The field name of the data source
	Name string `json:"name"`
	// [TableFieldTypeInteger], [TableFieldTypeNumber], [TableFieldTypeMeta]
	FieldType string `json:"field_type"`
	// The label of the field
	Label string `json:"label"`
	// The sum of the field's value
	Total float64 `json:"value"`
}

/*
Creates an interactive and customizable data search control
*/
type Browser struct {
	Table
	// The caption of the browser
	Title string `json:"title"`
	// Current search settings view value
	View string `json:"view"`
	// The list of values and labels of selectable views.
	// If it does not contain any items, its button will not be displayed
	Views []SelectOption `json:"views"`
	// Show or hide the data filtering and display settings options
	HideHeader bool `json:"hide_header"`
	// Display a menu of selectable views
	ShowDropdown bool `json:"show_dropdown"`
	// Displaying the table column visibility settings
	ShowColumns bool `json:"show_columns"`
	// Displaying the total values of the number fields
	ShowTotal bool `json:"show_total"`
	// Show or hide the bookmark button
	HideBookmark bool `json:"hide_bookmark"`
	// Show or hide the export button
	HideExport bool `json:"hide_export"`
	// Export limit of data rows. Limit the maximum character length of a URL. Ignored if the ExportURL value is not empty
	ExportLimit int64 `json:"export_limit"`
	// Specifies the url for downloading data. If it is not specified, then the built-in limited csv export results
	ExportURL string `json:"export_url"`
	// Specifies the name of the file downloaded from ExportURL. Default value: data.csv
	Download string `json:"download"`
	// Show or hide the help button
	HideHelp bool `json:"hide_help"`
	// Specifies the url for help. If it is not specified, then the built-in button event
	HelpURL string `json:"help_url"`
	// Editability of table rows
	ReadOnly bool `json:"readonly"`
	// The name of the columns to be displayed from the data source
	VisibleColumns map[string]bool `json:"visible_columns"`
	// List of filter criteria
	Filters []BrowserFilter `json:"filters"`
	// The index of the filter criteria to be edited
	FilterIndex int64 `json:"filter_index"`
	// Unfilterable table fields. These field names will not be included in the list
	HideFilters map[string]bool `json:"hide_filters"`
	// Multiple type column filter definitions
	MetaFields map[string]BrowserMetaField `json:"meta_fields"`
	// The texts of the labels of the controls
	Labels      ut.SM `json:"labels"`
	totalFields []BrowserTotalField
}

/*
Returns all properties of the [Browser]
*/
func (bro *Browser) Properties() ut.IM {
	return ut.MergeIM(
		bro.Table.Properties(),
		ut.IM{
			"title":           bro.Title,
			"view":            bro.View,
			"views":           bro.Views,
			"hide_header":     bro.HideHeader,
			"show_dropdown":   bro.ShowDropdown,
			"show_columns":    bro.ShowColumns,
			"show_total":      bro.ShowTotal,
			"hide_bookmark":   bro.HideBookmark,
			"hide_export":     bro.HideExport,
			"export_limit":    bro.ExportLimit,
			"export_url":      bro.ExportURL,
			"download":        bro.Download,
			"readonly":        bro.ReadOnly,
			"hide_help":       bro.HideHelp,
			"help_url":        bro.HelpURL,
			"visible_columns": bro.VisibleColumns,
			"filters":         bro.Filters,
			"filter_index":    bro.FilterIndex,
			"hide_filters":    bro.HideFilters,
			"meta_fields":     bro.MetaFields,
			"labels":          bro.Labels,
		})
}

/*
Returns the value of the property of the [Browser] with the specified name.
*/
func (bro *Browser) GetProperty(propName string) interface{} {
	return bro.Properties()[propName]
}

func (bro *Browser) validationProcessBooleanMap(propValue interface{}, existingMap map[string]bool) map[string]bool {
	// Helper functions to reduce code duplication
	makeEmptyMap := func() map[string]bool {
		return make(map[string]bool)
	}

	value := existingMap
	if len(value) == 0 {
		value = makeEmptyMap()
	}

	switch cols := propValue.(type) {
	case []map[string]bool:
		for _, values := range cols {
			for key, bvalue := range values {
				value[key] = bvalue
			}
		}
	case map[string]bool:
		value = cols
	case map[string]interface{}:
		for key, ivalue := range cols {
			value[key] = ut.ToBoolean(ivalue, false)
		}
	}
	return value
}

func (bro *Browser) validationFilters(propValue interface{}) []BrowserFilter {
	value := bro.Filters
	switch v := propValue.(type) {
	case []BrowserFilter:
		value = v
	case []interface{}:
		value = make([]BrowserFilter, 0, len(v))
		for _, filter := range v {
			if filterMap, ok := filter.(ut.IM); ok {
				value = append(value, BrowserFilter{
					Or:    ut.ToBoolean(filterMap["or"], false),
					Field: ut.ToString(filterMap["field"], ""),
					Comp:  ut.ToString(filterMap["comp"], ""),
					Value: ut.ToString(filterMap["value"], ""),
				})
			}
		}
	}
	if value == nil {
		value = make([]BrowserFilter, 0)
	}
	return value
}

func (bro *Browser) validationMetaFields(propValue interface{}) map[string]BrowserMetaField {
	fields := make(map[string]BrowserMetaField)
	switch mFields := propValue.(type) {
	case map[string]BrowserMetaField:
		for fname, fvalue := range mFields {
			fvalue.FieldType = bro.CheckEnumValue(fvalue.FieldType, TableFieldTypeString, TableMetaType)
			fields[fname] = fvalue
		}
	case ut.IM:
		for fname, fvalue := range mFields {
			if values, ok := fvalue.(ut.IM); ok {
				fieldType := bro.CheckEnumValue(ut.ToString(values["field_type"], ""), TableFieldTypeString, TableMetaType)
				fields[fname] = BrowserMetaField{
					FieldType: fieldType,
					Label:     ut.ToString(values["label"], ""),
				}
			}
		}
	}
	return fields
}

/*
It checks the value given to the property of the [Browser] and always returns a valid value
*/
func (bro *Browser) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"labels": func() interface{} {
			value := ut.ToSM(bro.Labels, ut.SM{})
			switch v := propValue.(type) {
			case ut.SM:
				value = ut.MergeSM(value, v)
			case ut.IM:
				value = ut.MergeSM(value, ut.IMToSM(v))
			}
			if len(value) == 0 {
				value = browserDefaultLabel
			}
			return value
		},
		"views": func() interface{} {
			return SelectOptionRangeValidation(propValue, bro.Views)
		},
		"visible_columns": func() interface{} {
			return bro.validationProcessBooleanMap(propValue, bro.VisibleColumns)
		},
		"hide_filters": func() interface{} {
			return bro.validationProcessBooleanMap(propValue, bro.HideFilters)
		},
		"filters": func() interface{} {
			return bro.validationFilters(propValue)
		},
		"filter_index": func() interface{} {
			value := ut.ToInteger(propValue, 0)
			if value < 0 || value > int64(len(bro.Filters)) {
				value = 0
			}
			return value
		},
		"meta_fields": func() interface{} {
			return bro.validationMetaFields(propValue)
		},
		"target": func() interface{} {
			bro.SetProperty("id", bro.Id)
			value := ut.ToString(propValue, bro.Id)
			if value != "this" && !strings.HasPrefix(value, "#") {
				value = "#" + value
			}
			return value
		},
	}

	if fn, found := pm[propName]; found {
		return fn()
	}
	if bro.Table.GetProperty(propName) != nil {
		return bro.Table.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [Browser] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (bro *Browser) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"title": func() interface{} {
			bro.Title = ut.ToString(propValue, browserDefaultLabel["browser_title"])
			return bro.Title
		},
		"view": func() interface{} {
			bro.View = ut.ToString(propValue, "default")
			return bro.View
		},
		"views": func() interface{} {
			bro.Views = bro.Validation(propName, propValue).([]SelectOption)
			return bro.Views
		},
		"fields": func() interface{} {
			bro.Fields = bro.Table.SetProperty(propName, propValue).([]TableField)
			bro.totalFields = bro.setTotalFields()
			return bro.Fields
		},
		"visible_columns": func() interface{} {
			bro.VisibleColumns = bro.Validation(propName, propValue).(map[string]bool)
			return bro.VisibleColumns
		},
		"filters": func() interface{} {
			bro.Filters = bro.Validation(propName, propValue).([]BrowserFilter)
			return bro.Filters
		},
		"filter_index": func() interface{} {
			bro.FilterIndex = bro.Validation(propName, propValue).(int64)
			return bro.FilterIndex
		},
		"hide_filters": func() interface{} {
			bro.HideFilters = bro.Validation(propName, propValue).(map[string]bool)
			return bro.HideFilters
		},
		"meta_fields": func() interface{} {
			bro.MetaFields = bro.Validation(propName, propValue).(map[string]BrowserMetaField)
			return bro.MetaFields
		},
		"hide_header": func() interface{} {
			bro.HideHeader = ut.ToBoolean(propValue, false)
			return bro.HideHeader
		},
		"show_dropdown": func() interface{} {
			bro.ShowDropdown = ut.ToBoolean(propValue, false)
			return bro.ShowDropdown
		},
		"show_columns": func() interface{} {
			bro.ShowColumns = ut.ToBoolean(propValue, false)
			return bro.ShowColumns
		},
		"show_total": func() interface{} {
			bro.ShowTotal = ut.ToBoolean(propValue, false)
			return bro.ShowTotal
		},
		"hide_bookmark": func() interface{} {
			bro.HideBookmark = ut.ToBoolean(propValue, false)
			return bro.HideBookmark
		},
		"hide_export": func() interface{} {
			bro.HideExport = ut.ToBoolean(propValue, false)
			return bro.HideExport
		},
		"export_limit": func() interface{} {
			bro.ExportLimit = ut.ToInteger(propValue, BrowserExportLimit)
			return bro.ExportLimit
		},
		"export_url": func() interface{} {
			bro.ExportURL = ut.ToString(propValue, "")
			return bro.ExportURL
		},
		"download": func() interface{} {
			bro.Download = ut.ToString(propValue, "data.csv")
			return bro.Download
		},
		"readonly": func() interface{} {
			bro.ReadOnly = ut.ToBoolean(propValue, false)
			return bro.ReadOnly
		},
		"hide_help": func() interface{} {
			bro.HideHelp = ut.ToBoolean(propValue, false)
			return bro.HideHelp
		},
		"help_url": func() interface{} {
			bro.HelpURL = ut.ToString(propValue, "")
			return bro.HelpURL
		},
		"labels": func() interface{} {
			bro.Labels = bro.Validation(propName, propValue).(ut.SM)
			return bro.Labels
		},
		"target": func() interface{} {
			bro.Target = bro.Validation(propName, propValue).(string)
			return bro.Target
		},
	}
	if _, found := pm[propName]; found {
		return bro.SetRequestValue(propName, pm[propName](), []string{})
	}
	if bro.Table.GetProperty(propName) != nil {
		return bro.Table.SetProperty(propName, propValue)
	}
	return propValue
}

func (bro *Browser) exportData() (re ResponseEvent) {

	fieldLabel := func(fieldName string) string {
		for _, field := range bro.Fields {
			if field.Name == fieldName {
				return ut.ToString(field.Label, field.Name)
			}
		}
		return fieldName
	}

	re = ResponseEvent{
		Trigger:     &BaseComponent{},
		TriggerName: bro.Name,
		Name:        BrowserEventExport,
		Header: ut.SM{
			HeaderReswap: SwapNone,
		},
	}

	sRows := make([][]string, 0)
	//labels
	sRow := make([]string, 0)
	for fieldName, visible := range bro.VisibleColumns {
		if visible {
			sRow = append(sRow, fieldLabel(fieldName))
		}
	}
	sRows = append(sRows, sRow)

	// table data
	for _, tRow := range bro.Rows {
		sRow = make([]string, 0)
		for fieldName, visible := range bro.VisibleColumns {
			if visible {
				sRow = append(sRow, ut.ToString(tRow[fieldName], ""))
			}
		}
		sRows = append(sRows, sRow)
	}

	var b bytes.Buffer
	writr := csv.NewWriter(&b)
	if err := writr.WriteAll(sRows); err == nil {
		encURL := base64.URLEncoding.EncodeToString(b.Bytes())
		if len(encURL) > int(bro.ExportLimit) {
			return ResponseEvent{
				Trigger: &Toast{
					Type:  ToastTypeError,
					Value: bro.msg("browser_export_error"),
				},
				TriggerName: bro.Name,
				Name:        BrowserEventExport,
				Header: ut.SM{
					HeaderRetarget: "#toast-msg",
					HeaderReswap:   SwapInnerHTML,
				},
			}
		}
		re.Header[HeaderRedirect] = "data:text/csv;base64," + encURL
	}

	return re
}

func (bro *Browser) filterEvent(evt ResponseEvent) (re ResponseEvent) {
	broEvt := ResponseEvent{
		Trigger: bro, TriggerName: bro.Name, Value: evt.Value,
		Name: BrowserEventChangeFilter,
		Header: ut.SM{
			HeaderRetarget: "#" + bro.Id,
		},
	}

	evtData := ut.ToIM(evt.Value, ut.IM{})
	filterIndex := ut.ToInteger(evtData["index"], 0)
	filterRow := ut.ToIM(evtData["row"], ut.IM{})
	filters := bro.GetProperty("filters").([]BrowserFilter)

	switch evt.Name {
	case TableEventFormDelete:
		if len(filters) > int(filterIndex) {
			filters = append(filters[:filterIndex], filters[filterIndex+1:]...)
		}
		bro.SetProperty("filter_index", 0)
		bro.SetProperty("filters", filters)

	case TableEventFormUpdate:
		filters[filterIndex].Field = ut.ToString(filterRow["field"], "")
		filters[filterIndex].Comp = ut.ToString(filterRow["comp"], "")
		filters[filterIndex].Value = filterRow["value"]
		bro.SetProperty("filter_index", 0)
		bro.SetProperty("filters", filters)

	case TableEventFormChange:
		if ut.ToString(evtData["field"], "") == "field" {
			rows := ut.ToIMA(evt.Trigger.GetProperty("rows"), []ut.IM{})
			field := ut.ToString(evtData["value"], "")
			rows[filterIndex]["comp"] = browserFilterComp[0].Value
			rows[filterIndex]["comp_options"] = bro.filterCompOptions(bro.getFilterType(field))
			rows[filterIndex]["value"] = bro.defaultFilterValue(bro.getFilterType(field))
			rows[filterIndex]["value_meta"] = bro.getFilterType(field)
			evt.Trigger.SetProperty("rows", rows)
		}
		broEvt = evt

	case TableEventFormCancel:
		broEvt.Name = TableEventFormCancel
		broEvt.Value = ut.IM{"row": filterRow, "index": filterIndex}
		bro.SetProperty("filter_index", 0)
		bro.SetProperty("filters", filters)

	case TableEventFormEdit:
		broEvt.Name = TableEventFormEdit
		bro.SetProperty("filter_index", evt.Value)
	}
	if bro.OnResponse != nil {
		return bro.OnResponse(broEvt)
	}
	return broEvt
}

func (bro *Browser) response(evt ResponseEvent) (re ResponseEvent) {
	broEvt := ResponseEvent{
		Trigger: bro, TriggerName: bro.Name, Value: evt.Value,
	}
	if evt.TriggerName != "btn_views" {
		bro.SetProperty("show_dropdown", false)
	}
	bro.SetProperty("filter_index", 0)
	switch evt.TriggerName {
	case "table":
		broEvt = evt

	case "btn_export":
		return bro.exportData()

	case "hide_header", "btn_search", "btn_bookmark", "btn_help", "btn_views", "btn_columns",
		"btn_filter", "btn_total", "menu_item", "col_item", "btn_ok", "edit_row":
		evtMap := map[string]func(){
			"hide_header": func() {
				broEvt.Name = BrowserEventChange
				bro.SetProperty("hide_header", !bro.HideHeader)
			},
			"btn_search": func() {
				broEvt.Name = BrowserEventSearch
			},
			"btn_bookmark": func() {
				broEvt.Name = BrowserEventBookmark
			},
			"btn_help": func() {
				broEvt.Name = BrowserEventHelp
			},
			"btn_views": func() {
				broEvt.Name = BrowserEventChange
				bro.SetProperty("show_dropdown", !bro.ShowDropdown)
			},
			"btn_columns": func() {
				broEvt.Name = BrowserEventChange
				bro.SetProperty("show_columns", !bro.ShowColumns)
			},
			"btn_filter": func() {
				broEvt.Name = BrowserEventAddFilter
				if len(bro.Fields) > 0 {
					filters := bro.GetProperty("filters").([]BrowserFilter)
					filters = append(filters, BrowserFilter{
						Field: bro.Fields[0].Name, Comp: browserFilterComp[0].Value,
						Value: bro.defaultFilterValue(bro.Fields[0].FieldType),
					})
					bro.SetProperty("filters", filters)
					bro.SetProperty("filter_index", int64(len(filters)))
				}
			},
			"btn_total": func() {
				broEvt.Name = BrowserEventShowTotal
				bro.SetProperty("show_total", true)
			},
			"btn_ok": func() {
				broEvt.Name = BrowserEventShowTotal
				bro.SetProperty("show_total", false)
			},
			"edit_row": func() {
				broEvt.Value = evt.Trigger.GetProperty("data")
				broEvt.Name = BrowserEventEditRow
			},
			"menu_item": func() {
				broEvt.Value = ut.ToString(evt.Trigger.GetProperty("data").(ut.IM)["key"], "")
				broEvt.Name = BrowserEventView
				bro.SetProperty("show_dropdown", false)
			},
			"col_item": func() {
				fieldName := ut.ToString(evt.Trigger.GetProperty("data").(ut.IM)["key"], "")
				oldValue := ut.ToBoolean(bro.VisibleColumns[fieldName], false)
				broEvt.Name = BrowserEventSetColumn
				broEvt.Value = fieldName
				bro.SetProperty("visible_columns", []map[string]bool{{fieldName: !oldValue}})
			},
		}
		evtMap[evt.TriggerName]()

	default:
	}
	if bro.OnResponse != nil {
		return bro.OnResponse(broEvt)
	}
	return broEvt
}

func (bro *Browser) defaultFilterValue(ftype string) interface{} {
	defvalue := ut.IM{
		TableFieldTypeNumber:   0,
		TableFieldTypeInteger:  0,
		TableFieldTypeDateTime: time.Now().Format("2006-01-02T15:04"),
		TableFieldTypeDate:     time.Now().Format("2006-01-02"),
		TableFieldTypeTime:     time.Now().Format("15:04"),
		TableFieldTypeBool:     "1",
	}
	if value, found := defvalue[ftype]; found {
		return value
	}
	return ""
}

func (bro *Browser) getFilterType(fieldName string) string {
	fieldType := TableFieldTypeString
	for _, field := range bro.Fields {
		if (field.Name == fieldName) && (field.FieldType != "") && (field.FieldType != TableFieldTypeLink) {
			return field.FieldType
		}
	}
	if mValue, found := bro.MetaFields[fieldName]; found && (mValue.FieldType != TableFieldTypeLink) {
		return mValue.FieldType
	}
	return fieldType
}

func (bro *Browser) setTotalFields() []BrowserTotalField {
	total := []BrowserTotalField{}
	for _, field := range bro.Fields {
		if slices.Contains([]string{TableFieldTypeInteger, TableFieldTypeNumber, TableFieldTypeMeta}, field.FieldType) {
			total = append(total,
				BrowserTotalField{Name: field.Name, FieldType: field.FieldType, Label: field.Label, Total: 0})
		}
	}
	return total
}

func (bro *Browser) setTotalValues() []BrowserTotalField {
	total := bro.setTotalFields()
	for _, row := range bro.Rows {
		for index, field := range total {
			if value, found := row[field.Name]; found && (field.FieldType != TableFieldTypeMeta) {
				total[index].Total += ut.ToFloat(value, 0)
			}
			if value, found := row[field.Name]; found && (field.FieldType == TableFieldTypeMeta) {
				if slices.Contains([]string{TableFieldTypeInteger, TableFieldTypeNumber}, ut.ToString(row[field.Name+"_meta"], "")) {
					total[index].Total += ut.ToFloat(value, 0)
				}
			}
		}
	}
	return total
}

func (bro *Browser) getComponentTable() *Table {
	fields := []TableField{
		{Column: &TableColumn{
			Id:        "edit_row",
			Header:    "",
			CellStyle: ut.SM{"width": "25px", "padding": "7px 3px 3px 8px"},
			Cell: func(row ut.IM, col TableColumn, value interface{}, rowIndex int64) template.HTML {
				var ico template.HTML
				ico, _ = bro.getComponent("edit_row", row)
				return ico
			}}},
	}
	for _, fd := range bro.Fields {
		if ut.ToBoolean(bro.VisibleColumns[fd.Name], false) {
			fields = append(fields, fd)
		}
	}
	tbl := &Table{
		BaseComponent: BaseComponent{
			Id:           bro.Id + "_table",
			Name:         "table",
			EventURL:     bro.EventURL,
			OnResponse:   bro.response,
			RequestValue: bro.RequestValue,
			RequestMap:   bro.RequestMap,
		},
		Fields:            fields,
		Rows:              bro.Rows,
		Pagination:        bro.Pagination,
		HidePaginatonSize: bro.HidePaginatonSize,
		PageSize:          bro.PageSize,
		CurrentPage:       bro.CurrentPage,
		RowKey:            bro.RowKey,
		TableFilter:       bro.TableFilter,
		FilterPlaceholder: bro.msg("browser_placeholder"),
		FilterValue:       bro.FilterValue,
		AddItem:           bro.AddItem,
		LabelAdd:          ut.ToString(bro.LabelAdd, bro.msg("browser_label_new")),
		AddIcon:           bro.AddIcon,
		RowSelected:       bro.RowSelected,
		TablePadding:      bro.TablePadding,
		SortCol:           bro.SortCol,
		SortAsc:           bro.SortAsc,
	}
	return tbl
}

func (bro *Browser) filterCompOptions(ftype string) []SelectOption {
	if !slices.Contains([]string{
		TableFieldTypeDate, TableFieldTypeDateTime, TableFieldTypeTime, TableFieldTypeInteger,
		TableFieldTypeNumber}, ftype) {
		return browserFilterComp[0:2]
	}
	return browserFilterComp
}

func (bro *Browser) filterTable() *Table {
	fieldOptions := func() (options []SelectOption) {
		options = []SelectOption{}
		metaField := false
		for _, field := range bro.Fields {
			if !bro.HideFilters[field.Name] {
				if field.FieldType == TableFieldTypeMeta {
					if !metaField {
						for fName, fValue := range bro.MetaFields {
							options = append(options, SelectOption{Value: fName, Text: fValue.Label})
						}
						metaField = true
					}
				} else {
					options = append(options, SelectOption{Value: field.Name, Text: field.Label})
				}
			}
		}
		return options
	}
	var fields []TableField = []TableField{
		{Name: "field", FieldType: TableFieldTypeString, Label: bro.msg("browser_filter"),
			Options: fieldOptions(), TriggerEvent: true, Required: true},
		{Name: "comp", FieldType: TableFieldTypeMeta, Label: bro.msg("browser_comparison"), Required: true},
		{Name: "value", FieldType: TableFieldTypeMeta, Label: bro.msg("browser_value")},
	}
	rows := []ut.IM{}
	for _, filter := range bro.Filters {
		rows = append(rows, ut.IM{
			"field": filter.Field,
			"comp":  filter.Comp, "comp_meta": TableFieldTypeString, "comp_options": bro.filterCompOptions(bro.getFilterType(filter.Field)),
			"value": filter.Value, "value_meta": bro.getFilterType(filter.Field),
		})
	}
	tbl := &Table{
		BaseComponent: BaseComponent{
			Id:           bro.Id + "_filter_table",
			Name:         "filter_table",
			EventURL:     bro.EventURL,
			OnResponse:   bro.filterEvent,
			RequestValue: bro.RequestValue,
			RequestMap:   bro.RequestMap,
		},
		Fields:      fields,
		Pagination:  PaginationTypeNone,
		PageSize:    5,
		TableFilter: false,
		Editable:    true,
		HideHeader:  true,
	}
	tbl.SetProperty("rows", rows)
	tbl.SetProperty("edit_index", bro.FilterIndex)
	return tbl
}

func (bro *Browser) getComponent(name string, data ut.IM) (html template.HTML, err error) {
	ccBtn := func(icoKey, label, bstyle, index string) *Button {
		btn := &Button{
			BaseComponent: BaseComponent{
				Id:           bro.Id + "_" + name + "_" + index,
				Name:         name,
				Data:         ut.IM{"index": index},
				EventURL:     bro.EventURL,
				Target:       bro.Id,
				OnResponse:   bro.response,
				RequestValue: bro.RequestValue,
				RequestMap:   bro.RequestMap,
			},
			ButtonStyle: bstyle,
			Label:       bro.msg(label),
			Icon:        icoKey,
		}
		if bstyle == ButtonStyleBorder {
			btn.HideLabel = true
			btn.Style = ut.SM{"padding": "8px 12px", "margin": "0 1px"}
		}
		return btn
	}
	ccLnk := func(icoKey, label, hrefURL, download string) *Link {
		return &Link{
			BaseComponent: BaseComponent{
				Id:    bro.Id + "_" + name + "_0",
				Name:  name,
				Style: ut.SM{"padding": "8px 12px", "margin": "0 1px"},
			},
			LinkStyle:  LinkStyleBorder,
			Label:      bro.msg(label),
			Icon:       icoKey,
			HideLabel:  true,
			Href:       hrefURL,
			Download:   download,
			LinkTarget: "_blank",
		}
	}
	ccLbl := func(key, icoKey, label string, class []string) *Label {
		return &Label{
			Value:    label,
			LeftIcon: icoKey,
			BaseComponent: BaseComponent{
				Id:           bro.Id + "_" + name + "_" + key,
				Name:         name,
				Data:         ut.IM{"key": key},
				EventURL:     bro.EventURL,
				Target:       bro.Id,
				OnResponse:   bro.response,
				RequestValue: bro.RequestValue,
				RequestMap:   bro.RequestMap,
				Class:        class,
			},
		}
	}
	ccIcon := func(icoKey, rowKey string, width, height float64) *Icon {
		return &Icon{
			BaseComponent: BaseComponent{
				Id:           bro.Id + "_" + name + "_" + rowKey,
				Name:         name,
				EventURL:     bro.EventURL,
				Target:       bro.Target,
				Data:         data,
				OnResponse:   bro.response,
				RequestValue: bro.RequestValue,
				RequestMap:   bro.RequestMap,
			},
			Value:  icoKey,
			Width:  width,
			Height: height,
		}
	}
	ccMap := map[string]func() ClientComponent{
		"hide_header": func() ClientComponent {
			btn := ccBtn(IconFilter, "browser_view", ButtonStylePrimary, "0")
			for _, view := range bro.Views {
				if view.Value == bro.View {
					btn.Label = view.Text
				}
			}
			btn.Full = true
			btn.Align = TextAlignLeft
			return btn
		},
		"btn_search": func() ClientComponent {
			return ccBtn(IconSearch, "browser_search", ButtonStyleBorder, "0")
		},
		"btn_bookmark": func() ClientComponent {
			return ccBtn(IconStar, "browser_bookmark", ButtonStyleBorder, "0")
		},
		"btn_export": func() ClientComponent {
			if bro.ExportURL != "" {
				return ccLnk(IconDownload, "browser_export", bro.ExportURL, bro.Download)
			}
			btn := ccBtn(IconDownload, "browser_export", ButtonStyleBorder, "0")
			btn.SetProperty("indicator", IndicatorNone)
			return btn
		},
		"btn_help": func() ClientComponent {
			if bro.HelpURL != "" {
				return ccLnk(IconQuestionCircle, "browser_help", bro.HelpURL, "")
			}
			return ccBtn(IconQuestionCircle, "browser_help", ButtonStyleBorder, "0")
		},
		"btn_views": func() ClientComponent {
			btn := ccBtn(IconEye, "browser_views", ButtonStyleBorder, "0")
			btn.Selected = bro.ShowDropdown
			btn.Indicator = IndicatorNone
			return btn
		},
		"btn_columns": func() ClientComponent {
			return ccBtn(IconColumns, "browser_columns", ButtonStyleBorder, "0")
		},
		"btn_filter": func() ClientComponent {
			return ccBtn(IconPlus, "browser_filter", ButtonStyleBorder, "0")
		},
		"btn_total": func() ClientComponent {
			btn := ccBtn(IconInfoCircle, "browser_total", ButtonStyleBorder, "0")
			btn.Disabled = (len(bro.Rows) == 0) || (len(bro.totalFields) == 0)
			return btn
		},
		"btn_ok": func() ClientComponent {
			btn := ccBtn(IconCheck, "browser_label_ok", ButtonStylePrimary, "0")
			btn.AutoFocus = true
			btn.Full = true
			return btn
		},
		"filter_table": func() ClientComponent {
			return bro.filterTable()
		},
		"menu_item": func() ClientComponent {
			key := ut.ToString(data["key"], "")
			label := ut.ToString(data["value"], "")
			icoKey := IconEye
			class := []string{}
			if key == bro.View {
				icoKey = IconCheck
				class = []string{"selected"}
			}
			return ccLbl(key, icoKey, label, class)
		},
		"col_item": func() ClientComponent {
			key := ut.ToString(data["key"], "")
			label := ut.ToString(data["value"], "")
			icoKey := IconSquareEmpty
			class := []string{"edit-col"}
			if ut.ToBoolean(bro.VisibleColumns[key], false) {
				icoKey = IconCheckSquare
				class = []string{"select-col"}
			}
			class = append(class, "base-col")
			return ccLbl(key, icoKey, label, class)
		},
		"total_label": func() ClientComponent {
			return &Label{
				Value: ut.ToString(data["label"], ""),
			}
		},
		"total_value": func() ClientComponent {
			return &NumberInput{
				Value:    ut.ToFloat(data["total"], 0),
				ReadOnly: true,
				Full:     true,
			}
		},
		"edit_row": func() ClientComponent {
			rowKey := ut.ToString(data[bro.RowKey], "")
			if rowKey != "" && !bro.ReadOnly {
				return ccIcon(IconEdit, rowKey, 24, 21.3)
			}
			return &Icon{
				Value: IconCaretRight,
				Width: 9, Height: 24,
			}
		},
		"table": func() ClientComponent {
			return bro.getComponentTable()
		},
	}
	cc := ccMap[name]()
	html, err = cc.Render()
	return html, err
}

func (bro *Browser) msg(labelID string) string {
	if label, found := bro.Labels[labelID]; found {
		return label
	}
	return labelID
}

/*
Based on the values, it will generate the html code of the [Browser] or return with an error message.
*/
func (bro *Browser) Render() (html template.HTML, err error) {
	bro.InitProps(bro)
	if bro.ShowTotal {
		bro.totalFields = bro.setTotalValues()
	}

	funcMap := map[string]any{
		"msg": func(labelID string) string {
			return bro.msg(labelID)
		},
		"styleMap": func() bool {
			return len(bro.Style) > 0
		},
		"showViews": func() bool {
			return len(bro.Views) > 0
		},
		"showFilters": func() bool {
			return len(bro.Filters) > 0
		},
		"customClass": func() string {
			return strings.Join(bro.Class, " ")
		},
		"browserComponent": func(name string) (template.HTML, error) {
			return bro.getComponent(name, ut.IM{})
		},
		"menuItem": func(key, value string) (template.HTML, error) {
			return bro.getComponent("menu_item", ut.IM{"key": key, "value": value})
		},
		"colItem": func(key, value string) (template.HTML, error) {
			return bro.getComponent("col_item", ut.IM{"key": key, "value": value})
		},
		"resultCount": func() int {
			return len(bro.Rows)
		},
		"totalFields": func() []BrowserTotalField {
			return bro.totalFields
		},
		"totalLabel": func(label string) (template.HTML, error) {
			return bro.getComponent("total_label", ut.IM{"label": label})
		},
		"totalValue": func(total float64) (template.HTML, error) {
			return bro.getComponent("total_value", ut.IM{"total": total})
		},
	}
	tpl := `<div id="{{ .Id }}" name="{{ .Name }}" class="row full {{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	><div class="panel">
	<div class="panel-title"><div class="cell title-cell"><span>{{ .Title }}</span></div></div>
	<div class="panel-container" >
	<div class="row full" ><div class="cell" >{{ browserComponent "hide_header" }}</div></div>
	{{ if ne .HideHeader true }}<div class="filter-panel" >
	<div class="row full" >
	<div class="cell" >{{ browserComponent "btn_search" }}</div>
	<div class="cell align-right" >
	{{ if ne .HideBookmark true }}{{ browserComponent "btn_bookmark" }}{{ end }}
	{{ if ne .HideExport true }}{{ browserComponent "btn_export" }}{{ end }}
	{{ if ne .HideHelp true }}{{ browserComponent "btn_help" }}{{ end }}
	</div></div>
	<div class="row full section-small-top" >
	<div class="cell" >
	<div class="dropdown-box" >
	{{ if showViews }}{{ browserComponent "btn_views" }}{{ end }}
	{{ if .ShowDropdown }}<div class="dropdown-content" >
	{{ range $index, $view := .Views }}<div class="drop-label" >{{ menuItem $view.Value $view.Text }}</div>{{ end }}
	</div>{{ end }}
	</div>
	{{ browserComponent "btn_columns" }}{{ browserComponent "btn_filter" }}{{ browserComponent "btn_total" }}
	</div>
	</div>
	{{ if .ShowColumns }}<div class="col-box" >
	{{ range $index, $field := .Table.Fields }}<div 
	class="cell col-cell" >{{ colItem $field.Name $field.Label }}</div>{{ end }}
	</div>{{ end }}
	{{ if showFilters }}<div class="row section-top">
	<div class="row full" style="margin-bottom: 1px;"><div class="cell result-title result-border" >{{ msg "browser_filters" }}</div></div>
	<div class="row full"><div class="cell" >{{ browserComponent "filter_table" }}</div></div>
	</div>{{ end }}
	</div>{{ end }}
	<div class="row full section-small-top" ><div class="row full result-border" >
	<div class="cell result-title" >{{ resultCount }} {{ msg "browser_result" }}</div>
	</div></div>
	<div class="row full" >{{ browserComponent "table" }}</div>
	</div></div>
	{{ if .ShowTotal }}<div class="modal"><div class="dialog"><div class="panel">
	<div class="panel-title">
	<div class="cell title-cell" ><span>{{ msg "browser_total" }}</span></div>
	</div>
	<div class="section" ><div class="row full container" >
	{{ range $index, $row := totalFields }}<div class="trow full">
	<div class="cell padding-tiny mobile">{{ totalLabel $row.Label }}</div>
	<div class="cell padding-tiny mobile">{{ totalValue $row.Total }}</div>
	</div>{{ end }}
	</div></div>
  <div class="section buttons" ><div class="row full container" ><div class="cell padding-small" >
	{{ browserComponent "btn_ok" }}
	</div></div></div>
	</div></div></div>{{ end }}
	</div>`

	return ut.TemplateBuilder("browser", tpl, funcMap, bro)
}

var testBrowserFields map[string]func() []TableField = map[string]func() []TableField{
	"customer": func() []TableField {
		return []TableField{
			{Name: "custnumber", Label: "Customer No."},
			{Name: "custname", Label: "Customer Name", FieldType: TableFieldTypeLink},
			{Name: "taxnumber", Label: "Taxnumber"},
			{Name: "custtype", Label: "Customer Type"},
			{Name: "account", Label: "Account No."},
			{Name: "notax", Label: "Tax-free", FieldType: TableFieldTypeBool},
			{Name: "terms", Label: "Payment per.", FieldType: TableFieldTypeNumber},
			{Name: "creditlimit", Label: "Credit line", FieldType: TableFieldTypeNumber},
			{Name: "discount", Label: "Discount%", FieldType: TableFieldTypeNumber},
			{Name: "notes", Label: "Comment"},
			{Name: "inactive", Label: "Inactive", FieldType: TableFieldTypeBool},
			{Name: "address", Label: "Address"},
		}
	},
	"meta": func() []TableField {
		return []TableField{
			{Name: "custnumber", Label: "Customer No."},
			{Name: "custname", Label: "Customer Name", FieldType: TableFieldTypeLink},
			{Name: "description", Label: "Description"},
			{Name: "deffield", Label: "Value", FieldType: TableFieldTypeMeta},
			{Name: "notes", Label: "Other data"},
		}
	},
	"contact": func() []TableField {
		return []TableField{
			{Name: "custnumber", Label: "Customer No."},
			{Name: "custname", Label: "Customer Name", FieldType: TableFieldTypeLink},
			{Name: "firstname", Label: "Firstname"},
			{Name: "surname", Label: "Surname"},
			{Name: "status", Label: "Status"},
			{Name: "phone", Label: "Phone"},
			{Name: "email", Label: "Email"},
			{Name: "notes", Label: "Comment"},
		}
	},
}

var testBrowserRows map[string]func() []ut.IM = map[string]func() []ut.IM{
	"customer": func() []ut.IM {
		return []ut.IM{
			{"account": nil, "address": "City1 street 1.", "creditlimit": 1000000,
				"custname":   "First Customer Co.",
				"custnumber": "DMCUST/00001", "custtype": "company", "discount": 2,
				"id": "customer-2", "inactive": 0, "notax": 0, "notes": nil,
				"row_id": 2, "taxnumber": "12345678-1-12", "terms": 8},
			{"account": nil, "address": "City3 street 3.", "creditlimit": 10, "custname": "Second Customer Name",
				"custnumber": "DMCUST/00002", "custtype": "private", "discount": 6,
				"id": "customer-3", "inactive": 0, "notax": 0, "notes": nil,
				"row_id": 3, "taxnumber": "12121212-1-12", "terms": 1},
			{
				"account": nil, "address": "City4 street 4.", "creditlimit": 30, "custname": "Third Customer Foundation",
				"custnumber": "DMCUST/00003", "custtype": "other", "discount": 0,
				"id": "customer-4", "inactive": 0, "notax": 1, "notes": nil,
				"row_id": 4, "taxnumber": "10101010-1-01", "terms": 4},
		}
	},
	"meta": func() []ut.IM {
		return []ut.IM{
			{"id": 1, "custnumber": "DMCUST/00001", "custname": "First Customer Co.",
				"deffield_meta": TableFieldTypeNumber,
				"description":   "Customer Float", "deffield": 20.5, "notes": ""},
			{"id": 2, "custnumber": "DMCUST/00001", "custname": "First Customer Co.",
				"deffield_meta": TableFieldTypeDate,
				"description":   "Customer Date", "deffield": "2022-01-01", "notes": "Comment"},
			{"id": 3, "custnumber": "DMCUST/00002", "custname": "Second Customer Name",
				"deffield_meta": TableFieldTypeBool,
				"description":   "Customer Bool", "deffield": true, "notes": ""},
			{"id": 4, "custnumber": "DMCUST/00002", "custname": "Second Customer Name",
				"deffield_meta": TableFieldTypeInteger,
				"description":   "Customer Integer", "deffield": 12345, "notes": ""},
			{"id": 5, "custnumber": "DMCUST/00003", "custname": "Third Customer Foundation",
				"deffield_meta": TableFieldTypeLink,
				"description":   "Customer Product", "deffield": "Big Product", "notes": ""},
		}
	},
	"contact": func() []ut.IM {
		return []ut.IM{
			{"id": 1, "custnumber": "DMCUST/00001", "custname": "First Customer Co.",
				"firstname": "Firstname 1", "surname": "Surname 1", "status": "Status 1",
				"phone": "123456", "email": "email@company.com", "notes": ""},
			{"id": 2, "custnumber": "DMCUST/00001", "custname": "First Customer Co.",
				"firstname": "Firstname 2", "surname": "Surname 2", "status": "Status 2",
				"phone": "123456", "email": "email@company.com", "notes": ""},
			{"id": 3, "custnumber": "DMCUST/00001", "custname": "First Customer Co.",
				"firstname": "Firstname 3", "surname": "Surname 3", "status": "Status 3",
				"phone": "123456", "email": "email@company.com", "notes": ""},
			{"id": 4, "custnumber": "DMCUST/00001", "custname": "First Customer Co.",
				"firstname": "Firstname 4", "surname": "Surname 4", "status": "Status 4",
				"phone": "123456", "email": "email@company.com", "notes": ""},
			{"id": 5, "custnumber": "DMCUST/00001", "custname": "First Customer Co.",
				"firstname": "Firstname 5", "surname": "Surname 5", "status": "Status 5",
				"phone": "123456", "email": "email@company.com", "notes": ""},
			{"id": 6, "custnumber": "DMCUST/00001", "custname": "First Customer Co.",
				"firstname": "Firstname 6", "surname": "Surname 6", "status": "Status 6",
				"phone": "123456", "email": "email@company.com", "notes": ""},
			{"id": 7, "custnumber": "DMCUST/00001", "custname": "First Customer Co.",
				"firstname": "Firstname 7", "surname": "Surname 7", "status": "Status 7",
				"phone": "123456", "email": "email@company.com", "notes": ""},
			{"id": 8, "custnumber": "DMCUST/00001", "custname": "First Customer Co.",
				"firstname": "Firstname 8", "surname": "Surname 8", "status": "Status 8",
				"phone": "123456", "email": "email@company.com", "notes": ""},
			{"id": 9, "custnumber": "DMCUST/00001", "custname": "First Customer Co.",
				"firstname": "Firstname 9", "surname": "Surname 9", "status": "Status 9",
				"phone": "123456", "email": "email@company.com", "notes": ""},
			{"id": 10, "custnumber": "DMCUST/00001", "custname": "First Customer Co.",
				"firstname": "Firstname 10", "surname": "Surname 10", "status": "Status 10",
				"phone": "123456", "email": "email@company.com", "notes": ""},
		}
	},
}

var testBrowserColumns map[string]func() map[string]bool = map[string]func() map[string]bool{
	"customer": func() map[string]bool {
		return map[string]bool{
			"custnumber": true, "custname": true, "address": true,
		}
	},
	"meta": func() map[string]bool {
		return map[string]bool{
			"custname": true, "description": true, "deffield": true,
		}
	},
	"contact": func() map[string]bool {
		return map[string]bool{
			"custname": true, "firstname": true, "surname": true, "phone": true,
		}
	},
}

var testBrowserFilters map[string]func() []BrowserFilter = map[string]func() []BrowserFilter{
	"customer": func() []BrowserFilter {
		return []BrowserFilter{
			{Field: "custname", Comp: "==", Value: "%Customer%"},
			{Field: "creditlimit", Comp: ">=", Value: 5},
			{Field: "inactive", Comp: "==", Value: 0},
		}
	},
	"meta": func() []BrowserFilter {
		return []BrowserFilter{
			{Field: "customer_date", Comp: ">=", Value: "2021-01-01"},
		}
	},
	"contact": func() []BrowserFilter {
		return []BrowserFilter{}
	},
}

var testBrowserMetaFields map[string]func() map[string]BrowserMetaField = map[string]func() map[string]BrowserMetaField{
	"customer": func() map[string]BrowserMetaField {
		return map[string]BrowserMetaField{}
	},
	"meta": func() map[string]BrowserMetaField {
		return map[string]BrowserMetaField{
			"customer_float":   {FieldType: TableFieldTypeNumber, Label: "Customer Float"},
			"customer_date":    {FieldType: TableFieldTypeDate, Label: "Customer Date"},
			"customer_bool":    {FieldType: TableFieldTypeBool, Label: "Customer Bool"},
			"customer_integer": {FieldType: TableFieldTypeInteger, Label: "Customer Integer"},
			"customer_product": {FieldType: TableFieldTypeLink, Label: "Customer Product"},
		}
	},
	"contact": func() map[string]BrowserMetaField {
		return map[string]BrowserMetaField{}
	},
}

var testBrowserResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
	switch evt.Name {
	case BrowserEventSearch, BrowserEventBookmark, BrowserEventHelp,
		TableEventAddItem, TableEventEditCell, TableEventRowSelected, BrowserEventEditRow:
		re = ResponseEvent{
			Trigger: &Toast{
				Type:    ToastTypeInfo,
				Value:   ut.ToString(evt.Value, evt.Name),
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

	case BrowserEventView:
		evt.Trigger.SetProperty("view", evt.Value)
		evt.Trigger.SetProperty("fields", testBrowserFields[ut.ToString(evt.Value, "")]())
		evt.Trigger.SetProperty("rows", testBrowserRows[ut.ToString(evt.Value, "")]())
		evt.Trigger.SetProperty("visible_columns", testBrowserColumns[ut.ToString(evt.Value, "")]())
		evt.Trigger.SetProperty("filters", testBrowserFilters[ut.ToString(evt.Value, "")]())
		evt.Trigger.SetProperty("meta_fields", testBrowserMetaFields[ut.ToString(evt.Value, "")]())
	}

	return evt
}

// [Browser] test and demo data
func TestBrowser(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeBrowser,
			Component: &Browser{
				Table: Table{
					BaseComponent: BaseComponent{
						Id:           id + "default",
						EventURL:     eventURL,
						OnResponse:   testBrowserResponse,
						RequestValue: requestValue,
						RequestMap:   requestMap,
					},
					Fields:            testBrowserFields["customer"](),
					Rows:              testBrowserRows["customer"](),
					TableFilter:       true,
					AddItem:           true,
					HidePaginatonSize: false,
					PageSize:          5,
				},
				View:  "customer",
				Title: "Customer Data",
				Views: []SelectOption{
					{Value: "customer", Text: "Customer Data"},
					{Value: "meta", Text: "Metadata"},
					{Value: "contact", Text: "Contact Info"},
				},
				VisibleColumns: testBrowserColumns["customer"](),
				Filters:        testBrowserFilters["customer"](),
				MetaFields:     testBrowserMetaFields["customer"](),
				HideFilters:    map[string]bool{"status": true},
			}},
		{
			Label:         "Meta data",
			ComponentType: ComponentTypeBrowser,
			Component: &Browser{
				Table: Table{
					BaseComponent: BaseComponent{
						Id:           id + "meta",
						EventURL:     eventURL,
						OnResponse:   testBrowserResponse,
						RequestValue: requestValue,
						RequestMap:   requestMap,
					},
					Fields:            testBrowserFields["meta"](),
					Rows:              testBrowserRows["meta"](),
					TableFilter:       true,
					AddItem:           false,
					HidePaginatonSize: false,
					PageSize:          5,
				},
				View: "meta",
				Views: []SelectOption{
					{Value: "customer", Text: "Customer Data"},
					{Value: "meta", Text: "Metadata"},
					{Value: "contact", Text: "Contact Info"},
				},
				VisibleColumns: testBrowserColumns["meta"](),
				Filters:        testBrowserFilters["meta"](),
				MetaFields:     testBrowserMetaFields["meta"](),
				ReadOnly:       true,
				ExportURL:      "/export",
				Download:       "export.csv",
				HelpURL:        "https://www.google.com",
			}},
		{
			Label:         "Contact data",
			ComponentType: ComponentTypeBrowser,
			Component: &Browser{
				Table: Table{
					BaseComponent: BaseComponent{
						Id:           id + "contact",
						EventURL:     eventURL,
						OnResponse:   testBrowserResponse,
						RequestValue: requestValue,
						RequestMap:   requestMap,
					},
					Fields:            testBrowserFields["contact"](),
					Rows:              testBrowserRows["contact"](),
					TableFilter:       true,
					AddItem:           true,
					HidePaginatonSize: false,
					PageSize:          5,
				},
				View: "contact",
				Views: []SelectOption{
					{Value: "customer", Text: "Customer Data"},
					{Value: "meta", Text: "Metadata"},
					{Value: "contact", Text: "Contact Info"},
				},
				VisibleColumns: testBrowserColumns["contact"](),
				Filters:        testBrowserFilters["contact"](),
				MetaFields:     testBrowserMetaFields["contact"](),
				ShowColumns:    true,
				ShowDropdown:   true,
				HideBookmark:   true,
				HideHelp:       true,
				HideExport:     true,
			}},
		{
			Label:         "Total",
			ComponentType: ComponentTypeBrowser,
			Component: &Browser{
				Table: Table{
					BaseComponent: BaseComponent{
						Id:           id + "total",
						EventURL:     eventURL,
						OnResponse:   testBrowserResponse,
						RequestValue: requestValue,
						RequestMap:   requestMap,
					},
					Fields:            testBrowserFields["customer"](),
					Rows:              testBrowserRows["customer"](),
					TableFilter:       true,
					AddItem:           true,
					HidePaginatonSize: false,
					PageSize:          5,
				},
				View:           "customer",
				VisibleColumns: testBrowserColumns["customer"](),
				Filters:        testBrowserFilters["customer"](),
				MetaFields:     testBrowserMetaFields["customer"](),
				ShowTotal:      true,
				HideHeader:     true,
			}},
	}
}
