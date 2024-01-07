package atom

import (
	"strings"

	bc "github.com/nervatura/component/component/base"
)

const (
	SelectEventChange = "change"
)

type SelectOption struct {
	Value string `json:"value"`
	Text  string `json:"text"`
}

type Select struct {
	bc.BaseComponent
	Value     string         `json:"value"`
	Options   []SelectOption `json:"options"`
	IsNull    bool           `json:"is_null"`
	Label     string         `json:"label"`
	Disabled  bool           `json:"disabled"`
	AutoFocus bool           `json:"auto_focus"`
	Full      bool           `json:"full"`
}

func (sel *Select) Properties() bc.IM {
	return bc.MergeIM(
		sel.BaseComponent.Properties(),
		bc.IM{
			"value":      sel.Value,
			"options":    sel.Options,
			"is_null":    sel.IsNull,
			"label":      sel.Label,
			"disabled":   sel.Disabled,
			"auto_focus": sel.AutoFocus,
			"full":       sel.Full,
		})
}

func (sel *Select) GetProperty(propName string) interface{} {
	return sel.Properties()[propName]
}

func (sel *Select) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"options": func() interface{} {
			if options, valid := propValue.([]SelectOption); valid && (options != nil) {
				return options
			}
			return []SelectOption{}
		},
		"value": func() interface{} {
			value := bc.ToString(propValue, "")
			if (value == "") && sel.IsNull {
				return value
			}
			valid := false
			for _, opt := range sel.Options {
				if opt.Value == value {
					valid = true
				}
			}
			if !valid {
				if len(sel.Options) > 0 {
					value = sel.Options[0].Value
				} else {
					value = ""
				}
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if sel.BaseComponent.GetProperty(propName) != nil {
		return sel.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

func (sel *Select) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			sel.Value = sel.Validation(propName, propValue).(string)
			return sel.Value
		},
		"options": func() interface{} {
			sel.Options = sel.Validation(propName, propValue).([]SelectOption)
			return sel.Options
		},
		"is_null": func() interface{} {
			sel.IsNull = bc.ToBoolean(propValue, false)
			return sel.IsNull
		},
		"label": func() interface{} {
			sel.Label = bc.ToString(propValue, "")
			return sel.Label
		},
		"disabled": func() interface{} {
			sel.Disabled = bc.ToBoolean(propValue, false)
			return sel.Disabled
		},
		"auto_focus": func() interface{} {
			sel.AutoFocus = bc.ToBoolean(propValue, false)
			return sel.AutoFocus
		},
		"full": func() interface{} {
			sel.Full = bc.ToBoolean(propValue, false)
			return sel.Full
		},
	}
	if _, found := pm[propName]; found {
		return sel.SetRequestValue(propName, pm[propName](), []string{})
	}
	if sel.BaseComponent.GetProperty(propName) != nil {
		return sel.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (sel *Select) OnRequest(te bc.TriggerEvent) (re bc.ResponseEvent) {
	value := sel.SetProperty("value", te.Values.Get(te.Name))
	evt := bc.ResponseEvent{
		Trigger: sel, TriggerName: sel.Name,
		Name:  SelectEventChange,
		Value: value,
	}
	if sel.OnResponse != nil {
		return sel.OnResponse(evt)
	}
	return evt
}

func (sel *Select) Render() (res string, err error) {
	sel.InitProps(sel)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(sel.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(sel.Class, " ")
		},
	}
	tpl := `<select id="{{ .Id }}" name="{{ .Name }}" value="{{ .Value }}"
	{{ if ne .EventURL "" }} hx-post="{{ .EventURL }}" hx-target="{{ .Target }}" hx-swap="{{ .Swap }}"{{ end }}
	{{ if ne .Indicator "" }} hx-indicator="#{{ .Indicator }}"{{ end }}
	{{ if .Disabled }} disabled{{ end }}
	{{ if .AutoFocus }} autofocus{{ end }}
	{{ if ne .Label "" }} aria-label="{{ .Label }}"{{ end }}
	 class="{{ customClass }}{{ if .Full }} full{{ end }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	>{{ if .IsNull }}<option {{ if eq .Value "" }}selected{{ end }} key="-1" value="" ></option>{{ end }}
	{{ range $index, $option := .Options }}<option {{ if eq .Value $.Value }}selected{{ end }} key="{{ $index }}" value="{{ $option.Value }}" >{{ $option.Text }}</option>{{ end }}
	</select>`

	if res, err = bc.TemplateBuilder("button", tpl, funcMap, sel); err == nil && sel.EventURL != "" {
		sel.SetProperty("request_map", sel)
	}
	return res, nil
}

func DemoSelect(demo bc.ClientComponent) []bc.DemoComponent {
	id := bc.ToString(demo.GetProperty("id"), "")
	eventURL := bc.ToString(demo.GetProperty("event_url"), "")
	requestValue := demo.GetProperty("request_value").(map[string]bc.IM)
	requestMap := demo.GetProperty("request_map").(map[string]bc.ClientComponent)
	return []bc.DemoComponent{
		{
			Label:         "Default",
			ComponentType: bc.ComponentTypeSelect,
			Component: &Select{
				BaseComponent: bc.BaseComponent{
					Id:           id + "_select_default",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value: "value1",
				Options: []SelectOption{
					{Value: "value1", Text: "Text 1"},
					{Value: "value2", Text: "Text 2"},
					{Value: "value3", Text: "Text 3"},
				},
				IsNull: true,
			}},
		{
			Label:         "Not null",
			ComponentType: bc.ComponentTypeSelect,
			Component: &Select{
				BaseComponent: bc.BaseComponent{
					Id:           id + "_select_not_null",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value: "value2",
				Options: []SelectOption{
					{Value: "value1", Text: "Text 1"},
					{Value: "value2", Text: "Text 2"},
					{Value: "value3", Text: "Text 3"},
				},
				IsNull: false,
				Full:   true,
			}},
	}
}
