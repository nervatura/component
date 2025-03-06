package component

import (
	"reflect"
	"testing"

	ut "github.com/nervatura/component/pkg/util"
)

func TestTestClient(t *testing.T) {
	for _, tt := range TestClient(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	testClientResponse(ResponseEvent{Trigger: &Client{}, Name: FormEventOK})
	testClientResponse(ResponseEvent{Trigger: &Client{
		BaseComponent: BaseComponent{
			Data: ut.IM{},
		},
	},
		Name: LoginEventLogin, Value: ut.IM{"database": "demo"}})
	testClientResponse(ResponseEvent{Trigger: &Client{}, Name: ClientEventLogOut})
	testClientResponse(ResponseEvent{Trigger: &Client{
		BaseComponent: BaseComponent{
			Data: ut.IM{},
		},
	},
		Name: ClientEventSideMenu, Value: "customer_simple"})
	testClientResponse(ResponseEvent{Trigger: &Client{
		BaseComponent: BaseComponent{
			Data: ut.IM{},
		},
	},
		Name: ClientEventSideMenu, Value: "editor_cancel"})
	testClientResponse(ResponseEvent{Trigger: &Client{
		BaseComponent: BaseComponent{
			Data: ut.IM{},
		},
	},
		Name: SearchEventSelected, Value: ut.IM{}})
	testClientResponse(ResponseEvent{Trigger: &Client{
		BaseComponent: BaseComponent{
			Data: ut.IM{},
		},
	},
		Name: BrowserEventEditRow, Value: ut.IM{}})
	testClientResponse(ResponseEvent{Trigger: &Client{
		BaseComponent: BaseComponent{
			Data: ut.IM{},
		},
	},
		Name: ClientEventModule, Value: "search"})
	testClientResponse(ResponseEvent{Trigger: &Client{
		BaseComponent: BaseComponent{
			Data: ut.IM{},
		},
	},
		Name: ClientEventModule, Value: "setting"})
	testClientResponse(ResponseEvent{Trigger: &Client{
		BaseComponent: BaseComponent{
			Data: ut.IM{},
		},
	},
		Name: ClientEventModule, Value: "info"})
	testClientResponse(ResponseEvent{Trigger: &Client{}, Name: FormEventChange})
}

func TestClient_Validation(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Version           string
		Ticket            Ticket
		Theme             string
		Lang              string
		SideBarVisibility string
		LoginDisabled     bool
		LoginURL          string
		LoginButtons      []LoginAuthButton
		HideSideBar       bool
		HideMenu          bool
		ClientLabels      func(lang string) ut.SM
		ClientMenu        func(labels ut.SM, config ut.IM) MenuBar
		ClientSideBar     func(moduleKey string, labels ut.SM, data ut.IM) SideBar
		ClientLogin       func(labels ut.SM, config ut.IM) Login
		ClientSearch      func(viewName string, labels ut.SM, searchData ut.IM) Search
		ClientBrowser     func(viewName string, labels ut.SM, searchData ut.IM) Browser
		ClientEditor      func(editorKey, viewName string, labels ut.SM, editorData ut.IM) Editor
		ClientModalForm   func(formKey string, labels ut.SM, data ut.IM) Form
		ClientForm        func(editorKey, formKey string, labels ut.SM, data ut.IM) Form
	}
	type args struct {
		propName  string
		propValue interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name:   "invalid_ticket",
			fields: fields{},
			args: args{
				propName: "ticket", propValue: "",
			},
			want: Ticket{},
		},
		{
			name:   "invalid",
			fields: fields{},
			args: args{
				propName: "invalid", propValue: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &Client{
				BaseComponent:     tt.fields.BaseComponent,
				Version:           tt.fields.Version,
				Ticket:            tt.fields.Ticket,
				Theme:             tt.fields.Theme,
				Lang:              tt.fields.Lang,
				SideBarVisibility: tt.fields.SideBarVisibility,
				LoginDisabled:     tt.fields.LoginDisabled,
				LoginURL:          tt.fields.LoginURL,
				LoginButtons:      tt.fields.LoginButtons,
				HideSideBar:       tt.fields.HideSideBar,
				HideMenu:          tt.fields.HideMenu,
				ClientLabels:      tt.fields.ClientLabels,
				ClientMenu:        tt.fields.ClientMenu,
				ClientSideBar:     tt.fields.ClientSideBar,
				ClientLogin:       tt.fields.ClientLogin,
				ClientSearch:      tt.fields.ClientSearch,
				ClientBrowser:     tt.fields.ClientBrowser,
				ClientEditor:      tt.fields.ClientEditor,
				ClientModalForm:   tt.fields.ClientModalForm,
				ClientForm:        tt.fields.ClientForm,
			}
			if got := cli.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_OnRequest(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Version           string
		Ticket            Ticket
		Theme             string
		Lang              string
		SideBarVisibility string
		LoginDisabled     bool
		LoginURL          string
		LoginButtons      []LoginAuthButton
		HideSideBar       bool
		HideMenu          bool
		ClientLabels      func(lang string) ut.SM
		ClientMenu        func(labels ut.SM, config ut.IM) MenuBar
		ClientSideBar     func(moduleKey string, labels ut.SM, data ut.IM) SideBar
		ClientLogin       func(labels ut.SM, config ut.IM) Login
		ClientSearch      func(viewName string, labels ut.SM, searchData ut.IM) Search
		ClientBrowser     func(viewName string, labels ut.SM, searchData ut.IM) Browser
		ClientEditor      func(editorKey, viewName string, labels ut.SM, editorData ut.IM) Editor
		ClientModalForm   func(formKey string, labels ut.SM, data ut.IM) Form
		ClientForm        func(editorKey, formKey string, labels ut.SM, data ut.IM) Form
	}
	type args struct {
		te TriggerEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "invalid",
			args: args{
				te: TriggerEvent{},
			},
		},
		{
			name: "valid",
			fields: fields{
				BaseComponent: BaseComponent{
					RequestMap: map[string]ClientComponent{
						"ID12345": &BaseComponent{},
					},
				},
			},
			args: args{
				te: TriggerEvent{
					Id: "ID12345",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &Client{
				BaseComponent:     tt.fields.BaseComponent,
				Version:           tt.fields.Version,
				Ticket:            tt.fields.Ticket,
				Theme:             tt.fields.Theme,
				Lang:              tt.fields.Lang,
				SideBarVisibility: tt.fields.SideBarVisibility,
				LoginDisabled:     tt.fields.LoginDisabled,
				LoginURL:          tt.fields.LoginURL,
				LoginButtons:      tt.fields.LoginButtons,
				HideSideBar:       tt.fields.HideSideBar,
				HideMenu:          tt.fields.HideMenu,
				ClientLabels:      tt.fields.ClientLabels,
				ClientMenu:        tt.fields.ClientMenu,
				ClientSideBar:     tt.fields.ClientSideBar,
				ClientLogin:       tt.fields.ClientLogin,
				ClientSearch:      tt.fields.ClientSearch,
				ClientBrowser:     tt.fields.ClientBrowser,
				ClientEditor:      tt.fields.ClientEditor,
				ClientModalForm:   tt.fields.ClientModalForm,
				ClientForm:        tt.fields.ClientForm,
			}
			cli.OnRequest(tt.args.te)
		})
	}
}

func TestClient_response(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Version           string
		Ticket            Ticket
		Theme             string
		Lang              string
		SideBarVisibility string
		LoginDisabled     bool
		LoginURL          string
		LoginButtons      []LoginAuthButton
		HideSideBar       bool
		HideMenu          bool
		ClientLabels      func(lang string) ut.SM
		ClientMenu        func(labels ut.SM, config ut.IM) MenuBar
		ClientSideBar     func(moduleKey string, labels ut.SM, data ut.IM) SideBar
		ClientLogin       func(labels ut.SM, config ut.IM) Login
		ClientSearch      func(viewName string, labels ut.SM, searchData ut.IM) Search
		ClientBrowser     func(viewName string, labels ut.SM, searchData ut.IM) Browser
		ClientEditor      func(editorKey, viewName string, labels ut.SM, editorData ut.IM) Editor
		ClientModalForm   func(formKey string, labels ut.SM, data ut.IM) Form
		ClientForm        func(editorKey, formKey string, labels ut.SM, data ut.IM) Form
	}
	type args struct {
		evt ResponseEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "modal",
			fields: fields{
				BaseComponent: BaseComponent{
					OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
						return evt
					},
				},
			},
			args: args{
				evt: ResponseEvent{
					Trigger:     &Form{},
					TriggerName: "modal",
				},
			},
		},
		{
			name:   "table",
			fields: fields{},
			args: args{
				evt: ResponseEvent{
					Trigger:     &Table{},
					TriggerName: "table",
					Name:        TableEventEditCell,
				},
			},
		},
		{
			name:   "table2",
			fields: fields{},
			args: args{
				evt: ResponseEvent{
					Trigger:     &Table{},
					TriggerName: "table",
					Name:        TableEventSort,
				},
			},
		},
		{
			name: "search",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
				},
			},
			args: args{
				evt: ResponseEvent{
					Trigger:     &Search{},
					TriggerName: "search",
					Name:        SearchEventSearch,
				},
			},
		},
		{
			name: "browser",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
				},
			},
			args: args{
				evt: ResponseEvent{
					Trigger:     &Browser{},
					TriggerName: "browser",
					Name:        BrowserEventSearch,
				},
			},
		},
		{
			name: "browser2",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
				},
			},
			args: args{
				evt: ResponseEvent{
					Trigger:     &Browser{},
					TriggerName: "browser",
					Name:        BrowserEventSetColumn,
				},
			},
		},
		{
			name: "browser3",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
				},
			},
			args: args{
				evt: ResponseEvent{
					Trigger:     &Browser{},
					TriggerName: "browser",
					Name:        BrowserEventAddFilter,
				},
			},
		},
		{
			name: "browser4",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
				},
			},
			args: args{
				evt: ResponseEvent{
					Trigger:     &Browser{},
					TriggerName: "browser",
					Name:        BrowserEventBookmark,
				},
			},
		},
		{
			name:   "editor",
			fields: fields{},
			args: args{
				evt: ResponseEvent{
					Trigger:     &Editor{},
					TriggerName: "editor",
					Name:        EditorEventField,
				},
			},
		},
		{
			name:   "editor_view",
			fields: fields{},
			args: args{
				evt: ResponseEvent{
					Trigger:     &Editor{},
					TriggerName: "editor",
					Name:        EditorEventView,
				},
			},
		},
		{
			name:   "form",
			fields: fields{},
			args: args{
				evt: ResponseEvent{
					Trigger:     &Form{},
					TriggerName: "form",
					Name:        FormEventOK,
				},
			},
		},
		{
			name:   "login",
			fields: fields{},
			args: args{
				evt: ResponseEvent{
					Trigger:     &Login{},
					TriggerName: "login",
					Name:        LoginEventLang,
				},
			},
		},
		{
			name: "login2",
			fields: fields{
				Theme: ThemeDark,
			},
			args: args{
				evt: ResponseEvent{
					Trigger:     &Login{},
					TriggerName: "login",
					Name:        LoginEventTheme,
				},
			},
		},
		{
			name: "login3",
			fields: fields{
				Theme: ThemeDark,
			},
			args: args{
				evt: ResponseEvent{
					Trigger: &Login{
						BaseComponent: BaseComponent{
							Data: ut.IM{
								"database": "demo",
							},
						},
					},
					TriggerName: "login",
					Name:        LoginEventLogin,
				},
			},
		},
		{
			name: "auth_login",
			fields: fields{
				Theme: ThemeDark,
			},
			args: args{
				evt: ResponseEvent{
					Trigger:     &Login{},
					TriggerName: "login",
					Name:        LoginEventAuth,
				},
			},
		},
		{
			name:   "side_menu",
			fields: fields{},
			args: args{
				evt: ResponseEvent{
					Trigger:     &SideBar{},
					TriggerName: "side_menu",
					Name:        SideBarEventItem,
				},
			},
		},
		{
			name:   "main_menu",
			fields: fields{},
			args: args{
				evt: ResponseEvent{
					Trigger:     &MenuBar{},
					TriggerName: "main_menu",
					Name:        MenuBarEventSide,
				},
			},
		},
		{
			name: "main_menu2",
			fields: fields{
				Theme: ThemeDark,
			},
			args: args{
				evt: ResponseEvent{
					Trigger:     &MenuBar{},
					TriggerName: "main_menu",
					Name:        MenuBarEventValue,
					Value:       "theme",
				},
			},
		},
		{
			name: "main_menu3",
			args: args{
				evt: ResponseEvent{
					Trigger:     &MenuBar{},
					TriggerName: "main_menu",
					Name:        MenuBarEventValue,
					Value:       "search",
				},
			},
		},
		{
			name: "main_menu4",
			args: args{
				evt: ResponseEvent{
					Trigger:     &MenuBar{},
					TriggerName: "main_menu",
					Name:        MenuBarEventValue,
					Value:       "logout",
				},
			},
		},
		{
			name: "invalid",
			args: args{
				evt: ResponseEvent{
					Trigger:     &BaseComponent{},
					TriggerName: "invalid",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &Client{
				BaseComponent:     tt.fields.BaseComponent,
				Version:           tt.fields.Version,
				Ticket:            tt.fields.Ticket,
				Theme:             tt.fields.Theme,
				Lang:              tt.fields.Lang,
				SideBarVisibility: tt.fields.SideBarVisibility,
				LoginDisabled:     tt.fields.LoginDisabled,
				LoginURL:          tt.fields.LoginURL,
				LoginButtons:      tt.fields.LoginButtons,
				HideSideBar:       tt.fields.HideSideBar,
				HideMenu:          tt.fields.HideMenu,
				ClientLabels:      tt.fields.ClientLabels,
				ClientMenu:        tt.fields.ClientMenu,
				ClientSideBar:     tt.fields.ClientSideBar,
				ClientLogin:       tt.fields.ClientLogin,
				ClientSearch:      tt.fields.ClientSearch,
				ClientBrowser:     tt.fields.ClientBrowser,
				ClientEditor:      tt.fields.ClientEditor,
				ClientModalForm:   tt.fields.ClientModalForm,
				ClientForm:        tt.fields.ClientForm,
			}
			cli.response(tt.args.evt)
		})
	}
}

