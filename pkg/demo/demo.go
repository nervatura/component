package demo

import (
	"fmt"
	"strings"

	ct "github.com/nervatura/component/pkg/component"
	ut "github.com/nervatura/component/pkg/util"
)

const (
	DemoEventChange   = "change"
	DemoEventTheme    = "theme"
	DemoEventViewSize = "view_size"

	ComponentGroupAtom     = "atom"
	ComponentGroupMolecule = "molecule"
	ComponentGroupModal    = "modal"
	ComponentGroupTemplate = "template"
	ComponentGroupPage     = "page"

	ViewSizeCentered = "centered"
	ViewSizeFull     = "full"
)

type Demo struct {
	ct.BaseComponent
	Title         string                `json:"title"`
	Theme         string                `json:"theme"`
	ViewSize      string                `json:"view_size"`
	SelectedGroup string                `json:"selected_group"`
	SelectedType  int64                 `json:"selected_type"`
	SelectedDemo  int64                 `json:"selected_demo"`
	DemoMap       map[string][]DemoView `json:"-"`
}

type DemoSession struct {
	Label     string             `json:"label"`
	Component ct.ClientComponent `json:"component"`
}

type DemoView struct {
	ComponentType string                                           `json:"component_type"`
	Session       []DemoSession                                    `json:"session"`
	Stories       func(demo ct.ClientComponent) []ct.TestComponent `json:"-"`
}

var ComponentGroup []string = []string{
	ComponentGroupAtom, ComponentGroupMolecule, ComponentGroupModal, ComponentGroupTemplate,
}
var ViewSize []string = []string{ViewSizeCentered, ViewSizeFull}

var DemoMap map[string][]DemoView = map[string][]DemoView{
	ComponentGroupAtom: {
		{ComponentType: ct.ComponentTypeButton, Stories: ct.TestButton},
		{ComponentType: ct.ComponentTypeDateTime, Stories: ct.TestDateTime},
		{ComponentType: ct.ComponentTypeIcon, Stories: ct.TestIcon},
		{ComponentType: ct.ComponentTypeInput, Stories: ct.TestInput},
		{ComponentType: ct.ComponentTypeLabel, Stories: ct.TestLabel},
		{ComponentType: ct.ComponentTypeNumberInput, Stories: ct.TestNumberInput},
		{ComponentType: ct.ComponentTypeSelect, Stories: ct.TestSelect},
		{ComponentType: ct.ComponentTypeToast, Stories: ct.TestToast},
	},
	ComponentGroupMolecule: {
		{ComponentType: ct.ComponentTypeTable, Stories: ct.TestTable},
		{ComponentType: ct.ComponentTypeMenuBar, Stories: ct.TestMenuBar},
		{ComponentType: ct.ComponentTypePagination, Stories: ct.TestPagination},
	},
	ComponentGroupModal: {
		{ComponentType: ct.ComponentTypeLogin, Stories: ct.TestLogin},
	},
	ComponentGroupTemplate: {
		{ComponentType: ct.ComponentTypeLocale, Stories: ct.TestLocale},
	},
}

var demoIcoMap map[string][]string = map[string][]string{
	ct.ThemeDark: {ct.ThemeLight, "Sun"}, ct.ThemeLight: {ct.ThemeDark, "Moon"},
	ViewSizeCentered: {ViewSizeFull, "Desktop"}, ViewSizeFull: {ViewSizeCentered, "Mobile"},
}

func NewDemo(eventURL, title string) *Demo {
	sto := &Demo{
		BaseComponent: ct.BaseComponent{
			Id:           ut.GetComponentID(),
			EventURL:     eventURL,
			RequestValue: map[string]ut.IM{},
			RequestMap:   map[string]ct.ClientComponent{},
		},
		Title:   title,
		DemoMap: DemoMap,
	}
	sto.InitDemoMap()
	return sto
}

func (sto *Demo) InitDemoMap() {
	for group, sg := range sto.DemoMap {
		for index, sv := range sg {
			sto.DemoMap[group][index].Session = make([]DemoSession, 0)
			for _, sc := range sv.Stories(sto) {
				sto.DemoMap[group][index].Session = append(
					sto.DemoMap[group][index].Session, DemoSession{Label: sc.Label, Component: sc.Component})
			}
		}
	}
}

