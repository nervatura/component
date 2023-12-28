package atom

import (
	"strings"

	bc "github.com/nervatura/component/component/base"
)

const (
	ButtonEventClick = "click"

	ButtonTypeDefault = ""
	ButtonTypePrimary = "primary"
	ButtonTypeBorder  = "border"
)

var ButtonType []string = []string{ButtonTypeDefault, ButtonTypePrimary, ButtonTypeBorder}

type Button struct {
	bc.BaseComponent
	Type           string
	Align          string
	Value          string
	Label          string
	LabelComponent bc.ClientComponent
	Icon           string
	Disabled       bool
	AutoFocus      bool
	Full           bool
	Small          bool
	Selected       bool
	HideLabel      bool
	Badge          int64
	ShowBadge      bool
}

func (btn *Button) Properties() bc.IM {
	return bc.MergeIM(
		btn.BaseComponent.Properties(),
		bc.IM{
			"type":            btn.Type,
			"align":           btn.Align,
			"value":           btn.Value,
			"label":           btn.Label,
			"label_component": btn.LabelComponent,
			"icon":            btn.Icon,
			"disabled":        btn.Disabled,
			"auto_focus":      btn.AutoFocus,
			"full":            btn.Full,
			"small":           btn.Small,
			"selected":        btn.Selected,
			"hide_label":      btn.HideLabel,
			"badge":           btn.Badge,
			"show_badge":      btn.ShowBadge,
		})
}

func (btn *Button) GetProperty(propName string) interface{} {
	return btn.Properties()[propName]
}

