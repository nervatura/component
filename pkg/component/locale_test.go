package component

import (
	"reflect"
	"testing"

	ut "github.com/nervatura/component/pkg/util"
)

func TestTestLocale(t *testing.T) {
	for _, tt := range TestLocale(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	testLocaleResponse(ResponseEvent{Name: LocalesEventError, Trigger: &Locale{}})
	testLocaleResponse(ResponseEvent{Name: LocalesEventSave, Trigger: &Locale{}})
	testLocaleResponse(ResponseEvent{Name: LocalesEventUndo, Trigger: &Locale{}})
}

func TestLocale_Validation(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Locales       []SelectOption
		TagKeys       []SelectOption
		FilterValue   string
		Dirty         bool
		AddItem       bool
		Labels        ut.SM
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
			name: "locales",
			args: args{
				propName:  "locales",
				propValue: []SelectOption{},
			},
			want: []SelectOption{{Value: "default", Text: "default"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc := &Locale{
				BaseComponent: tt.fields.BaseComponent,
				Locales:       tt.fields.Locales,
				TagKeys:       tt.fields.TagKeys,
				FilterValue:   tt.fields.FilterValue,
				Dirty:         tt.fields.Dirty,
				AddItem:       tt.fields.AddItem,
				Labels:        tt.fields.Labels,
			}
			if got := loc.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Locale.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocale_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Locales       []SelectOption
		TagKeys       []SelectOption
		FilterValue   string
		Dirty         bool
		AddItem       bool
		Labels        ut.SM
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
			loc := &Locale{
				BaseComponent: tt.fields.BaseComponent,
				Locales:       tt.fields.Locales,
				TagKeys:       tt.fields.TagKeys,
				FilterValue:   tt.fields.FilterValue,
				Dirty:         tt.fields.Dirty,
				AddItem:       tt.fields.AddItem,
				Labels:        tt.fields.Labels,
			}
			if got := loc.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Locale.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocale_response(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Locales       []SelectOption
		TagKeys       []SelectOption
		FilterValue   string
		Dirty         bool
		AddItem       bool
		Labels        ut.SM
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
			name: "values",
			args: args{
				evt: ResponseEvent{
					TriggerName: "values",
				},
			},
		},
		{
			name: "tag_keys",
			args: args{
				evt: ResponseEvent{
					TriggerName: "tag_keys",
				},
			},
		},
		{
			name: "locales",
			fields: fields{
				TagKeys: []SelectOption{{Value: "address", Text: "address"}},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "locales",
				},
			},
		},
		{
			name: "undo",
			args: args{
				evt: ResponseEvent{
					TriggerName: "undo",
				},
			},
		},
		{
			name: "update",
			args: args{
				evt: ResponseEvent{
					TriggerName: "update",
				},
			},
		},
		{
			name: "add_locales_missing",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{
						"locfile": ut.IM{"locales": ut.IM{}},
					},
					OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
						evt.Trigger = &BaseComponent{}
						return evt
					},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "add",
				},
			},
		},
		{
			name: "add_existing_lang",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{
						"lang_key": "en",
						"locfile":  ut.IM{"locales": ut.IM{}},
					},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "add",
				},
			},
		},
		{
			name: "add",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{
						"lang_key":  "de",
						"lang_name": "de",
						"locfile":   ut.IM{"locales": ut.IM{}},
					},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "add",
				},
			},
		},
		{
			name: "missing",
			args: args{
				evt: ResponseEvent{
					TriggerName: "missing",
				},
			},
		},
		{
			name: "tag_cell",
			args: args{
				evt: ResponseEvent{
					TriggerName: "tag_cell",
				},
			},
		},
		{
			name: "value_cell",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{
						"locales": "de",
						"locfile": ut.IM{"locales": ut.IM{"de": ut.IM{}}},
					},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "value_cell",
					Trigger: &Input{
						BaseComponent: BaseComponent{
							Data: ut.IM{
								"key": "abc",
							},
						},
					},
				},
			},
		},
		{
			name: "lang_key",
			args: args{
				evt: ResponseEvent{
					TriggerName: "lang_key",
				},
			},
		},
		{
			name: "lang_name",
			args: args{
				evt: ResponseEvent{
					TriggerName: "lang_name",
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
			name: "add_item",
			args: args{
				evt: ResponseEvent{
					TriggerName: "add_item",
				},
			},
		},
		{
			name: "default",
			args: args{
				evt: ResponseEvent{
					TriggerName: "default",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc := &Locale{
				BaseComponent: tt.fields.BaseComponent,
				Locales:       tt.fields.Locales,
				TagKeys:       tt.fields.TagKeys,
				FilterValue:   tt.fields.FilterValue,
				Dirty:         tt.fields.Dirty,
				AddItem:       tt.fields.AddItem,
				Labels:        tt.fields.Labels,
			}
			loc.response(tt.args.evt)
		})
	}
}

func TestLocale_getComponent(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Locales       []SelectOption
		TagKeys       []SelectOption
		FilterValue   string
		Dirty         bool
		AddItem       bool
		Labels        ut.SM
	}
	type args struct {
		name string
		data ut.IM
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "tag_keys",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{
						"locales": "de",
					},
				},
			},
			args: args{
				name: "tag_keys",
				data: ut.IM{},
			},
			wantErr: false,
		},
		{
			name: "missing",
			args: args{
				name: "missing",
				data: ut.IM{},
			},
			wantErr: false,
		},
		{
			name: "update",
			args: args{
				name: "update",
				data: ut.IM{},
			},
			wantErr: false,
		},
		{
			name: "undo",
			args: args{
				name: "undo",
				data: ut.IM{},
			},
			wantErr: false,
		},
		{
			name: "add_item",
			fields: fields{
				AddItem: true,
			},
			args: args{
				name: "add_item",
				data: ut.IM{},
			},
			wantErr: false,
		},
		{
			name: "add",
			args: args{
				name: "add",
				data: ut.IM{},
			},
			wantErr: false,
		},
		{
			name: "lang_key",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{
						"lang_key": "key",
					},
				},
			},
			args: args{
				name: "lang_key",
				data: ut.IM{},
			},
			wantErr: false,
		},
		{
			name: "lang_name",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{
						"lang_name": "name",
					},
				},
			},
			args: args{
				name: "lang_name",
				data: ut.IM{},
			},
			wantErr: false,
		},
		{
			name: "tag_cell",
			args: args{
				name: "tag_cell",
				data: ut.IM{},
			},
			wantErr: false,
		},
		{
			name: "value_cell",
			args: args{
				name: "value_cell",
				data: ut.IM{},
			},
			wantErr: false,
		},
		{
			name: "values",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{
						"locales":    "de",
						"tag_keys":   "tag",
						"locfile":    ut.IM{"locales": ut.IM{"de": ut.IM{"tag_key1": "value1"}}},
						"deflang":    ut.IM{"tag_key1": "value1", "tag_key2": "value2"},
						"tag_values": map[string][]string{"tag": {"tag_key1", "tag_key2"}},
					},
				},
			},
			args: args{
				name: "values",
				data: ut.IM{},
			},
			wantErr: false,
		},
		{
			name: "values_client",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{
						"locales":    "client",
						"tag_keys":   "tag",
						"locfile":    ut.IM{"locales": ut.IM{"de": ut.IM{"tag_key1": "value1"}}},
						"deflang":    ut.IM{"tag_key1": "value1", "tag_key2": "value2"},
						"tag_values": map[string][]string{"tag": {"tag_key1", "tag_key2"}},
					},
				},
			},
			args: args{
				name: "values",
				data: ut.IM{},
			},
			wantErr: false,
		},
		{
			name: "values_missing",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{
						"locales":    "de",
						"tag_keys":   "missing",
						"locfile":    ut.IM{"locales": ut.IM{"de": ut.IM{"tag_key1": "value1"}}},
						"deflang":    ut.IM{"tag_key1": "value1", "tag_key2": "value2"},
						"tag_values": map[string][]string{"tag": {"tag_key1", "tag_key2"}},
					},
				},
			},
			args: args{
				name: "values",
				data: ut.IM{},
			},
			wantErr: false,
		},
		{
			name: "values_filter",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{
						"locales":    "de",
						"tag_keys":   "tag",
						"locfile":    ut.IM{"locales": ut.IM{"de": ut.IM{"tag_key1": "value1", "tag_key2": "key1"}}},
						"deflang":    ut.IM{"tag_key1": "value1", "tag_key2": "value2"},
						"tag_values": map[string][]string{"tag": {"tag_key1", "tag_key2"}},
					},
				},
				FilterValue: "key1",
			},
			args: args{
				name: "values",
				data: ut.IM{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc := &Locale{
				BaseComponent: tt.fields.BaseComponent,
				Locales:       tt.fields.Locales,
				TagKeys:       tt.fields.TagKeys,
				FilterValue:   tt.fields.FilterValue,
				Dirty:         tt.fields.Dirty,
				AddItem:       tt.fields.AddItem,
				Labels:        tt.fields.Labels,
			}
			_, err := loc.getComponent(tt.args.name, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Locale.getComponent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
