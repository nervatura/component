package component

import (
	"fmt"
	"html/template"
	"slices"
	"strings"
	"time"

	ut "github.com/nervatura/component/pkg/util"
)

const (
	ComponentTypeClient = "client"

	ClientEventModule   = "client_module"
	ClientEventTheme    = "client_theme"
	ClientEventLogOut   = "client_logout"
	ClientEventSide     = "client_side"
	ClientEventSideMenu = "client_side_menu"
	ClientEventForm     = "client_form"
)

var ClientIcoMap map[string][]string = map[string][]string{
	ThemeDark: {ThemeLight, "Sun"}, ThemeLight: {ThemeDark, "Moon"},
}

/*
The [Client] Login Ticket struct represents a user authentication ticket with various properties.
*/
type Ticket struct {
	// Login session ID.
	SessionID string `json:"session_id"`
	// Login frontend client host address. Optional.
	Host string `json:"host"`
	// Login authentication method. Example: password, google, facebook, etc. Optional.
	AuthMethod string `json:"auth_method"`
	// Database is the login database name of the ticket. Optional.
	Database string `json:"database"`
	// User is the user information of the ticket.
	User ut.IM `json:"user,omitempty"`
	// Expiry is the expiration time of the ticket.
	Expiry time.Time `json:"expiry,omitempty"`
}

// expired reports whether the ticket is expired.
func (t *Ticket) expired() bool {
	if t.Expiry.IsZero() {
		return false
	}
	return t.Expiry.Round(0).Before(time.Now())
}

/*
Checking the validity of the login ticket. The ticket is valid if has a SessionID and User, and its validity time is less
than the current time or zero. The ticket is checked every time the Client component is displayed. In case of an invalid ticket,
the login form is automatically displayed. The LoginDisabled option disables the verification.
*/
func (t *Ticket) Valid() bool {
	return t != nil && t.SessionID != "" && t.User != nil && !t.expired()
}

/*
The Client component is a main application component that can be used to implement all the main functions of a
typical client application. It allows you to use the following components: [Login], [MenuBar], [SideBar], [Search],
[Browser], [Editor], modal and simple [Form].
*/
type Client struct {
	BaseComponent
	// Application version value
	Version string `json:"version"`
	/*
		The theme of the control.
		[Theme] variable constants: [ThemeLight], [ThemeDark]. Default value: [ThemeLight]
	*/
	Theme string `json:"theme"`
	// Current ui language
	Lang string `json:"lang"`
	// User authentication ticket
	Ticket Ticket `json:"ticket"`
	/*
		Specifies whether the login verification is disabled. By default, it will display the
		login form based on the Validation function of the login [Ticket]. Default value: false
	*/
	LoginDisabled bool `json:"login_disabled"`
	//	The redirect url of the login page. Default value: "/"
	LoginURL string `json:"login_url"`
	/*
		The buttons of the login page.
	*/
	LoginButtons []LoginAuthButton `json:"login_buttons"`
	/*
		Specifies whether the side bar is hidden. Default value: false
	*/
	HideSideBar bool `json:"hide_side_bar"`
	/*
		The visibility of the side bar if HideSideBar is false.
		[SideBarVisibility] variable constants: [SideBarVisibilityAuto], [SideBarVisibilityShow], [SideBarVisibilityHide].
		Default value: [SideBarVisibilityAuto]
	*/
	SideBarVisibility string `json:"sidebar_visibility"`
	// Specifies whether the main menu is hidden. Default value: false
	HideMenu bool `json:"hide_menu"`
	// Custom UI and any message text.
	ClientLabels func(lang string) ut.SM `json:"-"`
	/* Custom main menu. The function is called before each display of the Client component if it also affects the
	display of the [MenuBar] component.
	*/
	ClientMenu func(labels ut.SM, config ut.IM) MenuBar `json:"-"`
	/* Custom side bar. The function is called before each display of the Client component if it also affects the
	display of the [SideBar] component.
	*/
	ClientSideBar func(
		moduleKey string, labels ut.SM, data ut.IM) SideBar `json:"-"`
	/* Custom login form. The function is called before each display of the Client component if it also affects the
	display of the [Login] component.
	*/
	ClientLogin func(labels ut.SM, config ut.IM) Login `json:"-"`
	/* Custom search form. The function is called before each display of the Client component if it also affects the
	display of the [Search] component.
	*/
	ClientSearch func(viewName string, labels ut.SM, searchData ut.IM) Search `json:"-"`
	/* Custom advanced search component. The function is called before each display of the Client component if it also affects the
	display of the [Browser] component.
	*/
	ClientBrowser func(viewName string, labels ut.SM, searchData ut.IM) Browser `json:"-"`
	/* Custom editor component. The function is called before each display of the Client component if it also affects the
	display of the [Editor] component.
	*/
	ClientEditor func(editorKey, viewName string, labels ut.SM, editorData ut.IM) Editor `json:"-"`
	/* Custom modal form. The function is called before each display of the Client component if it also affects the
	display of the modal [Form] component.
	*/
	ClientModalForm func(formKey string, labels ut.SM, data ut.IM) Form `json:"-"`
	/* Custom simple form. The function is called before each display of the Client component if it also affects the
	display of the [Form] component.
	*/
	ClientForm func(editorKey, formKey string, labels ut.SM, data ut.IM) Form `json:"-"`
}

