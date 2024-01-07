package atom

import (
	"math"
	"strings"

	bc "github.com/nervatura/component/component/base"
)

const (
	NumberEventChange = "change"
)

type NumberInput struct {
	bc.BaseComponent
	Value     float64 `json:"value"`
	Integer   bool    `json:"integer"`
	Label     string  `json:"label"`
	SetMax    bool    `json:"set_max"`
	MaxValue  float64 `json:"max_value"`
	SetMin    bool    `json:"set_min"`
	MinValue  float64 `json:"min_value"`
	Disabled  bool    `json:"disabled"`
	ReadOnly  bool    `json:"read_only"`
	AutoFocus bool    `json:"auto_focus"`
	Full      bool    `json:"full"`
}

func (inp *NumberInput) Properties() bc.IM {
	return bc.MergeIM(
		inp.BaseComponent.Properties(),
		bc.IM{
			"value":      inp.Value,
			"integer":    inp.Integer,
			"label":      inp.Label,
			"set_max":    inp.SetMax,
			"max_value":  inp.MaxValue,
			"set_min":    inp.SetMin,
			"min_value":  inp.MinValue,
			"disabled":   inp.Disabled,
			"readonly":   inp.ReadOnly,
			"auto_focus": inp.AutoFocus,
			"full":       inp.Full,
		})
}

func (inp *NumberInput) GetProperty(propName string) interface{} {
	return inp.Properties()[propName]
}

func (inp *NumberInput) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			value := bc.ToFloat(propValue, 0)
			if inp.SetMax && (value > inp.MaxValue) {
				value = inp.MaxValue
			}
			if inp.SetMin && (value < inp.MinValue) {
				value = inp.MinValue
			}
			if inp.Integer {
				value = math.Floor(value)
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if inp.BaseComponent.GetProperty(propName) != nil {
		return inp.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

func (inp *NumberInput) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			inp.Value = inp.Validation(propName, propValue).(float64)
			return inp.Value
		},
		"integer": func() interface{} {
			inp.Integer = bc.ToBoolean(propValue, false)
			return inp.Integer
		},
		"set_max": func() interface{} {
			inp.SetMax = bc.ToBoolean(propValue, false)
			return inp.SetMax
		},
		"disabled": func() interface{} {
			inp.Disabled = bc.ToBoolean(propValue, false)
			return inp.Disabled
		},
		"set_min": func() interface{} {
			inp.SetMin = bc.ToBoolean(propValue, false)
			return inp.SetMin
		},
		"readonly": func() interface{} {
			inp.ReadOnly = bc.ToBoolean(propValue, false)
			return inp.ReadOnly
		},
		"auto_focus": func() interface{} {
			inp.AutoFocus = bc.ToBoolean(propValue, false)
			return inp.AutoFocus
		},
		"full": func() interface{} {
			inp.Full = bc.ToBoolean(propValue, false)
			return inp.Full
		},
		"label": func() interface{} {
			inp.Label = bc.ToString(propValue, "")
			return inp.Label
		},
		"max_value": func() interface{} {
			inp.MaxValue = bc.ToFloat(propValue, 0)
			return inp.MaxValue
		},
		"min_value": func() interface{} {
			inp.MinValue = bc.ToFloat(propValue, 0)
			return inp.MinValue
		},
	}
	if _, found := pm[propName]; found {
		return inp.SetRequestValue(propName, pm[propName](), []string{})
	}
	if inp.BaseComponent.GetProperty(propName) != nil {
		return inp.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (inp *NumberInput) OnRequest(te bc.TriggerEvent) (re bc.ResponseEvent) {
	value := inp.SetProperty("value", te.Values.Get(te.Name))
	evt := bc.ResponseEvent{
		Trigger:     inp,
		TriggerName: inp.Name,
		Name:        NumberEventChange,
		Value:       value,
	}
	if inp.OnResponse != nil {
		return inp.OnResponse(evt)
	}
	return evt
}

func (inp *NumberInput) Render() (res string, err error) {
	inp.InitProps(inp)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(inp.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(inp.Class, " ")
		},
	}
	tpl := `<input id="{{ .Id }}" name="{{ .Name }}" type="number" value="{{ .Value }}"
	{{ if ne .EventURL "" }} hx-post="{{ .EventURL }}" hx-target="{{ .Target }}" hx-swap="{{ .Swap }}"{{ end }}
	{{ if ne .Indicator "" }} hx-indicator="#{{ .Indicator }}"{{ end }}
	{{ if .ReadOnly }} readonly{{ end }}
	{{ if .Disabled }} disabled{{ end }}
	{{ if .AutoFocus }} autofocus{{ end }}
	{{ if ne .Label "" }} aria-label="{{ .Label }}"{{ end }}
	 class="{{ customClass }}{{ if .Full }} full{{ end }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	></input>`

	if res, err = bc.TemplateBuilder("number", tpl, funcMap, inp); err == nil && inp.EventURL != "" {
		inp.SetProperty("request_map", inp)
	}
	return res, nil
}

func DemoNumberInput(demo bc.ClientComponent) []bc.DemoComponent {
	id := bc.ToString(demo.GetProperty("id"), "")
	eventURL := bc.ToString(demo.GetProperty("event_url"), "")
	requestValue := demo.GetProperty("request_value").(map[string]bc.IM)
	requestMap := demo.GetProperty("request_map").(map[string]bc.ClientComponent)
	return []bc.DemoComponent{
		{
			Label:         "Default",
			ComponentType: bc.ComponentTypeNumberInput,
			Component: &NumberInput{
				BaseComponent: bc.BaseComponent{
					Id:           id + "_number_default",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value: 1.5,
			}},
		{
			Label:         "Integer and AutoFocus",
			ComponentType: bc.ComponentTypeNumberInput,
			Component: &NumberInput{
				BaseComponent: bc.BaseComponent{
					Id:           id + "_number_integer",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value:     0,
				Integer:   true,
				AutoFocus: true,
			}},
		{
			Label:         "Max(100) and min(50) value, full",
			ComponentType: bc.ComponentTypeNumberInput,
			Component: &NumberInput{
				BaseComponent: bc.BaseComponent{
					Id:           id + "_number_min_max",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value:    150,
				SetMax:   true,
				MaxValue: 100,
				SetMin:   true,
				MinValue: 50,
				Full:     true,
			}},
		{
			Label:         "ReadOnly",
			ComponentType: bc.ComponentTypeNumberInput,
			Component: &NumberInput{
				BaseComponent: bc.BaseComponent{
					Id: id + "_number_readonly",
				},
				Value:    1234,
				ReadOnly: true,
			}},
		{
			Label:         "Disabled",
			ComponentType: bc.ComponentTypeNumberInput,
			Component: &NumberInput{
				BaseComponent: bc.BaseComponent{
					Id: id + "_number_disabled",
				},
				Value:    1234,
				Disabled: true,
			}},
	}
}
