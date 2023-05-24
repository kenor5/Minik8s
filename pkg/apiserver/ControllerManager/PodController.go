package ControllerManager

import (
	"container/list"
	"encoding/json"
	"fmt"
	"minik8s/entity"
	"minik8s/tools/etcdctl"
	"minik8s/tools/log"
	"reflect"
)

func GetPodsByLabels(labels *map[string]string) *list.List {
	selectedPods := list.New()

	// 从etcd中拿出所有的Pod
	cli, err := etcdctl.NewClient()
	if err != nil {
		fmt.Println("etcd client connetc error")
	}
	defer cli.Close()

	out, _ := etcdctl.GetWithPrefix(cli, "Pod/")

	// 筛选出符合Label且状态为Running的Pod
	for _, data := range out.Kvs {
		pod := &entity.Pod{}
		err = json.Unmarshal(data.Value, pod)
		if err != nil {
			fmt.Println("pod unmarshal error")
		}
		fmt.Println("get etcd", pod)
		fmt.Println("lable1 ")
		fmt.Println(pod.Metadata.Labels)
		fmt.Println("label ")
		fmt.Println(*labels)
		fmt.Println(reflect.DeepEqual(pod.Metadata.Labels, *labels))
        // 判断Pod仍在运行(状态为Running)Selector和Label完全相等
		if pod.Status.Phase == entity.Running && reflect.DeepEqual(pod.Metadata.Labels, *labels){
			selectedPods.PushBack(pod)
		}
	}

	return selectedPods
}

func GetPodByName(PodName string) (*entity.Pod, error) {
		// 从etcd中拿出Pod
		cli, err := etcdctl.NewClient()
		if err != nil {
			fmt.Println("etcd client connetc error")
			return nil, err
		}
		defer cli.Close()
		out, err := etcdctl.Get(cli, "Pod/"+PodName)
		if err != nil {
			fmt.Println("No such Pod!")
		    return nil, err
		}

		// 解析Pod并返回
		pod := &entity.Pod{}
        json.Unmarshal(out.Kvs[0].Value, pod)
		return pod, nil
} 


// for debug
func PrintList(List *list.List) {
	if List.Len() == 0 {
		log.PrintE("list len is 0")
	} 

	for element := List.Front(); element != nil; element = element.Next() {
		value := element.Value
		fmt.Println(value)
	}
}
