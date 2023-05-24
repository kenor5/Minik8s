FROM python:3.8.13-slim-buster

MAINTAINER luoshicai <luoshicai@sjtu.edu.cn>

RUN pip install flask

RUN pip install -r ./requirement.txt

ADD ../serverless_server.py /serverless_server.py
ADD ./hello_function /hello_function.py

EXPOSE 8070

CMD ["python3", "/serverless_server.py"]
