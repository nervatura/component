package component

import (
	"html/template"
	"strings"
	"time"

	ut "github.com/nervatura/component/pkg/util"
)

// [DateTime] constants
const (
	ComponentTypeDateTime = "datetime"

	DateTimeEventChange = "datetime_change"

	DateTimeTypeDate     = "date"
	DateTimeTypeTime     = "time"
	DateTimeTypeDateTime = "datetime-local"
)

// [DateTime] Type values
var DateTimeType []string = []string{DateTimeTypeDate, DateTimeTypeTime, DateTimeTypeDateTime}

/*
Creates an HTML date, datetime or time input control

For example:

	&DateTime{
	  BaseComponent: BaseComponent{
	    Id: "id_datetime_picker",
	  },
	  Type:   DateTimeTypeDateTime,
	  Picker: true,
	}
*/
type DateTime struct {
	BaseComponent
	/* [DateTimeType] variable constants: [DateTimeTypeDate], [DateTimeTypeTime], [DateTimeTypeDateTime].
	Default value: [DateTimeTypeDate] */
	Type string `json:"type"`
	// Any valid value based on control type
	Value string `json:"value"`
	// The HTML aria-label attribute of the component
	Label string `json:"label"`
	// Allows entry of an empty value
	IsNull bool `json:"is_null"`
	// Show value picker when input
	Picker bool `json:"picker"`
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
Returns all properties of the [DateTime]
*/
func (dti *DateTime) Properties() ut.IM {
	return ut.MergeIM(
		dti.BaseComponent.Properties(),
		ut.IM{
			"type":       dti.Type,
			"value":      dti.Value,
			"label":      dti.Label,
			"is_null":    dti.IsNull,
			"picker":     dti.Picker,
			"disabled":   dti.Disabled,
			"readonly":   dti.ReadOnly,
			"auto_focus": dti.AutoFocus,
			"full":       dti.Full,
		})
}

/*
Returns the value of the property of the [DateTime] with the specified name.
*/
func (dti *DateTime) GetProperty(propName string) interface{} {
	return dti.Properties()[propName]
}

func (dti *DateTime) defaultValue(dtype string) (value string) {
	if dti.IsNull {
		return ""
	}
	switch dtype {
	case DateTimeTypeTime:
		value = time.Now().Format("15:04")
	case DateTimeTypeDate:
		value = time.Now().Format("2006-01-02")
	default:
		value = time.Now().Format("2006-01-02T15:04")
	}
	return value
}

/*
It checks the value given to the property of the [DateTime] and always returns a valid value
*/
func (dti *DateTime) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"type": func() interface{} {
			return dti.CheckEnumValue(ut.ToString(propValue, ""), DateTimeTypeDate, DateTimeType)
		},
		"value": func() interface{} {
			dtype := dti.Validation("type", dti.Type).(string)
			value := ut.ToString(propValue, dti.defaultValue(dtype))
			valueLength := map[string]int{
				DateTimeTypeTime: 5, DateTimeTypeDate: 10, DateTimeTypeDateTime: 16,
			}
			if len(value) > valueLength[dtype] {
				value = value[:valueLength[dtype]]
			}
			return value
		},
		"swap": func() interface{} {
			return dti.CheckEnumValue(ut.ToString(propValue, ""), SwapInnerHTML, Swap)
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if dti.BaseComponent.GetProperty(propName) != nil {
		return dti.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [DateTime] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (dti *DateTime) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"type": func() interface{} {
			dti.Type = dti.Validation(propName, propValue).(string)
			return dti.Type
		},
		"value": func() interface{} {
			dti.Value = dti.Validation(propName, propValue).(string)
			return dti.Value
		},
		"label": func() interface{} {
			dti.Label = ut.ToString(propValue, "")
			return dti.Label
		},
		"is_null": func() interface{} {
			dti.IsNull = ut.ToBoolean(propValue, false)
			return dti.IsNull
		},
		"picker": func() interface{} {
			dti.Picker = ut.ToBoolean(propValue, false)
			return dti.Picker
		},
		"disabled": func() interface{} {
			dti.Disabled = ut.ToBoolean(propValue, false)
			return dti.Disabled
		},
		"readonly": func() interface{} {
			dti.ReadOnly = ut.ToBoolean(propValue, false)
			return dti.ReadOnly
		},
		"auto_focus": func() interface{} {
			dti.AutoFocus = ut.ToBoolean(propValue, false)
			return dti.AutoFocus
		},
		"full": func() interface{} {
			dti.Full = ut.ToBoolean(propValue, false)
			return dti.Full
		},
		"swap": func() interface{} {
			dti.Swap = dti.Validation(propName, propValue).(string)
			return dti.Swap
		},
	}
	if _, found := pm[propName]; found {
		return dti.SetRequestValue(propName, pm[propName](), []string{})
	}
	if dti.BaseComponent.GetProperty(propName) != nil {
		return dti.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

/*
If the OnResponse function of the [DateTime] is implemented, the function calls it after the [TriggerEvent]
is processed, otherwise the function's return [ResponseEvent] is the processed [TriggerEvent].
*/
func (dti *DateTime) OnRequest(te TriggerEvent) (re ResponseEvent) {
	value := dti.SetProperty("value", te.Values.Get(te.Name))
	evt := ResponseEvent{
		Trigger: dti, TriggerName: dti.Name,
		Name:   DateTimeEventChange,
		Value:  value,
		Header: map[string]string{},
	}
	if (value != te.Values.Get(te.Name)) && (dti.Swap == SwapInnerHTML) {
		evt.Header[HeaderReswap] = SwapOuterHTML
	}
	if dti.OnResponse != nil {
		return dti.OnResponse(evt)
	}
	return evt
}

/*
Based on the values, it will generate the html code of the [DateTime] or return with an error message.
*/
func (dti *DateTime) Render() (html template.HTML, err error) {
	dti.InitProps(dti)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(dti.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(dti.Class, " ")
		},
	}
	tpl := `<input id="{{ .Id }}" name="{{ .Name }}" type="{{ .Type }}" value="{{ .Value }}"
	{{ if ne .EventURL "" }} hx-post="{{ .EventURL }}" hx-trigger="blur, keyup[keyCode==13]" hx-target="{{ .Target }}" {{ if ne .Sync "none" }} hx-sync="{{ .Sync }}"{{ end }} hx-swap="{{ .Swap }}"{{ end }}
	{{ if ne .Indicator "none" }} hx-indicator="#{{ .Indicator }}"{{ end }}
	{{ if .Picker }} onfocus="this.showPicker()"{{ end }}
	{{ if .ReadOnly }} readonly{{ end }}
	{{ if .Disabled }} disabled{{ end }}
	{{ if .AutoFocus }} autofocus{{ end }}
	{{ if ne .Label "" }} aria-label="{{ .Label }}"{{ end }}
	 class="{{ customClass }}{{ if .Full }} full{{ end }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	></input>`

	if html, err = ut.TemplateBuilder("datetime", tpl, funcMap, dti); err == nil && dti.EventURL != "" {
		dti.SetProperty("request_map", dti)
	}
	return html, nil
}

// [DateTime] test and demo data
func TestDateTime(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default, focus, null",
			ComponentType: ComponentTypeDateTime,
			Component: &DateTime{
				BaseComponent: BaseComponent{
					Id: id + "_datetime_default",
				},
				Type:      DateTimeTypeDate,
				Value:     "",
				IsNull:    true,
				AutoFocus: true,
			}},
		{
			Label:         "DateTime, picker, full",
			ComponentType: ComponentTypeDateTime,
			Component: &DateTime{
				BaseComponent: BaseComponent{
					Id: id + "_datetime_picker",
				},
				Type:   DateTimeTypeDateTime,
				Picker: true,
			}},
		{
			Label:         "Time",
			ComponentType: ComponentTypeDateTime,
			Component: &DateTime{
				BaseComponent: BaseComponent{
					Id:           id + "_datetime_time",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: DateTimeTypeTime,
			}},
	}
}
