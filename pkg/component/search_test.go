package component

import (
	"reflect"
	"testing"

	ut "github.com/nervatura/component/pkg/util"
)

func TestTestSearch(t *testing.T) {
	for _, tt := range TestSearch(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	testSearchResponse(ResponseEvent{Name: SearchEventSearch, Trigger: &Search{}})
	testSearchResponse(ResponseEvent{
		Name: SearchEventSelected, Trigger: &Search{},
		Value: ut.IM{
			"row": ut.IM{},
		},
	})
}

func TestSearch_Validation(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Rows              []ut.IM
		Fields            []TableField
		Title             string
		FilterPlaceholder string
		AutoFocus         bool
		Full              bool
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
			name: "rows",
			args: args{
				propName:  "rows",
				propValue: "",
			},
			want: []ut.IM{},
		},
		{
			name: "fields",
			fields: fields{
				Rows: []ut.IM{{"field": "value"}},
			},
			args: args{
				propName:  "fields",
				propValue: "",
			},
			want: []TableField{{
				Name: "field", FieldType: TableFieldTypeString, Label: "field",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sea := &Search{
				BaseComponent:     tt.fields.BaseComponent,
				Rows:              tt.fields.Rows,
				Fields:            tt.fields.Fields,
				Title:             tt.fields.Title,
				FilterPlaceholder: tt.fields.FilterPlaceholder,
				AutoFocus:         tt.fields.AutoFocus,
				Full:              tt.fields.Full,
			}
			if got := sea.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearch_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Rows              []ut.IM
		Fields            []TableField
		Title             string
		FilterPlaceholder string
		AutoFocus         bool
		Full              bool
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
			sea := &Search{
				BaseComponent:     tt.fields.BaseComponent,
				Rows:              tt.fields.Rows,
				Fields:            tt.fields.Fields,
				Title:             tt.fields.Title,
				FilterPlaceholder: tt.fields.FilterPlaceholder,
				AutoFocus:         tt.fields.AutoFocus,
				Full:              tt.fields.Full,
			}
			if got := sea.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearch_response(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Rows              []ut.IM
		Fields            []TableField
		Title             string
		FilterPlaceholder string
		AutoFocus         bool
		Full              bool
	}
	type args struct {
		evt ResponseEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "search_result",
			args: args{
				evt: ResponseEvent{
					TriggerName: "search_result",
				},
			},
		},
		{
			name: "search_result",
			args: args{
				evt: ResponseEvent{
					Name:        TableEventRowSelected,
					TriggerName: "search_result",
				},
			},
		},
		{
			name: "filter_value",
			args: args{
				evt: ResponseEvent{
					TriggerName: "filter_value",
				},
			},
		},
		{
			name: "btn_search",
			fields: fields{
				BaseComponent: BaseComponent{
					OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
						evt.Trigger = &Search{}
						return evt
					},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "btn_search",
				},
			},
		},
		{
			name: "invalid",
			args: args{
				evt: ResponseEvent{
					Trigger: &Table{
						BaseComponent: BaseComponent{
							Data: ut.IM{},
						},
					},
					TriggerName: "invalid",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sea := &Search{
				BaseComponent:     tt.fields.BaseComponent,
				Rows:              tt.fields.Rows,
				Fields:            tt.fields.Fields,
				Title:             tt.fields.Title,
				FilterPlaceholder: tt.fields.FilterPlaceholder,
				AutoFocus:         tt.fields.AutoFocus,
				Full:              tt.fields.Full,
			}
			sea.response(tt.args.evt)
		})
	}
}
