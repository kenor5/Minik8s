#!/bin/bash

etcdctl del "" --prefix
etcdctl put /coreos.com/network/config '{"Network": "10.0.0.0/16", "SubnetLen": 24, "SubnetMin": "10.0.1.0", "SubnetMax": "10.0.20.0", "Backend": {"Type": "vxlan"}}'