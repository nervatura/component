package component

import (
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [Link] constants
const (
	ComponentTypeLink = "link"

	LinkEventClick = "click"

	LinkStyleDefault = ""
	LinkStyleButton  = "button"
	LinkStylePrimary = "primary"
	LinkStyleBorder  = "border"
)

// [Link] LinkStyle values
var LinkStyle []string = []string{LinkStyleDefault, LinkStyleButton, LinkStylePrimary, LinkStyleBorder}

// It defines a hyperlink, which is used to link from one page to another.
type Link struct {
	BaseComponent
	// Specifies the URL of the page the link goes to
	Href string `json:"href"`
	// Specifies that the target will be downloaded when a user clicks on the hyperlink
	Download string `json:"download"`
	// Specifies what media/device the linked document is optimized for
	Media string `json:"media"`
	// Specifies a space-separated list of URLs to which, when the link is followed, post requests with the body ping will be sent by the browser (in the background).
	// Typically used for tracking.
	Ping string `json:"ping"`
	// Specifies which referrer information to send with the link
	// Valid values: "no-referrer", "no-referrer-when-downgrade", "origin", "origin-when-cross-origin", "same-origin", "strict-origin-when-cross-origin", "unsafe-url"
	ReferrerPolicy string `json:"referrerpolicy"`
	// Specifies the relationship between the current document and the linked document
	// Valid values: "alternate", "author", "bookmark", "external", "help", "license", "next", "nofollow", "noreferrer", "noopener", "prev", "search", "tag"
	Rel string `json:"rel"`
	// Specifies where to open the linked document
	// Valid values: "_blank", "_parent", "_self", "_top"
	LinkTarget string `json:"link_target"`
	// Specifies the media type of the linked document
	MediaType string `json:"media_type"`
	/* [LinkStyle] variable constants: [LinkStyleDefault], [LinkStyleButton], [LinkStylePrimary], [LinkStyleBorder].
	Default value: [LinkStyleDefault] */
	LinkStyle string `json:"link_style"`
	/* [TextAlign] variable constants: [TextAlignLeft], [TextAlignCenter], [TextAlignRight].
	Default value: [TextAlignCenter] */
	Align string `json:"align"`
	// The HTML title, aria-label attribute and link caption of the component
	Label string `json:"label"`
	// Any component to be displayed (e.g. [Label] component) instead of the label text
	LabelComponent ClientComponent `json:"label_component"`
	// Valid [Icon] component value. See more [IconKey] variable values.
	Icon string `json:"icon"`
	// Specifies that the link should be disabled
	Disabled bool `json:"disabled"`
	// Specifies that the link should automatically get focus when the page loads
	AutoFocus bool `json:"auto_focus"`
	// Full width link (100%)
	Full bool `json:"full"`
	// Sets the values of the small-link class style
	Small bool `json:"small"`
	// Sets the values of the selected class style
	Selected bool `json:"selected"`
	// The link label is visible or invisible
	HideLabel bool `json:"hide_label"`
	// The badge value of the button
	Badge int64 `json:"badge"`
	// The button badge is visible or invisible
	ShowBadge bool `json:"show_badge"`
}

/*
Returns all properties of the [Link]
*/
func (lnk *Link) Properties() ut.IM {
	return ut.MergeIM(
		lnk.BaseComponent.Properties(),
		ut.IM{
			"href":            lnk.Href,
			"download":        lnk.Download,
			"media":           lnk.Media,
			"ping":            lnk.Ping,
			"referrerpolicy":  lnk.ReferrerPolicy,
			"rel":             lnk.Rel,
			"link_target":     lnk.LinkTarget,
			"media_type":      lnk.MediaType,
			"link_style":      lnk.LinkStyle,
			"align":           lnk.Align,
			"label":           lnk.Label,
			"label_component": lnk.LabelComponent,
			"icon":            lnk.Icon,
			"disabled":        lnk.Disabled,
			"auto_focus":      lnk.AutoFocus,
			"full":            lnk.Full,
			"small":           lnk.Small,
			"selected":        lnk.Selected,
			"hide_label":      lnk.HideLabel,
			"badge":           lnk.Badge,
			"show_badge":      lnk.ShowBadge,
		})
}

/*
Returns the value of the property of the [Link] with the specified name.
*/
func (lnk *Link) GetProperty(propName string) interface{} {
	return lnk.Properties()[propName]
}

/*
It checks the value given to the property of the [Link] and always returns a valid value
*/
func (lnk *Link) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"link_style": func() interface{} {
			return lnk.CheckEnumValue(ut.ToString(propValue, ""), LinkStyleDefault, LinkStyle)
		},
		"align": func() interface{} {
			return lnk.CheckEnumValue(ut.ToString(propValue, ""), TextAlignCenter, TextAlign)
		},
		"indicator": func() interface{} {
			return lnk.CheckEnumValue(ut.ToString(propValue, IndicatorSpinner), IndicatorSpinner, Indicator)
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if lnk.BaseComponent.GetProperty(propName) != nil {
		return lnk.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [Link] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (lnk *Link) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"href": func() interface{} {
			lnk.Href = ut.ToString(propValue, "")
			return lnk.Href
		},
		"download": func() interface{} {
			lnk.Download = ut.ToString(propValue, "")
			return lnk.Download
		},
		"media": func() interface{} {
			lnk.Media = ut.ToString(propValue, "")
			return lnk.Media
		},
		"ping": func() interface{} {
			lnk.Ping = ut.ToString(propValue, "")
			return lnk.Ping
		},
		"referrerpolicy": func() interface{} {
			lnk.ReferrerPolicy = ut.ToString(propValue, "noreferrer")
			return lnk.ReferrerPolicy
		},
		"rel": func() interface{} {
			lnk.Rel = ut.ToString(propValue, "")
			return lnk.Rel
		},
		"link_target": func() interface{} {
			lnk.LinkTarget = ut.ToString(propValue, "_self")
			return lnk.LinkTarget
		},
		"media_type": func() interface{} {
			lnk.MediaType = ut.ToString(propValue, "")
			return lnk.MediaType
		},
		"link_style": func() interface{} {
			lnk.LinkStyle = lnk.Validation(propName, propValue).(string)
			return lnk.LinkStyle
		},
		"align": func() interface{} {
			lnk.Align = lnk.Validation(propName, propValue).(string)
			return lnk.Align
		},
		"label": func() interface{} {
			lnk.Label = ut.ToString(propValue, "")
			return lnk.Label
		},
		"label_component": func() interface{} {
			if cc, valid := propValue.(ClientComponent); valid {
				lnk.LabelComponent = cc
			}
			return lnk.LabelComponent
		},
		"icon": func() interface{} {
			lnk.Icon = ut.ToString(propValue, "")
			return lnk.Icon
		},
		"disabled": func() interface{} {
			lnk.Disabled = ut.ToBoolean(propValue, false)
			return lnk.Disabled
		},
		"auto_focus": func() interface{} {
			lnk.AutoFocus = ut.ToBoolean(propValue, false)
			return lnk.AutoFocus
		},
		"full": func() interface{} {
			lnk.Full = ut.ToBoolean(propValue, false)
			return lnk.Full
		},
		"small": func() interface{} {
			lnk.Small = ut.ToBoolean(propValue, false)
			return lnk.Small
		},
		"selected": func() interface{} {
			lnk.Selected = ut.ToBoolean(propValue, false)
			return lnk.Selected
		},
		"hide_label": func() interface{} {
			lnk.HideLabel = ut.ToBoolean(propValue, false)
			return lnk.HideLabel
		},
		"badge": func() interface{} {
			lnk.Badge = ut.ToInteger(propValue, 0)
			return lnk.Badge
		},
		"show_badge": func() interface{} {
			lnk.ShowBadge = ut.ToBoolean(propValue, false)
			return lnk.ShowBadge
		},
		"indicator": func() interface{} {
			lnk.Indicator = lnk.Validation(propName, propValue).(string)
			return lnk.Indicator
		},
	}
	if _, found := pm[propName]; found {
		return lnk.SetRequestValue(propName, pm[propName](), []string{})
	}
	if lnk.BaseComponent.GetProperty(propName) != nil {
		return lnk.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (lnk *Link) getComponent(name string) (string, error) {
	ccMap := map[string]func() ClientComponent{
		"icon": func() ClientComponent {
			return &Icon{Value: lnk.Icon, Width: 20}
		},
		"label": func() ClientComponent {
			return lnk.LabelComponent
		},
	}
	return ccMap[name]().Render()
}

/*
Based on the values, it will generate the html code of the [Link] or return with an error message.
*/
func (lnk *Link) Render() (res string, err error) {
	lnk.InitProps(lnk)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(lnk.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(lnk.Class, " ")
		},
		"linkComponent": func(name string) (string, error) {
			return lnk.getComponent(name)
		},
	}
	tpl := `<a id="{{ .Id }}" name="{{ .Name }}" 
	{{ if ne .Href "" }} href="{{ .Href }}" target="{{ .LinkTarget }}" referrerpolicy="{{ .ReferrerPolicy }}"{{ end }}
	{{ if ne .Download "" }} download="{{ .Download }}" type="{{ .MediaType }}"{{ end }}
	{{ if ne .Media "" }} media="{{ .Media }}"{{ end }}
	{{ if ne .Ping "" }} ping="{{ .Ping }}"{{ end }}
	{{ if ne .Rel "" }} rel="{{ .Rel }}"{{ end }}
	{{ if ne .LinkStyle "" }} link-type="{{ .LinkStyle }}"{{ end }}
	{{ if ne .Indicator "none" }} hx-indicator="#{{ .Indicator }}"{{ end }}
	{{ if .Disabled }} disabled{{ end }}
	{{ if .AutoFocus }} autofocus{{ end }}
	{{ if ne .Label "" }} aria-label="{{ .Label }}" title="{{ .Label }}"{{ end }}
	 class="{{ .Align }}{{ if .Small }} small-button{{ end }}{{ if .Full }} full{{ end }}{{ if .Selected }} selected{{ end }}{{ if .HideLabel }} hidelabel{{ end }} {{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	>{{ if and (ne .Icon "") (ne .Align "right") }}{{ linkComponent "icon" }}{{ end }}
	{{ if .LabelComponent }}{{ linkComponent "label" }}{{ else }}<span>{{ .Label }}</span>{{ end }}
	{{ if and (ne .Icon "") (eq .Align "right") }}{{ linkComponent "icon" }}{{ end }}
	{{ if and (ne .LinkStyle "") (.ShowBadge) }}<span class="right" ><span class="badge{{ if .Selected }} selected-badge{{ end }}" >{{ .Badge }}</span></span>{{ end }}
	</a>`

	return ut.TemplateBuilder("link", tpl, funcMap, lnk)
}

// [Link] test and demo data
func TestLink(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeLink,
			Component: &Link{
				BaseComponent: BaseComponent{
					Id: id + "_link_default",
				},
				LinkStyle:  LinkStyleDefault,
				Label:      "Default",
				Href:       "https://www.google.com",
				LinkTarget: "_blank",
			}},
		{
			Label:         "Button",
			ComponentType: ComponentTypeLink,
			Component: &Link{
				BaseComponent: BaseComponent{
					Id: id + "_link_button",
				},
				LinkStyle:  LinkStyleButton,
				Align:      TextAlignCenter,
				Label:      "Button",
				Href:       "https://www.google.com",
				LinkTarget: "_blank",
			}},
		{
			Label:         "Primary and icon",
			ComponentType: ComponentTypeLink,
			Component: &Link{
				BaseComponent: BaseComponent{
					Id: id + "_button_primary",
				},
				LinkStyle:  LinkStylePrimary,
				Align:      TextAlignCenter,
				Label:      "Primary",
				Icon:       "Check",
				Selected:   true,
				Href:       "https://www.google.com",
				LinkTarget: "_blank",
			}},
		{
			Label:         "Border full and badge",
			ComponentType: ComponentTypeLink,
			Component: &Link{
				BaseComponent: BaseComponent{
					Id: id + "_button_full",
				},
				LinkStyle:  LinkStyleBorder,
				Align:      TextAlignLeft,
				Label:      "Border full and badge",
				Full:       true,
				Badge:      0,
				ShowBadge:  true,
				Href:       "https://www.google.com",
				LinkTarget: "_blank",
			}},
		{
			Label:         "Small disabled",
			ComponentType: ComponentTypeLink,
			Component: &Link{
				LinkStyle: LinkStyleDefault,
				Label:     "Small disabled",
				Align:     TextAlignLeft,
				Small:     true,
				Disabled:  true,
			}},
		{
			Label:         "Label component",
			ComponentType: ComponentTypeButton,
			Component: &Link{
				BaseComponent: BaseComponent{
					Id: id + "_button_label",
				},
				LinkStyle:      LinkStyleButton,
				Align:          TextAlignCenter,
				Label:          "Label component",
				LabelComponent: &Icon{Value: "Search", Width: 32, Height: 32},
				Href:           "https://www.google.com",
				LinkTarget:     "_blank",
			}},
	}
}
