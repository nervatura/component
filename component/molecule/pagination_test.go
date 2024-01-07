package molecule

import (
	"reflect"
	"testing"

	bc "github.com/nervatura/component/component/base"
)

func TestDemoPagination(t *testing.T) {
	for _, tt := range DemoPagination(&bc.BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
}

func TestPagination_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Value         int64
		PageSize      int64
		PageCount     int64
		HidePageSize  bool
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
				propName: "value",
			},
			want: int64(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pgn := &Pagination{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				PageSize:      tt.fields.PageSize,
				PageCount:     tt.fields.PageCount,
				HidePageSize:  tt.fields.HidePageSize,
			}
			if got := pgn.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pagination.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPagination_Validation(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Value         int64
		PageSize      int64
		PageCount     int64
		HidePageSize  bool
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
			pgn := &Pagination{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				PageSize:      tt.fields.PageSize,
				PageCount:     tt.fields.PageCount,
				HidePageSize:  tt.fields.HidePageSize,
			}
			if got := pgn.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pagination.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPagination_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Value         int64
		PageSize      int64
		PageCount     int64
		HidePageSize  bool
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
			pgn := &Pagination{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				PageSize:      tt.fields.PageSize,
				PageCount:     tt.fields.PageCount,
				HidePageSize:  tt.fields.HidePageSize,
			}
			if got := pgn.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pagination.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPagination_response(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Value         int64
		PageSize      int64
		PageCount     int64
		HidePageSize  bool
	}
	type args struct {
		evt bc.ResponseEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "pagination_page_size",
			fields: fields{
				BaseComponent: bc.BaseComponent{
					OnResponse: func(evt bc.ResponseEvent) (re bc.ResponseEvent) {
						evt.Trigger = nil
						return evt
					},
				},
			},
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "pagination_page_size",
					Value:       int64(10),
				},
			},
		},
		{
			name: "value",
			args: args{
				evt: bc.ResponseEvent{
					TriggerName: "pagination_btn_last",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pgn := &Pagination{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				PageSize:      tt.fields.PageSize,
				PageCount:     tt.fields.PageCount,
				HidePageSize:  tt.fields.HidePageSize,
			}
			pgn.response(tt.args.evt)
		})
	}
}
