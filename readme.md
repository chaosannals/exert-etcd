# exert etcd

## 安装

直接 [Github仓库](https://github.com/etcd-io/etcd) 的 [Releases](https://github.com/etcd-io/etcd/releases) 里面下载编译后的二进制文件，添加环境变量 PATH 路径。

## 杂项

有几个开源的项目依赖这个 etcd 做配置管理。
比如向量数据库 milvus 2.0 开始使用 etcd 做配置管理。

使用的时候基本通过 docker ，最近发现 etcd 是 golang 编写，官方文档的单机部署没有简单的 Windows 平台构建和发布。居然通过 docker 做 Windows 兼容。但是可以在 Github 上下载到编译打包的 Windows exe 文件。

## 替代项目

[kubebrain](https://github.com/kubewharf/kubebrain) 这个项目亦是没有简单的构建，居然使用 make 做构建。同样 golang 项目，却使用 make 。很是吊诡。该项目还用 etcd 的API协议。


[eureka](https://github.com/Netflix/eureka) Java 项目，没有试，感觉小公司的破服务器用 JAVA 会卡死。

[Consul](https://github.com/hashicorp/consul) GO 项目