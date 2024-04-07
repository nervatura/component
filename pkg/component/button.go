package component

import (
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [Button] constants
const (
	ComponentTypeButton = "button"

	ButtonEventClick = "click"

	ButtonTypeDefault = ""
	ButtonTypePrimary = "primary"
	ButtonTypeBorder  = "border"
)

// [Button] Type values
var ButtonType []string = []string{ButtonTypeDefault, ButtonTypePrimary, ButtonTypeBorder}

/*
Creates an HTML button control

For example:

	&Button{
	  BaseComponent: BaseComponent{
	    Id:           "id_button_primary",
	    EventURL:     "/event",
	    OnResponse:   func(evt ResponseEvent) (re ResponseEvent) {
	      // return parent_component response
	      return evt
	    },
	    RequestValue: parent_component.GetProperty("request_value").(map[string]ut.IM),
	    RequestMap:   parent_component.GetProperty("request_map").(map[string]ClientComponent)
	  },
	  Type:     ButtonTypePrimary,
	  Align:    TextAlignCenter,
	  Label:    "Primary",
	  Icon:     "Check",
	  Selected: true
	}
*/
type Button struct {
	BaseComponent
	/* [ButtonType] variable constants: [ButtonTypeDefault], [ButtonTypePrimary], [ButtonTypeBorder].
	Default value: [ButtonTypeDefault] */
	Type string `json:"type"`
	/* [TextAlign] variable constants: [TextAlignLeft], [TextAlignCenter], [TextAlignRight].
	Default value: [TextAlignCenter] */
	Align string `json:"align"`
	// The HTML title, aria-label attribute and button caption of the component
	Label string `json:"label"`
	// Any component to be displayed (e.g. [Label] component) instead of the label text
	LabelComponent ClientComponent `json:"label_component"`
	// Valid [Icon] component value. See more [IconKey] variable values.
	Icon string `json:"icon"`
	// Specifies that the button should be disabled
	Disabled bool `json:"disabled"`
	// Specifies that the button should automatically get focus when the page loads
	AutoFocus bool `json:"auto_focus"`
	// Full width button (100%)
	Full bool `json:"full"`
	// Sets the values of the small-button class style
	Small bool `json:"small"`
	// Sets the values of the selected class style
	Selected bool `json:"selected"`
	// The button label is visible or invisible
	HideLabel bool `json:"hide_label"`
	// The badge value of the button
	Badge int64 `json:"badge"`
	// The button badge is visible or invisible
	ShowBadge bool `json:"show_badge"`
}

/*
Returns all properties of the [Button]
*/
func (btn *Button) Properties() ut.IM {
	return ut.MergeIM(
		btn.BaseComponent.Properties(),
		ut.IM{
			"type":            btn.Type,
			"align":           btn.Align,
			"label":           btn.Label,
			"label_component": btn.LabelComponent,
			"icon":            btn.Icon,
			"disabled":        btn.Disabled,
			"auto_focus":      btn.AutoFocus,
			"full":            btn.Full,
			"small":           btn.Small,
			"selected":        btn.Selected,
			"hide_label":      btn.HideLabel,
			"badge":           btn.Badge,
			"show_badge":      btn.ShowBadge,
		})
}

/*
Returns the value of the property of the [Button] with the specified name.
*/
func (btn *Button) GetProperty(propName string) interface{} {
	return btn.Properties()[propName]
}

/*
It checks the value given to the property of the [Button] and always returns a valid value
*/
func (btn *Button) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"type": func() interface{} {
			return btn.CheckEnumValue(ut.ToString(propValue, ""), ButtonTypeDefault, ButtonType)
		},
		"align": func() interface{} {
			return btn.CheckEnumValue(ut.ToString(propValue, ""), TextAlignCenter, TextAlign)
		},
		"indicator": func() interface{} {
			return btn.CheckEnumValue(ut.ToString(propValue, IndicatorSpinner), IndicatorSpinner, Indicator)
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if btn.BaseComponent.GetProperty(propName) != nil {
		return btn.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [Button] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (btn *Button) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"type": func() interface{} {
			btn.Type = btn.Validation(propName, propValue).(string)
			return btn.Type
		},
		"align": func() interface{} {
			btn.Align = btn.Validation(propName, propValue).(string)
			return btn.Align
		},
		"label": func() interface{} {
			btn.Label = ut.ToString(propValue, "")
			return btn.Label
		},
		"label_component": func() interface{} {
			if cc, valid := propValue.(ClientComponent); valid {
				btn.LabelComponent = cc
			}
			return btn.LabelComponent
		},
		"icon": func() interface{} {
			btn.Icon = ut.ToString(propValue, "")
			return btn.Icon
		},
		"disabled": func() interface{} {
			btn.Disabled = ut.ToBoolean(propValue, false)
			return btn.Disabled
		},
		"auto_focus": func() interface{} {
			btn.AutoFocus = ut.ToBoolean(propValue, false)
			return btn.AutoFocus
		},
		"full": func() interface{} {
			btn.Full = ut.ToBoolean(propValue, false)
			return btn.Full
		},
		"small": func() interface{} {
			btn.Small = ut.ToBoolean(propValue, false)
			return btn.Small
		},
		"selected": func() interface{} {
			btn.Selected = ut.ToBoolean(propValue, false)
			return btn.Selected
		},
		"hide_label": func() interface{} {
			btn.HideLabel = ut.ToBoolean(propValue, false)
			return btn.HideLabel
		},
		"badge": func() interface{} {
			btn.Badge = ut.ToInteger(propValue, 0)
			return btn.Badge
		},
		"show_badge": func() interface{} {
			btn.ShowBadge = ut.ToBoolean(propValue, false)
			return btn.ShowBadge
		},
		"indicator": func() interface{} {
			btn.Indicator = btn.Validation(propName, propValue).(string)
			return btn.Indicator
		},
	}
	if _, found := pm[propName]; found {
		return btn.SetRequestValue(propName, pm[propName](), []string{})
	}
	if btn.BaseComponent.GetProperty(propName) != nil {
		return btn.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

/*
If the OnResponse function of the [Button] is implemented, the function calls it after the [TriggerEvent]
is processed, otherwise the function's return [ResponseEvent] is the processed [TriggerEvent].
*/
func (btn *Button) OnRequest(te TriggerEvent) (re ResponseEvent) {
	value := btn.SetProperty("value", te.Values.Get(te.Id+"_value"))
	evt := ResponseEvent{
		Trigger: btn, TriggerName: btn.Name,
		Name:  ButtonEventClick,
		Value: value,
	}
	if btn.OnResponse != nil {
		return btn.OnResponse(evt)
	}
	return evt
}

func (btn *Button) getComponent(name string) (string, error) {
	ccMap := map[string]func() ClientComponent{
		"icon": func() ClientComponent {
			return &Icon{Value: btn.Icon, Width: 20}
		},
		"label": func() ClientComponent {
			return btn.LabelComponent
		},
	}
	return ccMap[name]().Render()
}

/*
Based on the values, it will generate the html code of the [Button] or return with an error message.
*/
func (btn *Button) Render() (res string, err error) {
	btn.InitProps(btn)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(btn.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(btn.Class, " ")
		},
		"buttonComponent": func(name string) (string, error) {
			return btn.getComponent(name)
		},
	}
	tpl := `<button id="{{ .Id }}" name="{{ .Name }}"
	{{ if or (eq .Type "primary") (eq .Type "border") }} button-type="{{ .Type }}"{{ end }}
	{{ if ne .EventURL "" }} hx-post="{{ .EventURL }}" hx-target="{{ .Target }}" hx-swap="{{ .Swap }}"{{ end }}
	{{ if ne .Indicator "" }} hx-indicator="#{{ .Indicator }}"{{ end }}
	{{ if .Disabled }} disabled{{ end }}
	{{ if .AutoFocus }} autofocus{{ end }}
	{{ if ne .Label "" }} aria-label="{{ .Label }}" title="{{ .Label }}"{{ end }}
	 class="{{ .Align }}{{ if .Small }} small-button{{ end }}{{ if .Full }} full{{ end }}{{ if .Selected }} selected{{ end }}{{ if .HideLabel }} hidelabel{{ end }} {{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	>{{ if and (ne .Icon "") (ne .Align "align-right") }}{{ buttonComponent "icon" }}{{ end }}
	{{ if .LabelComponent }}{{ buttonComponent "label" }}{{ else }}<span>{{ .Label }}</span>{{ end }}
	{{ if and (ne .Icon "") (eq .Align "align-right") }}{{ buttonComponent "icon" }}{{ end }}
	{{ if .ShowBadge }}<span class="right" ><span class="badge{{ if .Selected }} selected-badge{{ end }}" >{{ .Badge }}</span></span>{{ end }}
	</button>`

	if res, err = ut.TemplateBuilder("button", tpl, funcMap, btn); err == nil && btn.EventURL != "" {
		btn.SetProperty("request_map", btn)
	}
	return res, err
}

var demoBtnResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
	badge := ut.ToInteger(evt.Trigger.GetProperty("badge"), 0)
	evt.Trigger.SetProperty("badge", badge+1)
	evt.Trigger.SetProperty("show_badge", true)
	return evt
}

// [Button] test and demo data
func TestButton(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeButton,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id:           id + "_button_default",
					EventURL:     eventURL,
					OnResponse:   demoBtnResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type:  ButtonTypeDefault,
				Align: TextAlignCenter,
				Label: "Default",
			}},
		{
			Label:         "Right icon",
			ComponentType: ComponentTypeButton,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id:           id + "_button_right_icon",
					EventURL:     eventURL,
					OnResponse:   demoBtnResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type:      ButtonTypeDefault,
				Align:     TextAlignRight,
				Label:     "Right icon",
				Icon:      "InfoCircle",
				HideLabel: true,
			}},
		{
			Label:         "Primary and icon",
			ComponentType: ComponentTypeButton,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id:           id + "_button_primary",
					EventURL:     eventURL,
					OnResponse:   demoBtnResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type:     ButtonTypePrimary,
				Align:    TextAlignCenter,
				Label:    "Primary",
				Icon:     "Check",
				Selected: true,
			}},
		{
			Label:         "Border and selected",
			ComponentType: ComponentTypeButton,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id:           id + "_button_border",
					EventURL:     eventURL,
					OnResponse:   demoBtnResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type:     ButtonTypeBorder,
				Align:    TextAlignCenter,
				Label:    "Border selected",
				Selected: true,
			}},
		{
			Label:         "Border full and badge",
			ComponentType: ComponentTypeButton,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id:           id + "_button_full",
					EventURL:     eventURL,
					OnResponse:   demoBtnResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type:      ButtonTypeBorder,
				Align:     TextAlignCenter,
				Label:     "Border full and badge",
				Full:      true,
				Badge:     0,
				ShowBadge: true,
			}},
		{
			Label:         "Small disabled",
			ComponentType: ComponentTypeButton,
			Component: &Button{
				Type:     ButtonTypeDefault,
				Align:    TextAlignCenter,
				Label:    "Small disabled",
				Small:    true,
				Disabled: true,
			}},
		{
			Label:         "Label component",
			ComponentType: ComponentTypeButton,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id:           id + "_button_label",
					EventURL:     eventURL,
					OnResponse:   demoBtnResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Type:           ButtonTypeDefault,
				Align:          TextAlignCenter,
				Label:          "Label component",
				LabelComponent: &Icon{Value: "Print", Width: 32, Height: 32},
			}},
		{
			Label:         "Border and custom style",
			ComponentType: ComponentTypeButton,
			Component: &Button{
				BaseComponent: BaseComponent{
					Id:           id + "_button_custom",
					EventURL:     eventURL,
					OnResponse:   demoBtnResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Style: ut.SM{
						"border-color": "green", "color": "red", "border-radius": "3px",
					},
				},
				Type:  ButtonTypeBorder,
				Align: TextAlignCenter,
				Label: "Button style",
			}},
	}
}
