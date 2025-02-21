package component

import (
	"html/template"
	"slices"
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [Pagination] constants
const (
	ComponentTypePagination = "pagination"

	PaginationEventValue    = "pagination_value"
	PaginationEventPageSize = "pagination_page_size"
)

// [Pagination] PageSize values
var ValidPageSize []int64 = []int64{5, 10, 20, 50, 100}

/*
Creates a pagination control

For example:

	&Pagination{
	  BaseComponent: BaseComponent{
	    Id:           "id_table_page_size",
	    EventURL:     "/event",
	    RequestValue: parent_component.GetProperty("request_value").(map[string]ut.IM),
	    RequestMap:   parent_component.GetProperty("request_map").(map[string]ClientComponent),
	  },
	  Value:        2,
	  PageSize:     10,
	  PageCount:    3,
	  HidePageSize: false,
	}
*/
type Pagination struct {
	BaseComponent
	// Current page value
	Value int64 `json:"value"`
	// [PageSize] variable constants: 5, 10, 20, 50, 100. Default value: 10
	PageSize int64 `json:"page_size"`
	// The maximum value of the pagination
	PageCount int64 `json:"page_count"`
	// Show/hide page size selector
	HidePageSize bool `json:"hide_pageSize"`
}

/*
Returns all properties of the [Pagination]
*/
func (pgn *Pagination) Properties() ut.IM {
	return ut.MergeIM(
		pgn.BaseComponent.Properties(),
		ut.IM{
			"value":          pgn.Value,
			"page_size":      pgn.PageSize,
			"page_count":     pgn.PageCount,
			"hide_page_size": pgn.HidePageSize,
		})
}

/*
Returns the value of the property of the [Pagination] with the specified name.
*/
func (pgn *Pagination) GetProperty(propName string) interface{} {
	return pgn.Properties()[propName]
}

/*
It checks the value given to the property of the [Pagination] and always returns a valid value
*/
func (pgn *Pagination) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			value := ut.ToInteger(propValue, 1)
			if value > pgn.PageCount {
				value = pgn.PageCount
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
			pgn.SetProperty("id", pgn.Id)
			value := ut.ToString(propValue, pgn.Id)
			if value != "this" && !strings.HasPrefix(value, "#") {
				value = "#" + value
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if pgn.BaseComponent.GetProperty(propName) != nil {
		return pgn.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [Pagination] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (pgn *Pagination) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			pgn.Value = pgn.Validation(propName, propValue).(int64)
			return pgn.Value
		},
		"page_count": func() interface{} {
			pgn.PageCount = ut.ToInteger(propValue, 1)
			return pgn.PageCount
		},
		"page_size": func() interface{} {
			pgn.PageSize = pgn.Validation(propName, propValue).(int64)
			return pgn.PageSize
		},
		"hide_page_size": func() interface{} {
			pgn.HidePageSize = ut.ToBoolean(propValue, false)
			return pgn.HidePageSize
		},
		"target": func() interface{} {
			pgn.Target = pgn.Validation(propName, propValue).(string)
			return pgn.Target
		},
	}
	if _, found := pm[propName]; found {
		return pgn.SetRequestValue(propName, pm[propName](), []string{})
	}
	if pgn.BaseComponent.GetProperty(propName) != nil {
		return pgn.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (pgn *Pagination) response(evt ResponseEvent) (re ResponseEvent) {
	pgnEvt := ResponseEvent{Trigger: pgn, TriggerName: pgn.Name}
	switch evt.TriggerName {
	case "pagination_page_size":
		value := pgn.SetProperty("page_size", evt.Value).(int64)
		pgnEvt.Name = PaginationEventPageSize
		pgnEvt.Value = value
	case "pagination_input_value":
		value := pgn.SetProperty("value", evt.Value).(int64)
		pgnEvt.Name = PaginationEventValue
		pgnEvt.Value = value
	default:
		data := evt.Trigger.GetProperty("data").(ut.IM)
		value := pgn.SetProperty("value", data["value"]).(int64)
		pgnEvt.Name = PaginationEventValue
		pgnEvt.Value = value
		pgn.Value = value
	}
	if pgn.OnResponse != nil {
		return pgn.OnResponse(pgnEvt)
	}
	return pgnEvt
}

func (pgn *Pagination) getComponent(name string) (html template.HTML, err error) {
	ccBtn := func(label, value string, disabled bool, style ut.SM) *Button {
		return &Button{
			BaseComponent: BaseComponent{
				Id: pgn.Id + "_" + name, Name: name,
				Data: ut.IM{
					"value": value,
				},
				Style:        style,
				EventURL:     pgn.EventURL,
				Target:       pgn.Target,
				OnResponse:   pgn.response,
				RequestValue: pgn.RequestValue,
				RequestMap:   pgn.RequestMap,
			},
			ButtonStyle: ButtonStyleBorder,
			Label:       label, Disabled: disabled,
		}
	}
	ccMap := map[string]func() ClientComponent{
		"pagination_btn_first": func() ClientComponent {
			return ccBtn(
				"1", "1", !(pgn.Value > 1),
				ut.SM{"padding": "6px 6px 7px", "font-size": "15px", "margin": "1px 1px 2px 0px"},
			)
		},
		"pagination_btn_previous": func() ClientComponent {
			return ccBtn(
				"❮", ut.ToString(pgn.Value-1, "1"), !(pgn.Value > 1),
				ut.SM{"padding": "5px 6px 8px", "font-size": "15px", "margin": "1px 0 2px"},
			)
		},
		"pagination_btn_next": func() ClientComponent {
			return ccBtn(
				"❯", ut.ToString(pgn.Value+1, "1"), !(pgn.Value < pgn.PageCount),
				ut.SM{"padding": "5px 6px 8px", "font-size": "15px", "margin": "1px 1px 2px 0px"},
			)
		},
		"pagination_btn_last": func() ClientComponent {
			return ccBtn(
				ut.ToString(pgn.PageCount, ""), ut.ToString(pgn.PageCount, "1"),
				!(pgn.Value < pgn.PageCount),
				ut.SM{"padding": "6px 6px 7px", "font-size": "15px", "margin": "1px 0 2px"},
			)
		},
		"pagination_input_value": func() ClientComponent {
			inp := &NumberInput{
				BaseComponent: BaseComponent{
					Id:           pgn.Id + "_" + name,
					Name:         name,
					Style:        ut.SM{"padding": "7px", "width": "60px", "font-weight": "bold"},
					EventURL:     pgn.EventURL,
					Target:       pgn.Target,
					OnResponse:   pgn.response,
					RequestValue: pgn.RequestValue,
					RequestMap:   pgn.RequestMap,
				},
				Label:    "Page",
				Integer:  true,
				Disabled: (pgn.PageCount == 0),
				SetMin:   true, MinValue: 1,
				SetMax: true, MaxValue: float64(pgn.PageCount),
			}
			inp.SetProperty("value", pgn.Value)
			return inp
		},
		"pagination_page_size": func() ClientComponent {
			sel := &Select{
				BaseComponent: BaseComponent{
					Id: pgn.Id + "_" + name, Name: name,
					Style:        ut.SM{"padding": "7px"},
					EventURL:     pgn.EventURL,
					Target:       pgn.Target,
					OnResponse:   pgn.response,
					RequestValue: pgn.RequestValue,
					RequestMap:   pgn.RequestMap,
				},
				Label:    "Size",
				Disabled: (pgn.PageCount == 0),
				IsNull:   false, Value: ut.ToString(pgn.PageSize, ""),
				Options: []SelectOption{},
			}
			for _, ps := range ValidPageSize {
				sel.Options = append(sel.Options, SelectOption{
					Value: ut.ToString(ps, ""), Text: ut.ToString(ps, ""),
				})
			}
			return sel
		},
	}
	cc := ccMap[name]()
	html, err = cc.Render()
	return html, err
}

/*
Based on the values, it will generate the html code of the [Pagination] or return with an error message.
*/
func (pgn *Pagination) Render() (html template.HTML, err error) {
	pgn.InitProps(pgn)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(pgn.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(pgn.Class, " ")
		},
		"paginationComponent": func(name string) (template.HTML, error) {
			return pgn.getComponent(name)
		},
	}
	tpl := `<div id="{{ .Id }}" name="{{ .Name }}" class="row {{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	><div class="cell padding-small" >{{ paginationComponent "pagination_btn_first" }}{{ paginationComponent "pagination_btn_previous" }}</div>
	<div class="cell" >{{ paginationComponent "pagination_input_value" }}</div>
	<div class="cell padding-small" >{{ paginationComponent "pagination_btn_next" }}{{ paginationComponent "pagination_btn_last" }}</div>
	{{ if ne .HidePageSize true }}<div class="cell padding-small" >{{ paginationComponent "pagination_page_size" }}</div>{{ end }}
	</div>`

	return ut.TemplateBuilder("pagination", tpl, funcMap, pgn)
}

// [Pagination] test and demo data
func TestPagination(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypePagination,
			Component: &Pagination{
				BaseComponent: BaseComponent{
					Id:           id + "_table_default",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				PageSize:     5,
				HidePageSize: true,
			}},
		{
			Label:         "PageSize",
			ComponentType: ComponentTypePagination,
			Component: &Pagination{
				BaseComponent: BaseComponent{
					Id:           id + "_table_page_size",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value:        2,
				PageSize:     10,
				PageCount:    3,
				HidePageSize: false,
			}},
	}
}
