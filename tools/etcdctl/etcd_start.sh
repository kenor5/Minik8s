#!/bin/sh

etcdctl member list >/dev/null 2>&1

if [ $? != 0 ]; then
    nohup etcd &
else
    echo "etcd already running"
fi
