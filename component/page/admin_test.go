package page

import (
	"reflect"
	"testing"

	bc "github.com/nervatura/component/component/base"
)

func TestAdmin_Render(t *testing.T) {
	for _, tt := range DemoAdmin("/demo", "") {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	demoAdminResponse(bc.ResponseEvent{Name: AdminEventModule, Trigger: &Admin{}, Value: "client"})
	demoAdminResponse(bc.ResponseEvent{Name: AdminEventCreate, Trigger: &Admin{}})
	demoAdminResponse(bc.ResponseEvent{Name: AdminEventTheme, Trigger: &Admin{}})
}

func TestAdmin_Validation(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Version       string
		Theme         string
		Module        string
		View          string
		Token         string
		HelpURL       string
		ClientURL     string
		Labels        bc.SM
		Verified      func(token string) bool
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
			name: "base",
			args: args{
				propName:  "id",
				propValue: "BTNID",
			},
			want: "BTNID",
		},
		{
			name: "invalid",
			args: args{
				propName:  "invalid",
				propValue: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adm := &Admin{
				BaseComponent: tt.fields.BaseComponent,
				Version:       tt.fields.Version,
				Theme:         tt.fields.Theme,
				Module:        tt.fields.Module,
				View:          tt.fields.View,
				Token:         tt.fields.Token,
				HelpURL:       tt.fields.HelpURL,
				ClientURL:     tt.fields.ClientURL,
				Labels:        tt.fields.Labels,
				Verified:      tt.fields.Verified,
			}
			if got := adm.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Admin.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdmin_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Version       string
		Theme         string
		Module        string
		View          string
		Token         string
		HelpURL       string
		ClientURL     string
		Labels        bc.SM
		Verified      func(token string) bool
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
			name: "base",
			args: args{
				propName:  "id",
				propValue: "BTNID",
			},
			want: "BTNID",
		},
		{
			name: "invalid",
			args: args{
				propName:  "invalid",
				propValue: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adm := &Admin{
				BaseComponent: tt.fields.BaseComponent,
				Version:       tt.fields.Version,
				Theme:         tt.fields.Theme,
				Module:        tt.fields.Module,
				View:          tt.fields.View,
				Token:         tt.fields.Token,
				HelpURL:       tt.fields.HelpURL,
				ClientURL:     tt.fields.ClientURL,
				Labels:        tt.fields.Labels,
				Verified:      tt.fields.Verified,
			}
			if got := adm.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Admin.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdmin_msg(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Version       string
		Theme         string
		Module        string
		View          string
		Token         string
		HelpURL       string
		ClientURL     string
		Labels        bc.SM
		Verified      func(token string) bool
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
			name: "missing",
			args: args{
				labelID: "missing",
			},
			want: "missing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adm := &Admin{
				BaseComponent: tt.fields.BaseComponent,
				Version:       tt.fields.Version,
				Theme:         tt.fields.Theme,
				Module:        tt.fields.Module,
				View:          tt.fields.View,
				Token:         tt.fields.Token,
				HelpURL:       tt.fields.HelpURL,
				ClientURL:     tt.fields.ClientURL,
				Labels:        tt.fields.Labels,
				Verified:      tt.fields.Verified,
			}
			if got := adm.msg(tt.args.labelID); got != tt.want {
				t.Errorf("Admin.msg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdmin_response(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Version       string
		Theme         string
		Module        string
		View          string
		Token         string
		HelpURL       string
		ClientURL     string
		LocalesURL    string
		Labels        bc.SM
		Verified      func(token string) bool
	}
	type args struct {
		evt bc.ResponseEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "api_key",
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "api_key",
				},
			},
		},
		{
			name: "theme",
			fields: fields{
				BaseComponent: bc.BaseComponent{
					OnResponse: func(evt bc.ResponseEvent) (re bc.ResponseEvent) {
						evt.Trigger = &Admin{}
						return evt
					},
				},
			},
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "theme",
				},
			},
		},
		{
			name: "main_menu_help",
			fields: fields{
				HelpURL: "/help",
			},
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "main_menu",
					Value:       "help",
				},
			},
		},
		{
			name: "main_menu_client",
			fields: fields{
				ClientURL: "/client",
			},
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "main_menu",
					Value:       "client",
				},
			},
		},
		{
			name: "main_menu_locales",
			fields: fields{
				LocalesURL: "/locales",
			},
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "main_menu",
					Value:       "locales",
				},
			},
		},
		{
			name: "view_menu",
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "view_menu",
					Value:       "logout",
				},
			},
		},
		{
			name: "create",
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "create",
				},
			},
		},
		{
			name: "login",
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "login",
				},
			},
		},
		{
			name: "report_install",
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "report_install",
				},
			},
		},
		{
			name: "report_delete",
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "report_delete",
				},
			},
		},
		{
			name: "password_change",
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "password_change",
				},
			},
		},
		{
			name: "missing",
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "missing",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adm := &Admin{
				BaseComponent: tt.fields.BaseComponent,
				Version:       tt.fields.Version,
				Theme:         tt.fields.Theme,
				Module:        tt.fields.Module,
				View:          tt.fields.View,
				Token:         tt.fields.Token,
				HelpURL:       tt.fields.HelpURL,
				ClientURL:     tt.fields.ClientURL,
				LocalesURL:    tt.fields.LocalesURL,
				Labels:        tt.fields.Labels,
				Verified:      tt.fields.Verified,
			}
			adm.response(tt.args.evt)
		})
	}
}
