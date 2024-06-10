package component

import (
	"reflect"
	"testing"

	ut "github.com/nervatura/component/pkg/util"
)

func TestTestBrowser(t *testing.T) {
	for _, tt := range TestBrowser(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	testBrowserResponse(ResponseEvent{Name: BrowserEventSearch, Trigger: &Browser{}})
	testBrowserResponse(ResponseEvent{Name: BrowserEventView, Trigger: &Browser{}, Value: "meta"})
}

func TestBrowser_GetProperty(t *testing.T) {
	type fields struct {
		Table          Table
		View           string
		Views          []SelectOption
		HideHeader     bool
		ShowDropdown   bool
		ShowColumns    bool
		ShowTotal      bool
		HideBookmark   bool
		HideExport     bool
		HideHelp       bool
		ReadOnly       bool
		VisibleColumns map[string]bool
		Filters        []BrowserFilter
		MetaFields     map[string]BrowserMetaField
		Labels         ut.SM
		totalFields    []BrowserTotalField
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
				propName: "current_page",
			},
			want: int64(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bro := &Browser{
				Table:          tt.fields.Table,
				View:           tt.fields.View,
				Views:          tt.fields.Views,
				HideHeader:     tt.fields.HideHeader,
				ShowDropdown:   tt.fields.ShowDropdown,
				ShowColumns:    tt.fields.ShowColumns,
				ShowTotal:      tt.fields.ShowTotal,
				HideBookmark:   tt.fields.HideBookmark,
				HideExport:     tt.fields.HideExport,
				HideHelp:       tt.fields.HideHelp,
				ReadOnly:       tt.fields.ReadOnly,
				VisibleColumns: tt.fields.VisibleColumns,
				Filters:        tt.fields.Filters,
				MetaFields:     tt.fields.MetaFields,
				Labels:         tt.fields.Labels,
				totalFields:    tt.fields.totalFields,
			}
			if got := bro.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Browser.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrowser_Validation(t *testing.T) {
	type fields struct {
		Table          Table
		View           string
		Views          []SelectOption
		HideHeader     bool
		ShowDropdown   bool
		ShowColumns    bool
		ShowTotal      bool
		HideBookmark   bool
		HideExport     bool
		HideHelp       bool
		ReadOnly       bool
		VisibleColumns map[string]bool
		Filters        []BrowserFilter
		MetaFields     map[string]BrowserMetaField
		Labels         ut.SM
		totalFields    []BrowserTotalField
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
			name: "filters",
			args: args{
				propName:  "filters",
				propValue: nil,
			},
			want: []BrowserFilter{},
		},
		{
			name: "visible_columns",
			args: args{
				propName:  "visible_columns",
				propValue: []map[string]bool{{"fieldName": true}},
			},
			want: map[string]bool{"fieldName": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bro := &Browser{
				Table:          tt.fields.Table,
				View:           tt.fields.View,
				Views:          tt.fields.Views,
				HideHeader:     tt.fields.HideHeader,
				ShowDropdown:   tt.fields.ShowDropdown,
				ShowColumns:    tt.fields.ShowColumns,
				ShowTotal:      tt.fields.ShowTotal,
				HideBookmark:   tt.fields.HideBookmark,
				HideExport:     tt.fields.HideExport,
				HideHelp:       tt.fields.HideHelp,
				ReadOnly:       tt.fields.ReadOnly,
				VisibleColumns: tt.fields.VisibleColumns,
				Filters:        tt.fields.Filters,
				MetaFields:     tt.fields.MetaFields,
				Labels:         tt.fields.Labels,
				totalFields:    tt.fields.totalFields,
			}
			if got := bro.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Browser.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrowser_SetProperty(t *testing.T) {
	type fields struct {
		Table          Table
		View           string
		Views          []SelectOption
		HideHeader     bool
		ShowDropdown   bool
		ShowColumns    bool
		ShowTotal      bool
		HideBookmark   bool
		HideExport     bool
		HideHelp       bool
		ReadOnly       bool
		VisibleColumns map[string]bool
		Filters        []BrowserFilter
		MetaFields     map[string]BrowserMetaField
		Labels         ut.SM
		totalFields    []BrowserTotalField
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
			bro := &Browser{
				Table:          tt.fields.Table,
				View:           tt.fields.View,
				Views:          tt.fields.Views,
				HideHeader:     tt.fields.HideHeader,
				ShowDropdown:   tt.fields.ShowDropdown,
				ShowColumns:    tt.fields.ShowColumns,
				ShowTotal:      tt.fields.ShowTotal,
				HideBookmark:   tt.fields.HideBookmark,
				HideExport:     tt.fields.HideExport,
				HideHelp:       tt.fields.HideHelp,
				ReadOnly:       tt.fields.ReadOnly,
				VisibleColumns: tt.fields.VisibleColumns,
				Filters:        tt.fields.Filters,
				MetaFields:     tt.fields.MetaFields,
				Labels:         tt.fields.Labels,
				totalFields:    tt.fields.totalFields,
			}
			if got := bro.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Browser.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrowser_exportData(t *testing.T) {
	type fields struct {
		Table          Table
		View           string
		Views          []SelectOption
		HideHeader     bool
		ShowDropdown   bool
		ShowColumns    bool
		ShowTotal      bool
		HideBookmark   bool
		HideExport     bool
		ExportLimit    int64
		HideHelp       bool
		ReadOnly       bool
		VisibleColumns map[string]bool
		Filters        []BrowserFilter
		MetaFields     map[string]BrowserMetaField
		Labels         ut.SM
		totalFields    []BrowserTotalField
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "export",
			fields: fields{
				Table: Table{
					Fields: []TableField{
						{Name: "field", Label: "Label"},
					},
					Rows: []ut.IM{
						{"field": "value", "missing": "missing"},
						{"field": "value", "missing": "missing"},
					},
				},
				VisibleColumns: map[string]bool{
					"field": true, "missing": true,
				},
				ExportLimit: 2000,
			},
		},
		{
			name: "limit",
			fields: fields{
				Table: Table{
					Fields: []TableField{
						{Name: "field", Label: "Label"},
					},
					Rows: []ut.IM{
						{"field": "value", "missing": "missing"},
						{"field": "value", "missing": "missing"},
					},
				},
				VisibleColumns: map[string]bool{
					"field": true, "missing": true,
				},
				ExportLimit: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bro := &Browser{
				Table:          tt.fields.Table,
				View:           tt.fields.View,
				Views:          tt.fields.Views,
				HideHeader:     tt.fields.HideHeader,
				ShowDropdown:   tt.fields.ShowDropdown,
				ShowColumns:    tt.fields.ShowColumns,
				ShowTotal:      tt.fields.ShowTotal,
				HideBookmark:   tt.fields.HideBookmark,
				HideExport:     tt.fields.HideExport,
				ExportLimit:    tt.fields.ExportLimit,
				HideHelp:       tt.fields.HideHelp,
				ReadOnly:       tt.fields.ReadOnly,
				VisibleColumns: tt.fields.VisibleColumns,
				Filters:        tt.fields.Filters,
				MetaFields:     tt.fields.MetaFields,
				Labels:         tt.fields.Labels,
				totalFields:    tt.fields.totalFields,
			}
			bro.exportData()
		})
	}
}

func TestBrowser_response(t *testing.T) {
	type fields struct {
		Table          Table
		View           string
		Views          []SelectOption
		HideHeader     bool
		ShowDropdown   bool
		ShowColumns    bool
		ShowTotal      bool
		HideBookmark   bool
		HideExport     bool
		HideHelp       bool
		ReadOnly       bool
		VisibleColumns map[string]bool
		Filters        []BrowserFilter
		MetaFields     map[string]BrowserMetaField
		Labels         ut.SM
		totalFields    []BrowserTotalField
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
			name:   "table",
			fields: fields{},
			args: args{
				evt: ResponseEvent{
					TriggerName: "table",
					Name:        TableEventEditCell,
				},
			},
		},
		{
			name: "table_on",
			fields: fields{
				Table: Table{
					BaseComponent: BaseComponent{
						OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
							evt.Trigger = &Browser{}
							return evt
						},
					},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "table",
					Name:        TableEventEditCell,
				},
			},
		},
		{
			name: "btn_export",
			args: args{
				evt: ResponseEvent{
					TriggerName: "btn_export",
				},
			},
		},
		{
			name: "hide_header",
			args: args{
				evt: ResponseEvent{
					TriggerName: "hide_header",
				},
			},
		},
		{
			name: "btn_search",
			fields: fields{
				Table: Table{
					BaseComponent: BaseComponent{
						OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
							evt.Trigger = &Browser{}
							return evt
						},
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
			name: "btn_bookmark",
			args: args{
				evt: ResponseEvent{
					TriggerName: "btn_bookmark",
				},
			},
		},
		{
			name: "btn_help",
			args: args{
				evt: ResponseEvent{
					TriggerName: "btn_help",
				},
			},
		},
		{
			name: "btn_views",
			args: args{
				evt: ResponseEvent{
					TriggerName: "btn_views",
				},
			},
		},
		{
			name: "btn_columns",
			args: args{
				evt: ResponseEvent{
					TriggerName: "btn_columns",
				},
			},
		},
		{
			name: "btn_filter",
			fields: fields{
				Table: Table{
					Fields: []TableField{
						{Name: "field", Label: "Label"},
					},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "btn_filter",
				},
			},
		},
		{
			name: "btn_total",
			args: args{
				evt: ResponseEvent{
					TriggerName: "btn_total",
				},
			},
		},
		{
			name: "btn_ok",
			args: args{
				evt: ResponseEvent{
					TriggerName: "btn_ok",
				},
			},
		},
		{
			name: "edit_row",
			args: args{
				evt: ResponseEvent{
					TriggerName: "edit_row",
					Trigger: &BaseComponent{
						Data: ut.IM{"index": 1},
					},
				},
			},
		},
		{
			name: "menu_item",
			args: args{
				evt: ResponseEvent{
					TriggerName: "menu_item",
					Trigger: &BaseComponent{
						Data: ut.IM{"index": 1},
					},
				},
			},
		},
		{
			name: "col_item",
			args: args{
				evt: ResponseEvent{
					TriggerName: "col_item",
					Trigger: &BaseComponent{
						Data: ut.IM{"index": 1},
					},
				},
			},
		},
		{
			name: "filter_field",
			fields: fields{
				Filters: []BrowserFilter{
					{Field: "field", Comp: "==", Value: "value"},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "filter_field",
					Trigger: &BaseComponent{
						Data: ut.IM{"index": 0},
					},
				},
			},
		},
		{
			name: "filter_comp",
			fields: fields{
				Filters: []BrowserFilter{
					{Field: "field", Comp: "==", Value: "value"},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "filter_comp",
					Trigger: &BaseComponent{
						Data: ut.IM{"index": 0},
					},
				},
			},
		},
		{
			name: "filter_value",
			fields: fields{
				Filters: []BrowserFilter{
					{Field: "field", Comp: "==", Value: "value"},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "filter_value",
					Trigger: &BaseComponent{
						Data: ut.IM{"index": 0},
					},
				},
			},
		},
		{
			name: "filter_delete",
			fields: fields{
				Filters: []BrowserFilter{
					{Field: "field", Comp: "==", Value: "value"},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "filter_delete",
					Trigger: &BaseComponent{
						Data: ut.IM{"index": 0},
					},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bro := &Browser{
				Table:          tt.fields.Table,
				View:           tt.fields.View,
				Views:          tt.fields.Views,
				HideHeader:     tt.fields.HideHeader,
				ShowDropdown:   tt.fields.ShowDropdown,
				ShowColumns:    tt.fields.ShowColumns,
				ShowTotal:      tt.fields.ShowTotal,
				HideBookmark:   tt.fields.HideBookmark,
				HideExport:     tt.fields.HideExport,
				HideHelp:       tt.fields.HideHelp,
				ReadOnly:       tt.fields.ReadOnly,
				VisibleColumns: tt.fields.VisibleColumns,
				Filters:        tt.fields.Filters,
				MetaFields:     tt.fields.MetaFields,
				Labels:         tt.fields.Labels,
				totalFields:    tt.fields.totalFields,
			}
			bro.response(tt.args.evt)
		})
	}
}

func TestBrowser_defaultFilterValue(t *testing.T) {
	type fields struct {
		Table          Table
		View           string
		Views          []SelectOption
		HideHeader     bool
		ShowDropdown   bool
		ShowColumns    bool
		ShowTotal      bool
		HideBookmark   bool
		HideExport     bool
		HideHelp       bool
		ReadOnly       bool
		VisibleColumns map[string]bool
		Filters        []BrowserFilter
		MetaFields     map[string]BrowserMetaField
		Labels         ut.SM
		totalFields    []BrowserTotalField
	}
	type args struct {
		ftype string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "float",
			args: args{
				ftype: TableFieldTypeNumber,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bro := &Browser{
				Table:          tt.fields.Table,
				View:           tt.fields.View,
				Views:          tt.fields.Views,
				HideHeader:     tt.fields.HideHeader,
				ShowDropdown:   tt.fields.ShowDropdown,
				ShowColumns:    tt.fields.ShowColumns,
				ShowTotal:      tt.fields.ShowTotal,
				HideBookmark:   tt.fields.HideBookmark,
				HideExport:     tt.fields.HideExport,
				HideHelp:       tt.fields.HideHelp,
				ReadOnly:       tt.fields.ReadOnly,
				VisibleColumns: tt.fields.VisibleColumns,
				Filters:        tt.fields.Filters,
				MetaFields:     tt.fields.MetaFields,
				Labels:         tt.fields.Labels,
				totalFields:    tt.fields.totalFields,
			}
			if got := bro.defaultFilterValue(tt.args.ftype); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Browser.defaultFilterValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrowser_setTotalValues(t *testing.T) {
	type fields struct {
		Table          Table
		View           string
		Views          []SelectOption
		HideHeader     bool
		ShowDropdown   bool
		ShowColumns    bool
		ShowTotal      bool
		HideBookmark   bool
		HideExport     bool
		HideHelp       bool
		ReadOnly       bool
		VisibleColumns map[string]bool
		Filters        []BrowserFilter
		MetaFields     map[string]BrowserMetaField
		Labels         ut.SM
		totalFields    []BrowserTotalField
	}
	tests := []struct {
		name   string
		fields fields
		want   []BrowserTotalField
	}{
		{
			name: "meta",
			fields: fields{
				Table: Table{
					Fields: []TableField{
						{Name: "field", FieldType: TableFieldTypeMeta},
					},
					Rows: []ut.IM{
						{"field": 12, "field_meta": TableFieldTypeInteger},
					},
				},
			},
			want: []BrowserTotalField{
				{Name: "field", FieldType: "meta", Total: 12},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bro := &Browser{
				Table:          tt.fields.Table,
				View:           tt.fields.View,
				Views:          tt.fields.Views,
				HideHeader:     tt.fields.HideHeader,
				ShowDropdown:   tt.fields.ShowDropdown,
				ShowColumns:    tt.fields.ShowColumns,
				ShowTotal:      tt.fields.ShowTotal,
				HideBookmark:   tt.fields.HideBookmark,
				HideExport:     tt.fields.HideExport,
				HideHelp:       tt.fields.HideHelp,
				ReadOnly:       tt.fields.ReadOnly,
				VisibleColumns: tt.fields.VisibleColumns,
				Filters:        tt.fields.Filters,
				MetaFields:     tt.fields.MetaFields,
				Labels:         tt.fields.Labels,
				totalFields:    tt.fields.totalFields,
			}
			if got := bro.setTotalValues(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Browser.setTotalValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrowser_msg(t *testing.T) {
	type fields struct {
		Table          Table
		View           string
		Views          []SelectOption
		HideHeader     bool
		ShowDropdown   bool
		ShowColumns    bool
		ShowTotal      bool
		HideBookmark   bool
		HideExport     bool
		HideHelp       bool
		ReadOnly       bool
		VisibleColumns map[string]bool
		Filters        []BrowserFilter
		MetaFields     map[string]BrowserMetaField
		Labels         ut.SM
		totalFields    []BrowserTotalField
	}
	type args struct {
		labelID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "missing",
			args: args{
				labelID: "missing_id",
			},
			want: "missing_id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bro := &Browser{
				Table:          tt.fields.Table,
				View:           tt.fields.View,
				Views:          tt.fields.Views,
				HideHeader:     tt.fields.HideHeader,
				ShowDropdown:   tt.fields.ShowDropdown,
				ShowColumns:    tt.fields.ShowColumns,
				ShowTotal:      tt.fields.ShowTotal,
				HideBookmark:   tt.fields.HideBookmark,
				HideExport:     tt.fields.HideExport,
				HideHelp:       tt.fields.HideHelp,
				ReadOnly:       tt.fields.ReadOnly,
				VisibleColumns: tt.fields.VisibleColumns,
				Filters:        tt.fields.Filters,
				MetaFields:     tt.fields.MetaFields,
				Labels:         tt.fields.Labels,
				totalFields:    tt.fields.totalFields,
			}
			if got := bro.msg(tt.args.labelID); got != tt.want {
				t.Errorf("Browser.msg() = %v, want %v", got, tt.want)
			}
		})
	}
}
