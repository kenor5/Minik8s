package functioncontroller

import (
	"encoding/json"
	"minik8s/entity"
	"minik8s/tools/etcdctl"
	"minik8s/tools/log"
)

func GetFunction(functionName string) (*entity.Function, error) {
	// 从etcd中拿出对应的function
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("etcd client connetc error")
	}
	defer cli.Close()

	out, err := etcdctl.Get(cli, "Function/"+functionName)
	if err != nil {
		log.PrintE("No such Function!")
		return nil, err
	}

	// 解析Function并返回
	function := &entity.Function{}
    json.Unmarshal(out.Kvs[0].Value, function)
	return function, nil
}

func SetFunction(function *entity.Function) error {
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("etcd client connetc error")
		return err
	}	
	functionByte, err := json.Marshal(function)
	etcdctl.Put(cli, "Function/"+function.Metadata.Name, string(functionByte))
    return nil
}