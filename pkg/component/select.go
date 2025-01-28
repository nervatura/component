package component

import (
	"html/template"
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [Select] constants
const (
	ComponentTypeSelect = "select"

	SelectEventChange = "change"
)

// [Select] control item
type SelectOption struct {
	Value string `json:"value"`
	Text  string `json:"text"`
}

/*
Creates an HTML select control

For example:

	&Select{
	  BaseComponent: BaseComponent{
	    Id:           "id_select_default",
	    EventURL:     "/event",
	    RequestValue: parent_component.GetProperty("request_value").(map[string]ut.IM),
	    RequestMap:   parent_component.GetProperty("request_map").(map[string]ClientComponent),
	  },
	  Value: "value1",
	  Options: []SelectOption{
	    {Value: "value1", Text: "Text 1"},
	    {Value: "value2", Text: "Text 2"},
	    {Value: "value3", Text: "Text 3"},
	  },
	  IsNull: true,
	}
*/
type Select struct {
	BaseComponent
	// Value of the selected item in options
	Value string `json:"value"`
	// Items of optional values
	Options []SelectOption `json:"options"`
	// An empty value can also be selected
	IsNull bool `json:"is_null"`
	// The HTML aria-label attribute of the component
	Label string `json:"label"`
	// Specifies that the input should be disabled
	Disabled bool `json:"disabled"`
	// Specifies that the input element should automatically get focus when the page loads
	AutoFocus bool `json:"auto_focus"`
	// Full width input (100%)
	Full bool `json:"full"`
}

/*
Returns all properties of the [Select]
*/
func (sel *Select) Properties() ut.IM {
	return ut.MergeIM(
		sel.BaseComponent.Properties(),
		ut.IM{
			"value":      sel.Value,
			"options":    sel.Options,
			"is_null":    sel.IsNull,
			"label":      sel.Label,
			"disabled":   sel.Disabled,
			"auto_focus": sel.AutoFocus,
			"full":       sel.Full,
		})
}

/*
Returns the value of the property of the [Select] with the specified name.
*/
func (sel *Select) GetProperty(propName string) interface{} {
	return sel.Properties()[propName]
}

func SelectOptionRangeValidation(values interface{}, defaultValue []SelectOption) []SelectOption {
	if values, valid := values.([]SelectOption); valid && len(values) >= 0 {
		return values
	}
	if values, valid := values.([]interface{}); valid {
		for _, value := range values {
			if valueMap, valid := value.(ut.IM); valid {
				defaultValue = append(defaultValue, SelectOption{Value: ut.ToString(valueMap["value"], ""), Text: ut.ToString(valueMap["text"], "")})
			}
		}
	}
	return defaultValue
}

/*
It checks the value given to the property of the [Select] and always returns a valid value
*/
func (sel *Select) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"options": func() interface{} {
			return SelectOptionRangeValidation(propValue, []SelectOption{})
		},
		"value": func() interface{} {
			value := ut.ToString(propValue, "")
			if (value == "") && sel.IsNull {
				return value
			}
			valid := false
			for _, opt := range sel.Options {
				if opt.Value == value {
					valid = true
				}
			}
			if !valid {
				if len(sel.Options) > 0 {
					value = sel.Options[0].Value
				} else {
					value = ""
				}
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if sel.BaseComponent.GetProperty(propName) != nil {
		return sel.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [Select] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (sel *Select) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			sel.Value = sel.Validation(propName, propValue).(string)
			return sel.Value
		},
		"options": func() interface{} {
			sel.Options = sel.Validation(propName, propValue).([]SelectOption)
			sel.SetProperty("value", sel.Value)
			return sel.Options
		},
		"is_null": func() interface{} {
			sel.IsNull = ut.ToBoolean(propValue, false)
			return sel.IsNull
		},
		"label": func() interface{} {
			sel.Label = ut.ToString(propValue, "")
			return sel.Label
		},
		"disabled": func() interface{} {
			sel.Disabled = ut.ToBoolean(propValue, false)
			return sel.Disabled
		},
		"auto_focus": func() interface{} {
			sel.AutoFocus = ut.ToBoolean(propValue, false)
			return sel.AutoFocus
		},
		"full": func() interface{} {
			sel.Full = ut.ToBoolean(propValue, false)
			return sel.Full
		},
	}
	if _, found := pm[propName]; found {
		return sel.SetRequestValue(propName, pm[propName](), []string{})
	}
	if sel.BaseComponent.GetProperty(propName) != nil {
		return sel.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

/*
If the OnResponse function of the [Select] is implemented, the function calls it after the [TriggerEvent]
is processed, otherwise the function's return [ResponseEvent] is the processed [TriggerEvent].
*/
func (sel *Select) OnRequest(te TriggerEvent) (re ResponseEvent) {
	value := sel.SetProperty("value", te.Values.Get(te.Name))
	evt := ResponseEvent{
		Trigger: sel, TriggerName: sel.Name,
		Name:  SelectEventChange,
		Value: value,
	}
	if sel.OnResponse != nil {
		return sel.OnResponse(evt)
	}
	return evt
}

/*
Based on the values, it will generate the html code of the [Select] or return with an error message.
*/
func (sel *Select) Render() (html template.HTML, err error) {
	sel.InitProps(sel)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(sel.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(sel.Class, " ")
		},
	}
	tpl := `<select id="{{ .Id }}" name="{{ .Name }}" value="{{ .Value }}"
	{{ if ne .EventURL "" }} hx-post="{{ .EventURL }}" hx-target="{{ .Target }}" {{ if ne .Sync "none" }} hx-sync="{{ .Sync }}"{{ end }} hx-swap="{{ .Swap }}"{{ end }}
	{{ if ne .Indicator "none" }} hx-indicator="#{{ .Indicator }}"{{ end }}
	{{ if .Disabled }} disabled{{ end }}
	{{ if .AutoFocus }} autofocus{{ end }}
	{{ if ne .Label "" }} aria-label="{{ .Label }}"{{ end }}
	 class="{{ customClass }}{{ if .Full }} full{{ end }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	>{{ if .IsNull }}<option {{ if eq .Value "" }}selected{{ end }} key="-1" value="" ></option>{{ end }}
	{{ range $index, $option := .Options }}<option {{ if eq .Value $.Value }}selected{{ end }} key="{{ $index }}" value="{{ $option.Value }}" >{{ $option.Text }}</option>{{ end }}
	</select>`

	if html, err = ut.TemplateBuilder("button", tpl, funcMap, sel); err == nil && sel.EventURL != "" {
		sel.SetProperty("request_map", sel)
	}
	return html, nil
}

// [Select] test and demo data
func TestSelect(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeSelect,
			Component: &Select{
				BaseComponent: BaseComponent{
					Id:           id + "_select_default",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value: "value1",
				Options: []SelectOption{
					{Value: "value1", Text: "Text 1"},
					{Value: "value2", Text: "Text 2"},
					{Value: "value3", Text: "Text 3"},
				},
				IsNull: true,
			}},
		{
			Label:         "Not null",
			ComponentType: ComponentTypeSelect,
			Component: &Select{
				BaseComponent: BaseComponent{
					Id:           id + "_select_not_null",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value: "value2",
				Options: []SelectOption{
					{Value: "value1", Text: "Text 1"},
					{Value: "value2", Text: "Text 2"},
					{Value: "value3", Text: "Text 3"},
				},
				IsNull: false,
				Full:   true,
			}},
	}
}
