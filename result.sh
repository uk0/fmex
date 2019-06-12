#!/usr/bin/env bash
bash_path=$(cd `dirname $0`; pwd)

IFS=$'\n'

for line in `cat $bash_path/config`;do
    conf_host=`echo $line | awk -F "=" '{print$1}' `
    ssh -i pirkey root@${conf_host} "cat /root/build/result.log | head -1"
done