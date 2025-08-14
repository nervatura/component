package component

import (
	"html/template"
	"slices"
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [Form] constants
const (
	ComponentTypeForm = "form"

	FormEventOK     = "form_ok"
	FormEventCancel = "form_cancel"
	FormEventChange = "form_change"
)

// Simple input form component
type Form struct {
	BaseComponent
	// The caption of the editor
	Title string `json:"title"`
	// Valid [Icon] component value. See more [IconValues] variable values.
	Icon string `json:"icon"`
	// The contents of the view rows
	BodyRows []Row `json:"body_rows"`
	// The footer of the form
	FooterRows []Row `json:"footer_rows"`
	// The modal mode
	Modal bool `json:"modal"`
}

/*
Returns all properties of the [Form]
*/
func (frm *Form) Properties() ut.IM {
	return ut.MergeIM(
		frm.BaseComponent.Properties(),
		ut.IM{
			"title":       frm.Title,
			"icon":        frm.Icon,
			"body_rows":   frm.BodyRows,
			"footer_rows": frm.FooterRows,
			"modal":       frm.Modal,
		})
}

/*
Returns the value of the property of the [Form] with the specified name.
*/
func (frm *Form) GetProperty(propName string) interface{} {
	return frm.Properties()[propName]
}

/*
It checks the value given to the property of the [Form] and always returns a valid value
*/
func (frm *Form) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"icon": func() interface{} {
			return frm.CheckEnumValue(ut.ToString(propValue, ""), IconFileText, IconValues)
		},
		"body_rows": func() interface{} {
			if value, valid := propValue.([]Row); valid {
				return value
			}
			return []Row{}
		},
		"footer_rows": func() interface{} {
			if value, valid := propValue.([]Row); valid {
				return value
			}
			return []Row{}
		},
		"target": func() interface{} {
			frm.SetProperty("id", frm.Id)
			value := ut.ToString(propValue, frm.Id)
			if value != "this" && !strings.HasPrefix(value, "#") {
				value = "#" + value
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if frm.BaseComponent.GetProperty(propName) != nil {
		return frm.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [Form] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (frm *Form) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"title": func() interface{} {
			frm.Title = ut.ToString(propValue, "")
			return frm.Title
		},
		"icon": func() interface{} {
			frm.Icon = frm.Validation(propName, propValue).(string)
			return frm.Icon
		},
		"body_rows": func() interface{} {
			frm.BodyRows = frm.Validation(propName, propValue).([]Row)
			return frm.BodyRows
		},
		"footer_rows": func() interface{} {
			frm.FooterRows = frm.Validation(propName, propValue).([]Row)
			return frm.FooterRows
		},
		"modal": func() interface{} {
			frm.Modal = frm.Validation(propName, propValue).(bool)
			return frm.Modal
		},
		"target": func() interface{} {
			frm.Target = frm.Validation(propName, propValue).(string)
			return frm.Target
		},
	}
	if _, found := pm[propName]; found {
		return frm.SetRequestValue(propName, pm[propName](), []string{})
	}
	if frm.BaseComponent.GetProperty(propName) != nil {
		return frm.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

/*
If the OnResponse function of the [Form] is implemented, the function calls it after the [TriggerEvent]
is processed, otherwise the function's return [ResponseEvent] is the processed [TriggerEvent].
*/
func (frm *Form) OnRequest(te TriggerEvent) (re ResponseEvent) {
	evt := ResponseEvent{
		Trigger: frm, TriggerName: frm.Name,
		Name:  FormEventCancel,
		Value: ut.IM{"value": ut.IM{}, "data": frm.Data},
	}
	if te.Values.Has(FormEventOK) {
		evt.Name = FormEventOK
	}
	values := ut.IM{}
	for key, value := range te.Values {
		if len(value) > 0 && !slices.Contains([]string{FormEventOK, FormEventCancel, FormEventChange}, key) {
			values[key] = value[0]
		}
	}
	evt.Value = ut.IM{"value": values, "data": frm.Data}

	if frm.OnResponse != nil {
		return frm.OnResponse(evt)
	}
	return evt
}

func (frm *Form) triggerEvent(evt ResponseEvent) (re ResponseEvent) {
	frmEvt := ResponseEvent{
		Trigger:     evt.Trigger,
		TriggerName: frm.Name,
		Name:        FormEventChange,
		Value: ut.IM{
			"name": evt.TriggerName, "event": evt.Name, "value": evt.Value,
			"form": frm, "data": frm.Data,
		},
	}
	if evt.TriggerName == "btn_close" {
		evt.Trigger = frm
		frmEvt.Name = FormEventCancel
		frmEvt.Value = ut.IM{
			"name": evt.TriggerName, "event": FormEventCancel, "value": "btn_close",
			"form": frm, "data": frm.Data,
		}
	}
	if frm.OnResponse != nil {
		return frm.OnResponse(frmEvt)
	}
	return frmEvt
}

func (frm *Form) getComponent(name string, index int) (html template.HTML, err error) {
	checkFieldTrigger := func(row *Row) {
		for index, column := range row.Columns {
			if column.Value.FormTrigger {
				row.Columns[index].Value.Id = frm.Id + "_" + column.Value.Name + "_" + ut.ToString(index, "")
				row.Columns[index].Value.EventURL = frm.EventURL
				row.Columns[index].Value.Target = frm.Target
				row.Columns[index].Value.OnResponse = frm.triggerEvent
				row.Columns[index].Value.RequestValue = frm.RequestValue
				row.Columns[index].Value.RequestMap = frm.RequestMap
			}
		}
	}
	ccMap := map[string]func() ClientComponent{
		"title": func() ClientComponent {
			return &Label{
				Value:    frm.Title,
				LeftIcon: frm.Icon,
			}
		},
		"body_row": func() ClientComponent {
			row := &frm.BodyRows[index]
			checkFieldTrigger(row)
			return row
		},
		"footer_row": func() ClientComponent {
			row := &frm.FooterRows[index]
			checkFieldTrigger(row)
			return row
		},
		"btn_close": func() ClientComponent {
			return &Icon{
				BaseComponent: BaseComponent{
					Id:           frm.Id + "_" + name,
					Name:         name,
					EventURL:     frm.EventURL,
					Target:       frm.Target,
					OnResponse:   frm.triggerEvent,
					RequestValue: frm.RequestValue,
					RequestMap:   frm.RequestMap,
					Class:        []string{"close-icon"},
				},
				Value: "Times",
			}
		},
	}
	cc := ccMap[name]()
	html, err = cc.Render()
	return html, err
}

/*
Based on the values, it will generate the html code of the [Form] or return with an error message.
*/
func (frm *Form) Render() (html template.HTML, err error) {
	frm.InitProps(frm)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(frm.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(frm.Class, " ")
		},
		"inputComponent": func(name string, index int) (template.HTML, error) {
			return frm.getComponent(name, index)
		},
		"footerRows": func() bool {
			return len(frm.FooterRows) > 0
		},
	}
	tpl := `<div id="{{ .Id }}" name="{{ .Name }}" class="row full {{ customClass }}">
	<form id="{{ .Id }}" name="inputbox_form"
	{{ if ne .EventURL "" }} hx-post="{{ .EventURL }}" hx-target="{{ .Target }}" {{ if ne .Sync "none" }} hx-sync="{{ .Sync }}"{{ end }} hx-swap="{{ .Swap }}"{{ end }}
	{{ if ne .Indicator "none" }} hx-indicator="#{{ .Indicator }}"{{ end }}
	>{{ if .Modal }}<div class="modal"><div class="dialog" 
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}>{{ end }}
	<div class="editor" {{ if and (eq .Modal false) (styleMap) }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}>
	<div class="editor-title">
	<div class="cell">{{ inputComponent "title" 0 }}</div>
	{{ if .Modal }}<div class="cell align-right">{{ inputComponent "btn_close" 0 }}</div>{{ end }}</div>
	<div class="section-small container-small" >
	{{ range $row_index, $row := $.BodyRows }}{{ inputComponent "body_row" $row_index }}{{ end }}
	</div>
	{{ if footerRows }}<div class="section-small container-small buttons full" >
	{{ range $row_index, $row := $.FooterRows }}{{ inputComponent "footer_row" $row_index }}{{ end }}
	</div>{{ end }}
	</div>{{ if .Modal }}</div></div>{{ end }}
	</form></div>`

	if html, err = ut.TemplateBuilder("inputform", tpl, funcMap, frm); err == nil && frm.EventURL != "" {
		frm.SetProperty("request_map", frm)
	}
	return html, err
}

var testFormResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
	if evt.Name == FormEventChange {
		values := ut.ToIM(evt.Value, ut.IM{})
		if !slices.Contains([]string{ListEventEditItem, ListEventDelete}, ut.ToString(values["event"], "")) {
			return evt
		}
	}
	return ResponseEvent{
		Trigger: &Toast{
			Type:    ToastTypeInfo,
			Value:   evt.Name,
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

func testFormBodyRows(key string) (bodyRows []Row) {
	switch key {
	case "info":
		return []Row{
			{
				Columns: []RowColumn{
					{Label: "Info message label", Value: Field{
						Type: FieldTypeLabel,
						Value: ut.IM{
							"value": "Info message text",
							"style": ut.SM{
								"font-weight": "normal",
								"font-style":  "italic",
							},
						},
					}},
				},
				Full:         false,
				FieldCol:     false,
				BorderTop:    false,
				BorderBottom: false,
			},
		}
	case "warning":
		return []Row{
			{
				Columns: []RowColumn{
					{Label: "The data has changed, but has not been saved!",
						Value: Field{
							Type: FieldTypeLabel,
							Value: ut.IM{
								"value": "Do you want to save changes?",
								"style": ut.SM{
									"font-weight": "normal",
									"font-style":  "italic",
								},
							},
						}},
				},
				Full:         false,
				FieldCol:     false,
				BorderTop:    false,
				BorderBottom: false,
			},
		}
	case "string_input":
		return []Row{
			{
				Columns: []RowColumn{
					{Label: "Please enter a name",
						Value: Field{
							Type: FieldTypeString,
							Value: ut.IM{
								"name":       "string",
								"value":      "default value",
								"auto_focus": true,
							},
						}},
				},
				Full:         true,
				FieldCol:     false,
				BorderTop:    false,
				BorderBottom: false,
			},
		}
	case "multiple_input":
		return []Row{
			{
				Columns: []RowColumn{
					{Label: "Required field",
						Value: Field{
							Type: FieldTypeString,
							Value: ut.IM{
								"name":        "string",
								"placeholder": "Required field",
								"invalid":     true,
								"auto_focus":  true,
								"required":    true,
							},
						}},
					{Label: "Select field",
						Value: Field{
							Type: FieldTypeSelect,
							Value: ut.IM{
								"name": "select",
								"options": []SelectOption{
									{Text: "Option 1", Value: "option1"},
									{Text: "Option 2", Value: "option2"},
									{Text: "Option 3", Value: "option3"},
								},
								"value":   "option1",
								"is_null": true,
							},
						}},
					{Label: "Date and time field",
						Value: Field{
							Type: FieldTypeDateTime,
							Value: ut.IM{
								"name":    "datetime",
								"value":   "2025-01-01 15:00",
								"is_null": false,
							},
						}},
				},
				Full:         true,
				FieldCol:     false,
				BorderTop:    false,
				BorderBottom: false,
			},
			{
				Columns: []RowColumn{
					{Label: "Integer (0-100)",
						Value: Field{
							Type: FieldTypeInteger,
							Value: ut.IM{
								"name":      "integer",
								"value":     "80",
								"max_value": 100,
								"min_value": 0,
								"set_max":   true,
								"set_min":   true,
							},
						}},
					{Label: "Time field",
						Value: Field{
							Type: FieldTypeTime,
							Value: ut.IM{
								"name":    "time",
								"value":   "15:00",
								"is_null": false,
							},
						}},
					{Label: "Boolean",
						Value: Field{
							Type: FieldTypeBool,
							Value: ut.IM{
								"name":  "boolean",
								"value": true,
							},
							FormTrigger: true,
						}},
					{Label: "Color input",
						Value: Field{
							Type: FieldTypeColor,
							Value: ut.IM{
								"name":  "color",
								"value": "#845185",
							},
						}},
				},
				Full:         true,
				FieldCol:     false,
				BorderTop:    true,
				BorderBottom: false,
			},
			{
				Columns: []RowColumn{
					{Label: "Comment field",
						Value: Field{
							Type: FieldTypeText,
							Value: ut.IM{
								"name":        "comment",
								"placeholder": "Enter a comment",
								"rows":        3,
							},
						}},
				},
				Full:         true,
				FieldCol:     false,
				BorderTop:    true,
				BorderBottom: false,
			},
		}
	}
	// "list_selector"
	return []Row{
		{
			Columns: []RowColumn{
				{Value: Field{
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
						"pagination":          PaginationTypeBottom,
						"page_size":           5,
						"hide_paginaton_size": true,
						"edit_item":           true,
						"list_filter":         true,
						"delete_item":         true,
					},
					FormTrigger: true,
				}},
			},
			Full:         true,
			FieldCol:     false,
			BorderTop:    false,
			BorderBottom: false,
		},
	}
}

func testFormFooterRows(key string) (footerRows []Row) {
	switch key {
	case "info":
		return []Row{
			{
				Columns: []RowColumn{
					{Value: Field{
						Type:  FieldTypeLabel,
						Value: ut.IM{},
					}},
					{Value: Field{
						Type: FieldTypeButton,
						Value: ut.IM{
							"name":         FormEventOK,
							"type":         ButtonTypeSubmit,
							"button_style": ButtonStylePrimary,
							"icon":         IconCheck,
							"label":        "OK",
							"auto_focus":   true,
							"selected":     true,
						},
					}},
					{Value: Field{
						Type:  FieldTypeLabel,
						Value: ut.IM{},
					}},
				},
				Full:         true,
				FieldCol:     false,
				BorderTop:    false,
				BorderBottom: false,
			},
		}
	case "warning":
		return []Row{
			{
				Columns: []RowColumn{
					{Value: Field{
						Type: FieldTypeButton,
						Value: ut.IM{
							"name":         FormEventOK,
							"type":         ButtonTypeSubmit,
							"button_style": ButtonStylePrimary,
							"icon":         IconCheck,
							"label":        "OK",
							"auto_focus":   true,
							"selected":     true,
						},
					}},
					{Value: Field{
						Type: FieldTypeButton,
						Value: ut.IM{
							"name":         FormEventCancel,
							"type":         ButtonTypeSubmit,
							"button_style": ButtonStyleDefault,
							"icon":         IconTimes,
							"label":        "Cancel",
						},
					}},
				},
				Full:         true,
				FieldCol:     false,
				BorderTop:    false,
				BorderBottom: false,
			},
		}
	case "string_input":
		return []Row{
			{
				Columns: []RowColumn{
					{Value: Field{
						Type: FieldTypeButton,
						Value: ut.IM{
							"name":         FormEventOK,
							"type":         ButtonTypeSubmit,
							"button_style": ButtonStylePrimary,
							"icon":         IconCheck,
							"label":        "OK",
							"auto_focus":   false,
							"selected":     false,
						},
					}},
					{Value: Field{
						Type: FieldTypeButton,
						Value: ut.IM{
							"name":         FormEventCancel,
							"type":         ButtonTypeSubmit,
							"button_style": ButtonStyleDefault,
							"icon":         IconTimes,
							"label":        "Cancel",
						},
					}},
				},
				Full:         true,
				FieldCol:     false,
				BorderTop:    false,
				BorderBottom: false,
			},
		}
	case "multiple_input":
		return []Row{
			{
				Columns: []RowColumn{
					{Value: Field{
						Type: FieldTypeButton,
						Value: ut.IM{
							"name":         FormEventOK,
							"type":         ButtonTypeSubmit,
							"button_style": ButtonStylePrimary,
							"icon":         IconCheck,
							"label":        "OK",
							"auto_focus":   false,
							"selected":     false,
						},
					}},
					{Value: Field{
						Type: FieldTypeButton,
						Value: ut.IM{
							"name":         FormEventCancel,
							"type":         ButtonTypeSubmit,
							"button_style": ButtonStyleDefault,
							"icon":         IconTimes,
							"label":        "Cancel",
						},
					}},
				},
				Full:         true,
				FieldCol:     false,
				BorderTop:    false,
				BorderBottom: false,
			},
		}
	}
	return []Row{}
}

// [Form] test and demo data
func TestForm(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeList,
			Component: &Form{
				BaseComponent: BaseComponent{
					Id:           id + "default_inputform",
					EventURL:     eventURL,
					OnResponse:   testFormResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Title:      "Info message",
				BodyRows:   testFormBodyRows("info"),
				FooterRows: testFormFooterRows("info"),
				Modal:      false,
				Icon:       IconInfoCircle,
			},
		},
		{
			Label:         "Warning",
			ComponentType: ComponentTypeList,
			Component: &Form{
				BaseComponent: BaseComponent{
					Id:           id + "warning_inputform",
					EventURL:     eventURL,
					OnResponse:   testFormResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Data: ut.IM{
						"tag": "next_func",
					},
				},
				Title:      "Warning message",
				BodyRows:   testFormBodyRows("warning"),
				FooterRows: testFormFooterRows("warning"),
				Modal:      false,
				Icon:       IconExclamationTriangle,
			},
		},
		{
			Label:         "String input",
			ComponentType: ComponentTypeList,
			Component: &Form{
				BaseComponent: BaseComponent{
					Id:           id + "string_inputform",
					EventURL:     eventURL,
					OnResponse:   testFormResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Data: ut.IM{
						"tag": "next_func",
					},
				},
				Title:      "String input",
				BodyRows:   testFormBodyRows("string_input"),
				FooterRows: testFormFooterRows("string_input"),
				Modal:      false,
				Icon:       IconQuestionCircle,
			},
		},
		{
			Label:         "Multiple fields input",
			ComponentType: ComponentTypeList,
			Component: &Form{
				BaseComponent: BaseComponent{
					Id:           id + "multiple_inputform",
					EventURL:     eventURL,
					OnResponse:   testFormResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Data: ut.IM{
						"tag": "next_func",
					},
					Style: ut.SM{
						"min-width": "600px",
					},
				},
				Title:      "Multiple fields input",
				BodyRows:   testFormBodyRows("multiple_input"),
				FooterRows: testFormFooterRows("multiple_input"),
				Modal:      false,
				Icon:       IconEdit,
			},
		},
		{
			Label:         "List selector",
			ComponentType: ComponentTypeList,
			Component: &Form{
				BaseComponent: BaseComponent{
					Id:           id + "selector_inputform",
					EventURL:     eventURL,
					OnResponse:   testFormResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Title:      "List selector",
				BodyRows:   testFormBodyRows("list_selector"),
				FooterRows: testFormFooterRows("list_selector"),
				Modal:      false,
				Icon:       IconSearch,
			},
		},
	}
}