func TestClient_GetSearchVisibleColumns(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Version           string
		Ticket            Ticket
		Theme             string
		Lang              string
		SideBarVisibility string
		LoginDisabled     bool
		LoginURL          string
		LoginButtons      []LoginAuthButton
		HideSideBar       bool
		HideMenu          bool
		ClientLabels      func(lang string) ut.SM
		ClientMenu        func(labels ut.SM, config ut.IM) MenuBar
		ClientSideBar     func(moduleKey string, labels ut.SM, data ut.IM) SideBar
		ClientLogin       func(labels ut.SM, config ut.IM) Login
		ClientSearch      func(viewName string, labels ut.SM, searchData ut.IM) Search
		ClientBrowser     func(viewName string, labels ut.SM, searchData ut.IM) Browser
		ClientEditor      func(editorKey, viewName string, labels ut.SM, editorData ut.IM) Editor
		ClientModalForm   func(formKey string, labels ut.SM, data ut.IM) Form
		ClientForm        func(editorKey, formKey string, labels ut.SM, data ut.IM) Form
	}
	type args struct {
		icols map[string]bool
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantCols map[string]bool
	}{
		{
			name: "map_bool",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{
						"search": ut.IM{
							"view": "customer",
							"customer": ut.IM{
								"visible_columns": map[string]bool{
									"col1": true,
									"col2": false,
								},
							},
						},
					},
				},
			},
			args: args{
				icols: map[string]bool{
					"col1": true,
					"col2": false,
				},
			},
			wantCols: map[string]bool{
				"col1": true,
				"col2": false,
			},
		},
		{
			name: "map_im",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{
						"search": ut.IM{
							"view": "customer",
							"customer": ut.IM{
								"visible_columns": ut.IM{
									"col1": true,
									"col2": false,
								},
							},
						},
					},
				},
			},
			args: args{
				icols: map[string]bool{
					"col1": true,
					"col2": false,
				},
			},
			wantCols: map[string]bool{
				"col1": true,
				"col2": false,
			},
		},
		{
			name: "default",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
				},
			},
			args: args{
				icols: map[string]bool{
					"col1": true,
					"col2": false,
				},
			},
			wantCols: map[string]bool{
				"col1": true,
				"col2": false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &Client{
				BaseComponent:     tt.fields.BaseComponent,
				Version:           tt.fields.Version,
				Ticket:            tt.fields.Ticket,
				Theme:             tt.fields.Theme,
				Lang:              tt.fields.Lang,
				SideBarVisibility: tt.fields.SideBarVisibility,
				LoginDisabled:     tt.fields.LoginDisabled,
				LoginURL:          tt.fields.LoginURL,
				LoginButtons:      tt.fields.LoginButtons,
				HideSideBar:       tt.fields.HideSideBar,
				HideMenu:          tt.fields.HideMenu,
				ClientLabels:      tt.fields.ClientLabels,
				ClientMenu:        tt.fields.ClientMenu,
				ClientSideBar:     tt.fields.ClientSideBar,
				ClientLogin:       tt.fields.ClientLogin,
				ClientSearch:      tt.fields.ClientSearch,
				ClientBrowser:     tt.fields.ClientBrowser,
				ClientEditor:      tt.fields.ClientEditor,
				ClientModalForm:   tt.fields.ClientModalForm,
				ClientForm:        tt.fields.ClientForm,
			}
			if gotCols := cli.GetSearchVisibleColumns(tt.args.icols); !reflect.DeepEqual(gotCols, tt.wantCols) {
				t.Errorf("Client.GetSearchVisibleColumns() = %v, want %v", gotCols, tt.wantCols)
			}
		})
	}
}

