package page

import (
	"fmt"
	"strings"

	fm "github.com/nervatura/component/component/atom"
	bc "github.com/nervatura/component/component/base"
	mc "github.com/nervatura/component/component/molecule"
)

const (
	AdminEventModule        = "module"
	AdminEventTheme         = "theme"
	AdminEventChange        = "change"
	AdminEventCreate        = "create"
	AdminEventLogin         = "login"
	AdminEventPassword      = "password"
	AdminEventReportInstall = "report_install"
	AdminEventReportDelete  = "report_delete"
)

var adminDefaultLabel bc.SM = bc.SM{
	"admin_title":                   "Nervatura Admin",
	"admin_login":                   "Login",
	"admin_database":                "Database",
	"admin_client":                  "Client",
	"admin_locales":                 "Client languages",
	"admin_help":                    "Help",
	"admin_api_key":                 "API Key",
	"admin_alias":                   "Database alias name",
	"admin_demo":                    "Demo database",
	"admin_true":                    "TRUE",
	"admin_false":                   "FALSE",
	"admin_create":                  "Create a darabase",
	"admin_create_result_state":     "State",
	"admin_create_result_stamp":     "Time",
	"admin_create_result_message":   "Message",
	"admin_create_result_section":   "Section",
	"admin_create_result_datatype":  "Datatype",
	"admin_username":                "Username",
	"admin_password":                "Password",
	"admin_view_password":           "Password",
	"admin_view_report":             "Report",
	"admin_view_configuration":      "Configuration",
	"admin_view_logout":             "Logout",
	"admin_confirm":                 "Confirm",
	"admin_password_change":         "Password change",
	"admin_report_key":              "Report Key",
	"admin_report_install":          "INSTALL",
	"admin_report_delete":           "DELETE",
	"admin_report_list_reportkey":   "Report Key",
	"admin_report_list_installed":   "Installed",
	"admin_report_list_repname":     "Name",
	"admin_report_list_description": "Description",
	"admin_report_list_reptype":     "Type",
	"admin_report_list_filename":    "Filename",
	"admin_report_list_label":       "Label",
	"admin_env_list_key":            "Conf. Key",
	"admin_env_list_value":          "Conf. Value",
}

var adminIcoMap map[string][]string = map[string][]string{
	bc.ThemeDark: {bc.ThemeLight, "Sun"}, bc.ThemeLight: {bc.ThemeDark, "Moon"},
}

type Admin struct {
	bc.BaseComponent
	Version    string
	Theme      string
	Module     string
	View       string
	Token      string
	HelpURL    string
	ClientURL  string
	LocalesURL string
	Labels     bc.SM
	TokenLogin func(database, token string) bool // Token validation
}

func (adm *Admin) Properties() bc.IM {
	return bc.MergeIM(
		adm.BaseComponent.Properties(),
		bc.IM{
			"version":     adm.Version,
			"theme":       adm.Theme,
			"module":      adm.Module,
			"view":        adm.View,
			"token":       adm.Token,
			"help_url":    adm.HelpURL,
			"client_url":  adm.ClientURL,
			"locales_url": adm.LocalesURL,
			"labels":      adm.Labels,
		})
}

func (adm *Admin) GetProperty(propName string) interface{} {
	return adm.Properties()[propName]
}

