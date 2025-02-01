package component

import (
	"net/url"
	"reflect"
	"testing"

	ut "github.com/nervatura/component/pkg/util"
)

func TestInputBox_Render(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		InputType     string
		Value         string
		ValueOptions  []SelectOption
		Title         string
		Message       string
		Info          string
		Tag           string
		LabelOK       string
		LabelCancel   string
		DefaultOK     bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				BaseComponent: BaseComponent{
					Id:       "test_inputbox_cancel",
					EventURL: "/demo",
				},
				InputType: InputBoxTypeCancel,
			},
			wantErr: false,
		},
		{
			name: "input",
			fields: fields{
				InputType: InputBoxTypeString,
				Value:     "test",
			},
			wantErr: false,
		},
		{
			name: "text",
			fields: fields{
				InputType: InputBoxTypeText,
				Value:     "test",
			},
			wantErr: false,
		},
		{
			name: "color",
			fields: fields{
				InputType: InputBoxTypeColor,
				Value:     "#000000",
			},
			wantErr: false,
		},
		{
			name: "options",
			fields: fields{
				InputType: InputBoxTypeSelect,
				ValueOptions: []SelectOption{
					{Text: "Option 1", Value: "option1"},
					{Text: "Option 2", Value: "option2"},
					{Text: "Option 3", Value: "option3"},
				},
			},
			wantErr: false,
		},
		{
			name: "datetime",
			fields: fields{
				InputType: InputBoxTypeDateTime,
				Value:     "2025-01-01T12:00",
			},
			wantErr: false,
		},
		{
			name: "integer",
			fields: fields{
				InputType: InputBoxTypeInteger,
				Value:     "123",
			},
			wantErr: false,
		},
		{
			name: "number",
			fields: fields{
				InputType: InputBoxTypeNumber,
				Value:     "123.45",
			},
			wantErr: false,
		},
		{
			name: "date",
			fields: fields{
				InputType: InputBoxTypeDate,
				Value:     "2025-01-01",
			},
			wantErr: false,
		},
		{
			name: "time",
			fields: fields{
				InputType: InputBoxTypeTime,
				Value:     "12:00",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ibx := &InputBox{
				BaseComponent: tt.fields.BaseComponent,
				InputType:     tt.fields.InputType,
				Value:         tt.fields.Value,
				ValueOptions:  tt.fields.ValueOptions,
				Title:         tt.fields.Title,
				Message:       tt.fields.Message,
				Info:          tt.fields.Info,
				Tag:           tt.fields.Tag,
				LabelOK:       tt.fields.LabelOK,
				LabelCancel:   tt.fields.LabelCancel,
				DefaultOK:     tt.fields.DefaultOK,
			}
			_, err := ibx.Render()
			if (err != nil) != tt.wantErr {
				t.Errorf("InputBox.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestInputBox_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		InputType     string
		Value         string
		ValueOptions  []SelectOption
		Title         string
		Message       string
		Info          string
		Tag           string
		LabelOK       string
		LabelCancel   string
		DefaultOK     bool
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
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ibx := &InputBox{
				BaseComponent: tt.fields.BaseComponent,
				InputType:     tt.fields.InputType,
				Value:         tt.fields.Value,
				ValueOptions:  tt.fields.ValueOptions,
				Title:         tt.fields.Title,
				Message:       tt.fields.Message,
				Info:          tt.fields.Info,
				Tag:           tt.fields.Tag,
				LabelOK:       tt.fields.LabelOK,
				LabelCancel:   tt.fields.LabelCancel,
				DefaultOK:     tt.fields.DefaultOK,
			}
			if got := ibx.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InputBox.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInputBox_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		InputType     string
		Value         string
		ValueOptions  []SelectOption
		Title         string
		Message       string
		Info          string
		Tag           string
		LabelOK       string
		LabelCancel   string
		DefaultOK     bool
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
			ibx := &InputBox{
				BaseComponent: tt.fields.BaseComponent,
				InputType:     tt.fields.InputType,
				Value:         tt.fields.Value,
				ValueOptions:  tt.fields.ValueOptions,
				Title:         tt.fields.Title,
				Message:       tt.fields.Message,
				Info:          tt.fields.Info,
				Tag:           tt.fields.Tag,
				LabelOK:       tt.fields.LabelOK,
				LabelCancel:   tt.fields.LabelCancel,
				DefaultOK:     tt.fields.DefaultOK,
			}
			if got := ibx.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InputBox.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInputBox_Validation(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		InputType     string
		Value         string
		ValueOptions  []SelectOption
		Title         string
		Message       string
		Info          string
		Tag           string
		LabelOK       string
		LabelCancel   string
		DefaultOK     bool
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
			name: "invalid",
			args: args{
				propName:  "invalid",
				propValue: "",
			},
			want: "",
		},
		{
			name: "value_options",
			args: args{
				propName:  "value_options",
				propValue: []interface{}{ut.IM{"value": "value1", "text": "text1"}},
			},
			want: []SelectOption{{Value: "value1", Text: "text1"}},
		},
		{
			name: "input_type_string",
			args: args{propName: "input_type", propValue: "IBOX_STRING"},
			want: InputBoxTypeString,
		},
		{
			name: "input_type_int",
			args: args{propName: "input_type", propValue: 2},
			want: InputBoxTypeCancel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ibx := &InputBox{
				BaseComponent: tt.fields.BaseComponent,
				InputType:     tt.fields.InputType,
				Value:         tt.fields.Value,
				ValueOptions:  tt.fields.ValueOptions,
				Title:         tt.fields.Title,
				Message:       tt.fields.Message,
				Info:          tt.fields.Info,
				Tag:           tt.fields.Tag,
				LabelOK:       tt.fields.LabelOK,
				LabelCancel:   tt.fields.LabelCancel,
				DefaultOK:     tt.fields.DefaultOK,
			}
			if got := ibx.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InputBox.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTestInputBox(t *testing.T) {

	for _, tt := range TestInputBox(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}

	testInputBoxResponse(ResponseEvent{Name: ButtonEventClick, Trigger: &Button{
		BaseComponent: BaseComponent{
			Data: ut.IM{
				"value_options": []SelectOption{
					{Text: "Option 1", Value: "option1"},
					{Text: "Option 2", Value: "option2"},
					{Text: "Option 3", Value: "option3"},
				},
				"input_type": InputBoxTypeSelect,
			},
		},
	}})
	testInputBoxResponse(ResponseEvent{Name: ButtonEventClick, Trigger: &Button{
		BaseComponent: BaseComponent{
			Data: ut.IM{"input_type": InputBoxTypeCancel},
		},
	}})
	testInputBoxResponse(ResponseEvent{Name: InputBoxEventOK, Trigger: &InputBox{
		BaseComponent: BaseComponent{
			Data: ut.IM{"input_type": InputBoxTypeOK},
		},
	}})

	testInputBoxResponse(ResponseEvent{Name: ButtonEventClick, Trigger: &Button{
		BaseComponent: BaseComponent{
			Data: ut.IM{"input_type": InputBoxTypeString},
		},
	}})

}

func TestInputBox_OnRequest(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		InputType     string
		Value         string
		ValueOptions  []SelectOption
		Title         string
		Message       string
		Info          string
		Tag           string
		LabelOK       string
		LabelCancel   string
		DefaultOK     bool
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
			name: "base",
			args: args{
				te: TriggerEvent{
					Id: "id",
					Values: url.Values{
						"btn_ok": []string{"true"},
						"value":  []string{"test"},
					},
				},
			},
		},
		{
			name: "OnResponse",
			fields: fields{
				BaseComponent: BaseComponent{
					OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
						evt.Trigger = &InputBox{}
						return evt
					},
				},
			},
			args: args{
				te: TriggerEvent{
					Id: "id",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ibx := &InputBox{
				BaseComponent: tt.fields.BaseComponent,
				InputType:     tt.fields.InputType,
				Value:         tt.fields.Value,
				ValueOptions:  tt.fields.ValueOptions,
				Title:         tt.fields.Title,
				Message:       tt.fields.Message,
				Info:          tt.fields.Info,
				Tag:           tt.fields.Tag,
				LabelOK:       tt.fields.LabelOK,
				LabelCancel:   tt.fields.LabelCancel,
				DefaultOK:     tt.fields.DefaultOK,
			}
			ibx.OnRequest(tt.args.te)
		})
	}
}
