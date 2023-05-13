#!/bin/sh

data_dir="$PWD/etcdstore"
localurl=$(ifconfig | grep 192.168 | (awk -v OFS="" '{print "http://",$2,":2379"}'))
listen_client_url="${localurl},http://127.0.0.1:2379"
advertise_client_url="${localurl},http://localhost:2379"

etcdctl member list >/dev/null 2>&1

if [ $? != 0 ]; then
    nohup etcd --data-dir ${data_dir} --advertise-client-urls ${advertise_client_url} --listen-client-urls ${listen_client_url} &
else
    echo "etcd already running"
fi
