package molecule

import (
	"strings"

	fm "github.com/nervatura/component/component/atom"
	bc "github.com/nervatura/component/component/base"
)

const (
	MenuBarEventSide  = "side"
	MenuBarEventValue = "value"
)

type MenuBarItem struct {
	Value string
	Label string
	Icon  string
}

type MenuBar struct {
	bc.BaseComponent
	Value          string
	SideBar        bool
	SideVisibility string
	LabelHide      string
	LabelMenu      string
	Items          []MenuBarItem
}

func (mnb *MenuBar) Properties() bc.IM {
	return bc.MergeIM(
		mnb.BaseComponent.Properties(),
		bc.IM{
			"value":           mnb.Value,
			"side_bar":        mnb.SideBar,
			"side_visibility": mnb.SideVisibility,
			"label_hide":      mnb.LabelHide,
			"label_menu":      mnb.LabelMenu,
			"items":           mnb.Items,
		})
}

func (mnb *MenuBar) GetProperty(propName string) interface{} {
	return mnb.Properties()[propName]
}

func (mnb *MenuBar) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"side_visibility": func() interface{} {
			return mnb.CheckEnumValue(mnb.SideVisibility, bc.SideVisibilityAuto, bc.SideVisibility)
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
			value := bc.ToString(propValue, mnb.Id)
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

func (mnb *MenuBar) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"value": func() interface{} {
			mnb.Value = bc.ToString(propValue, "")
			return mnb.Value
		},
		"side_bar": func() interface{} {
			mnb.SideBar = bc.ToBoolean(propValue, false)
			return mnb.SideBar
		},
		"side_visibility": func() interface{} {
			mnb.SideVisibility = mnb.Validation(propName, propValue).(string)
			return mnb.SideVisibility
		},
		"label_hide": func() interface{} {
			mnb.LabelHide = bc.ToString(propValue, "Hide")
			return mnb.LabelHide
		},
		"label_menu": func() interface{} {
			mnb.LabelMenu = bc.ToString(propValue, "Menu")
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
		return pm[propName]()
	}
	if mnb.BaseComponent.GetProperty(propName) != nil {
		return mnb.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (mnb *MenuBar) response(evt bc.ResponseEvent) (re bc.ResponseEvent) {
	mnbEvt := bc.ResponseEvent{Trigger: mnb, TriggerName: mnb.Name}
	switch evt.TriggerName {
	case "item", "icon":
		value := evt.Trigger.GetProperty("data").(bc.IM)["item"].(MenuBarItem).Value
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
	ccMap := map[string]func() bc.ClientComponent{
		"sidebar": func() bc.ClientComponent {
			cclass := []string{"menu-label"}
			label := mnb.LabelMenu
			icon := "Bars"
			is := bc.SM{}
			if mnb.SideVisibility == bc.SideVisibilityShow {
				cclass = []string{"selected exit"}
				label = mnb.LabelHide
				icon = "Close"
				is = bc.SM{"width": "24px", "height": "24px"}
			}
			return &fm.Label{
				BaseComponent: bc.BaseComponent{
					Id:         mnb.Id + "_sidebar",
					Name:       name,
					EventURL:   mnb.EventURL,
					Target:     mnb.Target,
					OnResponse: mnb.response,
					Class:      cclass,
				},
				Value:     label,
				LeftIcon:  icon,
				IconStyle: is,
			}
		},
		"item": func() bc.ClientComponent {
			cclass := []string{"menu-label"}
			if item.Value == mnb.Value {
				cclass = []string{"selected"}
			}
			if item.Value == "logout" {
				cclass = append(cclass, "exit")
			}
			return &fm.Label{
				BaseComponent: bc.BaseComponent{
					Id:         mnb.Id + "_" + item.Value,
					Name:       name,
					EventURL:   mnb.EventURL,
					Target:     mnb.Target,
					Data:       bc.IM{"item": item},
					OnResponse: mnb.response,
					Class:      cclass,
				},
				Value:    item.Label,
				LeftIcon: item.Icon,
			}
		},
		"icon": func() bc.ClientComponent {
			cclass := []string{"menu-label"}
			if item.Value == mnb.Value {
				cclass = []string{"selected"}
			}
			if item.Value == "logout" {
				cclass = append(cclass, "exit")
			}
			return &fm.Icon{
				BaseComponent: bc.BaseComponent{
					Id:         mnb.Id + "_" + item.Value,
					Name:       name,
					EventURL:   mnb.EventURL,
					Target:     mnb.Target,
					Data:       bc.IM{"item": item},
					OnResponse: mnb.response,
					Class:      cclass,
				},
				Value: item.Icon,
			}
		},
	}
	cc := ccMap[name]()
	res, err = cc.Render()
	if err == nil {
		mnb.RequestMap = bc.MergeCM(mnb.RequestMap, cc.GetProperty("request_map").(map[string]bc.ClientComponent))
	}
	return res, err
}

func (mnb *MenuBar) InitProps() {
	for key, value := range mnb.Properties() {
		mnb.SetProperty(key, value)
	}
}

func (mnb *MenuBar) Render() (res string, err error) {
	mnb.InitProps()

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

	return bc.TemplateBuilder("menubar", tpl, funcMap, mnb)
}

func DemoMenuBar(eventURL, parentID string) []bc.DemoComponent {
	return []bc.DemoComponent{
		{
			Label:         "Default",
			ComponentType: bc.ComponentTypeMenuBar,
			Component: &MenuBar{
				BaseComponent: bc.BaseComponent{
					Id:       bc.GetComponentID(),
					EventURL: eventURL,
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
				SideVisibility: bc.SideVisibilityShow,
			}},
	}
}
