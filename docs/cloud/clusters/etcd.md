# ETCD集群安装

etcd 目前默认使用 2379 端口提供 HTTP API 服务，2380 端口提供 Peer 通信（这两个端口已经被 IANA 官方预留给 etcd），在之前的版本中，可能会分别使用 4001 和 7001，在使用的过程中需要注意这个区别。

虽然 etcd 也支持单点部署，但是在生产环境中推荐集群方式部署，一般 etcd 节点数会选择 3、5、7。etcd 会保证所有的节点都会保存数据，并保证数据的一致性和正确性。



## 安装

### 编译安装

因为 etcd 是 Golang 编写的，安装只需要下载对应的二进制文件，并放到合适的路径就行。如果在测试环境，启动一个单点的 etcd 服务，只需要运行 etcd 执行即可。

```bash
git clone https://github.com/etcd-io/etcd.git
cd etcd
./build
```

使用 build 脚本构建会在当前项目的 bin 目录生产 etcd 和 etcdctl 可执行程序。`etcd `就是 etcd server，而 `etcdctl `主要为 etcd server 提供指令行操作。

查看版本：

```bash
$ ./bin/etcd --version
etcd Version: 3.5.0-pre
Git SHA: ab4cc3cae
Go Version: go1.14.4
Go OS/Arch: darwin/amd64

$ ./bin/etcdctl version
etcdctl version: 3.5.0-pre
API version: 3.5
```

### 直接下载二进制文件

为了方便下载最新版本使用下面脚本进行安装：

```shell
$ download_url=`curl -q https://api.github.com/repos/etcd-io/etcd/releases | jq -r '.[0].assets[].browser_download_url'  | grep $(uname -s) -i | head -n 1`
$ wget $download_url
```

或者您直接下载: https://github.com/etcd-io/etcd/releases

### Docker 方式安装

docker 安装只需要pull下载相应的版本即可 `quay.io/coreos/etcd:v3.5.0` 。[国内加速参考](../docker/image_pull.md)



## 单点部署

### 直接启动

```shell 
./bin/etcd -advertise-client-urls http://0.0.0.0:2379 -listen-client-urls http://0.0.0.0:2379
```

### 容器部署

```shell
$ docker run --rm --name etcd -p 2379:2379 -e ETCDCTL_API=3 quay.io/coreos/etcd:v3.5.0 /usr/local/bin/etcd \
		-advertise-client-urls http://0.0.0.0:2379 \
		-listen-client-urls http://0.0.0.0:2379
    
$ docker ps -f name=etcd

$ docker exec -it etcd etcdctl version
etcdctl version: 3.3.8
API version: 3.3

