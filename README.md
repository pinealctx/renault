# Renault

Renault, useful tools.

## Install Renault

```shell
go get github.com/pinealctx/renault@master
```

## Commands 

### 初始化工作区

```shell
renault workspace init
renault w init
```

### 工作区新增新项目

```shell
renault workspace add --url=git@gl.codectn.com:hermes/user.git
renault w add --url=git@gl.codectn.com:hermes/user.git
```

### 同步工作区并拉取最新代码

```shell
renault workspace sync
renault w sync
```

### 初始化项目结构

```shell
renault project init --name=github.com/pinealctx/renault
```

## TODO

后续将会加入更多有利于工程化的工具链，对于仓库结构以及git管理起到辅助作用。

希望可以通过此类工具，能够提升日常开发效率。

将区分 workspace 以及 project 的概念

- workspace，用于维护工作区相关的功能，例如初始化工作区，添加工作区项目，批量拉取工作区的项目，批量推送工作区的项目等等
- project，用于维护指定项目的一些功能，如：初始化项目脚手架等等