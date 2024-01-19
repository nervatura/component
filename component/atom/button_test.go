package atom

import (
	"reflect"
	"testing"

	bc "github.com/nervatura/component/component/base"
)

func TestDemoButton(t *testing.T) {
	for _, tt := range DemoButton(&bc.BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	demoBtnResponse(bc.ResponseEvent{Trigger: &Button{}})
}

func TestButton_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent  bc.BaseComponent
		Type           string
		Align          string
		Label          string
		LabelComponent bc.ClientComponent
		Icon           string
		Disabled       bool
		AutoFocus      bool
		Full           bool
		Small          bool
		Selected       bool
		HideLabel      bool
		Badge          int64
		ShowBadge      bool
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
			btn := &Button{
				BaseComponent:  tt.fields.BaseComponent,
				Type:           tt.fields.Type,
				Align:          tt.fields.Align,
				Label:          tt.fields.Label,
				LabelComponent: tt.fields.LabelComponent,
				Icon:           tt.fields.Icon,
				Disabled:       tt.fields.Disabled,
				AutoFocus:      tt.fields.AutoFocus,
				Full:           tt.fields.Full,
				Small:          tt.fields.Small,
				Selected:       tt.fields.Selected,
				HideLabel:      tt.fields.HideLabel,
				Badge:          tt.fields.Badge,
				ShowBadge:      tt.fields.ShowBadge,
			}
			if got := btn.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Button.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestButton_Validation(t *testing.T) {
	type fields struct {
		BaseComponent  bc.BaseComponent
		Type           string
		Align          string
		Label          string
		LabelComponent bc.ClientComponent
		Icon           string
		Disabled       bool
		AutoFocus      bool
		Full           bool
		Small          bool
		Selected       bool
		HideLabel      bool
		Badge          int64
		ShowBadge      bool
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
			btn := &Button{
				BaseComponent:  tt.fields.BaseComponent,
				Type:           tt.fields.Type,
				Align:          tt.fields.Align,
				Label:          tt.fields.Label,
				LabelComponent: tt.fields.LabelComponent,
				Icon:           tt.fields.Icon,
				Disabled:       tt.fields.Disabled,
				AutoFocus:      tt.fields.AutoFocus,
				Full:           tt.fields.Full,
				Small:          tt.fields.Small,
				Selected:       tt.fields.Selected,
				HideLabel:      tt.fields.HideLabel,
				Badge:          tt.fields.Badge,
				ShowBadge:      tt.fields.ShowBadge,
			}
			if got := btn.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Button.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestButton_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent  bc.BaseComponent
		Type           string
		Align          string
		Label          string
		LabelComponent bc.ClientComponent
		Icon           string
		Disabled       bool
		AutoFocus      bool
		Full           bool
		Small          bool
		Selected       bool
		HideLabel      bool
		Badge          int64
		ShowBadge      bool
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
			btn := &Button{
				BaseComponent:  tt.fields.BaseComponent,
				Type:           tt.fields.Type,
				Align:          tt.fields.Align,
				Label:          tt.fields.Label,
				LabelComponent: tt.fields.LabelComponent,
				Icon:           tt.fields.Icon,
				Disabled:       tt.fields.Disabled,
				AutoFocus:      tt.fields.AutoFocus,
				Full:           tt.fields.Full,
				Small:          tt.fields.Small,
				Selected:       tt.fields.Selected,
				HideLabel:      tt.fields.HideLabel,
				Badge:          tt.fields.Badge,
				ShowBadge:      tt.fields.ShowBadge,
			}
			if got := btn.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Button.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestButton_OnRequest(t *testing.T) {
	type fields struct {
		BaseComponent  bc.BaseComponent
		Type           string
		Align          string
		Label          string
		LabelComponent bc.ClientComponent
		Icon           string
		Disabled       bool
		AutoFocus      bool
		Full           bool
		Small          bool
		Selected       bool
		HideLabel      bool
		Badge          int64
		ShowBadge      bool
	}
	type args struct {
		te bc.TriggerEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantRe bc.ResponseEvent
	}{
		{
			name: "base",
			args: args{
				te: bc.TriggerEvent{
					Id: "id",
				},
			},
			wantRe: bc.ResponseEvent{
				Trigger: &Button{},
				Name:    ButtonEventClick,
				Value:   "",
			},
		},
		{
			name: "OnResponse",
			fields: fields{
				BaseComponent: bc.BaseComponent{
					OnResponse: func(evt bc.ResponseEvent) (re bc.ResponseEvent) {
						evt.Trigger = &Button{}
						return evt
					},
				},
			},
			args: args{
				te: bc.TriggerEvent{
					Id: "id",
				},
			},
			wantRe: bc.ResponseEvent{
				Trigger: &Button{},
				Name:    ButtonEventClick,
				Value:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			btn := &Button{
				BaseComponent:  tt.fields.BaseComponent,
				Type:           tt.fields.Type,
				Align:          tt.fields.Align,
				Label:          tt.fields.Label,
				LabelComponent: tt.fields.LabelComponent,
				Icon:           tt.fields.Icon,
				Disabled:       tt.fields.Disabled,
				AutoFocus:      tt.fields.AutoFocus,
				Full:           tt.fields.Full,
				Small:          tt.fields.Small,
				Selected:       tt.fields.Selected,
				HideLabel:      tt.fields.HideLabel,
				Badge:          tt.fields.Badge,
				ShowBadge:      tt.fields.ShowBadge,
			}
			if gotRe := btn.OnRequest(tt.args.te); !reflect.DeepEqual(gotRe, tt.wantRe) {
				t.Errorf("Button.OnRequest() = %v, want %v", gotRe, tt.wantRe)
			}
		})
	}
}
