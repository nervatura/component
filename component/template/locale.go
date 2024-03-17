package template

import (
	"fmt"
	"sort"
	"strings"

	fm "github.com/nervatura/component/component/atom"
	bc "github.com/nervatura/component/component/base"
	mc "github.com/nervatura/component/component/molecule"
)

const (
	LocalesEventChange = "change"
	LocalesEventUndo   = "undo"
	LocalesEventSave   = "save"
	LocalesEventError  = "error"
)

var localeDefaultLabel bc.SM = bc.SM{
	"locale_title":         "Translation helper tool",
	"locale_missing":       "Missing values",
	"locale_update":        "Write changes to the json files",
	"locale_undo":          "Discard all changes",
	"locale_add":           "Add a new language",
	"locale_filter":        "Filter rows",
	"locale_lcode":         "Lang.code (e.g. de)",
	"locale_lname":         "Lang.name(e.g.Deutsche)",
	"locale_existing_lang": "Existing language code",
	"locale_tag":           "Tag",
	"locale_key":           "Fieldname",
	"locale_value":         "Value",
}

type Locale struct {
	bc.BaseComponent
	Locales     []fm.SelectOption `json:"locales"`
	TagKeys     []fm.SelectOption `json:"tag_keys"`
	FilterValue string            `json:"filter_value"`
	Dirty       bool              `json:"dirty"`
	AddItem     bool              `json:"add_item"`
	Labels      bc.SM             `json:"labels"`
}

func (loc *Locale) Properties() bc.IM {
	return bc.MergeIM(
		loc.BaseComponent.Properties(),
		bc.IM{
			"locales":      loc.Locales,
			"tag_keys":     loc.TagKeys,
			"filter_value": loc.FilterValue,
			"dirty":        loc.Dirty,
			"add_item":     loc.AddItem,
			"labels":       loc.Labels,
		})
}

func (loc *Locale) GetProperty(propName string) interface{} {
	return loc.Properties()[propName]
}

