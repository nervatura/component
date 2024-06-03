package component

import (
	"reflect"
	"testing"
)

func TestTestUpload(t *testing.T) {
	for _, tt := range TestUpload(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
}

func TestUpload_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Accept        string
		Placeholder   string
		ToastMessage  string
		Disabled      bool
		MaxLength     int64
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
			upl := &Upload{
				BaseComponent: tt.fields.BaseComponent,
				Accept:        tt.fields.Accept,
				Placeholder:   tt.fields.Placeholder,
				ToastMessage:  tt.fields.ToastMessage,
				Disabled:      tt.fields.Disabled,
				MaxLength:     tt.fields.MaxLength,
				Full:          tt.fields.Full,
			}
			if got := upl.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Upload.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpload_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Accept        string
		Placeholder   string
		ToastMessage  string
		Disabled      bool
		MaxLength     int64
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
			upl := &Upload{
				BaseComponent: tt.fields.BaseComponent,
				Accept:        tt.fields.Accept,
				Placeholder:   tt.fields.Placeholder,
				ToastMessage:  tt.fields.ToastMessage,
				Disabled:      tt.fields.Disabled,
				MaxLength:     tt.fields.MaxLength,
				Full:          tt.fields.Full,
			}
			if got := upl.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Upload.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpload_OnRequest(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Accept        string
		Placeholder   string
		ToastMessage  string
		Disabled      bool
		MaxLength     int64
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
						evt.Trigger = &Upload{}
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
			upl := &Upload{
				BaseComponent: tt.fields.BaseComponent,
				Accept:        tt.fields.Accept,
				Placeholder:   tt.fields.Placeholder,
				ToastMessage:  tt.fields.ToastMessage,
				Disabled:      tt.fields.Disabled,
				MaxLength:     tt.fields.MaxLength,
				Full:          tt.fields.Full,
			}
			upl.OnRequest(tt.args.te)
		})
	}
}
