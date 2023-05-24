#!/bin/bash

function_name=$1

dockerfile_content=$(cat <<EOF
FROM python:3.8.13-slim-buster

MAINTAINER luoshicai <luoshicai@sjtu.edu.cn>

RUN pip install flask

RUN pip install -r ./requirement.txt

ADD ../serverless_server.py /serverless_server.py
ADD ./${function_name} /${function_name}.py

EXPOSE 8070

CMD ["python3", "/serverless_server.py"]
EOF
)

echo "$dockerfile_content" > ./tools/serverless/${function_name}/Dockerfile

