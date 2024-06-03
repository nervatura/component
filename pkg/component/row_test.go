package component

import (
	"reflect"
	"testing"
)

func TestTestRow(t *testing.T) {
	for _, tt := range TestRow(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
}

func TestRow_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Columns       []RowColumn
		Full          bool
		BorderTop     bool
		BorderBottom  bool
		FieldCol      bool
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
				propName: "border_top",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := &Row{
				BaseComponent: tt.fields.BaseComponent,
				Columns:       tt.fields.Columns,
				Full:          tt.fields.Full,
				BorderTop:     tt.fields.BorderTop,
				BorderBottom:  tt.fields.BorderBottom,
				FieldCol:      tt.fields.FieldCol,
			}
			if got := row.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Row.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRow_Validation(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Columns       []RowColumn
		Full          bool
		BorderTop     bool
		BorderBottom  bool
		FieldCol      bool
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
			name: "invalid_columns",
			args: args{
				propName:  "columns",
				propValue: "",
			},
			want: []RowColumn{},
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
			row := &Row{
				BaseComponent: tt.fields.BaseComponent,
				Columns:       tt.fields.Columns,
				Full:          tt.fields.Full,
				BorderTop:     tt.fields.BorderTop,
				BorderBottom:  tt.fields.BorderBottom,
				FieldCol:      tt.fields.FieldCol,
			}
			if got := row.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Row.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRow_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Columns       []RowColumn
		Full          bool
		BorderTop     bool
		BorderBottom  bool
		FieldCol      bool
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
			row := &Row{
				BaseComponent: tt.fields.BaseComponent,
				Columns:       tt.fields.Columns,
				Full:          tt.fields.Full,
				BorderTop:     tt.fields.BorderTop,
				BorderBottom:  tt.fields.BorderBottom,
				FieldCol:      tt.fields.FieldCol,
			}
			if got := row.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Row.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}
