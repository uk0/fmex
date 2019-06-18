#!/bin/bash
bash_path=$(cd `dirname $0`; pwd)
IFS=$'\n'

for line in `cat $bash_path/config | head -1`;do
    conf_host=`echo $line | awk -F "=" '{print$1}' `
    ssh -i pirkey root@$conf_host "mkdir -p /root/build/"
    scp -i pirkey -r build/ root@$conf_host:/root/
    echo "${conf_host} copy files finished"
done

echo "${conf_host} Auto copy files finished"
ssh -i pirkey -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@${conf_host} "cd /root/build/ && ./scp_to_all.sh"
