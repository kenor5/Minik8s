package JobController

import (
	"fmt"
	"bytes"
	"minik8s/tools/log"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// 提交任务
func Sbatch(PodIp string, module_name string) (string, string, error){
	// 获取url
	url := "http://"+ PodIp + ":8090/sbatch"
    fmt.Println(url)
    
	log.PrintS("1")

	// 组装消息
	data := map[string]interface{}{
		"module_name": module_name,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return "", "", err
	}
	
	log.PrintS("2")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("NewRequest error:", err)
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")

	log.PrintS("3")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return "", "", err
	}	

	log.PrintS("4")

	// 读取响应体的内容
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(body)
	if err != nil {
		fmt.Println("ReadAll error:", err)
		return "", "", err
	}
	
	log.PrintS("5")

	// 提取指定字段的原始 JSON 数据
	var result map[string]json.RawMessage
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("JSON unmarshal error:", err)
		return "", "", err
	}

	log.PrintS("6")

	// 从原始 JSON 数据中提取指定字段的值
	var status, info string
	err = json.Unmarshal(result["status"], &status)
	if err != nil {
		fmt.Println("Status unmarshal error:", err)
		return "", "", err
	}
	err = json.Unmarshal(result["job_id"], &info)
	if err != nil {
		fmt.Println("Info unmarshal error:", err)
		return "", "", err
	}

	// 打印提取的字段值
	fmt.Println("Status:", status)
	fmt.Println("job_id:", info)
   
	return status, info, nil
}

// 查询任务结果
func Query(PodIp string, module_name string) (string, string, error){
	// 获取url
	url := "http://"+ PodIp + ":8090/query"
    fmt.Println(url)

	// 组装消息
	data := map[string]interface{}{
		"module_name": module_name,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return "", "",err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("NewRequest error:", err)
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return "", "", err
	}	

	// 读取响应体的内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ReadAll error:", err)
		return "", "", err
	}
	
	// 提取指定字段的原始 JSON 数据
	var result map[string]json.RawMessage
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("JSON unmarshal error:", err)
		return "", "", err
	}

	// 从原始 JSON 数据中提取指定字段的值
	var status, info string
	err = json.Unmarshal(result["status"], &status)
	if err != nil {
		fmt.Println("Status unmarshal error:", err)
		return "", "", err
	}
	err = json.Unmarshal(result["info"], &info)
	if err != nil {
		fmt.Println("Info unmarshal error:", err)
		return "", "", err
	}

	// 打印提取的字段值
	fmt.Println("Status:", status)
	fmt.Println("Info:", info)
   
	return status, info, nil
}