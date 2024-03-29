package component

import (
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

const (
	LabelEventClick = "click"
)

type Label struct {
	BaseComponent
	Value     string `json:"value"`
	Centered  bool   `json:"centered"`
	LeftIcon  string `json:"left_icon"`
	RightIcon string `json:"right_icon"`
	IconStyle ut.SM  `json:"icon_style"`
}

func (lbl *Label) Properties() ut.IM {
	return ut.MergeIM(
		lbl.BaseComponent.Properties(),
		ut.IM{
			"value":      lbl.Value,
			"centered":   lbl.Centered,
			"left_icon":  lbl.LeftIcon,
			"right_icon": lbl.RightIcon,
			"icon_style": lbl.IconStyle,
		})
}

func (lbl *Label) GetProperty(propName string) interface{} {
	return lbl.Properties()[propName]
}

func (lbl *Label) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"icon_style": func() interface{} {
			value := ut.SetSMValue(lbl.IconStyle, "", "")
			if smap, valid := propValue.(ut.SM); valid {
				value = ut.MergeSM(value, smap)
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if lbl.BaseComponent.GetProperty(propName) != nil {
		return lbl.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

func (lbl *Label) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			lbl.Value = ut.ToString(propValue, "")
			return lbl.Value
		},
		"centered": func() interface{} {
			lbl.Centered = ut.ToBoolean(propValue, false)
			return lbl.Centered
		},
		"left_icon": func() interface{} {
			lbl.LeftIcon = ut.ToString(propValue, "")
			return lbl.LeftIcon
		},
		"right_icon": func() interface{} {
			lbl.RightIcon = ut.ToString(propValue, "")
			return lbl.RightIcon
		},
		"icon_style": func() interface{} {
			lbl.IconStyle = lbl.Validation(propName, propValue).(ut.SM)
			return lbl.IconStyle
		},
	}
	if _, found := pm[propName]; found {
		return lbl.SetRequestValue(propName, pm[propName](), []string{})
	}
	if lbl.BaseComponent.GetProperty(propName) != nil {
		return lbl.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (lbl *Label) OnRequest(te TriggerEvent) (re ResponseEvent) {
	evt := ResponseEvent{
		Trigger:     lbl,
		TriggerName: lbl.Name,
		Name:        LabelEventClick,
		Value:       lbl.Value,
	}
	if lbl.OnResponse != nil {
		return lbl.OnResponse(evt)
	}
	return evt
}

func (lbl *Label) getComponent(name string) (string, error) {
	ccMap := map[string]func() ClientComponent{
		"left_icon": func() ClientComponent {
			return &Icon{
				BaseComponent: BaseComponent{Style: lbl.IconStyle},
				Value:         lbl.LeftIcon,
				Width:         20,
				Color:         ut.ToString(lbl.IconStyle["color"], ""),
			}
		},
		"right_icon": func() ClientComponent {
			return &Icon{
				BaseComponent: BaseComponent{Style: lbl.IconStyle},
				Value:         lbl.RightIcon,
				Width:         20,
				Color:         ut.ToString(lbl.IconStyle["color"], ""),
			}
		},
	}
	return ccMap[name]().Render()
}

func (lbl *Label) Render() (res string, err error) {
	lbl.InitProps(lbl)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(lbl.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(lbl.Class, " ")
		},
		"labelComponent": func(name string) (string, error) {
			return lbl.getComponent(name)
		},
	}
	tpl := `{{ if or (ne .LeftIcon "") (ne .RightIcon "") }}<div id="{{ .Id }}" name="{{ .Name }}"
	{{ if ne .EventURL "" }} hx-post="{{ .EventURL }}" hx-target="{{ .Target }}" hx-swap="{{ .Swap }}"{{ end }}
	{{ if ne .Indicator "" }} hx-indicator="#{{ .Indicator }}"{{ end }}
	 class="label row{{ if ne .EventURL "" }} label_link{{ else }} label_text{{ end }}{{ if and (ne .LeftIcon "") (.Centered) }} centered{{ end }} {{ customClass }}"
	>{{ if ne .LeftIcon "" }}
	<div class="cell label_icon_left">{{ labelComponent "left_icon" }}</div>
	<div class="cell label_info_left bold"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	>{{ .Value }}</div>
	{{ else }}
	<div class="cell label_info_right bold"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	>{{ .Value }}</div>
	<div class="cell label_icon_right">{{ labelComponent "right_icon" }}</div>
	{{ end }}</div>
	{{ else }}
	<span id="{{ .Id }}" name="{{ .Name }}"
	{{ if ne .EventURL "" }} hx-post="{{ .EventURL }}" hx-target="{{ .Target }}" hx-swap="{{ .Swap }}"{{ end }}
	{{ if ne .Indicator "" }} hx-indicator="#{{ .Indicator }}"{{ end }}
	 class="label bold{{ if ne .EventURL "" }} label_link{{ else }} label_text{{ end }}"
	 {{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	>{{ .Value }}</span>{{ end }}`

	if res, err = ut.TemplateBuilder("label", tpl, funcMap, lbl); err == nil && lbl.EventURL != "" {
		lbl.SetProperty("request_map", lbl)
	}
	return res, nil
}

var demoLblResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
	re = ResponseEvent{
		Trigger: &Toast{
			Type:    ToastTypeInfo,
			Value:   ut.ToString(evt.Trigger.GetProperty("data").(ut.IM)["toast_value"], ""),
			Timeout: 4,
		},
		TriggerName: evt.TriggerName,
		Name:        LabelEventClick,
		Header: ut.SM{
			HeaderRetarget: "#toast-msg",
			HeaderReswap:   "innerHTML",
		},
	}
	return re
}

func TestLabel(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default label",
			ComponentType: ComponentTypeLabel,
			Component: &Label{
				BaseComponent: BaseComponent{
					Id: id + "_label_default",
				},
				Value: "Label",
			}},
		{
			Label:         "Left icon",
			ComponentType: ComponentTypeLabel,
			Component: &Label{
				BaseComponent: BaseComponent{
					Id: id + "_label_left_icon",
				},
				Value:    "Label",
				LeftIcon: "InfoCircle",
			}},
		{
			Label:         "Right icon",
			ComponentType: ComponentTypeLabel,
			Component: &Label{
				BaseComponent: BaseComponent{
					Id: id + "_label_right_icon",
				},
				Value:     "Label",
				RightIcon: "InfoCircle",
			}},
		{
			Label:         "Centered",
			ComponentType: ComponentTypeLabel,
			Component: &Label{
				BaseComponent: BaseComponent{
					Id: id + "_label_centered",
				},
				Value:    "Label",
				LeftIcon: "InfoCircle",
				Centered: true,
			}},
		{
			Label:         "Label style",
			ComponentType: ComponentTypeLabel,
			Component: &Label{
				BaseComponent: BaseComponent{
					Id:    id + "_label_style",
					Style: ut.SM{"color": "red"},
				},
				Value:     "Label",
				LeftIcon:  "InfoCircle",
				IconStyle: ut.SM{"fill": "orange"},
			}},
		{
			Label:         "Label link",
			ComponentType: ComponentTypeLabel,
			Component: &Label{
				BaseComponent: BaseComponent{
					Id:       id + "_label_toast",
					EventURL: eventURL,
					Data: ut.IM{
						"toast_value": "Link value",
					},
					OnResponse:   demoLblResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value:    "Label link",
				LeftIcon: "Globe",
			}},
	}
}