func TestClient_Labels(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Version           string
		Ticket            Ticket
		Theme             string
		Lang              string
		SideBarVisibility string
		LoginDisabled     bool
		LoginURL          string
		LoginButtons      []LoginAuthButton
		HideSideBar       bool
		HideMenu          bool
		ClientLabels      func(lang string) ut.SM
		ClientMenu        func(labels ut.SM, config ut.IM) MenuBar
		ClientSideBar     func(moduleKey string, labels ut.SM, data ut.IM) SideBar
		ClientLogin       func(labels ut.SM, config ut.IM) Login
		ClientSearch      func(viewName string, labels ut.SM, searchData ut.IM) Search
		ClientBrowser     func(viewName string, labels ut.SM, searchData ut.IM) Browser
		ClientEditor      func(editorKey, viewName string, labels ut.SM, editorData ut.IM) Editor
		ClientModalForm   func(formKey string, labels ut.SM, data ut.IM) Form
		ClientForm        func(editorKey, formKey string, labels ut.SM, data ut.IM) Form
	}
	tests := []struct {
		name   string
		fields fields
		want   ut.SM
	}{
		{
			name: "default",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
				},
			},
			want: ut.SM{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &Client{
				BaseComponent:     tt.fields.BaseComponent,
				Version:           tt.fields.Version,
				Ticket:            tt.fields.Ticket,
				Theme:             tt.fields.Theme,
				Lang:              tt.fields.Lang,
				SideBarVisibility: tt.fields.SideBarVisibility,
				LoginDisabled:     tt.fields.LoginDisabled,
				LoginURL:          tt.fields.LoginURL,
				LoginButtons:      tt.fields.LoginButtons,
				HideSideBar:       tt.fields.HideSideBar,
				HideMenu:          tt.fields.HideMenu,
				ClientLabels:      tt.fields.ClientLabels,
				ClientMenu:        tt.fields.ClientMenu,
				ClientSideBar:     tt.fields.ClientSideBar,
				ClientLogin:       tt.fields.ClientLogin,
				ClientSearch:      tt.fields.ClientSearch,
				ClientBrowser:     tt.fields.ClientBrowser,
				ClientEditor:      tt.fields.ClientEditor,
				ClientModalForm:   tt.fields.ClientModalForm,
				ClientForm:        tt.fields.ClientForm,
			}
			if got := cli.Labels(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Labels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetSearchFilters(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Version           string
		Ticket            Ticket
		Theme             string
		Lang              string
		SideBarVisibility string
		LoginDisabled     bool
		LoginURL          string
		LoginButtons      []LoginAuthButton
		HideSideBar       bool
		HideMenu          bool
		ClientLabels      func(lang string) ut.SM
		ClientMenu        func(labels ut.SM, config ut.IM) MenuBar
		ClientSideBar     func(moduleKey string, labels ut.SM, data ut.IM) SideBar
		ClientLogin       func(labels ut.SM, config ut.IM) Login
		ClientSearch      func(viewName string, labels ut.SM, searchData ut.IM) Search
		ClientBrowser     func(viewName string, labels ut.SM, searchData ut.IM) Browser
		ClientEditor      func(editorKey, viewName string, labels ut.SM, editorData ut.IM) Editor
		ClientModalForm   func(formKey string, labels ut.SM, data ut.IM) Form
		ClientForm        func(editorKey, formKey string, labels ut.SM, data ut.IM) Form
	}
	type args struct {
		value     string
		cfFilters interface{}
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantFilters []BrowserFilter
	}{
		{
			name: "default",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
				},
			},
			args: args{
				value:     "test",
				cfFilters: nil,
			},
			wantFilters: []BrowserFilter{},
		},
		{
			name: "map_filters",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{
						"search": ut.IM{
							"view": "customer",
							"customer": ut.IM{
								"filters": []BrowserFilter{},
							},
						},
					},
				},
			},
			args: args{
				value:     "test",
				cfFilters: nil,
			},
			wantFilters: []BrowserFilter{},
		},
		{
			name: "map_interface",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{
						"search": ut.IM{
							"view": "customer",
							"customer": ut.IM{
								"filters": []interface{}{
									map[string]interface{}{
										"field": "col1",
										"value": "test",
									},
								},
							},
						},
					},
				},
			},
			args: args{
				value:     "test",
				cfFilters: nil,
			},
			wantFilters: []BrowserFilter{
				{
					Field: "col1",
					Value: "test",
				},
			},
		},
		{
			name: "default_interface",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
				},
			},
			args: args{
				value: "test",
				cfFilters: []interface{}{
					map[string]interface{}{
						"field": "col1",
						"value": "test",
					},
				},
			},
			wantFilters: []BrowserFilter{
				{
					Field: "col1",
					Value: "test",
				},
			},
		},
		{
			name: "default_filters",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
				},
			},
			args: args{
				value: "test",
				cfFilters: []BrowserFilter{
					{
						Field: "col1",
						Value: "test",
					},
				},
			},
			wantFilters: []BrowserFilter{
				{
					Field: "col1",
					Value: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &Client{
				BaseComponent:     tt.fields.BaseComponent,
				Version:           tt.fields.Version,
				Ticket:            tt.fields.Ticket,
				Theme:             tt.fields.Theme,
				Lang:              tt.fields.Lang,
				SideBarVisibility: tt.fields.SideBarVisibility,
				LoginDisabled:     tt.fields.LoginDisabled,
				LoginURL:          tt.fields.LoginURL,
				LoginButtons:      tt.fields.LoginButtons,
				HideSideBar:       tt.fields.HideSideBar,
				HideMenu:          tt.fields.HideMenu,
				ClientLabels:      tt.fields.ClientLabels,
				ClientMenu:        tt.fields.ClientMenu,
				ClientSideBar:     tt.fields.ClientSideBar,
				ClientLogin:       tt.fields.ClientLogin,
				ClientSearch:      tt.fields.ClientSearch,
				ClientBrowser:     tt.fields.ClientBrowser,
				ClientEditor:      tt.fields.ClientEditor,
				ClientModalForm:   tt.fields.ClientModalForm,
				ClientForm:        tt.fields.ClientForm,
			}
			if gotFilters := cli.GetSearchFilters(tt.args.value, tt.args.cfFilters); !reflect.DeepEqual(gotFilters, tt.wantFilters) {
				t.Errorf("Client.GetSearchFilters() = %v, want %v", gotFilters, tt.wantFilters)
			}
		})
	}
}

