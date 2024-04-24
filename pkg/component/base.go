/*
Server-side Go components

An easy way to create a server-side component in any programming language.
Detailed description and documentation:: https://nervatura.github.io/component
*/
package component

import (
	"net/url"
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// Common component constants
const (
	ThemeLight = "light"
	ThemeDark  = "dark"

	TextAlignLeft   = "align-left"
	TextAlignCenter = "center"
	TextAlignRight  = "align-right"

	VerticalAlignTop    = "top"
	VerticalAlignMiddle = "middle"
	VerticalAlignBottom = "bottom"
)

// [BaseComponent] constants
const (
	BaseEventValue = "value"

	// Replace the inner html of the target element
	SwapInnerHTML = "innerHTML"
	// Replace the entire target element with the response
	SwapOuterHTML = "outerHTML"
	// Insert the response before the target element
	SwapBeforeBegin = "beforebegin"
	// Insert the response before the first child of the target element
	SwapAfterBegin = "afterbegin"
	// Insert the response after the last child of the target element
	SwapBeforeEnd = "beforeend"
	// Insert the response after the target element
	SwapAfterEnd = "afterend"
	// Deletes the target element regardless of the response
	SwapDelete = "delete"
	// Does not append content from response
	SwapNone = "none"

	IndicatorNone    = ""
	IndicatorSpinner = "spinner"

	PaginationTypeTop    = "top"
	PaginationTypeBottom = "bottom"
	PaginationTypeAll    = "all"
	PaginationTypeNone   = "none"
)

// [ResponseEvent] Header map key constants
const (
	// Allows you to do a client-side redirect that does not do a full page reload
	HeaderLocation = "HX-Location"
	// Pushes a new url into the hidemo stack
	HeaderPushUrl = "HX-Push-Url"
	// Can be used to do a client-side redirect to a new location
	HeaderRedirect = "HX-Redirect"
	// If set to “true” the client-side will do a full refresh of the page
	HeaderRefresh = "HX-Refresh"
	// Replaces the current URL in the location bar
	HeaderReplaceUrl = "HX-Replace-Url"
	// Allows you to specify how the response will be swapped. See hx-swap for possible values
	HeaderReswap = "HX-Reswap"
	// A CSS selector that updates the target of the content update to a different element on the page
	HeaderRetarget = "HX-Retarget"
	/* A CSS selector that allows you to choose which part of the response is used to be swapped in.
	Overrides an existing hx-select on the triggering element */
	HeaderReselect = "HX-Reselect"
	// Allows you to trigger client-side events
	HeaderTrigger = "HX-Trigger"
	// Allows you to trigger client-side events after the settle step
	HeaderTriggerAfterSettle = "HX-Trigger-After-Settle"
	// Allows you to trigger client-side events after the swap step
	HeaderTriggerAfterSwap = "HX-Trigger-After-Swap"
)

// Component Theme values
var Theme []string = []string{ThemeLight, ThemeDark}

// Component TextAlign values
var TextAlign []string = []string{TextAlignLeft, TextAlignCenter, TextAlignRight}

// Component VerticalAlign values
var VerticalAlign []string = []string{VerticalAlignTop, VerticalAlignMiddle, VerticalAlignBottom}

// Component PaginationType values
var PaginationType []string = []string{PaginationTypeTop, PaginationTypeBottom, PaginationTypeAll, PaginationTypeNone}

// [BaseComponent] Swap values
var Swap []string = []string{SwapInnerHTML, SwapOuterHTML, SwapBeforeBegin, SwapAfterBegin,
	SwapBeforeEnd, SwapAfterEnd, SwapDelete, SwapNone}

// [BaseComponent] Indicator values
var Indicator []string = []string{IndicatorNone, IndicatorSpinner}

// [ResponseEvent] Header map keys
var Header []string = []string{
	HeaderLocation, HeaderPushUrl, HeaderRedirect, HeaderRefresh, HeaderReplaceUrl, HeaderReswap,
	HeaderRetarget, HeaderReselect, HeaderTrigger, HeaderTriggerAfterSettle, HeaderTriggerAfterSwap}

// Generic server-side component type. All components must implement these functions.
type ClientComponent interface {
	Properties() ut.IM /*
		Returns all properties of the component
	*/
	Validation(propName string, propValue interface{}) interface{} /*
		It checks the value given to the property and always returns a valid value
	*/
	GetProperty(propName string) interface{} /*
		Returns the value of the property with the specified name.
	*/
	SetProperty(propName string, propValue interface{}) interface{} /*
		Setting a property value safely. Checks the entered value.
		In case of an invalid value, the default value will be set.
	*/
	InitProps(cc ClientComponent) /*
		Checks the value of all properties of the component.
		If a value is missing or invalid, it will set the default value. The Render function calls it automatically.
	*/
	Render() (res string, err error) /*
		Based on the values, it will generate the component's html code or return with an error message.
		The InitProps function is automatically called at the beginning of the function.
	*/
	OnRequest(te TriggerEvent) (re ResponseEvent) /*
		This function is called in the event that the component also receives user interface events directly.
		As a rule, it is used by the basic components (button, input, etc.).
		Composite components receive these events not directly, but through their child components.
		If the component's OnResponse function is implemented, the function calls it after the [TriggerEvent] is processed,
		otherwise the function's return [ResponseEvent] is the processed [TriggerEvent].
	*/
}

// Test container for component test cases
type TestComponent struct {
	// The name of the test data
	Label string `json:"label"`
	// Type of tested component. Example: [ComponentTypeSelect]
	ComponentType string `json:"component_type"`
	// The tested component with the test data
	Component ClientComponent `json:"component"`
}

// Event data sent by the user interface in the htmx request
type TriggerEvent struct {
	// The id attribute of the component that receives the user event. htmx request header: HX-Trigger
	Id string `json:"id"`
	// The name attribute of the component that receives the user event. htmx request header: HX-Trigger-Name
	Name string `json:"name"`
	// The hx-target attribute of the component that receives the user event. htmx request header: HX-Target
	Target string `json:"target"`
	// The URL-encoded data of the request
	Values url.Values `json:"values"`
}

// Response data for a user event
type ResponseEvent struct {
	// The data processing and responding component
	Trigger ClientComponent `json:"trigger"`
	// The name of the component
	TriggerName string `json:"trigger_name"`
	// The name of the event. Example: [InputEventChange]
	Name string `json:"name"`
	// The value of the event.
	Value interface{} `json:"value"`
	/* htmx supports some htmx-specific response headers. See more [ResponseEvent] Header map key constants
	Example: Header: ut.SM{HeaderRetarget: "#toast-msg", HeaderReswap: SwapInnerHTML} */
	Header ut.SM `json:"header"`
}

// A component whose properties and functions are contained in all other components.
type BaseComponent struct {
	// Unique identifier of a component
	Id string `json:"id"`
	// The name attribute of the component
	Name string `json:"name"`
	// The htmx hx-post attribute can be sent form-data as an HTTP POST request to the EventURL value
	EventURL string `json:"event_url"`
	/*
		The htmx hx-target attribute allows you to target a different element for swapping than the one
		issuing the AJAX request. The value of this attribute can be:
		- "this" which indicates that the element that the hx-target attribute is on is the target
		- any [component ID]
		Default value: "this"
	*/
	Target string `json:"target"`
	/*
		The htmx hx-swap attribute allows you to specify how the response will be swapped in relative
		to the target of an AJAX request. If you do not specify the option,
		the default is [SwapOuterHTML]. See more [Swap] variable constants.
	*/
	Swap string `json:"swap"`
	/*
		The htmx hx-indicator can be used to show spinners or progress indicators while the request is in flight.
		[Indicator] variable constants: [IndicatorNone], [IndicatorSpinner]. Default value: [IndicatorNone]
	*/
	Indicator string `json:"indicator"`
	/*
		A list of custom class names for the component's HTML attribute. The names are added to the list of
		component class names
	*/
	Class []string `json:"class"`
	// The values of the style HTML attribute of the component. Example: ut.SM{"color":"red","padding":"8px"}
	Style ut.SM `json:"style"`
	/*
		Any additional data that can be associated with the component and must be stored with it.
		Example: ut.IM{"value":"value","number":12345}
	*/
	Data ut.IM `json:"data"`
	/*
		The current value of all properties of the component that will be set after the component is created.
		The saved session data is reloaded based on the map data.
	*/
	RequestValue map[string]ut.IM `json:"request_value"`
	/*
		A map in which all the components of the entire hierarchy can be found based on their ID, which have
		their own OnRequest function and can be called directly from the user interface.
	*/
	RequestMap map[string]ClientComponent `json:"-"`
	/*
		The control function of the component's parent component. If this is set, the component passes
		the [ResponseEvent] of its event to this function.
	*/
	OnResponse func(evt ResponseEvent) (re ResponseEvent) `json:"-"`
	init       bool
}

/*
Returns all properties of the [BaseComponent]
*/
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

/*
Checking the value of an enum type. In case of an invalid value, the function returns the value specified in defaultValue
*/
func (bcc *BaseComponent) CheckEnumValue(value, defaultValue string, enums []string) string {
	if ut.Contains(enums, value) {
		return value
	}
	return defaultValue
}

/*
It checks the value given to the property of the [BaseComponent] and always returns a valid value
*/
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

/*
Returns the value of the property of the [BaseComponent] with the specified name.
*/
func (bcc *BaseComponent) GetProperty(propName string) interface{} {
	return bcc.Properties()[propName]
}

/*
Setting a property of the [BaseComponent] value safely. Checks the entered value.
In case of an invalid value, the default value will be set. The function automatically
calls the [BaseComponent.SetRequestValue] function with the value.
*/
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

/*
The function sets the value of a changed property of the component in the RequestValue map.
[BaseComponent.SetProperty] calls this function as well at the end of processing.
*/
func (bcc *BaseComponent) SetRequestValue(propName string, propValue interface{}, staticFields []string) interface{} {
	if !ut.Contains(staticFields, propName) && bcc.Id != "" && !bcc.init {
		bcc.RequestValue = bcc.Validation("request_value", map[string]ut.IM{bcc.Id: {propName: propValue}}).(map[string]ut.IM)
	}
	return propValue
}

/*
Checks the value of all properties of the [BaseComponent]. If a value is missing or invalid, it will set the
default value. With the help of the RequestValue map, the function can also restore the saved session data of
the component. The Render function calls it automatically.
*/
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

/*
Based on the values, it will generate the html code of the [BaseComponent] or return with an error message.
The [BaseComponent.InitProps] function is automatically called at the beginning of the function.
*/
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

/*
If the OnResponse function of the [BaseComponent] is implemented, the function calls it after the [TriggerEvent]
is processed, otherwise the function's return [ResponseEvent] is the processed [TriggerEvent].
*/
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
