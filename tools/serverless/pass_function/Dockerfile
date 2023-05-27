FROM python:3.8.13-slim-buster

MAINTAINER luoshicai <luoshicai@sjtu.edu.cn>

RUN pip install flask

COPY . /src

RUN pip install -r /src/requirements.txt

EXPOSE 8070

CMD ["python3", "/src/serverless_server.py"]
