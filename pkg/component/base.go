/*
Server-side Go components

# Benefits of server-side components

Component based development is an approach to software development that focuses on the design and development
of reusable components. Server components are also reusable bits of code, but are compiled into HTML before
the browser sees them. The server-side components tend to perform better. When the page that a browser receives
contains everything it needs for presentation, it’s going to be able to deliver that presentation to the user
much quicker.

  - The development of a client-side application and component takes place in a very complex ecosystem. An
    average node_modules size can be hundreds of MB and contain hundreds or even over a thousand different packages.
    Each of these also means external dependencies of varying quality and reliability, which is also a big security
    risk. In addition, the constant updating and tracking of these different packages and the dozens of frameworks
    and technologies based on them requires a lot of resources.
    A server-side component has no external dependencies. They can be easily created within the technical capabilities
    of a given server-side language. Their maintenance needs are limited to their actual code, which is very small and
    much safer due to the lack of external dependencies.

  - The language of the client-side components is basically javascript, but most are server-side languages, so go is
    a much more efficient and safer programming language. JavaScript is originally an add-on to html code and browsers,
    which was originally created to increase the efficiency of the user interface and not to develop the codebase of
    complex programs. It is possible to partially replace it during development with, for example, the typescript language,
    but this also means additional dependencies and the complexity of the development ecosystem, the end result of which
    will be a javascript code base. This practically means that code written in a programming language is translated into
    the code of another language and the content to be displayed is created during its execution. There are many intermediate
    steps, used resources, potential for errors, security risks and uncertainties in the process. With the server-side
    components, it is possible to simply write the program code in an easy-to-use and safe language, the end result of which
    is the html content to be displayed.

  - Client-side components usually communicate with the server using a JSON-based REST API and receive the data to be
    displayed. This also means that the data retrieval must adapt to the data structure of the REST API, so the database
    data must first be converted to this structure, and then reprocessed on the client side for final display. In addition
    to possible changes to the data structure, this also means JSON encoding and decoding in all cases. The server-side components
    can directly access the database and use the data immediately in the data structure to be displayed. This also means faster
    rendering and better resource management for the server-side components.

# Nervatura components

Server components can be written in any server-side language. This enables you to write your client in the
same language as your server application’s logic.
On the user side, an application that is loaded in the browser in html syntax is a set of components
that are hierarchically related to each other. Any component of the application may be able to send a
request to the server, and depending on the processing of the request, any part of the application may
change. The entire page is not replaced or reloaded in the browser, only the required parts of the application.
The components do not use json data format to send data, all data is sent in URL-encoded form. All data of the
application is stored on the server, and the components do not contain javascript code.

  - Nervatura components use the htmx library for direct communication with the server. Htmx is small (~14k),
    dependency-free, browser-oriented javascript library that allows you to access modern browser
    features directly from HTML, rather than using javascript. The server-side components use only a small part
    of the possibilities of htmx. More information about htmx can be found on the https://htmx.org link.
    The [Application] component contains and automatically loads the appropriate version of htmx when used.

  - Nervatura components are not a framework, they use only the built-in packages of go and have no external
    dependencies. It is a library of components whose elements can be freely combined with each other and can
    be easily further developed. A Nervatura component is actually just a code implementation proposal that
    anyone can easily create a server-side component in any program language.

# Using the components

The basis of the components is a data-driven template, which contains the rules of their appearance and
behavior. The Nervatura components use the declarative template from the text/template package in the go
standard library. In traditional server-side programming, the same package is the basis for generating
the html code displayed on the client side.

The main difference is that templates do not contain the rules
of entire pages or certain parts of them, but the generated independent small html code fragments are
organized based on a different logic. These hierarchically organized, logically independent html code parts
are able to independently react and process user interactions and notify other components of the processing
results.

React's JSX or Lit's declarative template works on the same principle, but the same logic can be
implemented in any other server-side programming language's template package with similar functionality.

  - The [ClientComponent] interface contains the mandatory functions that are available in all Nervatura components,
    and each Nervatura component is also of the [ClientComponent] type in addition to its own type. The
    [BaseComponent] type is likewise a [ClientComponent], which, in addition to the mandatory functions of the
    interface, also contains properties and functions that all components based on it inherit and can use if necessary.

  - The components receive the data of the user interface event in a [TriggerEvent] type, and then return a
    [ResponseEvent] type data package after processing it. For complex components, the [TriggerEvent] data is always
    received by the smallest independent component in whose scope the event occurred. The component above it in the
    hierarchy always has the option to decide whether to forward the [ResponseEvent] response received from its child,
    or to continue processing it and replace it with its own [ResponseEvent] response. The html code fragment returned
    to the user interface as a response replaces only the part of the entire application that is covered by the scope
    of the processed response. The html code displayed in the browser is a dynamically changing display of the current
    data stored on the server side of the application.

  - The easiest way to create a new component is to add new properties and functions to the [BaseComponent]
    and override the functions of the [ClientComponent] interface if necessary. The [Application] component is a
    top-level element to which all other components belong. This element is completely never replaced, only some of
    its parts can change. Its task is to load and manage all static elements required for the operation and display
    of the components, such as style sheets and the htmx package. With the help of [TestComponent], we can define
    various test cases for the new component, which can be used to visually test the component in addition to unit
    tests.

# Examples and demo application

The logic of the components can be understood most easily from the code of the existing components.

  - The [NumberInput] component, for example, is a simple elemental component that directly receives events from
    the client interface via [NumberInput.OnRequest]. It validates and stores the new value received in [TriggerEvent],
    then returns it directly in the [ResponseEvent] event or transmits its new state to the parent component via the
    OnResponse function.

  - [Login] and [Locale] are complex components that do not receive their own OnRequest event, but receive the
    processed events of the client interface through their child components. In addition to the mandatory functions of
    [ClientComponent], they also have their own functions (for example, getComponent, response, msg) with which
    they initialize their own child components and handle the [ResponseEvent] events of the components.

  - The css files of the components are included in the component/pkg/static package. The index.css
    contains the reference to the style sheets of all components, so it is sufficient to specify this in the
    [Application] HeadLink property. The styles of the new components can be specified in additional css files and
    the styles of the existing components can also be overwritten.

The [pkg/github.com/nervatura/component/pkg/demo] package includes an example application that displays all
components with sample data. Applications can store component data in memory, but they can save it anywhere in json
format and load it back. The demo application can store session data in memory and as session files. The source code
of the example application also contains an example of using a session database.

The Admin interface of the application https://github.com/nervatura/nervatura is another example of
the use of server-side components (session and JWT token, database session and more).
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
