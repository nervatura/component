package component

import (
	"reflect"
	"testing"

	ut "github.com/nervatura/component/pkg/util"
)

func TestBaseComponent_GetProperty(t *testing.T) {
	type fields struct {
		Id         string
		Name       string
		EventURL   string
		Target     string
		Swap       string
		Indicator  string
		Style      ut.SM
		Data       ut.IM
		RequestMap map[string]ClientComponent
		OnResponse func(evt ResponseEvent) (re ResponseEvent)
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
			name: "ok",
			args: args{
				propName: "name",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bcc := &BaseComponent{
				Id:         tt.fields.Id,
				Name:       tt.fields.Name,
				EventURL:   tt.fields.EventURL,
				Target:     tt.fields.Target,
				Swap:       tt.fields.Swap,
				Indicator:  tt.fields.Indicator,
				Style:      tt.fields.Style,
				Data:       tt.fields.Data,
				RequestMap: tt.fields.RequestMap,
				OnResponse: tt.fields.OnResponse,
			}
			if got := bcc.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseComponent.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseComponent_SetProperty(t *testing.T) {
	type fields struct {
		Id         string
		Name       string
		EventURL   string
		Target     string
		Swap       string
		Indicator  string
		Style      ut.SM
		Data       ut.IM
		RequestMap map[string]ClientComponent
		OnResponse func(evt ResponseEvent) (re ResponseEvent)
	}
	type args struct {
		propName  string
		propValue interface{}
	}
	cc := &BaseComponent{}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "id",
			args: args{
				propName:  "id",
				propValue: "ID",
			},
			want: "ID",
		},
		{
			name: "name",
			args: args{
				propName:  "name",
				propValue: "name",
			},
			want: "name",
		},
		{
			name: "target",
			args: args{
				propName:  "target",
				propValue: "target",
			},
			want: "#target",
		},
		{
			name: "event_url",
			args: args{
				propName:  "event_url",
				propValue: "/post",
			},
			want: "/post",
		},
		{
			name: "swap",
			args: args{
				propName:  "swap",
				propValue: SwapOuterHTML,
			},
			want: SwapOuterHTML,
		},
		{
			name: "indicator",
			args: args{
				propName:  "indicator",
				propValue: "invalid",
			},
			want: IndicatorNone,
		},
		{
			name: "style",
			args: args{
				propName:  "style",
				propValue: ut.SM{"color": "red"},
			},
			want: ut.SM{"color": "red"},
		},
		{
			name: "data",
			args: args{
				propName:  "data",
				propValue: ut.IM{"value": 12345},
			},
			want: ut.IM{"value": 12345},
		},
		{
			name: "request_map",
			fields: fields{
				Id: "id",
			},
			args: args{
				propName:  "request_map",
				propValue: cc,
			},
			want: map[string]ClientComponent{"id": cc},
		},
		{
			name: "request_map",
			fields: fields{
				Id: "",
			},
			args: args{
				propName:  "request_map",
				propValue: cc,
			},
			want: map[string]ClientComponent{},
		},
		{
			name: "class",
			args: args{
				propName:  "class",
				propValue: []string{},
			},
			want: []string{},
		},
		{
			name: "not_found",
			args: args{
				propName:  "not_found",
				propValue: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bcc := &BaseComponent{
				Id:         tt.fields.Id,
				Name:       tt.fields.Name,
				EventURL:   tt.fields.EventURL,
				Target:     tt.fields.Target,
				Swap:       tt.fields.Swap,
				Indicator:  tt.fields.Indicator,
				Style:      tt.fields.Style,
				Data:       tt.fields.Data,
				RequestMap: tt.fields.RequestMap,
				OnResponse: tt.fields.OnResponse,
			}
			if got := bcc.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseComponent.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseComponent_Validation(t *testing.T) {
	type fields struct {
		Id         string
		Name       string
		EventURL   string
		Target     string
		Swap       string
		Indicator  string
		Style      ut.SM
		Data       ut.IM
		RequestMap map[string]ClientComponent
		OnResponse func(evt ResponseEvent) (re ResponseEvent)
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
			name: "not_found",
			args: args{
				propName:  "not_found",
				propValue: "",
			},
			want: "",
		},
		{
			name: "style_imap",
			args: args{
				propName:  "style",
				propValue: ut.IM{"color": "red"},
			},
			want: ut.SM{"color": "red"},
		},
		{
			name: "class_il",
			args: args{
				propName:  "class",
				propValue: []interface{}{"class1", "class2"},
			},
			want: []string{"class1", "class2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bcc := &BaseComponent{
				Id:         tt.fields.Id,
				Name:       tt.fields.Name,
				EventURL:   tt.fields.EventURL,
				Target:     tt.fields.Target,
				Swap:       tt.fields.Swap,
				Indicator:  tt.fields.Indicator,
				Style:      tt.fields.Style,
				Data:       tt.fields.Data,
				RequestMap: tt.fields.RequestMap,
				OnResponse: tt.fields.OnResponse,
			}
			if got := bcc.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseComponent.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseComponent_Render(t *testing.T) {
	type fields struct {
		Id           string
		Name         string
		EventURL     string
		Target       string
		Swap         string
		Indicator    string
		Class        []string
		Style        ut.SM
		Data         ut.IM
		RequestValue map[string]ut.IM
		RequestMap   map[string]ClientComponent
		OnResponse   func(evt ResponseEvent) (re ResponseEvent)
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Id: "id",
				RequestValue: map[string]ut.IM{
					"id": {
						"id": "id",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bcc := &BaseComponent{
				Id:           tt.fields.Id,
				Name:         tt.fields.Name,
				EventURL:     tt.fields.EventURL,
				Target:       tt.fields.Target,
				Swap:         tt.fields.Swap,
				Indicator:    tt.fields.Indicator,
				Class:        tt.fields.Class,
				Style:        tt.fields.Style,
				Data:         tt.fields.Data,
				RequestMap:   tt.fields.RequestMap,
				RequestValue: tt.fields.RequestValue,
				OnResponse:   tt.fields.OnResponse,
			}
			_, err := bcc.Render()
			if (err != nil) != tt.wantErr {
				t.Errorf("BaseComponent.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBaseComponent_OnRequest(t *testing.T) {
	type fields struct {
		Id         string
		Name       string
		EventURL   string
		Target     string
		Swap       string
		Indicator  string
		Class      []string
		Style      ut.SM
		Data       ut.IM
		RequestMap map[string]ClientComponent
		OnResponse func(evt ResponseEvent) (re ResponseEvent)
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
			name: "invalid",
			args: args{
				te: TriggerEvent{},
			},
		},
		{
			name: "valid",
			fields: fields{
				OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
					return re
				},
			},
			args: args{
				te: TriggerEvent{
					Id: "ID12345",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bcc := &BaseComponent{
				Id:         tt.fields.Id,
				Name:       tt.fields.Name,
				EventURL:   tt.fields.EventURL,
				Target:     tt.fields.Target,
				Swap:       tt.fields.Swap,
				Indicator:  tt.fields.Indicator,
				Class:      tt.fields.Class,
				Style:      tt.fields.Style,
				Data:       tt.fields.Data,
				RequestMap: tt.fields.RequestMap,
				OnResponse: tt.fields.OnResponse,
			}
			bcc.OnRequest(tt.args.te)
		})
	}
}