/*
Returns all properties of the [Client]
*/
func (cli *Client) Properties() ut.IM {
	return ut.MergeIM(
		cli.BaseComponent.Properties(),
		ut.IM{
			"version":            cli.Version,
			"ticket":             cli.Ticket,
			"theme":              cli.Theme,
			"lang":               cli.Lang,
			"sidebar_visibility": cli.SideBarVisibility,
			"login_disabled":     cli.LoginDisabled,
			"login_url":          cli.LoginURL,
			"login_buttons":      cli.LoginButtons,
			"hide_side_bar":      cli.HideSideBar,
			"hide_menu":          cli.HideMenu,
			"client_labels":      cli.ClientLabels,
			"client_menu":        cli.ClientMenu,
			"client_side_bar":    cli.ClientSideBar,
			"client_login":       cli.ClientLogin,
			"client_search":      cli.ClientSearch,
			"client_browser":     cli.ClientBrowser,
			"client_editor":      cli.ClientEditor,
			"client_modal_form":  cli.ClientModalForm,
			"client_form":        cli.ClientForm,
		})
}

/*
Returns the value of the property of the [Client] with the specified name.
*/
func (cli *Client) GetProperty(propName string) interface{} {
	return cli.Properties()[propName]
}

/*
It checks the value given to the property of the [Client] and always returns a valid value
*/
func (cli *Client) Validation(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"ticket": func() interface{} {
			if value, valid := propValue.(Ticket); valid {
				return value
			}
			return Ticket{}
		},
		"theme": func() interface{} {
			return cli.CheckEnumValue(ut.ToString(propValue, ""), ThemeLight, Theme)
		},
		"sidebar_visibility": func() interface{} {
			return cli.CheckEnumValue(ut.ToString(propValue, ""), SideBarVisibilityAuto, SideBarVisibility)
		},
		"login_buttons": func() interface{} {
			value := []LoginAuthButton{}
			if buttons, valid := propValue.([]LoginAuthButton); valid {
				value = buttons
			}
			return value
		},
		"target": func() interface{} {
			cli.SetProperty("id", cli.Id)
			value := ut.ToString(propValue, cli.Id)
			if value != "this" && !strings.HasPrefix(value, "#") {
				value = "#" + value
			}
			return value
		},
	}
	if _, found := pm[propName]; found {
		return pm[propName]()
	}
	if cli.BaseComponent.GetProperty(propName) != nil {
		return cli.BaseComponent.Validation(propName, propValue)
	}
	return propValue
}

/*
Setting a property of the [Client] value safely. Checks the entered value.
In case of an invalid value, the default value will be set.
*/
func (cli *Client) SetProperty(propName string, propValue interface{}) interface{} {
	pm := map[string]func() interface{}{
		"version": func() interface{} {
			cli.Version = ut.ToString(propValue, "1.0.0")
			return cli.Version
		},
		"theme": func() interface{} {
			cli.Theme = cli.Validation(propName, propValue).(string)
			return cli.Theme
		},
		"ticket": func() interface{} {
			cli.Ticket = cli.Validation(propName, propValue).(Ticket)
			return cli.Ticket
		},
		"lang": func() interface{} {
			cli.Lang = ut.ToString(propValue, "en")
			return cli.Lang
		},
		"sidebar_visibility": func() interface{} {
			cli.SideBarVisibility = cli.Validation(propName, propValue).(string)
			return cli.SideBarVisibility
		},
		"login_disabled": func() interface{} {
			cli.LoginDisabled = ut.ToBoolean(propValue, false)
			return cli.LoginDisabled
		},
		"login_url": func() interface{} {
			cli.LoginURL = ut.ToString(propValue, "/")
			return cli.LoginURL
		},
		"login_buttons": func() interface{} {
			cli.LoginButtons = cli.Validation(propName, propValue).([]LoginAuthButton)
			return cli.LoginButtons
		},
		"hide_side_bar": func() interface{} {
			cli.HideSideBar = ut.ToBoolean(propValue, false)
			return cli.HideSideBar
		},
		"hide_menu": func() interface{} {
			cli.HideMenu = ut.ToBoolean(propValue, false)
			return cli.HideMenu
		},
		"target": func() interface{} {
			cli.Target = cli.Validation(propName, propValue).(string)
			return cli.Target
		},
	}
	if _, found := pm[propName]; found {
		return cli.SetRequestValue(propName, pm[propName](), []string{})
	}
	if cli.BaseComponent.GetProperty(propName) != nil {
		return cli.BaseComponent.SetProperty(propName, propValue)
	}
	return propValue
}

/*
If the OnResponse function of the [Client] is implemented, the function calls it after the [TriggerEvent]
is processed, otherwise the function's return [ResponseEvent] is the processed [TriggerEvent].
*/
func (cli *Client) OnRequest(te TriggerEvent) (re ResponseEvent) {
	if cc, found := cli.RequestMap[te.Id]; found {
		return cc.OnRequest(te)
	}
	re = ResponseEvent{
		Trigger:     &BaseComponent{},
		TriggerName: te.Name,
		Name:        te.Name,
		Header: ut.SM{
			HeaderReswap:   SwapNone,
			HeaderRedirect: cli.LoginURL,
		},
	}
	return re
}

func (cli *Client) responseBrowser(evt ResponseEvent) (re ResponseEvent) {
	re = ResponseEvent{
		Trigger: cli, TriggerName: cli.Name, Value: evt.Value,
	}
	searchData := cli.getDataIM("search", cli.Data)
	searchView := cli.getDataIM(ut.ToString(searchData["view"], ""), searchData)
	if !slices.Contains([]string{
		BrowserEventSearch, BrowserEventSetColumn,
		BrowserEventAddFilter, BrowserEventChangeFilter, BrowserEventEditRow,
	}, evt.Name) {
		return evt
	}
	switch evt.Name {
	case BrowserEventSetColumn:
		searchView["visible_columns"] = evt.Trigger.GetProperty("visible_columns")
	case BrowserEventAddFilter, BrowserEventChangeFilter:
		searchView["filters"] = evt.Trigger.GetProperty("filters")
	default:
		re.Value = evt.Value
	}
	cli.SetProperty("data", ut.MergeIM(cli.Data, ut.IM{"search": searchData}))
	re.Name = evt.Name
	re.Header = ut.SM{
		HeaderRetarget: "#" + cli.Id,
	}
	if cli.OnResponse != nil {
		return cli.OnResponse(re)
	}
	return re
}

