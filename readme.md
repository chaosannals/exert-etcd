# exert etcd

有几个开源的项目依赖这个 etcd 做配置管理。
比如向量数据库 milvus 2.0 开始使用 etcd 做配置管理。

使用的时候基本通过 docker ，最近发现 etcd 是 golang 编写，却 没有简单的 Windows 平台构建和发布。居然通过 docker 做 Windows 兼容。有点吊诡。

## 替代项目

[kubebrain](https://github.com/kubewharf/kubebrain) 这个项目亦是没有简单的构建，居然使用 make 做构建。同样 golang 项目，却使用 make 。很是吊诡。该项目还用 etcd 的API协议。


