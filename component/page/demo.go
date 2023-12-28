package page

import (
	"strings"

	fm "github.com/nervatura/component/component/atom"
	bc "github.com/nervatura/component/component/base"
	md "github.com/nervatura/component/component/modal"
	mc "github.com/nervatura/component/component/molecule"
)

const (
	ComponentGroupAtom     = "atom"
	ComponentGroupMolecule = "molecule"
	ComponentGroupModal    = "modal"
	ComponentGroupPage     = "page"

	ViewSizeCentered = "centered"
	ViewSizeFull     = "full"
)

type Demo struct {
	bc.BaseComponent
	Title         string
	Theme         string
	ViewSize      string
	SelectedGroup string
	SelectedType  int64
	SelectedDemo  int64
	demoMap       map[string][]DemoView
}

type DemoSession struct {
	Label     string
	Component bc.ClientComponent
}

type DemoView struct {
	ComponentType string
	Stories       func(eventURL, parentID string) []bc.DemoComponent
	Session       []DemoSession
}

var ComponentGroup []string = []string{
	ComponentGroupAtom, ComponentGroupMolecule, ComponentGroupModal, ComponentGroupPage,
}
var ViewSize []string = []string{ViewSizeCentered, ViewSizeFull}

var DemoMap map[string][]DemoView = map[string][]DemoView{
	ComponentGroupAtom: {
		{ComponentType: bc.ComponentTypeButton, Stories: fm.DemoButton},
		{ComponentType: bc.ComponentTypeDateTime, Stories: fm.DemoDateTime},
		{ComponentType: bc.ComponentTypeIcon, Stories: fm.DemoIcon},
		{ComponentType: bc.ComponentTypeInput, Stories: fm.DemoInput},
		{ComponentType: bc.ComponentTypeLabel, Stories: fm.DemoLabel},
		{ComponentType: bc.ComponentTypeNumberInput, Stories: fm.DemoNumberInput},
		{ComponentType: bc.ComponentTypeSelect, Stories: fm.DemoSelect},
		{ComponentType: bc.ComponentTypeToast, Stories: fm.DemoToast},
	},
	ComponentGroupMolecule: {
		{ComponentType: bc.ComponentTypeTable, Stories: mc.DemoTable},
		{ComponentType: bc.ComponentTypeMenuBar, Stories: mc.DemoMenuBar},
		{ComponentType: bc.ComponentTypePagination, Stories: mc.DemoPagination},
	},
	ComponentGroupModal: {
		{ComponentType: bc.ComponentTypeLogin, Stories: md.DemoLogin},
	},
	ComponentGroupPage: {
		{ComponentType: bc.ComponentTypeAdmin, Stories: DemoAdmin},
	},
}

func NewDemo(eventURL, title string) *Demo {
	sto := &Demo{
		BaseComponent: bc.BaseComponent{
			Id:       bc.GetComponentID(),
			EventURL: eventURL,
		},
		Title:   title,
		demoMap: DemoMap,
	}
	sto.initDemoMap()
	return sto
}

func (sto *Demo) initDemoMap() {
	for group, sg := range sto.demoMap {
		for index, sv := range sg {
			sto.demoMap[group][index].Session = make([]DemoSession, 0)
			for _, sc := range sv.Stories(sto.EventURL, sto.Id) {
				sto.demoMap[group][index].Session = append(
					sto.demoMap[group][index].Session, DemoSession{Label: sc.Label, Component: sc.Component})
			}
		}
	}
}

func (sto *Demo) Properties() bc.IM {
	return bc.MergeIM(
		sto.BaseComponent.Properties(),
		bc.IM{
			"title":          sto.Title,
			"theme":          sto.Theme,
			"view_size":      sto.ViewSize,
			"selected_group": sto.SelectedGroup,
			"selected_type":  sto.SelectedType,
			"selected_demo":  sto.SelectedDemo,
		})
}

func (sto *Demo) GetProperty(propName string) interface{} {
	return sto.Properties()[propName]
}

