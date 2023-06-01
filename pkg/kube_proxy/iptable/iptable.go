package iptable

import (
	"fmt"
	"github.com/coreos/go-iptables/iptables"
	"github.com/google/uuid"
)

const (
	chainPreRouting  string = "PREROUTING"
	chainOUTPUT      string = "OUTPUT"
	chainPostRouting string = "POSTROUTING"
	chainMasquerade  string = "MASQUERADE"

	chainService         string = "KUBE-SERVICES"
	chainKubePostRouting string = "KUBE-POSTROUTING"
	chainKubeMarkMasq    string = "KUBE-MARK-MASQ"
)

type IpTable struct {
	hostIp	string
	flannelIp	string
	iptables *iptables.IPTables
}

func NewClient() (*IpTable, error) {
	fmt.Println("creating new iptable")
	tables, err := iptables.New()
	if err != nil {
		return nil, err
	}

	return &IpTable{iptables: tables}, nil
}

/*
reference to :
	https://juejin.cn/post/7134143215380201479

这是 k8s 官方的 iptable 规则，我们仿照它创建对应的chain，不需要的规则以 # 开头
#-N KUBE-MARK-DROP
-N KUBE-MARK-MASQ
#-N KUBE-NODEPORTS
-N KUBE-POSTROUTING
-N KUBE-SEP-FNO4E6JYD7EGUHTP
-N KUBE-SEP-SMDVMBZNJPO5AA7R
-N KUBE-SERVICES
-N KUBE-SVC-ELCM5PCEQWBTUJ2I
-A PREROUTING -m comment --comment "kubernetes service portals" -j KUBE-SERVICES
-A OUTPUT -m comment --comment "kubernetes service portals" -j KUBE-SERVICES
-A POSTROUTING -m comment --comment "kubernetes postrouting rules" -j KUBE-POSTROUTING
#-A KUBE-MARK-DROP -j MARK --set-xmark 0x8000/0x8000
-A KUBE-MARK-MASQ -j MARK --set-xmark 0x4000/0x4000
-A KUBE-POSTROUTING -m comment --comment "kubernetes service traffic requiring SNAT" -m mark --mark 0x4000/0x4000 -j MASQUERADE
-A KUBE-SEP-FNO4E6JYD7EGUHTP -s 10.244.1.7/32 -j KUBE-MARK-MASQ
-A KUBE-SEP-FNO4E6JYD7EGUHTP -p tcp -m tcp -j DNAT --to-destination 10.244.1.7:80
-A KUBE-SEP-SMDVMBZNJPO5AA7R -s 10.244.0.3/32 -j KUBE-MARK-MASQ
-A KUBE-SEP-SMDVMBZNJPO5AA7R -p tcp -m tcp -j DNAT --to-destination 10.244.0.3:80
-A KUBE-SERVICES ! -s 10.244.0.0/16 -d 10.1.190.219/32 -p tcp -m comment --comment "default/nginx-svc:http cluster IP" -m tcp --dport 8080 -j KUBE-MARK-MASQ
-A KUBE-SERVICES -d 10.1.190.219/32 -p tcp -m comment --comment "default/nginx-svc:http cluster IP" -m tcp --dport 8080 -j KUBE-SVC-ELCM5PCEQWBTUJ2I
-A KUBE-SVC-ELCM5PCEQWBTUJ2I -m statistic --mode random --probability 0.50000000000 -j KUBE-SEP-SMDVMBZNJPO5AA7R
-A KUBE-SVC-ELCM5PCEQWBTUJ2I -j KUBE-SEP-FNO4E6JYD7EGUHTP
*/

