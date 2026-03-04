package log

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	fiberlog "github.com/gofiber/fiber/v2/log"
	"github.com/rotisserie/eris"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	levelTrace = iota
	levelDebug
	levelInfo
	levelWarn
	levelError
	levelFatal
	levelPanic
)

type LoggerConfig interface {
	GetMaxSize() string
	GetMaxBackups() string
	GetLifeTime() string
	GetDirSave() string
	GetFilename() string
}

type Logger struct {
	Logger     *log.Logger
	Writer     io.Writer
	lumberjack *lumberjack.Logger // Сохраняем ссылку на lumberjack

	level int
}

func (lgr *Logger) logf(lvl int, format string, v ...interface{}) {
	if lvl < lgr.level {
		return
	}

	prefix := ""
	switch lvl {
	case levelTrace:
		prefix = "[TRACE] "
	case levelDebug:
		prefix = "[DEBUG] "
	case levelInfo:
		prefix = "[INFO]  "
	case levelWarn:
		prefix = "[WARN]  "
	case levelError:
		prefix = "[ERROR] "
	case levelFatal:
		prefix = "[FATAL] "
	case levelPanic:
		prefix = "[PANIC] "
	}

	msg := fmt.Sprintf(prefix+format, v...)
	lgr.Logger.Output(3, msg) // 3 — чтобы показать правильный caller (не наш метод)
}

func (lgr *Logger) log(lvl int, v ...interface{}) {
	if lvl < lgr.level {
		return
	}
	lgr.logf(lvl, "%v", v...)
}

// ────────────────────────────────────────────────
// Простые методы (без форматирования)
func (lgr *Logger) Trace(v ...interface{}) {
	lgr.log(levelTrace, v...)
}

func (lgr *Logger) Debug(v ...interface{}) {
	lgr.log(levelDebug, v...)
}

func (lgr *Logger) Info(v ...interface{}) {
	lgr.log(levelInfo, v...)
}

func (lgr *Logger) Warn(v ...interface{}) {
	lgr.log(levelWarn, v...)
}

func (lgr *Logger) Error(v ...interface{}) {
	lgr.log(levelError, v...)
}

func (lgr *Logger) Fatal(v ...interface{}) {
	lgr.log(levelFatal, v...)
	os.Exit(1)
}

func (lgr *Logger) Panic(v ...interface{}) {
	lgr.log(levelPanic, v...)
	panic(fmt.Sprint(v...))
}

// ────────────────────────────────────────────────
// Форматированные методы
func (lgr *Logger) Tracef(format string, v ...interface{}) {
	lgr.logf(levelTrace, format, v...)
}

func (lgr *Logger) Debugf(format string, v ...interface{}) {
	lgr.logf(levelDebug, format, v...)
}

func (lgr *Logger) Infof(format string, v ...interface{}) {
	lgr.logf(levelInfo, format, v...)
}

func (lgr *Logger) Warnf(format string, v ...interface{}) {
	lgr.logf(levelWarn, format, v...)
}

func (lgr *Logger) Errorf(format string, v ...interface{}) {
	lgr.logf(levelError, format, v...)
}

func (lgr *Logger) Fatalf(format string, v ...interface{}) {
	lgr.logf(levelFatal, format, v...)
	os.Exit(1)
}

func (lgr *Logger) Panicf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	lgr.logf(levelPanic, "%s", msg)
	panic(msg)
}

// ────────────────────────────────────────────────
// Structured-стиль (keysAndValues)
func (lgr *Logger) logw(lvl int, msg string, keysAndValues ...interface{}) {
	if lvl < lgr.level {
		return
	}

	var sb strings.Builder
	sb.WriteString(msg)

	if len(keysAndValues) > 0 {
		sb.WriteString(" ")
		for i := 0; i < len(keysAndValues); i += 2 {
			if i+1 < len(keysAndValues) {
				sb.WriteString(fmt.Sprintf("%v=%v ", keysAndValues[i], keysAndValues[i+1]))
			} else {
				sb.WriteString(fmt.Sprintf("%v ", keysAndValues[i]))
			}
		}
	}

	lgr.logf(lvl, "%s", sb.String())
}

func (lgr *Logger) Tracew(msg string, keysAndValues ...interface{}) {
	lgr.logw(levelTrace, msg, keysAndValues...)
}

func (lgr *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	lgr.logw(levelDebug, msg, keysAndValues...)
}

func (lgr *Logger) Infow(msg string, keysAndValues ...interface{}) {
	lgr.logw(levelInfo, msg, keysAndValues...)
}

func (lgr *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	lgr.logw(levelWarn, msg, keysAndValues...)
}

func (lgr *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	lgr.logw(levelError, msg, keysAndValues...)
}

func (lgr *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	lgr.logw(levelFatal, msg, keysAndValues...)
	os.Exit(1)
}

func (lgr *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	lgr.logw(levelPanic, msg, keysAndValues...)
	panic(msg)
}

// ────────────────────────────────────────────────
// Методы из ControlLogger / ConfigurableLogger
func (lgr *Logger) SetLevel(level fiberlog.Level) {
	// fiberlog.Level: Trace=0, Debug=1, Info=2, Warn=3, Error=4, Fatal=5, Panic=6
	lgr.level = int(level)
}

func (lgr *Logger) SetOutput(writer io.Writer) {
	lgr.Logger.SetOutput(writer)
	lgr.Writer = writer
}

// WithContext — в большинстве случаев просто возвращает себя
// (если нужен контекстный логгер — можно доработать позже)
func (lgr *Logger) WithContext(ctx context.Context) fiberlog.CommonLogger {
	// Можно добавить поля из контекста, если они есть
	// пока просто возвращаем тот же логгер
	return lgr
}

func InitLogger(cfg LoggerConfig) *Logger {
	filename := "./" + cfg.GetDirSave() + "/" + cfg.GetFilename()

	maxSize, err := strconv.Atoi(cfg.GetMaxSize())
	if err != nil {
		log.Println(err)
	}

	maxBackups, err := strconv.Atoi(cfg.GetMaxBackups())
	if err != nil {
		log.Println(err)
	}

	maxAge, err := strconv.Atoi(cfg.GetLifeTime())
	if err != nil {
		log.Println(err)
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
	}

	// Создаем мультирайтер, который будет писать и в файл, и в консоль
	multiWriter := io.MultiWriter(lumberjackLogger, os.Stdout)

	// Создаем новый логгер с мультирайтером
	logger := log.New(multiWriter, "", log.LstdFlags)

	internalLogger := &Logger{
		Logger:     logger,
		Writer:     multiWriter,
		lumberjack: lumberjackLogger, // Сохраняем ссылку
	}

	fiberlog.SetLogger(internalLogger)

	return internalLogger
}

func (lgr *Logger) CloseLogger() error {
	// Останавливаем Lumberjack (это остановит горутину millRun)
	if lgr.lumberjack != nil {
		if err := lgr.lumberjack.Close(); err != nil {
			return eris.Wrap(err, "failed to close lumberjack logger")
		}
	}

	// Закрываем другие writers если нужно
	if closer, ok := lgr.Writer.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			return eris.Wrap(err, "failed to close Writer")
		}
	}

	return nil
}
