package modal

import (
	"fmt"
	"strings"

	fm "github.com/nervatura/component/component/atom"
	bc "github.com/nervatura/component/component/base"
)

const (
	LoginEventChange = "change"
	LoginEventLogin  = "login"
	LoginEventTheme  = "theme"
	LoginEventLang   = "lang"
)

var loginDefaultLabel bc.SM = bc.SM{
	"title_login":    "Nervatura Client",
	"login_username": "Username",
	"login_password": "Password",
	"login_database": "Database",
	"login_lang":     "Language",
	"login_login":    "Login",
	"login_theme":    "Theme",
}

var loginThemeMap map[string][]string = map[string][]string{
	bc.ThemeDark: {bc.ThemeLight, "Sun"}, bc.ThemeLight: {bc.ThemeDark, "Moon"},
}

type Login struct {
	bc.BaseComponent
	Version string            `json:"version"`
	Lang    string            `json:"lang"`
	Theme   string            `json:"theme"`
	Labels  bc.SM             `json:"labels"`
	Locales []fm.SelectOption `json:"locales"`
}

func (lgn *Login) Properties() bc.IM {
	return bc.MergeIM(
		lgn.BaseComponent.Properties(),
		bc.IM{
			"version": lgn.Version,
			"lang":    lgn.Lang,
			"theme":   lgn.Theme,
			"labels":  lgn.Labels,
			"locales": lgn.Locales,
		})
}

func (lgn *Login) GetProperty(propName string) interface{} {
	return lgn.Properties()[propName]
}

