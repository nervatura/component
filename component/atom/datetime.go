package atom

import (
	"strings"
	"time"

	bc "github.com/nervatura/component/component/base"
)

const (
	DateTimeEventChange = "change"

	DateTimeTypeDate     = "date"
	DateTimeTypeTime     = "time"
	DateTimeTypeDateTime = "datetime-local"
)

var DateTimeType []string = []string{DateTimeTypeDate, DateTimeTypeTime, DateTimeTypeDateTime}

type DateTime struct {
	bc.BaseComponent
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

func (dti *DateTime) Properties() bc.IM {
	return bc.MergeIM(
		dti.BaseComponent.Properties(),
		bc.IM{
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
			return dti.CheckEnumValue(bc.ToString(propValue, ""), DateTimeTypeDate, DateTimeType)
		},
		"value": func() interface{} {
			dtype := dti.Validation("type", dti.Type).(string)
			value := bc.ToString(propValue, dti.defaultValue(dtype))
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
			dti.Label = bc.ToString(propValue, "")
			return dti.Label
		},
		"isnull": func() interface{} {
			dti.IsNull = bc.ToBoolean(propValue, false)
			return dti.IsNull
		},
		"picker": func() interface{} {
			dti.Picker = bc.ToBoolean(propValue, false)
			return dti.Picker
		},
		"disabled": func() interface{} {
			dti.Disabled = bc.ToBoolean(propValue, false)
			return dti.Disabled
		},
		"readonly": func() interface{} {
			dti.ReadOnly = bc.ToBoolean(propValue, false)
			return dti.ReadOnly
		},
		"auto_focus": func() interface{} {
			dti.AutoFocus = bc.ToBoolean(propValue, false)
			return dti.AutoFocus
		},
		"full": func() interface{} {
			dti.Full = bc.ToBoolean(propValue, false)
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

func (dti *DateTime) OnRequest(te bc.TriggerEvent) (re bc.ResponseEvent) {
	value := dti.SetProperty("value", te.Values.Get(te.Name))
	evt := bc.ResponseEvent{
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

	if res, err = bc.TemplateBuilder("datetime", tpl, funcMap, dti); err == nil && dti.EventURL != "" {
		dti.SetProperty("request_map", dti)
	}
	return res, nil
}

func DemoDateTime(demo bc.ClientComponent) []bc.DemoComponent {
	id := bc.ToString(demo.GetProperty("id"), "")
	eventURL := bc.ToString(demo.GetProperty("event_url"), "")
	requestValue := demo.GetProperty("request_value").(map[string]bc.IM)
	requestMap := demo.GetProperty("request_map").(map[string]bc.ClientComponent)
	return []bc.DemoComponent{
		{
			Label:         "Default, focus, null",
			ComponentType: bc.ComponentTypeInput,
			Component: &DateTime{
				BaseComponent: bc.BaseComponent{
					Id: id + "_datetime_default",
				},
				Type:      DateTimeTypeDate,
				Value:     "",
				IsNull:    true,
				AutoFocus: true,
			}},
		{
			Label:         "DateTime, picker, full",
			ComponentType: bc.ComponentTypeInput,
			Component: &DateTime{
				BaseComponent: bc.BaseComponent{
					Id: id + "_datetime_picker",
				},
				Type:   DateTimeTypeDateTime,
				Picker: true,
			}},
		{
			Label:         "Time",
			ComponentType: bc.ComponentTypeInput,
			Component: &DateTime{
				BaseComponent: bc.BaseComponent{
					Id:           id + "_datetime_time",
					EventURL:     eventURL,
					Swap:         bc.SwapOuterHTML,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: DateTimeTypeTime,
			}},
	}
}
