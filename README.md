## Fcoin hack Script 

### Status 


* 1st round 5000$ * 2  `10W FMEX` Date: 20190613
* 2nd round  5000$ * 4  `18W FMEX` Date: 20190614
* 3rd round  5000$ * 4  `20W FMEX` Date: 20190615
* 4th round  5000$ * 3  `15W FMEX` Date: 20190616
* 5th round  5000$ * 6  `30W FMEX` Date: 20190617


### Config 

* 修改`cmd` `host` 保持数量一致。
 
* cmd
```text
./linux_ssynflood --time "340" --name "smartguo_guo"
./linux_ssynfloodv2 --token_time "300" --buy_time "250"  --name "smartguo_lu"
./linux_ssynfloodv2 --token_time "350" --buy_time "250"  --name "smartguo_andrew"
./linux_ssynfloodv2 --token_time "300" --buy_time "250"  --name "wing_big"
./linux_ssynfloodv2 --token_time "300" --buy_time "250"  --name "wing_small"
./linux_ssynfloodv2 --token_time "350" --buy_time "250"  --name "zhangjianxin"
./linux_ssynflood --time "340"  --name "wing_small6"
./linux_ssynflood --time "340"  --name "wing_small5"
./linux_ssynflood --time "340"  --name "wing_small4"
./linux_ssynflood --time "340"  --name "wing_small3"
./linux_ssynflood --time "340"  --name "wing_small2"
./linux_ssynfloodv2 --token_time "350" --buy_time "250"  --name "wing_small8"
./linux_ssynfloodv2 --token_time "350" --buy_time "250"  --name "wing_small7"
./linux_ssynfloodv2 --token_time "350" --buy_time "250"  --name "wing_small1"
./linux_ssynfloodv2 --token_time "350" --buy_time "250"  --name "smartguo_xiu"
./linux_ssynflood --time "340"  --name "smartguo_pang"
```

* host 

```text
127.0.0.1
127.0.0.2
127.0.0.3
127.0.0.4
127.0.0.5
127.0.0.6
127.0.0.7
127.0.0.8
127.0.0.9
127.0.0.10
127.0.0.11
127.0.0.12
127.0.0.13
127.0.0.14
127.0.0.15
127.0.0.16
```

* 执行 `0.host_cmd_to_config` 生成配置文件 `config`
* 执行 `1.generate` 生成`Build `文件夹内的文件 其中主要检查 `Url.json` `cookie.json` 以及抽样检查一个`sh`脚本
* 执行 `2.scp` 分发脚本到每个机器
* 执行 `3.run` 将每个机器上的脚本启动
* 执行 `4.kill` 停止所有机器上的服务
* 执行 `5.result` 最后日志的100行进行查看
* 执行 `6.delete_log_sh` 清理过期的脚本以及日志 `！！！不会删除二进制执行文件`
* 执行 `9.status ` 查看机器上服务是否启动
* 执行 `10.scp_url_json` 分发url.json 以及 sh脚本 



### Update

* 更新抢购程序 `抛弃第一版，第二版更优`
* 参数调整 `等待分析结果`


### Quick Start


* 0.host_cmd_to_config 
* 检查Config 文件
* 1.generate
* 检查Build文件夹内容正常
* scp_to_master.sh 
* 将文件scp到Master 
* 3.run  启动
* 9.status 检查


* 测试流程 run `3.run`起来 查看状态`9.status` 以及 日志 `5.result` 杀死`4.kill` 清理日志 `6.delete_log_sh`