package component

import (
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [Row] constants
const (
	ComponentTypeRow = "row"
)

// [Row] column
type RowColumn struct {
	// The label of the column
	Label string `json:"label"`
	// The component of the column
	Value Field `json:"value"`
}

// Responsive data row component
type Row struct {
	BaseComponent
	// The contents of the row's columns
	Columns []RowColumn `json:"columns"`
	// Full width input (100%)
	Full         bool `json:"full"`
	BorderTop    bool `json:"border_top"`
	BorderBottom bool `json:"border_bottom"`
	// For single column only
	FieldCol bool `json:"field_col"`
}

/*
Returns all properties of the [Row]
*/
func (row *Row) Properties() ut.IM {
	return ut.MergeIM(
		row.BaseComponent.Properties(),
		ut.IM{
			"columns":       row.Columns,
			"full":          row.Full,
			"border_top":    row.BorderTop,
			"border_bottom": row.BorderBottom,
			"field_col":     row.FieldCol,
		})
}

/*
Returns the value of the property of the [Row] with the specified name.
*/
func (row *Row) GetProperty(propName string) interface{} {
	return row.Properties()[propName]
}

/*
It checks the value given to the property of the [Row] and always returns a valid value
*/
func (row *Row) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"columns": func() interface{} {
			if value, valid := propValue.([]RowColumn); valid {
				return value
			}
			return []RowColumn{}
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if row.BaseComponent.GetProperty(propName) != nil {
		return row.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [Row] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (row *Row) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"columns": func() interface{} {
			row.Columns = row.Validation(propName, propValue).([]RowColumn)
			return row.Columns
		},
		"full": func() interface{} {
			row.Full = ut.ToBoolean(propValue, false)
			return row.Full
		},
		"border_top": func() interface{} {
			row.BorderTop = ut.ToBoolean(propValue, false)
			return row.BorderTop
		},
		"border_bottom": func() interface{} {
			row.BorderBottom = ut.ToBoolean(propValue, false)
			return row.BorderBottom
		},
		"field_col": func() interface{} {
			row.FieldCol = ut.ToBoolean(propValue, false)
			return row.FieldCol
		},
	}
	if _, found := pm[propName]; found {
		return row.SetRequestValue(propName, pm[propName](), []string{})
	}
	if row.BaseComponent.GetProperty(propName) != nil {
		return row.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (row *Row) getComponent(index int, coltype string) (res string, err error) {
	ccMap := map[string]func() ClientComponent{
		"field": func() ClientComponent {
			field := row.Columns[index].Value
			field.Id = row.Id + "_" + ut.ToString(index, "") + "_" + field.Name
			field.EventURL = row.EventURL
			field.OnResponse = row.OnResponse
			field.RequestValue = row.RequestValue
			field.RequestMap = row.RequestMap
			return &field
		},
		"label": func() ClientComponent {
			return &Label{
				Value: row.Columns[index].Label,
			}
		},
	}
	cc := ccMap[coltype]()
	res, err = cc.Render()
	return res, err
}

/*
Based on the values, it will generate the html code of the [Row] or return with an error message.
*/
func (row *Row) Render() (res string, err error) {
	row.InitProps(row)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(row.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(row.Class, " ")
		},
		"rowComponent": func(index int, coltype string) (string, error) {
			return row.getComponent(index, coltype)
		},
		"fieldCol": func() bool {
			return (len(row.Columns) == 1) && row.FieldCol
		},
		"validCol": func() bool {
			return (len(row.Columns) >= 1) && (len(row.Columns) <= 4) && !row.FieldCol
		},
		"colClass": func() string {
			cols := []string{"m12 l12", "m6 l6", "m4 l4", "m3 l3"}
			return cols[len(row.Columns)-1]
		},
	}
	tpl := `<div id="{{ .Id }}" name="{{ .Name }}" 
	class="row section-tiny {{ customClass }}{{ if .Full }} full{{ end }}
	{{ if .BorderTop }} border-top{{ end }}{{ if .BorderBottom }} border-bottom{{ end }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	>{{ if fieldCol }}<div class="cell padding-small hide-small" style="width: 150px;" >{{ rowComponent 0 "label" }}</div>
	<div class="cell padding-small" ><div class="section-tiny-bottom hide-medium hide-large" >{{ rowComponent 0 "label" }}</div>
	{{ rowComponent 0 "field" }}</div>{{ end }}
	{{ if validCol }}{{ range $index, $col := .Columns }}
	<div class="cell padding-small s12 {{ colClass }}" >
	<div class="section-tiny-bottom">{{ rowComponent $index "label" }}</div>{{ rowComponent $index "field" }}</div>
	{{ end }}{{ end }}
	</div>`

	return ut.TemplateBuilder("row", tpl, funcMap, row)
}

// [Row] test and demo data
func TestRow(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default (field col)",
			ComponentType: ComponentTypeRow,
			Component: &Row{
				BaseComponent: BaseComponent{
					Id:           id + "_default",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Columns: []RowColumn{
					{Label: "Comment", Value: Field{
						Type: FieldTypeString,
						Value: ut.IM{
							"name":        "string",
							"placeholder": "placeholder text",
						},
					}},
				},
				Full:         true,
				BorderTop:    true,
				BorderBottom: true,
				FieldCol:     true,
			}},
		{
			Label:         "1 col",
			ComponentType: ComponentTypeRow,
			Component: &Row{
				BaseComponent: BaseComponent{
					Id:           id + "_1_col",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
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
			}},
		{
			Label:         "2 col",
			ComponentType: ComponentTypeRow,
			Component: &Row{
				BaseComponent: BaseComponent{
					Id:           id + "_2_col",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
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
			}},
		{
			Label:         "3 col",
			ComponentType: ComponentTypeRow,
			Component: &Row{
				BaseComponent: BaseComponent{
					Id:           id + "_3_col",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
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
			}},
		{
			Label:         "4 col",
			ComponentType: ComponentTypeRow,
			Component: &Row{
				BaseComponent: BaseComponent{
					Id:           id + "_4_col",
					EventURL:     eventURL,
					OnResponse:   testFieldResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Columns: []RowColumn{
					{Label: "Number", Value: Field{
						Type: FieldTypeNumber,
						Value: ut.IM{
							"name":  "number",
							"value": 1.5,
						},
					}},
					{Label: "Date", Value: Field{
						Type: FieldTypeDate,
						Value: ut.IM{
							"name":    "date",
							"is_null": false,
						},
					}},
					{Label: "Time", Value: Field{
						Type: FieldTypeTime,
						Value: ut.IM{
							"name":    "time",
							"is_null": true,
						},
					}},
					{Label: "Boolean", Value: Field{
						Type: FieldTypeBool,
						Value: ut.IM{
							"name":  "bool",
							"value": true,
						},
					}},
				},
				Full:         true,
				BorderTop:    true,
				BorderBottom: true,
				FieldCol:     false,
			}},
	}
}
