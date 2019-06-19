#!/bin/bash

bash_path=$(cd `dirname $0`; pwd)

IFS=$'\n'

bash_path=$bash_path/../config/

cd $bash_path

rm -rf ../build/*

   for line in `cat $bash_path/config`;do
    conf_host=`echo $line | awk -F "=" '{print$1}' `
    start_bash=`echo $line | awk -F "=" '{print $NF}' `

        echo "Job Task Is Generate $conf_host"
        echo "#!/bin/bash
cd /root/build/ && nohup  $start_bash >> result.log 2>&1 &
             " > ../build/$conf_host.sh
    done

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../build/fmex  ../

cp ${bash_path}../config/cookie.json ../build/
cp ${bash_path}../config/url.json ../build/
cp ${bash_path}../config/host ../build/
cp ${bash_path}../shell/2.scp.sh ../build/scp_to_all.sh
cp ${bash_path}../shell/3.run.sh ../build/run_all.sh
cp ${bash_path}../shell/9.status.sh ../build/show_all_status.sh
cp ${bash_path}../shell/4.kill.sh ../build/kill_all.sh
cp ${bash_path}../config/pirkey ../build/pirkey
chmod +x ../build/*.sh

