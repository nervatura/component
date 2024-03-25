package demo

import (
	"reflect"
	"testing"

	ct "github.com/nervatura/component/pkg/component"
	ut "github.com/nervatura/component/pkg/util"
)

func TestNewDemo(t *testing.T) {
	type args struct {
		eventURL string
		title    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "new",
			args: args{
				eventURL: "/",
				title:    "demo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewDemo(tt.args.eventURL, tt.args.title)
		})
	}
}

func TestDemo_Render(t *testing.T) {
	type fields struct {
		BaseComponent ct.BaseComponent
		Title         string
		Theme         string
		ViewSize      string
		SelectedGroup string
		SelectedType  int64
		SelectedDemo  int64
		DemoMap       map[string][]DemoView
	}
	demo := NewDemo("", "")
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "atom",
			fields: fields{
				BaseComponent: ct.BaseComponent{
					RequestValue: map[string]ut.IM{},
					RequestMap:   map[string]ct.ClientComponent{},
				},
				DemoMap: DemoMap,
			},
			wantErr: false,
		},
		{
			name: "molecule",
			fields: fields{
				BaseComponent: ct.BaseComponent{
					RequestValue: map[string]ut.IM{},
					RequestMap:   map[string]ct.ClientComponent{},
				},
				DemoMap:       demo.DemoMap,
				SelectedGroup: ComponentGroupMolecule,
				SelectedType:  0,
				SelectedDemo:  0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sto := &Demo{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Theme:         tt.fields.Theme,
				ViewSize:      tt.fields.ViewSize,
				SelectedGroup: tt.fields.SelectedGroup,
				SelectedType:  tt.fields.SelectedType,
				SelectedDemo:  tt.fields.SelectedDemo,
				DemoMap:       tt.fields.DemoMap,
			}
			_, err := sto.Render()
			if (err != nil) != tt.wantErr {
				t.Errorf("Demo.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDemo_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent ct.BaseComponent
		Title         string
		Theme         string
		ViewSize      string
		SelectedGroup string
		SelectedType  int64
		SelectedDemo  int64
		DemoMap       map[string][]DemoView
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
				propName: "selected_group",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sto := &Demo{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Theme:         tt.fields.Theme,
				ViewSize:      tt.fields.ViewSize,
				SelectedGroup: tt.fields.SelectedGroup,
				SelectedType:  tt.fields.SelectedType,
				SelectedDemo:  tt.fields.SelectedDemo,
				DemoMap:       tt.fields.DemoMap,
			}
			if got := sto.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Demo.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDemo_Validation(t *testing.T) {
	type fields struct {
		BaseComponent ct.BaseComponent
		Title         string
		Theme         string
		ViewSize      string
		SelectedGroup string
		SelectedType  int64
		SelectedDemo  int64
		DemoMap       map[string][]DemoView
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
			name: "selected_type",
			args: args{
				propName:  "selected_type",
				propValue: 100,
			},
			want: int64(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sto := &Demo{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Theme:         tt.fields.Theme,
				ViewSize:      tt.fields.ViewSize,
				SelectedGroup: tt.fields.SelectedGroup,
				SelectedType:  tt.fields.SelectedType,
				SelectedDemo:  tt.fields.SelectedDemo,
				DemoMap:       tt.fields.DemoMap,
			}
			if got := sto.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Demo.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDemo_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent ct.BaseComponent
		Title         string
		Theme         string
		ViewSize      string
		SelectedGroup string
		SelectedType  int64
		SelectedDemo  int64
		DemoMap       map[string][]DemoView
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
			name: "selected_type",
			fields: fields{
				SelectedType: 1,
			},
			args: args{
				propName:  "selected_type",
				propValue: 0,
			},
			want: int64(0),
		},
		{
			name: "selected_demo",
			fields: fields{
				SelectedDemo: 1,
			},
			args: args{
				propName:  "selected_demo",
				propValue: 0,
			},
			want: int64(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sto := &Demo{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Theme:         tt.fields.Theme,
				ViewSize:      tt.fields.ViewSize,
				SelectedGroup: tt.fields.SelectedGroup,
				SelectedType:  tt.fields.SelectedType,
				SelectedDemo:  tt.fields.SelectedDemo,
				DemoMap:       tt.fields.DemoMap,
			}
			if got := sto.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Demo.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDemo_response(t *testing.T) {
	type fields struct {
		BaseComponent ct.BaseComponent
		Title         string
		Theme         string
		ViewSize      string
		SelectedGroup string
		SelectedType  int64
		SelectedDemo  int64
		DemoMap       map[string][]DemoView
	}
	type args struct {
		evt ct.ResponseEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ok",
			fields: fields{
				BaseComponent: ct.BaseComponent{
					OnResponse: func(evt ct.ResponseEvent) (re ct.ResponseEvent) {
						return evt
					},
				},
			},
			args: args{
				evt: ct.ResponseEvent{},
			},
		},
		{
			name: "theme",
			fields: fields{
				Theme: ct.ThemeLight,
			},
			args: args{
				evt: ct.ResponseEvent{
					TriggerName: DemoEventTheme,
				},
			},
		},
		{
			name: "view_size",
			fields: fields{
				ViewSize: ViewSizeCentered,
			},
			args: args{
				evt: ct.ResponseEvent{
					TriggerName: DemoEventViewSize,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sto := &Demo{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Theme:         tt.fields.Theme,
				ViewSize:      tt.fields.ViewSize,
				SelectedGroup: tt.fields.SelectedGroup,
				SelectedType:  tt.fields.SelectedType,
				SelectedDemo:  tt.fields.SelectedDemo,
				DemoMap:       tt.fields.DemoMap,
			}
			sto.response(tt.args.evt)
		})
	}
}

func TestDemo_OnRequest(t *testing.T) {
	type fields struct {
		BaseComponent ct.BaseComponent
		Title         string
		Theme         string
		ViewSize      string
		SelectedGroup string
		SelectedType  int64
		SelectedDemo  int64
		DemoMap       map[string][]DemoView
	}
	type args struct {
		te ct.TriggerEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "invalid",
			args: args{
				te: ct.TriggerEvent{},
			},
		},
		{
			name: "valid",
			fields: fields{
				BaseComponent: ct.BaseComponent{
					RequestMap: map[string]ct.ClientComponent{
						"ID12345": &Demo{},
					},
				},
			},
			args: args{
				te: ct.TriggerEvent{
					Id: "ID12345",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sto := &Demo{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Theme:         tt.fields.Theme,
				ViewSize:      tt.fields.ViewSize,
				SelectedGroup: tt.fields.SelectedGroup,
				SelectedType:  tt.fields.SelectedType,
				SelectedDemo:  tt.fields.SelectedDemo,
				DemoMap:       tt.fields.DemoMap,
			}
			sto.OnRequest(tt.args.te)
		})
	}
}
