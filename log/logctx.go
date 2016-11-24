package log

import (
	"fmt"
	"runtime"

	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// 生成日志跟踪的唯一标识
func getTracerIDFromCtx(ctx context.Context) string {
	guid := "00000000-0000-0000-0000-000000000000"

	if ctx == nil {
		return guid
	}

	if meta, ok := metadata.FromContext(ctx); ok {
		if meta["tid"] != nil && len(meta["tid"]) > 0 {
			return meta["tid"][0]
		}
	}
	return guid
}

// CtxDebugf 包含上下文的Debug日志
func (l *Logger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	if Ldebug < l.Level {
		return
	}

	l.Output(getTracerIDFromCtx(ctx), Ldebug, 2, fmt.Sprintf(format, v...))
}

// CtxDebug 包含上下文的Debug日志
func (l *Logger) CtxDebug(ctx context.Context, v ...interface{}) {
	if Ldebug < l.Level {
		return
	}

	l.Output(getTracerIDFromCtx(ctx), Ldebug, 2, fmt.Sprintln(v...))
}

// CtxInfof 包含上下文的Info日志
func (l *Logger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	if Linfo < l.Level {
		return
	}
	l.Output(getTracerIDFromCtx(ctx), Linfo, 2, fmt.Sprintf(format, v...))
}

// CtxInfo 包含上下文的Info日志
func (l *Logger) CtxInfo(ctx context.Context, v ...interface{}) {
	if Linfo < l.Level {
		return
	}
	l.Output(getTracerIDFromCtx(ctx), Linfo, 2, fmt.Sprintln(v...))
}

// CtxWarnf 包含上下文的Warn日志
func (l *Logger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	l.Output(getTracerIDFromCtx(ctx), Lwarn, 2, fmt.Sprintf(format, v...))
}

// CtxWarn 包含上下文的Warn日志
func (l *Logger) CtxWarn(ctx context.Context, v ...interface{}) {
	l.Output(getTracerIDFromCtx(ctx), Lwarn, 2, fmt.Sprintln(v...))
}

// CtxErrorf 包含上下文的Error日志
func (l *Logger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	l.Output(getTracerIDFromCtx(ctx), Lerror, 2, fmt.Sprintf(format, v...))
}

// CtxError 包含上下文的Error日志
func (l *Logger) CtxError(ctx context.Context, v ...interface{}) {
	l.Output(getTracerIDFromCtx(ctx), Lerror, 2, fmt.Sprintln(v...))
}

// CtxFatal 包含上下文的Fatal日志
func (l *Logger) CtxFatal(ctx context.Context, v ...interface{}) {
	l.Output(getTracerIDFromCtx(ctx), Lfatal, 2, fmt.Sprint(v...))
	os.Exit(1)
}

// CtxFatalf 包含上下文的Fatal日志
func (l *Logger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	l.Output(getTracerIDFromCtx(ctx), Lfatal, 2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// CtxFatalln 包含上下文的Fatal日志
func (l *Logger) CtxFatalln(ctx context.Context, v ...interface{}) {
	l.Output(getTracerIDFromCtx(ctx), Lfatal, 2, fmt.Sprintln(v...))
	os.Exit(1)
}

// CtxPanic 包含上下文的Panic日志
func (l *Logger) CtxPanic(ctx context.Context, v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Output(getTracerIDFromCtx(ctx), Lpanic, 2, s)
	panic(s)
}

// CtxPanicf 包含上下文的Panic日志
func (l *Logger) CtxPanicf(ctx context.Context, format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Output(getTracerIDFromCtx(ctx), Lpanic, 2, s)
	panic(s)
}

// CtxPanicln 包含上下文的Panic日志
func (l *Logger) CtxPanicln(ctx context.Context, v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.Output(getTracerIDFromCtx(ctx), Lpanic, 2, s)
	panic(s)
}

// CtxStack 包含runtime的日志
func (l *Logger) CtxStack(ctx context.Context, v ...interface{}) {
	s := fmt.Sprint(v...)
	s += "\n"
	buf := make([]byte, 1024*1024)
	n := runtime.Stack(buf, true)
	s += string(buf[:n])
	s += "\n"
	l.Output(getTracerIDFromCtx(ctx), Lerror, 2, s)
}

// CtxPrint 控制台输出日志
func CtxPrint(ctx context.Context, v ...interface{}) {
	Std.Output(getTracerIDFromCtx(ctx), Linfo, 2, fmt.Sprint(v...))
}

// CtxPrintf 控制台输出日志
func CtxPrintf(ctx context.Context, format string, v ...interface{}) {
	Std.Output(getTracerIDFromCtx(ctx), Linfo, 2, fmt.Sprintf(format, v...))
}

// CtxPrintln 控制台输出日志
func CtxPrintln(ctx context.Context, v ...interface{}) {
	Std.Output(getTracerIDFromCtx(ctx), Linfo, 2, fmt.Sprintln(v...))
}

// CtxDebugf 控制台输出日志
func CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	if Ldebug < Std.Level {
		return
	}
	Std.Output(getTracerIDFromCtx(ctx), Ldebug, 2, fmt.Sprintf(format, v...))
}

// CtxDebug 控制台输出日志
func CtxDebug(ctx context.Context, v ...interface{}) {
	if Ldebug < Std.Level {
		return
	}
	Std.Output(getTracerIDFromCtx(ctx), Ldebug, 2, fmt.Sprintln(v...))
}

// CtxInfof 控制台输出日志
func CtxInfof(ctx context.Context, format string, v ...interface{}) {
	if Linfo < Std.Level {
		return
	}
	Std.Output(getTracerIDFromCtx(ctx), Linfo, 2, fmt.Sprintf(format, v...))
}

// CtxInfo 控制台输出日志
func CtxInfo(ctx context.Context, v ...interface{}) {
	if Linfo < Std.Level {
		return
	}
	Std.Output(getTracerIDFromCtx(ctx), Linfo, 2, fmt.Sprintln(v...))
}

// CtxWarnf 控制台输出日志
func CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	Std.Output(getTracerIDFromCtx(ctx), Lwarn, 2, fmt.Sprintf(format, v...))
}

// CtxWarn 控制台输出日志
func CtxWarn(ctx context.Context, v ...interface{}) {
	Std.Output(getTracerIDFromCtx(ctx), Lwarn, 2, fmt.Sprintln(v...))
}

// CtxErrorf 控制台输出日志
func CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	Std.Output(getTracerIDFromCtx(ctx), Lerror, 2, fmt.Sprintf(format, v...))
}

