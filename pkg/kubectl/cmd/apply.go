package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"minik8s/entity"
	"minik8s/pkg/kubectl/utils"
	"minik8s/tools/log"
	"minik8s/tools/yamlParser"
	"strings"
	"time"

	pb "minik8s/proto"
	"github.com/spf13/cobra"
)

var (
	filename string
	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "kubectl apply is to create object",
		Long: `kubectl apply is to create object
		   for example:
		   kubectl apply -f *.yaml`,
		Run: doApply,
	}
)

func init() {
	applyCmd.Flags().StringVarP(&filename, "filename", "f", "", "yaml name")
	err := applyCmd.MarkFlagRequired("filename")
	if err != nil {
		log.LOG("required filename")
		return
	}
}

func doApply(cmd *cobra.Command, args []string) {
	var (
		dirname                  string
		filenameWithoutExtention string
	)

	b, err := yamlParser.FileExists(filename)
	if !b || err != nil {
		log.LOG("file does not exist")
		return
	}

	// 把路径按照 ‘/’ 拆分开，获取没有 .yaml 后缀的文件名
	arr := strings.Split(filename, "/")
	fmt.Println(arr)
	for i := 0; i < len(arr)-1; i++ {
		dirname = dirname + arr[i] + "/"
	}
	if len(dirname) == 0 {
		dirname = "."
	}
	

	filenameWithoutExtention = strings.Split(arr[len(arr)-1], ".")[0]
	fmt.Println(filenameWithoutExtention)

	obj, err := utils.GetField(dirname, filenameWithoutExtention, "kind")
	fmt.Println(obj)
	if err != nil {
		log.LOG("file has no such field")
	}

	fmt.Println(dirname, filenameWithoutExtention)

	switch obj {
	case "Pod":
	case "pod":
		// 先 parse yaml 文件
		pod := &entity.Pod{}
		_, err := yamlParser.ParseYaml(pod, filename)
		if err != nil {
			log.LOG("parse pod failed")
			return
		}
		fmt.Println(pod)

		// 通过 rpc 连接 apiserver
		cli := NewClient()
		if cli == nil {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 把 pod 序列化成 string 传给 apiserver
		podByte, err := json.Marshal(pod)
		if err != nil {
			fmt.Println("parse pod error")
			return
		}

		res, err := cli.ApplyPod(ctx, &pb.ApplyPodRequest{
			Data: podByte,
		})

		fmt.Println("Create Pod, responce ", res)

	case "Deployment":
	case "deployment":
		deploy := &entity.Deployment{}
		_, err := yamlParser.ParseYaml(deploy, filename)
		if err != nil {
			log.LOG("parse deploy failed")
			return
		}

		// TODO

	case "Service":
	case "service":
		// TODO
	case "Node":
	case "node":
	default:
		log.LOG("there is no object named ")

	}
}
