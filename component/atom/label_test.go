package atom

import (
	"reflect"
	"testing"

	bc "github.com/nervatura/component/component/base"
)

func TestDemoLabel(t *testing.T) {
	for _, tt := range DemoLabel(&bc.BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	demoLblResponse(bc.ResponseEvent{Trigger: &Label{}})
}

func TestLabel_Validation(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Value         string
		Centered      bool
		LeftIcon      string
		RightIcon     string
		IconStyle     bc.SM
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
			lbl := &Label{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Centered:      tt.fields.Centered,
				LeftIcon:      tt.fields.LeftIcon,
				RightIcon:     tt.fields.RightIcon,
				IconStyle:     tt.fields.IconStyle,
			}
			if got := lbl.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Label.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLabel_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Value         string
		Centered      bool
		LeftIcon      string
		RightIcon     string
		IconStyle     bc.SM
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
			lbl := &Label{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Centered:      tt.fields.Centered,
				LeftIcon:      tt.fields.LeftIcon,
				RightIcon:     tt.fields.RightIcon,
				IconStyle:     tt.fields.IconStyle,
			}
			if got := lbl.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Label.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLabel_OnRequest(t *testing.T) {
	type fields struct {
		BaseComponent bc.BaseComponent
		Value         string
		Centered      bool
		LeftIcon      string
		RightIcon     string
		IconStyle     bc.SM
	}
	type args struct {
		te bc.TriggerEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantRe bc.ResponseEvent
	}{
		{
			name: "base",
			args: args{
				te: bc.TriggerEvent{
					Id: "id",
				},
			},
			wantRe: bc.ResponseEvent{
				Trigger: &Label{},
				Name:    LabelEventClick,
				Value:   "",
			},
		},
		{
			name: "OnResponse",
			fields: fields{
				BaseComponent: bc.BaseComponent{
					OnResponse: func(evt bc.ResponseEvent) (re bc.ResponseEvent) {
						evt.Trigger = &Label{}
						return evt
					},
				},
			},
			args: args{
				te: bc.TriggerEvent{
					Id: "id",
				},
			},
			wantRe: bc.ResponseEvent{
				Trigger: &Label{},
				Name:    LabelEventClick,
				Value:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lbl := &Label{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Centered:      tt.fields.Centered,
				LeftIcon:      tt.fields.LeftIcon,
				RightIcon:     tt.fields.RightIcon,
				IconStyle:     tt.fields.IconStyle,
			}
			if gotRe := lbl.OnRequest(tt.args.te); !reflect.DeepEqual(gotRe, tt.wantRe) {
				t.Errorf("Label.OnRequest() = %v, want %v", gotRe, tt.wantRe)
			}
		})
	}
}
