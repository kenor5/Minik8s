package ControllerManager

import (
	"container/list"
	"encoding/json"
	"fmt"
	"minik8s/entity"
	"minik8s/tools/etcdctl"
	"reflect"
)

func GetPodsByLabels(labels *map[string]string) *list.List {
	selectedPods := list.New()

	// 从etcd中拿出所有的Pod
	cli, err := etcdctl.NewClient()
	if err != nil {
		fmt.Println("etcd client connetc error")
	}

	out, _ := etcdctl.GetWithPrefix(cli, "Pod/")

	// 筛选出符合Label且状态为Running的Pod
	for _, data := range out.Kvs {
		pod := &entity.Pod{}
		err = json.Unmarshal(data.Value, pod)
		if err != nil {
			fmt.Println("pod unmarshal error")
		}
		//fmt.Println("get etcd", pod)

        // 判断Pod仍在运行(状态为Running)Selector和Label完全相等
		if pod.Status.Phase == entity.Running && reflect.DeepEqual(pod.Metadata.Labels, *labels){
			selectedPods.PushBack(pod)
		}
	}

	return selectedPods
}

// for debug
func PrintList(List *list.List) {
	for element := List.Front(); element != nil; element = element.Next() {
		value := element.Value
		fmt.Println(value)
	}
}
