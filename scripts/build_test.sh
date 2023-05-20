#!/bin/bash

# 定义变量
srv_bin_dir="${GOPATH}/src/minik8s/bin"
apiserver="${GOPATH}/src/minik8s/cmd/kube-apiserver/apiserver.go"
kubectl="${GOPATH}/src/minik8s/cmd/kubectl/kubectl.go"
kubelet="${GOPATH}/src/minik8s/cmd/kubelet/kubelet.go"

# 构建apiserver
echo "Building apiserver..."
go build -o ${srv_bin_dir}/kube-apiserver ${apiserver}

# 构建kubectl
echo "Building kubectl..."
go build -o ${srv_bin_dir}/kubectl ${kubectl}

# 构建kubelet
echo "Building kubelet..."
go build -o ${srv_bin_dir}/kubelet ${kubelet}

# 启动apiserver和kubelet
echo "Starting apiserver and kubelet..."

