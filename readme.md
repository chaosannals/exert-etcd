# exert etcd

默认端口 2379

项目简单，依赖库少，不一定冲突。项目依赖库越多，冲突可能性越大。

GO 版本的 Client 库 GRPC 有版本限制要使用低版本 GRPC 不然编译出错。
直接安装，默认装上的 GRPC 版本太高，亲测时 v1.62.1 => v1.26.0 的退版本。
v1.64.0 也是一样，这个问题拖这么久，从 2020 年到 2024 年都没解决。

```bash
# 这个库使用的库有点多，当依赖多个库使用不同版本时，尽量让这个库版本降低。
# 参考示例 go.mod 的版本时 2019 年的。如果直接使用 2024 冲突了。
go get google.golang.org/genproto
```

```
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
```

如果出现缺少类的版本冲突 go.mod 修改：
```
github.com/golang/protobuf v1.5.4 // indirect
```

如果使用到了 依赖 bbolt 的库也会气冲突。例如使用了 swag 后 github.com/coreos/bbolt  和 go.etcd.io/bbolt 会有冲突。

解决方案是在 go.mod 里面加入 

```
replace github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4
replace go.etcd.io/bbolt v1.3.4 => github.com/coreos/bbolt v1.3.4
```


```bash
# 安装客户端库
go get go.etcd.io/etcd/clientv3

# 使用低版本 GRPC
go get google.golang.org/grpc@v1.26.0
```

## 安装

直接 [Github仓库](https://github.com/etcd-io/etcd) 的 [Releases](https://github.com/etcd-io/etcd/releases) 里面下载编译后的二进制文件，添加环境变量 PATH 路径。

## 命令

```bat
@rem 查看 服务器 版本
etcd --version

@rem 启动服务端
etcd

@rem 查看 命令行客户端 版本
@rem 可以通过环境变量切换版本 ETCDCTL_API=2 ETCDCTL_API=3
etcdctl version

@rem 通过命令行客户端写入键值（需要另外启动 服务端）
etcdctl put /test/foo "hello etcd"

@rem 通过命令行客户端读取键值（需要另外启动 服务端）
etcdctl get /test/foo

@rem 只打印值
etcdctl get /test/foo --print-value-only

@rem 限制 2 个，前缀搜索，只打印键
etcdctl get --prefix --limit=2 /test --keys-only

@rem 删除键（必须全称）会返回数量 0 就是没删到
etcdctl del /test/foo

@rem 删除键，指定前缀，会返回数量 0 就是没删到
etcdctl del --prefix /test

@rem 监听值变化
etcdctl watch /test/foo

@rem 监听值变化可以多次watch
etcdctl watch -i
@rem 之后多次键入
watch /test/foo
watch /test/bar
@rem 等等 ...

@rem 监听值，指定前缀
etcdctl watch --prefix /test

@rem 租约，生成，此时会返回租约的 id
@rem 租约到期就会自动销毁，连同绑定的 键值
etcdctl lease grant 60

@rem 撤销租约
etcdctl lease revoke 694d8e68f0b2372c

@rem 租约列举
etcdctl lease list

@rem 绑定租约，由上面生成的 id 做为 --lease 参数
etcdctl put --lease=694d8e68f0b2372c /test/lll "hhhhh"

@rem 续租约，自动续上租约
etcdctl lease keep-alive 694d8e68f0b2372c

@rem 获取租约信息
etcdctl lease timetolive --keys 694d8e68f0b2372c
```

## 杂项

有几个开源的项目依赖这个 etcd 做配置管理。
比如向量数据库 milvus 2.0 开始使用 etcd 做配置管理。

使用的时候基本通过 docker ，最近发现 etcd 是 golang 编写，官方文档的单机部署没有简单的 Windows 平台构建和发布。居然通过 docker 做 Windows 兼容。但是可以在 Github 上下载到编译打包的 Windows exe 文件。

## 替代项目

[kubebrain](https://github.com/kubewharf/kubebrain) 这个项目亦是没有简单的构建，居然使用 make 做构建。同样 golang 项目，却使用 make 。很是吊诡。该项目还用 etcd 的API协议。


[eureka](https://github.com/Netflix/eureka) Java 项目，没有试，感觉小公司的破服务器用 JAVA 会卡死。

[Consul](https://github.com/hashicorp/consul) GO 项目