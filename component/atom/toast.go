package atom

import (
	"strings"

	bc "github.com/nervatura/component/component/base"
)

const (
	ToastTypeInfo    = "info"
	ToastTypeError   = "error"
	ToastTypeSuccess = "success"
)

var ToastType []string = []string{ToastTypeInfo, ToastTypeError, ToastTypeSuccess}

type Toast struct {
	bc.BaseComponent
	Type    string
	Value   string
	Timeout int64
}

func (tst *Toast) Properties() bc.IM {
	return bc.MergeIM(
		tst.BaseComponent.Properties(),
		bc.IM{
			"type":    tst.Type,
			"value":   tst.Value,
			"timeout": tst.Timeout,
		})
}

func (tst *Toast) GetProperty(propName string) interface{} {
	return tst.Properties()[propName]
}

func (tst *Toast) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"type": func() interface{} {
			return tst.CheckEnumValue(bc.ToString(propValue, ""), ToastTypeInfo, ToastType)
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if tst.BaseComponent.GetProperty(propName) != nil {
		return tst.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

func (tst *Toast) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"type": func() interface{} {
			tst.Type = tst.Validation(propName, propValue).(string)
			return tst.Type
		},
		"value": func() interface{} {
			tst.Value = bc.ToString(propValue, "")
			return tst.Value
		},
		"timeout": func() interface{} {
			tst.Timeout = bc.ToInteger(propValue, 0)
			return tst.Timeout
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if tst.BaseComponent.GetProperty(propName) != nil {
		return tst.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (tst *Toast) getComponent(name string) (string, error) {
	iconMap := bc.SM{
		"info":    "InfoCircle",
		"error":   "ExclamationTriangle",
		"success": "CheckSquare",
	}
	ccMap := map[string]func() bc.ClientComponent{
		"icon": func() bc.ClientComponent {
			return &Icon{
				BaseComponent: bc.BaseComponent{Style: bc.SM{"margin": "auto"}},
				Value:         iconMap[tst.Type],
				Width:         32,
				Height:        32,
			}
		},
	}
	return ccMap[name]().Render()
}

func (tst *Toast) InitProps() {
	for key, value := range tst.Properties() {
		tst.SetProperty(key, value)
	}
}

func (tst *Toast) Render() (res string, err error) {
	tst.InitProps()

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(tst.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(tst.Class, " ")
		},
		"toastComponent": func(name string) (string, error) {
			return tst.getComponent(name)
		},
	}
	tpl := `<div id="{{ .Id }}" name="{{ .Name }}" type="{{ .Type }}" 
	 class="toast {{ customClass }}" onclick="htmx.remove(htmx.find('#{{ .Id }}'))"
	{{ if gt .Timeout 0 }} remove-me="{{ .Timeout }}s"{{ end }}
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	><span class="toast-icon">{{ toastComponent "icon" }}</span>
	<span id="{{ .Id }}-value">{{ .Value }}</span>
	</div>`

	return bc.TemplateBuilder("toast", tpl, funcMap, tst)
}

var demoToastResponse func(evt bc.ResponseEvent) (re bc.ResponseEvent) = func(evt bc.ResponseEvent) (re bc.ResponseEvent) {
	data := evt.Trigger.GetProperty("data").(bc.IM)
	re = bc.ResponseEvent{
		Trigger: &Toast{
			Type:    bc.ToString(data["toast-type"], ""),
			Value:   bc.ToString(data["toast-value"], ""),
			Timeout: bc.ToInteger(data["toast-timeout"], 0),
		},
		TriggerName: evt.TriggerName,
		Name:        ButtonEventClick,
		Header: bc.SM{
			bc.HeaderRetarget: "#toast-msg",
			bc.HeaderReswap:   "innerHTML",
		},
	}
	return re
}

func DemoToast(eventURL, parentID string) []bc.DemoComponent {
	return []bc.DemoComponent{
		{
			Label:         "Info toast message",
			ComponentType: bc.ComponentTypeToast,
			Component: &Button{
				BaseComponent: bc.BaseComponent{
					Data: bc.IM{
						"toast-type": "info", "toast-value": "This is an info message.", "toast-timeout": "4",
					},
					EventURL:   eventURL,
					OnResponse: demoToastResponse,
				},
				Type:  ButtonTypeDefault,
				Align: bc.TextAlignCenter,
				Label: "Info message",
			}},

		{
			Label:         "Error message without timeout",
			ComponentType: bc.ComponentTypeToast,
			Component: &Button{
				BaseComponent: bc.BaseComponent{
					Data: bc.IM{
						"toast-type": "error", "toast-value": "<i>This is an error message.</i>", "toast-timeout": "0",
					},
					EventURL:   eventURL,
					OnResponse: demoToastResponse,
				},
				Type:  ButtonTypeDefault,
				Align: bc.TextAlignCenter,
				Label: "Error message",
			}},
		{
			Label:         "A long success message",
			ComponentType: bc.ComponentTypeToast,
			Component: &Button{
				BaseComponent: bc.BaseComponent{
					Data: bc.IM{
						"toast-type":    "success",
						"toast-value":   "This is an success message. This is an success message. This is an success message. This is an success message",
						"toast-timeout": "4",
					},
					EventURL:   eventURL,
					OnResponse: demoToastResponse,
				},
				Type:  ButtonTypeDefault,
				Align: bc.TextAlignCenter,
				Label: "Success message",
			}},
	}
}
