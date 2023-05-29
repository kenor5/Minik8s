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

	pb "minik8s/pkg/proto"

	"github.com/spf13/cobra"
	// "google.golang.org/grpc/channelz/service"
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
		fmt.Println("required filename")
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
		fmt.Println("file does not exist")
		return
	}

	// 把路径按照 ‘/’ 拆分开，获取没有 .yaml 后缀的文件名
	arr := strings.Split(filename, "/")
	for i := 0; i < len(arr)-1; i++ {
		dirname = dirname + arr[i] + "/"
	}
	if len(dirname) == 0 {
		dirname = "."
	}

	filenameWithoutExtention = strings.Split(arr[len(arr)-1], ".")[0]

	obj, err := utils.GetField(dirname, filenameWithoutExtention, "kind")
	if err != nil {
		fmt.Println("file has no such field")
	}

	fmt.Println("dirname: ", dirname, "filename without extention: ", filenameWithoutExtention)

	switch obj {
	case "Pod", "pod":
		err := applyPod(filename)
		if err != nil {
			fmt.Println(err)
		}

	case "Deployment", "deployment":
		err := applyDeployment(filename)
		if err != nil {
			fmt.Println(err)
			return
		}

	case "Service", "service":
		err := applyService(filename)
		if err != nil {
			fmt.Println(err)
		}

	case "Dns", "dns":
		err := applyDns(filename)
		if err != nil {
			log.PrintE(err)
		}
	case "Job", "job":
		err := applyJob(filename)
		if err != nil {
			return
		}
		if err != nil {
			fmt.Println(err)
		}

	case "Function", "function":
		err := applyFunction(filename)
		if err != nil {
			return
		}
		if err != nil {
			fmt.Println(err)
		}
	case "HPA", "hpa", "HorizontalPodAutoscaler":
		err := applyHPA(filename)
		if err != nil {
			return
		}

	case "Workflow", "workflow":
		err := applyWorkflow(filename)
		if err != nil {
			return
		}

	default:
		log.PrintE("there is no object named " + obj)
	}
}

