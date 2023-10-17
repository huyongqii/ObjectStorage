# [mgface.com](https://www.mgface.com)-分布式对象存储系统(DOSS)

### ①描述

一款被设计为具备高可用，可伸缩型的对象存储系统，它同时还兼具自动修复和故障转移等高级特性。这款OSS系统采用高性能语言golang开发。原生支持windows，linux系统的部署，同时系统还兼容在通用的云平台环境(Kebenetes)直接部署使用。

### ②功能介绍

1.  ✔:提供对象上传下载功能，当前支持HTTP协议访问
2.  ✔:对象具有多版本、去重关联管理
3.  ✔:对象具有冗余分片储存功能
4.  ✔:具有后台数据自动修复功能
5.  ✔:系统具有自动再平衡数据分布
6.  ✔:客户端具有故障自动切换(元数据服务metadata故障)
7.  ✔:支持全量元数据同步


### ③软件架构

![avatar](./doc/design/system_design.png)
用户通过公网访问到前端负载均衡器ingress(traefik)，由ingress直接路由apinode服务，把上传的对象元数据(对象大小，文件名，hash码，上传时间等)
存储到metanode服务中，并且把相关的数据写入到datanode，datanode写成功之后，把存储的相关元数据(存储位置，存储分片的hash码值，存储的大小等)也汇报给metanode。由metanode统一来管理各个对象的元数据。

### ④安装和使用

**开发模式**

1. 先启动startmetanode.go
2. 后启动startdatanode.go
   (保证启动≥3个，代码控制了对象需要分片，默认2个分片+1个校验片，可以修改以下代码)
   ```bash
   # cd com.mgface.disobj/apinode/api
   # vim initval.go 
   ```
3. 最后启动startapinode.go

**k8s(v1.19.3)部署模式**

1. 一键启动服务

```bash
# ./deployService.sh
```

2. 一键删除服务

```bash
# ./deleteService.sh
```


### ⑤备注

* 大文件传输处理:客户端直接选择其中任何一个datanode节点的IP地址，先把大文件POST传输到datanode上面临时存放，然后返回 该文件的sha256哈希码值和[大小(主要是为了metadata里面需要保持文件大小)]
  后续再次PATCH该对象，确认真正写入成功，否则 如果在一定时间没有PATCH的话，删除该临时文件

### ⑥计划功能

Ⅰ. ✔修复metadatanode停止之后，datanode节点不像其他master汇报心跳  
Ⅱ. ✔添加metadatanode自身元数据节点的元数据，方便datanode和apinode做负载均衡使用和故障容错  
Ⅲ. ✔修复metadatanode挂了之后，重新启动丢失元数据数据

### ⑦系统文档功能

**项目文档使用**

```bash
#安装go doc服务（go1.12版本移除了godoc，需要手动安装）
mkdir -p $GOPATH/src/golang.org/x
cd $GOPATH/src/golang.org/x
git clone https://github.com/golang/tools.git
cd tools/command/godoc
go install
#启动go doc服务
godoc -http=:6060  -goroot="C:\xxxx\GoglandProjects\distributedObjStorage"
```
