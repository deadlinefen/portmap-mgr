# tinyPortMapper-manager

[English](README_en.md)

使用[tinyPortMapper](https://github.com/wangyu-/tinyPortMapper)来实现基于DDNS+IPv6的反向代理。

tinyPortMapper-manager主要负责以下内容：

1. 对域名进行解析并监控IPv6地址的变化。
2. 根据配置文件中的任务信息批量启动与监控tinyPortManager进程。

为什么使用tinyPortMapper？

因为易用又可靠，比如我在使用socat配置幻兽帕鲁的ipv6到ipv4的映射时就始终没有成功: (，但是使用tinyPortMapper时就没有这个问题。

# 编译与使用

本项目使用make进行编译，通过指定config.toml配置文件来运行。

## 编译

```bash
go mod tidy
make manager
```

以上指令会在 `build/bin`下生成可执行文件 `pmmanager`。

注：当前版本可能由于没有正确使用[miekg/dns](https://github.com/miekg/dns)库的缘故，在编译时会比较慢，请耐心等待，但该问题不会影响使用性能。

## 使用

1. 下载或编译[tinyPortMapper](https://github.com/wangyu-/tinyPortMapper)。
2. 填写配置文件config.toml，详见[模板](config/config.toml)与[说明](doc/config_description.md)。
3. 使用以下命令行启动tinyPortMapper-manager：

```bash
/path/to/pmmanager -c /path/to/config/config.toml
```
