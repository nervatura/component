package component

import (
	"strings"
	"time"

	ut "github.com/nervatura/component/pkg/util"
)

const (
	DateTimeEventChange = "change"

	DateTimeTypeDate     = "date"
	DateTimeTypeTime     = "time"
	DateTimeTypeDateTime = "datetime-local"
)

var DateTimeType []string = []string{DateTimeTypeDate, DateTimeTypeTime, DateTimeTypeDateTime}

type DateTime struct {
	BaseComponent
	Type      string `json:"type"`
	Value     string `json:"value"`
	Label     string `json:"label"`
	IsNull    bool   `json:"is_null"`
	Picker    bool   `json:"picker"`
	Disabled  bool   `json:"disabled"`
	ReadOnly  bool   `json:"read_only"`
	AutoFocus bool   `json:"auto_focus"`
	Full      bool   `json:"full"`
}

func (dti *DateTime) Properties() ut.IM {
	return ut.MergeIM(
		dti.BaseComponent.Properties(),
		ut.IM{
			"type":       dti.Type,
			"value":      dti.Value,
			"label":      dti.Label,
			"isnull":     dti.IsNull,
			"picker":     dti.Picker,
			"disabled":   dti.Disabled,
			"readonly":   dti.ReadOnly,
			"auto_focus": dti.AutoFocus,
			"full":       dti.Full,
		})
}

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
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if dti.BaseComponent.GetProperty(propName) != nil {
		return dti.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

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
		"isnull": func() interface{} {
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
	}
	if _, found := pm[propName]; found {
		return dti.SetRequestValue(propName, pm[propName](), []string{})
	}
	if dti.BaseComponent.GetProperty(propName) != nil {
		return dti.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (dti *DateTime) OnRequest(te TriggerEvent) (re ResponseEvent) {
	value := dti.SetProperty("value", te.Values.Get(te.Name))
	evt := ResponseEvent{
		Trigger: dti, TriggerName: dti.Name,
		Name:  DateTimeEventChange,
		Value: value,
	}
	if dti.OnResponse != nil {
		return dti.OnResponse(evt)
	}
	return evt
}

func (dti *DateTime) Render() (res string, err error) {
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
	{{ if ne .EventURL "" }} hx-post="{{ .EventURL }}" hx-target="{{ .Target }}" hx-swap="{{ .Swap }}"{{ end }}
	{{ if ne .Indicator "" }} hx-indicator="#{{ .Indicator }}"{{ end }}
	{{ if .Picker }} onfocus="this.showPicker()"{{ end }}
	{{ if .ReadOnly }} readonly{{ end }}
	{{ if .Disabled }} disabled{{ end }}
	{{ if .AutoFocus }} autofocus{{ end }}
	{{ if ne .Label "" }} aria-label="{{ .Label }}"{{ end }}
	 class="{{ customClass }}{{ if .Full }} full{{ end }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	></input>`

	if res, err = ut.TemplateBuilder("datetime", tpl, funcMap, dti); err == nil && dti.EventURL != "" {
		dti.SetProperty("request_map", dti)
	}
	return res, nil
}

func TestDateTime(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default, focus, null",
			ComponentType: ComponentTypeInput,
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
			ComponentType: ComponentTypeInput,
			Component: &DateTime{
				BaseComponent: BaseComponent{
					Id: id + "_datetime_picker",
				},
				Type:   DateTimeTypeDateTime,
				Picker: true,
			}},
		{
			Label:         "Time",
			ComponentType: ComponentTypeInput,
			Component: &DateTime{
				BaseComponent: BaseComponent{
					Id:           id + "_datetime_time",
					EventURL:     eventURL,
					Swap:         SwapOuterHTML,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: DateTimeTypeTime,
			}},
	}
}