func applyPod(filename string) error {
	// 先 parse yaml 文件
	pod := &entity.Pod{}
	_, err := yamlParser.ParseYaml(pod, filename)
	if err != nil {
		fmt.Println("parse pod failed")
		return err
	}
	fmt.Println(pod)

	// 通过 rpc 连接 apiserver
	cli := NewClient()
	if cli == nil {
		return fmt.Errorf("fail to connect to apiserver")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 把 pod 序列化成 string 传给 apiserver
	podByte, err := json.Marshal(pod)
	if err != nil {
		fmt.Println("parse pod error")
		return err
	}

	res, err := cli.ApplyPod(ctx, &pb.ApplyPodRequest{
		Data: podByte,
	})

	fmt.Println("Create Pod, response ", res)
	return nil
}

func applyDeployment(filename string) error {
	deployment := &entity.Deployment{}
	_, err := yamlParser.ParseYaml(deployment, filename)
	if err != nil {
		fmt.Println("parse deployment failed")
		return err
	}

	cli := NewClient()
	if cli == nil {
		return fmt.Errorf("fail to connect to apiserver")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 把 deployment 序列化成 string 传给 apiserver
	deploymentByte, err := json.Marshal(deployment)
	if err != nil {
		fmt.Println("parse deployment error")
		return err
	}

	res, err := cli.ApplyDeployment(ctx, &pb.ApplyDeploymentRequest{
		Data: deploymentByte,
	})

	fmt.Printf("Create Deployment, response %v,error %v\n", res, err)
	return nil
}

func applyService(filename string) error {
	service := &entity.Service{}
	_, err := yamlParser.ParseYaml(service, filename)
	if err != nil {
		fmt.Println("parse service failed")
		return err
	}

	cli := NewClient()
	if cli == nil {
		return fmt.Errorf("fail to connect to apiserver")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 把 pod 序列化成 string 传给 apiserver
	podByte, err := json.Marshal(service)
	if err != nil {
		fmt.Println("parse service error")
		return err
	}

	res, err := cli.ApplyService(ctx, &pb.ApplyServiceRequest{
		Data: podByte,
	})

	fmt.Println("Create Service, response ", res)
	return nil
}

func applyHPA(filename string) error {
	hpa := &entity.HorizontalPodAutoscaler{}
	_, err := yamlParser.ParseYaml(hpa, filename)
	if err != nil {
		fmt.Println("parse hpa failed")
		return err
	}

	cli := NewClient()
	if cli == nil {
		return fmt.Errorf("fail to connect to apiserver")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 把 deployment 序列化成 string 传给 apiserver
	hpaByte, err := json.Marshal(hpa)
	if err != nil {
		fmt.Println("parse hpa error")
		return err
	}

	res, err := cli.ApplyHPA(ctx, &pb.ApplyHorizontalPodAutoscalerRequest{
		Data: hpaByte,
	})

	fmt.Printf("Create HPA for %s, response %v,error %v\n", hpa.Spec.ScaleTargetRef.Name, res, err)
	return nil
}

func applyDns(filename string) error {
	dns := &entity.Dns{}
	_, err := yamlParser.ParseYaml(dns, filename)
	if err != nil {
		fmt.Println("parse dns failed")
		return err
	}

	cli := NewClient()
	if cli == nil {
		return fmt.Errorf("fail to connect to apiserver")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dnsByte, err := json.Marshal(dns)
	if err != nil {
		fmt.Println("parse dns error")
		return err
	}

	res, err := cli.ApplyDns(ctx, &pb.ApplyDnsRequest{
		Data: dnsByte,
	})

	fmt.Println("Create Dns, response ", res)
	return nil
}

func applyJob(filename string) error {
	// 先 parse yaml 文件
	job := &entity.Job{}
	_, err := yamlParser.ParseYaml(job, filename)
	if err != nil {
		fmt.Println("parse pod failed")
		return err
	}
	fmt.Println(job)

	// 通过 rpc 连接 apiserver
	cli := NewClient()
	if cli == nil {
		return fmt.Errorf("fail to connect to apiserver")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 把 pod 序列化成 string 传给 apiserver
	jobByte, err := json.Marshal(job)
	if err != nil {
		fmt.Println("parse job error")
		return err
	}

	res, err := cli.ApplyJob(ctx, &pb.ApplyJobRequest{
		Data: jobByte,
	})

	fmt.Println("Create Job, response ", res)
	return nil
}

func applyFunction(filename string) error {
	// 先 parse yaml 文件
	function := &entity.Function{}
	_, err := yamlParser.ParseYaml(function, filename)
	if err != nil {
		fmt.Println("parse pod failed")
		return err
	}
	fmt.Println(function)

	// 通过 rpc 连接 apiserver
	cli := NewClient()
	if cli == nil {
		return fmt.Errorf("fail to connect to apiserver")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 把 pod 序列化成 string 传给 apiserver
	functionByte, err := json.Marshal(function)
	if err != nil {
		fmt.Println("parse job error")
		return err
	}

	res, err := cli.ApplyFunction(ctx, &pb.ApplyFunctionRequest{
		Data: functionByte,
	})

	fmt.Println("Create Job, response ", res)
	return nil
}

func applyWorkflow(filename string) error {
	// 先 parse yaml 文件
	workflow := &entity.Workflow{}
	_, err := yamlParser.ParseYaml(workflow, filename)
	if err != nil {
		fmt.Println("parse workflow failed")
		return err
	}
	fmt.Println(workflow)

	// 通过 rpc 连接 apiserver
	cli := NewClient()
	if cli == nil {
		return fmt.Errorf("fail to connect to apiserver")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 把 pod 序列化成 string 传给 apiserver
	workflowByte, err := json.Marshal(workflow)
	if err != nil {
		fmt.Println("parse workflow error")
		return err
	}

	res, err := cli.ApplyWorkflow(ctx, &pb.ApplyWorkflowRequest{
		Data: workflowByte,
	})

	fmt.Println("Create Workflow, response ", res)
	return nil
}