func (cli *IpTable) Init() error {
	// 如果有这样的规则，则先删除,防止它在其它规则后面
	// err := cli.iptables.DeleteIfExists("nat", chainPreRouting, "-j", chainService)
	// if err != nil {
	// 	return err
	// }

	// err = cli.iptables.DeleteIfExists("nat", chainOUTPUT, "-j", chainService)
	// if err != nil {
	// 	return err
	// }

	// err = cli.iptables.DeleteIfExists("nat", chainPostRouting,  "-j", chainKubePostRouting)
	// if err != nil {
	// 	return err
	// }

	// err = cli.iptables.DeleteIfExists("nat",
	// 	chainKubePostRouting,
	// 	"-m",
	// 	"mark",
	// 	"--mark",
	// 	"0x4000/0x4000",
	// 	"-j",
	// 	chainMasquerade)
	// if err != nil {
	// 	return err
	// }

	// err := cli.iptables.DeleteIfExists("nat",
	// 	chainKubeMarkMasq,
	// 	"-j",
	// 	"MARK",
	// 	"--set-xmark",
	// 	"0x4000/0x4000")
	// if err != nil {
	// 	return err
	// }

	
	// 创建新链和新规则
	
	// -N KUBE-MARK-MASQ
	exists, err := cli.iptables.ChainExists("nat", chainKubeMarkMasq)
	if err != nil {
		return err
	}
	if !exists {
		err := cli.iptables.NewChain("nat", chainKubeMarkMasq)
		if err != nil {
			return err
		}
	}

	exists, err = cli.iptables.ChainExists("nat", chainService)
	if err != nil {
		return err
	}
	if !exists {
		//-N KUBE-SERVICES
		err = cli.iptables.NewChain("nat", chainService)
		if err != nil {
			return err
		}
	}

	exists, err = cli.iptables.ChainExists("nat", chainKubePostRouting)
	if err != nil {
		return err
	}
	if !exists {
		//-N KUBE-POSTROUTING
		err = cli.iptables.NewChain("nat", chainKubePostRouting)
		if err != nil {
			return err
		}
	}

	exists, err = cli.iptables.Exists("nat", chainPreRouting, "-j", chainService)
	if err != nil {
		return err
	}
	if !exists {
		// -A PREROUTING  -j KUBE-SERVICES
		err = cli.iptables.Insert("nat", chainPreRouting, 1, "-j", chainService)
		if err != nil {
			return err
		}
	}

	exists, err = cli.iptables.Exists("nat", chainOUTPUT, "-j", chainService)
	if err != nil {
		return err
	}
	if !exists {
		// -A OUTPUT  -j KUBE-SERVICES
		err = cli.iptables.Insert("nat", chainOUTPUT, 1, "-j", chainService)
		if err != nil {
			return err
		}
	}

	exists, err = cli.iptables.Exists("nat", chainPostRouting, "-j", chainKubePostRouting)
	if err != nil {
		return err
	}
	if !exists {
		// -A POSTROUTING  -j KUBE-POSTROUTING
		err = cli.iptables.Insert("nat", chainPostRouting, 1, "-j", chainKubePostRouting)
		if err != nil {
			return err
		}
	}

	exists, err = cli.iptables.Exists("nat",
		chainKubeMarkMasq,
		"-j",
		"MARK",
		"--set-xmark",
		"0x4000/0x4000")
	if err != nil {
		return err
	}
	if !exists {
		// -A KUBE-MARK-MASQ -j MARK --set-xmark 0x4000/0x4000
		err = cli.iptables.Insert("nat",
			chainKubeMarkMasq,
			1,
			"-j",
			"MARK",
			"--set-xmark",
			"0x4000/0x4000")
		if err != nil {
			return err
		}
	}

	// 这里是为了把外部来的包的source换掉
	
	// exists, err = cli.iptables.Exists("nat",
	// 	chainKubePostRouting,
	// 	"-j",
	// 	chainMasquerade,
	// 	"--random-fully")
	// if err != nil {
	// 	return err
	// }
	// if !exists {
	// 	// -A KUBE-POSTROUTING -j MASQUERADE
	// 	err = cli.iptables.Insert("nat",
	// 		chainKubePostRouting,
	// 		1,
	// 		"-j",
	// 		chainMasquerade,
	// 		"--random-fully")
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// exists, err = cli.iptables.Exists("nat",
	// chainKubePostRouting,
	// "-j",
	// "MARK",
	// "--set-xmark",
	// "0x4000/0")
	// if err != nil {
	// 	return err
	// }
	// if !exists {
	// 	// -A KUBERBOAT-POSTROUTING -j MARK --set-xmark 0x4000/0
	// 	err = cli.iptables.Insert("nat",
	// 	chainKubePostRouting,
	// 	1,
	// 	"-j",
	// 	"MARK",
	// 	"--set-xmark",
	// 	"0x4000/0")
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	exists, err = cli.iptables.Exists("nat",
		chainKubePostRouting,
		"-m",
		"mark",
		
		"--mark",
		"0x4000/0x4000",
		"-j",
		"MASQUERADE")
	if err != nil {
		return err
	}
	if !exists {
		// -A KUBE-POSTROUTING -m mark ! --mark 0x4000/0x4000 -j RETURN
		err = cli.iptables.Insert("nat",
			chainKubePostRouting,
			1,
			"-m",
			"mark",
			
			"--mark",
			"0x4000/0x4000",
			"-j",
			"MASQUERADE")
		if err != nil {
			return err
		}
	}
	
	return nil

}

func (cli *IpTable) CreateServiceChain() string {
	// -N KUBE-SVC-ELCM5PCEQWBTUJ2I
	chainName := "KUBE-SVC-" + uuid.NewString()[:8]
	err := cli.iptables.NewChain("nat", chainName)
	if err != nil {
		return ""
	} else {
		return chainName
	}
}

func (cli *IpTable) CreatePodChain() string {
	// -N KUBE-SEP-FNO4E6JYD7EGUHTP
	chainName := "KUBE-SEP-" + uuid.NewString()[:8]
	err := cli.iptables.NewChain("nat", chainName)
	if err != nil {
		return ""
	} else {
		return chainName
	}
}

