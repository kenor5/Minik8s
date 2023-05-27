package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v\n", r)
	fmt.Fprintf(w, "Hello, World!")
}

// 默认处理函数
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Default Handler")
}

func main_4() {
	// 创建一个自定义的 ServeMux 对象
	mux := http.NewServeMux()

    // 注册默认路由
	mux.HandleFunc("/", defaultHandler)

	// 启动服务器
	go func() {
		fmt.Println("Server started on port 8080")
		http.ListenAndServe(":8080", mux)
	}()

	try(mux)

	// 阻塞主程序
	select {}
}

func try(mux *http.ServeMux) {
	mux.HandleFunc("/function/hello_function", helloHandler)
}