func (btn *Button) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"type": func() interface{} {
			return btn.CheckEnumValue(bc.ToString(propValue, ""), ButtonTypeDefault, ButtonType)
		},
		"align": func() interface{} {
			return btn.CheckEnumValue(bc.ToString(propValue, ""), bc.TextAlignCenter, bc.TextAlign)
		},
		"indicator": func() interface{} {
			return btn.CheckEnumValue(bc.ToString(propValue, bc.IndicatorSpinner), bc.IndicatorSpinner, bc.Indicator)
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if btn.BaseComponent.GetProperty(propName) != nil {
		return btn.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

func (btn *Button) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"type": func() interface{} {
			btn.Type = btn.Validation(propName, propValue).(string)
			return btn.Type
		},
		"align": func() interface{} {
			btn.Align = btn.Validation(propName, propValue).(string)
			return btn.Align
		},
		"value": func() interface{} {
			btn.Value = bc.ToString(propValue, "")
			return btn.Value
		},
		"label": func() interface{} {
			btn.Label = bc.ToString(propValue, "")
			return btn.Label
		},
		"label_component": func() interface{} {
			if cc, valid := propValue.(bc.ClientComponent); valid {
				btn.LabelComponent = cc
			}
			return btn.LabelComponent
		},
		"icon": func() interface{} {
			btn.Icon = bc.ToString(propValue, "")
			return btn.Icon
		},
		"disabled": func() interface{} {
			btn.Disabled = bc.ToBoolean(propValue, false)
			return btn.Disabled
		},
		"auto_focus": func() interface{} {
			btn.AutoFocus = bc.ToBoolean(propValue, false)
			return btn.AutoFocus
		},
		"full": func() interface{} {
			btn.Full = bc.ToBoolean(propValue, false)
			return btn.Full
		},
		"small": func() interface{} {
			btn.Small = bc.ToBoolean(propValue, false)
			return btn.Small
		},
		"selected": func() interface{} {
			btn.Selected = bc.ToBoolean(propValue, false)
			return btn.Selected
		},
		"hide_label": func() interface{} {
			btn.HideLabel = bc.ToBoolean(propValue, false)
			return btn.HideLabel
		},
		"badge": func() interface{} {
			btn.Badge = bc.ToInteger(propValue, 0)
			return btn.Badge
		},
		"show_badge": func() interface{} {
			btn.ShowBadge = bc.ToBoolean(propValue, false)
			return btn.ShowBadge
		},
		"indicator": func() interface{} {
			btn.Indicator = btn.Validation(propName, propValue).(string)
			return btn.Indicator
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if btn.BaseComponent.GetProperty(propName) != nil {
		return btn.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (btn *Button) OnRequest(te bc.TriggerEvent) (re bc.ResponseEvent) {
	value := btn.SetProperty("value", te.Values.Get(te.Id+"_value"))
	evt := bc.ResponseEvent{
		Trigger: btn, TriggerName: btn.Name,
		Name:  ButtonEventClick,
		Value: value,
	}
	if btn.OnResponse != nil {
		return btn.OnResponse(evt)
	}
	return evt
}

func (btn *Button) getComponent(name string) (string, error) {
	ccMap := map[string]func() bc.ClientComponent{
		"icon": func() bc.ClientComponent {
			return &Icon{Value: btn.Icon, Width: 20}
		},
		"label": func() bc.ClientComponent {
			return btn.LabelComponent
		},
	}
	return ccMap[name]().Render()
}

func (btn *Button) InitProps() {
	for key, value := range btn.Properties() {
		btn.SetProperty(key, value)
	}
}

func (btn *Button) Render() (res string, err error) {
	btn.InitProps()

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(btn.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(btn.Class, " ")
		},
		"buttonComponent": func(name string) (string, error) {
			return btn.getComponent(name)
		},
	}
	tpl := `<button id="{{ .Id }}" name="{{ .Name }}"
	{{ if ne .Value "" }} hx-include="[name='{{ .Id }}_value']"{{ end }}
	{{ if or (eq .Type "primary") (eq .Type "border") }} button-type="{{ .Type }}"{{ end }}
	{{ if ne .EventURL "" }} hx-post="{{ .EventURL }}" hx-target="{{ .Target }}" hx-swap="{{ .Swap }}"{{ end }}
	{{ if ne .Indicator "" }} hx-indicator="#{{ .Indicator }}"{{ end }}
	{{ if .Disabled }} disabled{{ end }}
	{{ if .AutoFocus }} autofocus{{ end }}
	{{ if ne .Label "" }} aria-label="{{ .Label }}" title="{{ .Label }}"{{ end }}
	 class="{{ .Align }}{{ if .Small }} small-button{{ end }}{{ if .Full }} full{{ end }}{{ if .Selected }} selected{{ end }}{{ if .HideLabel }} hidelabel{{ end }} {{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	>{{ if ne .Value "" }}<input type="hidden" name="{{ .Id }}_value" value="{{ .Value }}" >{{ end }}
	{{ if and (ne .Icon "") (ne .Align "align-right") }}{{ buttonComponent "icon" }}{{ end }}
	{{ if .LabelComponent }}{{ buttonComponent "label" }}{{ else }}<span>{{ .Label }}</span>{{ end }}
	{{ if and (ne .Icon "") (eq .Align "align-right") }}{{ buttonComponent "icon" }}{{ end }}
	{{ if .ShowBadge }}<span class="right" ><span class="badge{{ if .Selected }} selected-badge{{ end }}" >{{ .Badge }}</span></span>{{ end }}
	</button>`

	if res, err = bc.TemplateBuilder("button", tpl, funcMap, btn); err == nil && btn.EventURL != "" {
		bc.SetCMValue(btn.RequestMap, btn.Id, btn)
	}
	return res, err
}

var demoBtnResponse func(evt bc.ResponseEvent) (re bc.ResponseEvent) = func(evt bc.ResponseEvent) (re bc.ResponseEvent) {
	badge := bc.ToInteger(evt.Trigger.GetProperty("badge"), 0)
	evt.Trigger.SetProperty("badge", badge+1)
	evt.Trigger.SetProperty("show_badge", true)
	return evt
}

func DemoButton(eventURL, parentID string) []bc.DemoComponent {
	return []bc.DemoComponent{
		{
			Label:         "Default",
			ComponentType: bc.ComponentTypeButton,
			Component: &Button{
				BaseComponent: bc.BaseComponent{
					EventURL:   eventURL,
					OnResponse: demoBtnResponse,
				},
				Type:  ButtonTypeDefault,
				Align: bc.TextAlignCenter,
				Label: "Default",
				Value: "default",
			}},
		{
			Label:         "Right icon",
			ComponentType: bc.ComponentTypeButton,
			Component: &Button{
				BaseComponent: bc.BaseComponent{
					EventURL:   eventURL,
					OnResponse: demoBtnResponse,
				},
				Type:      ButtonTypeDefault,
				Align:     bc.TextAlignRight,
				Label:     "Right icon",
				Icon:      "InfoCircle",
				HideLabel: true,
			}},
		{
			Label:         "Primary and icon",
			ComponentType: bc.ComponentTypeButton,
			Component: &Button{
				BaseComponent: bc.BaseComponent{
					EventURL:   eventURL,
					OnResponse: demoBtnResponse,
				},
				Type:     ButtonTypePrimary,
				Align:    bc.TextAlignCenter,
				Label:    "Primary",
				Icon:     "Check",
				Selected: true,
			}},
		{
			Label:         "Border and selected",
			ComponentType: bc.ComponentTypeButton,
			Component: &Button{
				BaseComponent: bc.BaseComponent{
					EventURL:   eventURL,
					OnResponse: demoBtnResponse,
				},
				Type:     ButtonTypeBorder,
				Align:    bc.TextAlignCenter,
				Label:    "Border selected",
				Selected: true,
			}},
		{
			Label:         "Border full and badge",
			ComponentType: bc.ComponentTypeButton,
			Component: &Button{
				BaseComponent: bc.BaseComponent{
					EventURL:   eventURL,
					OnResponse: demoBtnResponse,
				},
				Type:      ButtonTypeBorder,
				Align:     bc.TextAlignCenter,
				Label:     "Border full and badge",
				Full:      true,
				Badge:     0,
				ShowBadge: true,
			}},
		{
			Label:         "Small disabled",
			ComponentType: bc.ComponentTypeButton,
			Component: &Button{
				Type:     ButtonTypeDefault,
				Align:    bc.TextAlignCenter,
				Label:    "Small disabled",
				Small:    true,
				Disabled: true,
			}},
		{
			Label:         "Label component",
			ComponentType: bc.ComponentTypeButton,
			Component: &Button{
				BaseComponent: bc.BaseComponent{
					EventURL:   eventURL,
					OnResponse: demoBtnResponse,
				},
				Type:           ButtonTypeDefault,
				Align:          bc.TextAlignCenter,
				Label:          "Label component",
				LabelComponent: &Icon{Value: "Print", Width: 32, Height: 32},
			}},
		{
			Label:         "Border and custom style",
			ComponentType: bc.ComponentTypeButton,
			Component: &Button{
				BaseComponent: bc.BaseComponent{
					EventURL:   eventURL,
					OnResponse: demoBtnResponse,
					Style: bc.SM{
						"border-color": "green", "color": "red", "border-radius": "3px",
					},
				},
				Type:  ButtonTypeBorder,
				Align: bc.TextAlignCenter,
				Label: "Button style",
			}},
	}
}
