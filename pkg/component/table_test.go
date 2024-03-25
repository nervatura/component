package component

import (
	"reflect"
	"testing"

	ut "github.com/nervatura/component/pkg/util"
)

func TestTestTable(t *testing.T) {
	for _, tt := range TestTable(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	demoTableResponse(ResponseEvent{Name: TableEventAddItem, Trigger: &Table{}})
	demoTableResponse(ResponseEvent{Name: TableEventFilterChange, Trigger: &Table{}})
}

func TestTable_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		RowKey            string
		Rows              []ut.IM
		Fields            []TableField
		Pagination        string
		CurrentPage       int64
		PageSize          int64
		HidePaginatonSize bool
		TableFilter       bool
		AddItem           bool
		FilterPlaceholder string
		FilterValue       string
		LabelYes          string
		LabelNo           string
		LabelAdd          string
		AddIcon           string
		TablePadding      string
		SortCol           string
		SortAsc           bool
		RowSelected       bool
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
			tbl := &Table{
				BaseComponent:     tt.fields.BaseComponent,
				RowKey:            tt.fields.RowKey,
				Rows:              tt.fields.Rows,
				Fields:            tt.fields.Fields,
				Pagination:        tt.fields.Pagination,
				CurrentPage:       tt.fields.CurrentPage,
				PageSize:          tt.fields.PageSize,
				HidePaginatonSize: tt.fields.HidePaginatonSize,
				TableFilter:       tt.fields.TableFilter,
				AddItem:           tt.fields.AddItem,
				FilterPlaceholder: tt.fields.FilterPlaceholder,
				FilterValue:       tt.fields.FilterValue,
				LabelYes:          tt.fields.LabelYes,
				LabelNo:           tt.fields.LabelNo,
				LabelAdd:          tt.fields.LabelAdd,
				AddIcon:           tt.fields.AddIcon,
				TablePadding:      tt.fields.TablePadding,
				SortCol:           tt.fields.SortCol,
				SortAsc:           tt.fields.SortAsc,
				RowSelected:       tt.fields.RowSelected,
			}
			if got := tbl.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Table.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_Validation(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		RowKey            string
		Rows              []ut.IM
		Fields            []TableField
		Pagination        string
		CurrentPage       int64
		PageSize          int64
		HidePaginatonSize bool
		TableFilter       bool
		AddItem           bool
		FilterPlaceholder string
		FilterValue       string
		LabelYes          string
		LabelNo           string
		LabelAdd          string
		AddIcon           string
		TablePadding      string
		SortCol           string
		SortAsc           bool
		RowSelected       bool
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
			name: "current_page",
			args: args{
				propName:  "current_page",
				propValue: "",
			},
			want: int64(1),
		},
		{
			name: "page_size",
			args: args{
				propName:  "page_size",
				propValue: 3,
			},
			want: int64(5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := &Table{
				BaseComponent:     tt.fields.BaseComponent,
				RowKey:            tt.fields.RowKey,
				Rows:              tt.fields.Rows,
				Fields:            tt.fields.Fields,
				Pagination:        tt.fields.Pagination,
				CurrentPage:       tt.fields.CurrentPage,
				PageSize:          tt.fields.PageSize,
				HidePaginatonSize: tt.fields.HidePaginatonSize,
				TableFilter:       tt.fields.TableFilter,
				AddItem:           tt.fields.AddItem,
				FilterPlaceholder: tt.fields.FilterPlaceholder,
				FilterValue:       tt.fields.FilterValue,
				LabelYes:          tt.fields.LabelYes,
				LabelNo:           tt.fields.LabelNo,
				LabelAdd:          tt.fields.LabelAdd,
				AddIcon:           tt.fields.AddIcon,
				TablePadding:      tt.fields.TablePadding,
				SortCol:           tt.fields.SortCol,
				SortAsc:           tt.fields.SortAsc,
				RowSelected:       tt.fields.RowSelected,
			}
			if got := tbl.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Table.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		RowKey            string
		Rows              []ut.IM
		Fields            []TableField
		Pagination        string
		CurrentPage       int64
		PageSize          int64
		HidePaginatonSize bool
		TableFilter       bool
		AddItem           bool
		FilterPlaceholder string
		FilterValue       string
		LabelYes          string
		LabelNo           string
		LabelAdd          string
		AddIcon           string
		TablePadding      string
		SortCol           string
		SortAsc           bool
		RowSelected       bool
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
			tbl := &Table{
				BaseComponent:     tt.fields.BaseComponent,
				RowKey:            tt.fields.RowKey,
				Rows:              tt.fields.Rows,
				Fields:            tt.fields.Fields,
				Pagination:        tt.fields.Pagination,
				CurrentPage:       tt.fields.CurrentPage,
				PageSize:          tt.fields.PageSize,
				HidePaginatonSize: tt.fields.HidePaginatonSize,
				TableFilter:       tt.fields.TableFilter,
				AddItem:           tt.fields.AddItem,
				FilterPlaceholder: tt.fields.FilterPlaceholder,
				FilterValue:       tt.fields.FilterValue,
				LabelYes:          tt.fields.LabelYes,
				LabelNo:           tt.fields.LabelNo,
				LabelAdd:          tt.fields.LabelAdd,
				AddIcon:           tt.fields.AddIcon,
				TablePadding:      tt.fields.TablePadding,
				SortCol:           tt.fields.SortCol,
				SortAsc:           tt.fields.SortAsc,
				RowSelected:       tt.fields.RowSelected,
			}
			if got := tbl.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Table.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_SortRows(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		RowKey            string
		Rows              []ut.IM
		Fields            []TableField
		Pagination        string
		CurrentPage       int64
		PageSize          int64
		HidePaginatonSize bool
		TableFilter       bool
		AddItem           bool
		FilterPlaceholder string
		FilterValue       string
		LabelYes          string
		LabelNo           string
		LabelAdd          string
		AddIcon           string
		TablePadding      string
		SortCol           string
		SortAsc           bool
		RowSelected       bool
	}
	type args struct {
		fieldName string
		fieldType string
		sortAsc   bool
	}
	rows := []ut.IM{
		{"string": "a", "number": 2},
		{"string": "b", "number": 1},
		{"string": "a", "number": 1},
		{"string": "d", "number": 2},
		{"string": "b", "number": 3},
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "string_asc",
			fields: fields{
				Rows: rows,
			},
			args: args{
				fieldName: "string",
				fieldType: TableFieldTypeString,
				sortAsc:   true,
			},
		},
		{
			name: "string",
			fields: fields{
				Rows: rows,
			},
			args: args{
				fieldName: "string",
				fieldType: TableFieldTypeString,
				sortAsc:   false,
			},
		},
		{
			name: "number",
			fields: fields{
				Rows: rows,
			},
			args: args{
				fieldName: "number",
				fieldType: TableFieldTypeNumber,
				sortAsc:   false,
			},
		},
		{
			name: "number_asc",
			fields: fields{
				Rows: rows,
			},
			args: args{
				fieldName: "number",
				fieldType: TableFieldTypeNumber,
				sortAsc:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := &Table{
				BaseComponent:     tt.fields.BaseComponent,
				RowKey:            tt.fields.RowKey,
				Rows:              tt.fields.Rows,
				Fields:            tt.fields.Fields,
				Pagination:        tt.fields.Pagination,
				CurrentPage:       tt.fields.CurrentPage,
				PageSize:          tt.fields.PageSize,
				HidePaginatonSize: tt.fields.HidePaginatonSize,
				TableFilter:       tt.fields.TableFilter,
				AddItem:           tt.fields.AddItem,
				FilterPlaceholder: tt.fields.FilterPlaceholder,
				FilterValue:       tt.fields.FilterValue,
				LabelYes:          tt.fields.LabelYes,
				LabelNo:           tt.fields.LabelNo,
				LabelAdd:          tt.fields.LabelAdd,
				AddIcon:           tt.fields.AddIcon,
				TablePadding:      tt.fields.TablePadding,
				SortCol:           tt.fields.SortCol,
				SortAsc:           tt.fields.SortAsc,
				RowSelected:       tt.fields.RowSelected,
			}
			tbl.SortRows(tt.args.fieldName, tt.args.fieldType, tt.args.sortAsc)
		})
	}
}

func TestTable_response(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		RowKey            string
		Rows              []ut.IM
		Fields            []TableField
		Pagination        string
		CurrentPage       int64
		PageSize          int64
		HidePaginatonSize bool
		TableFilter       bool
		AddItem           bool
		FilterPlaceholder string
		FilterValue       string
		LabelYes          string
		LabelNo           string
		LabelAdd          string
		AddIcon           string
		TablePadding      string
		SortCol           string
		SortAsc           bool
		RowSelected       bool
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
			name: "top_pagination",
			args: args{
				evt: ResponseEvent{
					TriggerName: "top_pagination",
					Name:        PaginationEventPageSize,
				},
			},
		},
		{
			name: "bottom_pagination",
			args: args{
				evt: ResponseEvent{
					TriggerName: "bottom_pagination",
					Name:        PaginationEventPageSize,
				},
			},
		},
		{
			name: "current_page",
			args: args{
				evt: ResponseEvent{
					TriggerName: "top_pagination",
				},
			},
		},
		{
			name: "filter",
			fields: fields{
				BaseComponent: BaseComponent{
					OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
						evt.Trigger = &Table{}
						return evt
					},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "filter",
				},
			},
		},
		{
			name: "btn_add",
			args: args{
				evt: ResponseEvent{
					TriggerName: "btn_add",
				},
			},
		},
		{
			name: "link_cell",
			args: args{
				evt: ResponseEvent{
					Trigger: &Table{
						BaseComponent: BaseComponent{
							Data: ut.IM{},
						},
					},
					TriggerName: "link_cell",
				},
			},
		},
		{
			name: "data_row",
			args: args{
				evt: ResponseEvent{
					Trigger: &Table{
						BaseComponent: BaseComponent{
							Data: ut.IM{},
						},
					},
					TriggerName: "data_row",
				},
			},
		},
		{
			name: "header_sort",
			fields: fields{
				Rows: []ut.IM{
					{"string": "a", "number": 2},
					{"string": "b", "number": 1},
					{"string": "a", "number": 1},
					{"string": "d", "number": 2},
					{"string": "b", "number": 3},
				},
				SortCol: "string",
			},
			args: args{
				evt: ResponseEvent{
					Trigger: &Table{
						BaseComponent: BaseComponent{
							Data: ut.IM{"fieldname": "string"},
						},
					},
					TriggerName: "header_sort",
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
			tbl := &Table{
				BaseComponent:     tt.fields.BaseComponent,
				RowKey:            tt.fields.RowKey,
				Rows:              tt.fields.Rows,
				Fields:            tt.fields.Fields,
				Pagination:        tt.fields.Pagination,
				CurrentPage:       tt.fields.CurrentPage,
				PageSize:          tt.fields.PageSize,
				HidePaginatonSize: tt.fields.HidePaginatonSize,
				TableFilter:       tt.fields.TableFilter,
				AddItem:           tt.fields.AddItem,
				FilterPlaceholder: tt.fields.FilterPlaceholder,
				FilterValue:       tt.fields.FilterValue,
				LabelYes:          tt.fields.LabelYes,
				LabelNo:           tt.fields.LabelNo,
				LabelAdd:          tt.fields.LabelAdd,
				AddIcon:           tt.fields.AddIcon,
				TablePadding:      tt.fields.TablePadding,
				SortCol:           tt.fields.SortCol,
				SortAsc:           tt.fields.SortAsc,
				RowSelected:       tt.fields.RowSelected,
			}
			tbl.response(tt.args.evt)
		})
	}
}
