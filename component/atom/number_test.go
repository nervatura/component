package atom

import (
	"reflect"
	"testing"

	bc "github.com/nervatura/component/component/base"
)

func TestDemoNumberInput(t *testing.T) {
	for _, tt := range DemoNumberInput(&bc.BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
}

func TestNumberInput_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Value         float64
		Integer       bool
		Label         string
		SetMax        bool
		MaxValue      float64
		SetMin        bool
		MinValue      float64
		Disabled      bool
		ReadOnly      bool
		AutoFocus     bool
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
			inp := &NumberInput{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Integer:       tt.fields.Integer,
				Label:         tt.fields.Label,
				SetMax:        tt.fields.SetMax,
				MaxValue:      tt.fields.MaxValue,
				SetMin:        tt.fields.SetMin,
				MinValue:      tt.fields.MinValue,
				Disabled:      tt.fields.Disabled,
				ReadOnly:      tt.fields.ReadOnly,
				AutoFocus:     tt.fields.AutoFocus,
				Full:          tt.fields.Full,
			}
			if got := inp.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NumberInput.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumberInput_Validation(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Value         float64
		Integer       bool
		Label         string
		SetMax        bool
		MaxValue      float64
		SetMin        bool
		MinValue      float64
		Disabled      bool
		ReadOnly      bool
		AutoFocus     bool
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
		{
			name: "min value",
			fields: fields{
				SetMin:   true,
				MinValue: 10,
			},
			args: args{
				propName:  "value",
				propValue: 1,
			},
			want: float64(10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inp := &NumberInput{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Integer:       tt.fields.Integer,
				Label:         tt.fields.Label,
				SetMax:        tt.fields.SetMax,
				MaxValue:      tt.fields.MaxValue,
				SetMin:        tt.fields.SetMin,
				MinValue:      tt.fields.MinValue,
				Disabled:      tt.fields.Disabled,
				ReadOnly:      tt.fields.ReadOnly,
				AutoFocus:     tt.fields.AutoFocus,
				Full:          tt.fields.Full,
			}
			if got := inp.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NumberInput.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumberInput_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Value         float64
		Integer       bool
		Label         string
		SetMax        bool
		MaxValue      float64
		SetMin        bool
		MinValue      float64
		Disabled      bool
		ReadOnly      bool
		AutoFocus     bool
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
			inp := &NumberInput{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Integer:       tt.fields.Integer,
				Label:         tt.fields.Label,
				SetMax:        tt.fields.SetMax,
				MaxValue:      tt.fields.MaxValue,
				SetMin:        tt.fields.SetMin,
				MinValue:      tt.fields.MinValue,
				Disabled:      tt.fields.Disabled,
				ReadOnly:      tt.fields.ReadOnly,
				AutoFocus:     tt.fields.AutoFocus,
				Full:          tt.fields.Full,
			}
			if got := inp.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NumberInput.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumberInput_OnRequest(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Value         float64
		Integer       bool
		Label         string
		SetMax        bool
		MaxValue      float64
		SetMin        bool
		MinValue      float64
		Disabled      bool
		ReadOnly      bool
		AutoFocus     bool
		Full          bool
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
				Trigger: &NumberInput{},
				Name:    NumberEventChange,
				Value:   float64(0),
			},
		},
		{
			name: "OnResponse",
			fields: fields{
				BaseComponent: bc.BaseComponent{
					OnResponse: func(evt bc.ResponseEvent) (re bc.ResponseEvent) {
						evt.Trigger = &NumberInput{}
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
				Trigger: &NumberInput{},
				Name:    NumberEventChange,
				Value:   float64(0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inp := &NumberInput{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Integer:       tt.fields.Integer,
				Label:         tt.fields.Label,
				SetMax:        tt.fields.SetMax,
				MaxValue:      tt.fields.MaxValue,
				SetMin:        tt.fields.SetMin,
				MinValue:      tt.fields.MinValue,
				Disabled:      tt.fields.Disabled,
				ReadOnly:      tt.fields.ReadOnly,
				AutoFocus:     tt.fields.AutoFocus,
				Full:          tt.fields.Full,
			}
			if gotRe := inp.OnRequest(tt.args.te); !reflect.DeepEqual(gotRe, tt.wantRe) {
				t.Errorf("NumberInput.OnRequest() = %v, want %v", gotRe, tt.wantRe)
			}
		})
	}
}
