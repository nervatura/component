package component

import (
	"html/template"
	"slices"
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [Editor] constants
const (
	ComponentTypeEditor = "editor"

	EditorEventView  = "edit_view"
	EditorEventField = "edit_field"
)

type EditorView struct {
	Key string `json:"key"`
	// The label of the view
	Label string `json:"label"`
	// Valid [Icon] component value. See more [IconKey] variable values.
	Icon string `json:"icon"`
	// The badge value of the view
	Badge string `json:"badge"`
}

// Generic input form component
type Editor struct {
	BaseComponent
	// The caption of the editor
	Title string `json:"title"`
	// Valid [Icon] component value. See more [IconKey] variable values.
	Icon string `json:"icon"`
	// Current view key
	View  string       `json:"view"`
	Views []EditorView `json:"views"`
	// The contents of the view rows
	Rows []Row `json:"rows"`
	// The contents of the view table
	Tables []Table `json:"tables"`
}

/*
Returns all properties of the [Editor]
*/
func (edi *Editor) Properties() ut.IM {
	return ut.MergeIM(
		edi.BaseComponent.Properties(),
		ut.IM{
			"title":  edi.Title,
			"icon":   edi.Icon,
			"view":   edi.View,
			"views":  edi.Views,
			"rows":   edi.Rows,
			"tables": edi.Tables,
		})
}

/*
Returns the value of the property of the [Editor] with the specified name.
*/
func (edi *Editor) GetProperty(propName string) interface{} {
	return edi.Properties()[propName]
}

/*
It checks the value given to the property of the [Editor] and always returns a valid value
*/
func (edi *Editor) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"views": func() interface{} {
			if value, valid := propValue.([]EditorView); valid {
				return value
			}
			return []EditorView{}
		},
		"rows": func() interface{} {
			if value, valid := propValue.([]Row); valid {
				return value
			}
			return []Row{}
		},
		"tables": func() interface{} {
			if value, valid := propValue.([]Table); valid {
				return value
			}
			return []Table{}
		},
		"target": func() interface{} {
			edi.SetProperty("id", edi.Id)
			value := ut.ToString(propValue, edi.Id)
			if value != "this" && !strings.HasPrefix(value, "#") {
				value = "#" + value
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if edi.BaseComponent.GetProperty(propName) != nil {
		return edi.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [Editor] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (edi *Editor) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"title": func() interface{} {
			edi.Title = ut.ToString(propValue, "")
			return edi.Title
		},
		"icon": func() interface{} {
			edi.Icon = ut.ToString(propValue, "FileText")
			return edi.Icon
		},
		"view": func() interface{} {
			edi.View = ut.ToString(propValue, "")
			return edi.View
		},
		"views": func() interface{} {
			edi.Views = edi.Validation(propName, propValue).([]EditorView)
			return edi.Views
		},
		"rows": func() interface{} {
			edi.Rows = edi.Validation(propName, propValue).([]Row)
			return edi.Rows
		},
		"tables": func() interface{} {
			edi.Tables = edi.Validation(propName, propValue).([]Table)
			return edi.Tables
		},
		"target": func() interface{} {
			edi.Target = edi.Validation(propName, propValue).(string)
			return edi.Target
		},
	}
	if _, found := pm[propName]; found {
		return edi.SetRequestValue(propName, pm[propName](), []string{})
	}
	if edi.BaseComponent.GetProperty(propName) != nil {
		return edi.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (edi *Editor) response(evt ResponseEvent) (re ResponseEvent) {
	admEvt := ResponseEvent{
		Trigger: edi, TriggerName: edi.Name, Value: evt.Value,
		Header: ut.SM{HeaderRetarget: "#" + edi.Id},
	}
	switch evt.TriggerName {
	case "tab_btn":
		data := evt.Trigger.GetProperty("data").(ut.IM)
		admEvt.Value = ut.ToString(data["key"], "")
		admEvt.Name = EditorEventView

	case "view_table":
		if slices.Contains([]string{TableEventAddItem, TableEventEditCell, TableEventRowSelected}, evt.Name) {
			admEvt.Name = EditorEventField
			admEvt.Value = ut.IM{"name": evt.Name, "value": evt.Value, "data": evt.Trigger.GetProperty("data").(ut.IM)}
		} else {
			admEvt = evt
		}

	default:
		admEvt.Name = EditorEventField
		admEvt.Value = ut.IM{"name": evt.TriggerName, "event": evt.Name, "value": evt.Value, "data": evt.Trigger.GetProperty("data").(ut.IM)}
	}
	if edi.OnResponse != nil {
		return edi.OnResponse(admEvt)
	}
	return admEvt
}

func (edi *Editor) getComponent(name string, view EditorView, index int) (html template.HTML, err error) {
	ccMap := map[string]func() ClientComponent{
		"title": func() ClientComponent {
			return &Label{
				Value:    edi.Title,
				LeftIcon: edi.Icon,
			}
		},
		"tab_btn": func() ClientComponent {
			return &Button{
				BaseComponent: BaseComponent{
					Id:           edi.Id + "_" + name + "_" + view.Key,
					Name:         name,
					Data:         ut.IM{"key": view.Key},
					EventURL:     edi.EventURL,
					Target:       edi.Id,
					OnResponse:   edi.response,
					RequestValue: edi.RequestValue,
					RequestMap:   edi.RequestMap,
					Style:        ut.SM{"border-radius": "0", "margin-top": "2px", "opacity": "1"},
				},
				ButtonStyle: ButtonStyleDefault,
				Label:       view.Label,
				Icon:        view.Icon,
				Badge:       view.Badge,
				Full:        true,
				Align:       TextAlignLeft,
				Selected:    (view.Key == edi.View),
				Disabled:    (view.Key == edi.View),
			}
		},
		"view_row": func() ClientComponent {
			row := &edi.Rows[index]
			row.Id = edi.Id + "_" + name + "_" + ut.ToString(index, "")
			row.Name = name
			row.EventURL = edi.EventURL
			row.Target = edi.Target
			row.OnResponse = edi.response
			row.RequestValue = edi.RequestValue
			row.RequestMap = edi.RequestMap
			return row
		},
		"view_table": func() ClientComponent {
			tbl := &edi.Tables[index]
			tbl.Id = edi.Id + "_" + name + "_" + ut.ToString(index, "")
			tbl.Name = name
			tbl.EventURL = edi.EventURL
			tbl.Target = edi.Target
			tbl.OnResponse = edi.response
			tbl.RequestValue = edi.RequestValue
			tbl.RequestMap = edi.RequestMap
			return tbl
		},
	}
	cc := ccMap[name]()
	html, err = cc.Render()
	return html, err
}

/*
Based on the values, it will generate the html code of the [Editor] or return with an error message.
*/
func (edi *Editor) Render() (html template.HTML, err error) {
	edi.InitProps(edi)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(edi.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(edi.Class, " ")
		},
		"editorComponent": func(name string) (template.HTML, error) {
			return edi.getComponent(name, EditorView{}, 0)
		},
		"viewComponent": func(name string, view EditorView, index int) (template.HTML, error) {
			return edi.getComponent(name, view, index)
		},
	}
	tpl := `<div id="{{ .Id }}"
	class="{{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}>
	<div class="page"><div class="editor">
	<div class="editor-title"><div class="cell">{{ editorComponent "title" }}</div></div>
	<div class="section-container" >
	{{ range $index, $view := .Views }}
	{{ viewComponent "tab_btn" $view 0 }}
	{{ if eq $view.Key $.View }}<div class="row-panel" >
	{{ range $row_index, $row := $.Rows }}{{ viewComponent "view_row" $view $row_index }}{{ end }}
	{{ range $tbl_index, $tbl := $.Tables }}{{ viewComponent "view_table" $view $tbl_index }}{{ end }}
	</div>{{ end }}
	{{ end }}
	</div></div>
	</div>
	</div>`

	return ut.TemplateBuilder("editor", tpl, funcMap, edi)
}

var testEditorResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
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

// [Editor] test and demo data
func TestEditor(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeList,
			Component: &Editor{
				BaseComponent: BaseComponent{
					Id:           id + "default_editor",
					EventURL:     eventURL,
					OnResponse:   testEditorResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Title: "Editor",
				View:  "main",
				Views: []EditorView{
					{
						Key:   "main",
						Label: "Main input",
						Icon:  "ShoppingCart",
					},
					{
						Key:   "item",
						Label: "Item rows",
						Icon:  "User",
						Badge: "3",
					},
				},
				Rows: []Row{
					{
						Columns: []RowColumn{
							{Label: "Select", Value: Field{
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
							{Label: "Selector", Value: Field{
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
						},
						Full:         true,
						BorderTop:    true,
						BorderBottom: true,
						FieldCol:     false,
					},
					{
						Columns: []RowColumn{
							{Label: "Button", Value: Field{
								Type: FieldTypeButton,
								Value: ut.IM{
									"name":         "button",
									"button_style": ButtonStylePrimary,
									"label":        "Primary",
									"icon":         "Check",
								},
							}},
							{Label: "DateTime", Value: Field{
								Type: FieldTypeDateTime,
								Value: ut.IM{
									"name":    "datetime",
									"is_null": false,
								},
							}},
							{Label: "Link", Value: Field{
								Type: FieldTypeLink,
								Value: ut.IM{
									"name":  "link",
									"value": "Product name",
								},
							}},
						},
						Full:         true,
						BorderTop:    true,
						BorderBottom: true,
						FieldCol:     false,
					},
					{
						Columns: []RowColumn{
							{Label: "Note", Value: Field{
								Type: FieldTypeText,
								Value: ut.IM{
									"name":  "text",
									"value": `Long text&#13;&#10;Next row...`,
								},
							}},
						},
						Full:         true,
						BorderTop:    true,
						BorderBottom: true,
						FieldCol:     false,
					},
				},
				Tables: []Table{},
			}},
		{
			Label:         "Item editor",
			ComponentType: ComponentTypeList,
			Component: &Editor{
				BaseComponent: BaseComponent{
					Id:           id + "item_editor",
					EventURL:     eventURL,
					OnResponse:   testEditorResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Title: "Editor",
				View:  "item",
				Views: []EditorView{
					{
						Key:   "main",
						Label: "Main input",
						Icon:  "ShoppingCart",
					},
					{
						Key:   "item",
						Label: "Item rows",
						Icon:  "User",
						Badge: "3",
					},
				},
				Rows: []Row{},
				Tables: []Table{
					{
						Fields: []TableField{
							{Name: "custnumber", Label: "Customer No."},
							{Name: "custname", Label: "Customer Name", FieldType: TableFieldTypeLink},
							{Name: "taxnumber", Label: "Taxnumber"},
							{Name: "address", Label: "Address"},
						},
						Rows:        testBrowserRows["customer"](),
						RowSelected: true,
					},
				},
			}},
	}
}
