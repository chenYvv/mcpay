package log

import (
	"bufio"
	"errors"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v5"
	elogrus "gopkg.in/sohlich/elogrus.v2"
	"net"
	"os"
	"path"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	sysTime "time"
)

// 日志等级
const (
	C_Kafka_Log_Topic      = "all-server-log-test2"
	C_Kafka_Log_User_Topic = "all-user-log-test2"
	LEVEL_FATA             = 1
	LEVEL_ERROR            = 2
	LEVEL_WARNING          = 3
	LEVEL_INFO             = 4
	LEVEL_DEBUG            = 6
)

// 日志模式
const (
	MODEL_PRO = iota
	MODEL_INFO
	MODEL_DEV
)

// 调用log的服务器名字
var serverName string
var gIp net.IP
var myLog = logrus.New()
var Logger = logrus.New()
var address string
var isOpen = false

func init() {
	ip, err := externalIP()
	if err == nil {
		gIp = ip
	} else {
		gIp = nil
	}
	Logger.SetLevel(logrus.TraceLevel)
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000000",
		ForceColors:     true,
	})
	myLog.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000000",
		ForceColors:     true,
	})
	myLog.SetOutput(os.Stdout)
	Logger.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	ConfigLocalFilesystemLogger("./logs", "log.out", sysTime.Hour*24*15, sysTime.Hour*24)
}

func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge sysTime.Duration, rotationTime sysTime.Duration) {
	baseLogPath := path.Join(logPath, logFileName)
	fmt.Println(baseLogPath)
	writer, err := rotatelogs.New(
		logPath+"/%Y%m%d/%H%M_"+logFileName,
		rotatelogs.WithLinkName(baseLogPath), // 生成软链，指向最新日志文件
		//rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		//rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
		rotatelogs.WithRotationSize(1024*1024*10), //10m
		rotatelogs.WithRotationCount(5),
	)
	if err != nil {
		//log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
		panic(err)
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		//logrus.TraceLevel: writer,
		//logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		//logrus.InfoLevel:  writer,
		//logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000000",
		ForceColors:     true,
	})
	Logger.AddHook(lfHook)
	myLog.AddHook(lfHook)
}

func SetLogName(name string) {
	serverName = name
}

func OpenSendLog(name string, open bool, elkAddress string) {
	serverName = name
	isOpen = open
	address = elkAddress
	if isOpen {
		ConfigESLogger(address, gIp.String(), serverName)
	}
}

var gLvl = 2

func callerPrettyFile() string {
	fName := ""
	pc, pathT, line, ok := runtime.Caller(gLvl) // 去掉两层，当前函数和日志的接口函数
	if ok {
		if f := runtime.FuncForPC(pc); f != nil {
			fName = f.Name()
		}
	}
	funcName := lastFname(fName)
	//pathT = getFilePath(pathT)
	pathT = getFilePathV2(pathT)
	return fmt.Sprintf("%s() %s:%d ", funcName, pathT, line)
}

func ConfigESLogger(esUrl string, esHOst string, index string) error {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(esUrl))
	if err != nil {
		Logger.Error(err.Error())
		return err
	}
	esHook, err := elogrus.NewElasticHook(client, esHOst, logrus.WarnLevel, index)
	if err != nil {
		Logger.Error(err.Error())
		return err
	}
	Logger.AddHook(esHook)

	esUserHook, err := elogrus.NewElasticHook(client, esHOst, logrus.WarnLevel, "用户操作记录测试")
	if err != nil {
		Logger.Error(err.Error())
		return err
	}
	myLog.AddHook(esUserHook)

	return nil
}

// SetLogLevel 设置日志等级
func SetLogLevel(lev logrus.Level) {
	Logger.SetLevel(lev)
}

func SetPathLvl(lvl int) { gLvl = lvl }

