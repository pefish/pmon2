# pmon2
go进程管理工具(golang process manager)，专门用于 go 常驻进程管理 （daemon manager）

<img src="http://p0.qhimg.com/t017d6cbb68aed4b693.png" style="max-width:680px" />

## 启动进程

```go
sudo pmon2 run [./二进制进程文件] [参数1]  ...
```

## 介绍

Go官方一直没有提供任何进程管理工具，对于 `Go` 服务的部署，简单的服务，我们使用 `linux` 内建命令 `nohup &`组合，或者使用系统自带进程管理工具, 例如 `systemd`，`init.d` 配置。或者借助第三方的进程管理工具，如：`python` 的 `supervisor` 或者 `nodejs` 的 `pm2`。

每种方式都有一定优劣，我们希望继承 `go` 语言部署集成的便捷易用思想，不需要再安装其他依赖软件，并且提升进程管理工具的体验。

和 `pm2` 不一样的是，`pmon2` 直接是启动的系统级常驻进程。进程直接由 `init` 来管理，因此，就算`pmon2`管理工具异常退出，也不会影响进程本身运行。

默认的，`pmon2` 所管理的进程如果发生异常崩溃，`pmon2` 会尝试重启该进程。如果你不希望某个进程自动重启，那么你可以设置一个 `--no-autorestart` 参数即可。具体请参考：常用命令使用章节。


## 如何安装

任意平台都可以选择以下方式安装：

```bash
go install github.com/pefish/pmon2/cmd/...@latest
```

目前 `Pmon2` 支持 `CentOS6`、`CentOS7`、`CentOS8`

[Releases](https://github.com/pefish/pmon2/releases) 中已经构建了可以直接安装的`rpm`包，可直接选择平台安装： 

```bash
# CentOS6
sudo yum install -y https://github.com/pefish/pmon2/releases/download/v1.12.0/pmon2-1.12.1-1.el6.x86_64.rpm

# CentOS7
sudo yum install -y https://github.com/pefish/pmon2/releases/download/v1.12.0/pmon2-1.12.1-1.el7.x86_64.rpm

# CentOS8
sudo yum install -y https://github.com/pefish/pmon2/releases/download/v1.12.0/pmon2-1.12.1-1.el8.x86_64.rpm
```

:exclamation::exclamation: **注意：** :exclamation::exclamation:

首次安装 `pmon2` 后，`pmon2` 服务没有自动启动，需要你手动启动该服务：

```bash
# centos6 使用 initctl
sudo initctl start pmon2

# centos7 使用 systemd
sudo systemctl start pmon2

# Debian/Ubuntu
sudo /usr/local/pmon2/bin/pmond &
```

## 命令介绍

#### 查看帮助

```sh
# 查看全局的帮助文档
sudo pmon2 help

# 查看某个具体命令 help
sudo pmon2 [command] help
```

帮助文档中能查看当前所有命令使用说明

```
Usage:
  pmon2 [command]

Available Commands:
  del         del process by id or name
  desc        print the process detail message
  exec        run one binary golang process file
  help        Help about any command
  ls          list all processes
  reload      reload some process
  start       start some process by id or name
  stop        stop running process
  log         display process log by id or name
  logf        display process log dynamic by id or name
  version     show pmon2 version
```

#### 运行进程 [run/exec]

```go
sudo pmon2 run [./二进制文件名] [参数1] [参数2] ...
```
启动进程，可以传入若干参数，参数说明如下：

```go
// 进程名称，如果不传递，则以二进制文件名作为默认名称
--name

// 进程运行日志，不配置则使用默认路径：/var/log/pmon2/
--log   -l

// 仅自定义日志目录，优先级低于 --log
--log_dir -d

// 进程自定义参数，多个参数以空格分割
--args  -a "-arg1=val1 -arg2=val2"

// 进程启动用户
--user  -u

// 不自动重启, 默认自动重启
--no-autorestart    -n
```

#### 示例：

```go
sudo pmon2 run ./bin/gin --args "-prjHome=`pwd`" --user ntt360
```
运行 `bin/gin` 二进制文件，并且传递参数：`-prjHome` 为当前目录，并且制定进程运行的用户为：`ntt360`

!!!warning "注意"
    pmon2启动的进程会另起一个隔离的上下文环境，自定义参数需要使用绝对路径，不能使用相对路径。

#### 查看列表  [ list/ls ]

```go
sudo pmon2 ls
```

#### 启动进程  [ start ]

```go
sudo pmon2 start [id or name]
```

#### 停止进程  [ stop ]

```go
sudo pmon2 stop [id or name]
```

#### 重载进程 [ reload ]

```go
sudo pmon2 reload [id or name]
```

#### 插件进程日志

```shell
# 查看最近进程的日志
sudo pmon2 log [id or name]

# 动态查看进程日志，类似系统tail -f xxx.log
sudo pmon2 logf [id or name]
```

仅仅重启配置文件，该命令需要所启动的进程配合使用，`reload` 命令默认仅仅发送 `SIGUSR2` 信号给启动的进程

如果你希望 `reload` 时自定义信号，那么请使用 `--sig` 参数：

```go
// 目前支持的信号：HUP、USR1、USR2
sudo pmon2 reload --sig HUP [id or name]
```

#### 删除进程  [ del/delete ]

```go
sudo pmon2 del [id or name]
```

#### 查看详情  [ show/desc ]

```go
sudo pmon2 show [id or name]
```
![](https://jscssimg-img.oss-cn-beijing.aliyuncs.com/89c3f649a583a852.png?t=1506950494)

## 常见问题

### 1.日志切割？

pmon2 自带一个 `logrotate` 日志切割配置文件，会默认切割 `/var/log/pmon2` 目录下的日志文件，如果你是自定义日志路径，请自行实现日志切割。

### 2.进程启动参数必须传绝对路径？

启动进程是，如果你传递的参数中存在路径，请使用绝对路径，`pmon2` 启动进程会新起一个新的沙盒环境以避免环境变量污染问题。

### 3. 平台支持

目前 `rpm` 适配 `CentOS6` 、 `CentOS7`、`CentOS8`， `Pmon2` 本身可运行在任何 `linux` 环境下，如有其它平台打包需求，请联系我们。

### 4. 命令行自动补全

pmon2 提供 bash 自动补全脚本，如果你发觉在 `sudo` 模式下命令无法自动补全，请安装 `bash-completion`，退出终端重新进入即可：

```bash
sudo yum intsall -y bash-completion
```

### 5. FATA stat /var/run/pmon2/pmon2.sock: no such file or directory

如果遇到如上报错，请尝试运行：

```bash
# centos6 使用 initctl
sudo initctl start pmon2

# centos7 使用 systemd
sudo systemctl start pmon2
```

原因请参考，安装启动部分说明。
