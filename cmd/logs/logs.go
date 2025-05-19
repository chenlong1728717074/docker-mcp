package logs

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

const (
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorReset   = "\033[0m"
)

func init() {
	// 创建基础配置
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    customLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeCaller:   customCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	// 创建 rotatelogs 实例
	logWriter, err := rotatelogs.New(
		"./logs/app-%Y-%m-%d.log",                 // 日志文件按日期命名
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 保留7天的日志
		rotatelogs.WithRotationTime(24*time.Hour), // 每天轮转一次
	)
	if err != nil {
		panic("failed to create rotatelogs: " + err.Error())
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevelAt(zap.DebugLevel)

	// 创建自定义编码器
	consoleEncoder := NewCustomEncoder(encoderConfig, true) // 控制台带颜色
	fileEncoder := NewCustomEncoder(encoderConfig, false)

	// 创建核心配置
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), atomicLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logWriter), atomicLevel),
	)

	// 创建Logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	Sugar = Logger.Sugar()
}

// 获取日志级别对应的颜色
func getLevelColor(level zapcore.Level) string {
	switch level {
	case zapcore.DebugLevel:
		return colorGreen
	case zapcore.InfoLevel:
		return colorBlue
	case zapcore.WarnLevel:
		return colorYellow
	case zapcore.ErrorLevel:
		return colorRed
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return colorMagenta
	default:
		return colorReset
	}
}

// 自定义编码器
type CustomEncoder struct {
	zapcore.Encoder
	withColor bool
	pool      buffer.Pool
}

func NewCustomEncoder(cfg zapcore.EncoderConfig, withColor bool) zapcore.Encoder {
	return &CustomEncoder{
		Encoder:   zapcore.NewJSONEncoder(cfg),
		pool:      buffer.NewPool(),
		withColor: withColor,
	}
}

// 自定义时间编码器
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// 自定义级别编码器
func customLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(strings.ToUpper(l.String()))
}

// 自定义调用者编码器
func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%d %s", os.Getpid(), caller.TrimmedPath()))
}

// 实现 EncodeEntry
func (e *CustomEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf := e.pool.Get()

	// 时间
	buf.AppendString(ent.Time.Format("2006-01-02 15:04:05.000"))
	buf.AppendByte(' ')

	// 带颜色的日志级别
	if e.withColor {
		levelColor := getLevelColor(ent.Level)
		buf.AppendString(levelColor)
		buf.AppendString(strings.ToUpper(ent.Level.String()))
		buf.AppendString(colorReset)
	} else {
		buf.AppendString(strings.ToUpper(ent.Level.String()))
	}
	buf.AppendByte(' ')

	// 调用者信息
	if ent.Caller.Defined {
		buf.AppendString(ent.Caller.TrimmedPath())
	}
	buf.AppendString(" : ")

	// 消息
	buf.AppendString(ent.Message)
	buf.AppendByte('\n')

	return buf, nil
}

// 基本日志方法
func Info(msg string, fields ...interface{}) {
	Sugar.Infof(msg, fields...)
}

func Error(msg string, fields ...interface{}) {
	Sugar.Errorf(msg, fields...)
}

func Debug(msg string, fields ...interface{}) {
	Sugar.Debugf(msg, fields...)
}

func Warn(msg string, fields ...interface{}) {
	Sugar.Warnf(msg, fields...)
}

func Fatal(msg string, fields ...interface{}) {
	Sugar.Fatalf(msg, fields...)
}

// 结构化日志方法
func InfoWithFields(msg string, fields map[string]interface{}) {
	Sugar.Infow(msg, fieldsToArgs(fields)...)
}

func ErrorWithFields(msg string, fields map[string]interface{}) {
	Sugar.Errorw(msg, fieldsToArgs(fields)...)
}

func DebugWithFields(msg string, fields map[string]interface{}) {
	Sugar.Debugw(msg, fieldsToArgs(fields)...)
}

func WarnWithFields(msg string, fields map[string]interface{}) {
	Sugar.Warnw(msg, fieldsToArgs(fields)...)
}

func FatalWithFields(msg string, fields map[string]interface{}) {
	Sugar.Fatalw(msg, fieldsToArgs(fields)...)
}

// 对象日志方法
func InfoObject(msg string, key string, obj interface{}) {
	Sugar.Infow(msg, key, obj)
}

func ErrorObject(msg string, key string, obj interface{}) {
	Sugar.Errorw(msg, key, obj)
}

func DebugObject(msg string, key string, obj interface{}) {
	Sugar.Debugw(msg, key, obj)
}

func WarnObject(msg string, key string, obj interface{}) {
	Sugar.Warnw(msg, key, obj)
}

func FatalObject(msg string, key string, obj interface{}) {
	Sugar.Fatalw(msg, key, obj)
}

// 多对象日志方法
func InfoObjects(msg string, keysAndObjects ...interface{}) {
	Sugar.Infow(msg, keysAndObjects...)
}

func ErrorObjects(msg string, keysAndObjects ...interface{}) {
	Sugar.Errorw(msg, keysAndObjects...)
}

func DebugObjects(msg string, keysAndObjects ...interface{}) {
	Sugar.Debugw(msg, keysAndObjects...)
}

func WarnObjects(msg string, keysAndObjects ...interface{}) {
	Sugar.Warnw(msg, keysAndObjects...)
}

func FatalObjects(msg string, keysAndObjects ...interface{}) {
	Sugar.Fatalw(msg, keysAndObjects...)
}

// 将map转换为zap的参数列表
func fieldsToArgs(fields map[string]interface{}) []interface{} {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return args
}
