package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"
)

type Level int8 //日志等级

type Fields map[string]interface{}	//字段类型

//日志等级常量
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

//日志等级转为对应字符串
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelPanic:
		return "panic"
	}
	return ""
}

//日志
type Logger struct {
	newLogger *log.Logger //标准库Logger
	ctx       context.Context	//上下文
	fields    Fields	//日志字段
	callers   []string	//调用栈信息
}

//创建日志
func NewLogger(w io.Writer, prefix string, flag int) *Logger {
	//创建一个标准库Logger
	//out设置日志信息写入的目的地
	//prefix会添加到生成的每一条日志前面
	//flag定义日志的属性(时间、文件等等)
	l := log.New(w, prefix, flag)
	return &Logger{
		newLogger: l,
	}
}

//复制一个新日志结构体
func (l *Logger) clone() *Logger {
	nl := *l
	return &nl
}

//返回带字段的新日志结构体
func (l *Logger) WithFields(f Fields) *Logger {
	ll := l.clone()
	if ll.fields == nil {
		ll.fields = make(Fields)
	}
	for k, v := range f {
		ll.fields[k] = v
	}
	return ll
}

//返回带上下文的新日志结构体
func (l *Logger) WithContext(ctx context.Context) *Logger {
	ll := l.clone()
	ll.ctx = ctx
	return ll
}

//返回带某一层函数调用栈信息的新日志结构体
func (l *Logger) WithCaller(skip int) *Logger {
	ll := l.clone()
	//Caller返回当前go程调用栈所执行的函数的调用栈标识符、文件名、该调用在文件中的行号
	//skip为上溯的栈帧数，0表示Caller的调用者的调用栈
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		//返回一个表示调用栈标识符pc对应的调用栈的函数
		f := runtime.FuncForPC(pc)
		ll.callers = []string{
			fmt.Sprintf("%s: %d: %s", file, line, f.Name()),
		}
	}
	return ll
}

//返回带所有函数调用栈信息的新日志结构体
func (l *Logger) WithCallerFrames() *Logger {
	maxCallerDepth := 25
	minCallerDepth := 1
	var callers []string
	pcs := make([]uintptr, maxCallerDepth)
	//把当前go程调用栈上的调用栈标识符填入切片pc中，返回写入到pc中的项数
	//跳过Callers所在的调用栈(即跳过WithCallerFrames()该函数)
	depth := runtime.Callers(minCallerDepth, pcs)
	//得到调用栈帧信息
	frames := runtime.CallersFrames(pcs[:depth])
	//遍历所有栈帧
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		//得到调用函数的文件名,行号和函数名
		s := fmt.Sprintf("%s: %d: %s", frame.File, frame.Line, frame.Function)
		callers = append(callers, s)
		//if !more {
		//	break
		//}
	}
	ll := l.clone()
	ll.callers = callers
	return ll
}


//添加日志共同套件
//为个人修改, 主要为了减少输出函数中代码重复
func (l *Logger) WithCommonUtils(ctx context.Context) *Logger {
	return l.WithContext(ctx).WithCaller(3)
}


//日志JSON格式化
func (l *Logger) JSONFormat(level Level, message string) map[string]interface{} {
	data := make(Fields, len(l.fields)+4)
	data["level"] = level.String()	//日志等级
	data["time"] = time.Now().Local().UnixNano()	//时间戳
	data["message"] = message		//日志信息
	data["callers"] = l.callers		//调用栈信息
	//其他日志字段
	if len(l.fields) > 0 {
		for k, v := range l.fields {
			if _, ok := data[k]; !ok {
				data[k] = v
			}
		}
	}
	return data
}

//输出日志
func (l Logger) Output(level Level, message string) {
	body, _ := json.Marshal(l.JSONFormat(level, message))
	content := string(body)
	switch level {
	case LevelDebug:
		fallthrough
	case LevelInfo:
		fallthrough
	case LevelWarn:
		fallthrough
	case LevelError:
		l.newLogger.Print(content)
	case LevelFatal:
		l.newLogger.Fatal(content)
	case LevelPanic:
		l.newLogger.Panic(content)
	}
}

//不同等级的日志输出
//日志中带上下文信息以及链路跟踪信息,并进行输出
func (l *Logger) Debug(ctx context.Context, v ...interface{}) {
	l.WithCommonUtils(ctx).Output(LevelDebug, fmt.Sprint(v...))
}

func (l *Logger) Debugf(ctx context.Context, format string, v ...interface{}) {
	l.WithCommonUtils(ctx).Output(LevelDebug, fmt.Sprintf(format, v...))
}

func (l *Logger) Info(ctx context.Context, v ...interface{}) {
	l.WithCommonUtils(ctx).Output(LevelInfo,fmt.Sprint(v...))
}

func (l *Logger) Infof(ctx context.Context,format string, v ...interface{}) {
	l.WithCommonUtils(ctx).Output(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(ctx context.Context, v ...interface{}) {
	l.WithCommonUtils(ctx).Output(LevelWarn, fmt.Sprint(v...))
}

func (l *Logger) Warnf(ctx context.Context, format string, v ...interface{}) {
	l.WithCommonUtils(ctx).Output(LevelWarn, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(ctx context.Context, v ...interface{}) {
	l.WithCommonUtils(ctx).Output(LevelError, fmt.Sprint(v...))
}

func (l *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.WithCommonUtils(ctx).Output(LevelError, fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(ctx context.Context, v ...interface{}) {
	l.WithCommonUtils(ctx).Output(LevelFatal, fmt.Sprint(v...))
}

func (l *Logger) Fatalf(ctx context.Context, format string, v ...interface{}) {
	l.WithCommonUtils(ctx).Output(LevelFatal, fmt.Sprintf(format, v...))
}

func (l *Logger) Panic(ctx context.Context, v ...interface{}) {
	l.WithCommonUtils(ctx).Output(LevelPanic, fmt.Sprint(v...))
}

func (l *Logger) Panicf(ctx context.Context, format string, v ...interface{}) {
	l.WithCommonUtils(ctx).Output(LevelPanic, fmt.Sprintf(format, v...))
}
