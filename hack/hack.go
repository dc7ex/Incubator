package hack

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
)

func Getenv(key string, def interface{}) interface{} {
	// 异常处理
	defer func() {
		if err := recover(); err != nil {
			os.Exit(1)
		}
	}()

	value := os.Getenv(key)

	if len(value) == 0 {
		return def
	} else {

		var val interface{}
		var err error

		switch def.(type) {
		case bool:
			val, err = strconv.ParseBool(value)
		case int:
			val, err = strconv.Atoi(value)
		case int64:
			val, err = strconv.ParseInt(value, 10, 64)
		case float32:
			val, err = strconv.ParseFloat(value, 32)
		case float64:
			val, err = strconv.ParseFloat(value, 64)
		default:
			val = value
		}

		if err != nil {
			panic(err)
		}

		return val
	}
}

func SignalChannel(call func()) {
	// 信号处理
	sc := make(chan os.Signal, 1)
	signal.Notify(
		sc,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGPIPE,
	)

	//  go func() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Error: ", err)
		}
	}()

	for {
		select {
		case sig := <-sc:
			if sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == syscall.SIGQUIT || sig == syscall.SIGHUP {
				log.Println("Shut down all kinds of connection")
				call()
				return
			} else if sig == syscall.SIGPIPE {
				log.Println("Ignore broken pipe signal")
			}
		}
	}
	//  }()
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Isset(arr []interface{}, index int) bool {
	return (len(arr) > index)
}

func StringHexArray(str string) []string {
	var buff []string
	var i int = 0
	var arr string = ""
	var hit bool = false

	for _, item := range str {
		i++
		if i > 1 {
			i = 0
			hit = true
		} else {
			hit = false
		}

		arr = arr + string(item)
		if hit {
			buff = append(buff, arr)
		}

		if i == 0 {
			arr = ""
		}
	}

	return buff
}

func Sha1s(s string) string {
	r := sha1.Sum([]byte(s))
	return hex.EncodeToString(r[:])
}

func Sha256s(s string) string {
	r := sha256.Sum256([]byte(s))
	return hex.EncodeToString(r[:])
}

func Uuid() string {
	return uuid.New().String()
}

func Md5s(s string) string {
	r := md5.Sum([]byte(s))
	return hex.EncodeToString(r[:])
}

func DayTimeSeries() int {
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	timeStr := time.Now().In(cstSh).Format("20060102")
	timeSeries, err := strconv.Atoi(timeStr)
	if err != nil {
		return 0
	}
	return timeSeries
}

func TimestampToDate(ts string) string {
	tsint64 := ToInt64(ts)
	if tsint64 == 0 {
		return ""
	}
	l, _ := time.LoadLocation("Asia/Shanghai")
	tm := time.Unix(tsint64, 0)
	return tm.In(l).Format("2006-01-02 15:04:05")
}

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

func StringToInt64(s string) int64 {
	if s == "" {
		return 0
	}

	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return val
}

func Unshift(slice, v interface{}) (interface{}, error) {
	var typ = reflect.TypeOf(slice)
	if typ.Kind() == reflect.Slice {
		var vv = reflect.ValueOf(slice)
		var tmp = reflect.MakeSlice(typ, vv.Len()+1, vv.Cap()+1)
		tmp.Index(0).Set(reflect.ValueOf(v))
		var dst = tmp.Slice(1, tmp.Len())
		reflect.Copy(dst, vv)
		return tmp.Interface(), nil
	}
	return nil, errors.New(`not a slice`)
}

func ToInt64(t interface{}) int64 {
	var i int64 = 0

	switch t.(type) {
	case int64:
		i = t.(int64)
		break
	case int:
		i = int64(t.(int))
		break
	case float64:
		i = int64(t.(float64))
		break
	case float32:
		i = int64(t.(float32))
		break
	case string:
		i = StringToInt64(t.(string))
		break
	}

	return i
}

func ToInt(t interface{}) int {
	var i int = 0

	switch t.(type) {
	case int64:
		//		strInt64 := strconv.FormatInt(t.(int64), 10)
		//		i, _ := strconv.Atoi(strInt64)
		i = int(t.(int64))
		break
	case int:
		i = t.(int)
		break
	case float64:
		i = int(t.(float64))
		break
	case float32:
		i = int(t.(float32))
		break
	case string:
		i, _ = strconv.Atoi(t.(string))
		break
	}

	return i
}

func ToFloat64(t interface{}) float64 {
	var i float64 = 0

	switch t.(type) {
	case int64:
		i = float64(t.(int64))
		break
	case int:
		i = float64(t.(int))
		break
	case float64:
		i = t.(float64)
		break
	case float32:
		i = float64(t.(float32))
		break
	case string:
		i, _ = strconv.ParseFloat(t.(string), 64)
		break
	}

	return i
}

func ToString(t interface{}) string {
	var i string = ""

	switch t.(type) {
	case int64:
		i = strconv.FormatInt(t.(int64), 10)
		break
	case int:
		i = strconv.Itoa(t.(int))
		break
	case float64:
		i = strconv.FormatFloat(t.(float64), 'f', -1, 64)
		break
	case float32:
		i = strconv.FormatFloat(float64(t.(float32)), 'f', -1, 32)
		break
	case string:
		i = t.(string)
		break
	}

	return i
}

func Int64InSlice(a int64, list []int64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func DateTimeSeries() int64 {
	l, _ := time.LoadLocation("Asia/Shanghai")
	timeStr := time.Now().In(l).Format("20060102150405")
	return StringToInt64(timeStr)
}

func TimestampToTime(ts int64) int64 {
	l, _ := time.LoadLocation("Asia/Shanghai")
	tm := time.Unix(ts, 0)
	timeStr := tm.In(l).Format("20060102150405")
	return StringToInt64(timeStr)
}

func StringToTimestamp(val string) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")                   //设置时区
	tt, _ := time.ParseInLocation("2006-01-02 15:04:05", val, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
	return tt.Unix()
}

func MapStringMerge(list ...map[string]string) (result map[string]string) {
	result = make(map[string]string)
	for _, item := range list {
		for k, v := range item {
			result[k] = v
		}
	}
	return
}

func TimestampToString(ts int64) string {
	l, _ := time.LoadLocation("Asia/Shanghai")
	tm := time.Unix(ts, 0)
	return tm.In(l).Format("2006-01-02 15:04:05")
}

// 格式化后的UTC时区转CST
func ToFormatTimeZone(template string, from string, value string) string {
	loc, _ := time.LoadLocation(from)                        //重要：获取时区
	theTime, _ := time.ParseInLocation(template, value, loc) //使用模板在对应时区转化为time.time类型
	l, _ := time.LoadLocation("Asia/Shanghai")
	return theTime.In(l).Format("2006-01-02 15:04:05")
}

func MapToJson(param map[string]interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func NewTimestamp() int64 {
	l, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(l).Unix()
}

func ListStringToJson(param []string) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func NewDateOffset(years, months, days int) int {
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	timeStr := time.Now().AddDate(years, months, days).In(cstSh).Format("20060102")
	timeSeries, err := strconv.Atoi(timeStr)
	if err != nil {
		return 0
	}
	return timeSeries
}

func YearsToString() string {
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(cstSh).Format("2006")
}

func DateToString() string {
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(cstSh).Format("20060102")
}
