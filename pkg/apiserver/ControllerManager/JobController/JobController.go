package JobController

import (
	"fmt"
	"time"
    "minik8s/entity"
	"minik8s/tools/log"
	"minik8s/tools/etcdctl"
	"encoding/json"
	"minik8s/pkg/apiserver/client"
	"minik8s/pkg/apiserver/ControllerManager"
	pb "minik8s/pkg/proto"
)

func SbatchAndQuery(JobName string, conn pb.KubeletApiServerServiceClient) {
	time.Sleep(2 * time.Second) // 休眠 2 秒,等待Pod可用
	// Job ServerPod的名字
	PodName := JobName + "-ServerPod"
	fmt.Println(PodName)

    // 确认状态,获取IP
	Pod, _ := ControllerManager.GetPodByName(PodName)
	if (Pod.Status.Phase != entity.Running) {
		panic("error")
	}
	PodIp := Pod.Status.PodIp

	// 提交任务
    status, info, err := Sbatch(PodIp, JobName)
    if err != nil {
		log.PrintE("sbatch job error!")
		return
	}
	if status != "Success" {
		log.PrintE("sbatch job error, status not success!")
		return
	}

	// 更新etcd中Job的状态
    job, err := GetJobByName(JobName)
    if err != nil {
		log.PrintE("Get job error!")
		return
	}
    job.JobStatus.Status = entity.Running
	job.JobStatus.JobID = info
	err = PutJobByName(JobName, job)
	if err != nil {
		log.PrintE("Put job error!")
		return
	}

	// 轮询Job的状态
	status = ""
	info = ""
	for {
		status, info, err = Query(PodIp, JobName)
		log.Print(status)

		if (status == "Error" || status == "Success") {
			break
		}

        // 否则状态为Running,继续轮询 
		time.Sleep(5 * time.Second) // 休眠 5 秒
	}
    
	// 不论Job成功失败，删除Pod
	podByte, err := json.Marshal(Pod)
	if err != nil {
		log.PrintE("parse job error")
		return
	}
	in := &pb.DeletePodRequest{
        Data : podByte,
	}
	err = client.KubeletDeletePod(conn, in)

	// 更新etcd中Job的状态
    job, err = GetJobByName(JobName)
    if err != nil {
		log.PrintE("Get job error!")
		return
	}
    job.JobStatus.Status = status
	job.JobStatus.Result = info
	err = PutJobByName(JobName, job)
	if err != nil {
		log.PrintE("Put job error!")
		return
	}    

	fmt.Println("Finished Job")    
}

func GetJobByName(JobName string) (*entity.Job, error) {
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("etcd client connetc error")
		return nil, err 
	}
	defer cli.Close()
	out, err := etcdctl.Get(cli, "Job/"+ JobName)
	if err != nil {
		log.PrintE("No such Job!")
	    return nil, err
	}
	
	// 解析Pod并返回
	job := &entity.Job{}
	json.Unmarshal(out.Kvs[0].Value, job)
	return job, nil
}

func PutJobByName(JobName string, job *entity.Job) error {
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("etcd client connetc error")
		return err 
	}
	defer cli.Close()

	jobByte, err := json.Marshal(job)
	err = etcdctl.Put(cli, "Job/"+JobName, string(jobByte))
	if err != nil {
		log.PrintE("put job err!")
	    return err
	}

	return nil	
}
