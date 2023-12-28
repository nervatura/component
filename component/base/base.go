package base

import (
	"net/url"
	"strings"
)

const (
	BaseEventValue = "value"
)

// IM is a map[string]interface{} type short alias
type IM = map[string]interface{}

// SM is a map[string]string type short alias
type SM = map[string]string

type ClientComponent interface {
	Properties() IM
	Validation(propName string, propValue interface{}) interface{}
	GetProperty(propName string) interface{}
	SetProperty(propName string, propValue interface{}) interface{}
	InitProps()
	Render() (res string, err error)
	OnRequest(te TriggerEvent) (re ResponseEvent)
}

type DemoComponent struct {
	Label         string
	ComponentType string
	Component     ClientComponent
}

type TriggerEvent struct {
	Id     string
	Name   string
	Target string
	Values url.Values
}

type ResponseEvent struct {
	Trigger     ClientComponent
	TriggerName string
	Name        string
	Value       interface{}
	Header      SM
}

type BaseComponent struct {
	Id         string
	Name       string
	EventURL   string
	Target     string
	Swap       string
	Indicator  string
	Class      []string
	Style      SM
	Data       IM
	RequestMap map[string]ClientComponent
	OnResponse func(evt ResponseEvent) (re ResponseEvent)
}

func (bcc *BaseComponent) Properties() IM {
	return IM{
		"id":          bcc.Id,
		"name":        bcc.Name,
		"event_url":   bcc.EventURL,
		"target":      bcc.Target,
		"swap":        bcc.Swap,
		"indicator":   bcc.Indicator,
		"class":       bcc.Class,
		"style":       bcc.Style,
		"data":        bcc.Data,
		"request_map": bcc.RequestMap,
	}
}

func (bcc *BaseComponent) CheckEnumValue(value, defaultValue string, enums []string) string {
	if Contains(enums, value) {
		return value
	}
	return defaultValue
}

func (bcc *BaseComponent) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"id": func() interface{} {
			return ToString(propValue, GetComponentID())
		},
		"name": func() interface{} {
			nid := bcc.SetProperty("id", bcc.Id)
			return ToString(propValue, ToString(nid, ""))
		},
		"target": func() interface{} {
			value := ToString(propValue, "this")
			if value != "this" && !strings.HasPrefix(value, "#") {
				value = "#" + value
			}
			return value
		},
		"swap": func() interface{} {
			return bcc.CheckEnumValue(ToString(propValue, ""), SwapOuterHTML, Swap)
		},
		"indicator": func() interface{} {
			return bcc.CheckEnumValue(ToString(propValue, ""), IndicatorNone, Indicator)
		},
		"class": func() interface{} {
			value := []string{}
			if class, valid := propValue.([]string); valid {
				value = class
			}
			return value
		},
		"style": func() interface{} {
			value := SetSMValue(bcc.Style, "", "")
			if smap, valid := propValue.(SM); valid {
				value = MergeSM(value, smap)
			}
			return value
		},
		"data": func() interface{} {
			value := SetIMValue(bcc.Data, "", "")
			if imap, valid := propValue.(IM); valid {
				value = MergeIM(value, imap)
			}
			return value
		},
		"request_map": func() interface{} {
			value := SetCMValue(bcc.RequestMap, "", nil)
			if cmap, valid := propValue.(map[string]ClientComponent); valid {
				value = MergeCM(value, cmap)
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	return propValue
}

func (bcc *BaseComponent) GetProperty(propName string) interface{} {
	return bcc.Properties()[propName]
}

func (bcc *BaseComponent) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"id": func() interface{} {
			bcc.Id = bcc.Validation(propName, propValue).(string)
			return bcc.Id
		},
		"name": func() interface{} {
			bcc.Name = bcc.Validation(propName, propValue).(string)
			return bcc.Name
		},
		"indicator": func() interface{} {
			bcc.Indicator = bcc.Validation(propName, propValue).(string)
			return bcc.Indicator
		},
		"event_url": func() interface{} {
			bcc.EventURL = ToString(propValue, "")
			return bcc.EventURL
		},
		"target": func() interface{} {
			bcc.Target = bcc.Validation(propName, propValue).(string)
			return bcc.Target
		},
		"swap": func() interface{} {
			bcc.Swap = bcc.Validation(propName, propValue).(string)
			return bcc.Swap
		},
		"class": func() interface{} {
			bcc.Class = bcc.Validation(propName, propValue).([]string)
			return bcc.Class
		},
		"style": func() interface{} {
			bcc.Style = bcc.Validation(propName, propValue).(SM)
			return bcc.Style
		},
		"data": func() interface{} {
			bcc.Data = bcc.Validation(propName, propValue).(IM)
			return bcc.Data
		},
		"request_map": func() interface{} {
			bcc.RequestMap = bcc.Validation(propName, propValue).(map[string]ClientComponent)
			return bcc.RequestMap
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	return propValue
}

func (bcc *BaseComponent) InitProps() {
	for key, value := range bcc.Properties() {
		bcc.SetProperty(key, value)
	}
}

func (bcc *BaseComponent) Render() (res string, err error) {
	bcc.InitProps()

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(bcc.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(bcc.Class, " ")
		},
	}
	tpl := `<div id="{{ .Id }}" class="{{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}>
	</div>`

	return TemplateBuilder("base", tpl, funcMap, bcc)
}

func (bcc *BaseComponent) OnRequest(te TriggerEvent) (re ResponseEvent) {
	re = ResponseEvent{
		Trigger:     bcc,
		TriggerName: te.Name,
		Name:        BaseEventValue,
	}
	if bcc.OnResponse != nil {
		return bcc.OnResponse(re)
	}
	return re
}
