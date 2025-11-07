package component

import (
	"bytes"
	"errors"
	"io"
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

func TestToIMA(t *testing.T) {
	type args struct {
		value    interface{}
		defValue []IM
	}
	tests := []struct {
		name string
		args args
		want []IM
	}{
		{
			name: "valid",
			args: args{
				value:    []IM{},
				defValue: []IM{},
			},
			want: []IM{},
		},
		{
			name: "sm",
			args: args{
				value:    []SM{{"field": "value"}},
				defValue: []IM{},
			},
			want: []IM{{"field": "value"}},
		},
		{
			name: "if",
			args: args{
				value:    []interface{}{SM{"field1": "value"}, IM{"field2": "value"}},
				defValue: []IM{},
			},
			want: []IM{{"field1": "value"}, {"field2": "value"}},
		},
		{
			name: "default",
			args: args{
				value:    nil,
				defValue: []IM{},
			},
			want: []IM{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToIMA(tt.args.value, tt.args.defValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToIMA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertFromReader(t *testing.T) {
	type args struct {
		data   io.Reader
		result interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				data:   bytes.NewBufferString(`{"key": "value"}`),
				result: &SM{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConvertFromReader(tt.args.data, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("ConvertFromReader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConvertToByte(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				data: SM{"key": "value"},
			},
			want: []byte(`{"key":"value"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToByte(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToByte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertFromByte(t *testing.T) {
	type args struct {
		data   []byte
		result interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				data:   []byte(`{"key": "value"}`),
				result: &SM{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConvertFromByte(tt.args.data, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("ConvertFromByte() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestToIM(t *testing.T) {
	type args struct {
		im       interface{}
		defValue IM
	}
	tests := []struct {
		name       string
		args       args
		wantResult IM
	}{
		{
			name: "ok",
			args: args{
				im:       IM{"key": "value"},
				defValue: IM{},
			},
			wantResult: IM{"key": "value"},
		},
		{
			name: "nil",
			args: args{
				im:       nil,
				defValue: IM{},
			},
			wantResult: IM{},
		},
		{
			name: "empty",
			args: args{
				im:       []IM{},
				defValue: IM{},
			},
			wantResult: IM{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := ToIM(tt.args.im, tt.args.defValue); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ToIM() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestToSM(t *testing.T) {
	type args struct {
		sm       interface{}
		defValue SM
	}
	tests := []struct {
		name       string
		args       args
		wantResult SM
	}{
		{
			name: "ok",
			args: args{
				sm:       SM{"key": "value"},
				defValue: SM{},
			},
			wantResult: SM{"key": "value"},
		},
		{
			name: "nil",
			args: args{
				sm:       nil,
				defValue: SM{},
			},
			wantResult: SM{},
		},
		{
			name: "empty",
			args: args{
				sm:       []SM{},
				defValue: SM{},
			},
			wantResult: SM{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := ToSM(tt.args.sm, tt.args.defValue); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ToSM() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestIMToSM(t *testing.T) {
	type args struct {
		im IM
	}
	tests := []struct {
		name string
		args args
		want SM
	}{
		{
			name: "string",
			args: args{
				im: IM{"key": "value"},
			},
			want: SM{"key": "value"},
		},
		{
			name: "object",
			args: args{
				im: IM{"key": IM{"key": "value"}},
			},
			want: SM{"key": `{"key":"value"}`},
		},
		{
			name: "nil",
			args: args{
				im: IM{"key": nil},
			},
			want: SM{"key": "null"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IMToSM(tt.args.im); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IMToSM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestILtoSL(t *testing.T) {
	type args struct {
		ivalue interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "string",
			args: args{
				ivalue: []string{"value1", "value2"},
			},
			want: []string{"value1", "value2"},
		},
		{
			name: "interface",
			args: args{
				ivalue: []interface{}{"value1", "value2"},
			},
			want: []string{"value1", "value2"},
		},
		{
			name: "int",
			args: args{
				ivalue: []int{1, 2},
			},
			want: []string{"1", "2"},
		},
		{
			name: "int64",
			args: args{
				ivalue: []int64{1, 2},
			},
			want: []string{"1", "2"},
		},
		{
			name: "float64",
			args: args{
				ivalue: []float64{1.1, 2.2},
			},
			want: []string{"1.1", "2.2"},
		},
		{
			name: "bool",
			args: args{
				ivalue: []bool{true, false},
			},
			want: []string{"true", "false"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ILtoSL(tt.args.ivalue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ILtoSL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringLimit(t *testing.T) {
	type args struct {
		value  string
		length int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				value:  "test",
				length: 5,
			},
			want: "test",
		},
		{
			name: "limit",
			args: args{
				value:  "testtest",
				length: 5,
			},
			want: "testt...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringLimit(tt.args.value, tt.args.length); got != tt.want {
				t.Errorf("StringLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToBoolMap(t *testing.T) {
	type args struct {
		im       interface{}
		defValue map[string]bool
	}
	tests := []struct {
		name       string
		args       args
		wantResult map[string]bool
	}{
		{
			name: "ok",
			args: args{
				im:       map[string]bool{"key": true},
				defValue: map[string]bool{},
			},
			wantResult: map[string]bool{"key": true},
		},
		{
			name: "nil",
			args: args{
				im:       nil,
				defValue: map[string]bool{},
			},
			wantResult: map[string]bool{},
		},
		{
			name: "empty",
			args: args{
				im:       []map[string]bool{},
				defValue: map[string]bool{},
			},
			wantResult: map[string]bool{},
		},
		{
			name: "interface",
			args: args{
				im:       IM{"key": true},
				defValue: map[string]bool{},
			},
			wantResult: map[string]bool{"key": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := ToBoolMap(tt.args.im, tt.args.defValue); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ToBoolMap() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestConvertToType(t *testing.T) {
	type args struct {
		data   any
		result any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				data:   SM{"key": "value"},
				result: &SM{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConvertToType(tt.args.data, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("ConvertToType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
