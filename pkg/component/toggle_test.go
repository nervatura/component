package component

import (
	"reflect"
	"testing"
)

func TestTestToggle(t *testing.T) {
	for _, tt := range TestToggle(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
}

func TestToggle_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Value         bool
		CheckBox      bool
		Border        bool
		Full          bool
		Disabled      bool
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
				propName: "disabled",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tgl := &Toggle{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				CheckBox:      tt.fields.CheckBox,
				Border:        tt.fields.Border,
				Full:          tt.fields.Full,
				Disabled:      tt.fields.Disabled,
			}
			if got := tgl.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Toggle.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToggle_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Value         bool
		CheckBox      bool
		Border        bool
		Full          bool
		Disabled      bool
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
			name: "missing",
			args: args{
				propName:  "missing",
				propValue: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tgl := &Toggle{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				CheckBox:      tt.fields.CheckBox,
				Border:        tt.fields.Border,
				Full:          tt.fields.Full,
				Disabled:      tt.fields.Disabled,
			}
			if got := tgl.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Toggle.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToggle_OnRequest(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Value         bool
		CheckBox      bool
		Border        bool
		Full          bool
		Disabled      bool
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
			name: "base",
			args: args{
				te: TriggerEvent{
					Id: "id",
				},
			},
		},
		{
			name: "OnResponse",
			fields: fields{
				BaseComponent: BaseComponent{
					OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
						evt.Trigger = &Toggle{}
						return evt
					},
				},
			},
			args: args{
				te: TriggerEvent{
					Id: "id",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tgl := &Toggle{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				CheckBox:      tt.fields.CheckBox,
				Border:        tt.fields.Border,
				Full:          tt.fields.Full,
				Disabled:      tt.fields.Disabled,
			}
			tgl.OnRequest(tt.args.te)
		})
	}
}