// SetLogModel /*
func SetLogModel(mod int) error {
	if mod <= MODEL_DEV {
		myLog.SetOutput(os.Stdout)
		Logger.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	}
	if mod <= MODEL_INFO {
		src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			Logger.Error(err.Error())
			return err
		}
		writer := bufio.NewWriter(src)
		myLog.SetOutput(writer)
		Logger.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	}
	if mod <= MODEL_PRO {
		src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			Logger.Error(err.Error())
			return err
		}
		writer := bufio.NewWriter(src)
		Logger.SetOutput(writer)
		myLog.SetOutput(writer)
	}
	return nil
}

// Listen 安全执行监听函数
func Listen(f interface{}, callback func(interface{}), param string) {
	fName := ""
	pc, pathT, line, ok := runtime.Caller(1) // 去掉两层，当前函数和日志的接口函数
	if ok {
		if f := runtime.FuncForPC(pc); f != nil {
			fName = f.Name()
		}
	}
	funcName := lastFname(fName)
	pathT = getFilePath(pathT)
	timer := sysTime.NewTicker(sysTime.Millisecond * 500)
	defer timer.Stop()
	success := make(chan bool)
	start := sysTime.Now()
	count := 0
	go func() {
		callback(f)
		close(success)
	}()
	for {
		select {
		case <-success:
			end := sysTime.Since(start).Nanoseconds() / 1000000.00

			if end >= 500 && end < 1000 {
				Logger.Info(fmt.Sprintf("执行严重超时 %s %s %d (%dms) %s", pathT, funcName, line, end, param))
			}
			if end >= 1000 && end < 2000 {
				Logger.Warn(fmt.Sprintf("执行严重超时 %s %s %d (%dms) %s", pathT, funcName, line, end, param))
			}
			if end >= 2000 {
				Logger.Error(fmt.Sprintf("执行严重超时 %s %s %d (%dms) %s", pathT, funcName, line, end, param))
			}
			return
		case <-timer.C:
			count++
			end := sysTime.Since(start).Nanoseconds() / 1000000.00
			if count >= 10 {
				Logger.Error(fmt.Sprintf("执行严重超时%d次提醒 %s %s %d (%dms) %s", count, pathT, funcName, line, end, param))
			} else {
				Logger.Info(fmt.Sprintf("执行严重超时%d次提醒 %s %s %d (%dms) %s", count, pathT, funcName, line, end, param))
			}
		}
	}
}

// TraceParam 计算函数所用时间
func TraceParam(param string) func() {
	fName := ""
	pc, pathT, line, ok := runtime.Caller(1) // 去掉两层，当前函数和日志的接口函数
	if ok {
		if f := runtime.FuncForPC(pc); f != nil {
			fName = f.Name()
		}
	}
	funcName := lastFname(fName)
	pathT = getFilePath(pathT)

	start := sysTime.Now()
	return func() {
		end := sysTime.Since(start).Nanoseconds() / 1000000.00
		if end >= 100 && end < 1000 {
			Debug("执行严重超时提醒 %s %s %d (%dms) %s", pathT, funcName, line, end, param)
		}
		if end >= 1000 && end < 2000 {
			Warn("执行严重超时提醒 %s %s %d (%dms) %s", pathT, funcName, line, end, param)
		}
		if end >= 2000 {
			Error("执行严重超时提醒 %s %s %d (%dms) %s", pathT, funcName, line, end, param)
		}
	}
}

// Trace 计算函数所用时间
func Trace() func() {
	fName := ""
	pc, pathT, line, ok := runtime.Caller(1) // 去掉两层，当前函数和日志的接口函数
	if ok {
		if f := runtime.FuncForPC(pc); f != nil {
			fName = f.Name()
		}
	}
	funcName := lastFname(fName)
	pathT = getFilePath(pathT)

	start := sysTime.Now()
	return func() {
		end := sysTime.Since(start).Nanoseconds() / 1000000.00
		if end >= 100 && end < 1000 {
			Info("执行严重超时提醒 %s %s %d (%dms)", pathT, funcName, line, end)
		}
		if end >= 1000 && end < 2000 {
			Warn("执行严重超时提醒 %s %s %d (%dms)", pathT, funcName, line, end)
		}
		if end >= 2000 {
			Error("执行严重超时提醒 %s %s %d (%dms)", pathT, funcName, line, end)
		}
	}
}