func TestClient_CleanComponent(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Version           string
		Ticket            Ticket
		Theme             string
		Lang              string
		SideBarVisibility string
		LoginDisabled     bool
		LoginURL          string
		LoginButtons      []LoginAuthButton
		HideSideBar       bool
		HideMenu          bool
		ClientLabels      func(lang string) ut.SM
		ClientMenu        func(labels ut.SM, config ut.IM) MenuBar
		ClientSideBar     func(moduleKey string, labels ut.SM, data ut.IM) SideBar
		ClientLogin       func(labels ut.SM, config ut.IM) Login
		ClientSearch      func(viewName string, labels ut.SM, searchData ut.IM) Search
		ClientBrowser     func(viewName string, labels ut.SM, searchData ut.IM) Browser
		ClientEditor      func(editorKey, viewName string, labels ut.SM, editorData ut.IM) Editor
		ClientModalForm   func(formKey string, labels ut.SM, data ut.IM) Form
		ClientForm        func(editorKey, formKey string, labels ut.SM, data ut.IM) Form
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "login",
			fields: fields{
				BaseComponent: BaseComponent{
					Id: "id",
					RequestValue: map[string]map[string]interface{}{
						"id_login": {},
					},
				},
			},
			args: args{
				name: "login",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &Client{
				BaseComponent:     tt.fields.BaseComponent,
				Version:           tt.fields.Version,
				Ticket:            tt.fields.Ticket,
				Theme:             tt.fields.Theme,
				Lang:              tt.fields.Lang,
				SideBarVisibility: tt.fields.SideBarVisibility,
				LoginDisabled:     tt.fields.LoginDisabled,
				LoginURL:          tt.fields.LoginURL,
				LoginButtons:      tt.fields.LoginButtons,
				HideSideBar:       tt.fields.HideSideBar,
				HideMenu:          tt.fields.HideMenu,
				ClientLabels:      tt.fields.ClientLabels,
				ClientMenu:        tt.fields.ClientMenu,
				ClientSideBar:     tt.fields.ClientSideBar,
				ClientLogin:       tt.fields.ClientLogin,
				ClientSearch:      tt.fields.ClientSearch,
				ClientBrowser:     tt.fields.ClientBrowser,
				ClientEditor:      tt.fields.ClientEditor,
				ClientModalForm:   tt.fields.ClientModalForm,
				ClientForm:        tt.fields.ClientForm,
			}
			cli.CleanComponent(tt.args.name)
		})
	}
}

