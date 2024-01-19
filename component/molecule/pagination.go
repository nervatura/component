package molecule

import (
	"strings"

	fm "github.com/nervatura/component/component/atom"
	bc "github.com/nervatura/component/component/base"
)

const (
	PaginationEventValue    = "value"
	PaginationEventPageSize = "page_size"

	PaginationTypeTop    = "top"
	PaginationTypeBottom = "bottom"
	PaginationTypeAll    = "all"
	PaginationTypeNone   = "none"
)

var PaginationType []string = []string{PaginationTypeTop, PaginationTypeBottom, PaginationTypeAll, PaginationTypeNone}
var ValidPageSize []int64 = []int64{5, 10, 20, 50, 100}

type Pagination struct {
	bc.BaseComponent
	Value        int64 `json:"value"`
	PageSize     int64 `json:"page_size"`
	PageCount    int64 `json:"page_count"`
	HidePageSize bool  `json:"hide_pageSize"`
}

func (pgn *Pagination) Properties() bc.IM {
	return bc.MergeIM(
		pgn.BaseComponent.Properties(),
		bc.IM{
			"value":          pgn.Value,
			"page_size":      pgn.PageSize,
			"page_count":     pgn.PageCount,
			"hide_page_size": pgn.HidePageSize,
		})
}

func (pgn *Pagination) GetProperty(propName string) interface{} {
	return pgn.Properties()[propName]
}

