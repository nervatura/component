package atom

import (
	"reflect"
	"testing"

	bc "github.com/nervatura/component/component/base"
)

func TestDemoToast(t *testing.T) {
	for _, tt := range DemoToast("/demo", "") {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	demoToastResponse(bc.ResponseEvent{Trigger: &Toast{}})
}

func TestToast_Render(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Type          string
		Value         string
		Timeout       int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "ok",
			fields:  fields{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tst := &Toast{
				BaseComponent: tt.fields.BaseComponent,
				Type:          tt.fields.Type,
				Value:         tt.fields.Value,
				Timeout:       tt.fields.Timeout,
			}
			_, err := tst.Render()
			if (err != nil) != tt.wantErr {
				t.Errorf("Toast.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestToast_Validation(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Type          string
		Value         string
		Timeout       int64
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
			tst := &Toast{
				BaseComponent: tt.fields.BaseComponent,
				Type:          tt.fields.Type,
				Value:         tt.fields.Value,
				Timeout:       tt.fields.Timeout,
			}
			if got := tst.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Toast.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToast_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Type          string
		Value         string
		Timeout       int64
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
			tst := &Toast{
				BaseComponent: tt.fields.BaseComponent,
				Type:          tt.fields.Type,
				Value:         tt.fields.Value,
				Timeout:       tt.fields.Timeout,
			}
			if got := tst.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Toast.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}