func (cli *Client) responseLogin(evt ResponseEvent) (re ResponseEvent) {
	re = ResponseEvent{
		Trigger: cli, TriggerName: cli.Name, Value: evt.Value,
	}
	re.Name = evt.Name
	re.Value = evt.Trigger.GetProperty("data")
	re.Header = ut.SM{
		HeaderRetarget: "#" + cli.Id,
	}
	if evt.Name == LoginEventLang {
		cli.SetProperty("lang", evt.Value)
	}
	if evt.Name == LoginEventTheme {
		re.Name = ClientEventTheme
		cli.SetProperty("theme", ClientIcoMap[cli.Theme][0])
	}
	if evt.Name == LoginEventAuth {
		re.Value = evt.Value
	}
	if evt.Name == LoginEventLogin && !ut.ToBoolean(evt.Trigger.GetProperty("hide_database"), false) {
		if values, found := re.Value.(ut.IM); found && ut.ToString(values["database"], "") != cli.Ticket.Database {
			cli.Ticket.Database = ut.ToString(values["database"], "")
		}
	}
	return re
}

func (cli *Client) responseMainMenu(evt ResponseEvent) (re ResponseEvent) {
	re = ResponseEvent{
		Trigger: cli, TriggerName: cli.Name, Value: evt.Value,
	}
	if evt.Name == MenuBarEventSide {
		sdValues := ut.SM{
			SideBarVisibilityAuto: SideBarVisibilityShow,
			SideBarVisibilityShow: SideBarVisibilityAuto,
			SideBarVisibilityHide: SideBarVisibilityShow,
		}
		re.Name = ClientEventSide
		cli.SetProperty("sidebar_visibility", sdValues[cli.SideBarVisibility])
	} else {
		switch evt.Value {
		case "theme":
			re.Name = ClientEventTheme
			cli.SetProperty("theme", ClientIcoMap[cli.Theme][0])
		case "logout":
			re.Name = ClientEventLogOut
			cli.SetProperty("token", "")
		default:
			re.Name = ClientEventModule
		}
	}
	return re
}

func (cli *Client) response(evt ResponseEvent) (re ResponseEvent) {
	admEvt := ResponseEvent{
		Trigger: cli, TriggerName: cli.Name, Value: evt.Value,
	}
	switch evt.TriggerName {

	case "modal":
		admEvt.Name = evt.Name
		data := cli.GetProperty("data").(ut.IM)
		delete(data, "modal")
		cli.SetProperty("data", data)

	case "table", "filter_table", "view_table":
		if !slices.Contains([]string{
			TableEventEditCell, TableEventAddItem,
		}, evt.Name) {
			return evt
		}
		admEvt.Value = evt.Value
		admEvt.Name = evt.Name
		admEvt.Header = ut.SM{
			HeaderRetarget: "#" + cli.Id,
		}

	case "search":
		admEvt.Name = evt.Name
		admEvt.Header = ut.SM{
			HeaderRetarget: "#" + cli.Id,
		}

	case "browser":
		return cli.responseBrowser(evt)

	case "editor":
		if evt.Name == EditorEventView {
			editorData := ut.ToIM(cli.Data["editor"], ut.IM{})
			editorData["view"] = evt.Value
			cli.SetProperty("view", evt.Value)
			cli.SetProperty("data", ut.MergeIM(cli.Data, ut.IM{"editor": editorData}))

		}
		admEvt.Name = evt.Name
		admEvt.Header = ut.SM{
			HeaderRetarget: "#" + cli.Id,
		}

	case "form":
		admEvt.Name = ClientEventForm
		evt.Value = ut.MergeIM(ut.ToIM(evt.Value, ut.IM{}), ut.IM{"event": evt.Name})
		admEvt.Header = ut.SM{
			HeaderRetarget: "#" + cli.Id,
		}
		if evt.Name != FormEventChange {
			editorData := ut.ToIM(cli.Data["editor"], ut.IM{})
			delete(editorData, "form")
			cli.SetProperty("data", ut.MergeIM(cli.Data, ut.IM{"editor": editorData}))
		}

	case "login":
		admEvt = cli.responseLogin(evt)

	case "side_menu":
		admEvt.Name = ClientEventSideMenu

	case "main_menu":
		admEvt = cli.responseMainMenu(evt)

	default:
	}
	if cli.OnResponse != nil {
		return cli.OnResponse(admEvt)
	}
	return admEvt
}

func (cli *Client) getDataIM(key string, data ut.IM) ut.IM {
	if _, found := data[key].(ut.IM); !found {
		data[key] = ut.IM{}
	}
	return data[key].(ut.IM)

}

/*
The GetSearchVisibleColumns function retrieves the visible columns from the search data.
*/
func (cli *Client) GetSearchVisibleColumns(icols map[string]bool) (cols map[string]bool) {
	searchData := cli.getDataIM("search", cli.Data)
	searchView := cli.getDataIM(ut.ToString(searchData["view"], ""), searchData)
	cols = map[string]bool{}
	if bcols, found := searchView["visible_columns"].(map[string]bool); found {
		return bcols
	}
	if vcols, found := searchView["visible_columns"].(ut.IM); found {
		return ut.ToBoolMap(vcols, map[string]bool{})
	}
	for key, ivalue := range icols {
		cols[key] = ut.ToBoolean(ivalue, false)
	}
	return cols
}

/*
The Labels function retrieves the labels that are used in the client.
*/
func (cli *Client) Labels() ut.SM {
	if cli.ClientLabels != nil {
		return cli.ClientLabels(cli.Lang)
	}
	return ut.SM{}
}

