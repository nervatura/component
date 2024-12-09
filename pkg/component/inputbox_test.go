package component

import (
	"reflect"
	"testing"

	ut "github.com/nervatura/component/pkg/util"
)

func TestInputBox_Render(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Value         string
		ValueOptions  []SelectOption
		Title         string
		Message       string
		Info          string
		Tag           string
		LabelOK       string
		LabelCancel   string
		ShowValue     bool
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
				ShowValue: true,
			},
			wantErr: false,
		},
		{
			name: "options",
			fields: fields{
				ShowValue: true,
				ValueOptions: []SelectOption{
					{Text: "Option 1", Value: "option1"},
					{Text: "Option 2", Value: "option2"},
					{Text: "Option 3", Value: "option3"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ibx := &InputBox{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				ValueOptions:  tt.fields.ValueOptions,
				Title:         tt.fields.Title,
				Message:       tt.fields.Message,
				Info:          tt.fields.Info,
				Tag:           tt.fields.Tag,
				LabelOK:       tt.fields.LabelOK,
				LabelCancel:   tt.fields.LabelCancel,
				ShowValue:     tt.fields.ShowValue,
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
		Value         string
		Title         string
		Message       string
		Info          string
		Tag           string
		LabelOK       string
		LabelCancel   string
		ShowValue     bool
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
				Value:         tt.fields.Value,
				Title:         tt.fields.Title,
				Message:       tt.fields.Message,
				Info:          tt.fields.Info,
				Tag:           tt.fields.Tag,
				LabelOK:       tt.fields.LabelOK,
				LabelCancel:   tt.fields.LabelCancel,
				ShowValue:     tt.fields.ShowValue,
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
		Value         string
		Title         string
		Message       string
		Info          string
		Tag           string
		LabelOK       string
		LabelCancel   string
		ShowValue     bool
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
				Value:         tt.fields.Value,
				Title:         tt.fields.Title,
				Message:       tt.fields.Message,
				Info:          tt.fields.Info,
				Tag:           tt.fields.Tag,
				LabelOK:       tt.fields.LabelOK,
				LabelCancel:   tt.fields.LabelCancel,
				ShowValue:     tt.fields.ShowValue,
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
		Value         string
		Title         string
		Message       string
		Info          string
		Tag           string
		LabelOK       string
		LabelCancel   string
		ShowValue     bool
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ibx := &InputBox{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Title:         tt.fields.Title,
				Message:       tt.fields.Message,
				Info:          tt.fields.Info,
				Tag:           tt.fields.Tag,
				LabelOK:       tt.fields.LabelOK,
				LabelCancel:   tt.fields.LabelCancel,
				ShowValue:     tt.fields.ShowValue,
				DefaultOK:     tt.fields.DefaultOK,
			}
			if got := ibx.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InputBox.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInputBox_response(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Value         string
		Title         string
		Message       string
		Info          string
		Tag           string
		LabelOK       string
		LabelCancel   string
		ShowValue     bool
		DefaultOK     bool
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
			name: "btn_ok",
			fields: fields{
				BaseComponent: BaseComponent{
					OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
						evt.Trigger = &Row{}
						return evt
					},
				},
				ShowValue: true,
				Tag:       "tag",
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "btn_ok",
				},
			},
		},
		{
			name: "input_value",
			fields: fields{
				ShowValue: true,
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "input_value",
				},
			},
		},
		{
			name: "btn_cancel",
			fields: fields{
				ShowValue: true,
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: "btn_cancel",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ibx := &InputBox{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Title:         tt.fields.Title,
				Message:       tt.fields.Message,
				Info:          tt.fields.Info,
				Tag:           tt.fields.Tag,
				LabelOK:       tt.fields.LabelOK,
				LabelCancel:   tt.fields.LabelCancel,
				ShowValue:     tt.fields.ShowValue,
				DefaultOK:     tt.fields.DefaultOK,
			}
			ibx.response(tt.args.evt)
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
			Data: ut.IM{"value_options": []SelectOption{
				{Text: "Option 1", Value: "option1"},
				{Text: "Option 2", Value: "option2"},
				{Text: "Option 3", Value: "option3"},
			}},
		},
	}})
	testInputBoxResponse(ResponseEvent{Name: ButtonEventClick, Trigger: &Button{}})
	testInputBoxResponse(ResponseEvent{Name: InputBoxEventOK, Trigger: &InputBox{}})
	testInputBoxResponse(ResponseEvent{Name: InputBoxEventValueChange, Trigger: &InputBox{}})
}
