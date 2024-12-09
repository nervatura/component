package component

import (
	"reflect"
	"testing"

	ut "github.com/nervatura/component/pkg/util"
)

func TestApplication_Render(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Title         string
		Theme         string
		Header        ut.SM
		Script        []string
		HeadLink      []HeadLink
		MainComponent ClientComponent
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "header",
			fields: fields{
				Header: ut.SM{"X-Session-Token": "TOKEN1234", "X-CSRF-Token": "CSRF1234"},
			},
			wantErr: false,
		},
		{
			name: "ok",
			fields: fields{
				MainComponent: &BaseComponent{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Theme:         tt.fields.Theme,
				Header:        tt.fields.Header,
				Script:        tt.fields.Script,
				HeadLink:      tt.fields.HeadLink,
				MainComponent: tt.fields.MainComponent,
			}
			_, err := app.Render()
			if (err != nil) != tt.wantErr {
				t.Errorf("Application.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestApplication_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Title         string
		Theme         string
		Header        ut.SM
		Script        []string
		HeadLink      []HeadLink
		MainComponent ClientComponent
	}
	type args struct {
		propName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "get",
			args: args{
				propName: "title",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Theme:         tt.fields.Theme,
				Header:        tt.fields.Header,
				Script:        tt.fields.Script,
				HeadLink:      tt.fields.HeadLink,
				MainComponent: tt.fields.MainComponent,
			}
			if got := app.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Application.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_Validation(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Title         string
		Theme         string
		Header        ut.SM
		Script        []string
		HeadLink      []HeadLink
		MainComponent ClientComponent
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
			name: "main",
			args: args{
				propName:  "main",
				propValue: &Login{},
			},
			want: &Login{},
		},
		{
			name: "main_nil",
			args: args{
				propName:  "main",
				propValue: nil,
			},
			want: nil,
		},
		{
			name: "header_imap",
			args: args{
				propName:  "header",
				propValue: ut.IM{"X-Session-Token": "TOKEN1234", "X-CSRF-Token": "CSRF1234"},
			},
			want: ut.SM{"X-Session-Token": "TOKEN1234", "X-CSRF-Token": "CSRF1234"},
		},
		{
			name: "script_il",
			args: args{
				propName:  "script",
				propValue: []interface{}{"script1", "script2"},
			},
			want: []string{"script1", "script2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Theme:         tt.fields.Theme,
				Header:        tt.fields.Header,
				Script:        tt.fields.Script,
				HeadLink:      tt.fields.HeadLink,
				MainComponent: tt.fields.MainComponent,
			}
			if got := app.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Application.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Title         string
		Theme         string
		Header        ut.SM
		Script        []string
		HeadLink      []HeadLink
		MainComponent ClientComponent
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
			name: "main",
			args: args{
				propName:  "main",
				propValue: &Login{},
			},
			want: &Login{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Theme:         tt.fields.Theme,
				Header:        tt.fields.Header,
				Script:        tt.fields.Script,
				HeadLink:      tt.fields.HeadLink,
				MainComponent: tt.fields.MainComponent,
			}
			if got := app.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Application.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_OnRequest(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Title         string
		Theme         string
		Header        ut.SM
		Script        []string
		HeadLink      []HeadLink
		MainComponent ClientComponent
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
						"ID12345": &Login{},
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
			app := &Application{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Theme:         tt.fields.Theme,
				Header:        tt.fields.Header,
				Script:        tt.fields.Script,
				HeadLink:      tt.fields.HeadLink,
				MainComponent: tt.fields.MainComponent,
			}
			app.OnRequest(tt.args.te)
		})
	}
}

func TestApplication_getComponent(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Title         string
		Theme         string
		Header        ut.SM
		Script        []string
		HeadLink      []HeadLink
		MainComponent ClientComponent
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "main",
			fields: fields{
				MainComponent: &BaseComponent{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Theme:         tt.fields.Theme,
				Header:        tt.fields.Header,
				Script:        tt.fields.Script,
				HeadLink:      tt.fields.HeadLink,
				MainComponent: tt.fields.MainComponent,
			}
			_, err := app.getComponent()
			if (err != nil) != tt.wantErr {
				t.Errorf("Application.getComponent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