/*
The GetSearchFilters function retrieves the filters from the search data.
*/
func (cli *Client) GetSearchFilters(value string, cfFilters interface{}) (filters []BrowserFilter) {
	searchData := cli.getDataIM("search", cli.Data)
	searchView := cli.getDataIM(ut.ToString(searchData["view"], ""), searchData)
	filters = []BrowserFilter{}
	if dfilters, found := searchView["filters"].([]BrowserFilter); found {
		return dfilters
	}
	toFilters := func(dfilters []interface{}) (f []BrowserFilter) {
		for _, filter := range dfilters {
			if filterMap, valid := filter.(ut.IM); valid {
				filters = append(filters, BrowserFilter{
					Or:    ut.ToBoolean(filterMap["or"], false),
					Field: ut.ToString(filterMap["field"], ""),
					Comp:  ut.ToString(filterMap["comp"], ""),
					Value: ut.ToString(filterMap["value"], value),
				})
			}
		}
		return filters
	}
	if dfilters, found := searchView["filters"].([]interface{}); found {
		return toFilters(dfilters)
	}
	if dfilters, found := cfFilters.([]interface{}); found && len(dfilters) > 0 {
		return toFilters(dfilters)
	}
	if dfilters, found := cfFilters.([]BrowserFilter); found && len(dfilters) > 0 {
		return dfilters
	}
	return filters
}

func (cli *Client) getComponent(name string) (html template.HTML, err error) {
	labels := cli.Labels()
	state, stateKey, stateData := cli.GetStateData()
	config := ut.MergeIM(ut.ToIM(cli.Data["config"], ut.IM{}),
		ut.IM{
			"version": cli.Version, "theme": cli.Theme, "lang": cli.Lang, "ticket": cli.Ticket,
			"login_disabled": cli.LoginDisabled, "hide_menu": cli.HideMenu,
			"hide_side_bar": cli.HideSideBar, "sidebar_visibility": cli.SideBarVisibility,
		})
	ccBase := func(data ut.IM) BaseComponent {
		return BaseComponent{
			Id:           cli.Id + "_" + name,
			Name:         name,
			EventURL:     cli.EventURL,
			Target:       cli.Target,
			OnResponse:   cli.response,
			RequestValue: cli.RequestValue,
			RequestMap:   cli.RequestMap,
			Data:         data,
		}
	}
	ccMap := map[string]func() ClientComponent{
		"main_menu": func() ClientComponent {
			mnu := MenuBar{}
			if cli.ClientMenu != nil {
				mnu = cli.ClientMenu(labels, config)
			}
			mnu.BaseComponent = ccBase(mnu.Data)
			mnu.SetProperty("side_bar", !cli.HideSideBar)
			mnu.SetProperty("sidebar_visibility", cli.SideBarVisibility)
			if state == "search" || state == "browser" {
				mnu.SetProperty("value", "search")
			} else {
				mnu.SetProperty("value", stateKey)
			}
			return &mnu
		},
		"side_menu": func() ClientComponent {
			moduleKey := stateKey
			if state == "search" || state == "browser" {
				moduleKey = state
			}
			if state == "form" {
				moduleKey = ut.ToString(stateData["key"], "")
			}
			sb := SideBar{}
			if cli.ClientSideBar != nil {
				sb = cli.ClientSideBar(moduleKey, labels,
					ut.MergeIM(stateData, ut.IM{"config": config}))
			}
			sb.BaseComponent = ccBase(sb.Data)
			sb.SetProperty("visibility", cli.SideBarVisibility)
			return &sb
		},
		"login": func() ClientComponent {
			lgn := Login{}
			if cli.ClientLogin != nil {
				lgn = cli.ClientLogin(labels, config)
			}
			lgn.BaseComponent = ccBase(lgn.Data)
			lgn.SetProperty("auth_buttons", cli.LoginButtons)
			if cli.Ticket.Database != "" {
				lgn.SetProperty("data", ut.MergeIM(
					ut.ToIM(lgn.GetProperty("data"), ut.IM{}), ut.IM{
						"database": cli.Ticket.Database,
					}))
			}
			lgn.SetProperty("theme", cli.GetProperty("theme"))
			return &lgn
		},
		"search": func() ClientComponent {
			sea := Search{}
			if cli.ClientSearch != nil {
				sea = cli.ClientSearch(stateKey, labels,
					ut.MergeIM(stateData, ut.IM{"config": config}))
			}
			sea.BaseComponent = ccBase(sea.Data)
			sea.SetProperty("data", ut.IM{"filter_value": ut.ToString(stateData["filter_value"], "")})
			return &sea
		},
		"browser": func() ClientComponent {
			bro := Browser{}
			if cli.ClientBrowser != nil {
				bro = cli.ClientBrowser(ut.ToString(stateData["view"], ""), labels,
					ut.MergeIM(stateData, ut.IM{"config": config}))
			}
			bro.BaseComponent = ccBase(bro.Data)
			bro.SetProperty("filters", cli.GetSearchFilters("", bro.Filters))
			bro.SetProperty("visible_columns", cli.GetSearchVisibleColumns(bro.VisibleColumns))
			return &bro
		},
		"editor": func() ClientComponent {
			edi := Editor{}
			if cli.ClientEditor != nil {
				edi = cli.ClientEditor(stateKey, ut.ToString(stateData["view"], ""), labels,
					ut.MergeIM(stateData, ut.IM{"config": config}))
			}
			edi.BaseComponent = ccBase(edi.Data)
			return &edi
		},
		"modal": func() ClientComponent {
			frm := Form{}
			modalData := ut.ToIM(cli.Data["modal"], ut.IM{})
			if cli.ClientModalForm != nil {
				frm = cli.ClientModalForm(ut.ToString(modalData["key"], ""), labels,
					ut.MergeIM(ut.ToIM(modalData["data"], ut.IM{}), ut.IM{"config": config}))
			}
			frm.BaseComponent = ccBase(frm.Data)
			frm.SetProperty("modal", true)
			frm.SetProperty("data", ut.ToIM(modalData["data"], ut.IM{}))
			return &frm
		},
		"form": func() ClientComponent {
			formData := ut.ToIM(stateData["form"], ut.IM{})
			frm := Form{}
			if cli.ClientForm != nil {
				frm = cli.ClientForm(
					ut.ToString(stateData["key"], ""), stateKey,
					labels, ut.MergeIM(ut.ToIM(formData["data"], ut.IM{}), ut.IM{"config": config}))
			}
			frm.BaseComponent = ccBase(frm.Data)
			frm.SetProperty("data", stateData)
			return &frm
		},
	}
	cc := ccMap[name]()
	html, err = cc.Render()
	return html, err
}

