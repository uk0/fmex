#!/bin/bash
bash_path=$(cd `dirname $0`; pwd)
IFS=$'\n'

for line in `cat $bash_path/config_test`;do
    conf_host=`echo $line | awk -F "=" '{print$1}' `
    ssh -i pirkey root@$conf_host "mkdir -p /root/build/"
    scp -i pirkey -r build/*  root@$conf_host:/root/build/
    echo "${conf_host} copy files finished"
done