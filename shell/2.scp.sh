#!/bin/bash
bash_path=$(cd `dirname $0`; pwd)
IFS=$'\n'

for conf_host in `cat $bash_path/host`;do
    ssh -i  pirkey -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$conf_host "mkdir -p /root/build/"
    ssh -i  pirkey -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$conf_host "ls -sort  /root/build/ "
    scp -i  pirkey -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null *.json fmex root@$conf_host:/root/build/
    scp -i  pirkey -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${conf_host}.sh  root@$conf_host:/root/build/
    ssh -i  pirkey -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null  root@$conf_host "ls  -sort /root/build/"
    echo "${conf_host} copy files finished"
done