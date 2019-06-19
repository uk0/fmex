#!/bin/bash
bash_path=$(cd `dirname $0`; pwd)
IFS=$'\n'

for conf_host in `cat $bash_path/host`;do
    ping $conf_host -c 5
done