package component

import (
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

const (
	ButtonEventClick = "click"

	ButtonTypeDefault = ""
	ButtonTypePrimary = "primary"
	ButtonTypeBorder  = "border"
)

var ButtonType []string = []string{ButtonTypeDefault, ButtonTypePrimary, ButtonTypeBorder}

type Button struct {
	BaseComponent
	Type           string          `json:"type"`
	Align          string          `json:"align"`
	Label          string          `json:"label"`
	LabelComponent ClientComponent `json:"label_component"`
	Icon           string          `json:"icon"`
	Disabled       bool            `json:"disabled"`
	AutoFocus      bool            `json:"auto_focus"`
	Full           bool            `json:"full"`
	Small          bool            `json:"small"`
	Selected       bool            `json:"selected"`
	HideLabel      bool            `json:"hide_label"`
	Badge          int64           `json:"badge"`
	ShowBadge      bool            `json:"show_badge"`
}

func (btn *Button) Properties() ut.IM {
	return ut.MergeIM(
		btn.BaseComponent.Properties(),
		ut.IM{
			"type":            btn.Type,
			"align":           btn.Align,
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
			return btn.CheckEnumValue(ut.ToString(propValue, ""), ButtonTypeDefault, ButtonType)
		},
		"align": func() interface{} {
			return btn.CheckEnumValue(ut.ToString(propValue, ""), TextAlignCenter, TextAlign)
		},
		"indicator": func() interface{} {
			return btn.CheckEnumValue(ut.ToString(propValue, IndicatorSpinner), IndicatorSpinner, Indicator)
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
		"label": func() interface{} {
			btn.Label = ut.ToString(propValue, "")
			return btn.Label
		},
		"label_component": func() interface{} {
			if cc, valid := propValue.(ClientComponent); valid {
				btn.LabelComponent = cc
			}
			return btn.LabelComponent
		},
		"icon": func() interface{} {
			btn.Icon = ut.ToString(propValue, "")
			return btn.Icon
		},
		"disabled": func() interface{} {
			btn.Disabled = ut.ToBoolean(propValue, false)
			return btn.Disabled
		},
		"auto_focus": func() interface{} {
			btn.AutoFocus = ut.ToBoolean(propValue, false)
			return btn.AutoFocus
		},
		"full": func() interface{} {
			btn.Full = ut.ToBoolean(propValue, false)
			return btn.Full
		},
		"small": func() interface{} {
			btn.Small = ut.ToBoolean(propValue, false)
			return btn.Small
		},
		"selected": func() interface{} {
			btn.Selected = ut.ToBoolean(propValue, false)
			return btn.Selected
		},
		"hide_label": func() interface{} {
			btn.HideLabel = ut.ToBoolean(propValue, false)
			return btn.HideLabel
		},
		"badge": func() interface{} {
			btn.Badge = ut.ToInteger(propValue, 0)
			return btn.Badge
		},
		"show_badge": func() interface{} {
			btn.ShowBadge = ut.ToBoolean(propValue, false)
			return btn.ShowBadge
		},
		"indicator": func() interface{} {
			btn.Indicator = btn.Validation(propName, propValue).(string)
			return btn.Indicator
		},
	}
	if _, found := pm[propName]; found {
		return btn.SetRequestValue(propName, pm[propName](), []string{})
	}
	if btn.BaseComponent.GetProperty(propName) != nil {
		return btn.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (btn *Button) OnRequest(te TriggerEvent) (re ResponseEvent) {
	value := btn.SetProperty("value", te.Values.Get(te.Id+"_value"))
	evt := ResponseEvent{
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
	ccMap := map[string]func() ClientComponent{
		"icon": func() ClientComponent {
			return &Icon{Value: btn.Icon, Width: 20}
		},
		"label": func() ClientComponent {
			return btn.LabelComponent
		},
	}
	return ccMap[name]().Render()
}

func (btn *Button) Render() (res string, err error) {
	btn.InitProps(btn)

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
	{{ if or (eq .Type "primary") (eq .Type "border") }} button-type="{{ .Type }}"{{ end }}
	{{ if ne .EventURL "" }} hx-post="{{ .EventURL }}" hx-target="{{ .Target }}" hx-swap="{{ .Swap }}"{{ end }}
	{{ if ne .Indicator "" }} hx-indicator="#{{ .Indicator }}"{{ end }}
	{{ if .Disabled }} disabled{{ end }}
	{{ if .AutoFocus }} autofocus{{ end }}
	{{ if ne .Label "" }} aria-label="{{ .Label }}" title="{{ .Label }}"{{ end }}
	 class="{{ .Align }}{{ if .Small }} small-button{{ end }}{{ if .Full }} full{{ end }}{{ if .Selected }} selected{{ end }}{{ if .HideLabel }} hidelabel{{ end }} {{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	>{{ if and (ne .Icon "") (ne .Align "align-right") }}{{ buttonComponent "icon" }}{{ end }}
	{{ if .LabelComponent }}{{ buttonComponent "label" }}{{ else }}<span>{{ .Label }}</span>{{ end }}
	{{ if and (ne .Icon "") (eq .Align "align-right") }}{{ buttonComponent "icon" }}{{ end }}
	{{ if .ShowBadge }}<span class="right" ><span class="badge{{ if .Selected }} selected-badge{{ end }}" >{{ .Badge }}</span></span>{{ end }}
	</button>`

	if res, err = ut.TemplateBuilder("button", tpl, funcMap, btn); err == nil && btn.EventURL != "" {
		btn.SetProperty("request_map", btn)
	}
	return res, err
}

var demoBtnResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
	badge := ut.ToInteger(evt.Trigger.GetProperty("badge"), 0)
	evt.Trigger.SetProperty("badge", badge+1)
	evt.Trigger.SetProperty("show_badge", true)
	return evt
}

func TestButton(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeButton,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id:           id + "_button_default",
					EventURL:     eventURL,
					OnResponse:   demoBtnResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type:  ButtonTypeDefault,
				Align: TextAlignCenter,
				Label: "Default",
			}},
		{
			Label:         "Right icon",
			ComponentType: ComponentTypeButton,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id:           id + "_button_right_icon",
					EventURL:     eventURL,
					OnResponse:   demoBtnResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type:      ButtonTypeDefault,
				Align:     TextAlignRight,
				Label:     "Right icon",
				Icon:      "InfoCircle",
				HideLabel: true,
			}},
		{
			Label:         "Primary and icon",
			ComponentType: ComponentTypeButton,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id:           id + "_button_primary",
					EventURL:     eventURL,
					OnResponse:   demoBtnResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type:     ButtonTypePrimary,
				Align:    TextAlignCenter,
				Label:    "Primary",
				Icon:     "Check",
				Selected: true,
			}},
		{
			Label:         "Border and selected",
			ComponentType: ComponentTypeButton,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id:           id + "_button_border",
					EventURL:     eventURL,
					OnResponse:   demoBtnResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type:     ButtonTypeBorder,
				Align:    TextAlignCenter,
				Label:    "Border selected",
				Selected: true,
			}},
		{
			Label:         "Border full and badge",
			ComponentType: ComponentTypeButton,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id:           id + "_button_full",
					EventURL:     eventURL,
					OnResponse:   demoBtnResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type:      ButtonTypeBorder,
				Align:     TextAlignCenter,
				Label:     "Border full and badge",
				Full:      true,
				Badge:     0,
				ShowBadge: true,
			}},
		{
			Label:         "Small disabled",
			ComponentType: ComponentTypeButton,
			Component: &Button{
				Type:     ButtonTypeDefault,
				Align:    TextAlignCenter,
				Label:    "Small disabled",
				Small:    true,
				Disabled: true,
			}},
		{
			Label:         "Label component",
			ComponentType: ComponentTypeButton,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id:           id + "_button_label",
					EventURL:     eventURL,
					OnResponse:   demoBtnResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type:           ButtonTypeDefault,
				Align:          TextAlignCenter,
				Label:          "Label component",
				LabelComponent: &Icon{Value: "Print", Width: 32, Height: 32},
			}},
		{
			Label:         "Border and custom style",
			ComponentType: ComponentTypeButton,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id:           id + "_button_custom",
					EventURL:     eventURL,
					OnResponse:   demoBtnResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Style: ut.SM{
						"border-color": "green", "color": "red", "border-radius": "3px",
					},
				},
				Type:  ButtonTypeBorder,
				Align: TextAlignCenter,
				Label: "Button style",
			}},
	}
}
