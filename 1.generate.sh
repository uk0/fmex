#!/bin/bash

bash_path=$(cd `dirname $0`; pwd)

IFS=$'\n'

cd $bash_path
rm -rf build/*

   for line in `cat $bash_path/config`;do
    conf_host=`echo $line | awk -F "=" '{print$1}' `
    start_bash=`echo $line | awk -F "=" '{print $NF}' `

        echo "Job Task Is Generate $conf_host"
        echo "#!/bin/bash
cd /root/build/ && nohup  $start_bash >> result.log 2>&1 &
             " > build/$conf_host.sh
    done

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/linux_ssynfloodv2  ./ver2.0/

cp ${bash_path}/cookie.json build/
cp ${bash_path}/url.json build/
cp ${bash_path}/host build/
cp ${bash_path}/2.scp.sh build/scp_to_all.sh
cp ${bash_path}/3.run.sh build/run_all.sh
cp ${bash_path}/9.status.sh build/show_all_status.sh
cp ${bash_path}/4.kill.sh build/kill_all.sh
cp ${bash_path}/pirkey build/pirkey
chmod +x build/*.sh

