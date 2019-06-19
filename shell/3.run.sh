#!/bin/bash
bash_path=$(cd `dirname $0`; pwd)
IFS=$'\n'

for conf_host in `cat $bash_path/host`;do
    ssh -i pirkey -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@${conf_host} "cd /root/build/ && sh ${conf_host}.sh"
    echo "${conf_host} job finished"
done