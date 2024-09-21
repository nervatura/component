package component

import (
	"fmt"
	"html/template"
	"strings"

	ut "github.com/nervatura/component/pkg/util"
)

// [Login] constants
const (
	ComponentTypeLogin = "login"

	LoginEventChange = "change"
	LoginEventLogin  = "login"
	LoginEventAuth   = "auth"
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

/*
Creates an application login control

Example component with the following main features:
  - server-side state management
  - selectable label languages
  - light and dark theme
  - modal appearance
  - [Input], [Select], [Label] and [Button] components

For example:

	&Login{
	  BaseComponent: BaseComponent{
	    Id:           "id_login_default",
	    EventURL:     "/event",
	    OnResponse:   func(evt ResponseEvent) (re ResponseEvent) {
	      // return parent_component response
	      return evt
	    },
	    RequestValue: parent_component.GetProperty("request_value").(map[string]ut.IM),
	    RequestMap:   parent_component.GetProperty("request_map").(map[string]ClientComponent),
	    Data: ut.IM{ "username": "admin", "database": "demo" },
	  },
	  Version: "6.0.0",
	  Lang:    "en",
	  Locales: []SelectOption{{Value: "en", Text: "English"}},
	  Theme:  ThemeLight,
	}
*/
type Login struct {
	BaseComponent
	// Application version value
	Version string `json:"version"`
	// Current ui language
	Lang string `json:"lang"`
	// Show or hide the database input control
	HideDatabase bool `json:"hide_database"`
	// Show or hide password login
	HidePassword bool `json:"hide_password"`
	/*
		The theme of the control.
		[Theme] variable constants: [ThemeLight], [ThemeDark]. Default value: [ThemeLight]
	*/
	Theme string `json:"theme"`
	// The texts of the labels of the controls
	Labels ut.SM `json:"labels"`
	// Selectable languages
	Locales []SelectOption `json:"locales"`
	// OAuth buttons
	AuthButtons []LoginAuthButton `json:"auth_buttons"`
}

// OAuth button parameters
type LoginAuthButton struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Icon  string `json:"icon"`
}

/*
Returns all properties of the [Login]
*/
func (lgn *Login) Properties() ut.IM {
	return ut.MergeIM(
		lgn.BaseComponent.Properties(),
		ut.IM{
			"version":       lgn.Version,
			"lang":          lgn.Lang,
			"hide_database": lgn.HideDatabase,
			"hide_password": lgn.HidePassword,
			"theme":         lgn.Theme,
			"labels":        lgn.Labels,
			"locales":       lgn.Locales,
			"auth_buttons":  lgn.AuthButtons,
		})
}

/*
Returns the value of the property of the [Login] with the specified name.
*/
func (lgn *Login) GetProperty(propName string) interface{} {
	return lgn.Properties()[propName]
}

