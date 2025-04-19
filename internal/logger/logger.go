package logger

import (
	"fmt"
	"os"
)

func Success(msg ...any) {
	fmt.Fprintln(os.Stderr, "\u001b[38;2;30;144;255m[Tera]\u001b[38;2;50;215;75m", fmt.Sprint(msg...))
}

func Error(msg ...any) {
	fmt.Fprintln(os.Stderr, "\u001b[38;2;30;144;255m[Tera]\u001b[38;2;255;69;58m", fmt.Sprint(msg...))
}

func Warning(msg ...any) {
	fmt.Fprintln(os.Stderr, "\u001b[38;2;30;144;255m[Tera]\u001b[38;2;254;215;9m", fmt.Sprint(msg...))
}

func Info(msg ...any) {
	fmt.Fprintln(os.Stderr, "\u001b[38;2;30;144;255m[Tera]\u001b[38;2;91;199;245m", fmt.Sprint(msg...))
}
