#!/bin/bash
bash_path=$(cd `dirname $0`; pwd)
IFS=$'\n'


${bash_path}/1.generate

echo "已经重新生成脚本！可以进行分发。"

ls -sort ${bash_path}/build/ | awk '{print$8,$NF}'

echo "请检查时间是否正确：`date '+%Y%m%d %H:%M:%S'`"

for line in `cat $bash_path/config | head -1`;do
    conf_host=`echo $line | awk -F "=" '{print$1}' `
    scp -i pirkey -r build/${conf_host}.sh  build/*.json root@$conf_host:/root/build/
    echo "${conf_host} copy files finished"
done