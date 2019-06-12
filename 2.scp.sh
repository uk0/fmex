#!/usr/bin/env bash
bash_path=$(cd `dirname $0`; pwd)
IFS=$'\n'

for line in `cat $bash_path/config`;do
    conf_host=`echo $line | awk -F "=" '{print$1}' `
    scp -i pirkey -r build/  root@$conf_host:/root/
    echo "${conf_host} copy files finished"
done