package functioncontroller

import (
	"bytes"
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

// 给Function发消息
func SendFunction(functionName string, r *http.Request) (string){
		// 发消息
		url := "http://127.0.0.1:8070/function/"+functionName
	
		req, err := http.NewRequest("POST", url, r.Body)
		if err != nil {
			fmt.Println("NewRequest error:", err)
			return "NewRequest error"
		}
	
		req.Header.Set("Content-Type", "application/json")
	
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Request error:", err)
			return "Request error"
		}
		defer resp.Body.Close()
	
		fmt.Println("Response Status:", resp.Status)
		fmt.Println("Response Body:")
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return "Error reading response body"
		}
		fmt.Println(buf.String())	

		return buf.String()
}