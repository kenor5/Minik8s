package kubelet

import (
	"fmt"
	"minik8s/configs"
	"minik8s/entity"
	"minik8s/tools/log"
	"os"
	"os/exec"
)


func CreateDns(dns *entity.Dns) error {

	// 修改nginx.conf
	err := modifyNginxConf(dns)
	if err != nil {
		log.PrintE("modify nginx.conf error")
	}

	// 重启nginx,使得nginx.conf生效
	cmd := exec.Command("/usr/bin/docker", "exec", configs.NginxContainerName, "nginx", "-s", "reload")
	err = cmd.Run()
	if err != nil {
		log.PrintE("nginx reload error")
		log.PrintE(err)
	}
	return nil
}

func modifyNginxConf(dns *entity.Dns) error {
	string2write := "server {\n\tlisten 80;\n\tserver_name " + dns.Metadata.Name + ";\n"

	for _, ip := range dns.Spec.Paths { 
		// **注意 在service ip：port 最后要加上 一个正斜杠 ‘/’ 这样会让匹配的subpath不出现在service ip后面 **
		// 例如 访问 http://example.com:80/aa时，如果不加正斜杠，会被转发到 http://10.20.0.2:8080/aa
		// 加了正斜杠 会被转发到 http://10.20.0.2:8080。 前者会导致错误
		string2write += "\tlocation " + ip.Path + " {\n\t\tproxy_pass http://" + ip.ServiceName + ":" + fmt.Sprintf("%d", ip.ServicePort) + "/;\n\t}\n"
	}
	string2write += "}\n"

	// 修改nginx.conf
	file, err := os.OpenFile(configs.NginxConfPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.PrintE("open nginx.conf error")
	}

	_, err = file.WriteString(string2write)
	if err != nil {
		log.PrintE("write nginx.conf error")
	}
	file.Close()
	
	return nil
}
