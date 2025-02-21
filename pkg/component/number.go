package component

import (
	"html/template"
	"math"
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [NumberInput] constants
const (
	ComponentTypeNumberInput = "number"

	NumberEventChange = "number_change"
)

/*
Creates an HTML number input control

For example:

	&NumberInput{
	  BaseComponent: BaseComponent{
	    Id:           "id_number_min_max",
	    EventURL:     "/event",
	    RequestValue: parent_component.GetProperty("request_value").(map[string]ut.IM),
	    RequestMap:   parent_component.GetProperty("request_map").(map[string]ClientComponent),
	  },
	  Value:    150,
	  SetMax:   true,
	  MaxValue: 100,
	  SetMin:   true,
	  MinValue: 50,
	  Full:     true,
	}
*/
type NumberInput struct {
	BaseComponent
	// Any valid value based on control type (float64 or integer)
	Value float64 `json:"value"`
	// Integer type input
	Integer bool `json:"integer"`
	// The HTML aria-label attribute of the component
	Label string `json:"label"`
	// Enable maximum value monitoring
	SetMax bool `json:"set_max"`
	// Maximum value that can be entered
	MaxValue float64 `json:"max_value"`
	// Enable minimum value monitoring
	SetMin bool `json:"set_min"`
	// Minimum value that can be entered
	MinValue float64 `json:"min_value"`
	// Specifies that the input should be disabled
	Disabled bool `json:"disabled"`
	// Specifies that the input field is read-only
	ReadOnly bool `json:"readonly"`
	// Specifies that the input element should automatically get focus when the page loads
	AutoFocus bool `json:"auto_focus"`
	// Full width input (100%)
	Full bool `json:"full"`
}

/*
Returns all properties of the [NumberInput]
*/
func (inp *NumberInput) Properties() ut.IM {
	return ut.MergeIM(
		inp.BaseComponent.Properties(),
		ut.IM{
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

/*
Returns the value of the property of the [NumberInput] with the specified name.
*/
func (inp *NumberInput) GetProperty(propName string) interface{} {
	return inp.Properties()[propName]
}

/*
It checks the value given to the property of the [NumberInput] and always returns a valid value
*/
func (inp *NumberInput) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			value := ut.ToFloat(propValue, 0)
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

/*
Setting a property of the [NumberInput] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (inp *NumberInput) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			inp.Value = inp.Validation(propName, propValue).(float64)
			return inp.Value
		},
		"integer": func() interface{} {
			inp.Integer = ut.ToBoolean(propValue, false)
			return inp.Integer
		},
		"set_max": func() interface{} {
			inp.SetMax = ut.ToBoolean(propValue, false)
			return inp.SetMax
		},
		"disabled": func() interface{} {
			inp.Disabled = ut.ToBoolean(propValue, false)
			return inp.Disabled
		},
		"set_min": func() interface{} {
			inp.SetMin = ut.ToBoolean(propValue, false)
			return inp.SetMin
		},
		"readonly": func() interface{} {
			inp.ReadOnly = ut.ToBoolean(propValue, false)
			return inp.ReadOnly
		},
		"auto_focus": func() interface{} {
			inp.AutoFocus = ut.ToBoolean(propValue, false)
			return inp.AutoFocus
		},
		"full": func() interface{} {
			inp.Full = ut.ToBoolean(propValue, false)
			return inp.Full
		},
		"label": func() interface{} {
			inp.Label = ut.ToString(propValue, "")
			return inp.Label
		},
		"max_value": func() interface{} {
			inp.MaxValue = ut.ToFloat(propValue, 0)
			return inp.MaxValue
		},
		"min_value": func() interface{} {
			inp.MinValue = ut.ToFloat(propValue, 0)
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

/*
If the OnResponse function of the [NumberInput] is implemented, the function calls it after the [TriggerEvent]
is processed, otherwise the function's return [ResponseEvent] is the processed [TriggerEvent].
*/
func (inp *NumberInput) OnRequest(te TriggerEvent) (re ResponseEvent) {
	value := inp.SetProperty("value", te.Values.Get(te.Name))
	evt := ResponseEvent{
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

/*
Based on the values, it will generate the html code of the [NumberInput] or return with an error message.
*/
func (inp *NumberInput) Render() (html template.HTML, err error) {
	inp.InitProps(inp)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(inp.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(inp.Class, " ")
		},
		"value": func() string {
			return ut.ToString(inp.Value, "0")
		},
	}
	tpl := `<input id="{{ .Id }}" name="{{ .Name }}" type="number" value="{{ value }}"
	{{ if .Integer }} step="1"{{ else }} step="any"{{ end }}
	{{ if ne .EventURL "" }} hx-post="{{ .EventURL }}" hx-target="{{ .Target }}" {{ if ne .Sync "none" }} hx-sync="{{ .Sync }}"{{ end }} hx-swap="{{ .Swap }}"{{ end }}
	{{ if ne .Indicator "none" }} hx-indicator="#{{ .Indicator }}"{{ end }}
	{{ if .ReadOnly }} readonly{{ end }}
	{{ if .Disabled }} disabled{{ end }}
	{{ if .AutoFocus }} autofocus{{ end }}
	{{ if ne .Label "" }} aria-label="{{ .Label }}"{{ end }}
	 class="{{ customClass }}{{ if .Full }} full{{ end }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	></input>`

	if html, err = ut.TemplateBuilder("number", tpl, funcMap, inp); err == nil && inp.EventURL != "" {
		inp.SetProperty("request_map", inp)
	}
	return html, nil
}

// [NumberInput] test and demo data
func TestNumberInput(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeNumberInput,
			Component: &NumberInput{
				BaseComponent: BaseComponent{
					Id:           id + "_number_default",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value: 1.5,
			}},
		{
			Label:         "Integer and AutoFocus",
			ComponentType: ComponentTypeNumberInput,
			Component: &NumberInput{
				BaseComponent: BaseComponent{
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
			ComponentType: ComponentTypeNumberInput,
			Component: &NumberInput{
				BaseComponent: BaseComponent{
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
			ComponentType: ComponentTypeNumberInput,
			Component: &NumberInput{
				BaseComponent: BaseComponent{
					Id: id + "_number_readonly",
				},
				Value:    1234,
				ReadOnly: true,
			}},
		{
			Label:         "Disabled",
			ComponentType: ComponentTypeNumberInput,
			Component: &NumberInput{
				BaseComponent: BaseComponent{
					Id: id + "_number_disabled",
				},
				Value:    1234,
				Disabled: true,
			}},
	}
}