func (sto *Demo) Properties() ut.IM {
	return ut.MergeIM(
		sto.BaseComponent.Properties(),
		ut.IM{
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
			return sto.CheckEnumValue(ut.ToString(propValue, ""), ct.ThemeLight, ct.Theme)
		},
		"view_size": func() interface{} {
			return sto.CheckEnumValue(ut.ToString(propValue, ""), ViewSizeCentered, ViewSize)
		},
		"selected_group": func() interface{} {
			return sto.CheckEnumValue(ut.ToString(propValue, ""), ComponentGroupAtom, ComponentGroup)
		},
		"selected_type": func() interface{} {
			value := ut.ToInteger(propValue, 0)
			if value > int64(len(sto.DemoMap[sto.SelectedGroup])-1) {
				value = 0
			}
			return value
		},
		"selected_demo": func() interface{} {
			value := ut.ToInteger(propValue, 0)
			if len(sto.DemoMap[sto.SelectedGroup]) == 0 || value > int64(len(sto.DemoMap[sto.SelectedGroup][sto.SelectedType].Session)-1) {
				value = 0
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return sto.SetRequestValue(propName, pm[propName](), []string{})
	}
	if sto.BaseComponent.GetProperty(propName) != nil {
		return sto.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

func (sto *Demo) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"title": func() interface{} {
			sto.Title = ut.ToString(propValue, "")
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
		return sto.SetRequestValue(propName, pm[propName](), []string{})
	}
	if sto.BaseComponent.GetProperty(propName) != nil {
		return sto.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (sto *Demo) OnRequest(te ct.TriggerEvent) (re ct.ResponseEvent) {
	if cc, found := sto.RequestMap[te.Id]; found {
		return cc.OnRequest(te)
	}
	re = ct.ResponseEvent{
		Trigger: &ct.Toast{
			Type:  ct.ToastTypeError,
			Value: fmt.Sprintf("Invalid parameter: %s", te.Id),
		},
		TriggerName: te.Name,
		Name:        te.Name,
		Header: ut.SM{
			ct.HeaderRetarget: "#toast-msg",
			ct.HeaderReswap:   "innerHTML",
		},
	}
	return re
}

func (sto *Demo) response(evt ct.ResponseEvent) (re ct.ResponseEvent) {
	stoEvt := ct.ResponseEvent{
		Trigger: sto, TriggerName: sto.Name,
	}
	var value interface{}
	switch evt.TriggerName {
	case "theme":
		stoEvt.Name = DemoEventTheme
		value = sto.SetProperty("theme", demoIcoMap[sto.Theme][0])

	case "view_size":
		stoEvt.Name = DemoEventViewSize
		value = sto.SetProperty("view_size", demoIcoMap[sto.ViewSize][0])

	default:
		sto.Name = DemoEventChange
		value = sto.SetProperty(evt.TriggerName, evt.Value)
	}
	stoEvt.Value = value
	if sto.OnResponse != nil {
		return sto.OnResponse(stoEvt)
	}
	return stoEvt
}

func (sto *Demo) getComponent(name string) (res string, err error) {
	propValue := ut.SM{
		"theme": sto.Theme, "view_size": sto.ViewSize,
		"selected_group": sto.SelectedGroup, "selected_type": ut.ToString(sto.SelectedType, ""),
		"selected_demo": ut.ToString(sto.SelectedDemo, ""),
	}
	ccBtn := func() *ct.Button {
		return &ct.Button{
			BaseComponent: ct.BaseComponent{
				Id: sto.Id + "_" + name, Name: name,
				Style:        ut.SM{"padding": "8px"},
				EventURL:     sto.EventURL,
				Target:       sto.Id,
				OnResponse:   sto.response,
				RequestValue: sto.RequestValue,
				RequestMap:   sto.RequestMap,
			},
			Type:           ct.ButtonTypePrimary,
			LabelComponent: &ct.Icon{Value: demoIcoMap[propValue[name]][1], Width: 18, Height: 18},
		}
	}
	ccSel := func() *ct.Select {
		return &ct.Select{
			BaseComponent: ct.BaseComponent{
				Id: sto.Id + "_" + name, Name: name,
				EventURL:     sto.EventURL,
				Target:       sto.Id,
				OnResponse:   sto.response,
				RequestValue: sto.RequestValue,
				RequestMap:   sto.RequestMap,
			},
			Value:   propValue[name],
			Options: []ct.SelectOption{},
		}
	}
	ccMap := map[string]func() ct.ClientComponent{
		"theme":     func() ct.ClientComponent { return ccBtn() },
		"view_size": func() ct.ClientComponent { return ccBtn() },
		"selected_group": func() ct.ClientComponent {
			sc := ccSel()
			for _, group := range ComponentGroup {
				sc.Options = append(
					sc.Options, ct.SelectOption{Value: group, Text: group})
			}
			return sc
		},
		"selected_type": func() ct.ClientComponent {
			sc := ccSel()
			for index, v := range sto.DemoMap[sto.SelectedGroup] {
				sc.Options = append(
					sc.Options, ct.SelectOption{Value: ut.ToString(index, ""), Text: v.ComponentType})
			}
			return sc
		},
		"selected_demo": func() ct.ClientComponent {
			sc := ccSel()
			for index, v := range sto.DemoMap[sto.SelectedGroup][sto.SelectedType].Session {
				sc.Options = append(
					sc.Options, ct.SelectOption{Value: ut.ToString(index, ""), Text: v.Label})
			}
			return sc
		},
	}
	cc := ccMap[name]()
	res, err = cc.Render()
	return res, err
}

func (sto *Demo) Render() (res string, err error) {
	sto.InitProps(sto)

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
			return (&ct.Label{
				BaseComponent: ct.BaseComponent{Style: ut.SM{"color": "brown"}},
				Value:         value,
			}).Render()
		},
		"clientComponent": func(cc ct.ClientComponent) (string, error) {
			res, err := cc.Render()
			return res, err
		},
		"stories": func() []DemoSession {
			return sto.DemoMap[sto.SelectedGroup][sto.SelectedType].Session
		},
		"demo": func() DemoSession {
			return sto.DemoMap[sto.SelectedGroup][sto.SelectedType].Session[sto.SelectedDemo]
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

	return ut.TemplateBuilder("demo", tpl, funcMap, sto)
}