/*
It checks the value given to the property of the [Login] and always returns a valid value
*/
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
		"auth_buttons": func() interface{} {
			if options, valid := propValue.([]LoginAuthButton); valid && (options != nil) {
				return options
			}
			return []LoginAuthButton{}
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

/*
Setting a property of the [Login] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
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
		"hide_database": func() interface{} {
			lgn.HideDatabase = ut.ToBoolean(propValue, false)
			return lgn.HideDatabase
		},
		"hide_password": func() interface{} {
			lgn.HidePassword = ut.ToBoolean(propValue, false)
			return lgn.HidePassword
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
		"auth_buttons": func() interface{} {
			lgn.AuthButtons = lgn.Validation(propName, propValue).([]LoginAuthButton)
			return lgn.AuthButtons
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

	case "auth":
		value := ut.ToString(evt.Trigger.GetProperty("data").(ut.IM)["id"], "")
		lgnEvt.Name = LoginEventAuth
		lgnEvt.Value = value

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

func (lgn *Login) getComponent(name string, authIdx int) (html template.HTML, err error) {
	var loginDisabled bool = ((ut.ToString(lgn.Data["username"], "") == "") || (ut.ToString(lgn.Data["database"], "") == ""))
	if lgn.HideDatabase {
		loginDisabled = (ut.ToString(lgn.Data["username"], "") == "")
	}

	ccInp := func(itype string) *Input {
		return &Input{
			BaseComponent: BaseComponent{
				Id: lgn.Id + "_" + name, Name: name,
				EventURL:     lgn.EventURL,
				Target:       lgn.Target,
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
	ccBtn := func(id, label, icon string) *Button {
		return &Button{
			BaseComponent: BaseComponent{
				Id:           lgn.Id + "_" + name + "_" + id,
				Name:         name,
				EventURL:     lgn.EventURL,
				Target:       lgn.Target,
				OnResponse:   lgn.response,
				RequestValue: lgn.RequestValue,
				RequestMap:   lgn.RequestMap,
				Data: ut.IM{
					"id": id,
				},
			},
			ButtonStyle: ButtonStylePrimary,
			Label:       label,
			Icon:        icon,
			Full:        true,
		}
	}
	ccMap := map[string]func() ClientComponent{
		"username": func() ClientComponent {
			return ccInp(InputTypeString)
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
			return ccInp(InputTypeString)
		},
		"login_database": func() ClientComponent {
			return ccLbl()
		},
		"auth": func() ClientComponent {
			btn := lgn.AuthButtons[authIdx]
			return ccBtn(btn.Id, btn.Label, btn.Icon)
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
				ButtonStyle: ButtonStylePrimary,
				Label:       lgn.Labels["login_"+name],
				Disabled:    loginDisabled,
				Full:        true, AutoFocus: true,
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
				ButtonStyle:    ButtonStyleBorder,
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
	html, err = cc.Render()
	return html, err
}

func (lgn *Login) msg(labelID string) string {
	if label, found := lgn.Labels[labelID]; found {
		return label
	}
	return labelID
}

/*
Based on the values, it will generate the html code of the [Login] or return with an error message.
*/
func (lgn *Login) Render() (html template.HTML, err error) {
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
		"loginComponent": func(name string) (template.HTML, error) {
			return lgn.getComponent(name, 0)
		},
		"authBtn": func(idx int) (template.HTML, error) {
			return lgn.getComponent("auth", idx)
		},
		"even": func(idx int) bool {
			return (idx%2 == 0)
		},
		"odd": func(idx int) bool {
			return !(idx%2 == 0)
		},
		"buttons": func() bool {
			return len(lgn.AuthButtons) > 0
		},
	}
	tpl := `<div id="{{ .Id }}" class="login-modal {{ customClass }}" theme="{{ .Theme }}" 
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}>
	<div class="middle"><div class="dialog">
	<div class="row title">
	<div class="cell title-cell login-title-cell" ><span>{{ msg "title_login" }}</span></div>
	<div class="cell version-cell" ><span>{{ .Version }}</span></div>
	</div>
	{{ if ne .HidePassword true }}
	<div class="row full section-small" >
	<div class="row full section-small" >
	<div class="cell label-cell padding-normal mobile" >{{ loginComponent "login_username" }}</div>
	<div class="cell container mobile" >{{ loginComponent "username" }}</div>
	</div>
	<div class="row full {{ if .HideDatabase }}section-small-bottom{{ end }}" >
	<div class="cell label-cell padding-normal mobile" >{{ loginComponent "login_password" }}</div>
	<div class="cell container mobile" >{{ loginComponent "password" }}</div>
	</div>
	{{ if ne .HideDatabase true }}
	<div class="row full section-small" >
	<div class="cell label-cell padding-normal mobile" >{{ loginComponent "login_database" }}</div>
	<div class="cell container mobile" >{{ loginComponent "database" }}</div>
	</div>
	{{ end }}
	</div>
	{{ end }}
	{{ if buttons }}<div class="row full section border-top" >
	{{ range $index, $auth := .AuthButtons }}
	{{ if even $index }}<div class="row full container-small section-small" >{{ end }}
	<div class="cell container-small mobile" >{{ authBtn $index }}</div>
	{{ if odd $index }}</div>{{ end }}
	{{ end }}
	</div>{{ end }}
  <div class="row full section buttons" >
	{{ if ne .HidePassword true }}<div class="cell section-small mobile" >{{ end }}
	<div class="cell container align-right" >
	{{ loginComponent "theme" }}
	</div>
	<div class="cell" >
	{{ loginComponent "lang" }}
	</div>
	{{ if ne .HidePassword true }}</div>{{ end }}
	{{ if ne .HidePassword true }}
	<div class="cell container section-small align-right mobile" >
	{{ loginComponent "login" }}
	</div>
	{{ end }}
	</div>
	</div></div></div>`

	return ut.TemplateBuilder("login", tpl, funcMap, lgn)
}

