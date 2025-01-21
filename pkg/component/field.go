package component

import (
	"html/template"

	ut "github.com/nervatura/component/pkg/util"
)

// [Field] constants
const (
	ComponentTypeField = "field"

	FieldTypeButton   = "button"
	FieldTypeUrlLink  = "url"
	FieldTypeString   = InputTypeString
	FieldTypeText     = InputTypeText
	FieldTypeColor    = InputTypeColor
	FieldTypePassword = InputTypePassword
	FieldTypeInteger  = "integer"
	FieldTypeNumber   = "float"
	FieldTypeDate     = DateTimeTypeDate
	FieldTypeTime     = DateTimeTypeTime
	FieldTypeDateTime = DateTimeTypeDateTime
	FieldTypeBool     = "bool"
	FieldTypeSelect   = "select"
	FieldTypeLink     = "link"
	FieldTypeUpload   = "upload"
	FieldTypeSelector = "selector"
	FieldTypeList     = "list"
)

// [Field] Type values
var FieldType []string = []string{
	FieldTypeButton, FieldTypeUrlLink, FieldTypeString, FieldTypeText, FieldTypeColor, FieldTypePassword,
	FieldTypeInteger, FieldTypeNumber, FieldTypeDate, FieldTypeTime, FieldTypeDateTime,
	FieldTypeBool, FieldTypeSelect, FieldTypeLink, FieldTypeUpload, FieldTypeSelector,
	FieldTypeList}

// Multi-type input component
type Field struct {
	BaseComponent
	// [FieldType] variable constants. Default value: [FieldTypeString]
	Type string `json:"type"`
	// Any valid value based on control type
	Value ut.IM `json:"value"`
}

/*
Returns all properties of the [Field]
*/
func (fld *Field) Properties() ut.IM {
	return ut.MergeIM(
		fld.BaseComponent.Properties(),
		ut.IM{
			"type":  fld.Type,
			"value": fld.Value,
		})
}

/*
Returns the value of the property of the [Field] with the specified name.
*/
func (fld *Field) GetProperty(propName string) interface{} {
	return fld.Properties()[propName]
}

