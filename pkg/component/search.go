package component

import (
	"html/template"
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

const (
	ComponentTypeSearch = "search"

	SearchEventSearch   = "search"
	SearchEventSelected = "selected"

	SearchDefaultPlaceholder = "Search conditions"
	SearchDefaultTitle       = "Search for data"
)

// Creates an interactive simple search control
type Search struct {
	BaseComponent
	// Data source of the table
	Rows []ut.IM `json:"rows"`
	// Table column definitions
	Fields []TableField `json:"fields"`
	// Search window title
	Title string `json:"title"`
	// Specifies a short hint that describes the expected value of the input element
	FilterPlaceholder string `json:"filter_placeholder"`
	// Filter input value
	FilterValue string `json:"filter_value"`
	// Specifies that the input element should automatically get focus when the page loads
	AutoFocus bool `json:"auto_focus"`
	// Full width input (100%)
	Full bool `json:"full"`
}

/*
Returns all properties of the [Search]
*/
func (sea *Search) Properties() ut.IM {
	return ut.MergeIM(
		sea.BaseComponent.Properties(),
		ut.IM{
			"rows":               sea.Rows,
			"fields":             sea.Fields,
			"title":              sea.Title,
			"filter_placeholder": sea.FilterPlaceholder,
			"filter_value":       sea.FilterValue,
			"auto_focus":         sea.AutoFocus,
			"full":               sea.Full,
		})
}

/*
Returns the value of the property of the [Search] with the specified name.
*/
func (sea *Search) GetProperty(propName string) interface{} {
	return sea.Properties()[propName]
}

/*
It checks the value given to the property of the [Search] and always returns a valid value
*/
func (sea *Search) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"rows": func() interface{} {
			return ut.ToIMA(propValue, []ut.IM{})
		},
		"fields": func() interface{} {
			fields := []TableField{}
			if fd, valid := propValue.([]TableField); valid && (fd != nil) {
				fields = fd
			}
			if len(fields) == 0 {
				if len(sea.Rows) > 0 {
					for field := range sea.Rows[0] {
						fields = append(fields,
							TableField{Name: field, FieldType: TableFieldTypeString, Label: field})
					}
				}
			}
			return fields
		},
		"target": func() interface{} {
			sea.SetProperty("id", sea.Id)
			value := ut.ToString(propValue, sea.Id)
			if value != "this" && !strings.HasPrefix(value, "#") {
				value = "#" + value
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if sea.BaseComponent.GetProperty(propName) != nil {
		return sea.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [Search] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (sea *Search) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"rows": func() interface{} {
			sea.Rows = sea.Validation(propName, propValue).([]ut.IM)
			return sea.Rows
		},
		"fields": func() interface{} {
			sea.Fields = sea.Validation(propName, propValue).([]TableField)
			return sea.Fields
		},
		"title": func() interface{} {
			sea.Title = ut.ToString(propValue, SearchDefaultTitle)
			return sea.Title
		},
		"filter_placeholder": func() interface{} {
			sea.FilterPlaceholder = ut.ToString(propValue, "")
			return sea.FilterPlaceholder
		},
		"filter_value": func() interface{} {
			sea.FilterValue = ut.ToString(propValue, "")
			return sea.FilterValue
		},
		"auto_focus": func() interface{} {
			sea.AutoFocus = ut.ToBoolean(propValue, false)
			return sea.AutoFocus
		},
		"full": func() interface{} {
			sea.Full = ut.ToBoolean(propValue, false)
			return sea.Full
		},
		"target": func() interface{} {
			sea.Target = sea.Validation(propName, propValue).(string)
			return sea.Target
		},
	}
	if _, found := pm[propName]; found {
		return sea.SetRequestValue(propName, pm[propName](), []string{})
	}
	if sea.BaseComponent.GetProperty(propName) != nil {
		return sea.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (sea *Search) response(evt ResponseEvent) (re ResponseEvent) {
	selEvt := ResponseEvent{Trigger: sea, TriggerName: sea.Name}
	switch evt.TriggerName {
	case "search_result":
		if evt.Name != TableEventRowSelected {
			return evt
		}
		selEvt.Name = SearchEventSelected
		selEvt.Value = evt.Value
		selEvt.Header = ut.SM{
			HeaderRetarget: "#" + sea.Id,
		}

	case "filter":
		selEvt.Name = SearchEventSearch
		sea.SetProperty("filter_value", evt.Value)

	case "btn_search":
		selEvt.Name = SearchEventSearch

	default:
	}
	if sea.OnResponse != nil {
		return sea.OnResponse(selEvt)
	}
	return selEvt
}

func (sea *Search) getComponent(name string) (html template.HTML, err error) {
	ccBtn := func(icon string) *Button {
		return &Button{
			BaseComponent: BaseComponent{
				Id:           sea.Id + "_" + name,
				Name:         name,
				Style:        ut.SM{"padding": "8px", "margin": "1px 0 2px 1px"},
				EventURL:     sea.EventURL,
				Target:       sea.Target,
				OnResponse:   sea.response,
				RequestValue: sea.RequestValue,
				RequestMap:   sea.RequestMap,
			},
			ButtonStyle: ButtonStyleBorder,
			Icon:        icon,
		}
	}
	ccMap := map[string]func() ClientComponent{
		"btn_search": func() ClientComponent {
			return ccBtn("Search")
		},
		"filter": func() ClientComponent {
			return &Input{
				BaseComponent: BaseComponent{
					Id: sea.Id + "_" + name, Name: name,
					Style:        ut.SM{"border-radius": "0", "margin": "1px 0 2px"},
					EventURL:     sea.EventURL,
					Target:       sea.Target,
					OnResponse:   sea.response,
					RequestValue: sea.RequestValue,
					RequestMap:   sea.RequestMap,
				},
				Type:        InputTypeString,
				Label:       sea.FilterPlaceholder,
				Placeholder: sea.FilterPlaceholder,
				Value:       sea.FilterValue,
				Full:        true,
				AutoFocus:   sea.AutoFocus,
			}
		},
		"search_result": func() ClientComponent {
			return &Table{
				BaseComponent: BaseComponent{
					Id:           sea.Id + "_" + name,
					Name:         name,
					EventURL:     sea.EventURL,
					OnResponse:   sea.response,
					RequestValue: sea.RequestValue,
					RequestMap:   sea.RequestMap,
				},
				Rows:              sea.Rows,
				Fields:            sea.Fields,
				Pagination:        PaginationTypeTop,
				PageSize:          5,
				HidePaginatonSize: true,
				TableFilter:       false,
				AddItem:           false,
				RowSelected:       true,
			}
		},
	}
	cc := ccMap[name]()
	html, err = cc.Render()
	return html, err
}

/*
Based on the values, it will generate the html code of the [Search] or return with an error message.
*/
func (sea *Search) Render() (html template.HTML, err error) {
	sea.InitProps(sea)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(sea.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(sea.Class, " ")
		},
		"searchComponent": func(name string) (template.HTML, error) {
			return sea.getComponent(name)
		},
	}
	tpl := `<div id="{{ .Id }}" name="{{ .Name }}" class="row full {{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	><div class="panel"><div class="panel-title">
	<div class="cell title-cell"><span>{{ .Title }}</span></div></div>
	<div class="section" >
	<div class="row full container" >
	<div class="cell">{{ searchComponent "filter" }}</div>
	<div class="cell" style="width: 20px;" >{{ searchComponent "btn_search" }}</div></div>
	<div class="row full container" >{{ searchComponent "search_result" }}</div>
	</div></div></div>`

	return ut.TemplateBuilder("search", tpl, funcMap, sea)
}

var testSearchResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
	toast := func(value string) ResponseEvent {
		return ResponseEvent{
			Trigger: &Toast{
				Type:    ToastTypeInfo,
				Value:   value,
				Timeout: 4,
			},
			TriggerName: evt.TriggerName,
			Name:        evt.Name,
			Header: ut.SM{
				HeaderRetarget: "#toast-msg",
				HeaderReswap:   SwapInnerHTML,
			},
		}
	}
	if evt.Name == SearchEventSearch {
		return toast(ut.ToString(evt.Trigger.GetProperty("filter_value"), ""))
	}

	row := evt.Value.(ut.IM)["row"].(ut.IM)
	return toast(ut.ToString(row["custname"], ""))
}

var testSearchFields []TableField = []TableField{
	{Name: "custname", Label: "Customer Name"},
	{Name: "custnumber", Label: "Customer Number"},
	{Name: "city", Label: "City"},
	{Name: "street", Label: "Street"},
}

var testSearchRows []ut.IM = []ut.IM{
	{
		"city":       "City1",
		"custname":   "First Customer Co.",
		"custnumber": "DMCUST/00001",
		"id":         "customer-2",
		"label":      "First Customer Co.",
		"street":     "street 1.",
	},
	{
		"city":       "City3",
		"custname":   "Second Customer Name",
		"custnumber": "DMCUST/00002",
		"id":         "customer-3",
		"label":      "Second Customer Name",
		"street":     "street 3.",
		"deleted":    1,
	},
	{
		"city":       "City4",
		"custname":   "Third Customer Foundation",
		"custnumber": "DMCUST/00003",
		"id":         "customer-4",
		"label":      "Third Customer Foundation",
		"street":     "street 4.",
	},
}

// [Search] test and demo data
func TestSearch(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeList,
			Component: &Search{
				BaseComponent: BaseComponent{
					Id:           id + "_search_default",
					EventURL:     eventURL,
					OnResponse:   testSearchResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Fields:    testSearchFields,
				Rows:      testSearchRows,
				AutoFocus: true,
			},
		},
	}
}
