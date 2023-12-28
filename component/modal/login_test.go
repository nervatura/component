package modal

import (
	"reflect"
	"testing"

	fm "github.com/nervatura/component/component/atom"
	bc "github.com/nervatura/component/component/base"
)

func TestDemoLogin(t *testing.T) {
	for _, tt := range DemoLogin("/demo", "") {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	demoLoginResponse(bc.ResponseEvent{Name: LoginEventLogin, Trigger: &Login{}})
	demoLoginResponse(bc.ResponseEvent{Name: LoginEventLang, Trigger: &Login{}})
	demoLoginResponse(bc.ResponseEvent{Name: LoginEventTheme, Trigger: &Login{}})
}

func TestLogin_Validation(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Version       string
		Lang          string
		Theme         string
		Labels        bc.SM
		Locales       []fm.SelectOption
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
		{
			name: "locales",
			args: args{
				propName:  "locales",
				propValue: "",
			},
			want: []fm.SelectOption{{Value: "en", Text: "en"}},
		},
		{
			name: "labels",
			args: args{
				propName:  "labels",
				propValue: "",
			},
			want: loginDefaultLabel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lgn := &Login{
				BaseComponent: tt.fields.BaseComponent,
				Version:       tt.fields.Version,
				Lang:          tt.fields.Lang,
				Theme:         tt.fields.Theme,
				Labels:        tt.fields.Labels,
				Locales:       tt.fields.Locales,
			}
			if got := lgn.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogin_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Version       string
		Lang          string
		Theme         string
		Labels        bc.SM
		Locales       []fm.SelectOption
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
			lgn := &Login{
				BaseComponent: tt.fields.BaseComponent,
				Version:       tt.fields.Version,
				Lang:          tt.fields.Lang,
				Theme:         tt.fields.Theme,
				Labels:        tt.fields.Labels,
				Locales:       tt.fields.Locales,
			}
			if got := lgn.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogin_response(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Version       string
		Lang          string
		Theme         string
		Labels        bc.SM
		Locales       []fm.SelectOption
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
			name: "change",
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "username",
					Name:        LoginEventChange,
				},
			},
		},
		{
			name: "theme",
			fields: fields{
				BaseComponent: bc.BaseComponent{
					OnResponse: func(evt bc.ResponseEvent) (re bc.ResponseEvent) {
						evt.Trigger = &Login{}
						return evt
					},
				},
			},
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "theme",
					Name:        LoginEventTheme,
				},
			},
		},
		{
			name: "login",
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "login",
					Name:        LoginEventLogin,
				},
			},
		},
		{
			name: "lang",
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "lang",
					Name:        LoginEventLang,
				},
			},
		},
		{
			name: "invalid",
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "invalid",
					Name:        LoginEventLang,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lgn := &Login{
				BaseComponent: tt.fields.BaseComponent,
				Version:       tt.fields.Version,
				Lang:          tt.fields.Lang,
				Theme:         tt.fields.Theme,
				Labels:        tt.fields.Labels,
				Locales:       tt.fields.Locales,
			}
			lgn.response(tt.args.evt)
		})
	}
}

func TestLogin_msg(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Version       string
		Lang          string
		Theme         string
		Labels        bc.SM
		Locales       []fm.SelectOption
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
				labelID: "missing_id",
			},
			want: "missing_id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lgn := &Login{
				BaseComponent: tt.fields.BaseComponent,
				Version:       tt.fields.Version,
				Lang:          tt.fields.Lang,
				Theme:         tt.fields.Theme,
				Labels:        tt.fields.Labels,
				Locales:       tt.fields.Locales,
			}
			if got := lgn.msg(tt.args.labelID); got != tt.want {
				t.Errorf("Login.msg() = %v, want %v", got, tt.want)
			}
		})
	}
}
