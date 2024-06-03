package component

import (
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [Toggle] constants
const (
	ComponentTypeToggle = "toggle"

	ToggleEventChange = "change"
)

// Creates an toggle or checkbox control
type Toggle struct {
	BaseComponent
	Value bool `json:"value"`
	// Checkbox or toggle icons
	CheckBox bool `json:"check_box"`
	// Cell border
	Border bool `json:"border"`
	// Full width cell (100%)
	Full bool `json:"full"`
	// Specifies that the input should be disabled
	Disabled bool `json:"disabled"`
}

/*
Returns all properties of the [Toggle]
*/
func (tgl *Toggle) Properties() ut.IM {
	return ut.MergeIM(
		tgl.BaseComponent.Properties(),
		ut.IM{
			"value":     tgl.Value,
			"check_box": tgl.CheckBox,
			"border":    tgl.Border,
			"full":      tgl.Full,
			"disabled":  tgl.Disabled,
		})
}

/*
Returns the value of the property of the [Toggle] with the specified name.
*/
func (tgl *Toggle) GetProperty(propName string) interface{} {
	return tgl.Properties()[propName]
}

/*
Setting a property of the [Toggle] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (tgl *Toggle) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			tgl.Value = ut.ToBoolean(propValue, false)
			return tgl.Value
		},
		"check_box": func() interface{} {
			tgl.CheckBox = ut.ToBoolean(propValue, false)
			return tgl.CheckBox
		},
		"border": func() interface{} {
			tgl.Border = ut.ToBoolean(propValue, false)
			return tgl.Border
		},
		"full": func() interface{} {
			tgl.Full = ut.ToBoolean(propValue, false)
			return tgl.Full
		},
		"disabled": func() interface{} {
			tgl.Disabled = ut.ToBoolean(propValue, false)
			return tgl.Disabled
		},
	}
	if _, found := pm[propName]; found {
		return tgl.SetRequestValue(propName, pm[propName](), []string{})
	}
	if tgl.BaseComponent.GetProperty(propName) != nil {
		return tgl.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

/*
If the OnResponse function of the [Toggle] is implemented, the function calls it after the [TriggerEvent]
is processed, otherwise the function's return [ResponseEvent] is the processed [TriggerEvent].
*/
func (tgl *Toggle) OnRequest(te TriggerEvent) (re ResponseEvent) {
	value := tgl.SetProperty("value", !tgl.Value)
	evt := ResponseEvent{
		Trigger:     tgl,
		TriggerName: tgl.Name,
		Name:        ToggleEventChange,
		Value:       value,
	}
	if tgl.OnResponse != nil {
		return tgl.OnResponse(evt)
	}
	return evt
}

func (tgl *Toggle) getComponent(name string) (string, error) {
	ccMap := map[string]func() ClientComponent{
		"toggleOn": func() ClientComponent {
			return &Icon{
				Value: "ToggleOn", Width: 40, Height: 32.6,
			}
		},
		"toggleOff": func() ClientComponent {
			return &Icon{
				Value: "ToggleOff", Width: 40, Height: 32.6,
			}
		},
		"check": func() ClientComponent {
			return &Icon{
				Value: "CheckSquare", Width: 32, Height: 32,
			}
		},
		"empty": func() ClientComponent {
			return &Icon{
				Value: "SquareEmpty", Width: 32, Height: 32,
			}
		},
	}
	return ccMap[name]().Render()
}

/*
Based on the values, it will generate the html code of the [Toggle] or return with an error message.
*/
func (tgl *Toggle) Render() (res string, err error) {
	tgl.InitProps(tgl)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(tgl.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(tgl.Class, " ")
		},
		"toggleComponent": func(name string) (string, error) {
			return tgl.getComponent(name)
		},
	}
	tpl := `<div id="{{ .Id }}" name="{{ .Name }}" 
	{{ if eq .Disabled false }}{{ if ne .EventURL "" }} hx-post="{{ .EventURL }}" hx-target="{{ .Target }}" hx-swap="{{ .Swap }}"{{ end }}
	{{ if ne .Indicator "" }} hx-indicator="#{{ .Indicator }}"{{ end }}{{ end }}
	 class="toggle {{ customClass }}{{ if .Full }} full{{ end }}{{ if .Disabled }} toggle-disabled{{ end }}
	{{ if .Border }} toggle-border{{ end }}{{ if .Value }} toggle-on{{ else }} toggle-off{{ end }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	>{{ if .CheckBox }}{{ if .Value }}{{ toggleComponent "check" }}{{ else }}{{ toggleComponent "empty" }}{{ end }}
	{{ else }}{{ if .Value }}{{ toggleComponent "toggleOn" }}{{ else }}{{ toggleComponent "toggleOff" }}{{ end }}{{ end }}
	</div>`

	if res, err = ut.TemplateBuilder("toggle", tpl, funcMap, tgl); err == nil && tgl.EventURL != "" {
		tgl.SetProperty("request_map", tgl)
	}
	return res, nil
}

// [Toggle] test and demo data
func TestToggle(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeToggle,
			Component: &Toggle{
				BaseComponent: BaseComponent{
					Id:           id + "_toggle_default",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value:    false,
				CheckBox: false,
				Border:   false,
				Full:     false,
				Disabled: false,
			}},
		{
			Label:         "Border and disabled",
			ComponentType: ComponentTypeToggle,
			Component: &Toggle{
				BaseComponent: BaseComponent{
					Id:           id + "_toggle_border",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value:    true,
				CheckBox: false,
				Border:   true,
				Full:     true,
				Disabled: true,
			}},
		{
			Label:         "CheckBox",
			ComponentType: ComponentTypeToggle,
			Component: &Toggle{
				BaseComponent: BaseComponent{
					Id:           id + "_toggle_check",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value:    true,
				CheckBox: true,
				Border:   false,
				Full:     false,
				Disabled: false,
			}},
		{
			Label:         "CheckBox false",
			ComponentType: ComponentTypeToggle,
			Component: &Toggle{
				BaseComponent: BaseComponent{
					Id:           id + "_toggle_check_false",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value:    false,
				CheckBox: true,
				Border:   true,
				Full:     false,
				Disabled: false,
			}},
	}
}
