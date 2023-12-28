package atom

import (
	"reflect"
	"testing"

	bc "github.com/nervatura/component/component/base"
)

func TestDemoSelect(t *testing.T) {
	for _, tt := range DemoSelect("/demo", "") {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
}

func TestSelect_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Value         string
		Options       []SelectOption
		IsNull        bool
		Label         string
		Disabled      bool
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
			sel := &Select{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Options:       tt.fields.Options,
				IsNull:        tt.fields.IsNull,
				Label:         tt.fields.Label,
				Disabled:      tt.fields.Disabled,
				AutoFocus:     tt.fields.AutoFocus,
				Full:          tt.fields.Full,
			}
			if got := sel.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Select.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelect_Validation(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Value         string
		Options       []SelectOption
		IsNull        bool
		Label         string
		Disabled      bool
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
			name: "options",
			args: args{
				propName:  "options",
				propValue: "",
			},
			want: []SelectOption{},
		},
		{
			name: "isnull",
			fields: fields{
				IsNull: true,
			},
			args: args{
				propName:  "value",
				propValue: "",
			},
			want: "",
		},
		{
			name: "invalid_value",
			fields: fields{
				Options: []SelectOption{},
			},
			args: args{
				propName:  "value",
				propValue: "abc",
			},
			want: "",
		},
		{
			name: "invalid_option",
			fields: fields{
				Options: []SelectOption{
					{Value: "aaa", Text: "aaa"},
				},
			},
			args: args{
				propName:  "value",
				propValue: "abc",
			},
			want: "aaa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel := &Select{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Options:       tt.fields.Options,
				IsNull:        tt.fields.IsNull,
				Label:         tt.fields.Label,
				Disabled:      tt.fields.Disabled,
				AutoFocus:     tt.fields.AutoFocus,
				Full:          tt.fields.Full,
			}
			if got := sel.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Select.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelect_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Value         string
		Options       []SelectOption
		IsNull        bool
		Label         string
		Disabled      bool
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
			sel := &Select{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Options:       tt.fields.Options,
				IsNull:        tt.fields.IsNull,
				Label:         tt.fields.Label,
				Disabled:      tt.fields.Disabled,
				AutoFocus:     tt.fields.AutoFocus,
				Full:          tt.fields.Full,
			}
			if got := sel.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Select.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelect_OnRequest(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Value         string
		Options       []SelectOption
		IsNull        bool
		Label         string
		Disabled      bool
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
				Trigger: &Select{},
				Name:    SelectEventChange,
				Value:   "",
			},
		},
		{
			name: "OnResponse",
			fields: fields{
				BaseComponent: bc.BaseComponent{
					OnResponse: func(evt bc.ResponseEvent) (re bc.ResponseEvent) {
						evt.Trigger = &Select{}
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
				Trigger: &Select{},
				Name:    SelectEventChange,
				Value:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel := &Select{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Options:       tt.fields.Options,
				IsNull:        tt.fields.IsNull,
				Label:         tt.fields.Label,
				Disabled:      tt.fields.Disabled,
				AutoFocus:     tt.fields.AutoFocus,
				Full:          tt.fields.Full,
			}
			if gotRe := sel.OnRequest(tt.args.te); !reflect.DeepEqual(gotRe, tt.wantRe) {
				t.Errorf("Select.OnRequest() = %v, want %v", gotRe, tt.wantRe)
			}
		})
	}
}
