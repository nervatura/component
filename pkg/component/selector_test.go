package component

import (
	"reflect"
	"testing"

	ut "github.com/nervatura/component/pkg/util"
)

func TestTestSelector(t *testing.T) {
	for _, tt := range TestSelector(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	testSelectorResponse(ResponseEvent{Name: SelectorEventSearch, Trigger: &Selector{}})
	testSelectorResponse(ResponseEvent{
		Name: SelectorEventLink, Trigger: &Selector{},
		Value: SelectOption{},
	})
	testSelectorResponse(ResponseEvent{
		Name: SelectorEventSelected, Trigger: &Selector{},
		Value: ut.IM{
			"row": ut.IM{},
		},
	})
}

func TestSelector_Validation(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Value             SelectOption
		Rows              []ut.IM
		Fields            []TableField
		Title             string
		FilterPlaceholder string
		FilterValue       string
		Link              bool
		IsNull            bool
		Disabled          bool
		AutoFocus         bool
		Full              bool
		ShowModal         bool
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
			name: "value",
			args: args{
				propName:  "value",
				propValue: "",
			},
			want: SelectOption{},
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
			sel := &Selector{
				BaseComponent:     tt.fields.BaseComponent,
				Value:             tt.fields.Value,
				Rows:              tt.fields.Rows,
				Fields:            tt.fields.Fields,
				Title:             tt.fields.Title,
				FilterPlaceholder: tt.fields.FilterPlaceholder,
				FilterValue:       tt.fields.FilterValue,
				Link:              tt.fields.Link,
				IsNull:            tt.fields.IsNull,
				Disabled:          tt.fields.Disabled,
				AutoFocus:         tt.fields.AutoFocus,
				Full:              tt.fields.Full,
				ShowModal:         tt.fields.ShowModal,
			}
			if got := sel.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Selector.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelector_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Value             SelectOption
		Rows              []ut.IM
		Fields            []TableField
		Title             string
		FilterPlaceholder string
		FilterValue       string
		Link              bool
		IsNull            bool
		Disabled          bool
		AutoFocus         bool
		Full              bool
		ShowModal         bool
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
			sel := &Selector{
				BaseComponent:     tt.fields.BaseComponent,
				Value:             tt.fields.Value,
				Rows:              tt.fields.Rows,
				Fields:            tt.fields.Fields,
				Title:             tt.fields.Title,
				FilterPlaceholder: tt.fields.FilterPlaceholder,
				FilterValue:       tt.fields.FilterValue,
				Link:              tt.fields.Link,
				IsNull:            tt.fields.IsNull,
				Disabled:          tt.fields.Disabled,
				AutoFocus:         tt.fields.AutoFocus,
				Full:              tt.fields.Full,
				ShowModal:         tt.fields.ShowModal,
			}
			if got := sel.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Selector.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelector_response(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Value             SelectOption
		Rows              []ut.IM
		Fields            []TableField
		Title             string
		FilterPlaceholder string
		FilterValue       string
		Link              bool
		IsNull            bool
		Disabled          bool
		AutoFocus         bool
		Full              bool
		ShowModal         bool
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
			name: "selector_result",
			args: args{
				evt: ResponseEvent{
					TriggerName: "selector_result",
				},
			},
		},
		{
			name: "selector_result",
			args: args{
				evt: ResponseEvent{
					Name:        TableEventRowSelected,
					TriggerName: "selector_result",
				},
			},
		},
		{
			name: "btn_delete",
			args: args{
				evt: ResponseEvent{
					TriggerName: "btn_delete",
				},
			},
		},
		{
			name: "btn_modal",
			args: args{
				evt: ResponseEvent{
					TriggerName: "btn_modal",
				},
			},
		},
		{
			name: "btn_close",
			args: args{
				evt: ResponseEvent{
					TriggerName: "btn_close",
				},
			},
		},
		{
			name: "filter",
			args: args{
				evt: ResponseEvent{
					TriggerName: "filter",
				},
			},
		},
		{
			name: "selector_text",
			args: args{
				evt: ResponseEvent{
					TriggerName: "selector_text",
				},
			},
		},
		{
			name: "btn_search",
			fields: fields{
				BaseComponent: BaseComponent{
					OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
						evt.Trigger = &Selector{}
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
			sel := &Selector{
				BaseComponent:     tt.fields.BaseComponent,
				Value:             tt.fields.Value,
				Rows:              tt.fields.Rows,
				Fields:            tt.fields.Fields,
				Title:             tt.fields.Title,
				FilterPlaceholder: tt.fields.FilterPlaceholder,
				FilterValue:       tt.fields.FilterValue,
				Link:              tt.fields.Link,
				IsNull:            tt.fields.IsNull,
				Disabled:          tt.fields.Disabled,
				AutoFocus:         tt.fields.AutoFocus,
				Full:              tt.fields.Full,
				ShowModal:         tt.fields.ShowModal,
			}
			sel.response(tt.args.evt)
		})
	}
}

func TestSelector_getComponent(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Value             SelectOption
		Rows              []ut.IM
		Fields            []TableField
		Title             string
		FilterPlaceholder string
		FilterValue       string
		Link              bool
		IsNull            bool
		Disabled          bool
		AutoFocus         bool
		Full              bool
		ShowModal         bool
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "btn_search",
			args: args{
				name: "btn_search",
			},
		},
		{
			name: "btn_close",
			args: args{
				name: "btn_close",
			},
		},
		{
			name: "filter",
			args: args{
				name: "filter",
			},
		},
		{
			name: "selector_result",
			args: args{
				name: "selector_result",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel := &Selector{
				BaseComponent:     tt.fields.BaseComponent,
				Value:             tt.fields.Value,
				Rows:              tt.fields.Rows,
				Fields:            tt.fields.Fields,
				Title:             tt.fields.Title,
				FilterPlaceholder: tt.fields.FilterPlaceholder,
				FilterValue:       tt.fields.FilterValue,
				Link:              tt.fields.Link,
				IsNull:            tt.fields.IsNull,
				Disabled:          tt.fields.Disabled,
				AutoFocus:         tt.fields.AutoFocus,
				Full:              tt.fields.Full,
				ShowModal:         tt.fields.ShowModal,
			}
			_, err := sel.getComponent(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Selector.getComponent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
