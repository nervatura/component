package component

import (
	"reflect"
	"testing"
)

func TestTestDateTime(t *testing.T) {
	for _, tt := range TestDateTime(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
}

func TestDateTime_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Type          string
		Value         string
		Label         string
		IsNull        bool
		Picker        bool
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
			dti := &DateTime{
				BaseComponent: tt.fields.BaseComponent,
				Type:          tt.fields.Type,
				Value:         tt.fields.Value,
				Label:         tt.fields.Label,
				IsNull:        tt.fields.IsNull,
				Picker:        tt.fields.Picker,
				Disabled:      tt.fields.Disabled,
				ReadOnly:      tt.fields.ReadOnly,
				AutoFocus:     tt.fields.AutoFocus,
				Full:          tt.fields.Full,
			}
			if got := dti.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateTime.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateTime_Validation(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Type          string
		Value         string
		Label         string
		IsNull        bool
		Picker        bool
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
			name: "datetime",
			fields: fields{
				Type: DateTimeTypeDate,
			},
			args: args{
				propName:  "value",
				propValue: "2023-01-01T",
			},
			want: "2023-01-01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dti := &DateTime{
				BaseComponent: tt.fields.BaseComponent,
				Type:          tt.fields.Type,
				Value:         tt.fields.Value,
				Label:         tt.fields.Label,
				IsNull:        tt.fields.IsNull,
				Picker:        tt.fields.Picker,
				Disabled:      tt.fields.Disabled,
				ReadOnly:      tt.fields.ReadOnly,
				AutoFocus:     tt.fields.AutoFocus,
				Full:          tt.fields.Full,
			}
			if got := dti.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateTime.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateTime_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Type          string
		Value         string
		Label         string
		IsNull        bool
		Picker        bool
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
			dti := &DateTime{
				BaseComponent: tt.fields.BaseComponent,
				Type:          tt.fields.Type,
				Value:         tt.fields.Value,
				Label:         tt.fields.Label,
				IsNull:        tt.fields.IsNull,
				Picker:        tt.fields.Picker,
				Disabled:      tt.fields.Disabled,
				ReadOnly:      tt.fields.ReadOnly,
				AutoFocus:     tt.fields.AutoFocus,
				Full:          tt.fields.Full,
			}
			if got := dti.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateTime.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateTime_OnRequest(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Type          string
		Value         string
		Label         string
		IsNull        bool
		Picker        bool
		Disabled      bool
		ReadOnly      bool
		AutoFocus     bool
		Full          bool
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
						evt.Trigger = &DateTime{
							Value:  "",
							IsNull: true,
						}
						return evt
					},
					Swap: SwapInnerHTML,
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
			dti := &DateTime{
				BaseComponent: tt.fields.BaseComponent,
				Type:          tt.fields.Type,
				Value:         tt.fields.Value,
				Label:         tt.fields.Label,
				IsNull:        tt.fields.IsNull,
				Picker:        tt.fields.Picker,
				Disabled:      tt.fields.Disabled,
				ReadOnly:      tt.fields.ReadOnly,
				AutoFocus:     tt.fields.AutoFocus,
				Full:          tt.fields.Full,
			}
			dti.OnRequest(tt.args.te)
		})
	}
}