/*
The CleanComponent function cleans the request value and request map for a given component name.
*/
func (cli *Client) CleanComponent(name string) {
	for key := range cli.RequestValue {
		if strings.HasPrefix(key, cli.Id+"_"+name) {
			delete(cli.RequestValue, key)
			delete(cli.RequestMap, key)
		}
	}
}

/*
The Msg function retrieves a label from the client's labels.
*/
func (cli *Client) Msg(labelID string) string {
	if label, found := cli.Labels()[labelID]; found {
		return label
	}
	return labelID
}

/*
The GetStateData function retrieves the state, value, and data from the client's data.
*/
func (cli *Client) GetStateData() (state, value string, data ut.IM) {
	cliData := ut.ToIM(cli.GetProperty("data"), ut.IM{})
	if editor, found := cliData["editor"].(ut.IM); found {
		if form, found := editor["form"].(ut.IM); found {
			return "form", ut.ToString(form["key"], ""), editor
		}
		return "editor", ut.ToString(editor["key"], ""), editor
	}
	search := ut.ToIM(cliData["search"], ut.IM{})
	if ut.ToBoolean(search["simple"], false) {
		return "search", ut.ToString(search["view"], ""), search
	}
	return "browser", ut.ToString(search["view"], ""), search
}

/*
The SetConfigValue function sets a value in the client's config.
*/
func (cli *Client) SetConfigValue(key string, value interface{}) {
	cli.Data["config"] = ut.MergeIM(ut.ToIM(cli.Data["config"], ut.IM{}), ut.IM{key: value})
	cli.SetProperty("data", cli.Data)
}

/*
The SetEditor function sets the editor data and view state.
*/
func (cli *Client) SetEditor(editorKey, viewName string, data ut.IM) {
	editorData := ut.MergeIM(data, ut.IM{"key": editorKey, "view": viewName})
	cli.Data["editor"] = editorData
	cli.SetProperty("data", ut.MergeIM(cli.Data, ut.IM{"editor": editorData}))
	cli.CleanComponent("login")
}

/*
The ResetEditor function clean the editor data and restores the last search state.
*/
func (cli *Client) ResetEditor() {
	delete(cli.Data, "editor")
	cli.SetProperty("data", cli.Data)
	cli.SetProperty("sidebar_visibility", SideBarVisibilityAuto)
	cli.CleanComponent("editor")
	cli.CleanComponent("login")
}

/*
The SetSearch function sets the search data and view state.
*/
func (cli *Client) SetSearch(viewName string, data ut.IM, simple bool) {
	cli.Data["search"] = ut.MergeIM(data, ut.IM{"view": viewName, "simple": simple})
	delete(cli.Data, "editor")
	cli.SetProperty("data", cli.Data)
	cli.SetProperty("sidebar_visibility", SideBarVisibilityAuto)
	cli.CleanComponent("editor")
	cli.CleanComponent("login")
}

/*
The SetForm function sets the form or modal data and view state.
*/
func (cli *Client) SetForm(formKey string, data ut.IM, index int64, modal bool) {
	if modal {
		cli.Data["modal"] = ut.IM{"key": formKey, "data": data}
	} else {
		editorData := ut.ToIM(cli.Data["editor"], ut.IM{})
		delete(editorData, "form")
		cli.Data["editor"] = ut.MergeIM(editorData, ut.IM{"form": ut.IM{"key": formKey, "data": data, "index": index}})
	}
	cli.SetProperty("data", cli.Data)
}

/*
Based on the values, it will generate the html code of the [Client] or return with an error message.
*/
func (cli *Client) Render() (html template.HTML, err error) {
	cli.InitProps(cli)

	funcMap := map[string]any{
		"styleMap": func() bool {
			return len(cli.Style) > 0
		},
		"customClass": func() string {
			return strings.Join(cli.Class, " ")
		},
		"clientComponent": func(name string) (template.HTML, error) {
			return cli.getComponent(name)
		},
		"validTicket": func() bool {
			fmt.Println("validTicket", cli.Ticket.Valid())
			return cli.Ticket.Valid()
		},
		"clientState": func() string {
			state, _, _ := cli.GetStateData()
			return state
		},
		"modalForm": func() bool {
			_, found := cli.Data["modal"].(ut.IM)
			return found
		},
	}
	tpl := `<div id="{{ .Id }}" theme="{{ .Theme }}" class="client {{ customClass }}"
	{{ if styleMap }} style="{{ range $key, $value := .Style }}{{ $key }}:{{ $value }};{{ end }}"{{ end }}>
	{{ if and (eq validTicket false) (eq .LoginDisabled false) }}
	{{ clientComponent "login" }}
	{{ end }}
	{{ if or (validTicket) (.LoginDisabled) }}
	{{ if eq .HideMenu false }}<div class="client-menubar">{{ clientComponent "main_menu" }}</div>{{ end }}
	<div theme="{{ .Theme }}" class="main" {{ if eq .HideMenu true }}style="top: 0;"{{ end }}>
  {{ if eq .HideSideBar false }}{{ clientComponent "side_menu" }}{{ end }}
	<div class="page" {{ if eq .HideSideBar true }}style="margin-left: 0;"{{ end }}>{{ clientComponent clientState }}</div>
	</div>
	{{ end }}
	{{ if modalForm }}{{ clientComponent "modal" }}{{ end }}
	</div>`

	return ut.TemplateBuilder("client", tpl, funcMap, cli)
}

