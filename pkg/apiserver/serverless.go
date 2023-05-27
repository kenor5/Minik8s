package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"minik8s/entity"
	"minik8s/pkg/apiserver/ControllerManager"
	fc "minik8s/pkg/apiserver/ControllerManager/FunctionController"
	wc "minik8s/pkg/apiserver/ControllerManager/WorkflowController"
	pb "minik8s/pkg/proto"
	"minik8s/tools/log"
	"net/http"

	"strconv"
	"time"
)

// serverless:添加路由和处理函数
func (master *ApiServer)AddRouter(functionName string) error {
	log.PrintS("Add Router:", functionName)
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
			Numstr := strconv.Itoa(podNum)
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
            fc.SetFunction(function)
			podIp = pod.Status.PodIp
            
		} else {
			log.PrintS("podNum != 0, send to pod")
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

func (master *ApiServer)MonitorFunction() {
	for {
        RunningFunctions := fc.GetRunningFunction()
		for _, function := range RunningFunctions {
			// 获取基本信息
			accessTimes := function.FunctionStatus.AccessTimes
            functionPods := function.FunctionStatus.FunctionPods
			podNum := len(functionPods)

			// 根据每30s访问数量计算Pod的期望数量
			targetNum := int((accessTimes + 5 - 1) / 5)
            
			// 调整至期望数量
			// 数量多了则删除Pod
			if (podNum > targetNum) {
				// 首先更新function
				targetFunctionPods := functionPods[:targetNum]
                function.FunctionStatus.FunctionPods = targetFunctionPods
				function.FunctionStatus.AccessTimes = 0
                fc.SetFunction(function)

				// 依次删除多余的Pod
				for i := podNum; i > targetNum; i-- {
					podName := functionPods[i-1].PodName
					pod, _ := ControllerManager.GetPodByName(podName)
					podByte, err := json.Marshal(pod)
					if err != nil {
						fmt.Println("parse pod error")
						return
					}
					in := &pb.DeletePodRequest{
                        Data: podByte,
					}
					master.DeletePod(in)
				}
			} else if (podNum < targetNum) {
                // 数量少了则创建Pod
                // 首先创建Pod
				newPod := function.FunctionStatus.PodTemplate
				newPodName := newPod.Metadata.Name
				for i := podNum; i < targetNum; i++ {
                    newPod.Metadata.Name = newPodName + "-" + strconv.Itoa(i)
					// 组装消息
					podByte, err := json.Marshal(newPod)
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
				}
				// 然后更新function
				PodsList := ControllerManager.GetPodsByLabels(&newPod.Metadata.Labels)
				targetFunctionPods := []entity.FunctionPod{}
				// 遍历列表
				for it := PodsList.Front(); it != nil; it = it.Next() {
					element := it.Value.(*entity.Pod)
                    targetFunctionPods = append(targetFunctionPods, entity.FunctionPod{
						PodName: element.Metadata.Name,
                        PodIp: element.Status.PodIp,
					})
				}
				function.FunctionStatus.FunctionPods = targetFunctionPods
				function.FunctionStatus.AccessTimes = 0
                fc.SetFunction(function)
			} else if (podNum == targetNum) {
				// 否则只清空AccessTime
				function.FunctionStatus.AccessTimes = 0
                fc.SetFunction(function)			
			}
		}

		time.Sleep(30 * time.Second) // 每30秒轮询一次
	}
}	

func (master *ApiServer)AddWorkflowRouter(workflowName string) error {
	log.PrintS("Add Workflow:", workflowName)
	// 处理函数
	workflowHandler := func(w http.ResponseWriter, r *http.Request){
		// 查找workflow
        workflow, _:= wc.GetWorkflow(workflowName)
        
		Next := workflow.StartAt

	    // 创建一个空的JSON对象，将返回值传给下个函数
	    data := new(bytes.Buffer)
		data.ReadFrom(r.Body)
		for Next != "End" {
			log.PrintS("Next Function:", Next)
            workflowNode, _ := wc.GetWorkflowNodeByName(Next, workflow)
		    
			if (workflowNode.Type == "Task") {
                // 如果是Task，转发给对应的函数
                data = fc.SendFunction(Next, data)
				log.PrintS("Response from workflow: ", data)
				Next = workflowNode.Next
				// 判断End是否为True
				if (workflowNode.End == "True") {
					Next = "End"
				}
			} else if (workflowNode.Type == "Choice") {
                // 如果是Choice，进行相应的判断
                Next = wc.SelectChoice(data.String(), workflowNode.Choices)
		    }
	    }
        
		// 返回结果
		fmt.Fprintf(w, "%v\n", data.String())
	}

    master.FunctionManager.Mux.HandleFunc("/workflow/"+workflowName, workflowHandler)

	return nil	
}