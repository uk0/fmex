#!/bin/bash
bash_path=$(cd `dirname $0`; pwd)
IFS=$'\n'

for conf_host in `cat $bash_path/host`;do
    ssh -i pirkey root@$conf_host "killall linux_ssynfloodv2"
    echo "${conf_host} kill linux_ssynflood"
done