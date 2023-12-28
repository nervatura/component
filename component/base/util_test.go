package base

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestToString(t *testing.T) {
	type args struct {
		value    interface{}
		defValue string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "string",
			args: args{
				value:    "test",
				defValue: "",
			},
			want: "test",
		},
		{
			name: "empty",
			args: args{
				value:    "",
				defValue: "",
			},
			want: "",
		},
		{
			name: "bool",
			args: args{
				value:    true,
				defValue: "",
			},
			want: "true",
		},
		{
			name: "int",
			args: args{
				value:    int(1),
				defValue: "",
			},
			want: "1",
		},
		{
			name: "int32",
			args: args{
				value:    int32(1),
				defValue: "",
			},
			want: "1",
		},
		{
			name: "int64",
			args: args{
				value:    int64(1),
				defValue: "",
			},
			want: "1",
		},
		{
			name: "float32",
			args: args{
				value:    float32(1),
				defValue: "",
			},
			want: "1",
		},
		{
			name: "float64",
			args: args{
				value:    float64(1.1),
				defValue: "",
			},
			want: "1.1",
		},
		{
			name: "time",
			args: args{
				value:    time.Now(),
				defValue: "",
			},
			want: time.Now().Format("2006-01-02T15:04:05-07:00"),
		},
		{
			name: "default",
			args: args{
				value:    []string{},
				defValue: "default",
			},
			want: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToString(tt.args.value, tt.args.defValue); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat(t *testing.T) {
	type args struct {
		value    interface{}
		defValue float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "float64",
			args: args{
				value:    float64(1.1),
				defValue: float64(0),
			},
			want: float64(1.1),
		},
		{
			name: "0",
			args: args{
				value:    float64(0),
				defValue: float64(0),
			},
			want: float64(0),
		},
		{
			name: "bool",
			args: args{
				value:    true,
				defValue: float64(0),
			},
			want: float64(1),
		},
		{
			name: "int",
			args: args{
				value:    int(1),
				defValue: float64(0),
			},
			want: float64(1),
		},
		{
			name: "int32",
			args: args{
				value:    int32(1),
				defValue: float64(0),
			},
			want: float64(1),
		},
		{
			name: "int64",
			args: args{
				value:    int64(1),
				defValue: float64(0),
			},
			want: float64(1),
		},
		{
			name: "float32",
			args: args{
				value:    float32(1),
				defValue: float64(0),
			},
			want: float64(1),
		},
		{
			name: "string",
			args: args{
				value:    "1",
				defValue: float64(0),
			},
			want: float64(1),
		},
		{
			name: "default",
			args: args{
				value:    []string{},
				defValue: float64(0),
			},
			want: float64(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToFloat(tt.args.value, tt.args.defValue); got != tt.want {
				t.Errorf("ToFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInteger(t *testing.T) {
	type args struct {
		value    interface{}
		defValue int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "int64",
			args: args{
				value:    int64(1),
				defValue: int64(0),
			},
			want: int64(1),
		},
		{
			name: "0",
			args: args{
				value:    int64(0),
				defValue: int64(0),
			},
			want: int64(0),
		},
		{
			name: "bool",
			args: args{
				value:    true,
				defValue: int64(0),
			},
			want: int64(1),
		},
		{
			name: "int",
			args: args{
				value:    int(1),
				defValue: int64(0),
			},
			want: int64(1),
		},
		{
			name: "int32",
			args: args{
				value:    int32(1),
				defValue: int64(0),
			},
			want: int64(1),
		},
		{
			name: "float32",
			args: args{
				value:    float32(1),
				defValue: int64(0),
			},
			want: int64(1),
		},
		{
			name: "float64",
			args: args{
				value:    float64(1),
				defValue: int64(0),
			},
			want: int64(1),
		},
		{
			name: "string",
			args: args{
				value:    "1",
				defValue: int64(0),
			},
			want: int64(1),
		},
		{
			name: "default",
			args: args{
				value:    []string{},
				defValue: int64(0),
			},
			want: int64(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInteger(tt.args.value, tt.args.defValue); got != tt.want {
				t.Errorf("ToInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToBoolean(t *testing.T) {
	type args struct {
		value    interface{}
		defValue bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "bool",
			args: args{
				value:    true,
				defValue: false,
			},
			want: true,
		},
		{
			name: "int",
			args: args{
				value:    int(1),
				defValue: false,
			},
			want: true,
		},
		{
			name: "int32",
			args: args{
				value:    int32(1),
				defValue: false,
			},
			want: true,
		},
		{
			name: "int64",
			args: args{
				value:    int64(1),
				defValue: false,
			},
			want: true,
		},
		{
			name: "float32",
			args: args{
				value:    float32(1),
				defValue: false,
			},
			want: true,
		},
		{
			name: "float64",
			args: args{
				value:    float64(1),
				defValue: false,
			},
			want: true,
		},
		{
			name: "string",
			args: args{
				value:    "true",
				defValue: false,
			},
			want: true,
		},
		{
			name: "default",
			args: args{
				value:    []string{},
				defValue: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBoolean(tt.args.value, tt.args.defValue); got != tt.want {
				t.Errorf("ToBoolean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToDateTime(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "time1",
			args: args{
				value: "2006-01-02T15:04:05-07:00",
			},
			wantErr: false,
		},
		{
			name: "time2",
			args: args{
				value: "2006-01-02T15:04:05-0700",
			},
			wantErr: false,
		},
		{
			name: "time3",
			args: args{
				value: "2006-01-02T15:04:05",
			},
			wantErr: false,
		},
		{
			name: "time4",
			args: args{
				value: "2006-01-02T15:04:05Z",
			},
			wantErr: false,
		},
		{
			name: "time5",
			args: args{
				value: "2006-01-02 15:04:05",
			},
			wantErr: false,
		},
		{
			name: "time6",
			args: args{
				value: "2006-01-02 15:04",
			},
			wantErr: false,
		},
		{
			name: "time7",
			args: args{
				value: "2006-01-02",
			},
			wantErr: false,
		},
		{
			name: "time8",
			args: args{
				value: "2006-01",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := StringToDateTime(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToDateTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRandString(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "create",
			args: args{
				length: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandString(tt.args.length); len(got) != tt.args.length {
				t.Errorf("RandString() = %v, want %v", len(got), tt.args.length)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		a []string
		x string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "found",
			args: args{
				a: []string{"abba", "baba", "aabb"},
				x: "baba",
			},
			want: true,
		},
		{
			name: "not_found",
			args: args{
				a: []string{"abba", "baba", "aabb"},
				x: "bbaa",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.a, tt.args.x); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetComponentID(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "get",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetComponentID()
		})
	}
}

func TestTemplateBuilder(t *testing.T) {
	type args struct {
		name    string
		tpl     string
		funcMap map[string]any
		data    any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "parse_error",
			args: args{
				tpl: `{{.???`,
			},
			wantErr: true,
		},
		{
			name: "execute_error",
			args: args{
				tpl: `<div>{{ data }}</div>`,
				funcMap: map[string]any{
					"data": func() (string, error) {
						return "", errors.New("error")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				tpl: `<div>{{ data }}</div>`,
				funcMap: map[string]any{
					"data": func() (string, error) {
						return "data", nil
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := TemplateBuilder(tt.args.name, tt.args.tpl, tt.args.funcMap, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TemplateBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSetIMValue(t *testing.T) {
	type args struct {
		imap  IM
		key   string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want IM
	}{
		{
			name: "ok",
			args: args{
				key:   "key",
				value: "value",
			},
			want: IM{"key": "value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetIMValue(tt.args.imap, tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetIMValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetSMValue(t *testing.T) {
	type args struct {
		smap  SM
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
		want SM
	}{
		{
			name: "ok",
			args: args{
				key:   "key",
				value: "value",
			},
			want: SM{"key": "value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetSMValue(tt.args.smap, tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetSMValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetCMValue(t *testing.T) {
	type args struct {
		cmap  map[string]ClientComponent
		key   string
		value ClientComponent
	}
	tests := []struct {
		name string
		args args
		want map[string]ClientComponent
	}{
		{
			name: "ok",
			args: args{
				key:   "key",
				value: nil,
			},
			want: map[string]ClientComponent{"key": nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetCMValue(tt.args.cmap, tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetCMValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeSM(t *testing.T) {
	type args struct {
		baseMap  SM
		valueMap SM
	}
	tests := []struct {
		name string
		args args
		want SM
	}{
		{
			name: "ok",
			args: args{
				baseMap:  SM{"key1": "value1"},
				valueMap: SM{"key2": "value2"},
			},
			want: SM{"key1": "value1", "key2": "value2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeSM(tt.args.baseMap, tt.args.valueMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeSM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeIM(t *testing.T) {
	type args struct {
		baseMap  IM
		valueMap IM
	}
	tests := []struct {
		name string
		args args
		want IM
	}{
		{
			name: "ok",
			args: args{
				baseMap:  IM{"key1": "value1"},
				valueMap: IM{"key2": "value2"},
			},
			want: IM{"key1": "value1", "key2": "value2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeIM(tt.args.baseMap, tt.args.valueMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeIM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeCM(t *testing.T) {
	type args struct {
		baseMap  map[string]ClientComponent
		valueMap map[string]ClientComponent
	}
	tests := []struct {
		name string
		args args
		want map[string]ClientComponent
	}{
		{
			name: "ok",
			args: args{
				baseMap:  map[string]ClientComponent{"key1": nil},
				valueMap: map[string]ClientComponent{"key2": nil},
			},
			want: map[string]ClientComponent{"key1": nil, "key2": nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeCM(tt.args.baseMap, tt.args.valueMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeCM() = %v, want %v", got, tt.want)
			}
		})
	}
}
