package servicecontroller

import (
	"encoding/json"
	"minik8s/entity"
	"minik8s/tools/etcdctl"
	"minik8s/tools/log"
)

func GetAllService() ([]*entity.Service, error) {
	// 从etcd中拿出所有的Service
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("etcd client connetc error")
	}
	defer cli.Close()

	out, err := etcdctl.GetWithPrefix(cli, "Service/")
	if err != nil {
		log.PrintE("No such Function!")
		return nil, err
	}
    
	var services []*entity.Service

	// 解析Service并返回
	for _, service := range out.Kvs{
		newService := &entity.Service{}
		err := json.Unmarshal(service.Value, newService)
		if err != nil {
			panic("podNew unmarshel err")
		}
		
		services = append(services, newService)
	}

	return services, nil 
}

func SetService(service *entity.Service) error {
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("etcd client connetc error")
		return err
	}
	defer cli.Close()
	serviceByte, err := json.Marshal(service)
	etcdctl.Put(cli, "Service/"+service.Metadata.Name, string(serviceByte))
	return nil
}