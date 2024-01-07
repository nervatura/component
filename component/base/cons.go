package base

const (
	ThemeLight = "light"
	ThemeDark  = "dark"

	ComponentTypeToast = "toast"

	ComponentTypeButton      = "button"
	ComponentTypeDateTime    = "datetime"
	ComponentTypeIcon        = "icon"
	ComponentTypeInput       = "input"
	ComponentTypeLabel       = "label"
	ComponentTypeNumberInput = "number"
	ComponentTypeSelect      = "select"

	ComponentTypeTable      = "table"
	ComponentTypePagination = "pagination"
	ComponentTypeMenuBar    = "menubar"

	ComponentTypeLogin = "login"

	TextAlignLeft   = "align-left"
	TextAlignCenter = "center"
	TextAlignRight  = "align-right"

	VerticalAlignTop    = "top"
	VerticalAlignMiddle = "middle"
	VerticalAlignBottom = "bottom"

	SideVisibilityAuto = "auto"
	SideVisibilityShow = "show"
	SideVisibilityHide = "hide"

	SwapInnerHTML   = "innerHTML"   // Replace the inner html of the target element
	SwapOuterHTML   = "outerHTML"   // Replace the entire target element with the response
	SwapBeforeBegin = "beforebegin" // Insert the response before the target element
	SwapAfterBegin  = "afterbegin"  // Insert the response before the first child of the target element
	SwapBeforeEnd   = "beforeend"   // Insert the response after the last child of the target element
	SwapAfterEnd    = "afterend"    // Insert the response after the target element
	SwapDelete      = "delete"      // Deletes the target element regardless of the response
	SwapNone        = "none"        // Does not append content from response

	HeaderLocation           = "HX-Location"             //allows you to do a client-side redirect that does not do a full page reload
	HeaderPushUrl            = "HX-Push-Url"             //pushes a new url into the hidemo stack
	HeaderRedirect           = "HX-Redirect"             //can be used to do a client-side redirect to a new location
	HeaderRefresh            = "HX-Refresh"              //if set to “true” the client-side will do a full refresh of the page
	HeaderReplaceUrl         = "HX-Replace-Url"          //replaces the current URL in the location bar
	HeaderReswap             = "HX-Reswap"               //allows you to specify how the response will be swapped. See hx-swap for possible values
	HeaderRetarget           = "HX-Retarget"             //a CSS selector that updates the target of the content update to a different element on the page
	HeaderReselect           = "HX-Reselect"             //a CSS selector that allows you to choose which part of the response is used to be swapped in. Overrides an existing hx-select on the triggering element
	HeaderTrigger            = "HX-Trigger"              //allows you to trigger client-side events
	HeaderTriggerAfterSettle = "HX-Trigger-After-Settle" //allows you to trigger client-side events after the settle step
	HeaderTriggerAfterSwap   = "HX-Trigger-After-Swap"   //allows you to trigger client-side events after the swap step

	IndicatorNone    = ""
	IndicatorSpinner = "spinner"
)

var Theme []string = []string{ThemeLight, ThemeDark}
var ComponentType []string = []string{
	ComponentTypeButton, ComponentTypeDateTime, ComponentTypeIcon, ComponentTypeInput,
	ComponentTypeLabel, ComponentTypeNumberInput, ComponentTypeSelect, ComponentTypeToast,
	ComponentTypeTable, ComponentTypePagination, ComponentTypeMenuBar,
	ComponentTypeLogin,
}
var TextAlign []string = []string{TextAlignLeft, TextAlignCenter, TextAlignRight}
var VerticalAlign []string = []string{VerticalAlignTop, VerticalAlignMiddle, VerticalAlignBottom}
var SideVisibility []string = []string{SideVisibilityAuto, SideVisibilityShow, SideVisibilityHide}

var Swap []string = []string{SwapInnerHTML, SwapOuterHTML, SwapBeforeBegin, SwapAfterBegin,
	SwapBeforeEnd, SwapAfterEnd, SwapDelete, SwapNone}
var Indicator []string = []string{IndicatorNone, IndicatorSpinner}
var Header []string = []string{
	HeaderLocation, HeaderPushUrl, HeaderRedirect, HeaderRefresh, HeaderReplaceUrl, HeaderReswap,
	HeaderRetarget, HeaderReselect, HeaderTrigger, HeaderTriggerAfterSettle, HeaderTriggerAfterSwap}
