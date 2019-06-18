#!/bin/bash
bash_path=$(cd `dirname $0`; pwd)
IFS=$'\n'

for line in `cat $bash_path/config_test`;do
    conf_host=`echo $line | awk -F "=" '{print$1}' `
    ssh -i pirkey -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@${conf_host} "cd /root/build/ && nohup sh ${conf_host}_test.sh &"
    echo "${conf_host} job finished"
done