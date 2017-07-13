FROM ubuntu:16.04

RUN apt-get update && apt-get install -y openssh-server apache2 golang git
RUN mkdir /var/run/sshd
RUN echo 'root:debian02' | chpasswd
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

# SSH login fix. Otherwise user is kicked off after login
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd

ENV NOTVISIBLE "in users profile"
RUN echo "export VISIBLE=now" >> /etc/profile

#  golang path setting
RUN mkdir -p ~/go; echo "export GOPATH=$HOME/go" >> ~/.bashrc
RUN echo "export PATH=$PATH:$HOME/go/bin:/usr/local/go/bin" >> ~/.bashrc
ADD ./go /root/go
RUN mkdir /root/log
WORKDIR /root/log

EXPOSE 22
EXPOSE 80
CMD ["/bin/bash"]
