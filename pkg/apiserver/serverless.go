package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"math/rand"
	"minik8s/entity"
	"minik8s/pkg/apiserver/ControllerManager"
	fc "minik8s/pkg/apiserver/ControllerManager/FunctionController"
	pb "minik8s/pkg/proto"
	"minik8s/tools/log"
	"net/http"

	// "net/http/httputil"
	// "net/url"
	"strconv"
	"time"
)

// serverless:添加路由和处理函数
func (master *ApiServer)AddRouter(functionName string) error {
	// 处理函数
	functionHandler := func(w http.ResponseWriter, r *http.Request){
		// 查找function
        function, _:= fc.GetFunction(functionName)
		podNum := len(function.FunctionStatus.FunctionPods)
		podIp := ""
		// 如果没有Pod，则冷启动创建Pod
		if (podNum == 0) {
			log.PrintS("podNum=0, crate pod")
            // 创建Pod
			pod := function.FunctionStatus.PodTemplate
			Numstr := strconv.Itoa(podNum+1)
			pod.Metadata.Name = pod.Metadata.Name + "-" + Numstr
			// 组装消息
			podByte, err := json.Marshal(pod)
			if err != nil {
				log.PrintE("parse pod error")
			}
	        in := &pb.ApplyPodRequest{
		        Data: podByte,
	        }			
			_, err = master.ApplyPod(in)
			if err != nil {
				log.PrintE("Apply pod error")
			}      
			time.Sleep(5 * time.Second) // 休眠 5 秒,等待Pod可用
			
			// 获取PodIp
			podptr, err := ControllerManager.GetPodByName(pod.Metadata.Name)
			if err != nil {
				log.PrintE("Apply pod error")
			}
			pod = *podptr
			
			functionPod := entity.FunctionPod{
				PodName: pod.Metadata.Name,
				PodIp: pod.Status.PodIp,
			}

			// 更新function
			function.FunctionStatus.FunctionPods = append(function.FunctionStatus.FunctionPods, functionPod)
            
			podIp = pod.Status.PodIp

		} else {
			log.PrintS("podNum!=0, send to pod")
			// 否则随机选取一个Pod转发请求
			rand.Seed(time.Now().UnixNano())
			randomNumber := rand.Intn(podNum)
			podIp = function.FunctionStatus.FunctionPods[randomNumber].PodIp
		}
		
		// 发消息
        log.PrintS(podIp)
		url := "http://"+podIp+":8070/function/"+functionName
	
		req, err := http.NewRequest("POST", url, r.Body)
		if err != nil {
			fmt.Println("NewRequest error:", err)
			return
		}
	
		req.Header.Set("Content-Type", "application/json")
	
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Request error:", err)
			return
		}
		defer resp.Body.Close()
	
		fmt.Println("Response Status:", resp.Status)
		fmt.Println("Response Body:")
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		fmt.Println(buf.String())
		fmt.Fprintf(w, "%v\n", buf.String())

		// 记录访问次数
		function.FunctionStatus.AccessTimes += 1
        fc.SetFunction(function)
	}

    master.FunctionManager.Mux.HandleFunc("/function/"+functionName, functionHandler)

	return nil
}