// TraceInfo 计算函数所用时间
func TraceInfo(str string) func() {
	fName := ""
	pc, pathT, line, ok := runtime.Caller(1) // 去掉两层，当前函数和日志的接口函数
	if ok {
		if f := runtime.FuncForPC(pc); f != nil {
			fName = f.Name()
		}
	}
	funcName := lastFname(fName)
	pathT = getFilePath(pathT)

	start := sysTime.Now()
	return func() {
		end := sysTime.Since(start).Nanoseconds() / 1000000.00
		if end >= 100 && end < 1000 {
			Info("%s 执行严重超时提醒 %s %s %d (%dms)", str, pathT, funcName, line, end)
		}
		if end >= 1000 && end < 2000 {
			Warn("%s 执行严重超时提醒 %s %s %d (%dms)", str, pathT, funcName, line, end)
		}
		if end >= 2000 {
			Error("%s 执行严重超时提醒 %s %s %d (%dms)", str, pathT, funcName, line, end)
		}
	}
}

func GetFunctionName(i interface{}, seps ...rune) string {
	// 获取函数名称
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	// 用 seps 进行分割
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	// fmt.Println(fields)

	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return ""
}

func UserInfoLog(userId int64, format string, args ...interface{}) {
	//defer Trace()()
	myLog.WithFields(logrus.Fields{
		"UserId": userId,
	}).Warn(callerPrettyFile() + fmt.Sprintf(format, args...))
}

// Fatal 危险的
func Fatal(format string, args ...interface{}) {
	//defer Trace()()
	Logger.Fatal(callerPrettyFile() + fmt.Sprintf(format, args...))
}

// Error 错误
func Error(format string, args ...interface{}) {
	//defer Trace()()
	Logger.Error(callerPrettyFile() + fmt.Sprintf(format, args...))
}

// Warn 警告
func Warn(format string, args ...interface{}) {
	//defer Trace()()
	Logger.Warn(callerPrettyFile() + fmt.Sprintf(format, args...))
}

// Info 提示
func Info(format string, args ...interface{}) {
	//defer Trace()()
	Logger.Info(callerPrettyFile() + fmt.Sprintf(format, args...))
}

// Debug 调试
func Debug(format string, args ...interface{}) {
	//defer Trace()()
	Logger.Debug(callerPrettyFile() + fmt.Sprintf(format, args...))
}

func PanicLn(args ...interface{}) {
	Logger.Errorln(" stack:", string(debug.Stack()), " args:", args)
	Logger.Panicln(args...)
}

func Test(format string, args ...interface{}) {
	//defer Trace()()
	Logger.Trace(callerPrettyFile() + fmt.Sprintf(format, args...))
}

func TraceError(x any) {
	Logger.Errorln(" stack:", string(debug.Stack()), " err:", x)
}

// func logFormat(msg string) string {
// 	fname := ""
// 	pc, path, line, ok := runtime.Caller(2) // 去掉两层，当前函数和日志的接口函数
// 	if ok {
// 		if f := runtime.FuncForPC(pc); f != nil {
// 			fname = f.Name()
// 		}
// 	}
// 	funcName := lastFname(fname)
// 	path = getFilePath(path)
// 	format := fmt.Sprintf(" %s %s %d ▶ %s", path, funcName, line, msg)
// 	//fmt.Println(format)
// 	return format
// }

func lastFname(fname string) string {
	flen := len(fname)
	n := strings.LastIndex(fname, ".")
	if n+1 < flen {
		return fname[n+1:]
	}
	return fname
}

func getFilePath(path string) string {
	s := strings.Split(path, "src")
	return s[0]
}

func getFilePathV2(file string) string {
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 3 {
				file = file[i+1:]
				break
			}
		}
	}
	return file
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
