package component

import (
	"html/template"
	"math"
	"slices"
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [List] constants
const (
	ComponentTypeList = "list"

	ListEventCurrentPage  = "current_page"
	ListEventFilterChange = "filter_change"
	ListEventAddItem      = "add_item"
	ListEventEditItem     = "edit_item"
	ListEventDelete       = "delete_item"
)

/*
Creates an interactive list control
*/
type List struct {
	BaseComponent
	// Data source of the list
	Rows []ut.IM `json:"rows"`
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
	// Show/hide list value filter input row
	ListFilter bool `json:"list_filter"`
	// Show/hide list add item button
	AddItem bool `json:"add_item"`
	// Enable edit event
	EditItem bool `json:"edit_item"`
	// Enable delete event
	DeleteItem bool `json:"delete_item"`
	// Specifies a short hint that describes the expected value of the input element
	FilterPlaceholder string `json:"filter_placeholder"`
	// Filter input value
	FilterValue string `json:"filter_value"`
	// The filter is case sensitive
	CaseSensitive bool `json:"case_sensitive"`
	// Add item button caption Default empty string
	LabelAdd string `json:"label_add"`
	// Valid [Icon] component value. See more [IconKey] variable values.
	AddIcon string `json:"add_icon"`
	// Valid [Icon] component value. See more [IconKey] variable values.
	EditIcon string `json:"edit_icon"`
	// Valid [Icon] component value. See more [IconKey] variable values.
	DeleteIcon string `json:"delete_icon"`
	// The field name containing the list label of the data source. Default: lslabel
	LabelField string `json:"label_field"`
	// The field name containing the list value of the data source. Default: lsvalue
	LabelValue string `json:"label_value"`
}

/*
Returns all properties of the [List]
*/
func (lst *List) Properties() ut.IM {
	return ut.MergeIM(
		lst.BaseComponent.Properties(),
		ut.IM{
			"rows":                lst.Rows,
			"pagination":          lst.Pagination,
			"current_page":        lst.CurrentPage,
			"page_size":           lst.PageSize,
			"hide_paginaton_size": lst.HidePaginatonSize,
			"list_filter":         lst.ListFilter,
			"add_item":            lst.AddItem,
			"edit_item":           lst.EditItem,
			"delete_item":         lst.DeleteItem,
			"filter_placeholder":  lst.FilterPlaceholder,
			"filter_value":        lst.FilterValue,
			"case_sensitive":      lst.CaseSensitive,
			"label_add":           lst.LabelAdd,
			"add_icon":            lst.AddIcon,
			"edit_icon":           lst.EditIcon,
			"delete_icon":         lst.DeleteIcon,
			"label_field":         lst.LabelField,
			"label_value":         lst.LabelValue,
		})
}

/*
Returns the value of the property of the [List] with the specified name.
*/
func (lst *List) GetProperty(propName string) interface{} {
	return lst.Properties()[propName]
}

/*
It checks the value given to the property of the [List] and always returns a valid value
*/
func (lst *List) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"rows": func() interface{} {
			return ut.ToIMA(propValue, []ut.IM{})
		},
		"pagination": func() interface{} {
			return lst.CheckEnumValue(ut.ToString(propValue, ""), PaginationTypeTop, PaginationType)
		},
		"current_page": func() interface{} {
			value := ut.ToInteger(propValue, 1)
			rows := lst.filterRows()
			pageCount := int64(math.Ceil(float64(len(rows)) / float64(lst.PageSize)))
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
			lst.SetProperty("id", lst.Id)
			value := ut.ToString(propValue, lst.Id)
			if value != "this" && !strings.HasPrefix(value, "#") {
				value = "#" + value
			}
			return value
		},
		"label_field": func() interface{} {
			return ut.ToString(propValue, "lslabel")
		},
		"label_value": func() interface{} {
			return ut.ToString(propValue, "lsvalue")
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if lst.BaseComponent.GetProperty(propName) != nil {
		return lst.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [List] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (lst *List) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"rows": func() interface{} {
			lst.Rows = lst.Validation(propName, propValue).([]ut.IM)
			return lst.Rows
		},
		"pagination": func() interface{} {
			lst.Pagination = lst.Validation(propName, propValue).(string)
			return lst.Pagination
		},
		"current_page": func() interface{} {
			lst.CurrentPage = lst.Validation(propName, propValue).(int64)
			return lst.CurrentPage
		},
		"page_size": func() interface{} {
			lst.PageSize = lst.Validation(propName, propValue).(int64)
			return lst.PageSize
		},
		"hide_paginaton_size": func() interface{} {
			lst.HidePaginatonSize = ut.ToBoolean(propValue, false)
			return lst.HidePaginatonSize
		},
		"list_filter": func() interface{} {
			lst.ListFilter = ut.ToBoolean(propValue, false)
			return lst.ListFilter
		},
		"add_item": func() interface{} {
			lst.AddItem = ut.ToBoolean(propValue, false)
			return lst.AddItem
		},
		"edit_item": func() interface{} {
			lst.EditItem = ut.ToBoolean(propValue, false)
			return lst.EditItem
		},
		"delete_item": func() interface{} {
			lst.DeleteItem = ut.ToBoolean(propValue, false)
			return lst.DeleteItem
		},
		"filter_placeholder": func() interface{} {
			lst.FilterPlaceholder = ut.ToString(propValue, "")
			return lst.FilterPlaceholder
		},
		"filter_value": func() interface{} {
			lst.FilterValue = ut.ToString(propValue, "")
			return lst.FilterValue
		},
		"case_sensitive": func() interface{} {
			lst.CaseSensitive = ut.ToBoolean(propValue, false)
			return lst.CaseSensitive
		},
		"label_add": func() interface{} {
			lst.LabelAdd = ut.ToString(propValue, "")
			return lst.LabelAdd
		},
		"add_icon": func() interface{} {
			lst.AddIcon = ut.ToString(propValue, "Plus")
			return lst.AddIcon
		},
		"edit_icon": func() interface{} {
			lst.EditIcon = ut.ToString(propValue, "Edit")
			return lst.EditIcon
		},
		"delete_icon": func() interface{} {
			lst.DeleteIcon = ut.ToString(propValue, "Times")
			return lst.DeleteIcon
		},
		"target": func() interface{} {
			lst.Target = lst.Validation(propName, propValue).(string)
			return lst.Target
		},
		"label_field": func() interface{} {
			lst.LabelField = lst.Validation(propName, propValue).(string)
			return lst.LabelField
		},
		"label_value": func() interface{} {
			lst.LabelValue = lst.Validation(propName, propValue).(string)
			return lst.LabelValue
		},
	}
	if _, found := pm[propName]; found {
		return lst.SetRequestValue(propName, pm[propName](), []string{})
	}
	if lst.BaseComponent.GetProperty(propName) != nil {
		return lst.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (lst *List) response(evt ResponseEvent) (re ResponseEvent) {
	lstEvt := ResponseEvent{
		Trigger: lst, TriggerName: lst.Name, Value: evt.Value,
	}
	switch evt.TriggerName {
	case "top_pagination", "bottom_pagination":
		if evt.Name == PaginationEventPageSize {
			lst.SetProperty("page_size", lstEvt.Value)
			lst.SetProperty("current_page", 1)
		} else {
			lstEvt.Name = ListEventCurrentPage
			lst.SetProperty("current_page", lstEvt.Value)
		}

	case "filter":
		lstEvt.Name = ListEventFilterChange
		lst.SetProperty("filter_value", lstEvt.Value)

	case "btn_add":
		lstEvt.Name = ListEventAddItem

	case "edit_item":
		lstEvt.Name = ListEventEditItem
		lstEvt.Value = evt.Trigger.GetProperty("data")

	case "delete_item":
		lstEvt.Name = ListEventDelete
		lstEvt.Value = evt.Trigger.GetProperty("data")

	default:
	}
	if lst.OnResponse != nil {
		return lst.OnResponse(lstEvt)
	}
	return lstEvt
}

func (lst *List) getComponent(name string, pageCount int64) (html template.HTML, err error) {
	ccPgn := func() *Pagination {
		return &Pagination{
			BaseComponent: BaseComponent{
				Id: lst.Id + "_" + name, Name: name,
				EventURL:     lst.EventURL,
				Target:       lst.Target,
				OnResponse:   lst.response,
				RequestValue: lst.RequestValue,
				RequestMap:   lst.RequestMap,
			},
			Value: lst.CurrentPage, PageSize: lst.PageSize,
			PageCount:    pageCount,
			HidePageSize: lst.HidePaginatonSize,
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
					Id: lst.Id + "_" + name, Name: name,
					Style:        ut.SM{"border-radius": "0", "margin": "1px 0 2px"},
					EventURL:     lst.EventURL,
					Target:       lst.Target,
					OnResponse:   lst.response,
					RequestValue: lst.RequestValue,
					RequestMap:   lst.RequestMap,
				},
				Type:        InputTypeString,
				Label:       lst.FilterPlaceholder,
				Placeholder: lst.FilterPlaceholder,
				Value:       lst.FilterValue,
				Full:        true,
			}
		},
		"btn_add": func() ClientComponent {
			return &Button{
				BaseComponent: BaseComponent{
					Id: lst.Id + "_" + name, Name: name,
					Style:        ut.SM{"padding": "8px 16px", "border-radius": "0", "margin": "1px 0 2px 1px"},
					EventURL:     lst.EventURL,
					Target:       lst.Target,
					OnResponse:   lst.response,
					RequestValue: lst.RequestValue,
					RequestMap:   lst.RequestMap,
				},
				ButtonStyle: ButtonStyleBorder,
				Icon:        lst.AddIcon, Label: lst.LabelAdd,
			}
		},
		"edit_icon": func() ClientComponent {
			return &Icon{
				Value: lst.EditIcon,
			}
		},
		"delete_icon": func() ClientComponent {
			return &Icon{
				Value: lst.DeleteIcon,
			}
		},
	}
	cc := ccMap[name]()
	html, err = cc.Render()
	return html, err
}

func (lst *List) filterRows() (rows []ut.IM) {
	rows = []ut.IM{}
	caseValue := func(value string) string {
		if !lst.CaseSensitive {
			return strings.ToLower(value)
		}
		return value
	}
	filterFields := func() (fields []string) {
		fields = []string{}
		if lst.LabelField != "" {
			fields = append(fields, lst.LabelField)
		}
		if lst.LabelValue != "" {
			fields = append(fields, lst.LabelValue)
		}
		return fields
	}
	getValidRow := func(row ut.IM, fields []string, filter string) bool {
		for _, field := range fields {
			if strings.Contains(caseValue(ut.ToString(row[field], "")), filter) {
				return true
			}
		}
		return false
	}
	if lst.FilterValue == "" {
		return lst.Rows
	}
	resFields := filterFields()
	for _, row := range lst.Rows {
		if getValidRow(row, resFields, caseValue(lst.FilterValue)) {
			rows = append(rows, row)
		}
	}
	return rows
}

/*
Based on the values, it will generate the html code of the [List] or return with an error message.
*/
func (lst *List) Render() (html template.HTML, err error) {
	lst.InitProps(lst)

	rows := lst.filterRows()
	pageCount := int64(math.Ceil(float64(len(rows)) / float64(lst.PageSize)))

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(lst.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(lst.Class, " ")
		},
		"topPagination": func() bool {
			return ((pageCount > 1) && ((lst.Pagination == PaginationTypeTop) || lst.Pagination == PaginationTypeAll))
		},
		"bottomPagination": func() bool {
			return ((pageCount > 1) && ((lst.Pagination == PaginationTypeBottom) || lst.Pagination == PaginationTypeAll))
		},
		"listComponent": func(name string) (template.HTML, error) {
			return lst.getComponent(name, pageCount)
		},
		"listRows": func() []ut.IM {
			if lst.Pagination != PaginationTypeNone {
				currentPage := lst.Validation("current_page", lst.CurrentPage).(int64)
				start := (currentPage - 1) * lst.PageSize
				end := currentPage * lst.PageSize
				if end > int64(len(rows)) {
					end = int64(len(rows))
				}
				return rows[start:end]
			}
			return rows
		},
		"rowID": func(row ut.IM, index int, event string) string {
			rowID := lst.Id + "_row_" + event + "_" + ut.ToString(index, "")
			lbl := &Label{BaseComponent: BaseComponent{
				Id: rowID, Name: event, Data: ut.IM{
					"row": row, "index": index,
				},
				OnResponse:   lst.response,
				RequestValue: lst.RequestValue,
				RequestMap:   lst.RequestMap,
			}}
			lbl.SetProperty("request_map", lbl)
			return rowID
		},
		"isValue": func(row ut.IM, fieldName string) bool {
			return (ut.ToString(row[fieldName], "") != "")
		},
		"rowValue": func(row ut.IM, fieldName string) string {
			return ut.ToString(row[fieldName], "")
		},
	}
	tpl := `<div id="{{ .Id }}" name="{{ .Name }}" class="responsive {{ customClass }}">
	{{ if or .ListFilter topPagination }}<div>
	{{ if topPagination }}<div>{{ listComponent "top_pagination" }}</div>{{ end }}
	{{ if .ListFilter }}<div class="row full">
	<div class="cell" >{{ listComponent "filter" }}</div>
	{{ if .AddItem }}<div class="cell" style="width: 20px;" >{{ listComponent "btn_add" }}</div>{{ end }}
	</div>{{ end }}</div>{{ end }}
	<ul class="list"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}>
	{{ range $index, $row := listRows }}
	<li class="list-row border-bottom">
	{{ if $.EditItem }}<div id="{{ rowID $row $index "edit_item" }}" class="list-edit-cell" 
	{{ if ne $.EventURL "" }} hx-post="{{ $.EventURL }}" hx-target="{{ $.Target }}" hx-swap="{{ $.Swap }}"{{ end }}
	{{ if ne $.Indicator "none" }} hx-indicator="#{{ $.Indicator }}"{{ end }}
	>{{ listComponent "edit_icon" }}</div>{{ end }}
	<div id="{{ rowID $row $index "edit_item" }}" class="list-value-cell {{ if $.EditItem }} cursor-pointer{{ end }}" 
	{{ if ne $.EventURL "" }} hx-post="{{ $.EventURL }}" hx-target="{{ $.Target }}" hx-swap="{{ $.Swap }}"{{ end }}
	{{ if ne $.Indicator "none" }} hx-indicator="#{{ $.Indicator }}"{{ end }}
	>
	{{ if isValue $row $.LabelField }}<div class="border-bottom list-label" ><span>{{ rowValue $row $.LabelField }}</span></div>{{ end }}
  {{ if isValue $row $.LabelValue }}<div class="list-value" ><span>{{ rowValue $row $.LabelValue }}</span></div>{{ end }}
	</div>
	{{ if $.DeleteItem }}<div id="{{ rowID $row $index "delete_item" }}" class="list-delete-cell" 
	{{ if ne $.EventURL "" }} hx-post="{{ $.EventURL }}" hx-target="{{ $.Target }}" hx-swap="{{ $.Swap }}"{{ end }}
	{{ if ne $.Indicator "none" }} hx-indicator="#{{ $.Indicator }}"{{ end }}
	>{{ listComponent "delete_icon" }}</div>{{ end }}
	</li>
	{{ end }}
	</ul>
	{{ if bottomPagination }}<div>{{ listComponent "bottom_pagination" }}</div>{{ end }}
	</div>`

	return ut.TemplateBuilder("list", tpl, funcMap, lst)
}

