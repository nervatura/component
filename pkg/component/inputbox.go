package component

import (
	"html/template"
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [InputBox] constants
const (
	ComponentTypeInputBox = "inputbox"

	InputBoxEventOK     = "inputbox_ok"
	InputBoxEventCancel = "inputbox_cancel"

	InputBoxTypeCancel   = "IBOX_CANCEL"
	InputBoxTypeOK       = "IBOX_OK"
	InputBoxTypeString   = "IBOX_STRING"
	InputBoxTypeText     = "IBOX_TEXT"
	InputBoxTypeColor    = "IBOX_COLOR"
	InputBoxTypeSelect   = "IBOX_SELECT"
	InputBoxTypeNumber   = "IBOX_NUMBER"
	InputBoxTypeInteger  = "IBOX_INTEGER"
	InputBoxTypeDate     = "IBOX_DATE"
	InputBoxTypeTime     = "IBOX_TIME"
	InputBoxTypeDateTime = "IBOX_DATETIME"
)

// [InputBox] Type values
var InputBoxType []string = []string{InputBoxTypeCancel, InputBoxTypeOK,
	InputBoxTypeString, InputBoxTypeText, InputBoxTypeColor, InputBoxTypeSelect,
	InputBoxTypeNumber, InputBoxTypeInteger, InputBoxTypeDate, InputBoxTypeTime, InputBoxTypeDateTime}

// Message and value request component
type InputBox struct {
	BaseComponent
	/* [InputBoxType] variable constants: [InputBoxTypeCancel], [InputBoxTypeOK], [InputBoxTypeInput],
	[InputBoxTypeSelect], [InputBoxTypeNumber], [InputBoxTypeInteger], [InputBoxTypeDate],
	[InputBoxTypeTime], [InputBoxTypeDateTime].
	Default value: [InputBoxTypeCancel] */
	InputType string `json:"input_type"`
	Value     string `json:"value"`
	/* [SelectOption] type values. */
	ValueOptions []SelectOption `json:"value_options"`
	Title        string         `json:"title"`
	Message      string         `json:"message"`
	Info         string         `json:"info"`
	Tag          string         `json:"tag"`
	LabelOK      string         `json:"label_ok"`
	LabelCancel  string         `json:"label_cancel"`
	DefaultOK    bool           `json:"default_ok"`
}

/*
Returns all properties of the [InputBox]
*/
func (ibx *InputBox) Properties() ut.IM {
	return ut.MergeIM(
		ibx.BaseComponent.Properties(),
		ut.IM{
			"input_type":    ibx.InputType,
			"value":         ibx.Value,
			"value_options": ibx.ValueOptions,
			"title":         ibx.Title,
			"message":       ibx.Message,
			"info":          ibx.Info,
			"tag":           ibx.Tag,
			"label_ok":      ibx.LabelOK,
			"label_cancel":  ibx.LabelCancel,
			"default_ok":    ibx.DefaultOK,
		})
}

/*
Returns the value of the property of the [InputBox] with the specified name.
*/
func (ibx *InputBox) GetProperty(propName string) interface{} {
	return ibx.Properties()[propName]
}

/*
It checks the value given to the property of the [InputBox] and always returns a valid value
*/
func (ibx *InputBox) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"input_type": func() interface{} {
			return ibx.CheckEnumValue(ut.ToString(propValue, ""), InputBoxTypeCancel, InputBoxType)
		},
		"value_options": func() interface{} {
			return SelectOptionRangeValidation(propValue, []SelectOption{})
		},
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
		"input_type": func() interface{} {
			ibx.InputType = ibx.Validation(propName, propValue).(string)
			return ibx.InputType
		},
		"value": func() interface{} {
			ibx.Value = ut.ToString(propValue, "")
			return ibx.Value
		},
		"value_options": func() interface{} {
			ibx.ValueOptions = ibx.Validation(propName, propValue).([]SelectOption)
			return ibx.ValueOptions
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

/*
If the OnResponse function of the [InputBox] is implemented, the function calls it after the [TriggerEvent]
is processed, otherwise the function's return [ResponseEvent] is the processed [TriggerEvent].
*/
func (ibx *InputBox) OnRequest(te TriggerEvent) (re ResponseEvent) {
	evt := ResponseEvent{
		Trigger: ibx, TriggerName: ibx.Name,
		Name: InputBoxEventCancel,
	}
	if te.Values.Has("btn_ok") {
		evt.Name = InputBoxEventOK
		if te.Values.Has("value") {
			ibx.SetProperty("value", te.Values.Get("value"))
		}
	}
	evt.Value = ut.SM{"value": ibx.Value, "tag": ibx.Tag}
	if ibx.OnResponse != nil {
		return ibx.OnResponse(evt)
	}
	return evt
}

func (ibx *InputBox) getComponent(name string) (html template.HTML, err error) {
	ccBtn := func(btnStyle, label, icon string, focus bool) *Button {
		return &Button{
			BaseComponent: BaseComponent{
				Id:   ibx.Id + "_" + name,
				Name: name,
			},
			Type:        ButtonTypeSubmit,
			ButtonStyle: btnStyle,
			Label:       label,
			Icon:        icon,
			Full:        true,
			AutoFocus:   focus,
			Selected:    focus,
		}
	}
	ccInp := func(value, typeValue string) *Input {
		inp := &Input{
			BaseComponent: BaseComponent{
				Id:    ibx.Id + "_" + name,
				Name:  "value",
				Style: ut.SM{"border-radius": "0", "margin": "1px 0 2px"},
			},
			Type:      typeValue,
			Label:     ibx.Message,
			Rows:      8,
			AutoFocus: true,
			Full:      true,
		}
		inp.SetProperty("value", value)
		return inp
	}
	ccNum := func(value float64, integer bool) *NumberInput {
		inp := &NumberInput{
			BaseComponent: BaseComponent{
				Id:    ibx.Id + "_" + name,
				Name:  "value",
				Style: ut.SM{"border-radius": "0", "margin": "1px 0 2px"},
			},
			Integer:   integer,
			Value:     value,
			AutoFocus: true,
		}
		return inp
	}
	ccSel := func(value string) *Select {
		sel := &Select{
			BaseComponent: BaseComponent{
				Id:    ibx.Id + "_" + name,
				Name:  "value",
				Style: ut.SM{"border-radius": "0", "margin": "1px 0 2px"},
			},
			IsNull:    false,
			AutoFocus: true,
			Options:   ibx.ValueOptions,
		}
		sel.SetProperty("value", value)
		return sel
	}
	ccDti := func(dateType, value string) *DateTime {
		dti := &DateTime{
			BaseComponent: BaseComponent{
				Id:    ibx.Id + "_" + name,
				Name:  "value",
				Style: ut.SM{"border-radius": "0", "margin": "1px 0 2px"},
			},
			Type:      dateType,
			Value:     value,
			AutoFocus: true,
			IsNull:    false,
		}
		return dti
	}
	ccMap := map[string]func() ClientComponent{
		"btn_ok": func() ClientComponent {
			return ccBtn(ButtonStylePrimary, ibx.LabelOK, IconCheck, ibx.DefaultOK)
		},
		"btn_cancel": func() ClientComponent {
			return ccBtn(ButtonStyleDefault, ibx.LabelCancel, IconTimes, false)
		},
		"string_value": func() ClientComponent {
			return ccInp(ibx.Value, InputTypeString)
		},
		"text_value": func() ClientComponent {
			return ccInp(ibx.Value, InputTypeText)
		},
		"color_value": func() ClientComponent {
			return ccInp(ibx.Value, InputTypeColor)
		},
		"select_value": func() ClientComponent {
			return ccSel(ibx.Value)
		},
		"number_value": func() ClientComponent {
			return ccNum(ut.ToFloat(ibx.Value, 0), false)
		},
		"integer_value": func() ClientComponent {
			return ccNum(ut.ToFloat(ut.ToInteger(ibx.Value, 0), 0), true)
		},
		"date_value": func() ClientComponent {
			return ccDti(DateTimeTypeDate, ibx.Value)
		},
		"time_value": func() ClientComponent {
			return ccDti(DateTimeTypeTime, ibx.Value)
		},
		"datetime_value": func() ClientComponent {
			return ccDti(DateTimeTypeDateTime, ibx.Value)
		},
	}
	cc := ccMap[name]()
	html, err = cc.Render()
	return html, err
}

/*
Based on the values, it will generate the html code of the [InputBox] or return with an error message.
*/
func (ibx *InputBox) Render() (html template.HTML, err error) {
	ibx.InitProps(ibx)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(ibx.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(ibx.Class, " ")
		},
		"inputComponent": func(name string) (template.HTML, error) {
			return ibx.getComponent(name)
		},
	}
	tpl := `<div id="{{ .Id }}" name="{{ .Name }}" class="row {{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}>
	<form id="{{ .Id }}" name="inputbox_form"
	{{ if ne .EventURL "" }} hx-post="{{ .EventURL }}" hx-target="{{ .Target }}" {{ if ne .Sync "none" }} hx-sync="{{ .Sync }}"{{ end }} hx-swap="{{ .Swap }}"{{ end }}
	{{ if ne .Indicator "none" }} hx-indicator="#{{ .Indicator }}"{{ end }}
	><div class="modal"><div class="dialog"><div class="panel">
	<div class="panel-title"><div class="cell title-cell"><span>{{ .Title }}</span></div></div>
	<div class="section" ><div class="row full container" >
	<div class="bold" >{{ .Message }}</div>
	{{ if ne .Info "" }}<div >{{ .Info }}</div>{{ end }}
	{{ if eq .InputType "IBOX_STRING" }}<div class="section-small-top" >{{ inputComponent "string_value" }}</div>{{ end }}
	{{ if eq .InputType "IBOX_TEXT" }}<div class="section-small-top" >{{ inputComponent "text_value" }}</div>{{ end }}
	{{ if eq .InputType "IBOX_COLOR" }}<div class="section-small-top" >{{ inputComponent "color_value" }}</div>{{ end }}
	{{ if eq .InputType "IBOX_SELECT" }}<div class="section-small-top" >{{ inputComponent "select_value" }}</div>{{ end }}
	{{ if eq .InputType "IBOX_NUMBER" }}<div class="section-small-top" >{{ inputComponent "number_value" }}</div>{{ end }}
	{{ if eq .InputType "IBOX_INTEGER" }}<div class="section-small-top" >{{ inputComponent "integer_value" }}</div>{{ end }}
	{{ if eq .InputType "IBOX_DATE" }}<div class="section-small-top" >{{ inputComponent "date_value" }}</div>{{ end }}
	{{ if eq .InputType "IBOX_TIME" }}<div class="section-small-top" >{{ inputComponent "time_value" }}</div>{{ end }}
	{{ if eq .InputType "IBOX_DATETIME" }}<div class="section-small-top" >{{ inputComponent "datetime_value" }}</div>{{ end }}
	</div></div>
	<div class="section buttons" ><div class="row full container" >
	<div class="cell padding-small {{ if ne .InputType "IBOX_OK" }}half{{ end }}" >{{ inputComponent "btn_ok" }}</div>
	{{ if ne .InputType "IBOX_OK" }}<div class="cell padding-small half" >{{ inputComponent "btn_cancel" }}</div>{{ end }}
	</div></div>
	</div></div></div>
	</form></div>`

	if html, err = ut.TemplateBuilder("inputbox", tpl, funcMap, ibx); err == nil && ibx.EventURL != "" {
		ibx.SetProperty("request_map", ibx)
	}
	return html, err
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
	re = ResponseEvent{
		Trigger: &InputBox{
			BaseComponent: BaseComponent{
				Id:           evt.Trigger.(*Button).Id,
				Data:         evt.Trigger.(*Button).Data,
				EventURL:     evt.Trigger.(*Button).EventURL,
				OnResponse:   evt.Trigger.(*Button).OnResponse,
				RequestValue: evt.Trigger.(*Button).RequestValue,
				RequestMap:   evt.Trigger.(*Button).RequestMap,
			},
			InputType:    ut.ToString(data["input_type"], ""),
			Value:        ut.ToString(data["value"], ""),
			ValueOptions: SelectOptionRangeValidation(data["value_options"], []SelectOption{}),
			Title:        ut.ToString(data["title"], ""),
			Message:      ut.ToString(data["message"], ""),
			Info:         ut.ToString(data["info"], ""),
			Tag:          ut.ToString(data["tag"], ""),
			DefaultOK:    ut.ToBoolean(data["default_ok"], false),
		},
		TriggerName: evt.TriggerName,
		Name:        evt.Trigger.(*Button).Name,
	}
	return re
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
						"input_type": InputBoxTypeCancel,
						"value":      "", "title": "Warning",
						"message":    "The data has changed, but has not been saved!",
						"info":       "Save changes?",
						"tag":        "next_func",
						"default_ok": true,
					},
					EventURL:     eventURL,
					OnResponse:   testInputBoxResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "OK and cancel message",
			}},
		{
			Label:         "InputBox info",
			ComponentType: ComponentTypeInputBox,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id: id + "_inputbox_info",
					Data: ut.IM{
						"input_type": InputBoxTypeOK,
						"title":      "Warning",
						"message":    "Info message",
						"info":       "Info message text",
						"default_ok": true,
					},
					EventURL:     eventURL,
					OnResponse:   testInputBoxResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "Info message",
			}},
		{
			Label:         "InputBox string value",
			ComponentType: ComponentTypeInputBox,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id: id + "_inputbox_string",
					Data: ut.IM{
						"input_type": InputBoxTypeString,
						"value":      "default value",
						"title":      "New fieldname",
						"message":    "Enter the value:",
						"info":       "",
						"tag":        "next_func",
						"default_ok": false,
					},
					EventURL:     eventURL,
					OnResponse:   testInputBoxResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "Input string message",
			}},
		{
			Label:         "InputBox text value",
			ComponentType: ComponentTypeInputBox,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id: id + "_inputbox_text",
					Data: ut.IM{
						"input_type": InputBoxTypeText,
						"value":      "default value",
						"title":      "New fieldname",
						"message":    "Enter the value:",
						"info":       "",
						"tag":        "next_func",
						"default_ok": false,
					},
					EventURL:     eventURL,
					OnResponse:   testInputBoxResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "Input text message",
			}},
		{
			Label:         "InputBox color value",
			ComponentType: ComponentTypeInputBox,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id: id + "_inputbox_color",
					Data: ut.IM{
						"input_type": InputBoxTypeColor,
						"value":      "#456200",
						"title":      "New fieldname",
						"message":    "Enter the value:",
						"info":       "",
						"tag":        "next_func",
						"default_ok": false,
					},
					EventURL:     eventURL,
					OnResponse:   testInputBoxResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "Input color message",
			}},
		{
			Label:         "InputBox options",
			ComponentType: ComponentTypeInputBox,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id: id + "_inputbox_options",
					Data: ut.IM{
						"input_type": InputBoxTypeSelect,
						"value":      "default value",
						"title":      "New fieldname",
						"message":    "Select the value:",
						"info":       "",
						"value_options": []SelectOption{
							{Text: "Option 1", Value: "option1"},
							{Text: "Option 2", Value: "option2"},
							{Text: "Option 3", Value: "option3"},
						},
						"tag":        "next_func",
						"default_ok": false,
					},
					EventURL:     eventURL,
					OnResponse:   testInputBoxResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "Select value",
			}},
		{
			Label:         "InputBox number value",
			ComponentType: ComponentTypeInputBox,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id: id + "_inputbox_number",
					Data: ut.IM{
						"input_type": InputBoxTypeNumber,
						"value":      "123.45",
						"title":      "New fieldname",
						"message":    "Enter the value:",
						"info":       "",
						"tag":        "next_func",
						"default_ok": false,
					},
					EventURL:     eventURL,
					OnResponse:   testInputBoxResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "Input number message",
			}},
		{
			Label:         "InputBox integer value",
			ComponentType: ComponentTypeInputBox,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id: id + "_inputbox_text",
					Data: ut.IM{
						"input_type": InputBoxTypeInteger,
						"value":      "123",
						"title":      "New fieldname",
						"message":    "Enter the value:",
						"info":       "",
						"tag":        "next_func",
						"default_ok": false,
					},
					EventURL:     eventURL,
					OnResponse:   testInputBoxResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "Input integer message",
			}},
		{
			Label:         "InputBox date value",
			ComponentType: ComponentTypeInputBox,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id: id + "_inputbox_date",
					Data: ut.IM{
						"input_type": InputBoxTypeDate,
						"value":      "2025-01-01",
						"title":      "New fieldname",
						"message":    "Enter the value:",
						"info":       "",
						"tag":        "next_func",
						"default_ok": false,
					},
					EventURL:     eventURL,
					OnResponse:   testInputBoxResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "Input date message",
			}},
		{
			Label:         "InputBox time value",
			ComponentType: ComponentTypeInputBox,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id: id + "_inputbox_time",
					Data: ut.IM{
						"input_type": InputBoxTypeTime,
						"value":      "12:00",
						"title":      "New fieldname",
						"message":    "Enter the value:",
						"info":       "",
						"tag":        "next_func",
						"default_ok": false,
					},
					EventURL:     eventURL,
					OnResponse:   testInputBoxResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "Input time message",
			}},
		{
			Label:         "InputBox datetime value",
			ComponentType: ComponentTypeInputBox,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id: id + "_inputbox_datetime",
					Data: ut.IM{
						"input_type": InputBoxTypeDateTime,
						"value":      "2025-01-01T12:00",
						"title":      "New fieldname",
						"message":    "Enter the value:",
						"info":       "",
						"tag":        "next_func",
						"default_ok": false,
					},
					EventURL:     eventURL,
					OnResponse:   testInputBoxResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				ButtonStyle: ButtonStyleDefault,
				Align:       TextAlignCenter,
				Label:       "Input datetime message",
			}},
	}
}
