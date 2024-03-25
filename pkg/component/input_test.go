package component

import (
	"reflect"
	"testing"
)

func TestTestInput(t *testing.T) {
	for _, tt := range TestInput(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	demoInputResponse(ResponseEvent{Trigger: &Input{}})
}

func TestInput_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Type          string
		Value         string
		Placeholder   string
		Label         string
		Disabled      bool
		ReadOnly      bool
		AutoFocus     bool
		Invalid       bool
		Accept        string
		MaxLength     int64
		Size          int64
		Full          bool
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
			inp := &Input{
				BaseComponent: tt.fields.BaseComponent,
				Type:          tt.fields.Type,
				Value:         tt.fields.Value,
				Placeholder:   tt.fields.Placeholder,
				Label:         tt.fields.Label,
				Disabled:      tt.fields.Disabled,
				ReadOnly:      tt.fields.ReadOnly,
				AutoFocus:     tt.fields.AutoFocus,
				Invalid:       tt.fields.Invalid,
				Accept:        tt.fields.Accept,
				MaxLength:     tt.fields.MaxLength,
				Size:          tt.fields.Size,
				Full:          tt.fields.Full,
			}
			if got := inp.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Input.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInput_Validation(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Type          string
		Value         string
		Placeholder   string
		Label         string
		Disabled      bool
		ReadOnly      bool
		AutoFocus     bool
		Invalid       bool
		Accept        string
		MaxLength     int64
		Size          int64
		Full          bool
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
			inp := &Input{
				BaseComponent: tt.fields.BaseComponent,
				Type:          tt.fields.Type,
				Value:         tt.fields.Value,
				Placeholder:   tt.fields.Placeholder,
				Label:         tt.fields.Label,
				Disabled:      tt.fields.Disabled,
				ReadOnly:      tt.fields.ReadOnly,
				AutoFocus:     tt.fields.AutoFocus,
				Invalid:       tt.fields.Invalid,
				Accept:        tt.fields.Accept,
				MaxLength:     tt.fields.MaxLength,
				Size:          tt.fields.Size,
				Full:          tt.fields.Full,
			}
			if got := inp.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Input.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInput_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Type          string
		Value         string
		Placeholder   string
		Label         string
		Disabled      bool
		ReadOnly      bool
		AutoFocus     bool
		Invalid       bool
		Accept        string
		MaxLength     int64
		Size          int64
		Full          bool
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
			inp := &Input{
				BaseComponent: tt.fields.BaseComponent,
				Type:          tt.fields.Type,
				Value:         tt.fields.Value,
				Placeholder:   tt.fields.Placeholder,
				Label:         tt.fields.Label,
				Disabled:      tt.fields.Disabled,
				ReadOnly:      tt.fields.ReadOnly,
				AutoFocus:     tt.fields.AutoFocus,
				Invalid:       tt.fields.Invalid,
				Accept:        tt.fields.Accept,
				MaxLength:     tt.fields.MaxLength,
				Size:          tt.fields.Size,
				Full:          tt.fields.Full,
			}
			if got := inp.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Input.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInput_OnRequest(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Type          string
		Value         string
		Placeholder   string
		Label         string
		Disabled      bool
		ReadOnly      bool
		AutoFocus     bool
		Invalid       bool
		Accept        string
		MaxLength     int64
		Size          int64
		Full          bool
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
				Trigger: &Input{},
				Name:    InputEventChange,
				Value:   "",
			},
		},
		{
			name: "OnResponse",
			fields: fields{
				BaseComponent: BaseComponent{
					OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
						evt.Trigger = &Input{}
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
				Trigger: &Input{},
				Name:    InputEventChange,
				Value:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inp := &Input{
				BaseComponent: tt.fields.BaseComponent,
				Type:          tt.fields.Type,
				Value:         tt.fields.Value,
				Placeholder:   tt.fields.Placeholder,
				Label:         tt.fields.Label,
				Disabled:      tt.fields.Disabled,
				ReadOnly:      tt.fields.ReadOnly,
				AutoFocus:     tt.fields.AutoFocus,
				Invalid:       tt.fields.Invalid,
				Accept:        tt.fields.Accept,
				MaxLength:     tt.fields.MaxLength,
				Size:          tt.fields.Size,
				Full:          tt.fields.Full,
			}
			if gotRe := inp.OnRequest(tt.args.te); !reflect.DeepEqual(gotRe, tt.wantRe) {
				t.Errorf("Input.OnRequest() = %v, want %v", gotRe, tt.wantRe)
			}
		})
	}
}
