package component

import (
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [Input] constants
const (
	ComponentTypeInput = "input"

	InputEventChange = "change"

	InputTypeText     = "text"
	InputTypeColor    = "color"
	InputTypeFile     = "file"
	InputTypePassword = "password"
)

// [Input] Type values
var InputType []string = []string{InputTypeText, InputTypeColor, InputTypeFile, InputTypePassword}

/*
Creates an HTML text, color, file or password type input control

For example:

	&Input{
	  BaseComponent: BaseComponent{
	    Id: "id_input_default",
	  },
	  Type:        InputTypeText,
	  Placeholder: "placeholder text",
	  AutoFocus:   true,
	}
*/
type Input struct {
	BaseComponent
	/* [InputType] variable constants:
	[InputTypeText], [InputTypeColor], [InputTypeFile], [InputTypePassword].
	Default value: [InputTypeText] */
	Type string `json:"type"`
	// Any valid value based on control type
	Value string `json:"value"`
	// Specifies a short hint that describes the expected value of the input element
	Placeholder string `json:"placeholder"`
	// The HTML aria-label attribute of the component
	Label string `json:"label"`
	// Specifies that the input should be disabled
	Disabled bool `json:"disabled"`
	// Specifies that the input field is read-only
	ReadOnly bool `json:"read_only"`
	// Specifies that the input element should automatically get focus when the page loads
	AutoFocus bool `json:"auto_focus"`
	// Sets the values of the invalid class style
	Invalid bool `json:"invalid"`
	/* Specifies a filter for what file types the user can pick from the file input dialog box
	(only for type=[InputTypeFile]) */
	Accept string `json:"accept"`
	// Specifies the maximum number of characters allowed in the input element
	MaxLength int64 `json:"max_length"`
	// Specifies the width, in characters, of the input element
	Size int64 `json:"size"`
	// Full width input (100%)
	Full bool `json:"full"`
}

/*
Returns all properties of the [Input]
*/
func (inp *Input) Properties() ut.IM {
	return ut.MergeIM(
		inp.BaseComponent.Properties(),
		ut.IM{
			"type":        inp.Type,
			"value":       inp.Value,
			"placeholder": inp.Placeholder,
			"label":       inp.Label,
			"disabled":    inp.Disabled,
			"readonly":    inp.ReadOnly,
			"auto_focus":  inp.AutoFocus,
			"invalid":     inp.Invalid,
			"accept":      inp.Accept,
			"max_length":  inp.MaxLength,
			"size":        inp.Size,
			"full":        inp.Full,
		})
}

/*
Returns the value of the property of the [Input] with the specified name.
*/
func (inp *Input) GetProperty(propName string) interface{} {
	return inp.Properties()[propName]
}

/*
It checks the value given to the property of the [Input] and always returns a valid value
*/
func (inp *Input) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"type": func() interface{} {
			return inp.CheckEnumValue(ut.ToString(propValue, ""), InputTypeText, InputType)
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
Setting a property of the [Input] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (inp *Input) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"type": func() interface{} {
			inp.Type = inp.Validation(propName, propValue).(string)
			return inp.Type
		},
		"value": func() interface{} {
			inp.Value = ut.ToString(propValue, "")
			return inp.Value
		},
		"placeholder": func() interface{} {
			inp.Placeholder = ut.ToString(propValue, "")
			return inp.Placeholder
		},
		"disabled": func() interface{} {
			inp.Disabled = ut.ToBoolean(propValue, false)
			return inp.Disabled
		},
		"readonly": func() interface{} {
			inp.ReadOnly = ut.ToBoolean(propValue, false)
			return inp.ReadOnly
		},
		"auto_focus": func() interface{} {
			inp.AutoFocus = ut.ToBoolean(propValue, false)
			return inp.AutoFocus
		},
		"invalid": func() interface{} {
			inp.Invalid = ut.ToBoolean(propValue, false)
			return inp.Invalid
		},
		"accept": func() interface{} {
			inp.Accept = ut.ToString(propValue, "")
			return inp.Accept
		},
		"max_length": func() interface{} {
			inp.MaxLength = ut.ToInteger(propValue, 0)
			return inp.MaxLength
		},
		"size": func() interface{} {
			inp.Size = ut.ToInteger(propValue, 0)
			return inp.Size
		},
		"full": func() interface{} {
			inp.Full = ut.ToBoolean(propValue, false)
			return inp.Full
		},
		"label": func() interface{} {
			inp.Label = ut.ToString(propValue, "")
			return inp.Label
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
If the OnResponse function of the [Input] is implemented, the function calls it after the [TriggerEvent]
is processed, otherwise the function's return [ResponseEvent] is the processed [TriggerEvent].
*/
func (inp *Input) OnRequest(te TriggerEvent) (re ResponseEvent) {
	value := inp.SetProperty("value", te.Values.Get(te.Name))
	evt := ResponseEvent{
		Trigger: inp, TriggerName: inp.Name,
		Name:  InputEventChange,
		Value: value,
	}
	if inp.OnResponse != nil {
		return inp.OnResponse(evt)
	}
	return evt
}

/*
Based on the values, it will generate the html code of the [Input] or return with an error message.
*/
func (inp *Input) Render() (res string, err error) {
	inp.InitProps(inp)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(inp.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(inp.Class, " ")
		},
	}
	tpl := `<input id="{{ .Id }}" name="{{ .Name }}" type="{{ .Type }}" value="{{ .Value }}"
	{{ if ne .EventURL "" }} hx-post="{{ .EventURL }}" hx-target="{{ .Target }}" hx-swap="{{ .Swap }}"{{ end }}
	{{ if ne .Indicator "" }} hx-indicator="#{{ .Indicator }}"{{ end }}
	{{ if ne .Placeholder "" }} placeholder="{{ .Placeholder }}"{{ end }}
	{{ if .ReadOnly }} readonly{{ end }}
	{{ if .Disabled }} disabled{{ end }}
	{{ if .AutoFocus }} autofocus{{ end }}
	{{ if ne .Label "" }} aria-label="{{ .Label }}"{{ end }}
	{{ if ne .Accept "" }} accept="{{ .Accept }}"{{ end }}
	{{ if gt .MaxLength 0 }} maxlength="{{ .MaxLength }}"{{ end }}
	{{ if gt .Size 0 }} size="{{ .Size }}"{{ end }}
	 class="{{ if .Full }} full{{ end }}{{ if .Invalid }} invalid{{ end }} {{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	></input>`

	if res, err = ut.TemplateBuilder("input", tpl, funcMap, inp); err == nil && inp.EventURL != "" {
		inp.SetProperty("request_map", inp)
	}
	return res, nil
}

var demoInputResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
	evt.Trigger.SetProperty("invalid", false)
	if evt.Value != "valid" {
		evt.Trigger.SetProperty("invalid", true)
	}
	return evt
}

// [Input] test and demo data
func TestInput(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default and AutoFocus",
			ComponentType: ComponentTypeInput,
			Component: &Input{
				BaseComponent: BaseComponent{
					Id: id + "_input_default",
				},
				Type:        InputTypeText,
				Placeholder: "placeholder text",
				AutoFocus:   true,
			}},
		{
			Label:         "Valid value: valid",
			ComponentType: ComponentTypeInput,
			Component: &Input{
				BaseComponent: BaseComponent{
					Id:           id + "_input_valid",
					EventURL:     eventURL,
					Swap:         SwapOuterHTML,
					OnResponse:   demoInputResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type:  InputTypeText,
				Value: "valid",
			}},
		{
			Label:         "ReadOnly",
			ComponentType: ComponentTypeInput,
			Component: &Input{
				BaseComponent: BaseComponent{
					Id: id + "_input_readonly",
				},
				Type:     InputTypeText,
				Value:    "hello",
				ReadOnly: true,
			}},
		{
			Label:         "Disabled",
			ComponentType: ComponentTypeInput,
			Component: &Input{
				BaseComponent: BaseComponent{
					Id: id + "_input_disabled",
				},
				Type:     InputTypeText,
				Value:    "hello",
				Disabled: true,
			}},
		{
			Label:         "Password full",
			ComponentType: ComponentTypeInput,
			Component: &Input{
				BaseComponent: BaseComponent{
					Id: id + "_input_password",
				},
				Type:  InputTypePassword,
				Value: "secret",
				Full:  true,
			}},
		{
			Label:         "File input",
			ComponentType: ComponentTypeInput,
			Component: &Input{
				BaseComponent: BaseComponent{
					Id: id + "_input_file",
				},
				Type:   InputTypeFile,
				Accept: ".jpg,.png",
			}},
		{
			Label:         "Color input",
			ComponentType: ComponentTypeInput,
			Component: &Input{
				BaseComponent: BaseComponent{
					Id: id + "_input_color",
				},
				Type:  InputTypeColor,
				Value: "#845185",
			}},
	}
}
