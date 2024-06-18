package component

import (
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [Field] constants
const (
	ComponentTypeSideBar = "sidebar"

	SideBarItemTypeState     = "state"
	SideBarItemTypeGroup     = "group"
	SideBarItemTypeElement   = "element"
	SideBarItemTypeStatic    = "static"
	SideBarItemTypeSeparator = "separator"

	SideBarVisibilityAuto = "auto"
	SideBarVisibilityShow = "show"
	SideBarVisibilityHide = "hide"

	SideBarEventItem  = "sidebar_item"
	SideBarEventGroup = "sidebar_group"
	SideBarEventState = "sidebar_state"
)

// [SideBar] Visibility values
var SideBarVisibility []string = []string{SideBarVisibilityAuto, SideBarVisibilityShow, SideBarVisibilityHide}

type SideBarItem interface {
	ItemType() string
	GetValue() string
	GetSelected() bool
}

type SideBarState struct {
	Name          string           `json:"name"`
	SelectedIndex int              `json:"selected_index"`
	Items         []SideBarElement `json:"items"`
}

func (sbst *SideBarState) ItemType() string {
	return SideBarItemTypeState
}

func (sbst *SideBarState) GetValue() string {
	if len(sbst.Items) > sbst.SelectedIndex {
		return sbst.Items[sbst.SelectedIndex].Value
	}
	return ""
}

func (sbst *SideBarState) GetSelected() bool {
	return false
}

type SideBarGroup struct {
	Name     string           `json:"name"`
	Value    string           `json:"value"`
	Label    string           `json:"label"`
	Align    string           `json:"align"`
	Icon     string           `json:"icon"`
	Selected bool             `json:"selected"`
	Items    []SideBarElement `json:"items"`
	Disabled bool             `json:"disabled"`
}

func (sbgr *SideBarGroup) ItemType() string {
	return SideBarItemTypeGroup
}

func (sbgr *SideBarGroup) GetValue() string {
	return sbgr.Value
}

func (sbgr *SideBarGroup) GetSelected() bool {
	return sbgr.Selected
}

type SideBarElement struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Label    string `json:"label"`
	Align    string `json:"align"`
	Icon     string `json:"icon"`
	Selected bool   `json:"selected"`
	Disabled bool   `json:"disabled"`
	NotFull  bool   `json:"not_full"`
}

func (sbe *SideBarElement) ItemType() string {
	return SideBarItemTypeElement
}

func (sbe *SideBarElement) GetValue() string {
	return sbe.Value
}

func (sbe *SideBarElement) GetSelected() bool {
	return sbe.Selected
}

