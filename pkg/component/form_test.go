package component

import (
	"net/url"
	"reflect"
	"testing"

	ut "github.com/nervatura/component/pkg/util"
)

func TestTestForm(t *testing.T) {
	for _, tt := range TestForm(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	testFormResponse(ResponseEvent{Trigger: &Form{}, Name: FormEventOK})
	testFormResponse(ResponseEvent{Trigger: &Form{}, Name: FormEventChange,
		Value: ut.IM{"name": "list", "event": ListEventCurrentPage}})
}

func TestForm_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Title         string
		Icon          string
		BodyRows      []Row
		FooterRows    []Row
		Modal         bool
		SubmitKey     string
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
			frm := &Form{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Icon:          tt.fields.Icon,
				BodyRows:      tt.fields.BodyRows,
				FooterRows:    tt.fields.FooterRows,
				Modal:         tt.fields.Modal,
				SubmitKey:     tt.fields.SubmitKey,
			}
			if got := frm.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Form.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForm_Validation(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Title         string
		Icon          string
		BodyRows      []Row
		FooterRows    []Row
		Modal         bool
		SubmitKey     string
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
			name: "body_rows",
			args: args{
				propName:  "body_rows",
				propValue: ut.IM{},
			},
			want: []Row{},
		},
		{
			name: "footer_rows",
			args: args{
				propName:  "footer_rows",
				propValue: ut.IM{},
			},
			want: []Row{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frm := &Form{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Icon:          tt.fields.Icon,
				BodyRows:      tt.fields.BodyRows,
				FooterRows:    tt.fields.FooterRows,
				Modal:         tt.fields.Modal,
				SubmitKey:     tt.fields.SubmitKey,
			}
			if got := frm.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Form.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForm_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Title         string
		Icon          string
		BodyRows      []Row
		FooterRows    []Row
		Modal         bool
		SubmitKey     string
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
			frm := &Form{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Icon:          tt.fields.Icon,
				BodyRows:      tt.fields.BodyRows,
				FooterRows:    tt.fields.FooterRows,
				Modal:         tt.fields.Modal,
				SubmitKey:     tt.fields.SubmitKey,
			}
			if got := frm.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Form.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForm_OnRequest(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Title         string
		Icon          string
		BodyRows      []Row
		FooterRows    []Row
		Modal         bool
		SubmitKey     string
	}
	type args struct {
		te TriggerEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "InputBoxEventOK",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
				},
				SubmitKey: "btn_ok",
			},
			args: args{
				te: TriggerEvent{
					Values: url.Values{
						"btn_ok": {},
						"string": {"test"},
					},
					Name: FormEventOK,
				},
			},
		},
		{
			name: "on_response",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
					OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
						evt.Trigger = &Form{}
						return evt
					},
				},
				SubmitKey: "btn_ok",
			},
			args: args{
				te: TriggerEvent{
					Values: url.Values{},
					Name:   FormEventCancel,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frm := &Form{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Icon:          tt.fields.Icon,
				BodyRows:      tt.fields.BodyRows,
				FooterRows:    tt.fields.FooterRows,
				Modal:         tt.fields.Modal,
				SubmitKey:     tt.fields.SubmitKey,
			}
			frm.OnRequest(tt.args.te)
		})
	}
}

func TestForm_triggerEvent(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Title         string
		Icon          string
		BodyRows      []Row
		FooterRows    []Row
		Modal         bool
		SubmitKey     string
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
			name: "InputBoxEventOK",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
				},
				SubmitKey: "btn_ok",
			},
			args: args{
				evt: ResponseEvent{
					Trigger:     &Toggle{},
					TriggerName: "bool",
					Name:        ToggleEventChange,
					Value:       ut.IM{"string": "test"},
				},
			},
		},
		{
			name: "on_response",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
					OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
						evt.Trigger = &Form{}
						return evt
					},
				},
				SubmitKey: "btn_ok",
			},
			args: args{
				evt: ResponseEvent{
					Trigger:     &Button{},
					TriggerName: "btn_close",
					Name:        ButtonEventClick,
					Value:       ut.IM{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frm := &Form{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Icon:          tt.fields.Icon,
				BodyRows:      tt.fields.BodyRows,
				FooterRows:    tt.fields.FooterRows,
				Modal:         tt.fields.Modal,
				SubmitKey:     tt.fields.SubmitKey,
			}
			frm.triggerEvent(tt.args.evt)
		})
	}
}

func TestForm_getComponent(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Title         string
		Icon          string
		BodyRows      []Row
		FooterRows    []Row
		Modal         bool
		SubmitKey     string
	}
	type args struct {
		name  string
		index int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "btn_close",
			fields: fields{
				BaseComponent: BaseComponent{
					Data: ut.IM{},
				},
			},
			args: args{
				name:  "btn_close",
				index: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frm := &Form{
				BaseComponent: tt.fields.BaseComponent,
				Title:         tt.fields.Title,
				Icon:          tt.fields.Icon,
				BodyRows:      tt.fields.BodyRows,
				FooterRows:    tt.fields.FooterRows,
				Modal:         tt.fields.Modal,
				SubmitKey:     tt.fields.SubmitKey,
			}
			_, err := frm.getComponent(tt.args.name, tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("Form.getComponent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
