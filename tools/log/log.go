package log

import (
	"fmt"
	"runtime"
)

// Print 打印正常内容文字为白色
func Print(a ...any) {
	fmt.Println(a...)
}

// PrintW 打印 warning 文字为黄色加粗
func PrintW(a ...any) {
	pc, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	fmt.Printf("%c[1;1;33m", 0x1B)
	fmt.Printf("[%s]\n%s:%d ", funcName, file, line)
	fmt.Print(a...)
	fmt.Printf("%c[0m\n", 0x1B)
}
func PrintfW(str string, a ...any) {
	pc, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	fmt.Printf("%c[1;1;33m", 0x1B)
	fmt.Printf("[%s]\n%s:%d ", funcName, file, line)
	fmt.Printf(str, a...)
	fmt.Printf("%c[0m\n", 0x1B)
}

// PrintE 打印错误，文字为红色加粗
func PrintE(a ...any) {
	pc, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	fmt.Printf("%c[1;1;31m", 0x1B)
	fmt.Printf("[%s]\n%s:%d ", funcName, file, line)
	fmt.Print(a...)
	fmt.Printf("%c[0m\n", 0x1B)
}

// PrintfE 打印错误，文字为红色加粗
func PrintfE(str string, a ...any) {
	pc, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	fmt.Printf("%c[1;1;31m", 0x1B)
	fmt.Printf("[%s]\n%s:%d ", funcName, file, line)
	fmt.Printf(str, a...)
	fmt.Printf("%c[0m\n", 0x1B)
}

// PrintS 打印成功信息，文字为绿色加粗
func PrintS(a ...any) {
	fmt.Printf("%c[1;1;32m", 0x1B)
	fmt.Print(a...)
	fmt.Printf("%c[0m\n", 0x1B)
}

func Printf(str string, a ...any) {
	pc, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	fmt.Printf("%c[1;1;32m", 0x1B)
	fmt.Printf("[%s]\n%s:%d ", funcName, file, line)
	fmt.Printf(str, a...)
	fmt.Printf("%c[0m\n", 0x1B)
}
