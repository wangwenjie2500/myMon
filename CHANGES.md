# 更新日志

## [version] - 2023-08-30

本次更新，主要是在代码规范、代码结构及代码测试上的更新，增加了少量必要的监控项。

### Added

- 新增了cron任务管理,不用托管到系统的定时任务,每分钟执行一次
- 新增了监控项`tps,Uptime`


### Update
- 修改了配置文件格式为toml
- 支持配置多个db实例监控
- 单实例监控超时控制在30S,超过30S直接退出
- 调整了一些指标的数据类型改为`orgin`指标如下`Uptime,Opened_files,Opened_tables,Threads_connected,Threads_running,Slave_SQL_Running,Slave_IO_Running,io_thread_delay,Seconds_Behind_Master,Is_Slave`

### Delete
- 删掉了process监控和快照功能