$ docker exec -it etcd etcdctl endpoint health
127.0.0.1:2379 is healthy: successfully committed proposal: took = 571.723µs
```

## 关键启动参数清单说明

| 参数名                       | 说明                                                         |
| ---------------------------- | ------------------------------------------------------------ |
| -name                        | 指定 etcd node 名称，可以使用 hostname。                     |
| –data-dir                    | 指定 etcd server 持久化数据存储目录路径。                    |
| –snapshot-count              | 指定有多少事务（transaction）被提交后，触发截取快照并持久化到磁盘。 |
| –heartbeat-interval          | 指定 Leader 多久发送一次心跳到 Followers。                   |
| –eletion-timeout             | 指定重新投票的超时时间，如果 Follow 在该时间间隔没有收到 Leader 发出的心跳包，则会触发重新投票。 |
| –listen-peer-urls            | 指定和 Cluster 其他 Node 通信的地址，比如：http://IP:2380，如果有多个，则使用逗号分隔。需要所有节点都能够访问，所以不要使用 localhost。 |
| –listen-client-urls          | 指定对外提供服务的地址，比如：http://IP:2379,http://127.0.0.1:2379。 |
| –advertise-client-urls       | 对外通告的该节点的客户端监听地址，会告诉集群中其他节点。     |
| –initial-advertise-peer-urls | 对外通告该节点的同伴（Peer）监听地址，这个值会告诉集群中其他节点。 |
| –initial-cluster             | 指定集群中所有节点的信息，通常为 IP:Port 信息，格式为：node1=http://ip1:2380,node2=http://ip2:2380,…。注意，这里的 node1 就是 --name 指定的名字，ip1:2380 就是 --initial-advertise-peer-urls 指定的值。 |
| –initial-cluster-state       | 新建集群时，这个值为 new；假如已经存在了集群，这个值为 existing。 |
| –initial-cluster-token       | 创建集群的 token，这个值每个集群保持唯一。这样的话，如果你要重新创建集群，即使配置和之前一样，也会再次生成新的集群和节点 UUID；否则会导致多个集群之间的冲突，造成未知的错误。 |
|                              |                                                              |




## 集群部署

etcd cluster 的部署，实际上就是多个主机上都部署 etcd server，然后将它们加入到一个 cluster 中。

 **注意：etcd cluster 必须具有时间同步服务器，否则会导致 Cluster 频繁进行 Leader Change。在 OpenShift 的 etcd cluster 中，会每隔 100ms 进行心跳检测。**

在安装和启动 etcd 服务进程的时候，各个 Node 都需要知道 Cluster 中其他 Nodes 的信息，一般是 IP:Port 信息。

根据用户是否提前知晓（规划）了每个 Node 的 IP 地址，有以下几种不同的集群部署方案：

-   1、静态配置

    在启动 etcd server 的时候，通过 --initial-cluster 参数配置好所有的节点信息。

-   2、公开的发现服务

    注册到已有的 etcd cluster：比如官方提供的 `discovery.etcd.io`

#### 1、静态配置集群

如果 etcd cluster 中的成员是已知的，且具有固定的 IP 地址，就可以静态的初始化一个集群。

环境准备：三台主机（10.0.2.1\10.0.2.2\10.0.2.3） 

每个 Node 都使用如下环境变量：

```shell
ETCD_INITIAL_CLUSTER="node1=http://10.0.2.1:2380,node2=http://10.0.2.2:2380,node3=http://10.0.2.3:2380"
ETCD_INITIAL_CLUSTER_STATE=new
```

或者使用如下指令行参数来指定集群成员：

```shell
--initial-cluster node1=http://10.0.2.1:2380,node2=http://10.0.2.2:2380,node3=http://10.0.2.3:2380
--initial-cluster-state new
```

初始化集群：

生成TOKEN：

```shell
ETCD_TOKEN=${uuidgen}
```

在每一个node上执行：

```shell
$ HOST_IP=`hostname -i`

$ etcd --name <NodeID> --initial-advertise-peer-urls http://$HOST_IP:2380 \
  --listen-peer-urls http://$HOST_IP:2380 \
  --listen-client-urls http://$HOST_IP:2379,http://127.0.0.1:2379 \
  --advertise-client-urls http://$HOST_IP:2380 \
  --initial-cluster-token $ETCD_TOKEN \ 
  --initial-cluster node1=http://10.0.2.1:2380,node2=http://10.0.2.2:2380,node3=http://10.0.2.3:2380 \
  --initial-cluster-state new
```

注：所有以 --initial-cluster* 开头的选项，在第一次运行（Bootstrap）后都被忽略。

#### 2、公开的发现服务创建集群

首先我们需要创建一个集群注册token

```shell
$ curl https://discovery.etcd.io/new?size=3
https://discovery.etcd.io/a81b5818e67a6ea83e9d4daea5ecbc92
```

初始化集群：在每一个node上执行

```shell
$ DISCOVERY=`curl https://discovery.etcd.io/new?size=3`
$ HOST_IP=`hostname -i`

