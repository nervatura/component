package component

import (
	"reflect"
	"testing"

	ut "github.com/nervatura/component/pkg/util"
)

func TestTestMenuBar(t *testing.T) {
	for _, tt := range TestMenuBar(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
}

func TestMenuBar_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Value             string
		SideBar           bool
		SideBarVisibility string
		LabelHide         string
		LabelMenu         string
		Items             []MenuBarItem
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
				BaseComponent:     tt.fields.BaseComponent,
				Value:             tt.fields.Value,
				SideBar:           tt.fields.SideBar,
				SideBarVisibility: tt.fields.SideBarVisibility,
				LabelHide:         tt.fields.LabelHide,
				LabelMenu:         tt.fields.LabelMenu,
				Items:             tt.fields.Items,
			}
			if got := mnb.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MenuBar.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMenuBar_Validation(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Value             string
		SideBar           bool
		SideBarVisibility string
		LabelHide         string
		LabelMenu         string
		Items             []MenuBarItem
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
				BaseComponent:     tt.fields.BaseComponent,
				Value:             tt.fields.Value,
				SideBar:           tt.fields.SideBar,
				SideBarVisibility: tt.fields.SideBarVisibility,
				LabelHide:         tt.fields.LabelHide,
				LabelMenu:         tt.fields.LabelMenu,
				Items:             tt.fields.Items,
			}
			if got := mnb.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MenuBar.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMenuBar_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Value             string
		SideBar           bool
		SideBarVisibility string
		LabelHide         string
		LabelMenu         string
		Items             []MenuBarItem
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
				BaseComponent:     tt.fields.BaseComponent,
				Value:             tt.fields.Value,
				SideBar:           tt.fields.SideBar,
				SideBarVisibility: tt.fields.SideBarVisibility,
				LabelHide:         tt.fields.LabelHide,
				LabelMenu:         tt.fields.LabelMenu,
				Items:             tt.fields.Items,
			}
			if got := mnb.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MenuBar.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMenuBar_response(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Value             string
		SideBar           bool
		SideBarVisibility string
		LabelHide         string
		LabelMenu         string
		Items             []MenuBarItem
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
			name: "item",
			fields: fields{
				BaseComponent: BaseComponent{
					OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
						evt.Trigger = &MenuBar{}
						return evt
					},
				},
				Value: "edit",
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "item",
					Trigger:     &MenuBar{BaseComponent: BaseComponent{Data: ut.IM{"item": MenuBarItem{}}}},
					Value:       "search",
				},
			},
		},
		{
			name: "sidebar",
			args: args{
				evt: ResponseEvent{
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
				BaseComponent:     tt.fields.BaseComponent,
				Value:             tt.fields.Value,
				SideBar:           tt.fields.SideBar,
				SideBarVisibility: tt.fields.SideBarVisibility,
				LabelHide:         tt.fields.LabelHide,
				LabelMenu:         tt.fields.LabelMenu,
				Items:             tt.fields.Items,
			}
			mnb.response(tt.args.evt)
		})
	}
}
