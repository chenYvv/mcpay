// Package helpers 存放辅助方法
package helpers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	mathrand "math/rand"
	"net/http"
	"net/smtp"
	url2 "net/url"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Empty 类似于 PHP 的 empty() 函数
func Empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// MicrosecondsStr 将 time.Duration 类型（nano seconds 为单位）
// 输出为小数点后 3 位的 ms （microsecond 毫秒，千分之一秒）
func MicrosecondsStr(elapsed time.Duration) string {
	return fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
}

/*
*产生6位随机数
 */
func Rand6() int {
	randomNumber := mathrand.Intn(999999) // 生成一个0到999999之间的随机数
	return randomNumber
}

// RandomNumber 生成长度为 length 随机数字字符串
func RandomNumber(length int) string {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

// FirstElement 安全地获取 args[0]，避免 panic: runtime error: index out of range
func FirstElement(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// RandomInt 生成指定范围内的随机整数
func RandomInt(max int) int {
	if max <= 0 {
		return 0
	}
	return int(time.Now().UnixNano() % int64(max))
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[RandomInt(len(charset))]
	}
	return string(result)
}

func IsMobile(userAgent string) bool {
	if len(userAgent) == 0 {
		return false
	}

	isMobile := false
	mobileKeywords := []string{"Mobile", "Android", "Silk/", "Kindle", "BlackBerry", "Opera Mini", "Opera Mobi"}

	for i := 0; i < len(mobileKeywords); i++ {
		if strings.Contains(userAgent, mobileKeywords[i]) {
			isMobile = true
			break
		}
	}

	return isMobile
}

func CreateFile(file string) error {
	if !FileExist(file) {
		file, err := os.Create(file)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

func CreatePath(path string) error {
	if !FileExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		fmt.Println(err)
		if err != nil {
			return err
		}
	}
	return nil
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	fmt.Println(err)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// HttpPost post请求
func HttpPost(reqUrl string, header map[string]string, body []byte) ([]byte, error) {
	c := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// 添加http header
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	} else {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	return resBody, nil
}

// HttpGet get请求
func HttpGet(reqUrl string, header map[string]string) ([]byte, error, int) {
	c := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, err, 0
	}

	// 添加http header
	for k, v := range header {
		req.Header.Set(k, v)
	}

	// 请求参数
	//q := req.URL.Query()
	//q.Add("wd", "csdn")
	//req.URL.RawQuery = q.Encode()

	resp, err := c.Do(req)
	if err != nil {
		return nil, err, 0
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	return resBody, nil, resp.StatusCode
}

func Josn(obj any) []byte {
	res, _ := json.Marshal(obj)
	return res
}

func Struct2Josn(obj any) string {
	res, _ := json.Marshal(obj)
	return string(res)
}

func Json2Struct(obj any) string {
	res, _ := json.Marshal(obj)
	return string(res)
}

func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// ZeroPadding 填充零
func ZeroPadding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{0}, padding) //剩余用0填充
	return append(cipherText, padText...)

}

// ZeroUnPadding 反填充
func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

func AESEncrypt(text string, key []byte) (string, error) {
	blockSize := aes.BlockSize //AES的分组大小为16位
	src := []byte(text)
	src = ZeroPadding(src, blockSize) //填充
	out := make([]byte, len(src))
	block, err := aes.NewCipher(key) //用aes创建一个加密器cipher
	if err != nil {
		return "", err
	}
	encrypted := cipher.NewCBCEncrypter(block, key) //CBC分组模式加密
	encrypted.CryptBlocks(out, src)                 //对src进行加密，加密结果放到dst里
	return hex.EncodeToString(out), nil
}

func AESDecrypt(text string, key []byte) (string, error) {
	src, err := hex.DecodeString(text) //转为[]byte
	if err != nil {
		return "", err
	}
	out := make([]byte, len(src))
	block, err := aes.NewCipher(key) //用aes创建一个加密器cipher
	if err != nil {
		return "", err
	}
	decrypted := cipher.NewCBCDecrypter(block, key) //CBC分组模式解密
	decrypted.CryptBlocks(out, src)                 //对src进行解密，解密结果放到dst里
	out = ZeroUnPadding(out)                        //反填充
	return string(out), nil
}

func PrintJson(data interface{}) {
	b, _ := json.Marshal(data)
	dst := &bytes.Buffer{}
	_ = json.Indent(dst, b, "", " ")
	fmt.Println(dst)
}

func InArray[E any](needle E, haystack []E) bool {
	for _, item := range haystack {
		if &item == &needle {
			return true
		}
	}
	return false
}

func InArrayString(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

func InArrayInt(needle int64, haystack []int64) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func ValidateEmail(email string) bool {
	// 正则表达式来匹配电子邮箱地址
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+$`)
	return emailRegex.MatchString(email)
}

func CurlGet(url string, params map[string]interface{}, header map[string]string) ([]byte, error) {
	params_ := make([]string, 0, len(params))
	for k, v := range params {
		params_ = append(params_, fmt.Sprintf("%v=%v", k, url2.QueryEscape(fmt.Sprintf("%v", v))))
	}
	req, err := http.NewRequest("GET", url+"?"+strings.Join(params_, "&"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "content-type: text/html; charset=utf-8")

	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	client := new(http.Client)
	client.Timeout = time.Second * 10
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Str2Int(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Printf("转换出错: %v\n", err)
	}
	return num
}

func Round(num float64, precision int) float64 {
	output := strconv.FormatFloat(num, 'f', precision, 64)
	rounded, _ := strconv.ParseFloat(output, 64)
	return rounded
}

// FileExists 检查指定路径的文件是否存在
func FileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

const (
	// 邮件服务器地址
	//SMTP_MAIL_HOST = "smtp.126.com"
	SMTP_MAIL_HOST = "smtp.gmail.com"
	// 端口
	//SMTP_MAIL_PORT = "25"
	SMTP_MAIL_PORT = "587"
	// 发送邮件用户账号
	SMTP_MAIL_USER = "wangzx6666@gmail.com"
	// 授权密码
	SMTP_MAIL_PWD = "gtjpwrcwumxngvzi"
	// 发送邮件昵称
	SMTP_MAIL_NICKNAME = "SMTPMail"
)

func SendEmail(email, code string) (err error) {
	address := []string{email}
	subject := "Compass"
	body := "Your verification code is: " + code

	// 通常身份应该是空字符串，填充用户名.
	auth := smtp.PlainAuth("", SMTP_MAIL_USER, SMTP_MAIL_PWD, SMTP_MAIL_HOST)
	contentType := "Content-Type: text/html; charset=UTF-8"
	for _, v := range address {
		s := fmt.Sprintf("To:%s\r\nFrom:%s<%s>\r\nSubject:%s\r\n%s\r\n\r\n%s",
			v, SMTP_MAIL_NICKNAME, SMTP_MAIL_USER, subject, contentType, body)
		msg := []byte(s)
		addr := fmt.Sprintf("%s:%s", SMTP_MAIL_HOST, SMTP_MAIL_PORT)
		err = smtp.SendMail(addr, auth, SMTP_MAIL_USER, []string{v}, msg)
		if err != nil {
			fmt.Printf(err.Error())
		}
	}
	return
}

// 乘法
func DecimalMul(a, b float64) (res float64) {
	res, _ = decimal.NewFromFloat(a).Mul(decimal.NewFromFloat(b)).Float64()
	return
}

// 除法
func DecimalDiv(a, b float64) (res float64) {
	// 创建 Decimal 对象
	decimalDividend := decimal.NewFromFloat(a)
	decimalDivisor := decimal.NewFromFloat(b)

	// 进行除法运算
	res, _ = decimalDividend.Div(decimalDivisor).Float64()
	return
}

func StringToFloat64(str string) (res float64) {
	res, _ = strconv.ParseFloat(str, 64)
	return
}

func JsonToMap(jsonStr string) (map[string]interface{}, error) {
	var resultMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &resultMap)
	if err != nil {
		return nil, err
	}
	return resultMap, nil
}

func Float64ToString(str float64) (res string) {
	res = strconv.FormatFloat(str, 'f', -1, 64)
	return
}

func InterfaceToString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
		// 可以添加更多类型的处理
	default:
		return fmt.Sprintf("%v", v)
	}
}

// 加法
func DecimalAdd(num1, num2 float64) (res float64) {
	// 创建 Decimal 对象
	decimalNum1 := decimal.NewFromFloat(num1)
	decimalNum2 := decimal.NewFromFloat(num2)

	// 加法运算
	sum := decimalNum1.Add(decimalNum2)

	// 将结果转换为 float64
	res, _ = sum.Float64()
	return
}

func GetBase64ImageType(base64Data string) string {
	// 从 base64 数据中提取前缀
	parts := strings.Split(base64Data, ",")
	if len(parts) < 2 {
		return "" // 无效的 base64 数据
	}

	prefix := parts[0]

	// 确定图像类型
	switch {
	case strings.HasPrefix(prefix, "data:image/jpeg"):
		return "jpeg"
	case strings.HasPrefix(prefix, "data:image/png"):
		return "png"
	case strings.HasPrefix(prefix, "data:image/gif"):
		return "gif"
	case strings.HasPrefix(prefix, "data:image/webp"):
		return "webp"
	// 添加其他支持的图像类型的判断条件
	default:
		return "" // 未知类型
	}
}

func QueryAsciiSortNoEmptyNoSign(params map[string]interface{}) string {
	keys := make([]string, 0, len(params))
	values := make([]string, 0, len(params))
	for k := range params {
		if params[k] != "" && params[k] != nil && k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, v := range keys {
		values = append(values, fmt.Sprintf("%v=%v", v, InterfaceToString(params[v])))
	}
	return strings.Join(values, "&")
}

func Fen2Yuan(a float64) float64 {
	return DecimalDiv(a, 100)
}

func Yuan2Fen(a float64) float64 {
	return DecimalMul(a, 100)
}

func GenerateOrderNo(prefix string) string {
	src := mathrand.NewSource(time.Now().UnixNano())
	r := mathrand.New(src)
	timestamp := time.Now().Format("20060102150405") // yyyyMMddHHmmss
	millis := time.Now().UnixNano() / 1e6
	random := r.Intn(100000) // 0 ~ 99999

	return fmt.Sprintf("%s%s%d%05d", prefix, timestamp, millis%1000, random)
}

func QueryAsciiSortNoEmpty(params map[string]interface{}) string {
	keys := make([]string, 0, len(params))
	values := make([]string, 0, len(params))
	for k := range params {
		if params[k] != "" && params[k] != nil {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, v := range keys {
		values = append(values, fmt.Sprintf("%v=%v", v, InterfaceToString(params[v])))
	}
	return strings.Join(values, "&")
}
