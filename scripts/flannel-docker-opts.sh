#!/bin/bash

# 读取当前机器上Flannel的配置信息，写入docker的配置文件中，一台机器上只需要执行一次
# 使用请确保/lib/systemd/system/docker.service文件中有:
# ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock --config-file=/etc/docker/daemon.json

FLANNEL_ENV_FILE="/run/flannel/subnet.env"
DOCKER_OPTS_FILE="/etc/docker/daemon.json"

if [[ -f "$FLANNEL_ENV_FILE" ]]; then
    # Read flannel environment variables
    source "$FLANNEL_ENV_FILE"
    DOCKER_OPTS="{\"bip\": \"${FLANNEL_SUBNET}\", \"mtu\": ${FLANNEL_MTU}"
    if [[ "$FLANNEL_IPMASQ" == "true" ]]; then
        DOCKER_OPTS="$DOCKER_OPTS, \"ip-masq\": true"
    else
        DOCKER_OPTS="$DOCKER_OPTS, \"ip-masq\": false"
    fi

    if [[ "$FLANNEL_IPTABLES" == "true" ]]; then
        DOCKER_OPTS="$DOCKER_OPTS, \"iptables\": true}"
    else
        DOCKER_OPTS="$DOCKER_OPTS, \"iptables\": false}"
    fi

    # Write docker options to file
    echo "$DOCKER_OPTS" > "$DOCKER_OPTS_FILE"
else
    echo "Flannel environment file $FLANNEL_ENV_FILE not found"
fi
