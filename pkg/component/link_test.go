package component

import (
	"reflect"
	"testing"
)

func TestTestLink(t *testing.T) {
	for _, tt := range TestLink(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
}

func TestLink_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent  BaseComponent
		Href           string
		Download       string
		Media          string
		Ping           string
		ReferrerPolicy string
		Rel            string
		LinkTarget     string
		MediaType      string
		LinkStyle      string
		Align          string
		Label          string
		LabelComponent ClientComponent
		Icon           string
		Disabled       bool
		AutoFocus      bool
		Full           bool
		Small          bool
		Selected       bool
		HideLabel      bool
		Badge          int64
		ShowBadge      bool
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
				propName: "disabled",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lnk := &Link{
				BaseComponent:  tt.fields.BaseComponent,
				Href:           tt.fields.Href,
				Download:       tt.fields.Download,
				Media:          tt.fields.Media,
				Ping:           tt.fields.Ping,
				ReferrerPolicy: tt.fields.ReferrerPolicy,
				Rel:            tt.fields.Rel,
				LinkTarget:     tt.fields.LinkTarget,
				MediaType:      tt.fields.MediaType,
				LinkStyle:      tt.fields.LinkStyle,
				Align:          tt.fields.Align,
				Label:          tt.fields.Label,
				LabelComponent: tt.fields.LabelComponent,
				Icon:           tt.fields.Icon,
				Disabled:       tt.fields.Disabled,
				AutoFocus:      tt.fields.AutoFocus,
				Full:           tt.fields.Full,
				Small:          tt.fields.Small,
				Selected:       tt.fields.Selected,
				HideLabel:      tt.fields.HideLabel,
				Badge:          tt.fields.Badge,
				ShowBadge:      tt.fields.ShowBadge,
			}
			if got := lnk.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Link.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLink_Validation(t *testing.T) {
	type fields struct {
		BaseComponent  BaseComponent
		Href           string
		Download       string
		Media          string
		Ping           string
		ReferrerPolicy string
		Rel            string
		LinkTarget     string
		MediaType      string
		LinkStyle      string
		Align          string
		Label          string
		LabelComponent ClientComponent
		Icon           string
		Disabled       bool
		AutoFocus      bool
		Full           bool
		Small          bool
		Selected       bool
		HideLabel      bool
		Badge          int64
		ShowBadge      bool
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
			lnk := &Link{
				BaseComponent:  tt.fields.BaseComponent,
				Href:           tt.fields.Href,
				Download:       tt.fields.Download,
				Media:          tt.fields.Media,
				Ping:           tt.fields.Ping,
				ReferrerPolicy: tt.fields.ReferrerPolicy,
				Rel:            tt.fields.Rel,
				LinkTarget:     tt.fields.LinkTarget,
				MediaType:      tt.fields.MediaType,
				LinkStyle:      tt.fields.LinkStyle,
				Align:          tt.fields.Align,
				Label:          tt.fields.Label,
				LabelComponent: tt.fields.LabelComponent,
				Icon:           tt.fields.Icon,
				Disabled:       tt.fields.Disabled,
				AutoFocus:      tt.fields.AutoFocus,
				Full:           tt.fields.Full,
				Small:          tt.fields.Small,
				Selected:       tt.fields.Selected,
				HideLabel:      tt.fields.HideLabel,
				Badge:          tt.fields.Badge,
				ShowBadge:      tt.fields.ShowBadge,
			}
			if got := lnk.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Link.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLink_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent  BaseComponent
		Href           string
		Download       string
		Media          string
		Ping           string
		ReferrerPolicy string
		Rel            string
		LinkTarget     string
		MediaType      string
		LinkStyle      string
		Align          string
		Label          string
		LabelComponent ClientComponent
		Icon           string
		Disabled       bool
		AutoFocus      bool
		Full           bool
		Small          bool
		Selected       bool
		HideLabel      bool
		Badge          int64
		ShowBadge      bool
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
			lnk := &Link{
				BaseComponent:  tt.fields.BaseComponent,
				Href:           tt.fields.Href,
				Download:       tt.fields.Download,
				Media:          tt.fields.Media,
				Ping:           tt.fields.Ping,
				ReferrerPolicy: tt.fields.ReferrerPolicy,
				Rel:            tt.fields.Rel,
				LinkTarget:     tt.fields.LinkTarget,
				MediaType:      tt.fields.MediaType,
				LinkStyle:      tt.fields.LinkStyle,
				Align:          tt.fields.Align,
				Label:          tt.fields.Label,
				LabelComponent: tt.fields.LabelComponent,
				Icon:           tt.fields.Icon,
				Disabled:       tt.fields.Disabled,
				AutoFocus:      tt.fields.AutoFocus,
				Full:           tt.fields.Full,
				Small:          tt.fields.Small,
				Selected:       tt.fields.Selected,
				HideLabel:      tt.fields.HideLabel,
				Badge:          tt.fields.Badge,
				ShowBadge:      tt.fields.ShowBadge,
			}
			if got := lnk.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Link.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}
