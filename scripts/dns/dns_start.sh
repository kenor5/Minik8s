#!/bin/bash

docker run -d \
    -p 80:80 \
    -v /root/go/src/minik8s/configs/nginx/default.conf:/etc/nginx/conf.d/default.conf \
    -v /root/go/src/minik8s/configs/nginx/log:/var/log/nginx \
    --rm  \
    --name n_test nginx