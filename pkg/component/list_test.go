package component

import (
	"reflect"
	"testing"

	ut "github.com/nervatura/component/pkg/util"
)

func TestTestList(t *testing.T) {
	for _, tt := range TestList(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	demoListResponse(ResponseEvent{Name: ListEventAddItem, Trigger: &List{}})
	demoListResponse(ResponseEvent{Name: ListEventFilterChange, Trigger: &List{}})
}

func TestList_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Rows              []ut.IM
		Pagination        string
		CurrentPage       int64
		PageSize          int64
		HidePaginatonSize bool
		ListFilter        bool
		AddItem           bool
		EditItem          bool
		DeleteItem        bool
		FilterPlaceholder string
		FilterValue       string
		LabelAdd          string
		AddIcon           string
		EditIcon          string
		DeleteIcon        string
		LabelField        string
		LabelValue        string
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
			lst := &List{
				BaseComponent:     tt.fields.BaseComponent,
				Rows:              tt.fields.Rows,
				Pagination:        tt.fields.Pagination,
				CurrentPage:       tt.fields.CurrentPage,
				PageSize:          tt.fields.PageSize,
				HidePaginatonSize: tt.fields.HidePaginatonSize,
				ListFilter:        tt.fields.ListFilter,
				AddItem:           tt.fields.AddItem,
				EditItem:          tt.fields.EditItem,
				DeleteItem:        tt.fields.DeleteItem,
				FilterPlaceholder: tt.fields.FilterPlaceholder,
				FilterValue:       tt.fields.FilterValue,
				LabelAdd:          tt.fields.LabelAdd,
				AddIcon:           tt.fields.AddIcon,
				EditIcon:          tt.fields.EditIcon,
				DeleteIcon:        tt.fields.DeleteIcon,
				LabelField:        tt.fields.LabelField,
				LabelValue:        tt.fields.LabelValue,
			}
			if got := lst.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_Validation(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Rows              []ut.IM
		Pagination        string
		CurrentPage       int64
		PageSize          int64
		HidePaginatonSize bool
		ListFilter        bool
		AddItem           bool
		EditItem          bool
		DeleteItem        bool
		FilterPlaceholder string
		FilterValue       string
		LabelAdd          string
		AddIcon           string
		EditIcon          string
		DeleteIcon        string
		LabelField        string
		LabelValue        string
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
			lst := &List{
				BaseComponent:     tt.fields.BaseComponent,
				Rows:              tt.fields.Rows,
				Pagination:        tt.fields.Pagination,
				CurrentPage:       tt.fields.CurrentPage,
				PageSize:          tt.fields.PageSize,
				HidePaginatonSize: tt.fields.HidePaginatonSize,
				ListFilter:        tt.fields.ListFilter,
				AddItem:           tt.fields.AddItem,
				EditItem:          tt.fields.EditItem,
				DeleteItem:        tt.fields.DeleteItem,
				FilterPlaceholder: tt.fields.FilterPlaceholder,
				FilterValue:       tt.fields.FilterValue,
				LabelAdd:          tt.fields.LabelAdd,
				AddIcon:           tt.fields.AddIcon,
				EditIcon:          tt.fields.EditIcon,
				DeleteIcon:        tt.fields.DeleteIcon,
				LabelField:        tt.fields.LabelField,
				LabelValue:        tt.fields.LabelValue,
			}
			if got := lst.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Rows              []ut.IM
		Pagination        string
		CurrentPage       int64
		PageSize          int64
		HidePaginatonSize bool
		ListFilter        bool
		AddItem           bool
		EditItem          bool
		DeleteItem        bool
		FilterPlaceholder string
		FilterValue       string
		LabelAdd          string
		AddIcon           string
		EditIcon          string
		DeleteIcon        string
		LabelField        string
		LabelValue        string
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
			lst := &List{
				BaseComponent:     tt.fields.BaseComponent,
				Rows:              tt.fields.Rows,
				Pagination:        tt.fields.Pagination,
				CurrentPage:       tt.fields.CurrentPage,
				PageSize:          tt.fields.PageSize,
				HidePaginatonSize: tt.fields.HidePaginatonSize,
				ListFilter:        tt.fields.ListFilter,
				AddItem:           tt.fields.AddItem,
				EditItem:          tt.fields.EditItem,
				DeleteItem:        tt.fields.DeleteItem,
				FilterPlaceholder: tt.fields.FilterPlaceholder,
				FilterValue:       tt.fields.FilterValue,
				LabelAdd:          tt.fields.LabelAdd,
				AddIcon:           tt.fields.AddIcon,
				EditIcon:          tt.fields.EditIcon,
				DeleteIcon:        tt.fields.DeleteIcon,
				LabelField:        tt.fields.LabelField,
				LabelValue:        tt.fields.LabelValue,
			}
			if got := lst.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_response(t *testing.T) {
	type fields struct {
		BaseComponent     BaseComponent
		Rows              []ut.IM
		Pagination        string
		CurrentPage       int64
		PageSize          int64
		HidePaginatonSize bool
		ListFilter        bool
		AddItem           bool
		EditItem          bool
		DeleteItem        bool
		FilterPlaceholder string
		FilterValue       string
		LabelAdd          string
		AddIcon           string
		EditIcon          string
		DeleteIcon        string
		LabelField        string
		LabelValue        string
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
						evt.Trigger = &List{}
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
			name: "edit_item",
			args: args{
				evt: ResponseEvent{
					Trigger: &List{
						BaseComponent: BaseComponent{
							Data: ut.IM{},
						},
					},
					TriggerName: "edit_item",
				},
			},
		},
		{
			name: "delete_item",
			args: args{
				evt: ResponseEvent{
					Trigger: &List{
						BaseComponent: BaseComponent{
							Data: ut.IM{},
						},
					},
					TriggerName: "delete_item",
				},
			},
		},
		{
			name: "invalid",
			args: args{
				evt: ResponseEvent{
					Trigger: &List{
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
			lst := &List{
				BaseComponent:     tt.fields.BaseComponent,
				Rows:              tt.fields.Rows,
				Pagination:        tt.fields.Pagination,
				CurrentPage:       tt.fields.CurrentPage,
				PageSize:          tt.fields.PageSize,
				HidePaginatonSize: tt.fields.HidePaginatonSize,
				ListFilter:        tt.fields.ListFilter,
				AddItem:           tt.fields.AddItem,
				EditItem:          tt.fields.EditItem,
				DeleteItem:        tt.fields.DeleteItem,
				FilterPlaceholder: tt.fields.FilterPlaceholder,
				FilterValue:       tt.fields.FilterValue,
				LabelAdd:          tt.fields.LabelAdd,
				AddIcon:           tt.fields.AddIcon,
				EditIcon:          tt.fields.EditIcon,
				DeleteIcon:        tt.fields.DeleteIcon,
				LabelField:        tt.fields.LabelField,
				LabelValue:        tt.fields.LabelValue,
			}
			lst.response(tt.args.evt)
		})
	}
}
