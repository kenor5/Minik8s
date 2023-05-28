#!/bin/bash

etcdctl del "Pod" --prefix
etcdctl del "Service" --prefix
etcdctl del "Deploy" --prefix
etcdctl del "Node" --prefix
etcdctl del "Dns" --prefix
etcdctl del "Job" --prefix
etcdctl del "Function" --prefix

# etcdctl put /coreos.com/network/config '{"Network": "10.0.0.0/16", "SubnetLen": 24, "SubnetMin": "10.0.1.0", "SubnetMax": "10.0.20.0", "Backend": {"Type": "vxlan"}}'