$ etcd --name <NodeID> --initial-advertise-peer-urls http://$HOST_IP:2380 \
  --listen-peer-urls http://$HOST_IP:2380 \
  --listen-client-urls http://$HOST_IP:2379,http://127.0.0.1:2379 \
  --advertise-client-urls http://$HOST_IP:2380 \
  --discovery ${DISCOVERY} \
  --initial-cluster-token $ETCD_TOKEN \
  --initial-cluster-state new
```



## 使用TLS加密通信

etcd 支持基于 TLS 加密的集群内部、集群外部（客户端与集群之间）的安全通信，每个集群节点都应该拥有被共享 CA 签名的证书：

**生成密钥对、证书签名**

```shell
openssl genrsa -out radon.key 2048
export SAN_CFG=$(printf "\n[SAN]\nsubjectAltName=IP:127.0.0.1,IP:10.0.2.1,DNS:radon.gmem.cc")

openssl req -new -sha256 -key radon.key -out radon.csr \
    -subj "/C=CN/ST=BeiJing/O=Gmem Studio/CN=Server Radon" \
    -reqexts SAN -config <(cat /etc/ssl/openssl.cnf <(echo $SAN_CFG))
```

**执行签名**

```shell
openssl x509 -req -sha256 -in radon.csr  -out radon.crt -CA ../ca.crt -CAkey ../ca.key -CAcreateserial -days 3650 \
     -extensions SAN -extfile <(echo "${SAN_CFG}")
```

**初始化集群命令需要修改为：**

```shell
etcd --name radon --initial-advertise-peer-urls https://10.0.2.1:2380
  --listen-peer-urls https://10.0.2.1:2380
  --listen-client-urls https://10.0.2.1:2379,https://127.0.0.1:2379
  --advertise-client-urls https://10.0.2.1:2380
  --initial-cluster-token etcd.gmem.cc
  --initial-cluster node1=http://10.0.2.1:2380,node2=http://10.0.2.2:2380,node3=http://10.0.2.3:2380      # 指定集群成员列表
  --initial-cluster-state new                                                                              # 初始化新集群时使用  
  --initial-cluster-state existing                                                                        # 加入已有集群时使用 
  
  --client-cert-auth  # 客户端 TLS 相关参数
  --trusted-ca-file=/usr/share/ca-certificates/GmemCA.crt
  --cert-file=/opt/etcd/cert/radon.crt
  --key-file=/opt/etcd/cert/radon.key
  
  --peer-client-cert-auth # 集群内部 TLS 相关参数
  --peer-trusted-ca-file=/usr/share/ca-certificates/GmemCA.crt
  --peer-cert-file=/opt/etcd/cert/radon.crt
  --peer-key-file=/opt/etcd/cert/radon.key
```

**集群健康检测**

```shell
ETCDCTL_API=3 /k8s/etcd/bin/etcdctl  endpoint health \
  --write-out=table \
  --cacert=/k8s/kubernetes/ssl/ca.pem \
  --cert=/k8s/kubernetes/ssl/server.pem \
  --key=/k8s/kubernetes/ssl/server-key.pem \
  --endpoints=https://192.168.0.108:2379,https://192.168.0.109:2379,https://192.168.0.110:2379
```

## 参考资料

快速安装指导： http://play.etcd.io/install#top-top

基本操作api: https://github.com/coreos/etcd/blob/master/Documentation/v2/api.md

集群配置api: https://github.com/coreos/etcd/blob/master/Documentation/v2/members_api.md

鉴权认证api: https://github.com/coreos/etcd/blob/master/Documentation/v2/auth_api.md

配置项：https://github.com/coreos/etcd/blob/master/Documentation/op-guide/configuration.md

https://coreos.com/etcd/docs/latest/runtime-configuration.html
https://coreos.com/etcd/docs/latest/clustering.html
https://coreos.com/etcd/docs/latest/runtime-configuration.html
https://coreos.com/etcd/docs/latest/
https://coreos.com/etcd/docs/latest/admin_guide.html#disaster-recovery

