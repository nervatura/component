package atom

import (
	"strings"

	bc "github.com/nervatura/component/component/base"
)

const (
	InputEventChange = "change"

	InputTypeText     = "text"
	InputTypeColor    = "color"
	InputTypeFile     = "file"
	InputTypePassword = "password"
)

var InputType []string = []string{InputTypeText, InputTypeColor, InputTypeFile, InputTypePassword}

type Input struct {
	bc.BaseComponent
	Type        string
	Value       string
	Placeholder string
	Label       string
	Disabled    bool
	ReadOnly    bool
	AutoFocus   bool
	Invalid     bool
	Accept      string
	MaxLength   int64
	Size        int64
	Full        bool
}

func (inp *Input) Properties() bc.IM {
	return bc.MergeIM(
		inp.BaseComponent.Properties(),
		bc.IM{
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

func (inp *Input) GetProperty(propName string) interface{} {
	return inp.Properties()[propName]
}

func (inp *Input) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"type": func() interface{} {
			return inp.CheckEnumValue(bc.ToString(propValue, ""), InputTypeText, InputType)
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

func (inp *Input) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"type": func() interface{} {
			inp.Type = inp.Validation(propName, propValue).(string)
			return inp.Type
		},
		"value": func() interface{} {
			inp.Value = bc.ToString(propValue, "")
			return inp.Value
		},
		"placeholder": func() interface{} {
			inp.Placeholder = bc.ToString(propValue, "")
			return inp.Placeholder
		},
		"disabled": func() interface{} {
			inp.Disabled = bc.ToBoolean(propValue, false)
			return inp.Disabled
		},
		"readonly": func() interface{} {
			inp.ReadOnly = bc.ToBoolean(propValue, false)
			return inp.ReadOnly
		},
		"auto_focus": func() interface{} {
			inp.AutoFocus = bc.ToBoolean(propValue, false)
			return inp.AutoFocus
		},
		"invalid": func() interface{} {
			inp.Invalid = bc.ToBoolean(propValue, false)
			return inp.Invalid
		},
		"accept": func() interface{} {
			inp.Accept = bc.ToString(propValue, "")
			return inp.Accept
		},
		"max_length": func() interface{} {
			inp.MaxLength = bc.ToInteger(propValue, 0)
			return inp.MaxLength
		},
		"size": func() interface{} {
			inp.Size = bc.ToInteger(propValue, 0)
			return inp.Size
		},
		"full": func() interface{} {
			inp.Full = bc.ToBoolean(propValue, false)
			return inp.Full
		},
		"label": func() interface{} {
			inp.Label = bc.ToString(propValue, "")
			return inp.Label
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if inp.BaseComponent.GetProperty(propName) != nil {
		return inp.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (inp *Input) OnRequest(te bc.TriggerEvent) (re bc.ResponseEvent) {
	value := inp.SetProperty("value", te.Values.Get(te.Name))
	evt := bc.ResponseEvent{
		Trigger: inp, TriggerName: inp.Name,
		Name:  InputEventChange,
		Value: value,
	}
	if inp.OnResponse != nil {
		return inp.OnResponse(evt)
	}
	return evt
}

func (inp *Input) InitProps() {
	for key, value := range inp.Properties() {
		inp.SetProperty(key, value)
	}
}

func (inp *Input) Render() (res string, err error) {
	inp.InitProps()

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

	if res, err = bc.TemplateBuilder("input", tpl, funcMap, inp); err == nil && inp.EventURL != "" {
		bc.SetCMValue(inp.RequestMap, inp.Id, inp)
	}
	return res, nil
}

var demoInputResponse func(evt bc.ResponseEvent) (re bc.ResponseEvent) = func(evt bc.ResponseEvent) (re bc.ResponseEvent) {
	evt.Trigger.SetProperty("invalid", false)
	if evt.Value != "valid" {
		evt.Trigger.SetProperty("invalid", true)
	}
	return evt
}

func DemoInput(eventURL, parentID string) []bc.DemoComponent {
	return []bc.DemoComponent{
		{
			Label:         "Default and AutoFocus",
			ComponentType: bc.ComponentTypeInput,
			Component: &Input{
				Type:        InputTypeText,
				Placeholder: "placeholder text",
				AutoFocus:   true,
			}},
		{
			Label:         "Valid value: valid",
			ComponentType: bc.ComponentTypeInput,
			Component: &Input{
				BaseComponent: bc.BaseComponent{
					EventURL:   eventURL,
					Swap:       bc.SwapOuterHTML,
					OnResponse: demoInputResponse,
				},
				Type:  InputTypeText,
				Value: "valid",
			}},
		{
			Label:         "ReadOnly",
			ComponentType: bc.ComponentTypeInput,
			Component: &Input{
				Type:     InputTypeText,
				Value:    "hello",
				ReadOnly: true,
			}},
		{
			Label:         "Disabled",
			ComponentType: bc.ComponentTypeInput,
			Component: &Input{
				Type:     InputTypeText,
				Value:    "hello",
				Disabled: true,
			}},
		{
			Label:         "Password full",
			ComponentType: bc.ComponentTypeInput,
			Component: &Input{
				Type:  InputTypePassword,
				Value: "secret",
				Full:  true,
			}},
		{
			Label:         "File input",
			ComponentType: bc.ComponentTypeInput,
			Component: &Input{
				Type:   InputTypeFile,
				Accept: ".jpg,.png",
			}},
		{
			Label:         "Color input",
			ComponentType: bc.ComponentTypeInput,
			Component: &Input{
				Type:  InputTypeColor,
				Value: "#845185",
			}},
	}
}
