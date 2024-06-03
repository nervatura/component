package component

import (
	"reflect"
	"testing"
)

func TestTestIcon(t *testing.T) {
	for _, tt := range TestIcon(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	testIcoResponse(ResponseEvent{Trigger: &Icon{}})
}

func TestIcon_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Value         string
		Width         float64
		Height        float64
		Color         string
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
				propName: "id",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ico := &Icon{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Width:         tt.fields.Width,
				Height:        tt.fields.Height,
				Color:         tt.fields.Color,
			}
			if got := ico.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Icon.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIcon_Validation(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Value         string
		Width         float64
		Height        float64
		Color         string
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
			ico := &Icon{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Width:         tt.fields.Width,
				Height:        tt.fields.Height,
				Color:         tt.fields.Color,
			}
			if got := ico.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Icon.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIcon_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Value         string
		Width         float64
		Height        float64
		Color         string
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
			ico := &Icon{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Width:         tt.fields.Width,
				Height:        tt.fields.Height,
				Color:         tt.fields.Color,
			}
			if got := ico.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Icon.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIcon_OnRequest(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Value         string
		Width         float64
		Height        float64
		Color         string
	}
	type args struct {
		te TriggerEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantRe ResponseEvent
	}{
		{
			name: "base",
			args: args{
				te: TriggerEvent{
					Id: "id",
				},
			},
			wantRe: ResponseEvent{
				Trigger: &Icon{},
				Name:    IconEventClick,
				Value:   "",
			},
		},
		{
			name: "OnResponse",
			fields: fields{
				BaseComponent: BaseComponent{
					OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
						evt.Trigger = &Icon{}
						return evt
					},
				},
			},
			args: args{
				te: TriggerEvent{
					Id: "id",
				},
			},
			wantRe: ResponseEvent{
				Trigger: &Icon{},
				Name:    IconEventClick,
				Value:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ico := &Icon{
				BaseComponent: tt.fields.BaseComponent,
				Value:         tt.fields.Value,
				Width:         tt.fields.Width,
				Height:        tt.fields.Height,
				Color:         tt.fields.Color,
			}
			if gotRe := ico.OnRequest(tt.args.te); !reflect.DeepEqual(gotRe, tt.wantRe) {
				t.Errorf("Icon.OnRequest() = %v, want %v", gotRe, tt.wantRe)
			}
		})
	}
}
