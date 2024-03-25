package component

import (
	"net/url"
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

const (
	BaseEventValue = "value"
)

type ClientComponent interface {
	Properties() ut.IM
	Validation(propName string, propValue interface{}) interface{}
	GetProperty(propName string) interface{}
	SetProperty(propName string, propValue interface{}) interface{}
	InitProps(cc ClientComponent)
	Render() (res string, err error)
	OnRequest(te TriggerEvent) (re ResponseEvent)
}

type TestComponent struct {
	Label         string          `json:"label"`
	ComponentType string          `json:"component_type"`
	Component     ClientComponent `json:"component"`
}

type TriggerEvent struct {
	Id     string     `json:"id"`
	Name   string     `json:"name"`
	Target string     `json:"target"`
	Values url.Values `json:"values"`
}

type ResponseEvent struct {
	Trigger     ClientComponent `json:"trigger"`
	TriggerName string          `json:"trigger_name"`
	Name        string          `json:"name"`
	Value       interface{}     `json:"value"`
	Header      ut.SM           `json:"header"`
}

type BaseComponent struct {
	Id           string                                     `json:"id"`
	Name         string                                     `json:"name"`
	EventURL     string                                     `json:"event_url"`
	Target       string                                     `json:"target"`
	Swap         string                                     `json:"swap"`
	Indicator    string                                     `json:"indicator"`
	Class        []string                                   `json:"class"`
	Style        ut.SM                                      `json:"style"`
	Data         ut.IM                                      `json:"data"`
	RequestValue map[string]ut.IM                           `json:"request_value"`
	RequestMap   map[string]ClientComponent                 `json:"-"`
	OnResponse   func(evt ResponseEvent) (re ResponseEvent) `json:"-"`
	init         bool
}

func (bcc *BaseComponent) Properties() ut.IM {
	return ut.IM{
		"id":            bcc.Id,
		"name":          bcc.Name,
		"event_url":     bcc.EventURL,
		"target":        bcc.Target,
		"swap":          bcc.Swap,
		"indicator":     bcc.Indicator,
		"class":         bcc.Class,
		"style":         bcc.Style,
		"data":          bcc.Data,
		"request_value": bcc.RequestValue,
		"request_map":   bcc.RequestMap,
	}
}

func (bcc *BaseComponent) CheckEnumValue(value, defaultValue string, enums []string) string {
	if ut.Contains(enums, value) {
		return value
	}
	return defaultValue
}

func (bcc *BaseComponent) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"id": func() interface{} {
			return ut.ToString(propValue, ut.GetComponentID())
		},
		"name": func() interface{} {
			nid := bcc.SetProperty("id", bcc.Id)
			return ut.ToString(propValue, ut.ToString(nid, ""))
		},
		"target": func() interface{} {
			value := ut.ToString(propValue, "this")
			if value != "this" && !strings.HasPrefix(value, "#") {
				value = "#" + value
			}
			return value
		},
		"swap": func() interface{} {
			return bcc.CheckEnumValue(ut.ToString(propValue, ""), SwapOuterHTML, Swap)
		},
		"indicator": func() interface{} {
			return bcc.CheckEnumValue(ut.ToString(propValue, ""), IndicatorNone, Indicator)
		},
		"class": func() interface{} {
			value := []string{}
			if class, valid := propValue.([]string); valid {
				value = class
			}
			return value
		},
		"style": func() interface{} {
			value := ut.SetSMValue(bcc.Style, "", "")
			if smap, valid := propValue.(ut.SM); valid {
				value = ut.MergeSM(value, smap)
			}
			return value
		},
		"data": func() interface{} {
			value := ut.SetIMValue(bcc.Data, "", "")
			if imap, valid := propValue.(ut.IM); valid {
				value = ut.MergeIM(value, imap)
			}
			return value
		},
		"request_value": func() interface{} {
			if bcc.RequestValue == nil {
				bcc.RequestValue = make(map[string]ut.IM)
			}
			if _, found := bcc.RequestValue[bcc.Id]; !found && (bcc.Id != "") {
				bcc.RequestValue[bcc.Id] = ut.IM{}
			}
			value := bcc.RequestValue
			if rvmap, valid := propValue.(map[string]ut.IM); valid {
				if ccmap, found := rvmap[bcc.Id]; found {
					value[bcc.Id] = ut.MergeIM(value[bcc.Id], ccmap)
				}
			}
			return value
		},
		"request_map": func() interface{} {
			if bcc.RequestMap == nil {
				bcc.RequestMap = map[string]ClientComponent{}
			}
			if _, valid := propValue.(ClientComponent); valid && (bcc.Id != "") {
				return true
			}
			return false
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
			bcc.EventURL = ut.ToString(propValue, "")
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
			bcc.Style = bcc.Validation(propName, propValue).(ut.SM)
			return bcc.Style
		},
		"data": func() interface{} {
			bcc.Data = bcc.Validation(propName, propValue).(ut.IM)
			return bcc.Data
		},
		"request_map": func() interface{} {
			if bcc.Validation(propName, propValue).(bool) {
				bcc.RequestMap[bcc.Id] = propValue.(ClientComponent)
			}
			return bcc.RequestMap
		},
	}
	if _, found := pm[propName]; found {
		return bcc.SetRequestValue(propName, pm[propName](), []string{"request_map"})
	}
	return propValue
}

func (bcc *BaseComponent) SetRequestValue(propName string, propValue interface{}, staticFields []string) interface{} {
	if !ut.Contains(staticFields, propName) && bcc.Id != "" && !bcc.init {
		bcc.RequestValue = bcc.Validation("request_value", map[string]ut.IM{bcc.Id: {propName: propValue}}).(map[string]ut.IM)
	}
	return propValue
}

func (bcc *BaseComponent) InitProps(cc ClientComponent) {
	bcc.init = true
	for key, value := range cc.Properties() {
		cc.SetProperty(key, value)
	}
	requestValue := cc.Validation("request_value", bcc.RequestValue).(map[string]ut.IM)
	if rq, found := requestValue[bcc.Id]; found {
		for key, value := range rq {
			cc.SetProperty(key, value)
		}
	}
	bcc.init = false
}

func (bcc *BaseComponent) Render() (res string, err error) {
	bcc.InitProps(bcc)

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

	return ut.TemplateBuilder("base", tpl, funcMap, bcc)
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
