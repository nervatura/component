package component

import (
	"reflect"
	"testing"
)

func TestTestButton(t *testing.T) {
	for _, tt := range TestButton(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	testBtnResponse(ResponseEvent{Trigger: &Button{}})
}

func TestButton_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent  BaseComponent
		Type           string
		ButtonStyle    string
		Align          string
		Label          string
		LabelComponent ClientComponent
		Icon           string
		Disabled       bool
		AutoFocus      bool
		Full           bool
		Small          bool
		Selected       bool
		HideLabel      bool
		Badge          string
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
				ButtonStyle:    tt.fields.ButtonStyle,
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
			}
			if got := btn.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Button.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestButton_Validation(t *testing.T) {
	type fields struct {
		BaseComponent  BaseComponent
		Type           string
		ButtonStyle    string
		Align          string
		Label          string
		LabelComponent ClientComponent
		Icon           string
		Disabled       bool
		AutoFocus      bool
		Full           bool
		Small          bool
		Selected       bool
		HideLabel      bool
		Badge          string
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
				ButtonStyle:    tt.fields.ButtonStyle,
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
			}
			if got := btn.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Button.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestButton_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent  BaseComponent
		Type           string
		ButtonStyle    string
		Align          string
		Label          string
		LabelComponent ClientComponent
		Icon           string
		Disabled       bool
		AutoFocus      bool
		Full           bool
		Small          bool
		Selected       bool
		HideLabel      bool
		Badge          string
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
				ButtonStyle:    tt.fields.ButtonStyle,
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
			}
			if got := btn.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Button.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestButton_OnRequest(t *testing.T) {
	type fields struct {
		BaseComponent  BaseComponent
		Type           string
		ButtonStyle    string
		Align          string
		Label          string
		LabelComponent ClientComponent
		Icon           string
		Disabled       bool
		AutoFocus      bool
		Full           bool
		Small          bool
		Selected       bool
		HideLabel      bool
		Badge          string
	}
	type args struct {
		te TriggerEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantRe ResponseEvent
	}{
		{
			name: "base",
			args: args{
				te: TriggerEvent{
					Id: "id",
				},
			},
			wantRe: ResponseEvent{
				Trigger: &Button{},
				Name:    ButtonEventClick,
				Value:   "id",
			},
		},
		{
			name: "OnResponse",
			fields: fields{
				BaseComponent: BaseComponent{
					OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
						evt.Trigger = &Button{}
						return evt
					},
				},
			},
			args: args{
				te: TriggerEvent{
					Id: "id",
				},
			},
			wantRe: ResponseEvent{
				Trigger: &Button{},
				Name:    ButtonEventClick,
				Value:   "id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			btn := &Button{
				BaseComponent:  tt.fields.BaseComponent,
				Type:           tt.fields.Type,
				ButtonStyle:    tt.fields.ButtonStyle,
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
			}
			if gotRe := btn.OnRequest(tt.args.te); !reflect.DeepEqual(gotRe, tt.wantRe) {
				t.Errorf("Button.OnRequest() = %v, want %v", gotRe, tt.wantRe)
			}
		})
	}
}