func (sto *Demo) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"theme": func() interface{} {
			return sto.CheckEnumValue(bc.ToString(propValue, ""), bc.ThemeLight, bc.Theme)
		},
		"view_size": func() interface{} {
			return sto.CheckEnumValue(bc.ToString(propValue, ""), ViewSizeCentered, ViewSize)
		},
		"selected_group": func() interface{} {
			return sto.CheckEnumValue(bc.ToString(propValue, ""), ComponentGroupAtom, ComponentGroup)
		},
		"selected_type": func() interface{} {
			value := bc.ToInteger(propValue, 0)
			if value > int64(len(sto.demoMap[sto.SelectedGroup])-1) {
				value = 0
			}
			return value
		},
		"selected_demo": func() interface{} {
			value := bc.ToInteger(propValue, 0)
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if sto.BaseComponent.GetProperty(propName) != nil {
		return sto.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

func (sto *Demo) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"title": func() interface{} {
			sto.Title = bc.ToString(propValue, "")
			return sto.Title
		},
		"theme": func() interface{} {
			sto.Theme = sto.Validation(propName, propValue).(string)
			return sto.Theme
		},
		"view_size": func() interface{} {
			sto.ViewSize = sto.Validation(propName, propValue).(string)
			return sto.ViewSize
		},
		"selected_group": func() interface{} {
			value := sto.Validation(propName, propValue).(string)
			if sto.SelectedGroup != value {
				sto.SelectedGroup = value
				sto.SelectedType = 0
				sto.SelectedDemo = 0
			}
			return sto.SelectedGroup
		},
		"selected_type": func() interface{} {
			value := sto.Validation(propName, propValue).(int64)
			if sto.SelectedType != value {
				sto.SelectedType = value
				sto.SelectedDemo = 0
			}
			return sto.SelectedType
		},
		"selected_demo": func() interface{} {
			value := sto.Validation(propName, propValue).(int64)
			if sto.SelectedDemo != value {
				sto.SelectedDemo = value
			}
			return sto.SelectedDemo
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if sto.BaseComponent.GetProperty(propName) != nil {
		return sto.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (sto *Demo) response(evt bc.ResponseEvent) (re bc.ResponseEvent) {
	value := sto.SetProperty(evt.TriggerName, evt.Value)
	return bc.ResponseEvent{
		Trigger: sto, TriggerName: sto.Name, Name: evt.TriggerName, Value: value,
	}
}

func (sto *Demo) getComponent(name string) (res string, err error) {
	btnMap := map[string][]string{
		bc.ThemeDark: {bc.ThemeLight, "Sun"}, bc.ThemeLight: {bc.ThemeDark, "Moon"},
		ViewSizeCentered: {ViewSizeFull, "Desktop"}, ViewSizeFull: {ViewSizeCentered, "Mobile"},
	}
	propValue := bc.SM{
		"theme": sto.Theme, "view_size": sto.ViewSize,
		"selected_group": sto.SelectedGroup, "selected_type": bc.ToString(sto.SelectedType, ""),
		"selected_demo": bc.ToString(sto.SelectedDemo, ""),
	}
	ccBtn := func() *fm.Button {
		return &fm.Button{
			BaseComponent: bc.BaseComponent{
				Id: sto.Id + "_" + name, Name: name,
				Style:      bc.SM{"padding": "8px"},
				EventURL:   sto.EventURL,
				Target:     sto.Id,
				OnResponse: sto.response,
			},
			Type:           fm.ButtonTypePrimary,
			Value:          btnMap[propValue[name]][0],
			LabelComponent: &fm.Icon{Value: btnMap[propValue[name]][1], Width: 18, Height: 18},
		}
	}
	ccSel := func() *fm.Select {
		return &fm.Select{
			BaseComponent: bc.BaseComponent{
				Id: sto.Id + "_" + name, Name: name,
				EventURL:   sto.EventURL,
				Target:     sto.Id,
				OnResponse: sto.response,
			},
			Value:   propValue[name],
			Options: []fm.SelectOption{},
		}
	}
	ccMap := map[string]func() bc.ClientComponent{
		"theme":     func() bc.ClientComponent { return ccBtn() },
		"view_size": func() bc.ClientComponent { return ccBtn() },
		"selected_group": func() bc.ClientComponent {
			sc := ccSel()
			for _, group := range ComponentGroup {
				sc.Options = append(
					sc.Options, fm.SelectOption{Value: group, Text: group})
			}
			return sc
		},
		"selected_type": func() bc.ClientComponent {
			sc := ccSel()
			for index, v := range sto.demoMap[sto.SelectedGroup] {
				sc.Options = append(
					sc.Options, fm.SelectOption{Value: bc.ToString(index, ""), Text: v.ComponentType})
			}
			return sc
		},
		"selected_demo": func() bc.ClientComponent {
			sc := ccSel()
			for index, v := range sto.demoMap[sto.SelectedGroup][sto.SelectedType].Session {
				sc.Options = append(
					sc.Options, fm.SelectOption{Value: bc.ToString(index, ""), Text: v.Label})
			}
			return sc
		},
	}
	cc := ccMap[name]()
	res, err = cc.Render()
	if err == nil {
		sto.RequestMap = bc.MergeCM(sto.RequestMap, cc.GetProperty("request_map").(map[string]bc.ClientComponent))
	}
	return res, err
}

func (sto *Demo) InitProps() {
	for key, value := range sto.Properties() {
		sto.SetProperty(key, value)
	}
}

func (sto *Demo) Render() (res string, err error) {
	sto.InitProps()

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(sto.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(sto.Class, " ")
		},
		"demoComponent": func(name string) (string, error) {
			return sto.getComponent(name)
		},
		"label": func(value string) (string, error) {
			return (&fm.Label{
				BaseComponent: bc.BaseComponent{Style: bc.SM{"color": "brown"}},
				Value:         value,
			}).Render()
		},
		"clientComponent": func(cc bc.ClientComponent) (string, error) {
			res, err := cc.Render()
			if err == nil {
				sto.RequestMap = bc.MergeCM(sto.RequestMap, cc.GetProperty("request_map").(map[string]bc.ClientComponent))
			}
			return res, err
		},
		"stories": func() []DemoSession {
			return sto.demoMap[sto.SelectedGroup][sto.SelectedType].Session
		},
		"demo": func() DemoSession {
			return sto.demoMap[sto.SelectedGroup][sto.SelectedType].Session[sto.SelectedDemo]
		},
	}
	tpl := `<div id="{{ .Id }}" theme="{{ .Theme }}" class="demo row mobile {{ .ViewSize }} {{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}>
	<div class="menubar">
	<div class="section-small">
	<div class="cell menu-label padding-small bold">{{ .Title }}</div>
	<div class="cell mobile">
	<div class="cell padding-tiny">{{ demoComponent "theme" }}</div>
	<div class="cell padding-tiny">{{ demoComponent "view_size" }}</div>
	<div class="cell padding-tiny">{{ demoComponent "selected_group" }}</div>
	<div class="cell padding-tiny">{{ demoComponent "selected_type" }}</div>
	</div>
	{{ if ne $.SelectedGroup "atom" }}
	<div class="cell padding-tiny mobile">{{ demoComponent "selected_demo" }}</div>
	{{ end }}
	</div></div>
	<div class="row full section">
	{{ if eq $.SelectedGroup "atom" }}
	{{ range $index, $se := stories }}
	<div class="row full"><div class="cell bold italic padding-normal">{{ label $se.Label }}</div></div>
	<div class="row full"><div class="cell padding-normal">{{ clientComponent $se.Component }}</div></div>
	{{ end }}
	{{ else }}
	{{ $st := demo }}
	<div class="row full"><div class="cell bold italic padding-normal">{{ label $st.Label }}</div></div>
	<div class="row full"><div class="cell padding-normal">{{ clientComponent $st.Component }}</div></div>
	{{ end }}
	</div></div>`

	return bc.TemplateBuilder("demo", tpl, funcMap, sto)
}
