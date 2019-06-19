shc -f 1.generate.sh -o 1.generate.run
shc -f 2.scp.sh -o 2.scp.run
shc -f 3.run.sh -o 3.run.run
shc -f 4.kill.sh -o  4.kill.run
shc -f 5.result.sh -o 5.result.run
shc -f 6.delete_log_sh.sh -o 6.delete_log_sh.run
shc -f 9.status.sh -o 9.status.run
shc -f 0.host_cmd_to_config.sh -o 0.host_cmd_to_config.run

rm -rf *.c