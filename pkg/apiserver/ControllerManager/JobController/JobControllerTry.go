// package main
// import(
// 	"fmt"
// 	"minik8s/pkg/apiserver/ControllerManager"
// )

// func main() {
//     pod, err := ControllerManager.GetPodByName("helloworld-ServerPod")
// 	if err != nil {
// 	    panic("error")
// 	}
// 	fmt.Println(pod)
// }

// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"io/ioutil"
// 	"fmt"
// 	"net/http"
// )

// func main() {
// 	url := "http://10.0.4.2:8090/query"
// 	data := map[string]interface{}{
// 		"module_name": "helloworld",
// 	}

// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		fmt.Println("JSON marshal error:", err)
// 		return
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		fmt.Println("NewRequest error:", err)
// 		return
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Request error:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	fmt.Println("Response status:", resp)

// 		// 读取响应体的内容
// 		body, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			fmt.Println("ReadAll error:", err)
// 			return
// 		}
	
// 		// 将响应体的内容转换为字符串并打印出来
// 		responseBody := string(body)
// 		fmt.Println("Response body:", responseBody)
// }

// package main

// import (
// 	"fmt"
// 	"bytes"
// 	"encoding/json"
// 	"io/ioutil"
// 	"net/http"
// )

// func main() {
// 	PodIp := "10.0.4.2"
// 	module_name := "helloworld"

// 	// 获取url
// 	url := "http://"+ PodIp + ":8090/query"
//     fmt.Println(url)

// 	// 组装消息
// 	data := map[string]interface{}{
// 		"module_name": module_name,
// 	}
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		fmt.Println("JSON marshal error:", err)
// 		return
// 	}
	
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		fmt.Println("NewRequest error:", err)
// 		return
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	// 发送请求
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Request error:", err)
// 		return
// 	}	

// 	// 读取响应体的内容
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("ReadAll error:", err)
// 		return
// 	}
	
// 	// 提取指定字段的原始 JSON 数据
// 	var result map[string]json.RawMessage
// 	err = json.Unmarshal(body, &result)
// 	if err != nil {
// 		fmt.Println("JSON unmarshal error:", err)
// 		return
// 	}

// 	// 从原始 JSON 数据中提取指定字段的值
// 	var status, info string
// 	err = json.Unmarshal(result["status"], &status)
// 	if err != nil {
// 		fmt.Println("Status unmarshal error:", err)
// 		return
// 	}
// 	err = json.Unmarshal(result["info"], &info)
// 	if err != nil {
// 		fmt.Println("Info unmarshal error:", err)
// 		return
// 	}

// 	// 打印提取的字段值
// 	fmt.Println("Status:", status)
// 	fmt.Println("Info:", info)
// }

package JobController