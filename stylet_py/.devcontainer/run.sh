#!bin/sh
apt-get update
apt-get -y install curl
apt-get install -y openssh-server
cp .devcontainer/sshd_config /etc/ssh/
service ssh start
pip3 install -r .devcontainer/requirements.txt

# NOTE: To do the ssh between containers, you need to set password with:
# passwd root 