var testLoginLabels map[string]ut.SM = map[string]ut.SM{
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

var testLoginResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
	toast := func(value string) ResponseEvent {
		return ResponseEvent{
			Trigger: &Toast{
				Type:    ToastTypeSuccess,
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
	switch evt.Name {
	case LoginEventLogin:
		data := evt.Trigger.GetProperty("data").(ut.IM)
		labels := evt.Trigger.GetProperty("labels").(ut.SM)
		value := fmt.Sprintf(`%s: %s, %s: %s`,
			labels["login_username"], data["username"], labels["login_database"], data["database"])
		return toast(value)
	case LoginEventAuth:
		return toast(ut.ToString(evt.Value, ""))
	case LoginEventLang:
		value := ut.ToString(evt.Value, "en")
		labels := ut.MergeSM(nil, testLoginLabels[value])
		evt.Trigger.SetProperty("labels", labels)
	default:
	}
	return evt
}

// [Login] test and demo data
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
					OnResponse:   testLoginResponse,
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
				Labels: ut.MergeSM(nil, testLoginLabels["en"]),
			}},
		{
			Label:         "Hide database",
			ComponentType: ComponentTypeLogin,
			Component: &Login{
				BaseComponent: BaseComponent{
					Id:           id + "_login_nodb",
					EventURL:     eventURL,
					OnResponse:   testLoginResponse,
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
				Theme:        ThemeLight,
				Labels:       ut.MergeSM(nil, testLoginLabels["en"]),
				HideDatabase: true,
			}},
		{
			Label:         "Auth buttons",
			ComponentType: ComponentTypeLogin,
			Component: &Login{
				BaseComponent: BaseComponent{
					Id:           id + "_login_auth",
					EventURL:     eventURL,
					OnResponse:   testLoginResponse,
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
				Theme:        ThemeLight,
				Labels:       ut.MergeSM(nil, testLoginLabels["en"]),
				HideDatabase: true,
				AuthButtons: []LoginAuthButton{
					{Id: "google", Label: "Google", Icon: "Google"},
					{Id: "facebook", Label: "Facebook", Icon: "Facebook"},
					{Id: "github", Label: "Github", Icon: "Github"},
					{Id: "microsoft", Label: "Microsoft", Icon: "Microsoft"},
				},
			}},
		{
			Label:         "Hide password",
			ComponentType: ComponentTypeLogin,
			Component: &Login{
				BaseComponent: BaseComponent{
					Id:           id + "_login_hide_passw",
					EventURL:     eventURL,
					OnResponse:   testLoginResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Data:         ut.IM{},
				},
				Version: "6.0.0",
				Lang:    "en",
				Locales: []SelectOption{
					{Value: "en", Text: "English"},
					{Value: "de", Text: "Deutsch"},
				},
				Theme:        ThemeLight,
				Labels:       ut.MergeSM(nil, testLoginLabels["en"]),
				HideDatabase: true,
				AuthButtons: []LoginAuthButton{
					{Id: "google", Label: "Google", Icon: "Google"},
					{Id: "facebook", Label: "Facebook", Icon: "Facebook"},
					{Id: "github", Label: "Github", Icon: "Github"},
					{Id: "microsoft", Label: "Microsoft", Icon: "Microsoft"},
				},
				HidePassword: true,
			}},
	}
}
