package component

import (
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [Toast] constants
const (
	ComponentTypeToast = "toast"

	ToastTypeInfo    = "info"
	ToastTypeError   = "error"
	ToastTypeSuccess = "success"
)

// [Toast] ToastType values
var ToastType []string = []string{ToastTypeInfo, ToastTypeError, ToastTypeSuccess}

/*
Creates a toast message control

For example:

	&Toast{
	  Type:    ToastTypeInfo,
	  Value:   "Info text",
	  Timeout: 4,
	}
*/
type Toast struct {
	BaseComponent
	/* [ToastType] variable constants: [ToastTypeInfo], [ToastTypeError], [ToastTypeSuccess].
	Default value: [ToastTypeInfo] */
	Type string `json:"type"`
	// Message text value
	Value string `json:"value"`
	/* Allows you to remove the element after a specified interval
	Its value must be given in seconds. Default value: 0*/
	Timeout int64 `json:"timeout"`
}

/*
Returns all properties of the [Toast]
*/
func (tst *Toast) Properties() ut.IM {
	return ut.MergeIM(
		tst.BaseComponent.Properties(),
		ut.IM{
			"type":    tst.Type,
			"value":   tst.Value,
			"timeout": tst.Timeout,
		})
}

/*
Returns the value of the property of the [Toast] with the specified name.
*/
func (tst *Toast) GetProperty(propName string) interface{} {
	return tst.Properties()[propName]
}

/*
It checks the value given to the property of the [Toast] and always returns a valid value
*/
func (tst *Toast) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"type": func() interface{} {
			return tst.CheckEnumValue(ut.ToString(propValue, ""), ToastTypeInfo, ToastType)
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

/*
Setting a property of the [Toast] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (tst *Toast) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"type": func() interface{} {
			tst.Type = tst.Validation(propName, propValue).(string)
			return tst.Type
		},
		"value": func() interface{} {
			tst.Value = ut.ToString(propValue, "")
			return tst.Value
		},
		"timeout": func() interface{} {
			tst.Timeout = ut.ToInteger(propValue, 0)
			return tst.Timeout
		},
	}
	if _, found := pm[propName]; found {
		return tst.SetRequestValue(propName, pm[propName](), []string{})
	}
	if tst.BaseComponent.GetProperty(propName) != nil {
		return tst.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (tst *Toast) getComponent(name string) (string, error) {
	iconMap := ut.SM{
		"info":    "InfoCircle",
		"error":   "ExclamationTriangle",
		"success": "CheckSquare",
	}
	ccMap := map[string]func() ClientComponent{
		"icon": func() ClientComponent {
			return &Icon{
				BaseComponent: BaseComponent{Style: ut.SM{"margin": "auto"}},
				Value:         iconMap[tst.Type],
				Width:         32,
				Height:        32,
			}
		},
	}
	return ccMap[name]().Render()
}

/*
Based on the values, it will generate the html code of the [Toast] or return with an error message.
*/
func (tst *Toast) Render() (res string, err error) {
	tst.InitProps(tst)

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

	return ut.TemplateBuilder("toast", tpl, funcMap, tst)
}

var testToastResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
	data := evt.Trigger.GetProperty("data").(ut.IM)
	re = ResponseEvent{
		Trigger: &Toast{
			Type:    ut.ToString(data["toast-type"], ""),
			Value:   ut.ToString(data["toast-value"], ""),
			Timeout: ut.ToInteger(data["toast-timeout"], 0),
		},
		TriggerName: evt.TriggerName,
		Name:        ButtonEventClick,
		Header: ut.SM{
			HeaderRetarget: "#toast-msg",
			HeaderReswap:   SwapInnerHTML,
		},
	}
	return re
}

// [Toast] test and demo data
func TestToast(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Info toast message",
			ComponentType: ComponentTypeToast,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id: id + "_toast_default",
					Data: ut.IM{
						"toast-type": "info", "toast-value": "This is an info message.", "toast-timeout": "4",
					},
					EventURL:     eventURL,
					OnResponse:   testToastResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "Info message",
			}},

		{
			Label:         "Error message without timeout",
			ComponentType: ComponentTypeToast,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id: id + "_toast_error",
					Data: ut.IM{
						"toast-type": "error", "toast-value": "<i>This is an error message.</i>", "toast-timeout": "0",
					},
					EventURL:     eventURL,
					OnResponse:   testToastResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "Error message",
			}},
		{
			Label:         "A long success message",
			ComponentType: ComponentTypeToast,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id: id + "_toast_success",
					Data: ut.IM{
						"toast-type":    "success",
						"toast-value":   "This is an success message. This is an success message. This is an success message. This is an success message",
						"toast-timeout": "4",
					},
					EventURL:     eventURL,
					OnResponse:   testToastResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "Success message",
			}},
	}
}