func (pgn *Pagination) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			value := bc.ToInteger(propValue, 1)
			if value > pgn.PageCount {
				value = pgn.PageCount
			}
			if value < 1 {
				value = 1
			}
			return value
		},
		"page_size": func() interface{} {
			value := bc.ToInteger(propValue, 10)
			pageSize := []string{}
			for _, ps := range ValidPageSize {
				pageSize = append(pageSize, bc.ToString(ps, ""))
			}
			if !bc.Contains(pageSize, bc.ToString(value, "")) {
				value = ValidPageSize[0]
			}
			return value
		},
		"target": func() interface{} {
			pgn.SetProperty("id", pgn.Id)
			value := bc.ToString(propValue, pgn.Id)
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

func (pgn *Pagination) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			pgn.Value = pgn.Validation(propName, propValue).(int64)
			return pgn.Value
		},
		"page_count": func() interface{} {
			pgn.PageCount = bc.ToInteger(propValue, 1)
			return pgn.PageCount
		},
		"page_size": func() interface{} {
			pgn.PageSize = pgn.Validation(propName, propValue).(int64)
			return pgn.PageSize
		},
		"hide_page_size": func() interface{} {
			pgn.HidePageSize = bc.ToBoolean(propValue, false)
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

func (pgn *Pagination) response(evt bc.ResponseEvent) (re bc.ResponseEvent) {
	pgnEvt := bc.ResponseEvent{Trigger: pgn, TriggerName: pgn.Name}
	if evt.TriggerName == "pagination_page_size" {
		value := pgn.SetProperty("page_size", evt.Value).(int64)
		pgnEvt.Name = PaginationEventPageSize
		pgnEvt.Value = value
	} else {
		data := evt.Trigger.GetProperty("data").(bc.IM)
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

func (pgn *Pagination) getComponent(name string) (res string, err error) {
	ccBtn := func(label, value string, disabled bool, style bc.SM) *fm.Button {
		return &fm.Button{
			BaseComponent: bc.BaseComponent{
				Id: pgn.Id + "_" + name, Name: name,
				Data: bc.IM{
					"value": value,
				},
				Style:        style,
				EventURL:     pgn.EventURL,
				Target:       pgn.Target,
				OnResponse:   pgn.response,
				RequestValue: pgn.RequestValue,
				RequestMap:   pgn.RequestMap,
			},
			Type:  fm.ButtonTypeBorder,
			Label: label, Disabled: disabled,
		}
	}
	ccMap := map[string]func() bc.ClientComponent{
		"pagination_btn_first": func() bc.ClientComponent {
			return ccBtn(
				"1", "1", !(pgn.Value > 1),
				bc.SM{"padding": "6px 6px 7px", "font-size": "15px", "margin": "1px 1px 2px 0px"},
			)
		},
		"pagination_btn_previous": func() bc.ClientComponent {
			return ccBtn(
				"&#10094;", bc.ToString(pgn.Value-1, "1"), !(pgn.Value > 1),
				bc.SM{"padding": "5px 6px 8px", "font-size": "15px", "margin": "1px 0 2px"},
			)
		},
		"pagination_btn_next": func() bc.ClientComponent {
			return ccBtn(
				"&#10095;", bc.ToString(pgn.Value+1, "1"), !(pgn.Value < pgn.PageCount),
				bc.SM{"padding": "5px 6px 8px", "font-size": "15px", "margin": "1px 1px 2px 0px"},
			)
		},
		"pagination_btn_last": func() bc.ClientComponent {
			return ccBtn(
				bc.ToString(pgn.PageCount, ""), bc.ToString(pgn.PageCount, "1"),
				!(pgn.Value < pgn.PageCount),
				bc.SM{"padding": "6px 6px 7px", "font-size": "15px", "margin": "1px 0 2px"},
			)
		},
		"pagination_input_value": func() bc.ClientComponent {
			return &fm.NumberInput{
				BaseComponent: bc.BaseComponent{
					Id:           pgn.Id + "_" + name,
					Name:         name,
					Style:        bc.SM{"padding": "7px", "width": "60px", "font-weight": "bold"},
					EventURL:     pgn.EventURL,
					Target:       pgn.Target,
					OnResponse:   pgn.response,
					RequestValue: pgn.RequestValue,
					RequestMap:   pgn.RequestMap,
				},
				Value:    float64(pgn.Value),
				Label:    "Page",
				Integer:  true,
				Disabled: (pgn.PageCount == 0),
				SetMin:   true, MinValue: 1,
				SetMax: true, MaxValue: float64(pgn.PageCount),
			}
		},
		"pagination_page_size": func() bc.ClientComponent {
			sel := &fm.Select{
				BaseComponent: bc.BaseComponent{
					Id: pgn.Id + "_" + name, Name: name,
					Style:        bc.SM{"padding": "7px"},
					EventURL:     pgn.EventURL,
					Target:       pgn.Target,
					OnResponse:   pgn.response,
					RequestValue: pgn.RequestValue,
					RequestMap:   pgn.RequestMap,
				},
				Label:    "Size",
				Disabled: (pgn.PageCount == 0),
				IsNull:   false, Value: bc.ToString(pgn.PageSize, ""),
				Options: []fm.SelectOption{},
			}
			for _, ps := range ValidPageSize {
				sel.Options = append(sel.Options, fm.SelectOption{
					Value: bc.ToString(ps, ""), Text: bc.ToString(ps, ""),
				})
			}
			return sel
		},
	}
	cc := ccMap[name]()
	res, err = cc.Render()
	return res, err
}

func (pgn *Pagination) Render() (res string, err error) {
	pgn.InitProps(pgn)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(pgn.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(pgn.Class, " ")
		},
		"paginationComponent": func(name string) (string, error) {
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

	return bc.TemplateBuilder("pagination", tpl, funcMap, pgn)
}

func DemoPagination(demo bc.ClientComponent) []bc.DemoComponent {
	id := bc.ToString(demo.GetProperty("id"), "")
	eventURL := bc.ToString(demo.GetProperty("event_url"), "")
	requestValue := demo.GetProperty("request_value").(map[string]bc.IM)
	requestMap := demo.GetProperty("request_map").(map[string]bc.ClientComponent)
	return []bc.DemoComponent{
		{
			Label:         "Default",
			ComponentType: bc.ComponentTypePagination,
			Component: &Pagination{
				BaseComponent: bc.BaseComponent{
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
			ComponentType: bc.ComponentTypePagination,
			Component: &Pagination{
				BaseComponent: bc.BaseComponent{
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
