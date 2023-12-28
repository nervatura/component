package molecule

import (
	"reflect"
	"testing"

	bc "github.com/nervatura/component/component/base"
)

func TestDemoMenuBar(t *testing.T) {
	for _, tt := range DemoMenuBar("/demo", "") {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
}

func TestMenuBar_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent  bc.BaseComponent
		Value          string
		SideBar        bool
		SideVisibility string
		LabelHide      string
		LabelMenu      string
		Items          []MenuBarItem
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
				propName: "value",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mnb := &MenuBar{
				BaseComponent:  tt.fields.BaseComponent,
				Value:          tt.fields.Value,
				SideBar:        tt.fields.SideBar,
				SideVisibility: tt.fields.SideVisibility,
				LabelHide:      tt.fields.LabelHide,
				LabelMenu:      tt.fields.LabelMenu,
				Items:          tt.fields.Items,
			}
			if got := mnb.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MenuBar.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMenuBar_Validation(t *testing.T) {
	type fields struct {
		BaseComponent  bc.BaseComponent
		Value          string
		SideBar        bool
		SideVisibility string
		LabelHide      string
		LabelMenu      string
		Items          []MenuBarItem
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
			mnb := &MenuBar{
				BaseComponent:  tt.fields.BaseComponent,
				Value:          tt.fields.Value,
				SideBar:        tt.fields.SideBar,
				SideVisibility: tt.fields.SideVisibility,
				LabelHide:      tt.fields.LabelHide,
				LabelMenu:      tt.fields.LabelMenu,
				Items:          tt.fields.Items,
			}
			if got := mnb.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MenuBar.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMenuBar_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent  bc.BaseComponent
		Value          string
		SideBar        bool
		SideVisibility string
		LabelHide      string
		LabelMenu      string
		Items          []MenuBarItem
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
			mnb := &MenuBar{
				BaseComponent:  tt.fields.BaseComponent,
				Value:          tt.fields.Value,
				SideBar:        tt.fields.SideBar,
				SideVisibility: tt.fields.SideVisibility,
				LabelHide:      tt.fields.LabelHide,
				LabelMenu:      tt.fields.LabelMenu,
				Items:          tt.fields.Items,
			}
			if got := mnb.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MenuBar.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMenuBar_response(t *testing.T) {
	type fields struct {
		BaseComponent  bc.BaseComponent
		Value          string
		SideBar        bool
		SideVisibility string
		LabelHide      string
		LabelMenu      string
		Items          []MenuBarItem
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
			name: "item",
			fields: fields{
				BaseComponent: bc.BaseComponent{
					OnResponse: func(evt bc.ResponseEvent) (re bc.ResponseEvent) {
						evt.Trigger = &MenuBar{}
						return evt
					},
				},
				Value: "edit",
			},
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "item",
					Trigger:     &MenuBar{BaseComponent: bc.BaseComponent{Data: bc.IM{"item": MenuBarItem{}}}},
					Value:       "search",
				},
			},
		},
		{
			name: "sidebar",
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "sidebar",
					Trigger:     &MenuBar{},
					Value:       "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mnb := &MenuBar{
				BaseComponent:  tt.fields.BaseComponent,
				Value:          tt.fields.Value,
				SideBar:        tt.fields.SideBar,
				SideVisibility: tt.fields.SideVisibility,
				LabelHide:      tt.fields.LabelHide,
				LabelMenu:      tt.fields.LabelMenu,
				Items:          tt.fields.Items,
			}
			mnb.response(tt.args.evt)
		})
	}
}
