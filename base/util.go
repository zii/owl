package base

import (
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// if production enviroment
var RELEASE = os.Getenv("RELEASE") == "1"

var (
	digestRegexp = regexp.MustCompile("^[0-9]+$")
)

func ValidPhoneNumber(s string) bool {
	if s == "" {
		return false
	}
	if s[0] == '+' {
		return true
	}
	return digestRegexp.MatchString(s)
}

func AddIntSet(array []int, d int) []int {
	for _, v := range array {
		if v == d {
			return array
		}
	}
	return append(array, d)
}

func AddInt64Set(array []int64, d int64) []int64 {
	for _, v := range array {
		if v == d {
			return array
		}
	}
	return append(array, d)
}

func IntsRemove(array []int, d int) []int {
	var n = len(array)
	if n == 0 {
		return array
	}
	var p = -1
	for i, v := range array {
		if v == d {
			p = i
			break
		}
	}
	if p < 0 {
		return array
	}
	var prev, next []int
	if p == 0 {
		next = array[1:]
	} else if p >= n-1 {
		prev = array[:n-1]
	} else {
		prev = array[:p]
		next = array[p+1:]
	}
	return append(prev, next...)
}

func Int64sRemove(array []int64, d int64) []int64 {
	var n = len(array)
	if n == 0 {
		return array
	}
	var out = make([]int64, 0, len(array))
	for _, a := range array {
		if a != d {
			out = append(out, a)
		}
	}
	return out
}

func InStrings(src string, array []string) bool {
	for _, s := range array {
		if src == s {
			return true
		}
	}
	return false
}

func InInt64s(src int64, array []int64) bool {
	for _, v := range array {
		if src == v {
			return true
		}
	}
	return false
}

func CmpInt64s(a, b []int64) int {
	for i, av := range a {
		if i >= len(b) {
			return 1
		}
		bv := b[i]
		if av > bv {
			return 1
		} else if av < bv {
			return -1
		}
	}
	if len(b) > len(a) {
		return -1
	}
	return 0
}

func InInt32s(src int32, array []int32) bool {
	for _, v := range array {
		if src == v {
			return true
		}
	}
	return false
}

func InInts(src int, array []int) bool {
	for _, v := range array {
		if src == v {
			return true
		}
	}
	return false
}

func Int32sToInts(src []int32) []int {
	var out []int
	for _, v := range src {
		out = append(out, int(v))
	}
	return out
}

func IntsToInt32s(src []int) []int32 {
	var out []int32
	for _, v := range src {
		out = append(out, int32(v))
	}
	return out
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func ToJson(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func JoinStrings(sl []string) string {
	return strings.Join(sl, ",")
}

func JoinInts(il []int) string {
	var s string
	for _, v := range il {
		if len(s) > 0 {
			s += ","
		}
		s += strconv.Itoa(v)
	}
	return s
}

func ToInts(s string) []int {
	vs := strings.Split(s, ",")
	var out []int
	for _, v := range vs {
		id, err := strconv.Atoi(v)
		if err == nil {
			out = append(out, id)
		}
	}
	return out
}

func JoinInt64s(il []int64) string {
	var s string
	for _, v := range il {
		if len(s) > 0 {
			s += ","
		}
		s += strconv.FormatInt(v, 10)
	}
	return s
}

// 归纳成统一的语言代码 zh/en/hant
func FormatLanguage(s string) string {
	if strings.HasPrefix(s, "zh") && strings.Contains(s, "Hant") {
		return "hant"
	}
	if len(s) > 2 {
		s = strings.ToLower(s[:2])
	}
	return s
}

// 是中国语言
func IsZh(lang string) bool {
	return lang == "zh" || lang == "hant"
}

func FormatPlatform(p string) string {
	if strings.HasPrefix(p, "android") {
		return "android"
	} else if p == "iOS" {
		return "ios"
	}
	return ""
}

// 返回true, 如果字符串由全数字组成
func IsDigit(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// 返回true, 如果字符串是字母加数字
func IsAlnum(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z') {
			return false
		}
	}
	return true
}

// 数组转集合
func Int64sToSet(a []int64) map[int64]struct{} {
	var out = make(map[int64]struct{}, len(a))
	for _, d := range a {
		out[d] = struct{}{}
	}
	return out
}

// 截取字符串
func TruncateString(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "..."
}

// 1.0.0-35
type AppleVersion struct {
	Major int
	Minor int
	Patch int
	Build int
}

func ParseAppleVersion(v string) *AppleVersion {
	var out = &AppleVersion{}
	tup := strings.Split(v, "-")
	if len(tup) == 0 {
		return out
	}
	if len(tup) >= 2 {
		out.Build, _ = strconv.Atoi(tup[1])
	}
	fs := strings.Split(tup[0], ".")
	fn := len(fs)
	if fn > 0 {
		out.Major, _ = strconv.Atoi(fs[0])
	}
	if fn > 1 {
		out.Minor, _ = strconv.Atoi(fs[1])
	}
	if fn > 2 {
		out.Patch, _ = strconv.Atoi(fs[2])
	}
	return out
}

// 比较苹果app_version
func CmpAppleVersion(v1, v2 string) int {
	av1 := ParseAppleVersion(v1)
	av2 := ParseAppleVersion(v2)
	if av1.Major != av2.Major {
		return av1.Major - av2.Major
	}
	if av1.Minor != av2.Minor {
		return av1.Minor - av2.Minor
	}
	if av1.Patch != av2.Patch {
		return av1.Patch - av2.Patch
	}
	if av1.Build != av2.Build {
		return av1.Build - av2.Build
	}
	return 0
}

// 1.0.0
type AndroidVersion struct {
	Major int
	Minor int
	Patch int
}

func ParseAndroidVersion(v string) *AndroidVersion {
	var out = &AndroidVersion{}
	fs := strings.Split(v, ".")
	fn := len(fs)
	if fn > 0 {
		out.Major, _ = strconv.Atoi(fs[0])
	}
	if fn > 1 {
		out.Minor, _ = strconv.Atoi(fs[1])
	}
	if fn > 2 {
		out.Patch, _ = strconv.Atoi(fs[2])
	}
	return out
}

// 比较安卓app_version
func CmpAndroidVersion(v1, v2 string) int {
	av1 := ParseAndroidVersion(v1)
	av2 := ParseAndroidVersion(v2)
	if av1.Major != av2.Major {
		return av1.Major - av2.Major
	}
	if av1.Minor != av2.Minor {
		return av1.Minor - av2.Minor
	}
	if av1.Patch != av2.Patch {
		return av1.Patch - av2.Patch
	}
	return 0
}

var VALID_ADCODE_PATTERN = regexp.MustCompile("^[0-9]{6,8}$")

// 粗略检测行政区划编码格式(比如河南: 410000)
func VerifyAdCode(code string) bool {
	return VALID_ADCODE_PATTERN.MatchString(code)
}

// 格式化城市行政编码
func FormatCityCode(code string) string {
	if len(code) > 4 {
		return code[:4] + "00"
	}
	return code
}