var testClientLabels func(lang string) ut.SM = func(lang string) ut.SM {
	labelMap := map[string]ut.SM{
		"en": {
			"theme_dark": "Dark", "theme_light": "Light",
			"mnu_menu": "Menu", "mnu_hide": "Hide", "mnu_search": "Search",
			"mnu_setting": "Setting", "mnu_help": "Help", "mnu_logout": "Logout",
			"title_login":          "Demo Client",
			"login_username":       "Username",
			"login_password":       "Password",
			"login_database":       "Database",
			"login_lang":           "Language",
			"login_login":          "Login",
			"login_theme":          "Theme",
			"login_help":           "Help",
			"mnu_info":             "Info",
			"mnu_search_simple":    "Simple Search",
			"mnu_search_browser":   "Browser Search",
			"editor_save":          "Save",
			"editor_delete":        "Delete",
			"editor_cancel":        "Cancel",
			"customer_new":         "New Customer",
			"browser_title":        "Data browser",
			"browser_view":         "Data view",
			"browser_views":        "Views",
			"browser_columns":      "Columns",
			"browser_total":        "Total",
			"browser_filter":       "Filter",
			"browser_comparison":   "?",
			"browser_filters":      "Filter criterias",
			"browser_search":       "Search",
			"browser_bookmark":     "Bookmark",
			"browser_export":       "Export",
			"browser_help":         "Help",
			"browser_result":       "record(s) found",
			"browser_placeholder":  "Filter",
			"browser_label_new":    "NEW",
			"browser_label_delete": "Delete",
			"browser_label_yes":    "YES",
			"browser_label_no":     "NO",
			"browser_label_ok":     "OK",
			"browser_export_error": "There is too much data, it cannot be exported. Use the Filters to limit the displayed result!",
			"editor_title":         "Demo Editor",
			"settings_title":       "Settings",
			"info_title":           "Info",
		},
		"zh": {
			"theme_dark": "暗黑", "theme_light": "明亮",
			"mnu_menu": "菜单", "mnu_hide": "隐藏", "mnu_search": "搜索",
			"mnu_setting": "设置", "mnu_help": "帮助", "mnu_logout": "退出",
			"title_login":    "Demo Client",
			"login_username": "用户名",
			"login_password": "密码",
			"login_database": "数据库",
			"login_lang":     "语言",
			"login_login":    "登录",
		},
	}
	return labelMap[lang]
}

func testClientLogin(labels ut.SM, config ut.IM) Login {
	theme := ut.ToString(config["theme"], "light")
	version := ut.ToString(config["version"], "1.0.0")
	lang := ut.ToString(config["lang"], "en")
	login := Login{
		Locales: []SelectOption{
			{Value: "en", Text: "English"},
			{Value: "zh", Text: "Chinese"},
		},
		AuthButtons:  []LoginAuthButton{},
		Version:      version,
		Theme:        theme,
		Labels:       labels,
		Lang:         lang,
		HideDatabase: false,
		HidePassword: false,
		ShowHelp:     false,
		HelpURL:      "",
	}
	login.SetProperty("data", ut.IM{
		"username": "admin", "database": "demo",
	})
	return login
}

func testClientMenu(labels ut.SM, config ut.IM) MenuBar {
	theme := ut.ToString(config["theme"], "light")
	hideExit := ut.ToBoolean(config["login_disabled"], false)
	mnu := MenuBar{
		Items: []MenuBarItem{
			{Value: "theme", Label: labels["theme_"+ClientIcoMap[theme][0]], Icon: ClientIcoMap[theme][1]},
			{Value: "search", Label: labels["mnu_search"], Icon: IconSearch},
			{Value: "setting", Label: labels["mnu_setting"], Icon: IconCog},
			{Value: "info", Label: labels["mnu_info"], Icon: IconInfoCircle},
		},
		LabelMenu: labels["mnu_menu"],
		LabelHide: labels["mnu_hide"],
		SideBar:   true,
	}
	if !hideExit {
		mnu.Items = append(mnu.Items, MenuBarItem{
			Value: "logout", Label: labels["mnu_logout"], Icon: IconExit,
		})
	}
	return mnu
}

func testClientSideBar(moduleKey string, labels ut.SM, data ut.IM) SideBar {
	searchSideBar := func(labels ut.SM, data ut.IM) []SideBarItem {
		return []SideBarItem{
			&SideBarSeparator{},
			&SideBarElement{
				Name:     "customer_simple",
				Value:    "customer_simple",
				Label:    labels["mnu_search_simple"],
				Icon:     IconBolt,
				Selected: (ut.ToString(data["view"], "") == "customer_simple"),
			},
			&SideBarSeparator{},
			&SideBarElement{
				Name:     "customer_browser",
				Value:    "customer_browser",
				Label:    labels["mnu_search_browser"],
				Icon:     IconSearch,
				Selected: (ut.ToString(data["view"], "") == "customer_browser"),
			},
		}
	}
	itemMap := map[string]func() []SideBarItem{
		"search": func() []SideBarItem {
			return searchSideBar(labels, data)
		},
		"browser": func() []SideBarItem {
			return searchSideBar(labels, data)
		},
		"customer": func() []SideBarItem {
			return []SideBarItem{
				&SideBarSeparator{},
				&SideBarElement{
					Name:    "editor_cancel",
					Value:   "editor_cancel",
					Label:   labels["browser_title"],
					Icon:    IconReply,
					NotFull: true,
				},
				&SideBarSeparator{},
				&SideBarSeparator{},
				&SideBarElement{
					Name:  "editor_save",
					Value: "editor_save",
					Label: labels["editor_save"],
					Icon:  IconUpload,
				},
				&SideBarElement{
					Name:  "editor_delete",
					Value: "editor_delete",
					Label: labels["editor_delete"],
					Icon:  IconTimes,
				},
				&SideBarSeparator{},
				&SideBarElement{
					Name:     "editor_new",
					Value:    "editor_new",
					Label:    labels["customer_new"],
					Icon:     IconUser,
					Disabled: true,
				},
			}
		},
	}
	sb := SideBar{
		Items: []SideBarItem{},
	}
	if item, found := itemMap[moduleKey]; found {
		sb.Items = item()
	}
	return sb
}

