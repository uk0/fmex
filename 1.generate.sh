#!/bin/bash

bash_path=$(cd `dirname $0`; pwd)

IFS=$'\n'

cd $bash_path

   for line in `cat $bash_path/config`;do
    conf_host=`echo $line | awk -F "=" '{print$1}' `
    start_bash=`echo $line | awk -F "=" '{print $NF}' `

        echo "Job Task Is Generate $conf_host"
        echo "#!/bin/bash
cd /root/build/ && nohup ./linux_ssynflood $start_bash >> result.log 2>&1 &
             " > build/$conf_host.sh
    done

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o linux_ssynflood   ./
cp ${bash_path}/linux_ssynflood build/
cp ${bash_path}/cookie.json build/
chmod +x build/*.sh

