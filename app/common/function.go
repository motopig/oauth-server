package common

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Map map[string]interface{}

func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

func Float2str(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func Str2int(str string) int {
	val, _ := strconv.Atoi(str)
	return val
}

func Str2int32(str string) int32 {
	i, _ := strconv.ParseInt(str, 10, 32)
	return int32(i)
}

func Str2int64(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}

func Hex2int32(str string) int32 {
	val, _ := strconv.ParseInt(str, 16, 32)
	return int32(val)
}

func Hex2int64(str string) int64 {
	val, _ := strconv.ParseInt(str, 16, 64)
	return val
}

func Str2float64(str string) float64 {
	val, _ := strconv.ParseFloat(str, 64)
	return val
}

func Int2str(i int) string {
	return strconv.Itoa(i)
}

func Int2Hex(num int64) string {
	return strconv.FormatInt(num, 16)
}

func Int32str(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

func Int642Str(i int64) string {
	return strconv.FormatInt(i, 10)
}

func String2Interface(str []string) []interface{} {

	s := make([]interface{}, len(str))
	for i, v := range str {
		s[i] = v
	}
	return s
}

func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func ParseEmail(str string) string {
	if VerifyEmailFormat(str) {
		strs := strings.Split(str, "@")
		return strs[0]
	}
	return str
}

func Md5V2(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func PasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SmartPrint(i interface{}) {
	var kv = make(map[string]interface{})
	vValue := reflect.ValueOf(i)
	vType := reflect.TypeOf(i)
	for i := 0; i < vValue.NumField(); i++ {
		kv[vType.Field(i).Name] = vValue.Field(i)
	}
	fmt.Println("获取到数据:")
	for k, v := range kv {
		fmt.Print(k)
		fmt.Print(":")
		fmt.Print(v)
		fmt.Println()
	}
}
