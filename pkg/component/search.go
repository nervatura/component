package component

import (
	"html/template"
	"slices"
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

const (
	ComponentTypeSearch = "search"

	SearchEventSearch   = "search_search"
	SearchEventSelected = "search_selected"
	SearchEventHelp     = "search_help"

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
	// Pagination component [PageSize] variable constants: 5, 10, 20, 50, 100. Default value: 10
	PageSize int64 `json:"page_size"`
	// [Pagination] component show/hide page size selector
	HidePaginatonSize bool `json:"hide_paginaton_size"`
	// Specifies a short hint that describes the expected value of the input element
	FilterPlaceholder string `json:"filter_placeholder"`
	// Specifies that the input element should automatically get focus when the page loads
	AutoFocus bool `json:"auto_focus"`
	// Full width input (100%)
	Full bool `json:"full"`
	// Show or hide the help button
	ShowHelp bool `json:"hide_help"`
	// Specifies the url for help. If it is not specified, then the built-in button event
	HelpURL string `json:"help_url"`
}

/*
Returns all properties of the [Search]
*/
func (sea *Search) Properties() ut.IM {
	return ut.MergeIM(
		sea.BaseComponent.Properties(),
		ut.IM{
			"rows":                sea.Rows,
			"fields":              sea.Fields,
			"title":               sea.Title,
			"page_size":           sea.PageSize,
			"hide_paginaton_size": sea.HidePaginatonSize,
			"filter_placeholder":  sea.FilterPlaceholder,
			"auto_focus":          sea.AutoFocus,
			"full":                sea.Full,
			"show_help":           sea.ShowHelp,
			"help_url":            sea.HelpURL,
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
		"page_size": func() interface{} {
			sea.PageSize = sea.Validation(propName, propValue).(int64)
			return sea.PageSize
		},
		"hide_paginaton_size": func() interface{} {
			sea.HidePaginatonSize = ut.ToBoolean(propValue, false)
			return sea.HidePaginatonSize
		},
		"filter_placeholder": func() interface{} {
			sea.FilterPlaceholder = ut.ToString(propValue, "")
			return sea.FilterPlaceholder
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
		"show_help": func() interface{} {
			sea.ShowHelp = ut.ToBoolean(propValue, false)
			return sea.ShowHelp
		},
		"help_url": func() interface{} {
			sea.HelpURL = ut.ToString(propValue, "")
			return sea.HelpURL
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

	case "filter_value":
		selEvt.Name = SearchEventSearch
		selEvt.Value = evt.Value
		sea.SetProperty("data", ut.IM{evt.TriggerName: evt.Value})

	case "btn_search":
		selEvt.Value = ut.ToString(sea.Data["filter_value"], "")
		selEvt.Name = SearchEventSearch

	case "btn_help":
		selEvt.Name = SearchEventHelp

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
	ccLnk := func() *Link {
		return &Link{
			BaseComponent: BaseComponent{
				Id:    sea.Id + "_" + name,
				Name:  name,
				Style: ut.SM{"padding": "8px", "margin": "1px 0 2px 1px"},
			},
			LinkStyle:  LinkStyleBorder,
			Icon:       "QuestionCircle",
			HideLabel:  true,
			Href:       sea.HelpURL,
			LinkTarget: "_blank",
		}
	}
	ccInp := func(value string) *Input {
		inp := &Input{
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
			Full:        true,
			AutoFocus:   sea.AutoFocus,
		}
		inp.SetProperty("value", value)
		return inp
	}
	ccMap := map[string]func() ClientComponent{
		"btn_search": func() ClientComponent {
			return ccBtn("Search")
		},
		"btn_help": func() ClientComponent {
			if sea.HelpURL != "" {
				return ccLnk()
			}
			return ccBtn("QuestionCircle")
		},
		"filter_value": func() ClientComponent {
			return ccInp(ut.ToString(sea.Data["filter_value"], ""))
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
				PageSize:          sea.PageSize,
				HidePaginatonSize: sea.HidePaginatonSize,
				TableFilter:       false,
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
	<div class="cell">{{ searchComponent "filter_value" }}</div>
	<div class="cell" style="width: 20px;" >{{ searchComponent "btn_search" }}</div>
	{{ if .ShowHelp }}<div class="cell" style="width: 20px;" >{{ searchComponent "btn_help" }}</div>{{ end }}
	</div>
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
		return toast(ut.ToString(evt.Value, ""))
	}
	if evt.Name == SearchEventHelp {
		return toast("Help!!!")
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
	},
	{
		"city":       "City4",
		"custname":   "Third Customer Foundation",
		"custnumber": "DMCUST/00003",
		"id":         "customer-4",
		"label":      "Third Customer Foundation",
		"street":     "street 4.",
	},
	{
		"city":       "City5",
		"custname":   "First Customer Co.",
		"custnumber": "DMCUST/00004",
		"id":         "customer-5",
		"label":      "First Customer Co.",
		"street":     "street 1.",
	},
	{
		"city":       "City3",
		"custname":   "Second Customer Name",
		"custnumber": "DMCUST/00005",
		"id":         "customer-6",
		"label":      "Second Customer Name",
		"street":     "street 3.",
	},
	{
		"city":       "City4",
		"custname":   "Third Customer Foundation",
		"custnumber": "DMCUST/00006",
		"id":         "customer-7",
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
				Fields:            testSearchFields,
				Rows:              testSearchRows,
				FilterPlaceholder: "Customer Name, Number, City, Steet...",
				PageSize:          5,
				AutoFocus:         true,
			},
		},
		{
			Label:         "Filter default value",
			ComponentType: ComponentTypeList,
			Component: &Search{
				BaseComponent: BaseComponent{
					Id:           id + "_search_filter_default",
					EventURL:     eventURL,
					OnResponse:   testSearchResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Data: ut.IM{
						"filter_value": "DMCUST",
					},
				},
				Fields:    testSearchFields,
				Rows:      testSearchRows,
				AutoFocus: true,
				ShowHelp:  true,
			},
		},
		{
			Label:         "Help URL",
			ComponentType: ComponentTypeList,
			Component: &Search{
				BaseComponent: BaseComponent{
					Id:           id + "_search_help",
					EventURL:     eventURL,
					OnResponse:   testSearchResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Fields:    testSearchFields,
				Rows:      testSearchRows,
				AutoFocus: true,
				ShowHelp:  true,
				HelpURL:   "https://www.google.com",
			},
		},
	}
}
