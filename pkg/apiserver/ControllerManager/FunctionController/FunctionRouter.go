package functioncontroller

import (
	"fmt"
	"net/http"
	// "minik8s/pkg/apiserver/ControllerManager"
)

/****************************************
**             主数据结构              ***
*****************************************/
type FunctionController struct {
	Mux *http.ServeMux // 创建一个自定义的 ServeMux 对象
} 

func NewFunctionController() *FunctionController {
	newFunctionController := &FunctionController{
        Mux : http.NewServeMux(),
	}
	return newFunctionController
}

/*******************************************************
** Serverless服务器，暴露端口让用户可以通过路由访问函数 **
********************************************************/
func (fc *FunctionController)FunctionServer() {
	// 注册默认路由
	fc.Mux.HandleFunc("/", defaultHandler)

	// 启动服务器
	fmt.Println("[Serverless] Server started on port 8070")
	http.ListenAndServe(":8070", fc.Mux)
}

// 默认处理函数
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Default Handler")
}

