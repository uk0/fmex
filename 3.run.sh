#!/usr/bin/env bash
bash_path=$(cd `dirname $0`; pwd)
IFS=$'\n'

for line in `cat $bash_path/config`;do
    conf_host=`echo $line | awk -F "=" '{print$1}' `
    ssh -i pirkey -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@${conf_host} "cd /root/build/ && sh ${conf_host}.sh"
    echo "${conf_host} job finished"
done