package component

import (
	"fmt"
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

const (
	LoginEventChange = "change"
	LoginEventLogin  = "login"
	LoginEventTheme  = "theme"
	LoginEventLang   = "lang"
)

var loginDefaultLabel ut.SM = ut.SM{
	"title_login":    "Nervatura Client",
	"login_username": "Username",
	"login_password": "Password",
	"login_database": "Database",
	"login_lang":     "Language",
	"login_login":    "Login",
	"login_theme":    "Theme",
}

var loginThemeMap map[string][]string = map[string][]string{
	ThemeDark: {ThemeLight, "Sun"}, ThemeLight: {ThemeDark, "Moon"},
}

type Login struct {
	BaseComponent
	Version string         `json:"version"`
	Lang    string         `json:"lang"`
	Theme   string         `json:"theme"`
	Labels  ut.SM          `json:"labels"`
	Locales []SelectOption `json:"locales"`
}

func (lgn *Login) Properties() ut.IM {
	return ut.MergeIM(
		lgn.BaseComponent.Properties(),
		ut.IM{
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
			return lgn.CheckEnumValue(ut.ToString(propValue, ""), ThemeLight, Theme)
		},
		"labels": func() interface{} {
			value := ut.SetSMValue(lgn.Labels, "", "")
			if smap, valid := propValue.(ut.SM); valid {
				value = ut.MergeSM(value, smap)
			}
			if len(value) == 0 {
				value = loginDefaultLabel
			}
			return value
		},
		"locales": func() interface{} {
			value := lgn.Locales
			if loc, valid := propValue.([]SelectOption); valid && len(loc) > 0 {
				value = loc
			}
			if len(value) == 0 {
				lang := ut.ToString(lgn.Lang, "en")
				value = []SelectOption{{Value: lang, Text: lang}}
			}
			return value
		},
		"target": func() interface{} {
			lgn.SetProperty("id", lgn.Id)
			value := ut.ToString(propValue, lgn.Id)
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
			lgn.Version = ut.ToString(propValue, "")
			return lgn.Version
		},
		"lang": func() interface{} {
			lgn.Lang = ut.ToString(propValue, "en")
			return lgn.Lang
		},
		"locales": func() interface{} {
			lgn.Locales = lgn.Validation(propName, propValue).([]SelectOption)
			return lgn.Locales
		},
		"theme": func() interface{} {
			lgn.Theme = lgn.Validation(propName, propValue).(string)
			return lgn.Theme
		},
		"labels": func() interface{} {
			lgn.Labels = lgn.Validation(propName, propValue).(ut.SM)
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

func (lgn *Login) response(evt ResponseEvent) (re ResponseEvent) {
	lgnEvt := ResponseEvent{
		Trigger: lgn, TriggerName: lgn.Name, Value: evt.Value,
	}
	switch evt.TriggerName {

	case "username", "password", "database":
		lgnEvt.Name = LoginEventChange
		lgn.SetProperty("data", ut.IM{evt.TriggerName: lgnEvt.Value})

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
	var loginDisabled bool = ((ut.ToString(lgn.Data["username"], "") == "") || (ut.ToString(lgn.Data["database"], "") == ""))
	ccInp := func(itype string) *Input {
		return &Input{
			BaseComponent: BaseComponent{
				Id: lgn.Id + "_" + name, Name: name,
				EventURL:     lgn.EventURL,
				Target:       lgn.Target,
				Swap:         SwapOuterHTML,
				OnResponse:   lgn.response,
				RequestValue: lgn.RequestValue,
				RequestMap:   lgn.RequestMap,
			},
			Type:  itype,
			Label: lgn.Labels["login_"+name],
			Value: ut.ToString(lgn.Data[name], ""),
			Full:  true,
		}
	}
	ccLbl := func() *Label {
		return &Label{
			Value: lgn.Labels[name],
		}
	}
	ccMap := map[string]func() ClientComponent{
		"username": func() ClientComponent {
			return ccInp(InputTypeText)
		},
		"login_username": func() ClientComponent {
			return ccLbl()
		},
		"password": func() ClientComponent {
			return ccInp(InputTypePassword)
		},
		"login_password": func() ClientComponent {
			return ccLbl()
		},
		"database": func() ClientComponent {
			return ccInp(InputTypeText)
		},
		"login_database": func() ClientComponent {
			return ccLbl()
		},
		"login": func() ClientComponent {
			return &Button{
				BaseComponent: BaseComponent{
					Id: lgn.Id + "_" + name, Name: name,
					EventURL:     lgn.EventURL,
					Target:       lgn.Target,
					OnResponse:   lgn.response,
					RequestValue: lgn.RequestValue,
					RequestMap:   lgn.RequestMap,
				},
				Type:     ButtonTypePrimary,
				Label:    lgn.Labels["login_"+name],
				Disabled: loginDisabled,
				Full:     true, AutoFocus: true,
			}
		},
		"theme": func() ClientComponent {
			return &Button{
				BaseComponent: BaseComponent{
					Id: lgn.Id + "_" + name, Name: name,
					EventURL:     lgn.EventURL,
					Target:       lgn.Target,
					OnResponse:   lgn.response,
					RequestValue: lgn.RequestValue,
					RequestMap:   lgn.RequestMap,
				},
				Type:           ButtonTypeBorder,
				Label:          lgn.Labels["login_"+name],
				LabelComponent: &Icon{Value: loginThemeMap[lgn.Theme][1], Width: 18, Height: 18},
			}
		},
		"lang": func() ClientComponent {
			return &Select{
				BaseComponent: BaseComponent{
					Id: lgn.Id + "_" + name, Name: name,
					EventURL:     lgn.EventURL,
					Target:       lgn.Target,
					OnResponse:   lgn.response,
					RequestValue: lgn.RequestValue,
					RequestMap:   lgn.RequestMap,
				},
				Label:   lgn.Labels["login_"+name],
				IsNull:  false,
				Value:   ut.ToString(lgn.Lang, ""),
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

	return ut.TemplateBuilder("login", tpl, funcMap, lgn)
}

var demoLabels map[string]ut.SM = map[string]ut.SM{
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

var demoLoginResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
	switch evt.Name {
	case LoginEventLogin:
		data := evt.Trigger.GetProperty("data").(ut.IM)
		labels := evt.Trigger.GetProperty("labels").(ut.SM)
		value := fmt.Sprintf(`%s: %s, %s: %s`,
			labels["login_username"], data["username"], labels["login_database"], data["database"])
		re = ResponseEvent{
			Trigger: &Toast{
				Type:    ToastTypeSuccess,
				Value:   value,
				Timeout: 4,
			},
			TriggerName: evt.TriggerName,
			Name:        evt.Name,
			Header: ut.SM{
				HeaderRetarget: "#toast-msg",
				HeaderReswap:   "innerHTML",
			},
		}
		return re
	case LoginEventLang:
		value := ut.ToString(evt.Value, "en")
		labels := ut.MergeSM(nil, demoLabels[value])
		evt.Trigger.SetProperty("labels", labels)
	default:
	}
	return evt
}

func TestLogin(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	return []TestComponent{
		{
			Label:         "Default",
			ComponentType: ComponentTypeLogin,
			Component: &Login{
				BaseComponent: BaseComponent{
					Id:           id + "_login_default",
					EventURL:     eventURL,
					OnResponse:   demoLoginResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Data: ut.IM{
						"username": "admin",
						"database": "demo",
					},
				},
				Version: "6.0.0",
				Lang:    "en",
				Locales: []SelectOption{
					{Value: "en", Text: "English"},
					{Value: "de", Text: "Deutsch"},
				},
				Theme:  ThemeLight,
				Labels: ut.MergeSM(nil, demoLabels["en"]),
			}},
	}
}
