package component

import (
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [MenuBar] constants
const (
	ComponentTypeMenuBar = "menubar"

	MenuBarEventSide  = "side"
	MenuBarEventValue = "value"

	SideVisibilityAuto = "auto"
	SideVisibilityShow = "show"
	SideVisibilityHide = "hide"
)

// [MenuBar] SideVisibility values
var SideVisibility []string = []string{SideVisibilityAuto, SideVisibilityShow, SideVisibilityHide}

// [MenuBar] item data
type MenuBarItem struct {
	// Menu key/id value
	Value string `json:"value"`
	// Menu caption
	Label string `json:"label"`
	// Valid [Icon] component value. See more [IconKey] variable values.
	Icon string `json:"icon"`
}

/*
Creates a clickable menu control

For example:

	&MenuBar{
	  BaseComponent: BaseComponent{
	    Id:           "id_menubar_default",
	    EventURL:     "/event",
	    RequestValue: parent_component.GetProperty("request_value").(map[string]ut.IM),
	    RequestMap:   parent_component.GetProperty("request_map").(map[string]ClientComponent),
	  },
	  Items: []MenuBarItem{
	    {Value: "search", Label: "Search", Icon: "Search"},
	    {Value: "edit", Label: "Edit", Icon: "Edit"},
	    {Value: "setting", Label: "Setting", Icon: "Cog"},
	  },
	  Value:          "search",
	  SideBar:        true,
	  SideVisibility: SideVisibilityAuto,
	}
*/
type MenuBar struct {
	BaseComponent
	// Value of a selected/active menu item
	Value string `json:"value"`
	// Show/hide sidebar button in mobile view
	SideBar bool `json:"side_bar"`
	/* [SideVisibility] variable constants:
	[SideVisibilityAuto], [SideVisibilityShow], [SideVisibilityHide].
	Default value: [SideVisibilityAuto] */
	SideVisibility string `json:"side_visibility"`
	// Sidebar button caption. Default: Hide
	LabelHide string `json:"label_hide"`
	// Sidebar button caption. Default: Menu
	LabelMenu string `json:"label_menu"`
	// Menu items data
	Items []MenuBarItem `json:"items"`
}

/*
Returns all properties of the [MenuBar]
*/
func (mnb *MenuBar) Properties() ut.IM {
	return ut.MergeIM(
		mnb.BaseComponent.Properties(),
		ut.IM{
			"value":           mnb.Value,
			"side_bar":        mnb.SideBar,
			"side_visibility": mnb.SideVisibility,
			"label_hide":      mnb.LabelHide,
			"label_menu":      mnb.LabelMenu,
			"items":           mnb.Items,
		})
}

/*
Returns the value of the property of the [MenuBar] with the specified name.
*/
func (mnb *MenuBar) GetProperty(propName string) interface{} {
	return mnb.Properties()[propName]
}

/*
It checks the value given to the property of the [MenuBar] and always returns a valid value
*/
func (mnb *MenuBar) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"side_visibility": func() interface{} {
			return mnb.CheckEnumValue(mnb.SideVisibility, SideVisibilityAuto, SideVisibility)
		},
		"items": func() interface{} {
			items := []MenuBarItem{}
			if it, valid := propValue.([]MenuBarItem); valid && (it != nil) {
				items = it
			}
			return items
		},
		"target": func() interface{} {
			mnb.SetProperty("id", mnb.Id)
			value := ut.ToString(propValue, mnb.Id)
			if value != "this" && !strings.HasPrefix(value, "#") {
				value = "#" + value
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if mnb.BaseComponent.GetProperty(propName) != nil {
		return mnb.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [MenuBar] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (mnb *MenuBar) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			mnb.Value = ut.ToString(propValue, "")
			return mnb.Value
		},
		"side_bar": func() interface{} {
			mnb.SideBar = ut.ToBoolean(propValue, false)
			return mnb.SideBar
		},
		"side_visibility": func() interface{} {
			mnb.SideVisibility = mnb.Validation(propName, propValue).(string)
			return mnb.SideVisibility
		},
		"label_hide": func() interface{} {
			mnb.LabelHide = ut.ToString(propValue, "Hide")
			return mnb.LabelHide
		},
		"label_menu": func() interface{} {
			mnb.LabelMenu = ut.ToString(propValue, "Menu")
			return mnb.LabelMenu
		},
		"items": func() interface{} {
			mnb.Items = mnb.Validation(propName, propValue).([]MenuBarItem)
			return mnb.Items
		},
		"target": func() interface{} {
			mnb.Target = mnb.Validation(propName, propValue).(string)
			return mnb.Target
		},
	}
	if _, found := pm[propName]; found {
		return mnb.SetRequestValue(propName, pm[propName](), []string{})
	}
	if mnb.BaseComponent.GetProperty(propName) != nil {
		return mnb.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (mnb *MenuBar) response(evt ResponseEvent) (re ResponseEvent) {
	mnbEvt := ResponseEvent{Trigger: mnb, TriggerName: mnb.Name}
	switch evt.TriggerName {
	case "item", "icon":
		value := evt.Trigger.GetProperty("data").(ut.IM)["item"].(MenuBarItem).Value
		mnbEvt.Name = MenuBarEventValue
		mnbEvt.Value = value
		if mnb.Value != "" {
			mnb.SetProperty("value", value)
		}
	default:
		mnbEvt.Name = MenuBarEventSide
	}
	if mnb.OnResponse != nil {
		return mnb.OnResponse(mnbEvt)
	}
	return mnbEvt
}

func (mnb *MenuBar) getComponent(name string, item MenuBarItem) (res string, err error) {
	ccMap := map[string]func() ClientComponent{
		"sidebar": func() ClientComponent {
			cclass := []string{"menu-label"}
			label := mnb.LabelMenu
			icon := "Bars"
			is := ut.SM{}
			if mnb.SideVisibility == SideVisibilityShow {
				cclass = []string{"selected exit"}
				label = mnb.LabelHide
				icon = "Close"
				is = ut.SM{"width": "24px", "height": "24px"}
			}
			return &Label{
				BaseComponent: BaseComponent{
					Id:           mnb.Id + "_sidebar",
					Name:         name,
					EventURL:     mnb.EventURL,
					Target:       mnb.Target,
					OnResponse:   mnb.response,
					RequestValue: mnb.RequestValue,
					RequestMap:   mnb.RequestMap,
					Class:        cclass,
				},
				Value:     label,
				LeftIcon:  icon,
				IconStyle: is,
			}
		},
		"item": func() ClientComponent {
			cclass := []string{"menu-label"}
			if item.Value == mnb.Value {
				cclass = []string{"selected"}
			}
			if item.Value == "logout" {
				cclass = append(cclass, "exit")
			}
			return &Label{
				BaseComponent: BaseComponent{
					Id:           mnb.Id + "_" + item.Value,
					Name:         name,
					EventURL:     mnb.EventURL,
					Target:       mnb.Target,
					Data:         ut.IM{"item": item},
					OnResponse:   mnb.response,
					RequestValue: mnb.RequestValue,
					RequestMap:   mnb.RequestMap,
					Class:        cclass,
				},
				Value:    item.Label,
				LeftIcon: item.Icon,
			}
		},
		"icon": func() ClientComponent {
			cclass := []string{"menu-label"}
			if item.Value == mnb.Value {
				cclass = []string{"selected"}
			}
			if item.Value == "logout" {
				cclass = append(cclass, "exit")
			}
			return &Icon{
				BaseComponent: BaseComponent{
					Id:           mnb.Id + "_" + item.Value,
					Name:         name,
					EventURL:     mnb.EventURL,
					Target:       mnb.Target,
					Data:         ut.IM{"item": item},
					OnResponse:   mnb.response,
					RequestValue: mnb.RequestValue,
					RequestMap:   mnb.RequestMap,
					Class:        cclass,
				},
				Value: item.Icon,
			}
		},
	}
	cc := ccMap[name]()
	res, err = cc.Render()
	return res, err
}

/*
Based on the values, it will generate the html code of the [MenuBar] or return with an error message.
*/
func (mnb *MenuBar) Render() (res string, err error) {
	mnb.InitProps(mnb)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(mnb.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(mnb.Class, " ")
		},
		"sideBar": func() (string, error) {
			return mnb.getComponent("sidebar", MenuBarItem{})
		},
		"menuItem": func(item MenuBarItem) (string, error) {
			return mnb.getComponent("item", item)
		},
		"menuIcon": func(item MenuBarItem) (string, error) {
			return mnb.getComponent("icon", item)
		},
		"reverse": func(idx int) MenuBarItem {
			reverseIndex := len(mnb.Items) - 1 - idx
			return mnb.Items[reverseIndex]
		},
	}
	tpl := `<div id="{{ .Id }}" name="{{ .Name }}" class="menubar {{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	><div class="cell">
	{{ if .SideBar }}<div id="mnu_sidebar" class="menuitem sidebar">{{ sideBar }}</div>{{ end }}
	{{ range $index, $item := .Items }}
	<div id="mnu_{{ $item.Value }}_large" class="hide-small hide-medium menuitem">{{ menuItem $item }}</div>
	{{ end }}
	</div>
	<div class="cell container">
	{{ range $index, $item := .Items }}{{ $reverseItem := reverse $index }}
	<div id="mnu_{{ $item.Value }}_medium" class="right hide-large menuitem">
	<span class="hide-small">{{ menuItem $reverseItem }}</span>
	<span class="menu-label hide-medium">{{ menuIcon $reverseItem }}</span>
	</div>
	{{ end }}
	</div>
	</div>`

	return ut.TemplateBuilder("menubar", tpl, funcMap, mnb)
}

// [MenuBar] test and demo data
func TestMenuBar(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeMenuBar,
			Component: &MenuBar{
				BaseComponent: BaseComponent{
					Id:           id + "_menubar_default",
					EventURL:     eventURL,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Items: []MenuBarItem{
					{Value: "search", Label: "Search", Icon: "Search"},
					{Value: "edit", Label: "Edit", Icon: "Edit"},
					{Value: "setting", Label: "Setting", Icon: "Cog"},
					{Value: "bookmark", Label: "Bookmark", Icon: "Star"},
					{Value: "help", Label: "Help", Icon: "QuestionCircle"},
					{Value: "logout", Label: "Logout", Icon: "Exit"},
				},
				Value:          "search",
				SideBar:        true,
				SideVisibility: SideVisibilityShow,
			}},
	}
}