func (adm *Admin) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"theme": func() interface{} {
			return adm.CheckEnumValue(bc.ToString(propValue, ""), bc.ThemeLight, bc.Theme)
		},
		"labels": func() interface{} {
			value := bc.SetSMValue(adm.Labels, "", "")
			if smap, valid := propValue.(bc.SM); valid {
				value = bc.MergeSM(value, smap)
			}
			if len(value) == 0 {
				value = adminDefaultLabel
			}
			return value
		},
		"target": func() interface{} {
			adm.SetProperty("id", adm.Id)
			value := bc.ToString(propValue, adm.Id)
			if value != "this" && !strings.HasPrefix(value, "#") {
				value = "#" + value
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if adm.BaseComponent.GetProperty(propName) != nil {
		return adm.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

func (adm *Admin) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"version": func() interface{} {
			adm.Version = bc.ToString(propValue, "1.0.0")
			return adm.Version
		},
		"theme": func() interface{} {
			adm.Theme = adm.Validation(propName, propValue).(string)
			return adm.Theme
		},
		"module": func() interface{} {
			adm.Module = bc.ToString(propValue, "database")
			return adm.Module
		},
		"view": func() interface{} {
			adm.View = bc.ToString(propValue, "password")
			return adm.View
		},
		"token": func() interface{} {
			adm.Token = bc.ToString(propValue, "")
			return adm.Token
		},
		"help_url": func() interface{} {
			adm.HelpURL = bc.ToString(propValue, "")
			return adm.HelpURL
		},
		"client_url": func() interface{} {
			adm.ClientURL = bc.ToString(propValue, "")
			return adm.ClientURL
		},
		"locales_url": func() interface{} {
			adm.LocalesURL = bc.ToString(propValue, "")
			return adm.LocalesURL
		},
		"labels": func() interface{} {
			adm.Labels = adm.Validation(propName, propValue).(bc.SM)
			return adm.Labels
		},
		"target": func() interface{} {
			adm.Target = adm.Validation(propName, propValue).(string)
			return adm.Target
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if adm.BaseComponent.GetProperty(propName) != nil {
		return adm.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

func (adm *Admin) response(evt bc.ResponseEvent) (re bc.ResponseEvent) {
	admEvt := bc.ResponseEvent{
		Trigger: adm, TriggerName: adm.Name, Value: evt.Value,
	}
	switch evt.TriggerName {

	case "api_key", "alias", "demo", "username", "password", "database", "confirm", "report_key":
		admEvt.Name = AdminEventChange
		adm.SetProperty("data", bc.IM{evt.TriggerName: admEvt.Value})

	case "theme":
		admEvt.Name = AdminEventTheme
		adm.SetProperty("theme", admEvt.Value)

	case "main_menu":
		admEvt.Name = AdminEventModule
		adm.SetProperty("module", admEvt.Value)
		if admEvt.Value == "help" && adm.HelpURL != "" {
			admEvt.Header = bc.MergeSM(admEvt.Header,
				bc.SM{bc.HeaderRedirect: adm.HelpURL})
		}
		if admEvt.Value == "client" && adm.ClientURL != "" {
			admEvt.Header = bc.MergeSM(admEvt.Header,
				bc.SM{bc.HeaderRedirect: adm.ClientURL})
		}
		if admEvt.Value == "locales" && adm.LocalesURL != "" {
			admEvt.Header = bc.MergeSM(admEvt.Header,
				bc.SM{bc.HeaderRedirect: adm.LocalesURL})
		}

	case "view_menu":
		admEvt.Name = AdminEventChange
		adm.SetProperty("view", admEvt.Value)
		if admEvt.Value == "logout" {
			adm.SetProperty("token", "")
		}

	case "create":
		admEvt.Name = AdminEventCreate

	case "login":
		admEvt.Name = AdminEventLogin

	case "report_install":
		admEvt.Name = AdminEventReportInstall

	case "report_delete":
		admEvt.Name = AdminEventReportDelete

	case "password_change":
		admEvt.Name = AdminEventPassword

	default:
	}
	if adm.OnResponse != nil {
		return adm.OnResponse(admEvt)
	}
	return admEvt
}

func (adm *Admin) getComponent(name string) (res string, err error) {
	ccLbl := func() *fm.Label {
		return &fm.Label{
			Value: adm.msg(name),
		}
	}
	ccSel := func(options []fm.SelectOption) *fm.Select {
		return &fm.Select{
			BaseComponent: bc.BaseComponent{
				Id: adm.Id + "_" + name, Name: name,
				EventURL:   adm.EventURL,
				Target:     adm.Target,
				OnResponse: adm.response,
			},
			Label:   adm.msg("admin_" + name),
			IsNull:  false,
			Value:   bc.ToString(adm.Data[name], ""),
			Options: options,
		}
	}
	ccInp := func(itype string) *fm.Input {
		return &fm.Input{
			BaseComponent: bc.BaseComponent{
				Id: adm.Id + "_" + name, Name: name,
				EventURL:   adm.EventURL,
				Target:     adm.Target,
				Swap:       bc.SwapOuterHTML,
				OnResponse: adm.response,
			},
			Type:  itype,
			Label: adm.msg("admin_" + name),
			Value: bc.ToString(adm.Data[name], ""),
			Full:  true,
		}
	}
	ccBtn := func(btnType, label string, disabled bool) *fm.Button {
		return &fm.Button{
			BaseComponent: bc.BaseComponent{
				Id: adm.Id + "_" + name, Name: name,
				EventURL:   adm.EventURL,
				Target:     adm.Target,
				OnResponse: adm.response,
			},
			Type:     btnType,
			Label:    label,
			Disabled: disabled,
		}
	}
	ccMenu := func(items []mc.MenuBarItem, value, class string) *mc.MenuBar {
		return &mc.MenuBar{
			BaseComponent: bc.BaseComponent{
				Id: adm.Id + "_" + name, Name: name,
				EventURL:   adm.EventURL,
				Target:     adm.Target,
				OnResponse: adm.response,
				Class:      []string{class},
			},
			Items:   items,
			Value:   value,
			SideBar: false,
		}
	}
	ccTbl := func(rowKey string, rows []bc.IM, fields []mc.TableField) *mc.Table {
		return &mc.Table{
			BaseComponent: bc.BaseComponent{
				Id: adm.Id + "_" + name, Name: name,
				EventURL: adm.EventURL,
			},
			Rows:        rows,
			Fields:      fields,
			Pagination:  mc.PaginationTypeTop,
			PageSize:    10,
			RowKey:      rowKey,
			TableFilter: true,
			AddItem:     false,
		}
	}
	ccMap := map[string]func() bc.ClientComponent{
		"main_menu": func() bc.ClientComponent {
			return ccMenu(
				[]mc.MenuBarItem{
					{Value: "database", Label: adm.msg("admin_database"), Icon: "Database"},
					{Value: "login", Label: adm.msg("admin_login"), Icon: "Edit"},
					{Value: "client", Label: adm.msg("admin_client"), Icon: "Globe"},
					{Value: "locales", Label: adm.msg("admin_locales"), Icon: "User"},
					{Value: "help", Label: adm.msg("admin_help"), Icon: "QuestionCircle"},
				},
				bc.ToString(adm.GetProperty("module"), ""), "border-top")
		},
		"view_menu": func() bc.ClientComponent {
			return ccMenu(
				[]mc.MenuBarItem{
					{Value: "password", Label: adm.msg("admin_view_password"), Icon: "Key"},
					{Value: "report", Label: adm.msg("admin_view_report"), Icon: "ChartBar"},
					{Value: "configuration", Label: adm.msg("admin_view_configuration"), Icon: "Cog"},
					{Value: "logout", Label: adm.msg("admin_view_logout"), Icon: "Exit"},
				},
				bc.ToString(adm.GetProperty("view"), ""), "border-top")
		},
		"admin_api_key": func() bc.ClientComponent {
			return ccLbl()
		},
		"api_key": func() bc.ClientComponent {
			return ccInp(fm.InputTypeText)
		},
		"admin_alias": func() bc.ClientComponent {
			return ccLbl()
		},
		"alias": func() bc.ClientComponent {
			return ccInp(fm.InputTypeText)
		},
		"admin_demo": func() bc.ClientComponent {
			return ccLbl()
		},
		"demo": func() bc.ClientComponent {
			return ccSel([]fm.SelectOption{
				{Value: "true", Text: adm.msg("admin_true")},
				{Value: "false", Text: adm.msg("admin_false")},
			})
		},
		"create": func() bc.ClientComponent {
			disabled := (bc.ToString(adm.Data["alias"], "") == "") || (bc.ToString(adm.Data["api_key"], "") == "")
			return ccBtn(fm.ButtonTypePrimary, adm.msg("admin_"+name), disabled)
		},
		"theme": func() bc.ClientComponent {
			themeBtn := ccBtn(fm.ButtonTypePrimary, "", false)
			themeBtn.Style = bc.SM{"padding": "4px"}
			themeBtn.Value = adminIcoMap[adm.Theme][0]
			themeBtn.LabelComponent = &fm.Icon{Value: adminIcoMap[adm.Theme][1], Width: 18, Height: 18}
			return themeBtn
		},
		"create_result": func() bc.ClientComponent {
			fields := []mc.TableField{
				{Column: &mc.TableColumn{
					Id:        "state",
					Header:    adm.msg("admin_create_result_state"),
					CellStyle: bc.SM{"text-align": "center"},
					Cell: func(row bc.IM, col mc.TableColumn, value interface{}) string {
						icoKey := "InfoCircle"
						color := "orange"
						if value == "err" {
							icoKey = "ExclamationTriangle"
							color = "red"
						}
						res, _ := (&fm.Icon{Value: icoKey, Color: color}).Render()
						return fmt.Sprintf(
							`<span class="cell-label">%s</span>%s`, col.Header, res)
					}}},
				{Name: "stamp", FieldType: mc.TableFieldTypeTime, Label: adm.msg("admin_create_result_stamp")},
				{Name: "section", FieldType: mc.TableFieldTypeString, Label: adm.msg("admin_create_result_section")},
				{Name: "datatype", FieldType: mc.TableFieldTypeString, Label: adm.msg("admin_create_result_datatype")},
				{Column: &mc.TableColumn{
					Id:     "message",
					Header: adm.msg("admin_create_result_message"),
					Cell: func(row bc.IM, col mc.TableColumn, value interface{}) string {
						style := ""
						if row["state"] == "err" {
							style = `style="color:red;"`
						}
						return fmt.Sprintf(
							`<span class="cell-label">%s</span>
							<span %s >%s</span>`, col.Header, style, bc.ToString(value, ""))
					}}},
			}
			rows := []bc.IM{}

			if result, valid := adm.Data["create_result"].([]bc.IM); valid {
				rows = result
			}
			return ccTbl("stamp", rows, fields)
		},
		"admin_username": func() bc.ClientComponent {
			return ccLbl()
		},
		"username": func() bc.ClientComponent {
			return ccInp(fm.InputTypeText)
		},
		"admin_password": func() bc.ClientComponent {
			return ccLbl()
		},
		"password": func() bc.ClientComponent {
			return ccInp(fm.InputTypePassword)
		},
		"admin_confirm": func() bc.ClientComponent {
			return ccLbl()
		},
		"confirm": func() bc.ClientComponent {
			return ccInp(fm.InputTypePassword)
		},
		"admin_database": func() bc.ClientComponent {
			return ccLbl()
		},
		"database": func() bc.ClientComponent {
			return ccInp(fm.InputTypeText)
		},
		"login": func() bc.ClientComponent {
			disabled := (bc.ToString(adm.Data["username"], "") == "") || (bc.ToString(adm.Data["database"], "") == "")
			return ccBtn(fm.ButtonTypePrimary, adm.msg("admin_"+name), disabled)
		},
		"password_change": func() bc.ClientComponent {
			disabled := (bc.ToString(adm.Data["username"], "") == "") || (bc.ToString(adm.Data["password"], "") == "") || (bc.ToString(adm.Data["confirm"], "") == "")
			return ccBtn(fm.ButtonTypePrimary, adm.msg("admin_"+name), disabled)
		},
		"admin_report_key": func() bc.ClientComponent {
			return ccLbl()
		},
		"report_key": func() bc.ClientComponent {
			return ccInp(fm.InputTypeText)
		},
		"report_install": func() bc.ClientComponent {
			disabled := (bc.ToString(adm.Data["report_key"], "") == "")
			return ccBtn(fm.ButtonTypePrimary, adm.msg("admin_"+name), disabled)
		},
		"report_delete": func() bc.ClientComponent {
			disabled := (bc.ToString(adm.Data["report_key"], "") == "")
			return ccBtn(fm.ButtonTypePrimary, adm.msg("admin_"+name), disabled)
		},
		"report_list": func() bc.ClientComponent {
			fields := []mc.TableField{
				{Name: "reportkey", FieldType: mc.TableFieldTypeString, Label: adm.msg("admin_report_list_reportkey")},
				{Column: &mc.TableColumn{
					Id:        "installed",
					Header:    adm.msg("admin_report_list_installed"),
					CellStyle: bc.SM{"text-align": "center"},
					Cell: func(row bc.IM, col mc.TableColumn, value interface{}) string {
						icoKey := "Times"
						color := "red"
						if bc.ToBoolean(value, false) {
							icoKey = "Check"
							color = "green"
						}
						res, _ := (&fm.Icon{Value: icoKey, Color: color}).Render()
						return fmt.Sprintf(
							`<span class="cell-label">%s</span>%s`, col.Header, res)
					}}},
				{Name: "repname", FieldType: mc.TableFieldTypeString, Label: adm.msg("admin_report_list_repname")},
				{Name: "description", FieldType: mc.TableFieldTypeString, Label: adm.msg("admin_report_list_description")},
				{Name: "reptype", FieldType: mc.TableFieldTypeString, Label: adm.msg("admin_report_list_reptype")},
				{Name: "filename", FieldType: mc.TableFieldTypeString, Label: adm.msg("admin_report_list_filename")},
				{Name: "label", FieldType: mc.TableFieldTypeString, Label: adm.msg("admin_report_list_label")},
			}
			rows := []bc.IM{}

			if result, valid := adm.Data["report_list"].([]bc.IM); valid {
				rows = result
			}
			return ccTbl("reportkey", rows, fields)
		},
		"env_list": func() bc.ClientComponent {
			fields := []mc.TableField{
				{Name: "key", FieldType: mc.TableFieldTypeString, Label: adm.msg("admin_env_list_key")},
				{Name: "value", FieldType: mc.TableFieldTypeString, Label: adm.msg("admin_env_list_value")},
			}
			rows := []bc.IM{}

			if envList, valid := adm.Data["env_list"].([]bc.SM); valid {
				for key, value := range envList {
					rows = append(rows, bc.IM{"key": key, "value": value})
				}
			}
			return ccTbl("key", rows, fields)
		},
	}
	cc := ccMap[name]()
	res, err = cc.Render()
	if err == nil {
		adm.RequestMap = bc.MergeCM(adm.RequestMap, cc.GetProperty("request_map").(map[string]bc.ClientComponent))
	}
	return res, err
}

func (adm *Admin) InitProps() {
	for key, value := range adm.Properties() {
		adm.SetProperty(key, value)
	}
}

func (adm *Admin) msg(labelID string) string {
	if label, found := adm.Labels[labelID]; found {
		return label
	}
	return labelID
}

func (adm *Admin) Render() (res string, err error) {
	adm.InitProps()

	funcMap := map[string]any{
		"msg": func(labelID string) string {
			return adm.msg(labelID)
		},
		"styleMap": func() bool {
			return len(adm.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(adm.Class, " ")
		},
		"adminComponent": func(name string) (string, error) {
			return adm.getComponent(name)
		},
		"showResult": func() bool {
			if _, valid := adm.Data["create_result"].([]bc.IM); valid {
				return true
			}
			return false
		},
		"tokenLogin": func() bool {
			if adm.TokenLogin != nil {
				database := bc.ToString(adm.Data["database"], "")
				return adm.TokenLogin(database, adm.Token)
			}
			return false
		},
	}
	tpl := `<div id="{{ .Id }}" theme="{{ .Theme }}" class="admin row mobile {{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}>
	<div class="row title">
	<div class="cell">
	<div class="cell title-cell" ><span>{{ msg "admin_title" }}</span></div>
	<div class="cell">{{ adminComponent "theme" }}</div>
	</div>
	<div class="cell version-cell" ><span>{{ .Version }}</span></div>
	</div>
	{{ adminComponent "main_menu" }}
	{{ if eq .Module "database" }}
	<div class="row full section-small" >
	<div class="row full container-small section-small-bottom" >
	<div class="cell mobile">
	<div class="cell padding-small mobile" >
	<div class="section-small" >{{ adminComponent "admin_api_key" }}</div>
	{{ adminComponent "api_key" }}
	</div>
	<div class="cell padding-small mobile" >
	<div class="section-small" >{{ adminComponent "admin_alias" }}</div>
	{{ adminComponent "alias" }}
	</div>
	<div class="cell padding-small mobile" >
	<div class="section-small" >{{ adminComponent "admin_demo" }}</div>
	{{ adminComponent "demo" }}
	</div>
	</div></div>
	<div class="row full section-small" >
	<div class="cell container center" >
	{{ adminComponent "create" }}
	</div>
	</div>
	{{ if showResult }}
	<div class="container section-small" >
	{{ adminComponent "create_result" }}
	</div>
	{{ end }}
	</div>
	{{ end }}
	{{ if eq .Module "login" }}
	{{ if eq tokenLogin false }}
	<div class="row full section-small" >
	<div class="row full container-small section-small-bottom" >
	<div class="cell padding-small mobile">
	<div class="section-small" >{{ adminComponent "admin_username" }}</div>
	{{ adminComponent "username" }}
	</div>
	<div class="cell padding-small mobile">
	<div class="section-small" >{{ adminComponent "admin_password" }}</div>
	{{ adminComponent "password" }}
	</div>
	<div class="cell padding-small mobile">
	<div class="section-small" >{{ adminComponent "admin_database" }}</div>
	{{ adminComponent "database" }}
	</div>
	</div></div>
	<div class="row full section-small-bottom" >
	<div class="cell container center section-small-bottom" >
	{{ adminComponent "login" }}
	</div>
	</div>
	{{ else }}
	{{ adminComponent "view_menu" }}
	{{ if eq .View "password" }}
	<div class="row full section-small" >
	<div class="row full container-small section-small-bottom" >
	<div class="cell padding-small mobile">
	<div class="section-small" >{{ adminComponent "admin_username" }}</div>
	{{ adminComponent "username" }}
	</div>
	<div class="cell padding-small mobile">
	<div class="section-small" >{{ adminComponent "admin_password" }}</div>
	{{ adminComponent "password" }}
	</div>
	<div class="cell padding-small mobile">
	<div class="section-small" >{{ adminComponent "admin_confirm" }}</div>
	{{ adminComponent "confirm" }}
	</div>
	</div></div>
	<div class="row full section-small-bottom" >
	<div class="cell container center section-small-bottom" >
	{{ adminComponent "password_change" }}
	</div>
	</div>
	{{ end }}
	{{ if eq .View "report" }}
	<div class="row full section-small" >
	<div class="row full padding-small mobile" >
	<div class="cell padding-small mobile" >
	{{ adminComponent "admin_report_key" }}
	</div>
	<div class="cell padding-small mobile" >
	{{ adminComponent "report_key" }}
	</div>
	<div class="cell mobile" >
	<div class="cell padding-small" >
	{{ adminComponent "report_install" }}
	</div>
	<div class="cell padding-small" >
	{{ adminComponent "report_delete" }}
	</div>
	</div>
	</div>
	<div class="container section-small" >
	{{ adminComponent "report_list" }}
	</div>
	</div>
	{{ end }}
	{{ if eq .View "configuration" }}
	<div class="row full section-small" >
	<div class="container section-small" >
	{{ adminComponent "env_list" }}
	</div>
	</div>
	{{ end }}
	{{ end }}
	{{ end }}
	</div>`

	return bc.TemplateBuilder("admin", tpl, funcMap, adm)
}

var testCreateResult []bc.IM = []bc.IM{
	{"database": "demo", "message": "Start process", "stamp": "2023-12-22 17:03:26", "state": "create"},
	{"message": "The existing table is dropped...", "stamp": "2023-12-22 17:03:26", "state": "create"},
	{"message": "Creating the tables...", "stamp": "2023-12-22 17:03:26", "state": "create"},
	{"message": "Creating indexes ...", "stamp": "2023-12-22 17:03:26", "state": "create"},
	{"message": "Data initialization ...", "stamp": "2023-12-22 17:03:26", "state": "err"},
	{"message": "Rebuilding the database...", "stamp": "2023-12-22 17:03:26", "state": "create"},
	{"result": "csv_custpos_en,csv_vat_en,ntr_bank_en,ntr_cash_in_en,ntr_cash_out_en,ntr_customer_en,ntr_custpos_en,ntr_delivery_in_en,ntr_delivery_out_en,ntr_delivery_transfer_en,ntr_employee_en,ntr_formula_en,ntr_inventory_en,ntr_invoice_en,ntr_offer_in_en,ntr_offer_out_en,ntr_order_in_en,ntr_order_out_en,ntr_product_en,ntr_production_en,ntr_project_en,ntr_receipt_en,ntr_rental_in_en,ntr_rental_out_en,ntr_tool_en,ntr_vat_en,ntr_waybill_in_en,ntr_waybill_out_en,ntr_worksheet_en,sample", "section": "report templates", "stamp": "2023-12-22 17:03:26", "state": "demo"},
	{"datatype": "groups", "result": "138,139,140,141,142,143", "stamp": "2023-12-22 17:03:26", "state": "demo"},
	{"datatype": "deffield", "result": "47,48,49,50,51", "section": "customer", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "customer", "result": "2,3,4", "section": "customer", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "address", "result": "3,4,5,6", "section": "customer", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "contact", "result": "2,3,4,5", "section": "customer", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "event", "result": "1,2,3,4", "section": "customer", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "deffield", "result": "52", "section": "employee", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "employee", "result": "2,3,4,5", "section": "employee", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "address", "result": "7,8,9,10", "section": "employee", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "contact", "result": "6,7,8,9", "section": "employee", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "event", "result": "5", "section": "employee", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "deffield", "result": "53,54,55", "section": "product", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "product", "result": "1,2,3,4,5,6,7,8,9,10,11,12,13", "section": "product", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "barcode", "result": "1,2,3", "section": "product", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "price", "result": "1,2", "section": "product", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "event", "result": "6,7,8", "section": "product", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "deffield", "result": "56,57", "section": "project", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "project", "result": "1", "section": "project", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "address", "result": "11", "section": "project", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "contact", "result": "10", "section": "project", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "event", "result": "9", "section": "project", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "deffield", "result": "58,59", "section": "tool", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "tool", "result": "1,2,3", "section": "tool", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "event", "result": "10", "section": "tool", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "place", "result": "4", "stamp": "2023-12-22 17:03:27", "state": "demo"},
	{"datatype": "trans", "result": "1,2,3,4,5,6,7,8,9", "section": "document(offer,order,invoice,rent,worksheet)", "stamp": "2023-12-22 17:03:28", "state": "demo"},
	{"datatype": "item", "result": "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33", "stamp": "2023-12-22 17:03:28", "state": "demo"},
	{"datatype": "link", "result": "1,2,3,4", "stamp": "2023-12-22 17:03:28", "state": "demo"},
	{"datatype": "trans", "result": "10,11", "section": "payment(bank,petty cash)", "stamp": "2023-12-22 17:03:28", "state": "demo"},
	{"datatype": "payment", "result": "1,2,3,4", "stamp": "2023-12-22 17:03:28", "state": "demo"},
	{"datatype": "link", "result": "5,6,7,8", "stamp": "2023-12-22 17:03:28", "state": "demo"},
	{"datatype": "trans", "result": "12,13,14,15,16,17,18,19,20", "section": "pstock control(tool movement,delivery,stock transfer,correction,formula,production)", "stamp": "2023-12-22 17:03:28", "state": "demo"},
	{"datatype": "movement", "result": "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44", "stamp": "2023-12-22 17:03:29", "state": "demo"},
	{"datatype": "link", "result": "9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28", "stamp": "2023-12-22 17:03:29", "state": "demo"},
	{"datatype": "ui_menu", "result": "1,2", "section": "sample menus", "stamp": "2023-12-22 17:03:29", "state": "demo"},
	{"datatype": "ui_menufields", "result": "1,2,3", "stamp": "2023-12-22 17:03:29", "state": "demo"},
	{"message": "The database was created successfully!", "stamp": "2023-12-22 17:03:29", "state": "log"},
}

var testReportList []bc.IM = []bc.IM{
	{"description": "Accounts Payable and Receivable.", "filename": "csv_custpos_en.json", "installed": true, "label": "", "repname": "Payments Due List - CSV output", "reportkey": "csv_custpos_en", "reptype": "csv"},
	{"description": "Recoverable and payable VAT summary grouped by currency.", "filename": "csv_vat_en.json", "installed": true, "label": "", "repname": "VAT Summary - CSV output.", "reportkey": "csv_vat_en", "reptype": "csv"},
	{"description": "Bank Statement", "filename": "ntr_bank_en.json", "installed": false, "label": "bank", "repname": "Bank Statement EN", "reportkey": "ntr_bank_en", "reptype": "pdf"},
	{"description": "Deposit Voucher", "filename": "ntr_cash_in_en.json", "installed": true, "label": "cash", "repname": "Petty Cash Voucher EN", "reportkey": "ntr_cash_in_en", "reptype": "pdf"},
	{"description": "Payment Voucher", "filename": "ntr_cash_out_en.json", "installed": true, "label": "cash", "repname": "Petty Cash Voucher EN", "reportkey": "ntr_cash_out_en", "reptype": "pdf"},
	{"description": "Customer Information Sheet", "filename": "ntr_customer_en.json", "installed": true, "label": "", "repname": "Customer Sheet", "reportkey": "ntr_customer_en", "reptype": "pdf"},
	{"description": "Accounts Payable and Receivable", "filename": "ntr_custpos_en.json", "installed": true, "label": "", "repname": "Payments Due List", "reportkey": "ntr_custpos_en", "reptype": "pdf"},
	{"description": "Goods Received Note", "filename": "ntr_delivery_in_en.json", "installed": false, "label": "delivery", "repname": "Goods Received Note EN", "reportkey": "ntr_delivery_in_en", "reptype": "pdf"},
	{"description": "Delivery Note", "filename": "ntr_delivery_out_en.json", "installed": true, "label": "delivery", "repname": "Delivery Note EN", "reportkey": "ntr_delivery_out_en", "reptype": "pdf"},
	{"description": "Stock Transfer", "filename": "ntr_delivery_transfer_en.json", "installed": true, "label": "delivery", "repname": "Stock Transfer EN", "reportkey": "ntr_delivery_transfer_en", "reptype": "pdf"},
	{"description": "Employee Information Sheet", "filename": "ntr_employee_en.json", "installed": true, "label": "", "repname": "Employee Sheet", "reportkey": "ntr_employee_en", "reptype": "pdf"},
	{"description": "Formula Sheet", "filename": "ntr_formula_en.json", "installed": true, "label": "formula", "repname": "Formula Sheet EN", "reportkey": "ntr_formula_en", "reptype": "pdf"},
	{"description": "Inventory Control", "filename": "ntr_inventory_en.json", "installed": true, "label": "inventory", "repname": "Inventory Control EN", "reportkey": "ntr_inventory_en", "reptype": "pdf"},
	{"description": "Customer invoice", "filename": "ntr_invoice_en.json", "installed": true, "label": "invoice", "repname": "Invoice EN", "reportkey": "ntr_invoice_en", "reptype": "pdf"},
	{"description": "Supplier Offer", "filename": "ntr_offer_in_en.json", "installed": true, "label": "offer", "repname": "Request EN", "reportkey": "ntr_offer_in_en", "reptype": "pdf"},
	{"description": "Customer Offer", "filename": "ntr_offer_out_en.json", "installed": true, "label": "offer", "repname": "Offer EN", "reportkey": "ntr_offer_out_en", "reptype": "pdf"},
	{"description": "Supplier Order", "filename": "ntr_order_in_en.json", "installed": true, "label": "order", "repname": "Order EN", "reportkey": "ntr_order_in_en", "reptype": "pdf"},
	{"description": "Customer Order", "filename": "ntr_order_out_en.json", "installed": true, "label": "order", "repname": "Order EN", "reportkey": "ntr_order_out_en", "reptype": "pdf"},
	{"description": "Product Information Sheet", "filename": "ntr_product_en.json", "installed": true, "label": "", "repname": "Product Sheet", "reportkey": "ntr_product_en", "reptype": "pdf"},
	{"description": "Production Sheet", "filename": "ntr_production_en.json", "installed": true, "label": "production", "repname": "Production Sheet EN", "reportkey": "ntr_production_en", "reptype": "pdf"},
	{"description": "Project Information Sheet", "filename": "ntr_project_en.json", "installed": true, "label": "", "repname": "Project Sheet", "reportkey": "ntr_project_en", "reptype": "pdf"},
	{"description": "Receipt", "filename": "ntr_receipt_en.json", "installed": true, "label": "receipt", "repname": "Receipt EN", "reportkey": "ntr_receipt_en", "reptype": "pdf"},
	{"description": "Supplier Rental", "filename": "ntr_rental_in_en.json", "installed": true, "label": "rent", "repname": "Rental Data EN", "reportkey": "ntr_rental_in_en", "reptype": "pdf"},
	{"description": "Customer Rental", "filename": "ntr_rental_out_en.json", "installed": true, "label": "rent", "repname": "Rental Data EN", "reportkey": "ntr_rental_out_en", "reptype": "pdf"},
	{"description": "Tool Information Sheet", "filename": "ntr_tool_en.json", "installed": true, "label": "", "repname": "Tool Sheet", "reportkey": "ntr_tool_en", "reptype": "pdf"},
	{"description": "Recoverable and payable VAT summary grouped by currency.", "filename": "ntr_vat_en.json", "installed": true, "label": "", "repname": "VAT Summary", "reportkey": "ntr_vat_en", "reptype": "pdf"},
	{"description": "Incoming Movement", "filename": "ntr_waybill_in_en.json", "installed": true, "label": "waybill", "repname": "Tool Movement EN", "reportkey": "ntr_waybill_in_en", "reptype": "pdf"},
	{"description": "Outgoing Movement", "filename": "ntr_waybill_out_en.json", "installed": true, "label": "waybill", "repname": "Tool Movement EN", "reportkey": "ntr_waybill_out_en", "reptype": "pdf"},
	{"description": "Worksheet data", "filename": "ntr_worksheet_en.json", "installed": true, "label": "worksheet", "repname": "Worksheet EN", "reportkey": "ntr_worksheet_en", "reptype": "pdf"},
	{"description": "Test PDF", "filename": "sample.json", "installed": true, "label": "", "repname": "Test PDF", "reportkey": "sample", "reptype": "pdf"},
}

var testEnvList bc.SM = bc.SM{
	"NT_ALIAS_DEFAULT": "", "NT_ALIAS_DEMO": "sqlite://file:data/demo.db?cache=shared&mode=rwc",
	"NT_API_KEY": "EXAMPLE_API_KEY", "NT_APP_LOG_FILE": "", "NT_CLIENT_CONFIG": "data/client_config_loc.json",
	"NT_CORS_ALLOW_CREDENTIALS": "false", "NT_CORS_ALLOW_HEADERS": "", "NT_CORS_ALLOW_METHODS": "",
	"NT_CORS_ALLOW_ORIGINS": "", "NT_CORS_ENABLED": "true", "NT_CORS_EXPOSE_HEADERS": "", "NT_CORS_MAX_AGE": "0",
	"NT_DOCS_URL": "https://nervatura.github.io/nervatura/", "NT_FONT_DIR": "", "NT_FONT_FAMILY": "",
	"NT_GRPC_ENABLED": "true", "NT_GRPC_PORT": "9200", "NT_GRPC_TLS_ENABLED": "false",
	"NT_HASHTABLE": "ref17890714", "NT_HTTP_ENABLED": "true", "NT_HTTP_HOME": "/admin", "NT_HTTP_LOG_FILE": "",
	"NT_HTTP_PORT": "5000", "NT_HTTP_READ_TIMEOUT": "30", "NT_HTTP_TLS_ENABLED": "false", "NT_HTTP_WRITE_TIMEOUT": "30",
	"NT_PASSWORD_LOGIN": "true", "NT_REPORT_DIR": "",
	"NT_SECURITY_ALLOWED_HOSTS": "", "NT_SECURITY_ALLOWED_HOSTS_ARE_REGEX": "false", "NT_SECURITY_BROWSER_XSS_FILTER": "false",
	"NT_SECURITY_CONTENT_SECURITY_POLICY": "", "NT_SECURITY_CONTENT_TYPE_NOSNIFF": "false", "NT_SECURITY_CUSTOM_FRAME_OPTIONS_VALUE": "",
	"NT_SECURITY_DEVELOPMENT": "false", "NT_SECURITY_ENABLED": "false", "NT_SECURITY_EXPECT_CT_HEADER": "",
	"NT_SECURITY_FEATURE_POLICY": "", "NT_SECURITY_FORCE_STS_HEADER": "false", "NT_SECURITY_FRAME_DENY": "false",
	"NT_SECURITY_HOSTS_PROXY_HEADERS": "", "NT_SECURITY_PROXY_HEADERS": "", "NT_SECURITY_PUBLIC_KEY": "",
	"NT_SECURITY_REFERRER_POLICY": "", "NT_SECURITY_SSL_HOST": "", "NT_SECURITY_SSL_REDIRECT": "false",
	"NT_SECURITY_SSL_TEMPORARY_REDIRECT": "false", "NT_SECURITY_STS_INCLUDE_SUBDOMAINS": "false",
	"NT_SECURITY_STS_PRELOAD": "false", "NT_SECURITY_STS_SECONDS": "0",
	"NT_SMTP_HOST": "", "NT_SMTP_PASSWORD": "", "NT_SMTP_PORT": "465", "NT_SMTP_TLS_MIN_VERSION": "0", "NT_SMTP_USER": "",
	"NT_TLS_CERT_FILE": "", "NT_TLS_KEY_FILE": "",
	"NT_TOKEN_EXP": "6", "NT_TOKEN_ISS": "nervatura", "NT_TOKEN_PRIVATE_KEY": "cshl03Pncui68cGCFiOiO6kTtrSXvQza",
	"NT_TOKEN_PRIVATE_KID": "af566fd801211238a181c131a0725336157b6d5049a28ab9857151af843bdebf",
	"NT_TOKEN_PUBLIC_KEY":  "", "NT_TOKEN_PUBLIC_KEY_URL": "", "NT_TOKEN_PUBLIC_KID": "PUBLIC_KID",
	"SQL_CONN_MAX_LIFETIME": "15", "SQL_MAX_IDLE_CONNS": "3", "SQL_MAX_OPEN_CONNS": "10", "VERSION": "dev",
}

var demoAdminResponse func(evt bc.ResponseEvent) (re bc.ResponseEvent) = func(evt bc.ResponseEvent) (re bc.ResponseEvent) {
	toastMsg := func(value string) (re bc.ResponseEvent) {
		return bc.ResponseEvent{
			Trigger: &fm.Toast{
				Type:    fm.ToastTypeInfo,
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
	}
	switch evt.Name {
	case AdminEventModule:
		if (evt.Value == "client" || evt.Value == "help") && (len(evt.Header) == 0) {
			return toastMsg(bc.ToString(evt.Value, ""))
		}
	case AdminEventCreate, AdminEventLogin, AdminEventPassword, AdminEventReportInstall, AdminEventReportDelete:
		return toastMsg(evt.Name)
	}
	return evt
}

func DemoAdmin(eventURL, parentID string) []bc.DemoComponent {
	return []bc.DemoComponent{
		{
			Label:         "Default",
			ComponentType: bc.ComponentTypeAdmin,
			Component: &Admin{
				BaseComponent: bc.BaseComponent{
					Id:         bc.GetComponentID(),
					EventURL:   eventURL,
					OnResponse: demoAdminResponse,
					Data: bc.IM{
						"alias": "demo",
					},
				},
				Module:     "database",
				HelpURL:    "https://nervatura.github.io/nervatura/",
				ClientURL:  "/client",
				LocalesURL: "/locales",
			}},
		{
			Label:         "Database result",
			ComponentType: bc.ComponentTypeAdmin,
			Component: &Admin{
				BaseComponent: bc.BaseComponent{
					Id:         bc.GetComponentID(),
					EventURL:   eventURL,
					OnResponse: demoAdminResponse,
					Data: bc.IM{
						"api_key":       "DEMO_API_KEY",
						"alias":         "demo",
						"create_result": testCreateResult,
					},
				},
				Module: "database",
			}},
		{
			Label:         "Login",
			ComponentType: bc.ComponentTypeAdmin,
			Component: &Admin{
				BaseComponent: bc.BaseComponent{
					Id:         bc.GetComponentID(),
					EventURL:   eventURL,
					OnResponse: demoAdminResponse,
					Data: bc.IM{
						"username": "admin",
						"database": "demo",
					},
				},
				Module: "login",
			}},
		{
			Label:         "Password change",
			ComponentType: bc.ComponentTypeAdmin,
			Component: &Admin{
				BaseComponent: bc.BaseComponent{
					Id:         bc.GetComponentID(),
					EventURL:   eventURL,
					OnResponse: demoAdminResponse,
					Data: bc.IM{
						"username": "admin",
						"database": "demo",
					},
				},
				Module:     "login",
				Token:      "TOKEN0123456789",
				TokenLogin: func(database, token string) bool { return (token != "") },
				View:       "password",
			}},
		{
			Label:         "Report",
			ComponentType: bc.ComponentTypeAdmin,
			Component: &Admin{
				BaseComponent: bc.BaseComponent{
					Id:         bc.GetComponentID(),
					EventURL:   eventURL,
					OnResponse: demoAdminResponse,
					Data: bc.IM{
						"username":    "admin",
						"database":    "demo",
						"report_list": testReportList,
					},
				},
				Module:     "login",
				Token:      "TOKEN0123456789",
				TokenLogin: func(database, token string) bool { return (token != "") },
				View:       "report",
			}},
		{
			Label:         "Configuration",
			ComponentType: bc.ComponentTypeAdmin,
			Component: &Admin{
				BaseComponent: bc.BaseComponent{
					Id:         bc.GetComponentID(),
					EventURL:   eventURL,
					OnResponse: demoAdminResponse,
					Data: bc.IM{
						"username": "admin",
						"database": "demo",
						"env_list": testEnvList,
					},
				},
				Module:     "login",
				Token:      "TOKEN0123456789",
				TokenLogin: func(database, token string) bool { return (token != "") },
				View:       "configuration",
			}},
	}
}
