/* Component helper functions
 */
package component

import (
	"crypto/rand"
	"html/template"
	"math/big"
	"strconv"
	"strings"
	"time"
)

// IM is a map[string]interface{} type short alias
type IM = map[string]interface{}

// SM is a map[string]string type short alias
type SM = map[string]string

// ToString - safe string conversion
func ToString(value interface{}, defValue string) string {
	if stringValue, valid := value.(string); valid {
		if stringValue == "" {
			return defValue
		}
		return stringValue
	}
	if boolValue, valid := value.(bool); valid {
		return strconv.FormatBool(boolValue)
	}
	if intValue, valid := value.(int); valid {
		return strconv.Itoa(intValue)
	}
	if intValue, valid := value.(int32); valid {
		return strconv.Itoa(int(intValue))
	}
	if intValue, valid := value.(int64); valid {
		return strconv.FormatInt(intValue, 10)
	}
	if floatValue, valid := value.(float32); valid {
		return strconv.FormatFloat(float64(floatValue), 'f', -1, 64)
	}
	if floatValue, valid := value.(float64); valid {
		return strconv.FormatFloat(floatValue, 'f', -1, 64)
	}
	if timeValue, valid := value.(time.Time); valid {
		return timeValue.Format("2006-01-02T15:04:05-07:00")
	}
	return defValue
}

// ToFloat - safe float64 conversion
func ToFloat(value interface{}, defValue float64) float64 {
	if floatValue, valid := value.(float64); valid {
		if floatValue == 0 {
			return defValue
		}
		return floatValue
	}
	if boolValue, valid := value.(bool); valid {
		if boolValue {
			return 1
		}
	}
	if intValue, valid := value.(int); valid {
		return float64(intValue)
	}
	if intValue, valid := value.(int32); valid {
		return float64(intValue)
	}
	if intValue, valid := value.(int64); valid {
		return float64(intValue)
	}
	if floatValue, valid := value.(float32); valid {
		return float64(floatValue)
	}
	if stringValue, valid := value.(string); valid {
		floatValue, err := strconv.ParseFloat(stringValue, 64)
		if err == nil {
			return float64(floatValue)
		}
	}
	return defValue
}

// ToInteger - safe int64 conversion
func ToInteger(value interface{}, defValue int64) int64 {
	if intValue, valid := value.(int64); valid {
		if intValue == 0 {
			return defValue
		}
		return intValue
	}
	if boolValue, valid := value.(bool); valid {
		if boolValue {
			return 1
		}
	}
	if intValue, valid := value.(int); valid {
		return int64(intValue)
	}
	if intValue, valid := value.(int32); valid {
		return int64(intValue)
	}
	if floatValue, valid := value.(float32); valid {
		return int64(floatValue)
	}
	if floatValue, valid := value.(float64); valid {
		return int64(floatValue)
	}
	if stringValue, valid := value.(string); valid {
		intValue, err := strconv.ParseInt(stringValue, 10, 64)
		if err == nil {
			return int64(intValue)
		}
	}
	return defValue
}

// ToBoolean - safe bool conversion
func ToBoolean(value interface{}, defValue bool) bool {
	if boolValue, valid := value.(bool); valid {
		return boolValue
	}
	if intValue, valid := value.(int); valid {
		if intValue == 1 {
			return true
		}
	}
	if intValue, valid := value.(int32); valid {
		if intValue == 1 {
			return true
		}
	}
	if intValue, valid := value.(int64); valid {
		if intValue == 1 {
			return true
		}
	}
	if floatValue, valid := value.(float32); valid {
		if floatValue == 1 {
			return true
		}
	}
	if floatValue, valid := value.(float64); valid {
		if floatValue == 1 {
			return true
		}
	}
	if stringValue, valid := value.(string); valid {
		boolValue, err := strconv.ParseBool(stringValue)
		if err == nil {
			return boolValue
		}
	}
	return defValue
}

func ToIMA(value interface{}, defValue []IM) []IM {
	if imaValue, valid := value.([]IM); valid {
		return imaValue
	}
	var iRows []IM = []IM{}
	if smaValue, valid := value.([]SM); valid {
		for _, sm := range smaValue {
			for key, svalue := range sm {
				iRows = append(iRows, IM{key: svalue})
			}
		}
		return iRows
	}
	if ifaValue, valid := value.([]interface{}); valid {
		for _, ifRow := range ifaValue {
			if im, valid := ifRow.(IM); valid {
				iRows = append(iRows, im)
			}
			if sm, valid := ifRow.(SM); valid {
				for key, svalue := range sm {
					iRows = append(iRows, IM{key: svalue})
				}
			}
		}
		return iRows
	}
	return defValue
}

// StringToDateTime - parse string to datetime
func StringToDateTime(value string) (time.Time, error) {
	tm, err := time.Parse("2006-01-02T15:04:05-07:00", value)
	if err != nil {
		tm, err = time.Parse("2006-01-02T15:04:05-0700", value)
	}
	if err != nil {
		tm, err = time.Parse("2006-01-02T15:04:05", value)
	}
	if err != nil {
		tm, err = time.Parse("2006-01-02T15:04:05Z", value)
	}
	if err != nil {
		tm, err = time.Parse("2006-01-02 15:04:05", value)
	}
	if err != nil {
		tm, err = time.Parse("2006-01-02 15:04", value)
	}
	if err != nil {
		tm, err = time.Parse("2006-01-02", value)
	}
	return tm, err
}

func RandString(length int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	var b strings.Builder
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars[:51]))))
	b.WriteRune(chars[n.Int64()])
	for i := 1; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b.WriteRune(chars[n.Int64()])
	}
	return b.String()
}

func GetComponentID() string {
	return "ID" + RandString(16)
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func TemplateBuilder(name, tpl string, funcMap map[string]any, data any) (html template.HTML, err error) {
	var tmp *template.Template
	if tmp, err = template.New(name).Funcs(funcMap).Parse(tpl); err != nil {
		return "", err
	}

	var sb strings.Builder
	if err = tmp.Execute(&sb, data); err != nil {
		return "", err
	}

	return template.HTML(strings.ReplaceAll(sb.String(), "\n\t", "")), err
}

// SetIMValue - safe IM value setting
func SetIMValue(imap IM, key string, value interface{}) IM {
	if imap == nil {
		imap = IM{}
	}
	if key != "" {
		imap[key] = value
	}
	return imap
}

// SetSMValue - safe SM value setting
func SetSMValue(smap SM, key string, value string) SM {
	if smap == nil {
		smap = SM{}
	}
	if key != "" {
		smap[key] = value
	}
	return smap
}

func MergeSM(baseMap, valueMap SM) SM {
	for key, svalue := range valueMap {
		baseMap = SetSMValue(baseMap, key, svalue)
	}
	return baseMap
}

func MergeIM(baseMap, valueMap IM) IM {
	for key, ivalue := range valueMap {
		baseMap = SetIMValue(baseMap, key, ivalue)
	}
	return baseMap
}