type SideBarStatic struct {
	Label string `json:"label"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

func (sbs *SideBarStatic) ItemType() string {
	return SideBarItemTypeStatic
}

func (sbs *SideBarStatic) GetValue() string {
	return ""
}

func (sbs *SideBarStatic) GetSelected() bool {
	return false
}

type SideBarSeparator struct{}

func (sbsp *SideBarSeparator) ItemType() string {
	return SideBarItemTypeSeparator
}

func (sbsp *SideBarSeparator) GetValue() string {
	return ""
}

func (sbsp *SideBarSeparator) GetSelected() bool {
	return false
}

/*
Creates a clickable side menu control
*/
type SideBar struct {
	BaseComponent
	Items      []SideBarItem `json:"items"`
	Visibility string        `json:"visibility"`
}

/*
Returns all properties of the [SideBar]
*/
func (sb *SideBar) Properties() ut.IM {
	return ut.MergeIM(
		sb.BaseComponent.Properties(),
		ut.IM{
			"items":      sb.Items,
			"visibility": sb.Visibility,
		})
}

/*
Returns the value of the property of the [SideBar] with the specified name.
*/
func (sb *SideBar) GetProperty(propName string) interface{} {
	return sb.Properties()[propName]
}

/*
It checks the value given to the property of the [SideBar] and always returns a valid value
*/
func (sb *SideBar) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"items": func() interface{} {
			if value, valid := propValue.([]SideBarItem); valid {
				return value
			}
			return []SideBarItem{}
		},
		"visibility": func() interface{} {
			return sb.CheckEnumValue(ut.ToString(propValue, ""), SideBarVisibilityAuto, SideBarVisibility)
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if sb.BaseComponent.GetProperty(propName) != nil {
		return sb.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [SideBar] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (sb *SideBar) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"items": func() interface{} {
			sb.Items = sb.Validation(propName, propValue).([]SideBarItem)
			return sb.Items
		},
		"visibility": func() interface{} {
			sb.Visibility = sb.Validation(propName, propValue).(string)
			return sb.Visibility
		},
	}
	if _, found := pm[propName]; found {
		return sb.SetRequestValue(propName, pm[propName](), []string{})
	}
	if sb.BaseComponent.GetProperty(propName) != nil {
		return sb.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (sb *SideBar) response(evt ResponseEvent) (re ResponseEvent) {
	sbEvt := ResponseEvent{Trigger: sb, TriggerName: sb.Name, Name: SideBarEventItem,
		Header: ut.SM{HeaderRetarget: "#" + sb.Id}}
	itIndex := ut.ToInteger(evt.Trigger.GetProperty("data").(ut.IM)["index"], 0)

	switch sb.Items[itIndex].ItemType() {
	case SideBarItemTypeGroup:
		groupIndex := ut.ToInteger(evt.Trigger.GetProperty("data").(ut.IM)["group_index"], -1)
		gr := sb.Items[itIndex].(*SideBarGroup)
		if groupIndex >= 0 {
			sbEvt.Value = gr.Items[groupIndex].Value
		} else {
			sbEvt.Name = SideBarEventGroup
			sbEvt.Value = gr.GetValue()
			gr.Selected = !gr.GetSelected()
		}

	case SideBarItemTypeState:
		groupIndex := ut.ToInteger(evt.Trigger.GetProperty("data").(ut.IM)["group_index"], -1)
		st := sb.Items[itIndex].(*SideBarState)
		sbEvt.Name = SideBarEventState
		sbEvt.Value = st.GetValue()
		st.SelectedIndex = int(groupIndex)

	default:
		sbEvt.Value = sb.Items[itIndex].GetValue()
	}
	if sb.OnResponse != nil {
		return sb.OnResponse(sbEvt)
	}
	return sbEvt
}

func (sb *SideBar) getComponent(index, groupIndex int) (res string, err error) {
	ccBtn := func(idx int, name string) *Button {
		btn := &Button{
			BaseComponent: BaseComponent{
				Id:   sb.Id + "_" + ut.ToString(idx, ""),
				Name: name,
				Data: ut.IM{"index": idx},
				Style: ut.SM{
					"border-radius": "0",
					"border-color":  "rgba(var(--accent-1c), 0.2)"},
				EventURL:     sb.EventURL,
				Target:       sb.Target,
				OnResponse:   sb.response,
				RequestValue: sb.RequestValue,
				RequestMap:   sb.RequestMap,
			},
			ButtonStyle: ButtonStylePrimary,
		}
		return btn
	}
	ccMap := map[string]func(it interface{}) ClientComponent{
		SideBarItemTypeState: func(it interface{}) ClientComponent {
			st := it.(*SideBarState)
			el := it.(*SideBarState).Items[groupIndex]
			btn := ccBtn(index, el.Name)
			btn.SetProperty("id", btn.Id+"_"+ut.ToString(groupIndex, ""))
			btn.SetProperty("icon", el.Icon)
			btn.SetProperty("label", el.Label)
			btn.SetProperty("align", ut.ToString(el.Align, TextAlignCenter))
			btn.SetProperty("selected", (st.SelectedIndex == groupIndex))
			btn.SetProperty("full", true)
			btn.SetProperty("disabled", el.Disabled)
			btn.SetProperty("data", ut.MergeIM(btn.Data, ut.IM{
				"group_index": groupIndex,
			}))
			return btn
		},
		SideBarItemTypeGroup: func(it interface{}) ClientComponent {
			var btn *Button
			if groupIndex >= 0 {
				el := it.(*SideBarGroup).Items[groupIndex]
				btn = ccBtn(index, el.Name)
				btn.SetProperty("id", btn.Id+"_"+ut.ToString(groupIndex, ""))
				btn.SetProperty("icon", el.Icon)
				btn.SetProperty("label", el.Label)
				btn.SetProperty("align", ut.ToString(el.Align, TextAlignLeft))
				btn.SetProperty("selected", false)
				btn.SetProperty("full", true)
				btn.SetProperty("disabled", el.Disabled)
				btn.SetProperty("data", ut.MergeIM(btn.Data, ut.IM{
					"group_index": groupIndex,
				}))
				btn.SetProperty("style", ut.MergeSM(btn.Style, ut.SM{
					"color": "rgb(var(--functional-blue))",
					"fill":  "rgb(var(--functional-blue))",
				}))
			} else {
				gr := it.(*SideBarGroup)
				btn = ccBtn(index, gr.Name)
				btn.SetProperty("icon", gr.Icon)
				btn.SetProperty("label", gr.Label)
				btn.SetProperty("align", ut.ToString(gr.Align, TextAlignLeft))
				btn.SetProperty("full", true)
				btn.SetProperty("selected", gr.Selected)
				btn.SetProperty("disabled", gr.Disabled)
			}
			return btn
		},
		SideBarItemTypeElement: func(it interface{}) ClientComponent {
			el := it.(*SideBarElement)
			btn := ccBtn(index, el.Name)
			btn.SetProperty("icon", el.Icon)
			btn.SetProperty("label", el.Label)
			btn.SetProperty("align", ut.ToString(el.Align, TextAlignLeft))
			btn.SetProperty("full", !el.NotFull)
			btn.SetProperty("selected", el.Selected)
			btn.SetProperty("disabled", el.Disabled)
			return btn
		},
		SideBarItemTypeStatic: func(it interface{}) ClientComponent {
			lbl := &Label{
				Value:    it.(*SideBarStatic).Label,
				LeftIcon: it.(*SideBarStatic).Icon,
			}
			if it.(*SideBarStatic).Color != "" {
				lbl.Style = ut.SM{"color": it.(*SideBarStatic).Color}
				lbl.IconStyle = ut.SM{"fill": it.(*SideBarStatic).Color}
			}
			return lbl
		},
	}
	cc := ccMap[sb.Items[index].ItemType()](sb.Items[index])
	res, err = cc.Render()
	return res, err
}

/*
Based on the values, it will generate the html code of the [SideBar] or return with an error message.
*/
func (sb *SideBar) Render() (res string, err error) {
	sb.InitProps(sb)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(sb.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(sb.Class, " ")
		},
		"sidebarType": func(index int) string {
			return sb.Items[index].ItemType()
		},
		"selectedComponent": func(index int) bool {
			return sb.Items[index].GetSelected()
		},
		"validState": func(index int) bool {
			return index <= 1
		},
		"sidebarComponent": func(index, groupIndex int) (string, error) {
			return sb.getComponent(index, groupIndex)
		},
	}
	tpl := `<div id="{{ .Id }}" name="{{ .Name }}" 
	class="sidebar {{ customClass }}{{ if ne .Visibility "auto" }} {{ .Visibility  }}{{ end }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	>{{ range $index, $item := .Items }}
	{{ $stype := sidebarType $index  }}
	{{ if eq $stype "separator" }}<hr id="separator_{{ $index }}" class="separator" />{{ end }}
	{{ if eq $stype "static" }}<div class="row full"><div id="static_{{ $index }}" class="static-label" >{{ sidebarComponent $index -1 }}</div></div>{{ end }}
	{{ if eq $stype "element" }}{{ sidebarComponent $index -1 }}{{ end }}
	{{ if eq $stype "group" }}<div class="row full">{{ sidebarComponent $index -1 }}</div>
	{{ if selectedComponent $index }}<div class="row full sidebar-group" >
	{{ range $groupIndex, $groupItem := $item.Items }}{{ sidebarComponent $index $groupIndex }}{{ end }}
	</div>{{ end }}
	{{ end }}
	{{ if eq $stype "state" }}<div class="row full container">
	{{ range $groupIndex, $groupItem := $item.Items }}{{ if validState $groupIndex }}
	<div class="cell half">{{ sidebarComponent $index $groupIndex }}</div>
	{{ end }}{{ end }}
	</div>{{ end }}
	{{ end }}
	</div>`

	return ut.TemplateBuilder("sidebar", tpl, funcMap, sb)
}

var testSidebarResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
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
	if evt.Name == SideBarEventItem {
		return toast(ut.ToString(evt.Value, ""))
	}
	return evt
}

// [Sidebar] test and demo data
func TestSidebar(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeList,
			Component: &SideBar{
				BaseComponent: BaseComponent{
					Id:           id + "_default",
					EventURL:     eventURL,
					OnResponse:   testSidebarResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				Items: []SideBarItem{
					&SideBarState{
						Name:          "menu_state",
						SelectedIndex: 0,
						Items: []SideBarElement{
							{
								Name:  "state_new",
								Value: "state_new",
								Label: "New",
								Icon:  "Plus",
							},
							{
								Name:  "state_edit",
								Value: "state_edit",
								Label: "Edit",
								Icon:  "Edit",
							},
						},
					},
					&SideBarSeparator{},
					&SideBarStatic{
						Icon: "ExclamationTriangle", Label: "Color Label", Color: "red",
					},
					&SideBarElement{
						Name:  "left_element",
						Value: "copy",
						Label: "Menu item",
						Align: TextAlignLeft,
						Icon:  "Copy",
					},
					&SideBarElement{
						Name:  "right_element",
						Value: "right",
						Label: "Right item",
						Align: TextAlignRight,
						Icon:  "ArrowRight",
					},
					&SideBarElement{
						Name:     "selected_element",
						Value:    "save",
						Label:    "Menu item",
						Icon:     "Check",
						Selected: true,
					},
					&SideBarSeparator{},
					&SideBarElement{
						Name:    "back_element",
						Value:   "back",
						Label:   "Back",
						Icon:    "Reply",
						NotFull: true,
					},
					&SideBarSeparator{},
					&SideBarElement{
						Name:     "disabled_element",
						Label:    "Disabled item",
						Icon:     "Times",
						Align:    TextAlignCenter,
						Disabled: true,
					},
					&SideBarSeparator{},
					&SideBarStatic{
						Icon: "ExclamationTriangle", Label: "Static Label",
					},
					&SideBarSeparator{},
					&SideBarGroup{
						Name:     "search_1",
						Value:    "search_1",
						Label:    "Search 1",
						Icon:     "FileText",
						Selected: true,
						Items: []SideBarElement{
							{
								Name:  "quick_search",
								Value: "quick_search",
								Label: "Quick Search",
								Icon:  "Bolt",
							},
							{
								Name:  "browser",
								Value: "browser",
								Label: "Browser Search",
								Icon:  "Search",
							},
						},
					},
					&SideBarGroup{
						Name:     "search_2",
						Value:    "search_2",
						Label:    "Search 2",
						Icon:     "FileText",
						Selected: false,
						Items: []SideBarElement{
							{
								Name:  "quick_search",
								Value: "quick_search",
								Label: "Quick Search",
								Icon:  "Bolt",
							},
							{
								Name:  "browser",
								Value: "browser",
								Label: "Browser Search",
								Icon:  "Search",
							},
						},
					},
				},
				Visibility: SideBarVisibilityAuto,
			}},
	}
}
