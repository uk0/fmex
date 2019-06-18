#!/bin/bash
bash_path=$(cd `dirname $0`; pwd)
IFS=$'\n'

cmds=(`cat $bash_path/shell`)
hosts=(`cat $bash_path/host`)

yes | rm $bash_path/config


if (( "${#cmds[@]}"== "${#hosts[@]}"))
then
for (( i=0; i<${#cmds[@]}; i++ )) ; do
    echo ${hosts[$i]}"="${cmds[$i]} >> $bash_path/config
done
else
  echo "length is error cmd : [${#cmds[@]}] != hosts: [${#hosts[@]}] "
fi