func testClientSearch(viewName string, labels ut.SM, searchData ut.IM) Search {
	search := Search{
		Fields:            testSearchFields,
		Title:             labels["mnu_search_simple"],
		FilterPlaceholder: "Customer Name, Number, City, Steet...",
		AutoFocus:         true,
		Full:              true,
		Rows:              testSearchRows,
	}
	return search
}

func testClientBrowser(viewName string, labels ut.SM, searchData ut.IM) Browser {
	bro := Browser{
		Table: Table{
			TableFilter:       true,
			HidePaginatonSize: false,
			Fields:            testBrowserFields["customer"](),
			Rows:              testBrowserRows["customer"](),
			AddItem:           true,
			PageSize:          10,
		},
		Title:        labels["mnu_search_browser"],
		View:         viewName,
		HideBookmark: true,
		HideHelp:     true,
		ExportLimit:  65000,
		Labels:       labels,
		Download:     fmt.Sprintf("%s.csv", viewName),
		ReadOnly:     false,
	}
	bro.SetProperty("visible_columns", testBrowserColumns["customer"]())
	bro.SetProperty("hide_filters", map[string]bool{})
	bro.SetProperty("filters", testBrowserFilters["customer"]())
	return bro
}

func testClientEditor(editorKey, viewName string, labels ut.SM, editorData ut.IM) Editor {
	edi := Editor{
		Title:  labels["editor_title"],
		Icon:   IconEdit,
		Views:  testEditorViews(),
		Rows:   testEditorRows(viewName),
		Tables: testEditorTable(viewName),
	}
	edi.SetProperty("view", viewName)
	return edi
}

func testClientForm(editorKey, formKey string, labels ut.SM, data ut.IM) (form Form) {
	return Form{
		Title:      labels["settings_title"],
		BodyRows:   testFormBodyRows("multiple_input"),
		FooterRows: testFormFooterRows("multiple_input"),
		Icon:       IconEdit,
	}
}

func testClientModalForm(formKey string, labels ut.SM, data ut.IM) (form Form) {
	return Form{
		Title:      labels["info_title"],
		BodyRows:   testFormBodyRows("info"),
		FooterRows: testFormFooterRows("info"),
		Icon:       IconInfoCircle,
	}
}

var testClientResponse func(evt ResponseEvent) (re ResponseEvent) = func(evt ResponseEvent) (re ResponseEvent) {
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
	client := evt.Trigger.(*Client)
	switch evt.Name {
	case ClientEventTheme, ClientEventSide, LoginEventLang,
		BrowserEventChangeFilter, BrowserEventAddFilter, BrowserEventSetColumn,
		EditorEventView, FormEventOK, FormEventCancel:
		return evt
	case LoginEventLogin:
		values := ut.ToIM(evt.Value, ut.IM{})
		client.Ticket = Ticket{
			SessionID:  "1234567890",
			AuthMethod: "password",
			Database:   ut.ToString(values["database"], ""),
			User:       ut.IM{"username": ut.ToString(values["username"], "")},
			Expiry:     time.Now().Add(time.Hour * 24),
		}
		client.SetSearch("customer_simple", ut.IM{}, true)
		return evt
	case ClientEventLogOut:
		client.Ticket = Ticket{}
		return evt
	case ClientEventSideMenu:
		value := ut.ToString(evt.Value, "")
		if value == "customer_browser" || value == "customer_simple" {
			client.SetSearch(value, ut.IM{}, (value == "customer_simple"))
			return evt
		}
		if value == "editor_cancel" {
			client.ResetEditor()
			return evt
		}
	case SearchEventSelected:
		client.SetEditor("customer", "main", ut.IM{})
		return evt
	case BrowserEventEditRow:
		client.SetEditor("customer", "main", ut.IM{})
		return evt
	case ClientEventModule:
		value := ut.ToString(evt.Value, "")
		if value == "search" {
			client.ResetEditor()
			client.HideSideBar = false
			return evt
		}
		if value == "setting" {
			client.SetForm("setting", ut.IM{}, 0, false)
			client.HideSideBar = true
			return evt
		}
		if value == "info" {
			client.SetForm("info", ut.IM{}, 0, true)
			return evt
		}
	}
	return toast(evt.Name)
}

