package logger

import (
	"fmt"
	"os"
)

func log(color string, a ...any) {
	fmt.Fprint(os.Stderr, "\u001b[38;2;30;144;255m[Tera] ", color)
	for _, el := range a {
		fmt.Fprintf(os.Stderr, "%v ", el)
	}
	fmt.Fprint(os.Stderr, "\u001b[0m")
}

func logf(format string, color string, a ...any) {
	fmt.Fprint(os.Stderr, "\u001b[38;2;30;144;255m[Tera] ", color)
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprint(os.Stderr, "\u001b[0m")
}

func logln(color string, a ...any) {
	fmt.Fprint(os.Stderr, "\u001b[38;2;30;144;255m[Tera] ", color)
	for _, el := range a {
		fmt.Fprintf(os.Stderr, "%v ", el)
	}
	fmt.Fprintln(os.Stderr, "\u001b[0m")
}

func Success(msg ...any) {
	log("\u001b[38;2;50;215;75m", msg...)
}

func Successf(format string, msg ...any) {
	logf(format, "\u001b[38;2;50;215;75m", msg...)
}

func Successln(msg ...any) {
	logln("\u001b[38;2;50;215;75m", msg...)
}

func Error(msg ...any) {
	log("\u001b[38;2;255;69;58m", msg...)
}

func Errorf(format string, msg ...any) {
	logf(format, "\u001b[38;2;255;69;58m", msg...)
}

func Errorln(msg ...any) {
	logln("\u001b[38;2;255;69;58m", msg...)
}

func Warning(msg ...any) {
	log("\u001b[38;2;254;215;9m", msg...)
}

func Warningf(format string, msg ...any) {
	logf(format, "\u001b[38;2;254;215;9m", msg...)
}

func Warningln(msg ...any) {
	logln("\u001b[38;2;254;215;9m", msg...)
}

func Info(msg ...any) {
	log("\u001b[38;2;91;199;245m", msg...)
}

func Infof(format string, msg ...any) {
	logf(format, "\u001b[38;2;91;199;245m", msg...)
}

func Infoln(msg ...any) {
	logln("\u001b[38;2;91;199;245m", msg...)
}