func (loc *Locale) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"labels": func() interface{} {
			value := bc.SetSMValue(loc.Labels, "", "")
			if smap, valid := propValue.(bc.SM); valid {
				value = bc.MergeSM(value, smap)
			}
			if len(value) == 0 {
				value = localeDefaultLabel
			}
			return value
		},
		"locales": func() interface{} {
			value := loc.Locales
			if lang, valid := propValue.([]fm.SelectOption); valid && len(lang) > 0 {
				value = lang
			}
			if len(value) == 0 {
				lang := bc.ToString(loc.Data["locales"], "default")
				value = []fm.SelectOption{{Value: lang, Text: lang}}
			}
			return value
		},
		"tag_keys": func() interface{} {
			value := loc.TagKeys
			if key, valid := propValue.([]fm.SelectOption); valid && len(key) > 0 {
				value = key
			}
			return value
		},
		"target": func() interface{} {
			loc.SetProperty("id", loc.Id)
			value := bc.ToString(propValue, loc.Id)
			if value != "this" && !strings.HasPrefix(value, "#") {
				value = "#" + value
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if loc.BaseComponent.GetProperty(propName) != nil {
		return loc.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

func (loc *Locale) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"locales": func() interface{} {
			loc.Locales = loc.Validation(propName, propValue).([]fm.SelectOption)
			return loc.Locales
		},
		"tag_keys": func() interface{} {
			loc.TagKeys = loc.Validation(propName, propValue).([]fm.SelectOption)
			return loc.TagKeys
		},
		"filter_value": func() interface{} {
			loc.FilterValue = bc.ToString(propValue, "")
			return loc.FilterValue
		},
		"dirty": func() interface{} {
			loc.Dirty = bc.ToBoolean(propValue, false)
			return loc.Dirty
		},
		"add_item": func() interface{} {
			loc.AddItem = bc.ToBoolean(propValue, false)
			return loc.AddItem
		},
		"labels": func() interface{} {
			loc.Labels = loc.Validation(propName, propValue).(bc.SM)
			return loc.Labels
		},
		"target": func() interface{} {
			loc.Target = loc.Validation(propName, propValue).(string)
			return loc.Target
		},
	}
	if _, found := pm[propName]; found {
		return loc.SetRequestValue(propName, pm[propName](), []string{})
	}
	if loc.BaseComponent.GetProperty(propName) != nil {
		return loc.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (loc *Locale) response(evt bc.ResponseEvent) (re bc.ResponseEvent) {
	locEvt := bc.ResponseEvent{
		Trigger: loc, TriggerName: loc.Name, Value: evt.Value,
	}
	switch evt.TriggerName {
	case "values":
		return evt

	case "tag_keys", "locales":
		locEvt.Name = LocalesEventChange
		loc.SetProperty("data", bc.IM{evt.TriggerName: locEvt.Value})
		loc.SetProperty("filter_value", "")
		if evt.TriggerName == "locales" {
			loc.SetProperty("data", bc.IM{"tag_keys": loc.TagKeys[0].Value})
		}

	case "undo":
		locEvt.Name = LocalesEventUndo

	case "update":
		locEvt.Name = LocalesEventSave

	case "add":
		lang_key := bc.ToString(loc.GetProperty("data").(bc.IM)["lang_key"], "")
		lang_name := bc.ToString(loc.GetProperty("data").(bc.IM)["lang_name"], "")
		locales := loc.Data["locfile"].(bc.IM)["locales"].(bc.IM)
		if _, found := locales[lang_key].(bc.IM); found || lang_key == "en" {
			locEvt.Name = LocalesEventError
			locEvt.Value = loc.msg("locale_existing_lang")
		} else if lang_key == "" || lang_name == "" {
			locEvt.Name = LocalesEventError
			locEvt.Value = loc.msg("locale_missing")
		} else {
			locEvt.Name = LocalesEventChange
			locales[lang_key] = bc.IM{
				"key":    lang_key,
				lang_key: lang_name,
			}
			langs := append(loc.Locales, fm.SelectOption{Value: lang_key, Text: lang_key})
			loc.SetProperty("locales", langs)
			loc.SetProperty("add_item", false)
			loc.SetProperty("data", bc.IM{"lang_key": ""})
			loc.SetProperty("data", bc.IM{"lang_name": ""})
		}

	case "missing":
		locEvt.Name = LocalesEventChange
		loc.SetProperty("data", bc.IM{"tag_keys": "missing"})
		loc.SetProperty("filter_value", "")

	case "tag_cell":
		locEvt.Name = LocalesEventChange
		loc.SetProperty("data", bc.IM{"tag_keys": locEvt.Value})
		loc.SetProperty("filter_value", "")

	case "value_cell":
		locEvt.Name = LocalesEventChange
		key := bc.ToString(evt.Trigger.GetProperty("data").(bc.IM)["key"], "")
		locales := loc.Data["locfile"].(bc.IM)["locales"].(bc.IM)
		lang := bc.ToString(loc.Data["locales"], "")
		if langValues, found := locales[lang].(bc.IM); found {
			langValues[key] = locEvt.Value
		}
		loc.SetProperty("dirty", true)

	case "lang_key", "lang_name":
		locEvt.Name = LocalesEventChange
		loc.SetProperty("data", bc.IM{evt.TriggerName: locEvt.Value})
		loc.SetProperty("dirty", true)

	case "filter":
		locEvt.Name = LocalesEventChange
		loc.SetProperty("filter_value", locEvt.Value)

	case "add_item":
		loc.SetProperty("add_item", !loc.AddItem)

	default:
	}
	if loc.OnResponse != nil {
		return loc.OnResponse(locEvt)
	}
	return locEvt
}

func (loc *Locale) getComponent(name string, data bc.IM) (res string, err error) {
	ccSel := func(options []fm.SelectOption) *fm.Select {
		sel := &fm.Select{
			BaseComponent: bc.BaseComponent{
				Id: loc.Id + "_" + name, Name: name,
				EventURL:     loc.EventURL,
				Target:       loc.Target,
				OnResponse:   loc.response,
				RequestValue: loc.RequestValue,
				RequestMap:   loc.RequestMap,
			},
			IsNull:  false,
			Options: options,
		}
		sel.SetProperty("value", loc.Data[name])
		return sel
	}
	ccBtn := func(icoKey, label, btype string) *fm.Button {
		return &fm.Button{
			BaseComponent: bc.BaseComponent{
				Id: loc.Id + "_" + name, Name: name,
				Style:        bc.SM{"padding": "8px"},
				EventURL:     loc.EventURL,
				Target:       loc.Id,
				OnResponse:   loc.response,
				RequestValue: loc.RequestValue,
				RequestMap:   loc.RequestMap,
			},
			Type:           btype,
			Label:          loc.msg(label),
			LabelComponent: &fm.Icon{Value: icoKey, Width: 18, Height: 18},
		}
	}
	ccInp := func(label, placeholder, value string) *fm.Input {
		inp := &fm.Input{
			BaseComponent: bc.BaseComponent{
				Id:           loc.Id + "_" + name + "_" + bc.ToString(data["key"], ""),
				Name:         name,
				EventURL:     loc.EventURL,
				Target:       loc.Target,
				Swap:         bc.SwapOuterHTML,
				OnResponse:   loc.response,
				RequestValue: loc.RequestValue,
				RequestMap:   loc.RequestMap,
				Data:         data,
			},
			Type:        fm.InputTypeText,
			Label:       label,
			Placeholder: placeholder,
			Full:        true,
		}
		inp.SetProperty("value", value)
		return inp
	}
	ccTbl := func(rowKey string, rows []bc.IM, fields []mc.TableField) *mc.Table {
		tbl := &mc.Table{
			BaseComponent: bc.BaseComponent{
				Id: loc.Id + "_" + name, Name: name,
				EventURL:     loc.EventURL,
				OnResponse:   loc.response,
				RequestValue: loc.RequestValue,
				RequestMap:   loc.RequestMap,
			},
			Rows:        rows,
			Fields:      fields,
			Pagination:  mc.PaginationTypeTop,
			PageSize:    10,
			RowKey:      rowKey,
			TableFilter: false,
			AddItem:     false,
		}
		return tbl
	}
	ccMap := map[string]func() bc.ClientComponent{
		"locales": func() bc.ClientComponent {
			return ccSel(loc.Locales)
		},
		"tag_keys": func() bc.ClientComponent {
			keys := loc.TagKeys
			if bc.ToString(loc.Data["locales"], "") != "default" {
				keys = append(keys, fm.SelectOption{Value: "missing", Text: "missing"})
			}
			return ccSel(keys)
		},
		"missing": func() bc.ClientComponent {
			return ccBtn("QuestionCircle", "locale_missing", fm.ButtonTypeDefault)
		},
		"filter": func() bc.ClientComponent {
			return ccInp(loc.msg("locale_filter"), loc.msg("locale_filter"), loc.FilterValue)
		},
		"update": func() bc.ClientComponent {
			return ccBtn("Check", "locale_update", fm.ButtonTypePrimary)
		},
		"undo": func() bc.ClientComponent {
			return ccBtn("Undo", "locale_undo", fm.ButtonTypePrimary)
		},
		"add_item": func() bc.ClientComponent {
			icon := "Plus"
			if loc.AddItem {
				icon = "ArrowUp"
			}
			return ccBtn(icon, "locale_add", fm.ButtonTypeDefault)
		},
		"add": func() bc.ClientComponent {
			return ccBtn("Plus", "locale_add", fm.ButtonTypeDefault)
		},
		"lang_key": func() bc.ClientComponent {
			lang_key := bc.ToString(loc.GetProperty("data").(bc.IM)["lang_key"], "")
			inp := ccInp(loc.msg("locale_lcode"), loc.msg("locale_lcode"), lang_key)
			inp.Full = false
			inp.MaxLength = 5
			inp.Style = bc.SM{"text-transform": "lowercase"}
			return inp
		},
		"lang_name": func() bc.ClientComponent {
			lang_name := bc.ToString(loc.GetProperty("data").(bc.IM)["lang_name"], "")
			return ccInp(loc.msg("locale_lname"), loc.msg("locale_lname"), lang_name)
		},
		"tag_cell": func() bc.ClientComponent {
			return &fm.Label{
				BaseComponent: bc.BaseComponent{
					Id:           loc.Id + "_" + bc.ToString(data["key"], ""),
					Name:         name,
					EventURL:     loc.EventURL,
					Target:       loc.Target,
					Data:         data,
					OnResponse:   loc.response,
					RequestValue: loc.RequestValue,
					RequestMap:   loc.RequestMap,
				},
				Value: bc.ToString(data["tag"], ""),
			}
		},
		"value_cell": func() bc.ClientComponent {
			return ccInp(bc.ToString(data["default"], ""), bc.ToString(data["default"], ""), bc.ToString(data["value"], ""))
		},
		"values": func() bc.ClientComponent {
			toValue := func(lang, key string, locales bc.IM) string {
				if langValues, found := locales[lang].(bc.IM); found {
					return bc.ToString(langValues[key], "")
				}
				return ""
			}

			lang := bc.ToString(loc.Data["locales"], "")
			tag := bc.ToString(loc.Data["tag_keys"], "")
			rows := []bc.IM{}
			deflang := loc.Data["deflang"].(bc.IM)
			locales := loc.Data["locfile"].(bc.IM)["locales"].(bc.IM)
			fields := []mc.TableField{
				{Column: &mc.TableColumn{
					Id:     "tag",
					Header: loc.msg("locale_tag"),
					Cell: func(row bc.IM, col mc.TableColumn, value interface{}) string {
						linkLabel := fmt.Sprintf(
							`<span class="cell-label">%s</span>`, value)
						var link string
						link, _ = loc.getComponent("tag_cell", row)
						return linkLabel + link
					}}},
				{Name: "key", FieldType: mc.TableFieldTypeString, Label: loc.msg("locale_key")},
			}
			if lang == "default" {
				fields = append(fields, mc.TableField{
					Name: "default", FieldType: mc.TableFieldTypeString, Label: loc.msg("locale_value")})
			} else {
				fields = append(fields,
					mc.TableField{Column: &mc.TableColumn{
						Id:     "value",
						Header: loc.msg("locale_value"),
						Cell: func(row bc.IM, col mc.TableColumn, value interface{}) string {
							input, _ := loc.getComponent("value_cell", row)
							return input
						}}},
				)
			}
			var keys []string
			if loc.FilterValue != "" {
				keys = make([]string, 0)
				for rowKey := range deflang {
					if strings.Contains(rowKey, loc.FilterValue) ||
						strings.Contains(bc.ToString(deflang[rowKey], ""), loc.FilterValue) {
						keys = append(keys, rowKey)
					} else if lang != "default" {
						if strings.Contains(toValue(lang, rowKey, locales), loc.FilterValue) {
							keys = append(keys, rowKey)
						}
					}
				}
			} else if tag == "missing" {
				keys = make([]string, 0)
				for rowKey := range deflang {
					if _, found := locales[lang].(bc.IM)[rowKey]; !found {
						keys = append(keys, rowKey)
					}
				}
			} else {
				keys = loc.Data["tag_values"].(map[string][]string)[tag]
			}
			sort.Strings(keys)
			for _, key := range keys {
				if len(strings.Split(key, "_")) > 1 {
					rows = append(rows, bc.IM{
						"tag":     strings.Split(key, "_")[0],
						"key":     key,
						"default": bc.ToString(deflang[key], ""),
						"value":   toValue(lang, key, locales),
					})
				}
			}
			return ccTbl("key", rows, fields)
		},
	}
	cc := ccMap[name]()
	res, err = cc.Render()
	return res, err
}

func (loc *Locale) msg(labelID string) string {
	if label, found := loc.Labels[labelID]; found {
		return label
	}
	return labelID
}

func (loc *Locale) Render() (res string, err error) {
	loc.InitProps(loc)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(loc.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(loc.Class, " ")
		},
		"localeComponent": func(name string) (string, error) {
			return loc.getComponent(name, bc.IM{})
		},
		"lang": func() string {
			return bc.ToString(loc.Data["locales"], "")
		},
	}
	tpl := `<div id="{{ .Id }}" name="{{ .Name }}" class="row full {{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}
	><div class="row full section container-small">
	<div class="row full" >
	<div class="cell mobile" >
	<div class="cell mobile">
	<div class="cell padding-small" >{{ localeComponent "locales" }}</div>
	<div class="cell padding-small" >{{ localeComponent "tag_keys" }}</div>
	{{ if ne lang "default" }}<div class="cell padding-small" >{{ localeComponent "missing" }}</div>{{ end }}
	</div>
	<div class="cell padding-small mobile" >{{ localeComponent "filter" }}</div>
	</div>
	<div class="cell mobile" >
	<div class="right">
	{{ if eq .Dirty true }}
	<div class="cell padding-small" >{{ localeComponent "update" }}</div>
	<div class="cell padding-small" >{{ localeComponent "undo" }}</div>
	{{ end }}
	<div class="cell padding-small" >{{ localeComponent "add_item" }}</div>
	</div>
	</div>
	</div>
	{{ if eq .AddItem true }}
	<div class="row" >
	<div class="cell padding-small mobile" >{{ localeComponent "lang_key" }}</div>
	<div class="cell padding-small mobile" >{{ localeComponent "lang_name" }}</div>
	<div class="cell padding-small" >{{ localeComponent "add" }}</div>
	</div>
	{{ end }}
	</div>
	<div class="container section-small" >
	{{ localeComponent "values" }}
	</div>
	</div>`

	return bc.TemplateBuilder("locales", tpl, funcMap, loc)
}

func localeLocfile() bc.IM {
	return bc.IM{
		"locales": bc.IM{
			"de": bc.IM{
				"de":             "Deutsche",
				"key":            "de",
				"login_database": "Datenbank",
				"login_password": "Passwort",
				"login_username": "Nutzername",
			},
			"jp": bc.IM{
				"jp":             "日本語",
				"key":            "jp",
				"login_database": "データベース",
				"login_password": "パスワード",
				"login_username": "ユーザー名",
			},
		},
	}
}

var demoLocaleResponse func(evt bc.ResponseEvent) (re bc.ResponseEvent) = func(evt bc.ResponseEvent) (re bc.ResponseEvent) {
	switch evt.Name {
	case LocalesEventError:
		re = bc.ResponseEvent{
			Trigger: &fm.Toast{
				Type:  fm.ToastTypeError,
				Value: bc.ToString(evt.Value, ""),
			},
			TriggerName: evt.TriggerName,
			Name:        evt.Name,
			Header: bc.SM{
				bc.HeaderRetarget: "#toast-msg",
				bc.HeaderReswap:   "innerHTML",
			},
		}
		return re

	case LocalesEventSave:
		evt.Trigger.SetProperty("dirty", false)

	case LocalesEventUndo:
		evt.Trigger.SetProperty("data",
			bc.IM{"locfile": localeLocfile(),
				"locales": "default", "tag_keys": "address", "lang_key": "", "lang_name": "",
			})
		evt.Trigger.SetProperty("locales", []fm.SelectOption{
			{Value: "default", Text: "Default"},
			{Value: "de", Text: "Deutsch"},
			{Value: "jp", Text: "Japanese"},
		})
		evt.Trigger.SetProperty("filter_value", "")
		evt.Trigger.SetProperty("add_item", false)
		evt.Trigger.SetProperty("dirty", false)
	}
	return evt
}

func DemoLocale(demo bc.ClientComponent) []bc.DemoComponent {
	id := bc.ToString(demo.GetProperty("id"), "")
	eventURL := bc.ToString(demo.GetProperty("event_url"), "")
	requestValue := demo.GetProperty("request_value").(map[string]bc.IM)
	requestMap := demo.GetProperty("request_map").(map[string]bc.ClientComponent)
	return []bc.DemoComponent{
		{
			Label:         "Default",
			ComponentType: bc.ComponentTypeLocale,
			Component: &Locale{
				BaseComponent: bc.BaseComponent{
					Id:           id + "_locale_default",
					EventURL:     eventURL,
					OnResponse:   demoLocaleResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Data: bc.IM{
						"deflang": bc.IM{
							"key":               "en",
							"en":                "English",
							"address_view":      "Address Data",
							"address_country":   "Country",
							"address_state":     "State",
							"address_zipcode":   "Zipcode",
							"address_city":      "City",
							"address_street":    "Street",
							"address_notes":     "Comment",
							"login_username":    "Username",
							"login_password":    "Password",
							"login_database":    "Database",
							"login_lang":        "Language",
							"login_login":       "Login",
							"login_server":      "Server URL",
							"login_engine_err":  "Invalid database type!",
							"login_version_err": "Invalid service version!",
						},
						"locales":  "default",
						"tag_keys": "address",
						"tag_values": map[string][]string{
							"address": {
								"address_view", "address_country", "address_state", "address_zipcode",
								"address_city", "address_street", "address_notes",
							},
							"login": {
								"login_username", "login_password", "login_database", "login_lang",
								"login_login", "login_server", "login_engine_err", "login_version_err",
							},
						},
						"locfile": localeLocfile(),
					},
				},
				Locales: []fm.SelectOption{
					{Value: "default", Text: "Default"},
					{Value: "de", Text: "Deutsch"},
					{Value: "jp", Text: "Japanese"},
				},
				TagKeys: []fm.SelectOption{
					{Value: "address", Text: "address"},
					{Value: "login", Text: "login"},
				},
			}},
	}
}
