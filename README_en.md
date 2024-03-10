# tinyPortMapper-manager

`README_en.md` has been translated using Baidu's Wenxin Yiyan (文心一言) translation tool.

Using [tinyPortMapper](https://github.com/wangyu-/tinyPortMapper) to implement a reverse proxy based on DDNS+IPv6.

The main responsibilities of tinyPortMapper-manager include:

1. Resolving domain names and monitoring changes in IPv6 addresses.
2. Batch starting and monitoring tinyPortManager processes based on the task information in the configuration file.

Why choose tinyPortMapper?

The main reasons are its ease of use and reliability. For example, when I tried to configure IPv6 to IPv4 mapping for PalWorld using socat, I was unsuccessful : (. However, with tinyPortMapper, I didn't encounter any such issues.

# Compilation and Usage

This project uses make for compilation and runs by specifying the config.toml configuration file.

## Compilation

```bash
go mod tidy
make manager
```

The above instructions will generate the executable file `pmmanager` under `build/bin`.

Note: The current version may compile slower due to improper use of the [miekg/dns](https://github.com/miekg/dns) library. Please be patient during the compilation process, but this issue will not affect the performance of the usage.

## Usage

1. Download or compile [tinyPortMapper](https://github.com/wangyu-/tinyPortMapper).

2. Fill out the configuration file `config.toml`, refer to the [template](config/config.toml) and [instructions](doc/config_description_en.md) for more details.
3. Use the following command to start tinyPortMapper-manager:

```bash
/path/to/pmmanager -c /path/to/config/config.toml
```
