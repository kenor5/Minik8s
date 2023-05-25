#!/bin/bash

pw="$PWD/scripts"

$pw/docker_container_clear.sh
$pw/iptable_clear.sh
$pw/etcd_clear.sh

sudo systemctl restart docker
