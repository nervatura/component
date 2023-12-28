package atom

import (
	"strings"

	bc "github.com/nervatura/component/component/base"
)

const (
	LabelEventClick = "click"
)

type Label struct {
	bc.BaseComponent
	Value     string
	Centered  bool
	LeftIcon  string
	RightIcon string
	IconStyle bc.SM
}

func (lbl *Label) Properties() bc.IM {
	return bc.MergeIM(
		lbl.BaseComponent.Properties(),
		bc.IM{
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
			value := bc.SetSMValue(lbl.IconStyle, "", "")
			if smap, valid := propValue.(bc.SM); valid {
				value = bc.MergeSM(value, smap)
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
			lbl.Value = bc.ToString(propValue, "")
			return lbl.Value
		},
		"centered": func() interface{} {
			lbl.Centered = bc.ToBoolean(propValue, false)
			return lbl.Centered
		},
		"left_icon": func() interface{} {
			lbl.LeftIcon = bc.ToString(propValue, "")
			return lbl.LeftIcon
		},
		"right_icon": func() interface{} {
			lbl.RightIcon = bc.ToString(propValue, "")
			return lbl.RightIcon
		},
		"icon_style": func() interface{} {
			lbl.IconStyle = lbl.Validation(propName, propValue).(bc.SM)
			return lbl.IconStyle
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if lbl.BaseComponent.GetProperty(propName) != nil {
		return lbl.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (lbl *Label) OnRequest(te bc.TriggerEvent) (re bc.ResponseEvent) {
	evt := bc.ResponseEvent{
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
	ccMap := map[string]func() bc.ClientComponent{
		"left_icon": func() bc.ClientComponent {
			return &Icon{
				BaseComponent: bc.BaseComponent{Style: lbl.IconStyle},
				Value:         lbl.LeftIcon,
				Width:         20,
				Color:         bc.ToString(lbl.IconStyle["color"], ""),
			}
		},
		"right_icon": func() bc.ClientComponent {
			return &Icon{
				BaseComponent: bc.BaseComponent{Style: lbl.IconStyle},
				Value:         lbl.RightIcon,
				Width:         20,
				Color:         bc.ToString(lbl.IconStyle["color"], ""),
			}
		},
	}
	return ccMap[name]().Render()
}

func (lbl *Label) InitProps() {
	for key, value := range lbl.Properties() {
		lbl.SetProperty(key, value)
	}
}

func (lbl *Label) Render() (res string, err error) {
	lbl.InitProps()

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

	if res, err = bc.TemplateBuilder("label", tpl, funcMap, lbl); err == nil && lbl.EventURL != "" {
		bc.SetCMValue(lbl.RequestMap, lbl.Id, lbl)
	}
	return res, nil
}

var demoLblResponse func(evt bc.ResponseEvent) (re bc.ResponseEvent) = func(evt bc.ResponseEvent) (re bc.ResponseEvent) {
	re = bc.ResponseEvent{
		Trigger: &Toast{
			Type:    ToastTypeInfo,
			Value:   bc.ToString(evt.Trigger.GetProperty("data").(bc.IM)["toast_value"], ""),
			Timeout: 4,
		},
		TriggerName: evt.TriggerName,
		Name:        LabelEventClick,
		Header: bc.SM{
			bc.HeaderRetarget: "#toast-msg",
			bc.HeaderReswap:   "innerHTML",
		},
	}
	return re
}

func DemoLabel(eventURL, parentID string) []bc.DemoComponent {
	return []bc.DemoComponent{
		{
			Label:         "Default label",
			ComponentType: bc.ComponentTypeLabel,
			Component: &Label{
				Value: "Label",
			}},
		{
			Label:         "Left icon",
			ComponentType: bc.ComponentTypeLabel,
			Component: &Label{
				Value:    "Label",
				LeftIcon: "InfoCircle",
			}},
		{
			Label:         "Right icon",
			ComponentType: bc.ComponentTypeLabel,
			Component: &Label{
				Value:     "Label",
				RightIcon: "InfoCircle",
			}},
		{
			Label:         "Centered",
			ComponentType: bc.ComponentTypeLabel,
			Component: &Label{
				Value:    "Label",
				LeftIcon: "InfoCircle",
				Centered: true,
			}},
		{
			Label:         "Label style",
			ComponentType: bc.ComponentTypeLabel,
			Component: &Label{
				BaseComponent: bc.BaseComponent{
					Style: bc.SM{"color": "red"},
				},
				Value:     "Label",
				LeftIcon:  "InfoCircle",
				IconStyle: bc.SM{"fill": "orange"},
			}},
		{
			Label:         "Label link",
			ComponentType: bc.ComponentTypeLabel,
			Component: &Label{
				BaseComponent: bc.BaseComponent{
					Id:       bc.GetComponentID(),
					EventURL: eventURL,
					Data: bc.IM{
						"toast_value": "Link value",
					},
					OnResponse: demoLblResponse,
				},
				Value:    "Label link",
				LeftIcon: "Globe",
			}},
	}
}