func TestClient_Msg(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Version           string
		Ticket            Ticket
		Theme             string
		Lang              string
		SideBarVisibility string
		LoginDisabled     bool
		LoginURL          string
		LoginButtons      []LoginAuthButton
		HideSideBar       bool
		HideMenu          bool
		ClientLabels      func(lang string) ut.SM
		ClientMenu        func(labels ut.SM, config ut.IM) MenuBar
		ClientSideBar     func(moduleKey string, labels ut.SM, data ut.IM) SideBar
		ClientLogin       func(labels ut.SM, config ut.IM) Login
		ClientSearch      func(viewName string, labels ut.SM, searchData ut.IM) Search
		ClientBrowser     func(viewName string, labels ut.SM, searchData ut.IM) Browser
		ClientEditor      func(editorKey, viewName string, labels ut.SM, editorData ut.IM) Editor
		ClientModalForm   func(formKey string, labels ut.SM, data ut.IM) Form
		ClientForm        func(editorKey, formKey string, labels ut.SM, data ut.IM) Form
	}
	type args struct {
		labelID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "default",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
				},
			},
			args: args{
				labelID: "login",
			},
			want: "login",
		},
		{
			name: "result value",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
				},
				ClientLabels: func(lang string) ut.SM {
					return ut.SM{
						"login": "login",
					}
				},
			},
			args: args{
				labelID: "login",
			},
			want: "login",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &Client{
				BaseComponent:     tt.fields.BaseComponent,
				Version:           tt.fields.Version,
				Ticket:            tt.fields.Ticket,
				Theme:             tt.fields.Theme,
				Lang:              tt.fields.Lang,
				SideBarVisibility: tt.fields.SideBarVisibility,
				LoginDisabled:     tt.fields.LoginDisabled,
				LoginURL:          tt.fields.LoginURL,
				LoginButtons:      tt.fields.LoginButtons,
				HideSideBar:       tt.fields.HideSideBar,
				HideMenu:          tt.fields.HideMenu,
				ClientLabels:      tt.fields.ClientLabels,
				ClientMenu:        tt.fields.ClientMenu,
				ClientSideBar:     tt.fields.ClientSideBar,
				ClientLogin:       tt.fields.ClientLogin,
				ClientSearch:      tt.fields.ClientSearch,
				ClientBrowser:     tt.fields.ClientBrowser,
				ClientEditor:      tt.fields.ClientEditor,
				ClientModalForm:   tt.fields.ClientModalForm,
				ClientForm:        tt.fields.ClientForm,
			}
			if got := cli.Msg(tt.args.labelID); got != tt.want {
				t.Errorf("Client.Msg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_SetConfigValue(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Version           string
		Ticket            Ticket
		Theme             string
		Lang              string
		SideBarVisibility string
		LoginDisabled     bool
		LoginURL          string
		LoginButtons      []LoginAuthButton
		HideSideBar       bool
		HideMenu          bool
		ClientLabels      func(lang string) ut.SM
		ClientMenu        func(labels ut.SM, config ut.IM) MenuBar
		ClientSideBar     func(moduleKey string, labels ut.SM, data ut.IM) SideBar
		ClientLogin       func(labels ut.SM, config ut.IM) Login
		ClientSearch      func(viewName string, labels ut.SM, searchData ut.IM) Search
		ClientBrowser     func(viewName string, labels ut.SM, searchData ut.IM) Browser
		ClientEditor      func(editorKey, viewName string, labels ut.SM, editorData ut.IM) Editor
		ClientModalForm   func(formKey string, labels ut.SM, data ut.IM) Form
		ClientForm        func(editorKey, formKey string, labels ut.SM, data ut.IM) Form
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "default",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
				},
			},
			args: args{
				key:   "key",
				value: "value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &Client{
				BaseComponent:     tt.fields.BaseComponent,
				Version:           tt.fields.Version,
				Ticket:            tt.fields.Ticket,
				Theme:             tt.fields.Theme,
				Lang:              tt.fields.Lang,
				SideBarVisibility: tt.fields.SideBarVisibility,
				LoginDisabled:     tt.fields.LoginDisabled,
				LoginURL:          tt.fields.LoginURL,
				LoginButtons:      tt.fields.LoginButtons,
				HideSideBar:       tt.fields.HideSideBar,
				HideMenu:          tt.fields.HideMenu,
				ClientLabels:      tt.fields.ClientLabels,
				ClientMenu:        tt.fields.ClientMenu,
				ClientSideBar:     tt.fields.ClientSideBar,
				ClientLogin:       tt.fields.ClientLogin,
				ClientSearch:      tt.fields.ClientSearch,
				ClientBrowser:     tt.fields.ClientBrowser,
				ClientEditor:      tt.fields.ClientEditor,
				ClientModalForm:   tt.fields.ClientModalForm,
				ClientForm:        tt.fields.ClientForm,
			}
			cli.SetConfigValue(tt.args.key, tt.args.value)
		})
	}
}