var testListRows []ut.IM = []ut.IM{
	{"lslabel": "Label 1", "lsvalue": "Value row 1"},
	{"lslabel": "Label 2", "lsvalue": "Value Row 2"},
	{"lslabel": "", "lsvalue": "Value row 3"},
	{"lslabel": "", "lsvalue": "Value row 6"},
	{"lslabel": "Label 5", "lsvalue": "Value row 6"},
	{"lslabel": "Label 6", "lsvalue": "Value row 6"},
	{"lslabel": "Label 7", "lsvalue": "Value Row 7"},
	{"lslabel": "Label 8", "lsvalue": "Value row 8"},
	{"lslabel": "Label 9", "lsvalue": "Value row 9"},
}

var testListResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
	toastType := ut.SM{
		ListEventAddItem:  ToastTypeSuccess,
		ListEventEditItem: ToastTypeInfo,
		ListEventDelete:   ToastTypeError,
	}
	switch evt.Name {
	case ListEventAddItem, ListEventEditItem, ListEventDelete:
		re = ResponseEvent{
			Trigger: &Toast{
				Type:    toastType[evt.Name],
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

// [List] test and demo data
func TestList(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeList,
			Component: &List{
				BaseComponent: BaseComponent{
					Id: id + "_list_default",
				},
				Rows:       testListRows,
				Pagination: PaginationTypeNone,
				PageSize:   5,
			}},
		{
			Label:         "Top Pagination",
			ComponentType: ComponentTypeList,
			Component: &List{
				BaseComponent: BaseComponent{
					Id:           id + "_list_top_pagination",
					EventURL:     eventURL,
					OnResponse:   testListResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Indicator:    IndicatorSpinner,
				},
				Rows:        testListRows,
				Pagination:  PaginationTypeTop,
				PageSize:    5,
				CurrentPage: 10,
				AddItem:     true,
				EditItem:    true,
				DeleteItem:  true,
			}},
		{
			Label:         "Bottom Pagination",
			ComponentType: ComponentTypeList,
			Component: &List{
				BaseComponent: BaseComponent{
					Id:           id + "_list_bottom_pagination",
					EventURL:     eventURL,
					OnResponse:   testListResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Rows:              testListRows,
				Pagination:        PaginationTypeBottom,
				PageSize:          5,
				CurrentPage:       1,
				HidePaginatonSize: true,
				ListFilter:        true,
				FilterPlaceholder: "Placeholder text",
				AddIcon:           "User",
				EditIcon:          "Check",
				DeleteIcon:        "Close",
				AddItem:           true,
				EditItem:          true,
				DeleteItem:        true,
			}},
		{
			Label:         "Filtered",
			ComponentType: ComponentTypeList,
			Component: &List{
				BaseComponent: BaseComponent{
					Id:           id + "_list_filtered",
					EventURL:     eventURL,
					OnResponse:   testListResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Rows:          testListRows,
				Pagination:    PaginationTypeTop,
				PageSize:      5,
				ListFilter:    true,
				FilterValue:   "6",
				CaseSensitive: false,
				LabelAdd:      "Add new",
				AddItem:       true,
				EditItem:      true,
				DeleteItem:    true,
			}},
		{
			Label:         "Filtered CaseSensitive",
			ComponentType: ComponentTypeList,
			Component: &List{
				BaseComponent: BaseComponent{
					Id:           id + "_list_filtered",
					EventURL:     eventURL,
					OnResponse:   testListResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Rows:          testListRows,
				Pagination:    PaginationTypeTop,
				PageSize:      5,
				ListFilter:    true,
				FilterValue:   "Row",
				CaseSensitive: true,
				LabelAdd:      "Add new",
				AddItem:       true,
				EditItem:      true,
				DeleteItem:    true,
			}},
	}
}
