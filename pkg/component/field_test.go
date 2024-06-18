package component

import (
	"reflect"
	"testing"

	ut "github.com/nervatura/component/pkg/util"
)

func TestTestField(t *testing.T) {
	for _, tt := range TestField(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	testFieldResponse(ResponseEvent{Trigger: &Button{}, TriggerName: "button"})
	testFieldResponse(ResponseEvent{Trigger: &Label{}, TriggerName: "link"})
	testFieldResponse(ResponseEvent{
		Trigger: &Selector{}, TriggerName: "selector", Name: SelectorEventLink, Value: SelectOption{}})
	testFieldResponse(ResponseEvent{
		Trigger: &Selector{}, TriggerName: "selector", Name: SelectorEventSearch, Value: SelectOption{}})
	testFieldResponse(ResponseEvent{
		Trigger: &Selector{}, TriggerName: "selector", Name: SelectorEventSelected,
		Value: ut.IM{"row": ut.IM{}}})
	testFieldResponse(ResponseEvent{Trigger: &List{}, TriggerName: "list", Value: ut.IM{"row": ut.IM{}}})
}

func TestField_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Type          string
		Value         ut.IM
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
				propName: "type",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fld := &Field{
				BaseComponent: tt.fields.BaseComponent,
				Type:          tt.fields.Type,
				Value:         tt.fields.Value,
			}
			if got := fld.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Field.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_Validation(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Type          string
		Value         ut.IM
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
			fld := &Field{
				BaseComponent: tt.fields.BaseComponent,
				Type:          tt.fields.Type,
				Value:         tt.fields.Value,
			}
			if got := fld.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Field.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Type          string
		Value         ut.IM
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
			fld := &Field{
				BaseComponent: tt.fields.BaseComponent,
				Type:          tt.fields.Type,
				Value:         tt.fields.Value,
			}
			if got := fld.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Field.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}
