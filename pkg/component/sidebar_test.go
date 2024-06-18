package component

import (
	"reflect"
	"testing"

	ut "github.com/nervatura/component/pkg/util"
)

func TestTestSidebar(t *testing.T) {
	for _, tt := range TestSidebar(&BaseComponent{EventURL: "/demo"}) {
		t.Run(tt.Label, func(t *testing.T) {
			tt.Component.Render()
		})
	}
	testSidebarResponse(ResponseEvent{Name: SideBarEventItem, Trigger: &Selector{}})
	testSidebarResponse(ResponseEvent{Name: SideBarEventGroup, Trigger: &Selector{}})
}

func TestSideBar_GetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Items         []SideBarItem
		Visibility    string
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
				propName: "visibility",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &SideBar{
				BaseComponent: tt.fields.BaseComponent,
				Items:         tt.fields.Items,
				Visibility:    tt.fields.Visibility,
			}
			if got := sb.GetProperty(tt.args.propName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SideBar.GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSideBar_Validation(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Items         []SideBarItem
		Visibility    string
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
			name: "items_invalid",
			args: args{
				propName:  "items",
				propValue: "",
			},
			want: []SideBarItem{},
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
			sb := &SideBar{
				BaseComponent: tt.fields.BaseComponent,
				Items:         tt.fields.Items,
				Visibility:    tt.fields.Visibility,
			}
			if got := sb.Validation(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SideBar.Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSideBar_SetProperty(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Items         []SideBarItem
		Visibility    string
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
			sb := &SideBar{
				BaseComponent: tt.fields.BaseComponent,
				Items:         tt.fields.Items,
				Visibility:    tt.fields.Visibility,
			}
			if got := sb.SetProperty(tt.args.propName, tt.args.propValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SideBar.SetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSideBar_response(t *testing.T) {
	type fields struct {
		BaseComponent BaseComponent
		Items         []SideBarItem
		Visibility    string
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
			name: "element",
			fields: fields{
				BaseComponent: BaseComponent{
					OnResponse: func(evt ResponseEvent) (re ResponseEvent) {
						evt.Trigger = &SideBar{}
						return evt
					},
				},
				Items: []SideBarItem{
					&SideBarElement{},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: SideBarItemTypeElement,
					Trigger: &Button{
						BaseComponent: BaseComponent{
							Data: ut.IM{"index": 0},
						},
					},
				},
			},
		},
		{
			name: "group_head",
			fields: fields{
				Items: []SideBarItem{
					&SideBarGroup{
						Items: []SideBarElement{
							{},
						},
					},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: SideBarItemTypeGroup,
					Trigger: &Button{
						BaseComponent: BaseComponent{
							Data: ut.IM{"index": 0},
						},
					},
				},
			},
		},
		{
			name: "group_item",
			fields: fields{
				Items: []SideBarItem{
					&SideBarGroup{
						Items: []SideBarElement{
							{},
						},
					},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: SideBarItemTypeGroup,
					Trigger: &Button{
						BaseComponent: BaseComponent{
							Data: ut.IM{"index": 0, "group_index": 0},
						},
					},
				},
			},
		},
		{
			name: "state",
			fields: fields{
				Items: []SideBarItem{
					&SideBarState{
						Items: []SideBarElement{
							{},
						},
					},
				},
			},
			args: args{
				evt: ResponseEvent{
					TriggerName: SideBarItemTypeGroup,
					Trigger: &Button{
						BaseComponent: BaseComponent{
							Data: ut.IM{"index": 0, "group_index": 0},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &SideBar{
				BaseComponent: tt.fields.BaseComponent,
				Items:         tt.fields.Items,
				Visibility:    tt.fields.Visibility,
			}
			sb.response(tt.args.evt)
		})
	}
}

func TestSideBarState_GetValue(t *testing.T) {
	type fields struct {
		Name          string
		SelectedIndex int
		Items         []SideBarElement
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "value",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sbst := &SideBarState{
				Name:          tt.fields.Name,
				SelectedIndex: tt.fields.SelectedIndex,
				Items:         tt.fields.Items,
			}
			if got := sbst.GetValue(); got != tt.want {
				t.Errorf("SideBarState.GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSideBarState_GetSelected(t *testing.T) {
	type fields struct {
		Name          string
		SelectedIndex int
		Items         []SideBarElement
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "selected",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sbst := &SideBarState{
				Name:          tt.fields.Name,
				SelectedIndex: tt.fields.SelectedIndex,
				Items:         tt.fields.Items,
			}
			if got := sbst.GetSelected(); got != tt.want {
				t.Errorf("SideBarState.GetSelected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSideBarElement_GetSelected(t *testing.T) {
	type fields struct {
		Name     string
		Value    string
		Label    string
		Align    string
		Icon     string
		Selected bool
		Disabled bool
		NotFull  bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "selected",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sbe := &SideBarElement{
				Name:     tt.fields.Name,
				Value:    tt.fields.Value,
				Label:    tt.fields.Label,
				Align:    tt.fields.Align,
				Icon:     tt.fields.Icon,
				Selected: tt.fields.Selected,
				Disabled: tt.fields.Disabled,
				NotFull:  tt.fields.NotFull,
			}
			if got := sbe.GetSelected(); got != tt.want {
				t.Errorf("SideBarElement.GetSelected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSideBarStatic_GetValue(t *testing.T) {
	type fields struct {
		Label string
		Icon  string
		Color string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "value",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sbs := &SideBarStatic{
				Label: tt.fields.Label,
				Icon:  tt.fields.Icon,
				Color: tt.fields.Color,
			}
			if got := sbs.GetValue(); got != tt.want {
				t.Errorf("SideBarStatic.GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSideBarStatic_GetSelected(t *testing.T) {
	type fields struct {
		Label string
		Icon  string
		Color string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "selected",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sbs := &SideBarStatic{
				Label: tt.fields.Label,
				Icon:  tt.fields.Icon,
				Color: tt.fields.Color,
			}
			if got := sbs.GetSelected(); got != tt.want {
				t.Errorf("SideBarStatic.GetSelected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSideBarSeparator_GetValue(t *testing.T) {
	tests := []struct {
		name string
		sbsp *SideBarSeparator
		want string
	}{
		{
			name: "value",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sbsp := &SideBarSeparator{}
			if got := sbsp.GetValue(); got != tt.want {
				t.Errorf("SideBarSeparator.GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSideBarSeparator_GetSelected(t *testing.T) {
	tests := []struct {
		name string
		sbsp *SideBarSeparator
		want bool
	}{
		{
			name: "selected",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sbsp := &SideBarSeparator{}
			if got := sbsp.GetSelected(); got != tt.want {
				t.Errorf("SideBarSeparator.GetSelected() = %v, want %v", got, tt.want)
			}
		})
	}
}
