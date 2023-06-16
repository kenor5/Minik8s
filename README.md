# minik8s

## 一、项目总体架构

![image](https://github.com/kenor5/Minik8s/assets/75160010/ad757052-6d8a-4425-a9e9-f234d0fe0e2e)

Minik8s主要包括如下组件：

**Etcd:**

运行在master节点上，记录集群中的所有持久化数据，记录集群中的各种状态 (States)。

**Kubectl：**

运行在master节点，主要负责接受并解析用户指令，发送给Apiserver；获得Apiserver响应后输出到控制台；主要支持apply、delete、get、describe指令，此外还支持add、update指令。

**Apiserver：**

运行在Master节点，是minik8s的中心组件，只有Apiserver可以读写Etcd，集群中的所有通信都依赖于Apiserver暴露的API。

**Scheduler(NodeController):**

运行在Master节点，管理所有的Node，监控Node的状态，完成调度功能。

**ControllerManager：**

  运行在master节点，负责管理对应的Api对象，监控对应的Api对象的状态并且确保其与期望状态一致。包含的Controller如下所示：
  - PodController
  - ServiceController
  - DeploymenyController  
  - FunctionController
  - JobController
  - ScaleController
  - 
**KubeProxy**

运行在每个Node节点上，接受Apiserver创建/删除/更新Service的请求，修改节点上的iptable

**Kubelet**

运行在Node节点，接受Apiserver的请求，管理和创建Pod；监控本地Pod内部运行状态，在Pod状态更新时通知Apiserver更新etcd中Pod状态信息。

## 二、项目开发管理

**gitee仓库目录**

[https://gitee.com/minik8s-group/minik8s.git](https://gitee.com/minik8s-group/minik8s.git)

**项目目录**

![image](https://github.com/kenor5/Minik8s/assets/75160010/181ff64b-77c0-42a1-8086-4e96a1cfca77)

**项目分支介绍**
项目中maser为主分支，每个阶段不同的功能在新的对应分支如serverless、autoscale、Deployment等分支开发。

![image](https://github.com/kenor5/Minik8s/assets/75160010/cbe12fc6-c02b-4275-b0de-e80c5cbe2e43)



## 三、相关依赖和库

**主要使用的库**

* 操作Docker使用的库：
     github.com/docker/docker
* 操作Etcd使用的库： 
   go.etcd.io/etcd/client/v3
* 解析yaml使用的库：
   github.com/spf13/cobra
   github.com/spf13/viper
* 组件通讯的库： 
      google.golang.org/grpc
      google.golang.org/grpc/credentials/insecure
      net/http
* 配置iptable使用的库：
      github.com/coreos/go-iptables/iptables
* 资源监控使用的库：
  github.com/Prometheus
  github.com/docker/go-connections/nat
* 一些其它go基础库：
  "context"
  "encoding/json"
  "fmt"
  "strings"
  "strconv"
  "time"
  "os"
  "sync"
  "regexp"
  "testing"
  
## 四、各功能实现原理

### Node注册、管理和调度策略

[node演示视频](https://www.aliyundrive.com/s/ohUP1FZseWf)

**1.Node的配置文件**

每个节点上的Kubelet在启动时，会读取对应的Node配置文件，如下所示：
```
Name: master
Labels :
  hostName: master
  os: linux
```
除了该Node的名字外，配置文件主要使用Label字段标识节点的一些特点，如os等，这个Label在调度时会用到，如果Pod的NodeSelector字段不为空，则会将该Pod调度到Label符合该Pod要求的Node上。

**2.Node的加入**

当Kubelet启动后，会向ApiServer发消息进行注册，此时运行命令：
```kubectl get node```
可以看到pending状态的Node:

![image](https://github.com/kenor5/Minik8s/assets/75160010/b9a9e011-8593-4556-ba20-0da4534bb6e8)

运行命令：
```
kubectl add node [NodeName]
```
可以将指定的状态为Pending的Node加入集群中(下图一)，如果NodeName为空，则将所有的状态为Pending的Node加入集群中(下图二)。加入集群后，Node的状态变为Running并且NodeController会定期向Node发送心跳监控Node的状态：

![image](https://github.com/kenor5/Minik8s/assets/75160010/fbcca32f-beb8-4ca9-8dc0-756da8049bb1)

![image](https://github.com/kenor5/Minik8s/assets/75160010/3a4f2fad-91c5-497c-9141-8a6b500c5a87)

**3.Node的删除**

运行命令：
```
kubectl delete node [NodeName]
```
会将指定Node删除，将Node的状态变为Pending，并且调度也不会再考虑该Node：

![image](https://github.com/kenor5/Minik8s/assets/75160010/5b098c30-23e0-4069-8f3f-549a05ed0021)


**4.Node的监控**

运行在master节点的NodeController每隔30s会向状态为Running的Node发送心跳，如果没有回应，会将Node的状态标记为dead，不再将Pod调度到该节点上：

![image](https://github.com/kenor5/Minik8s/assets/75160010/b8cfdae1-2ab7-4943-8012-911a8f1f31b0)


**5.调度策略**

NodeController还负责调度功能的实现，我们的Minik8s支持两种调度策略：
- RoundRobin：如果待创建的Pod对于运行在哪个Node没有特殊要求，即NodeSelector字段为空，则会采用RoundRobin的方式调度到状态为Runnning的Node上。
- NodeSelector：如果带创建的Pod的NodeSelector字段不为空，则在调度时会将该字段与Node的Label进行匹配，选出所有符合条件的Node，再使用RoundRobin进行调度。

### Pod创建、管理和内部通信

**演示视频**

[Pod演示视频](https://www.aliyundrive.com/s/gNgs7AamLFk)

**Pod的创建和删除**

Pod的创建分情况处理:
* Kubectl执行apply pod创建，写入Pod元数据信息到etcd，状态标记为Pending，使用对应调度策略，通知Node创建pod。
* 若是动态扩缩或者replicaset创建的新的Pod，则由Apiserver将Pod信息写入etcd，状态标记为Pending，待Apiserver的watch监控Pod机制检测到其状态为Pending，则使用对应调度策略，通知Node创建pod，并更新状态。
在Kubelete中维护一个Podmanager数据结构，存储着Pod的元数据信息和其中的container信息。在Kubelete创建后更新Podmanager中的信息并通知master更新Pod状态为Running，容器命名方式为pod名称与yaml中指定的容器名组合而成:podName_ContainerName。

Pod的删除分情况处理：
* Kubectl直接执行del pod指令，Apiserver更新etcd中Pod信息，通知相应Node的Kubelete删除Pod相关容器；
* 如果是动态扩缩而减少的Pod，Autoscale程序更新etcd中相关Pod信息，标记状态为Succeed，之后Apiserver的watch监控Pod 机制检测到其状态为Succeed，则根据Pod的HostIP信息通知Node实际删除Pod，删除成功后删除etcd中Pod的hostip信息。

**Pod的监控和管理**

* Node上的监控：Kubelete会运行一个监控程序，根据Podmanager中的Pod列表遍历检查当前Pod中的容器是否存活，若有Container容器exited状态，则通知master更新Pod状态信息，删除该Pod剩下的存活Container，重启该Pod，重启成功后通知master更新Pod状态；
* Maseter上的监控:运行一个Pod监控，通过etcd watch监控信息，当由于扩缩策略增加新的Pod时，watch机制检测到变化，其etcd中写入pod状态标记为Pending，需要根据Node策略选择一个Node启动管理该Pod；当watch监测到etcd中Pod标记为Succeed且存在HostIP，则需要根据Hostip通知该Node删除Pod。

**Pod的内部通信**

每个Pod都创建一个pause的容器，pod内其它容器都通过network=pauseContainer的网络模式共享pause容器的网络地址空间，实现Pod内部容器通过localhost互相通信

## Pod间（同Node、不同Node）、service通信

**Pod-Pod**

我们使用Flannel插件并修改Docker的配置完成Pod的全局唯一IP分配以及Pod之间的通信。
每台机器的相关环境如下所示：
```
节点名称 IP地址 安装软件
master 192.168.1.5 etcd、flannel、docker
node1 192.168.1.4 flannel、docker
node2 192.168.1.6 flannel、docker
```

master节点的etcd中存放flannel的配置，并且暴露和监听端口：
```
// 启动Etcd时暴露端口
etcd -name etcd1 -data-dir /root/Tools/etcd/data --advertise-client-urls http://192.168.1.5:2379,http://127.0.0.1:2379 --listen-client-urls http://192.168.1.5:2379,http://127.0.0.1:2379
// 存放flannel配置信息
etcdctl put /coreos.com/network/config '{"Network": "10.0.0.0/16", "SubnetLen": 24, "SubnetMin": "10.0.1.0", "SubnetMax": "10.0.20.0", "Backend": {"Type": "vxlan"}}'
```
node节点的flannel连上master的etcd取得flannel的配置，被分配一个全局不冲突的子网：
```
sudo flanneld --etcd-endpoints=http://192.168.163.132:2379
```
每个节点上使用flannel.1作为docker的启动网桥而不是docker0这个默认的网桥：

![image](https://github.com/kenor5/Minik8s/assets/75160010/50868a7f-e509-4e12-8017-04e9debd2aa2)

这样就实现了全局唯一的IP分配和不同节点上Pod的互相通信。

**Service-Pod**

[service-pod通信演示视频](https://www.aliyundrive.com/s/mYZngsktLax)

Service 启动时，会根据 labels 选择所有符合条件的 Pod，然后修改机器的 iptable，将流量导航到对应的 pod上

```
kind: Service
metadata:
  name: fileserver-service1
spec:
  selector:
    app: myApp
  ports:
  - name: http
    port: 8080
    targetPort: 8080
  
  type: ClusterIP
  clusterIP: "10.20.0.2"
```

具体实现：对于 service yaml 文件中 ports 字段中的每一个 port，在 iptable 中创建一条 svc-chain，每一条 svc-chain 对应若干条 sep-chain，sep-chain再导航到对应的 PodIp 和端口。当一条流量到来，会依次经过这几条链，到达对应的 pod，当然这里的 pod 访问支持是由上面 pod 间通信保障的。

### Deployment的创建和管理

[Deployment演示视频](https://www.aliyundrive.com/s/Aa7JmbnhXqc)

**Deployment创建**

当Apiserver收Kubectl的创建指令，根据yaml中的template 创建Pod元数据写入etcd，etcdwatch监测到etcd中Pod信息更新，就会根据Scheduler策略调度选择Node启动Pod并管理,注意事项为yaml文件中metadata.name命名需要为xxx-deployment，以便生成唯一的Deployment name标识；

**Deployment管理**

master节点DeploymentController运行一个WatchMonitorDeployment 程序，通过Watch监听"Deployment/"前缀的etcd中key,当新建Deploym或者HPA更新了Deployment的信息，监听程序就根据 Spec.replicas和Status.replicas的差异判断需要增加还是删除多余的replica，更新和deployment相关的Pod元数据，通过Pod的watch监听机制间接更新到指定的replica数。

**DeployMent replica更新策略**

Deployment的yaml文件中可以通过spec.update字段指定删除多余replica的策略，当spec.update=newer 或为空时，优先删除后创建的较新的replica；当spec.update=older时优先删除创建时间更早的较旧的replica Pod。

### Autoscale的管理和动态扩缩容策略

[HPA演示视频](https://www.aliyundrive.com/s/WZm6ABwZFSL)

基于cAdvisor和promtheus实现对CPU和Memory资源进行使用情况监控，默认每30s通过master的promtheus查询、计算过去一分钟的CPU或Memory使用情况，监控时间间隔可以通过scaleInterval字段指定，不能低于promtheus采集数据的间隔15s。

执行HPA流程：
* Node节点会在kubelet启动时启动一个cAdvisor容器服务，Master节点会在启动Apiserver后启动一个Promtheus容器服务，当Node注册加入Master时，会更新Promtheus的配置文件将Node的IP和8080端口加入Promtheus监控的一个Target中开启对Node中container的监控。
* 通过yaml配置文件启用扩缩策略，yaml中配置策略主要限制 minReplicas、maxReplicas、metrics中监控对象可以选择CPU或Memory，在behavior字段可以设置扩缩速率的不同策略，详细配置文件见附录中HPA.yaml文件。
* 在每个scaleInterval时间段内，控制器管理器都会根据每个 HorizontalPodAutoscaler 定义中指定的指标查询资源利用率。 选择扩缩策略对应的 Pod， 并从资源指标 API（针对每个 Pod 的资源指标）。对于按 Pod 统计的资源指标（如 CPU），控制器从资源指标 API 中获取每一个 HorizontalPodAutoscaler 指定的 Pod 的度量值，如果设置了目标使用率，控制器获取每个 Pod 中的容器资源使用情况， 并计算资源使用率。如果设置了 target 值，将直接使用原始数据，不再计算百分比。 
* 接下来，控制器根据平均的资源使用率或原始值计算出扩缩的比例，进而计算出目标副本。当目标副本和当前副本数相差倍数较大时会成倍(乘除)扩缩，当数量较接近时会缓慢(加减)扩缩，replica数量增加时不会超过maxReplica，数量减少时不会低于minReplica。

HPA限制扩缩速率策略配置与选择：

```
behavior:
  scaleDown:
    policies:
    - type: Percent
      value: 10
      periodSeconds: 60
    - type: Pods
      value: 5
      periodSeconds: 60
    selectPolicy: Min
```

支持在HPA的yaml通过scaleDown和scaleUp指定不同的扩缩容策略，将上面的 behavior 配置添加到 HPA 中，可限制 Pod 被 HPA 删除速率为每分钟 10%，为了确保每分钟删除的 Pod 数不超过 5 个，可以添加第二个缩容策略，大小固定为 5，并将 `selectPolicy` 设置为最小值`Min`。 将 `selectPolicy `设置为 `Min` 意味着 autoscaler 会选择影响 Pod 数量最小的策略。
若未指定上面的behavior策略则默认策略为：当期望的replica数和当前replica相差多倍时，倍数扩缩到接近预期的replica数，当较接近replica时则使用加减策略微调整replica数量。

### DNS与转发

[DNS演示视频](https://www.aliyundrive.com/s/Js85v73W5Xq)

使用 nginx  容器作为反向代理，当一个请求来到时， 会先在 hosts 文件中解析到 nginx 的IP，通过 nginx 的IP和80端口访问到 nginx 容器，nginx容器会将对应的流量导航到对应的 service ip 上。

```
kind: Dns
metadata:
  name: dns-test1
spec: 
  host: example2.com 
  paths:
  - path: /path1
    serviceName: fileserver-service1
    servicePort: 8080
  - path: /path2
    serviceName: fileserver-service2
    servicePort: 8080
```

**具体实现**

当部署一个 dns 后，会根据 serviceName 在 etcd 中找到对应的 service ip，然后修改 nginx 的配置文件（这个文件是通过挂载卷的方式挂载在容器内的），再执行 nginx reload，这样就可以使 nginx 更新转发目录。

**流程：**

先按照正常情况启动pod和service
```
kubectl apply -f ./test/pod3.yaml
kubectl apply -f ./test/service_test.yaml
```

启动nginx容器
```
./script/dns/dns_start.sh
```

启动dns
```
./cmd/kubectl/kubectl/go apply -f ./test/dns1.yaml
```

这样就能通过域名访问到服务了
```
curl dns-example:80/path1
```

会返回hello world。具体流程为

-> 通过/etc/hosts文件解析dns-example 为127.0.0.1

-> 再通过nginx后台将127.0.0.1:80/a 转发到10.20.0.2:8080 
     这里的10.20.0.2是service 的clusterIP

-> 然后经过iptable 的nat 表的转换，就访问到pod Ip

因此如果在dns的yaml文件中修改了 dns 的server_name， 需要到/etc/hosts文件中 添加一条解析 127.0.0.1 server_name

### 容错和控制面重启

[容错演示视频](https://www.aliyundrive.com/s/WV9oLHZ4JWE)

**master节点容错和重启**

maseter中所有controller管理为无状态，所有元数据信息在apply之前都会写入到etcd，重启后所有信息都可以直接从etcd中恢复，当删除workload时也不会删除etcd中的信息，保证可以恢复到重启前的master控制面状态；

**Node节点容错和重启**

Pod元数据未持久化到硬盘，存储在内存中。主要由Podmanager存储和管理Node上运行的Pod相关的Container信息，在Kubelet重启后注册到maseter和Apiserver建立连接时，maseter会通过Pod中存储的HosIP信息返回给Node和该Node相关的Pod信息，Podmanager会利用该信息重新初始化Podmanager中的内容，建立本地Pod到container的映射关系。

### GPU功能实现

[GPU演示视频](https://www.aliyundrive.com/s/zzTLZf1C6zd)

整个GPU功能的核心是用Python写了一个Flask架构的Http服务器并打包成镜像，该服务器运行在8090端口，提供两个服务：
* 上传并执行任务：访问8090/sbatch，接受一个参数“module_name”，该函数会把运行服务器的容器目录下的“./Data/{module_name}.cu”文件与“./Data/{module_name}.slurm”文件上传至交我算平台，并且提交任务，使用一个全局变量记录job_id，如果任务成功提交，返回成功信息和job_id，否则返回失败信息。
* 查询任务状态：访问8090/query，为了保证任务之间的隔离性，每次提交任务使用不同的Pod，所以全局变量job_id记录的是该pod对应任务的id，故查询不需要参数。该函数会通过记录的job_id查询任务的状态，如果为pending则返回pending；如果正在计算但还没有结果（判断方式为.out和.err文件是否产生）则返回Running；如果计算结束，则根据.err文件是否为空来判断任务成功或失败，并将.out和.err文件拉到本地。
  将该flask程序打包成镜像，每当GPU相关的请求到来时，对应的Job抽象会创建一个使用该镜像容器的Pod，并将yaml文件中的cuda程序和slurm脚本的路径挂载到容器的./Data目录下。
   JobController发现该Pod创建成功后，首先会发送8090:/sbatch，提交任务；然后启动一个go routine，每隔5s向该Pod发送/query请求，并将状态同步到Etcd中。当监控线程发现任务执行完毕后，会将结果拉到本地，并且删除对应的Pod。
   因为这个涉及到与集群外的应用进行交互，有时候因为网络问题，使用交我算的人多等问题，可能出现连不上的情况，因此我们采用了Retry的机制，如果连接交我算/上传文件/提交任务/查询任务的请求失败或超时会重试3次。


### Serverless功能实现

[serverless演示视频](https://www.aliyundrive.com/s/KNc4wN3nxFv)

**镜像的动态创建**

每个函数的镜像都包含一个相同的flask架构的服务器，如下所示：

![image](https://github.com/kenor5/Minik8s/assets/75160010/e9a61515-1baf-4d44-9b24-3c2b04b1b683)

该服务器接收到请求后会动态导入当前目录下名为<module_name>的模块，并且执行该模块中的名为<module_name>的函数，返回该函数的执行结果。
      用户编写好函数以及在requirement.txt中写好自己函数需要pip install的依赖后，在function.yaml中填写路径（functionPath和requirementPath）即可。
      ApiServer在接收到创建函数的请求时，首先会动态的创建镜像，根据创建好的镜像写好Function对象的PodTemplate字段存入Etcd中，并且添加路由（http:浮动IP:8070/function/function_name）。

**HTTP请求的处理**

当用户通过postman等当时向`http://浮动IP:8070/function/function_name` 发送请求时，对应的FunctionHandler首先会判断当前有没有Pod正在运行，如果没有，则进入冷启动流程；否则，将请求负载均衡到正在运行的Pod上。

**Scale-to-0的实现**

每访问一次函数，FunctionController都会记录下来并且持久化到Etcd中。每30s，FunctionController会检查当前函数在这30s内被访问了多少次，并根据访问次数调整运行的Pod的数量。比如，如果用户写明一个Pod为10次/30s，当前有一个Pod正在运行，且在过去30s内有18个请求，则FunctionController就会增加一个Pod以达到用户的要求。同样，如果30s内没有请求到来，则FunctionController会关闭所有的Pod，实现Scale-to-0。

**函数的更新**

修改了functionPath对应的函数内容后，运行kubectl update function [functionName]更新对应的function。
在function更新时，首先会更新镜像，然后判断该function有无正在运行的Pod，如果没有则结束；如果该function有正在运行的Pod，则会删除旧的Pod，并重新创建相同数量的使用新镜像的Pod。

**函数的删除**

首先会判断该函数有没有正在运行的Pod，如果有的话会首先删除该函数所有运行的Pod，然后在Etcd中标记该函数的状态为已删除，当函数删除后再次访问路由则会转入默认的路由，不返回任何值，更不会创建Pod。

**WorkFlow的实现**

总的来说，WorkFlow的实现是建立在Function的实现之上的，核心在于如何处理Choice以及如何衔接Function（即将上一个函数的返回值作为参数传递给下一个参数）。核心代码如下所示：

![image](https://github.com/kenor5/Minik8s/assets/75160010/69080a29-573a-460c-ac24-3ddd07d3cb2f)

使用一个For循环，如果workflowNode的类型为“Task”则调用SendFunction函数交由对应的函数处理，获取JSON形式的返回值，并且判断Workflow是否结束。如果workflowNode的类型为“Choice”，则将上一个函数的返回值data和workflowNode.Choices作为参数传入SelectChoice函数选择下一个节点。
  因为Workflow是基于Function的进一步抽象，所以其扩容以及Scale-to-0与Function相同，不再赘述。

### 附录和补充

**Etcd中数据存储格式**

不同的数据使用不同的Kind前缀加其资源名字作为key存储，value为其结构体中的完整内容序列化转化为字符串，以`Pod/和Deployment/`为例：
```
Pod/nginx-deployment-9594276-320b7
{"kind":"Pod","metadata":{"name":"nginx-deployment-9594276-320b7","labels":{"app":"nginx"},"uid":"93e5d72f-0661-49b9-a722-24f14a2cf41c"},"spec":{"volumes":null,"containers":[{"name":"nginx","image":"nginx:1.14.2","ports":[{"containerPort":"80"}],"resources":{}}]},"status":{"host_ip":"127.0.0.1","phase":"Failed","pod_ip":"172.17.0.2"}}
Deployment/nginx-deployment
{"kind":"Deployment","metadata":{"name":"nginx-deployment","labels":{"app":"nginx"}},"spec":{"selector":{"matchLabels":{"app":"nginx"}},"template":{"metadata":{"name":"","labels":{"app":"nginx"}},"spec":{"volumes":null,"containers":[{"name":"nginx","image":"nginx:1.14.2","ports":[{"containerPort":"80"}],"resources":{}}]}},"replicas":3},"status":{"StartTime":"2023-05-19T19:56:07.506025657-07:00","replicas":0}}
```
