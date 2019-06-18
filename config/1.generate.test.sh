#!/bin/bash

bash_path=$(cd `dirname $0`; pwd)

IFS=$'\n'

cd $bash_path
rm -rf build/*.sh

   for line in `cat $bash_path/config_test`;do
    conf_host=`echo $line | awk -F "=" '{print$1}' `
    start_bash=`echo $line | awk -F "=" '{print $NF}' `
    if [ ! -f "build/${conf_host}_test.sh" ]; then
        echo "#!/bin/bash" > build/${conf_host}_test.sh
    fi
    echo "Job Task Is Generate $conf_host"
    echo "cd /root/build/ && sh $start_bash " >> build/${conf_host}_test.sh
    done


CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/linux_ssynflood  ./
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/linux_ssynfloodv2  ./ver2.0/

cp ${bash_path}/cookie.json build/
cp ${bash_path}/url.json build/
chmod +x build/*_test.sh

