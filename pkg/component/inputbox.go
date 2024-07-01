package component

import (
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [InputBox] constants
const (
	ComponentTypeInputBox = "inputbox"

	InputBoxEventOK          = "input_ok"
	InputBoxEventCancel      = "input_cancel"
	InputBoxEventValueChange = "input_value"
)

// Message and value request component
type InputBox struct {
	BaseComponent
	Value       string `json:"value"`
	Title       string `json:"title"`
	Message     string `json:"message"`
	Info        string `json:"info"`
	Tag         string `json:"tag"`
	LabelOK     string `json:"label_ok"`
	LabelCancel string `json:"label_cancel"`
	ShowValue   bool   `json:"show_value"`
	DefaultOK   bool   `json:"default_ok"`
}

/*
Returns all properties of the [InputBox]
*/
func (ibx *InputBox) Properties() ut.IM {
	return ut.MergeIM(
		ibx.BaseComponent.Properties(),
		ut.IM{
			"value":        ibx.Value,
			"title":        ibx.Title,
			"message":      ibx.Message,
			"info":         ibx.Info,
			"tag":          ibx.Tag,
			"label_ok":     ibx.LabelOK,
			"label_cancel": ibx.LabelCancel,
			"show_value":   ibx.ShowValue,
			"default_ok":   ibx.DefaultOK,
		})
}

/*
Returns the value of the property of the [InputBox] with the specified name.
*/
func (ibx *InputBox) GetProperty(propName string) interface{} {
	return ibx.Properties()[propName]
}

/*
It checks the value given to the property of the [Pagination] and always returns a valid value
*/
func (ibx *InputBox) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"target": func() interface{} {
			ibx.SetProperty("id", ibx.Id)
			value := ut.ToString(propValue, ibx.Id)
			if value != "this" && !strings.HasPrefix(value, "#") {
				value = "#" + value
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if ibx.BaseComponent.GetProperty(propName) != nil {
		return ibx.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [InputBox] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (ibx *InputBox) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			ibx.Value = ut.ToString(propValue, "")
			return ibx.Value
		},
		"title": func() interface{} {
			ibx.Title = ut.ToString(propValue, "Warning")
			return ibx.Title
		},
		"message": func() interface{} {
			ibx.Message = ut.ToString(propValue, "")
			return ibx.Message
		},
		"info": func() interface{} {
			ibx.Info = ut.ToString(propValue, "")
			return ibx.Info
		},
		"tag": func() interface{} {
			ibx.Tag = ut.ToString(propValue, "")
			return ibx.Tag
		},
		"label_cancel": func() interface{} {
			ibx.LabelCancel = ut.ToString(propValue, "Cancel")
			return ibx.LabelCancel
		},
		"label_ok": func() interface{} {
			ibx.LabelOK = ut.ToString(propValue, "OK")
			return ibx.LabelOK
		},
		"show_value": func() interface{} {
			ibx.ShowValue = ut.ToBoolean(propValue, false)
			return ibx.ShowValue
		},
		"default_ok": func() interface{} {
			ibx.DefaultOK = ut.ToBoolean(propValue, false)
			return ibx.DefaultOK
		},
		"target": func() interface{} {
			ibx.Target = ibx.Validation(propName, propValue).(string)
			return ibx.Target
		},
	}
	if _, found := pm[propName]; found {
		return ibx.SetRequestValue(propName, pm[propName](), []string{})
	}
	if ibx.BaseComponent.GetProperty(propName) != nil {
		return ibx.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (ibx *InputBox) response(evt ResponseEvent) (re ResponseEvent) {
	ibxEvt := ResponseEvent{Trigger: ibx, TriggerName: ibx.Name, Value: ibx.Tag}
	if ibx.ShowValue {
		if ibx.Tag != "" {
			ibxEvt.Value = ut.SM{"value": ibx.Value, "tag": ibx.Tag}
		} else {
			ibxEvt.Value = ibx.Value
		}
	}
	switch evt.TriggerName {
	case "btn_ok":
		ibxEvt.Name = InputBoxEventOK
	case "input_value":
		ibxEvt.Name = InputBoxEventValueChange
		ibx.SetProperty("value", evt.Value)
	default:
		ibxEvt.Name = InputBoxEventCancel
	}
	if ibx.OnResponse != nil {
		return ibx.OnResponse(ibxEvt)
	}
	return ibxEvt
}

func (ibx *InputBox) getComponent(name string) (res string, err error) {
	ccBtn := func(btnStyle, label, icon string, focus bool) *Button {
		return &Button{
			BaseComponent: BaseComponent{
				Id:           ibx.Id + "_" + name,
				Name:         name,
				EventURL:     ibx.EventURL,
				Target:       ibx.Target,
				OnResponse:   ibx.response,
				RequestValue: ibx.RequestValue,
				RequestMap:   ibx.RequestMap,
			},
			ButtonStyle: btnStyle,
			Label:       label,
			Icon:        icon,
			Full:        true,
			AutoFocus:   focus,
			Selected:    focus,
		}
	}
	ccMap := map[string]func() ClientComponent{
		"btn_ok": func() ClientComponent {
			return ccBtn(ButtonStylePrimary, ibx.LabelOK, "Check", ibx.DefaultOK)
		},
		"btn_cancel": func() ClientComponent {
			return ccBtn(ButtonStyleDefault, ibx.LabelCancel, "Times", false)
		},
		"input_value": func() ClientComponent {
			return &Input{
				BaseComponent: BaseComponent{
					Id:           ibx.Id + "_" + name,
					Name:         name,
					Style:        ut.SM{"border-radius": "0", "margin": "1px 0 2px"},
					EventURL:     ibx.EventURL,
					Target:       ibx.Target,
					OnResponse:   ibx.response,
					RequestValue: ibx.RequestValue,
					RequestMap:   ibx.RequestMap,
				},
				Type:      InputTypeString,
				Label:     ibx.Message,
				Value:     ibx.Value,
				AutoFocus: true,
				Full:      true,
			}
		},
	}
	cc := ccMap[name]()
	res, err = cc.Render()
	return res, err
}

/*
Based on the values, it will generate the html code of the [InputBox] or return with an error message.
*/
func (ibx *InputBox) Render() (res string, err error) {
	ibx.InitProps(ibx)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(ibx.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(ibx.Class, " ")
		},
		"inputComponent": func(name string) (string, error) {
			return ibx.getComponent(name)
		},
	}
	tpl := `<div id="{{ .Id }}" name="{{ .Name }}" class="row {{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	><div class="modal"><div class="dialog"><div class="panel">
	<div class="panel-title"><div class="cell title-cell"><span>{{ .Title }}</span></div></div>
	<div class="section" ><div class="row full container" >
	<div class="bold" >{{ .Message }}</div>
	{{ if ne .Info "" }}<div >{{ .Info }}</div>{{ end }}
	{{ if .ShowValue }}<div class="section-small-top" >{{ inputComponent "input_value" }}</div>{{ end }}
	</div></div>
	<div class="section buttons" ><div class="row full container" >
	<div class="cell padding-small half" >{{ inputComponent "btn_cancel" }}</div>
	<div class="cell padding-small half" >{{ inputComponent "btn_ok" }}</div>
	</div></div>
	</div></div></div>
	</div>`

	return ut.TemplateBuilder("inputbox", tpl, funcMap, ibx)
}

var testInputBoxResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
	data := evt.Trigger.GetProperty("data").(ut.IM)
	if (evt.Name == InputBoxEventOK) || (evt.Name == InputBoxEventCancel) {
		return ResponseEvent{
			Trigger: &Button{
				BaseComponent: BaseComponent{
					Id:           evt.Trigger.(*InputBox).Id,
					Data:         evt.Trigger.(*InputBox).Data,
					EventURL:     evt.Trigger.(*InputBox).EventURL,
					OnResponse:   evt.Trigger.(*InputBox).OnResponse,
					RequestValue: evt.Trigger.(*InputBox).RequestValue,
					RequestMap:   evt.Trigger.(*InputBox).RequestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "Input message",
			},
			TriggerName: evt.TriggerName,
			Name:        evt.Trigger.(*InputBox).Name,
		}
	}
	if evt.Name == ButtonEventClick {
		return ResponseEvent{
			Trigger: &InputBox{
				BaseComponent: BaseComponent{
					Id:           evt.Trigger.(*Button).Id,
					Data:         evt.Trigger.(*Button).Data,
					EventURL:     evt.Trigger.(*Button).EventURL,
					OnResponse:   evt.Trigger.(*Button).OnResponse,
					RequestValue: evt.Trigger.(*Button).RequestValue,
					RequestMap:   evt.Trigger.(*Button).RequestMap,
				},
				Value:     ut.ToString(data["value"], ""),
				Title:     ut.ToString(data["title"], ""),
				Message:   ut.ToString(data["message"], ""),
				Info:      ut.ToString(data["info"], ""),
				Tag:       ut.ToString(data["tag"], ""),
				ShowValue: ut.ToBoolean(data["show_value"], false),
				DefaultOK: ut.ToBoolean(data["default_ok"], false),
			},
			TriggerName: evt.TriggerName,
			Name:        evt.Trigger.(*Button).Name,
		}
	}
	return evt
}

// [InputBox] test and demo data
func TestInputBox(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "InputBox message",
			ComponentType: ComponentTypeInputBox,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id: id + "_inputbox_default",
					Data: ut.IM{
						"value": "", "title": "Warning",
						"message":    "The data has changed, but has not been saved!",
						"info":       "Save changes?",
						"tag":        "next_func",
						"default_ok": true, "show_value": false,
					},
					EventURL:     eventURL,
					OnResponse:   testInputBoxResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "Input message",
			}},
		{
			Label:         "InputBox value",
			ComponentType: ComponentTypeInputBox,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id: id + "_inputbox_default",
					Data: ut.IM{
						"value":      "default value",
						"title":      "New fieldname",
						"message":    "Enter the value:",
						"info":       "",
						"tag":        "next_func",
						"default_ok": false, "show_value": true,
					},
					EventURL:     eventURL,
					OnResponse:   testInputBoxResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "Input message",
			}},
	}
}