// [Client] test and demo data
func TestClient(cc ClientComponent) []TestComponent {
	id := ut.ToString(cc.GetProperty("id"), "")
	eventURL := ut.ToString(cc.GetProperty("event_url"), "")
	requestValue := cc.GetProperty("request_value").(map[string]ut.IM)
	requestMap := cc.GetProperty("request_map").(map[string]ClientComponent)
	ztm, _ := time.Parse(time.DateTime, "0000-00-01 00:00:00")
	return []TestComponent{
		{
			Label:         "Login",
			ComponentType: ComponentTypeClient,
			Component: &Client{
				BaseComponent: BaseComponent{
					Id:           id + "login",
					EventURL:     eventURL,
					OnResponse:   testClientResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
				},
				LoginButtons: []LoginAuthButton{
					{Id: "google", Label: "Google", Icon: "google"},
					{Id: "facebook", Label: "Facebook", Icon: "facebook"},
				},
				Ticket: Ticket{
					SessionID: "SES012345",
					Database:  "demo",
				},
				ClientLabels:    testClientLabels,
				ClientLogin:     testClientLogin,
				ClientMenu:      testClientMenu,
				ClientSideBar:   testClientSideBar,
				ClientSearch:    testClientSearch,
				ClientBrowser:   testClientBrowser,
				ClientEditor:    testClientEditor,
				ClientForm:      testClientForm,
				ClientModalForm: testClientModalForm,
			},
		},
		{
			Label:         "Simple Search",
			ComponentType: ComponentTypeClient,
			Component: &Client{
				BaseComponent: BaseComponent{
					Id:           id + "simple_search",
					EventURL:     eventURL,
					OnResponse:   testClientResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Data: ut.IM{
						"search": ut.IM{
							"view":   "customer_simple",
							"simple": true,
						},
					},
				},
				Ticket: Ticket{
					SessionID:  "1234567890",
					AuthMethod: "password",
					Database:   "demo",
					User:       ut.IM{"username": "admin"},
					Expiry:     ztm,
				},
				ClientLabels:    testClientLabels,
				ClientLogin:     testClientLogin,
				ClientMenu:      testClientMenu,
				ClientSideBar:   testClientSideBar,
				ClientSearch:    testClientSearch,
				ClientBrowser:   testClientBrowser,
				ClientEditor:    testClientEditor,
				ClientForm:      testClientForm,
				ClientModalForm: testClientModalForm,
			},
		},
		{
			Label:         "Search and disabled login",
			ComponentType: ComponentTypeClient,
			Component: &Client{
				BaseComponent: BaseComponent{
					Id:           id + "browser_search",
					EventURL:     eventURL,
					OnResponse:   testClientResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Data: ut.IM{
						"search": ut.IM{
							"view":   "customer_browser",
							"simple": false,
						},
					},
				},
				LoginDisabled:   true,
				ClientLabels:    testClientLabels,
				ClientLogin:     testClientLogin,
				ClientMenu:      testClientMenu,
				ClientSideBar:   testClientSideBar,
				ClientSearch:    testClientSearch,
				ClientBrowser:   testClientBrowser,
				ClientEditor:    testClientEditor,
				ClientForm:      testClientForm,
				ClientModalForm: testClientModalForm,
			},
		},
		{
			Label:         "Editor",
			ComponentType: ComponentTypeClient,
			Component: &Client{
				BaseComponent: BaseComponent{
					Id:           id + "editor",
					EventURL:     eventURL,
					OnResponse:   testClientResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Data: ut.IM{
						"search": ut.IM{
							"view":   "customer_simple",
							"simple": true,
						},
						"editor": ut.IM{
							"key":  "customer",
							"view": "main",
						},
					},
				},
				Ticket: Ticket{
					SessionID:  "1234567890",
					AuthMethod: "password",
					Database:   "demo",
					User:       ut.IM{"username": "admin"},
					Expiry:     time.Now().Add(time.Hour * 24),
				},
				ClientLabels:    testClientLabels,
				ClientLogin:     testClientLogin,
				ClientMenu:      testClientMenu,
				ClientSideBar:   testClientSideBar,
				ClientSearch:    testClientSearch,
				ClientBrowser:   testClientBrowser,
				ClientEditor:    testClientEditor,
				ClientForm:      testClientForm,
				ClientModalForm: testClientModalForm,
			},
		},
		{
			Label:         "Form and hidden side bar",
			ComponentType: ComponentTypeClient,
			Component: &Client{
				BaseComponent: BaseComponent{
					Id:           id + "form",
					EventURL:     eventURL,
					OnResponse:   testClientResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Data: ut.IM{
						"search": ut.IM{
							"view":   "customer_simple",
							"simple": true,
						},
						"editor": ut.IM{
							"key": "setting",
							"form": ut.IM{
								"key": "setting",
							},
						},
					},
				},
				Ticket: Ticket{
					SessionID:  "1234567890",
					AuthMethod: "password",
					Database:   "demo",
					User:       ut.IM{"username": "admin"},
					Expiry:     time.Now().Add(time.Hour * 24),
				},
				HideSideBar:     true,
				ClientLabels:    testClientLabels,
				ClientLogin:     testClientLogin,
				ClientMenu:      testClientMenu,
				ClientSideBar:   testClientSideBar,
				ClientSearch:    testClientSearch,
				ClientBrowser:   testClientBrowser,
				ClientEditor:    testClientEditor,
				ClientForm:      testClientForm,
				ClientModalForm: testClientModalForm,
			},
		},
		{
			Label:         "Modal Form",
			ComponentType: ComponentTypeClient,
			Component: &Client{
				BaseComponent: BaseComponent{
					Id:           id + "modal_form",
					EventURL:     eventURL,
					OnResponse:   testClientResponse,
					RequestValue: requestValue,
					RequestMap:   requestMap,
					Data: ut.IM{
						"search": ut.IM{
							"view":   "customer_simple",
							"simple": true,
						},
						"editor": ut.IM{
							"key": "setting",
							"form": ut.IM{
								"key": "setting",
							},
						},
						"modal": ut.IM{
							"key": "info",
						},
					},
				},
				Ticket: Ticket{
					SessionID:  "1234567890",
					AuthMethod: "password",
					Database:   "demo",
					User:       ut.IM{"username": "admin"},
					Expiry:     time.Now().Add(time.Hour * 24),
				},
				ClientLabels:    testClientLabels,
				ClientLogin:     testClientLogin,
				ClientMenu:      testClientMenu,
				ClientSideBar:   testClientSideBar,
				ClientSearch:    testClientSearch,
				ClientBrowser:   testClientBrowser,
				ClientEditor:    testClientEditor,
				ClientForm:      testClientForm,
				ClientModalForm: testClientModalForm,
			},
		},
	}
}
