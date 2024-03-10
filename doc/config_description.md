# config.toml

[English](config_description_en.md)

配置文件包含四个部分：mapper、resolution、jobs、log。其中mapper为tinyPortMapper相关配置，resolution为域名解析相关配置，jobs为端口映射任务配置，log为日志配置。

## mapper

1. ****bin****
   tinyPortMapper可执行文件的所在路径。
2. **file-directory**
   用于保存tinyPortMapper产生日志文件的**目录**，由于可能会有多个映射关系，因此这里填写的是日志的目录，文件命名格式为 `fromPort-toIp:toPort.log`。

## resolution

1. domain
   需要被解析的域名，当前版本只支持一个域名。
2. dns
   用于解析的dns服务器地址列表，为了防止频繁访问导致节点被拉黑，建议多填几个，manager会依次遍历dns服务器列表进行域名解析。
3. ttl
   每过ttl秒，从dns列表中选择一个dns服务器进行一次域名解析。

## jobs

可以填写多个job，每个job对应一个端口映射任务，每个任务以 `[jobs.job_name]`标识。

1. from-port
   对应tinyPortMapper的local-port
2. to-ip
   对应tinyPortMapper的remote-ip
3. to-port
   对应tinyPortMapper的remote-port
4. type
   对应tinyPortMapper的映射类型，`"t"`对应tcp，`"u"`对应udp，`tu`对应tcp与udp

## log

1. level
   日志等级，debug > info > warn > error > fatal > panic，日常建议使用info
2. path
   日志路径
3. to-stdout-only
   仅输出到stdout，日常建议使用false
4. also-to-stderr
   同时输出到stderr
