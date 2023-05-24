FROM python:3.8.13-slim-buster

MAINTAINER luoshicai <luoshicai@sjtu.edu.cn>

# 安装Flask
RUN pip install flask

RUN pip install invoke decorator fabric

# define a mount point
#RUN #ssh-keygen -q -t rsa -N '' -f /id_rsa
ADD ./server.py /server.py


EXPOSE 8090

# define command for this image
#CMD ["/bin/bash", "a.sh"]
CMD ["python3","/server.py"]