/*
It checks the value given to the property of the [Field] and always returns a valid value
*/
func (fld *Field) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"type": func() interface{} {
			return fld.CheckEnumValue(ut.ToString(propValue, ""), FieldTypeString, FieldType)
		},
		"value": func() interface{} {
			value := ut.ToIM(fld.Value, ut.IM{})
			if imap, valid := propValue.(ut.IM); valid {
				value = ut.MergeIM(value, imap)
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if fld.BaseComponent.GetProperty(propName) != nil {
		return fld.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [Field] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (fld *Field) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"type": func() interface{} {
			fld.Type = fld.Validation(propName, propValue).(string)
			return fld.Type
		},
		"value": func() interface{} {
			fld.Value = fld.Validation(propName, propValue).(ut.IM)
			return fld.Value
		},
	}
	if _, found := pm[propName]; found {
		return fld.SetRequestValue(propName, pm[propName](), []string{})
	}
	if fld.BaseComponent.GetProperty(propName) != nil {
		return fld.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (fld *Field) getComponent() (html template.HTML, err error) {
	ccInp := func() *Input {
		inp := &Input{
			BaseComponent: BaseComponent{
				Id:           fld.Id + "_" + fld.Type,
				EventURL:     fld.EventURL,
				OnResponse:   fld.OnResponse,
				RequestValue: fld.RequestValue,
				RequestMap:   fld.RequestMap,
			},
			Type: fld.Type,
			Full: true,
		}
		return inp
	}
	ccNum := func() *NumberInput {
		inp := &NumberInput{
			BaseComponent: BaseComponent{
				Id:           fld.Id + "_" + fld.Type,
				EventURL:     fld.EventURL,
				OnResponse:   fld.OnResponse,
				RequestValue: fld.RequestValue,
				RequestMap:   fld.RequestMap,
			},
			Full: true,
		}
		return inp
	}
	ccDti := func() *DateTime {
		dti := &DateTime{
			BaseComponent: BaseComponent{
				Id:           fld.Id + "_" + fld.Type,
				EventURL:     fld.EventURL,
				OnResponse:   fld.OnResponse,
				RequestValue: fld.RequestValue,
				RequestMap:   fld.RequestMap,
			},
			Type: fld.Type,
			Full: true,
		}
		return dti
	}
	setProperty := func(cc ClientComponent) {
		for propName, propValue := range fld.Value {
			cc.SetProperty(propName, propValue)
		}
	}
	ccMap := map[string]func() ClientComponent{
		FieldTypeString: func() ClientComponent {
			inp := ccInp()
			setProperty(inp)
			return inp
		},
		FieldTypeText: func() ClientComponent {
			inp := ccInp()
			setProperty(inp)
			return inp
		},
		FieldTypePassword: func() ClientComponent {
			inp := ccInp()
			setProperty(inp)
			return inp
		},
		FieldTypeUpload: func() ClientComponent {
			inp := &Upload{
				BaseComponent: BaseComponent{
					Id:           fld.Id + "_" + fld.Type,
					EventURL:     fld.EventURL,
					OnResponse:   fld.OnResponse,
					RequestValue: fld.RequestValue,
					RequestMap:   fld.RequestMap,
				},
				Full: true,
			}
			setProperty(inp)
			return inp
		},
		FieldTypeSelector: func() ClientComponent {
			inp := &Selector{
				BaseComponent: BaseComponent{
					Id:           fld.Id + "_" + fld.Type,
					EventURL:     fld.EventURL,
					OnResponse:   fld.OnResponse,
					RequestValue: fld.RequestValue,
					RequestMap:   fld.RequestMap,
				},
				Full: true,
			}
			setProperty(inp)
			return inp
		},
		FieldTypeList: func() ClientComponent {
			inp := &List{
				BaseComponent: BaseComponent{
					Id:           fld.Id + "_" + fld.Type,
					EventURL:     fld.EventURL,
					OnResponse:   fld.OnResponse,
					RequestValue: fld.RequestValue,
					RequestMap:   fld.RequestMap,
				},
			}
			setProperty(inp)
			return inp
		},
		FieldTypeColor: func() ClientComponent {
			inp := ccInp()
			setProperty(inp)
			return inp
		},
		FieldTypeButton: func() ClientComponent {
			btn := &Button{
				BaseComponent: BaseComponent{
					Id:           fld.Id + "_" + fld.Type,
					EventURL:     fld.EventURL,
					OnResponse:   fld.OnResponse,
					RequestValue: fld.RequestValue,
					RequestMap:   fld.RequestMap,
				},
				Full: true,
			}
			setProperty(btn)
			return btn
		},
		FieldTypeUrlLink: func() ClientComponent {
			btn := &Link{
				BaseComponent: BaseComponent{
					Id: fld.Id + "_" + fld.Type,
				},
				Full: true,
			}
			setProperty(btn)
			return btn
		},
		FieldTypeNumber: func() ClientComponent {
			inp := ccNum()
			setProperty(inp)
			if _, found := fld.Value["value"]; found {
				inp.SetProperty("", fld.Value["value"])
			}
			return inp
		},
		FieldTypeInteger: func() ClientComponent {
			inp := ccNum()
			inp.Integer = true
			setProperty(inp)
			if _, found := fld.Value["value"]; found {
				inp.SetProperty("", fld.Value["value"])
			}
			return inp
		},
		FieldTypeDate: func() ClientComponent {
			dti := ccDti()
			setProperty(dti)
			if _, found := fld.Value["value"]; found {
				dti.SetProperty("", fld.Value["value"])
			}
			return dti
		},
		FieldTypeTime: func() ClientComponent {
			dti := ccDti()
			setProperty(dti)
			if _, found := fld.Value["value"]; found {
				dti.SetProperty("", fld.Value["value"])
			}
			return dti
		},
		FieldTypeDateTime: func() ClientComponent {
			dti := ccDti()
			setProperty(dti)
			if _, found := fld.Value["value"]; found {
				dti.SetProperty("", fld.Value["value"])
			}
			return dti
		},
		FieldTypeSelect: func() ClientComponent {
			options := SelectOptionRangeValidation(fld.Value["options"], []SelectOption{})
			sel := &Select{
				BaseComponent: BaseComponent{
					Id:           fld.Id + "_" + fld.Type,
					EventURL:     fld.EventURL,
					OnResponse:   fld.OnResponse,
					RequestValue: fld.RequestValue,
					RequestMap:   fld.RequestMap,
				},
				Options: options,
				Full:    true,
			}
			setProperty(sel)
			if _, found := fld.Value["value"]; found {
				sel.SetProperty("", fld.Value["value"])
			}
			return sel
		},
		FieldTypeLink: func() ClientComponent {
			lbl := &Label{
				BaseComponent: BaseComponent{
					Id:           fld.Id + "_" + fld.Type,
					EventURL:     fld.EventURL,
					OnResponse:   fld.OnResponse,
					RequestValue: fld.RequestValue,
					RequestMap:   fld.RequestMap,
				},
				Border: true,
				Full:   true,
			}
			setProperty(lbl)
			return lbl
		},
		FieldTypeBool: func() ClientComponent {
			tgl := &Toggle{
				BaseComponent: BaseComponent{
					Id:           fld.Id + "_" + fld.Type,
					EventURL:     fld.EventURL,
					OnResponse:   fld.OnResponse,
					RequestValue: fld.RequestValue,
					RequestMap:   fld.RequestMap,
				},
				Border: true,
				Full:   true,
			}
			setProperty(tgl)
			return tgl
		},
	}
	cc := ccMap[fld.Type]()
	html, err = cc.Render()
	return html, err
}

/*
Based on the values, it will generate the html code of the [Field] or return with an error message.
*/
func (fld *Field) Render() (html template.HTML, err error) {
	fld.InitProps(fld)
	return fld.getComponent()
}

var testFieldResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
	toast := func(value string) ResponseEvent {
		return ResponseEvent{
			Trigger: &Toast{
				Type:    ToastTypeInfo,
				Value:   value,
				Timeout: 4,
			},
			TriggerName: evt.TriggerName,
			Name:        evt.Name,
			Header: ut.SM{
				HeaderRetarget: "#toast-msg",
				HeaderReswap:   SwapInnerHTML,
			},
		}
	}
	switch evt.TriggerName {
	case "button":
		badge := ut.ToInteger(evt.Trigger.GetProperty("badge"), 0)
		evt.Trigger.SetProperty("badge", badge+1)
		evt.Trigger.SetProperty("show_badge", true)
	case "link":
		return toast(ut.ToString(evt.Value, ""))
	case "list":
		row := evt.Value.(ut.IM)["row"].(ut.IM)
		return toast(ut.ToString(row["lsvalue"], ""))
	case "selector":
		switch evt.Name {
		case SelectorEventLink:
			return toast(evt.Value.(SelectOption).Text)
		case SelectorEventSearch:
			return toast(ut.ToString(evt.Trigger.GetProperty("filter_value"), ""))
		case SelectorEventSelected:
			if value, valid := evt.Value.(ut.IM); valid {
				if row, valid := value["row"].(ut.IM); valid {
					evt.Trigger.SetProperty("value", SelectOption{
						Value: ut.ToString(row["id"], ""),
						Text:  ut.ToString(row["custname"], ""),
					})
				}
			}
		}
	}
	return evt
}

// [Field] test and demo data
func TestField(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default (string)",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id:           id + "_string",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: FieldTypeString,
				Value: ut.IM{
					"name":        "string",
					"placeholder": "placeholder text",
				},
			}},
		{
			Label:         "Long text",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id:           id + "_text",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: FieldTypeText,
				Value: ut.IM{
					"name":  "text",
					"value": `Long text&#13;&#10;Next row...`,
				},
			}},
		{
			Label:         "Password, readonly",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id:           id + "_password",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: FieldTypePassword,
				Value: ut.IM{
					"name":     "password",
					"value":    "secret",
					"readonly": true,
				},
			}},
		{
			Label:         "Color input",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id:           id + "_password",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: FieldTypeColor,
				Value: ut.IM{
					"name":  "color",
					"value": "#845185",
				},
			}},

		{
			Label:         "Number",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id:           id + "_number",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: FieldTypeNumber,
				Value: ut.IM{
					"name":  "number",
					"value": 1.5,
				},
			}},
		{
			Label:         "Integer min(0), max(100)",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id:           id + "_integer",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: FieldTypeInteger,
				Value: ut.IM{
					"name":      "integer",
					"value":     25,
					"set_max":   true,
					"max_value": 100,
					"set_min":   true,
					"min_value": 0,
				},
			}},
		{
			Label:         "Button, primary and icon",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id:           id + "_button",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: FieldTypeButton,
				Value: ut.IM{
					"name":         "button",
					"button_style": ButtonStylePrimary,
					"label":        "Primary",
					"icon":         "Check",
				},
			}},
		{
			Label:         "URL link button, border and icon",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id: id + "_url_link",
				},
				Type: FieldTypeUrlLink,
				Value: ut.IM{
					"name":        "url_link",
					"link_style":  LinkStyleBorder,
					"label":       "Search",
					"icon":        "Search",
					"href":        "https://www.google.com",
					"link_target": "_blank",
				},
			}},
		{
			Label:         "Date (not null)",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id:           id + "_date",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: FieldTypeDate,
				Value: ut.IM{
					"name":    "date",
					"value":   "2024-12-24",
					"is_null": false,
				},
			}},
		{
			Label:         "Time",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id:           id + "_time",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: FieldTypeTime,
				Value: ut.IM{
					"name":    "time",
					"is_null": true,
					"value":   "12:24",
				},
			}},
		{
			Label:         "DateTime (not null)",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id:           id + "_datetime",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: FieldTypeDateTime,
				Value: ut.IM{
					"name":    "datetime",
					"is_null": false,
					"value":   "2024-12-24T12:24",
				},
			}},
		{
			Label:         "Select",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id:           id + "_select",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: FieldTypeSelect,
				Value: ut.IM{
					"name":    "select",
					"is_null": true,
					"options": []SelectOption{
						{Value: "value1", Text: "Text 1"},
						{Value: "value2", Text: "Text 2"},
						{Value: "value3", Text: "Text 3"},
					},
					"value": "value2",
				},
			}},
		{
			Label:         "Link value",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id:           id + "_link",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: FieldTypeLink,
				Value: ut.IM{
					"name":  "link",
					"value": "Product name",
				},
			}},
		{
			Label:         "Bool value",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id:           id + "_bool",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: FieldTypeBool,
				Value: ut.IM{
					"name":  "bool",
					"value": true,
				},
			}},
		{
			Label:         "Upload",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id:           id + "_upload",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: FieldTypeUpload,
				Value: ut.IM{
					"name":       "upload",
					"accept":     "image/*",
					"max_length": 50,
				},
			}},
		{
			Label:         "Selector",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id:           id + "_selector",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: FieldTypeSelector,
				Value: ut.IM{
					"name":    "selector",
					"value":   SelectOption{Value: "12345", Text: "Customer Name"},
					"fields":  testSelectorFields,
					"rows":    testSelectorRows,
					"link":    true,
					"is_null": true,
				},
			}},
		{
			Label:         "List",
			ComponentType: ComponentTypeField,
			Component: &Field{
				BaseComponent: BaseComponent{
					Id:           id + "_list",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type: FieldTypeList,
				Value: ut.IM{
					"name": "list",
					"rows": []ut.IM{
						{"lslabel": "Label 1", "lsvalue": "Value row 1"},
						{"lslabel": "Label 2", "lsvalue": "Value row 2"},
						{"lslabel": "", "lsvalue": "Value row 3"},
						{"lslabel": "", "lsvalue": "Value row 6"},
						{"lslabel": "", "lsvalue": "Value row 6"},
						{"lslabel": "", "lsvalue": "Value row 6"},
						{"lslabel": "Label 7", "lsvalue": "Value row 7"},
						{"lslabel": "Label 8", "lsvalue": "Value row 8"},
						{"lslabel": "Label 9", "lsvalue": "Value row 9"},
					},
					"pagination":  PaginationTypeNone,
					"delete_item": true,
					"edit_icon":   "Plus",
				},
			}},
	}
}