// CtxError 控制台输出日志
func CtxError(ctx context.Context, v ...interface{}) {
	Std.Output(getTracerIDFromCtx(ctx), Lerror, 2, fmt.Sprintln(v...))
}

// CtxFatal 控制台输出日志
func CtxFatal(ctx context.Context, v ...interface{}) {
	Std.Output(getTracerIDFromCtx(ctx), Lfatal, 2, fmt.Sprint(v...))
	os.Exit(1)
}

// CtxFatalf 控制台输出日志
func CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	Std.Output(getTracerIDFromCtx(ctx), Lfatal, 2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// CtxFatalln 控制台输出日志
func CtxFatalln(ctx context.Context, v ...interface{}) {
	Std.Output(getTracerIDFromCtx(ctx), Lfatal, 2, fmt.Sprintln(v...))
	os.Exit(1)
}

// CtxPanic 控制台输出日志
func CtxPanic(ctx context.Context, v ...interface{}) {
	s := fmt.Sprint(v...)
	Std.Output(getTracerIDFromCtx(ctx), Lpanic, 2, s)
	panic(s)
}

// CtxPanicf 控制台输出日志
func CtxPanicf(ctx context.Context, format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	Std.Output(getTracerIDFromCtx(ctx), Lpanic, 2, s)
	panic(s)
}

// CtxPanicln 控制台输出日志
func CtxPanicln(ctx context.Context, v ...interface{}) {
	s := fmt.Sprintln(v...)
	Std.Output(getTracerIDFromCtx(ctx), Lpanic, 2, s)
	panic(s)
}

// CtxStack 控制台输出日志
func CtxStack(ctx context.Context, v ...interface{}) {
	s := fmt.Sprint(v...)
	s += "\n"
	buf := make([]byte, 1024*1024)
	n := runtime.Stack(buf, true)
	s += string(buf[:n])
	s += "\n"
	Std.Output(getTracerIDFromCtx(ctx), Lerror, 2, s)
}
