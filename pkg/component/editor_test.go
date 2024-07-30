package component

import (
	"reflect"
	"testing"

	ut "github.com/nervatura/component/pkg/util"
)

func TestTestEditor(t *testing.T) {
	for _, tt := range TestEditor(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	testEditorResponse(ResponseEvent{Trigger: &Editor{}})
}

func TestEditor_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Title         string
		Icon          string
		View          string
		Views         []EditorView
		Rows          []Row
		Tables        []Table
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
				propName: "title",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edi := &Editor{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Icon:          tt.fields.Icon,
				View:          tt.fields.View,
				Views:         tt.fields.Views,
				Rows:          tt.fields.Rows,
				Tables:        tt.fields.Tables,
			}
			if got := edi.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Editor.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEditor_Validation(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Title         string
		Icon          string
		View          string
		Views         []EditorView
		Rows          []Row
		Tables        []Table
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
			name: "invalid_views",
			args: args{
				propName:  "views",
				propValue: "",
			},
			want: []EditorView{},
		},
		{
			name: "invalid_rows",
			args: args{
				propName:  "rows",
				propValue: "",
			},
			want: []Row{},
		},
		{
			name: "invalid_tables",
			args: args{
				propName:  "tables",
				propValue: "",
			},
			want: []Table{},
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
			edi := &Editor{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Icon:          tt.fields.Icon,
				View:          tt.fields.View,
				Views:         tt.fields.Views,
				Rows:          tt.fields.Rows,
				Tables:        tt.fields.Tables,
			}
			if got := edi.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Editor.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEditor_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Title         string
		Icon          string
		View          string
		Views         []EditorView
		Rows          []Row
		Tables        []Table
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
			edi := &Editor{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Icon:          tt.fields.Icon,
				View:          tt.fields.View,
				Views:         tt.fields.Views,
				Rows:          tt.fields.Rows,
				Tables:        tt.fields.Tables,
			}
			if got := edi.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Editor.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEditor_response(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Title         string
		Icon          string
		View          string
		Views         []EditorView
		Rows          []Row
		Tables        []Table
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
			name: "tab_btn",
			fields: fields{
				BaseComponent: BaseComponent{
					OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
						evt.Trigger = &Editor{}
						return evt
					},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "tab_btn",
					Trigger: &BaseComponent{
						Data: ut.IM{"key": "value"},
					},
				},
			},
		},
		{
			name: "view_table",
			args: args{
				evt: ResponseEvent{
					TriggerName: "view_table",
					Name:        TableEventEditCell,
					Trigger: &BaseComponent{
						Data: ut.IM{"key": "value"},
					},
				},
			},
		},
		{
			name: "view_table_sort",
			args: args{
				evt: ResponseEvent{
					TriggerName: "view_table",
					Name:        TableEventSort,
					Trigger: &BaseComponent{
						Data: ut.IM{"key": "value"},
					},
				},
			},
		},
		{
			name: "default",
			args: args{
				evt: ResponseEvent{
					TriggerName: "view_row",
					Trigger: &BaseComponent{
						Data: ut.IM{"key": "value"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edi := &Editor{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Icon:          tt.fields.Icon,
				View:          tt.fields.View,
				Views:         tt.fields.Views,
				Rows:          tt.fields.Rows,
				Tables:        tt.fields.Tables,
			}
			edi.response(tt.args.evt)
		})
	}
}
