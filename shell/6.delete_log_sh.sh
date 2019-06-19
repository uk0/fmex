#!/bin/bash
bash_path=$(cd `dirname $0`; pwd)

IFS=$'\n'

bash_path=$bash_path/../config/

cd ${bash_path}

for conf_host in `cat $bash_path/host`;do
    ssh -i pirkey root@${conf_host} "rm -rf  /root/build/*.log"
done