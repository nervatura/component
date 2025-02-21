package component

import (
	"html/template"
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [Selector] constants
const (
	ComponentTypeSelector = "selector"

	SelectorEventFilterChange = "selector_filter_change"
	SelectorEventLink         = "selector_link"
	SelectorEventSearch       = "selector_search"
	SelectorEventDelete       = "selector_delete"
	SelectorEventShowModal    = "selector_modal"
	SelectorEventSelected     = "selector_selected"

	SelectorDefaultPlaceholder = "Search conditions"
	SelectorDefaultTitle       = "Search for data"
)

// Creates an interactive input control
type Selector struct {
	BaseComponent
	// Current page value
	Value SelectOption `json:"value"`
	// Data source of the table
	Rows []ut.IM `json:"rows"`
	// Table column definitions
	Fields []TableField `json:"fields"`
	// Search window title
	Title string `json:"title"`
	// Specifies a short hint that describes the expected value of the input element
	FilterPlaceholder string `json:"filter_placeholder"`
	// The displayed text is a link
	Link bool `json:"link"`
	// Show/hide clear button
	IsNull bool `json:"is_null"`
	// Specifies that the input should be disabled
	Disabled bool `json:"disabled"`
	// Specifies that the input element should automatically get focus when the page loads
	AutoFocus bool `json:"auto_focus"`
	// Full width input (100%)
	Full bool `json:"full"`
	// Displaying the search window
	ShowModal bool `json:"show_modal"`
	// It prevents the modal search window and only the simple SelectorEventShowModal event is triggered
	CustomModal bool `json:"custom_modal"`
	// Icon for the modal button. See more [IconValues] variable values. Default: Search
	ModalIcon string `json:"modal_icon"`
}

/*
Returns all properties of the [Selector]
*/
func (sel *Selector) Properties() ut.IM {
	return ut.MergeIM(
		sel.BaseComponent.Properties(),
		ut.IM{
			"value":              sel.Value,
			"rows":               sel.Rows,
			"fields":             sel.Fields,
			"title":              sel.Title,
			"link":               sel.Link,
			"filter_placeholder": sel.FilterPlaceholder,
			"is_null":            sel.IsNull,
			"disabled":           sel.Disabled,
			"auto_focus":         sel.AutoFocus,
			"full":               sel.Full,
			"show_modal":         sel.ShowModal,
			"custom_modal":       sel.CustomModal,
			"modal_icon":         sel.ModalIcon,
		})
}

/*
Returns the value of the property of the [Selector] with the specified name.
*/
func (sel *Selector) GetProperty(propName string) interface{} {
	return sel.Properties()[propName]
}

/*
It checks the value given to the property of the [Selector] and always returns a valid value
*/
func (sel *Selector) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			if value, valid := propValue.(SelectOption); valid {
				return value
			}
			if valueOptions, found := propValue.(ut.IM); found {
				return SelectOption{
					Value: ut.ToString(valueOptions["value"], ""),
					Text:  ut.ToString(valueOptions["text"], ""),
				}
			}
			return SelectOption{}
		},
		"rows": func() interface{} {
			return ut.ToIMA(propValue, []ut.IM{})
		},
		"fields": func() interface{} {
			fields := []TableField{}
			if fd, valid := propValue.([]TableField); valid && (fd != nil) {
				fields = fd
			}
			if len(fields) == 0 {
				if len(sel.Rows) > 0 {
					for field := range sel.Rows[0] {
						fields = append(fields,
							TableField{Name: field, FieldType: TableFieldTypeString, Label: field})
					}
				}
			}
			return fields
		},
		"modal_icon": func() interface{} {
			return sel.CheckEnumValue(ut.ToString(propValue, ""), IconSearch, IconValues)
		},
		"target": func() interface{} {
			sel.SetProperty("id", sel.Id)
			value := ut.ToString(propValue, sel.Id)
			if value != "this" && !strings.HasPrefix(value, "#") {
				value = "#" + value
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
Setting a property of the [Selector] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (sel *Selector) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			sel.Value = sel.Validation(propName, propValue).(SelectOption)
			return sel.Value
		},
		"rows": func() interface{} {
			sel.Rows = sel.Validation(propName, propValue).([]ut.IM)
			return sel.Rows
		},
		"fields": func() interface{} {
			sel.Fields = sel.Validation(propName, propValue).([]TableField)
			return sel.Fields
		},
		"title": func() interface{} {
			sel.Title = ut.ToString(propValue, SelectorDefaultTitle)
			return sel.Title
		},
		"filter_placeholder": func() interface{} {
			sel.FilterPlaceholder = ut.ToString(propValue, "")
			return sel.FilterPlaceholder
		},
		"link": func() interface{} {
			sel.Link = ut.ToBoolean(propValue, false)
			return sel.Link
		},
		"is_null": func() interface{} {
			sel.IsNull = ut.ToBoolean(propValue, false)
			return sel.IsNull
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
		"show_modal": func() interface{} {
			sel.ShowModal = ut.ToBoolean(propValue, false)
			return sel.ShowModal
		},
		"custom_modal": func() interface{} {
			sel.CustomModal = ut.ToBoolean(propValue, false)
			return sel.CustomModal
		},
		"modal_icon": func() interface{} {
			sel.ModalIcon = sel.Validation(propName, propValue).(string)
			return sel.ModalIcon
		},
		"target": func() interface{} {
			sel.Target = sel.Validation(propName, propValue).(string)
			return sel.Target
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

func (sel *Selector) response(evt ResponseEvent) (re ResponseEvent) {
	selEvt := ResponseEvent{Trigger: sel, TriggerName: sel.Name, Value: sel.Value}
	switch evt.TriggerName {
	case "selector_result":
		if evt.Name != TableEventRowSelected {
			return evt
		}
		selEvt.Name = SelectorEventSelected
		selEvt.Value = evt.Value
		selEvt.Header = ut.SM{
			HeaderRetarget: "#" + sel.Id,
		}
		sel.SetProperty("show_modal", false)

	case "btn_delete":
		value := sel.SetProperty("value", SelectOption{})
		selEvt.Name = SelectorEventDelete
		selEvt.Value = value

	case "btn_modal":
		selEvt.Name = SelectorEventShowModal
		sel.SetProperty("show_modal", !sel.CustomModal)

	case "btn_close":
		selEvt.Name = SelectorEventShowModal
		sel.SetProperty("show_modal", false)

	case "filter_value":
		selEvt.Name = SelectorEventSearch
		selEvt.Value = evt.Value
		sel.SetProperty("data", ut.IM{evt.TriggerName: evt.Value})

	case "selector_text":
		selEvt.Name = SelectorEventLink

	case "btn_search":
		selEvt.Value = ut.ToString(sel.Data["filter_value"], "")
		selEvt.Name = SelectorEventSearch

	default:
	}
	if sel.OnResponse != nil {
		return sel.OnResponse(selEvt)
	}
	return selEvt
}

func (sel *Selector) getComponent(name string) (html template.HTML, err error) {
	ccBtn := func(icon string, focus bool) *Button {
		return &Button{
			BaseComponent: BaseComponent{
				Id:           sel.Id + "_" + name,
				Name:         name,
				Style:        ut.SM{"padding": "8px", "margin": "1px 0 2px 1px"},
				EventURL:     sel.EventURL,
				Target:       sel.Target,
				OnResponse:   sel.response,
				RequestValue: sel.RequestValue,
				RequestMap:   sel.RequestMap,
			},
			ButtonStyle: ButtonStyleBorder,
			Icon:        icon,
			Disabled:    sel.Disabled,
			AutoFocus:   focus,
		}
	}
	ccInp := func(value string) *Input {
		inp := &Input{
			BaseComponent: BaseComponent{
				Id: sel.Id + "_" + name, Name: name,
				Style:        ut.SM{"border-radius": "0", "margin": "1px 0 2px"},
				EventURL:     sel.EventURL,
				Target:       sel.Target,
				OnResponse:   sel.response,
				RequestValue: sel.RequestValue,
				RequestMap:   sel.RequestMap,
			},
			Type:        InputTypeString,
			Label:       sel.FilterPlaceholder,
			Placeholder: sel.FilterPlaceholder,
			AutoFocus:   true,
			Full:        true,
		}
		inp.SetProperty("value", value)
		return inp
	}
	ccMap := map[string]func() ClientComponent{
		"btn_modal": func() ClientComponent {
			return ccBtn(sel.ModalIcon, sel.AutoFocus)
		},
		"btn_delete": func() ClientComponent {
			return ccBtn(IconTimes, false)
		},
		"btn_search": func() ClientComponent {
			return ccBtn(IconSearch, false)
		},
		"selector_text": func() ClientComponent {
			lbl := &Label{
				BaseComponent: BaseComponent{
					Id: sel.Id + "_" + name,
				},
				Border: true,
				Full:   true,
			}
			if sel.Link && !sel.Disabled && (sel.Value.Text != "") {
				lbl.BaseComponent = BaseComponent{
					Id:           sel.Id + "_" + name,
					Name:         name,
					EventURL:     sel.EventURL,
					Target:       sel.Target,
					OnResponse:   sel.response,
					RequestValue: sel.RequestValue,
					RequestMap:   sel.RequestMap,
				}
			}
			lbl.SetProperty("value", sel.Value.Text)
			return lbl
		},
		"btn_close": func() ClientComponent {
			return &Icon{
				BaseComponent: BaseComponent{
					Id:           sel.Id + "_" + name,
					Name:         name,
					EventURL:     sel.EventURL,
					Target:       sel.Target,
					OnResponse:   sel.response,
					RequestValue: sel.RequestValue,
					RequestMap:   sel.RequestMap,
					Class:        []string{"close-icon"},
				},
				Value: "Times",
			}
		},
		"filter_value": func() ClientComponent {
			return ccInp(ut.ToString(sel.Data["filter_value"], ""))
		},
		"selector_result": func() ClientComponent {
			return &Table{
				BaseComponent: BaseComponent{
					Id:           sel.Id + "_" + name,
					Name:         name,
					EventURL:     sel.EventURL,
					OnResponse:   sel.response,
					RequestValue: sel.RequestValue,
					RequestMap:   sel.RequestMap,
				},
				Rows:              sel.Rows,
				Fields:            sel.Fields,
				Pagination:        PaginationTypeTop,
				PageSize:          5,
				HidePaginatonSize: true,
				TableFilter:       false,
				AddItem:           false,
				RowSelected:       true,
			}
		},
	}
	cc := ccMap[name]()
	html, err = cc.Render()
	return html, err
}

/*
Based on the values, it will generate the html code of the [Selector] or return with an error message.
*/
func (sel *Selector) Render() (html template.HTML, err error) {
	sel.InitProps(sel)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(sel.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(sel.Class, " ")
		},
		"selectorComponent": func(name string) (template.HTML, error) {
			return sel.getComponent(name)
		},
	}
	tpl := `<div id="{{ .Id }}" name="{{ .Name }}" class="selector row {{ customClass }}{{ if .Full }} full{{ end }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	>{{ if eq .Disabled false }}<div class="cell" style="width: 39px;" >{{ selectorComponent "btn_modal" }}</div>{{ end }}
	{{ if and (eq .Disabled false) (.IsNull) }}<div class="cell" style="width: 39px;" >{{ selectorComponent "btn_delete" }}</div>{{ end }}
	<div class="cell" >{{ selectorComponent "selector_text" }}</div>
	{{ if .ShowModal }}<div class="modal"><div class="dialog"><div class="panel">
	<div class="panel-title">
	<div class="cell title-cell"><span>{{ .Title }}</span></div>
	<div class="cell align-right">{{ selectorComponent "btn_close" }}</div>
	</div>
	<div class="section" >
	<div class="row full container" >
	<div class="cell">{{ selectorComponent "filter_value" }}</div>
	<div class="cell" style="width: 20px;" >{{ selectorComponent "btn_search" }}</div>
	</div>
	<div class="row full container" >{{ selectorComponent "selector_result" }}</div>
	</div>
	</div></div></div>{{ end }}
	</div>`

	return ut.TemplateBuilder("selector", tpl, funcMap, sel)
}

var testSelectorResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
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
	switch evt.Name {
	case SelectorEventShowModal:
		customModal := ut.ToBoolean(evt.Trigger.GetProperty("custom_modal"), false)
		if customModal {
			return toast("Custom modal event triggered")
		}
	case SelectorEventSearch:
		return toast(ut.ToString(evt.Value, ""))
	case SelectorEventLink:
		return toast(evt.Value.(SelectOption).Text)
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
	return evt
}

var testSelectorFields []TableField = []TableField{
	{Name: "custname", Label: "Customer Name"},
	{Name: "custnumber", Label: "Customer Number"},
	{Name: "city", Label: "City"},
	{Name: "street", Label: "Street"},
}

var testSelectorRows []ut.IM = []ut.IM{
	{
		"city":       "City1",
		"custname":   "First Customer Co.",
		"custnumber": "DMCUST/00001",
		"id":         "customer-2",
		"label":      "First Customer Co.",
		"street":     "street 1.",
	},
	{
		"city":       "City3",
		"custname":   "Second Customer Name",
		"custnumber": "DMCUST/00002",
		"id":         "customer-3",
		"label":      "Second Customer Name",
		"street":     "street 3.",
		"deleted":    1,
	},
	{
		"city":       "City4",
		"custname":   "Third Customer Foundation",
		"custnumber": "DMCUST/00003",
		"id":         "customer-4",
		"label":      "Third Customer Foundation",
		"street":     "street 4.",
	},
}

// [Selector] test and demo data
func TestSelector(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeList,
			Component: &Selector{
				BaseComponent: BaseComponent{
					Id:           id + "_selector_default",
					EventURL:     eventURL,
					OnResponse:   testSelectorResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Data: ut.IM{
						"filter_value": "default value",
					},
				},
				Value:  SelectOption{Value: "12345", Text: "Customer Name"},
				Fields: testSelectorFields,
				Rows:   testSelectorRows,
			},
		},
		{
			Label:         "Disabled",
			ComponentType: ComponentTypeList,
			Component: &Selector{
				BaseComponent: BaseComponent{
					Id:           id + "_selector_disabled",
					EventURL:     eventURL,
					OnResponse:   testSelectorResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value:    SelectOption{Value: "12345", Text: "Customer Name"},
				Link:     true,
				Full:     true,
				Disabled: true,
			},
		},
		{
			Label:         "Is null, link",
			ComponentType: ComponentTypeList,
			Component: &Selector{
				BaseComponent: BaseComponent{
					Id:           id + "_selector_link",
					EventURL:     eventURL,
					OnResponse:   testSelectorResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value:  SelectOption{Value: "12345", Text: "Customer Name"},
				Fields: testSelectorFields,
				Rows:   testSelectorRows,
				Link:   true,
				IsNull: true,
				Full:   true,
			},
		},
		{
			Label:         "Custom modal event and icon",
			ComponentType: ComponentTypeList,
			Component: &Selector{
				BaseComponent: BaseComponent{
					Id:           id + "_selector_custom",
					EventURL:     eventURL,
					OnResponse:   testSelectorResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Value:       SelectOption{Value: "12345", Text: "Customer Name"},
				Link:        true,
				IsNull:      false,
				Full:        true,
				CustomModal: true,
				ModalIcon:   IconBolt,
			},
		},
	}
}