func (lgn *Login) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"theme": func() interface{} {
			return lgn.CheckEnumValue(bc.ToString(propValue, ""), bc.ThemeLight, bc.Theme)
		},
		"labels": func() interface{} {
			value := bc.SetSMValue(lgn.Labels, "", "")
			if smap, valid := propValue.(bc.SM); valid {
				value = bc.MergeSM(value, smap)
			}
			if len(value) == 0 {
				value = loginDefaultLabel
			}
			return value
		},
		"locales": func() interface{} {
			value := lgn.Locales
			if loc, valid := propValue.([]fm.SelectOption); valid && len(loc) > 0 {
				value = loc
			}
			if len(value) == 0 {
				lang := bc.ToString(lgn.Lang, "en")
				value = []fm.SelectOption{{Value: lang, Text: lang}}
			}
			return value
		},
		"target": func() interface{} {
			lgn.SetProperty("id", lgn.Id)
			value := bc.ToString(propValue, lgn.Id)
			if value != "this" && !strings.HasPrefix(value, "#") {
				value = "#" + value
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if lgn.BaseComponent.GetProperty(propName) != nil {
		return lgn.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

func (lgn *Login) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"version": func() interface{} {
			lgn.Version = bc.ToString(propValue, "")
			return lgn.Version
		},
		"lang": func() interface{} {
			lgn.Lang = bc.ToString(propValue, "en")
			return lgn.Lang
		},
		"locales": func() interface{} {
			lgn.Locales = lgn.Validation(propName, propValue).([]fm.SelectOption)
			return lgn.Locales
		},
		"theme": func() interface{} {
			lgn.Theme = lgn.Validation(propName, propValue).(string)
			return lgn.Theme
		},
		"labels": func() interface{} {
			lgn.Labels = lgn.Validation(propName, propValue).(bc.SM)
			return lgn.Labels
		},
		"target": func() interface{} {
			lgn.Target = lgn.Validation(propName, propValue).(string)
			return lgn.Target
		},
	}
	if _, found := pm[propName]; found {
		return lgn.SetRequestValue(propName, pm[propName](), []string{})
	}
	if lgn.BaseComponent.GetProperty(propName) != nil {
		return lgn.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (lgn *Login) response(evt bc.ResponseEvent) (re bc.ResponseEvent) {
	lgnEvt := bc.ResponseEvent{
		Trigger: lgn, TriggerName: lgn.Name, Value: evt.Value,
	}
	switch evt.TriggerName {

	case "username", "password", "database":
		lgnEvt.Name = LoginEventChange
		lgn.SetProperty("data", bc.IM{evt.TriggerName: lgnEvt.Value})

	case "theme":
		lgnEvt.Name = LoginEventTheme
		lgn.SetProperty("theme", loginThemeMap[lgn.Theme][0])

	case "login":
		lgnEvt.Name = LoginEventLogin

	case "lang":
		lgnEvt.Name = LoginEventLang
		lgn.SetProperty("lang", lgnEvt.Value)

	default:
	}
	if lgn.OnResponse != nil {
		return lgn.OnResponse(lgnEvt)
	}
	return lgnEvt
}

func (lgn *Login) getComponent(name string) (res string, err error) {
	var loginDisabled bool = ((bc.ToString(lgn.Data["username"], "") == "") || (bc.ToString(lgn.Data["database"], "") == ""))
	ccInp := func(itype string) *fm.Input {
		return &fm.Input{
			BaseComponent: bc.BaseComponent{
				Id: lgn.Id + "_" + name, Name: name,
				EventURL:     lgn.EventURL,
				Target:       lgn.Target,
				Swap:         bc.SwapOuterHTML,
				OnResponse:   lgn.response,
				RequestValue: lgn.RequestValue,
				RequestMap:   lgn.RequestMap,
			},
			Type:  itype,
			Label: lgn.Labels["login_"+name],
			Value: bc.ToString(lgn.Data[name], ""),
			Full:  true,
		}
	}
	ccLbl := func() *fm.Label {
		return &fm.Label{
			Value: lgn.Labels[name],
		}
	}
	ccMap := map[string]func() bc.ClientComponent{
		"username": func() bc.ClientComponent {
			return ccInp(fm.InputTypeText)
		},
		"login_username": func() bc.ClientComponent {
			return ccLbl()
		},
		"password": func() bc.ClientComponent {
			return ccInp(fm.InputTypePassword)
		},
		"login_password": func() bc.ClientComponent {
			return ccLbl()
		},
		"database": func() bc.ClientComponent {
			return ccInp(fm.InputTypeText)
		},
		"login_database": func() bc.ClientComponent {
			return ccLbl()
		},
		"login": func() bc.ClientComponent {
			return &fm.Button{
				BaseComponent: bc.BaseComponent{
					Id: lgn.Id + "_" + name, Name: name,
					EventURL:     lgn.EventURL,
					Target:       lgn.Target,
					OnResponse:   lgn.response,
					RequestValue: lgn.RequestValue,
					RequestMap:   lgn.RequestMap,
				},
				Type:     fm.ButtonTypePrimary,
				Label:    lgn.Labels["login_"+name],
				Disabled: loginDisabled,
				Full:     true, AutoFocus: true,
			}
		},
		"theme": func() bc.ClientComponent {
			return &fm.Button{
				BaseComponent: bc.BaseComponent{
					Id: lgn.Id + "_" + name, Name: name,
					EventURL:     lgn.EventURL,
					Target:       lgn.Target,
					OnResponse:   lgn.response,
					RequestValue: lgn.RequestValue,
					RequestMap:   lgn.RequestMap,
				},
				Type:           fm.ButtonTypeBorder,
				Label:          lgn.Labels["login_"+name],
				LabelComponent: &fm.Icon{Value: loginThemeMap[lgn.Theme][1], Width: 18, Height: 18},
			}
		},
		"lang": func() bc.ClientComponent {
			return &fm.Select{
				BaseComponent: bc.BaseComponent{
					Id: lgn.Id + "_" + name, Name: name,
					EventURL:     lgn.EventURL,
					Target:       lgn.Target,
					OnResponse:   lgn.response,
					RequestValue: lgn.RequestValue,
					RequestMap:   lgn.RequestMap,
				},
				Label:   lgn.Labels["login_"+name],
				IsNull:  false,
				Value:   bc.ToString(lgn.Lang, ""),
				Options: lgn.Locales,
			}
		},
	}
	cc := ccMap[name]()
	res, err = cc.Render()
	return res, err
}

func (lgn *Login) msg(labelID string) string {
	if label, found := lgn.Labels[labelID]; found {
		return label
	}
	return labelID
}

func (lgn *Login) Render() (res string, err error) {
	lgn.InitProps(lgn)

	funcMap := map[string]any{
		"msg": func(labelID string) string {
			return lgn.msg(labelID)
		},
		"styleMap": func() bool {
			return len(lgn.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(lgn.Class, " ")
		},
		"loginComponent": func(name string) (string, error) {
			return lgn.getComponent(name)
		},
	}
	tpl := `<div id="{{ .Id }}" class="modal {{ customClass }}" theme="{{ .Theme }}" 
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}>
	<div class="middle"><div class="dialog">
	<div class="row title">
	<div class="cell title-cell" ><span>{{ msg "title_login" }}</span></div>
	<div class="cell version-cell" ><span>{{ .Version }}</span></div>
	</div>
	<div class="row full section-small" >
	<div class="row full section-small" >
	<div class="cell label-cell padding-normal mobile" >{{ loginComponent "login_username" }}</div>
	<div class="cell container mobile" >{{ loginComponent "username" }}</div>
	</div>
	<div class="row full" >
	<div class="cell label-cell padding-normal mobile" >{{ loginComponent "login_password" }}</div>
	<div class="cell container mobile" >{{ loginComponent "password" }}</div>
	</div>
	<div class="row full section-small" >
	<div class="cell label-cell padding-normal mobile" >{{ loginComponent "login_database" }}</div>
	<div class="cell container mobile" >{{ loginComponent "database" }}</div>
	</div>
	</div>
	<div class="row full section buttons" >
	<div class="cell section-small mobile" >
	<div class="cell container" >
	{{ loginComponent "theme" }}
	</div>
	<div class="cell" >
	{{ loginComponent "lang" }}
	</div>
	</div>
	<div class="cell container section-small align-right mobile" >
	{{ loginComponent "login" }}
	</div>
	</div>
	</div></div></div>`

	return bc.TemplateBuilder("login", tpl, funcMap, lgn)
}

var demoLabels map[string]bc.SM = map[string]bc.SM{
	"en": loginDefaultLabel,
	"de": {
		"title_login":    "Nervatura Client",
		"login_username": "Nutzername",
		"login_password": "Passwort",
		"login_database": "Datenbank",
		"login_lang":     "Sprache",
		"login_login":    "Einloggen",
		"login_theme":    "Thema",
	},
}

var demoLoginResponse func(evt bc.ResponseEvent) (re bc.ResponseEvent) = func(evt bc.ResponseEvent) (re bc.ResponseEvent) {
	switch evt.Name {
	case LoginEventLogin:
		data := evt.Trigger.GetProperty("data").(bc.IM)
		labels := evt.Trigger.GetProperty("labels").(bc.SM)
		value := fmt.Sprintf(`%s: %s, %s: %s`,
			labels["login_username"], data["username"], labels["login_database"], data["database"])
		re = bc.ResponseEvent{
			Trigger: &fm.Toast{
				Type:    fm.ToastTypeSuccess,
				Value:   value,
				Timeout: 4,
			},
			TriggerName: evt.TriggerName,
			Name:        evt.Name,
			Header: bc.SM{
				bc.HeaderRetarget: "#toast-msg",
				bc.HeaderReswap:   "innerHTML",
			},
		}
		return re
	case LoginEventLang:
		value := bc.ToString(evt.Value, "en")
		labels := bc.MergeSM(nil, demoLabels[value])
		evt.Trigger.SetProperty("labels", labels)
	default:
	}
	return evt
}

func DemoLogin(demo bc.ClientComponent) []bc.DemoComponent {
	id := bc.ToString(demo.GetProperty("id"), "")
	eventURL := bc.ToString(demo.GetProperty("event_url"), "")
	requestValue := demo.GetProperty("request_value").(map[string]bc.IM)
	requestMap := demo.GetProperty("request_map").(map[string]bc.ClientComponent)
	return []bc.DemoComponent{
		{
			Label:         "Default",
			ComponentType: bc.ComponentTypeLogin,
			Component: &Login{
				BaseComponent: bc.BaseComponent{
					Id:           id + "_login_default",
					EventURL:     eventURL,
					OnResponse:   demoLoginResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Data: bc.IM{
						"username": "admin",
						"database": "demo",
					},
				},
				Version: "6.0.0",
				Lang:    "en",
				Locales: []fm.SelectOption{
					{Value: "en", Text: "English"},
					{Value: "de", Text: "Deutsch"},
				},
				Theme:  bc.ThemeLight,
				Labels: bc.MergeSM(nil, demoLabels["en"]),
			}},
	}
}