// AddServiceRules serivce chain -> specific service chain
func (cli *IpTable) AddServiceRules(clusterIp string, serviceChainName string, port uint32) error {
	//-A KUBE-SERVICES ! -s 10.244.0.0/16 -d 10.1.190.219/32 -p tcp  -m tcp --dport 8080 -j KUBE-MARK-MASQ
	//-A KUBE-SERVICES -d 10.1.190.219/32 -p tcp -m tcp --dport 8080 -j KUBE-SVC-ELCM5PCEQWBTUJ2
	

	err := cli.iptables.AppendUnique(
		"nat",
		chainService,
		"-p",
		"tcp",
		"-d",
		clusterIp,
		"-m",
		"tcp",
		"--dport",
		fmt.Sprint(port),
		"-j",
		"KUBE-MARK-MASQ",
	)
	if err != nil {
		return err
	}

	err = cli.iptables.AppendUnique(
		"nat",
		chainService,
		"-p",
		"tcp",
		"-d",
		clusterIp,
		"-m",
		"tcp",
		"--dport",
		fmt.Sprint(port),
		"-j",
		serviceChainName,
	)
	if err != nil {
		return err
	}


	return nil
}

func (cli *IpTable) RemoveServiceChain(
	serviceChainName string,
	clusterIp string,
	port uint32,
) error {

	err := cli.iptables.DeleteIfExists(
		"nat",
		chainService,
		"-p",
		"tcp",
		"-d",
		clusterIp,
		"-m",
		"tcp",
		"--dport",
		fmt.Sprint(port),
		"-j",
		serviceChainName,
	)
	if err != nil {
		return err
	}

	err = cli.iptables.ClearAndDeleteChain("nat", serviceChainName)
	if err != nil {
		return err
	}
	return nil
}

// AddPodRules sepcific pod chain -> exact pod
func (cli *IpTable) AddPodRules(
	podChainName string,
	podIp string,
	targetPort uint32,
) error {

	//-A KUBE-SEP-FNO4E6JYD7EGUHTP -s 10.244.1.7/32 -j KUBE-MARK-MASQ
	//-A KUBE-SEP-FNO4E6JYD7EGUHTP -p tcp -m tcp -j DNAT --to-destination 10.244.1.7:80
	//-A KUBE-SEP-SMDVMBZNJPO5AA7R -s 10.244.0.3/32 -j KUBE-MARK-MASQ
	//-A KUBE-SEP-SMDVMBZNJPO5AA7R -p tcp -m tcp -j DNAT --to-destination 10.244.0.3:80

	err := cli.iptables.AppendUnique(
		"nat",
		podChainName,
		"-s",
		podIp,
		"-j",
		chainKubeMarkMasq,
	)
	if err != nil {
		return err
	}

	err = cli.iptables.AppendUnique(
		"nat",
		podChainName,
		"-p",
		"tcp",
		"-j",
		"DNAT",
		"--to-destination",
		fmt.Sprintf("%s:%d", podIp, targetPort),
	)
	if err != nil {
		return err
	}

	return nil
}

func (cli *IpTable) RemovePodChain(podChainName string) error {
	err := cli.iptables.ClearAndDeleteChain("nat", podChainName)
	if err != nil {
		return err
	}
	return nil
}

// ApplyPodRules specific service chain -> specific pod chain
func (cli *IpTable) ApplyPodRules(
	serviceChainName string,
	podChainName string,
	nth int,
) error {
	// -A KUBE-SVC-ELCM5PCEQWBTUJ2I -m statistic --mode random --probability 0.50000000000 -j KUBE-SEP-SMDVMBZNJPO5AA7R
	// -A KUBE-SVC-ELCM5PCEQWBTUJ2I -j KUBE-SEP-FNO4E6JYD7EGUHTP

	// 这里设置每n个包匹配一个
	err := cli.iptables.AppendUnique(
		"nat",
		serviceChainName,
		"-m",
		"statistic",
		"--mode",
		"nth",
		"--every",
		fmt.Sprint(nth),
		"--packet",
		"0",
		"-j",
		podChainName,
	)
	if err != nil {
		return err
	}

	return nil
}

func (cli *IpTable) RemovePodRules(
	serviceChainName string,
	podChainName string,
	nth int,
) error {
	// -A KUBE-SVC-ELCM5PCEQWBTUJ2I -m statistic --mode random --probability 0.50000000000 -j KUBE-SEP-SMDVMBZNJPO5AA7R
	// -A KUBE-SVC-ELCM5PCEQWBTUJ2I -j KUBE-SEP-FNO4E6JYD7EGUHTP

	// 这里设置每n个包匹配一个
	err := cli.iptables.DeleteIfExists(
		"nat",
		serviceChainName,
		"-m",
		"statistic",
		"--mode",
		"nth",
		"--every",
		fmt.Sprint(nth),
		"--packet",
		
		"0",
		"-j",
		podChainName,
	)
	if err != nil {
		return err
		// return nil
	}

	return nil
}