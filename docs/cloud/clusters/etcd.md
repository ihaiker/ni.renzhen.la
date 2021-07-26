# ETCD集群安装

etcd 目前默认使用 2379 端口提供 HTTP API 服务，2380 端口提供 Peer 通信（这两个端口已经被 IANA 官方预留给 etcd），在之前的版本中，可能会分别使用 4001 和 7001，在使用的过程中需要注意这个区别。

虽然 etcd 也支持单点部署，但是在生产环境中推荐集群方式部署，一般 etcd 节点数会选择 3、5、7。etcd 会保证所有的节点都会保存数据，并保证数据的一致性和正确性。

单点部署
编译部署
因为 etcd 是 Golang 编写的，安装只需要下载对应的二进制文件，并放到合适的路径就行。如果在测试环境，启动一个单点的 etcd 服务，只需要运行 etcd 执行即可。

git clone https://github.com/etcd-io/etcd.git
cd etcd
./build

使用 build 脚本构建会在当前项目的 bin 目录生产 etcd 和 etcdctl 可执行程序。etcd 就是 etcd Server，而 etcdctl 主要为 etcd Server 提供指令行操作。

查看版本：

$ ./bin/etcd --version
etcd Version: 3.5.0-pre
Git SHA: ab4cc3cae
Go Version: go1.14.4
Go OS/Arch: darwin/amd64

$ ./bin/etcdctl version
etcdctl version: 3.5.0-pre
API version: 3.5

启动 etcd Server：

$ ./bin/etcd
{"level":"info","ts":"2020-10-04T07:39:14.751+0800","caller":"etcdmain/etcd.go:69","msg":"Running: ","args":["./bin/etcd"]}
{"level":"info","ts":"2020-10-04T07:39:14.751+0800","caller":"etcdmain/etcd.go:94","msg":"failed to detect default host","error":"default host not supported on darwin_amd64"}
{"level":"warn","ts":"2020-10-04T07:39:14.751+0800","caller":"etcdmain/etcd.go:99","msg":"'data-dir' was empty; using default","data-dir":"default.etcd"}
{"level":"info","ts":"2020-10-04T07:39:14.751+0800","caller":"embed/etcd.go:113","msg":"configuring peer listeners","listen-peer-urls":["http://localhost:2380"]}
{"level":"info","ts":"2020-10-04T07:39:14.752+0800","caller":"embed/etcd.go:121","msg":"configuring client listeners","listen-client-urls":["http://localhost:2379"]}
{"level":"info","ts":"2020-10-04T07:39:14.753+0800","caller":"embed/etcd.go:266","msg":"starting an etcd server","etcd-version":"3.5.0-pre","git-sha":"ab4cc3cae","go-version":"go1.14.4","go-os":"darwin","go-arch":"amd64","max-cpu-set":4,"max-cpu-available":4,"member-initialized":false,"name":"default","data-dir":"default.etcd","wal-dir":"","wal-dir-dedicated":"","member-dir":"default.etcd/member","force-new-cluster":false,"heartbeat-interval":"100ms","election-timeout":"1s","initial-election-tick-advance":true,"snapshot-count":100000,"snapshot-catchup-entries":5000,"initial-advertise-peer-urls":["http://localhost:2380"],"listen-peer-urls":["http://localhost:2380"],"advertise-client-urls":["http://localhost:2379"],"listen-client-urls":["http://localhost:2379"],"listen-metrics-urls":[],"cors":["*"],"host-whitelist":["*"],"initial-cluster":"default=http://localhost:2380","initial-cluster-state":"new","initial-cluster-token":"etcd-cluster","quota-size-bytes":2147483648,"pre-vote":false,"initial-corrupt-check":false,"corrupt-check-time-interval":"0s","auto-compaction-mode":"periodic","auto-compaction-retention":"0s","auto-compaction-interval":"0s","discovery-url":"","discovery-proxy":""}
{"level":"info","ts":"2020-10-04T07:39:14.764+0800","caller":"etcdserver/backend.go:78","msg":"opened backend db","path":"default.etcd/member/snap/db","took":"9.908726ms"}
{"level":"info","ts":"2020-10-04T07:39:14.862+0800","caller":"etcdserver/raft.go:444","msg":"starting local member","local-member-id":"8e9e05c52164694d","cluster-id":"cdf818194e3a8c32"}
{"level":"info","ts":"2020-10-04T07:39:14.862+0800","caller":"raft/raft.go:1528","msg":"8e9e05c52164694d switched to configuration voters=()"}
{"level":"info","ts":"2020-10-04T07:39:14.862+0800","caller":"raft/raft.go:701","msg":"8e9e05c52164694d became follower at term 0"}
{"level":"info","ts":"2020-10-04T07:39:14.862+0800","caller":"raft/raft.go:383","msg":"newRaft 8e9e05c52164694d [peers: [], term: 0, commit: 0, applied: 0, lastindex: 0, lastterm: 0]"}
{"level":"info","ts":"2020-10-04T07:39:14.863+0800","caller":"raft/raft.go:701","msg":"8e9e05c52164694d became follower at term 1"}
{"level":"info","ts":"2020-10-04T07:39:14.863+0800","caller":"raft/raft.go:1528","msg":"8e9e05c52164694d switched to configuration voters=(10276657743932975437)"}
{"level":"warn","ts":"2020-10-04T07:39:14.888+0800","caller":"auth/store.go:1231","msg":"simple token is not cryptographically signed"}
{"level":"info","ts":"2020-10-04T07:39:14.912+0800","caller":"etcdserver/quota.go:94","msg":"enabled backend quota with default value","quota-name":"v3-applier","quota-size-bytes":2147483648,"quota-size":"2.1 GB"}
{"level":"info","ts":"2020-10-04T07:39:14.924+0800","caller":"etcdserver/server.go:752","msg":"starting etcd server","local-member-id":"8e9e05c52164694d","local-server-version":"3.5.0-pre","cluster-version":"to_be_decided"}
{"level":"info","ts":"2020-10-04T07:39:14.925+0800","caller":"etcdserver/server.go:640","msg":"started as single-node; fast-forwarding election ticks","local-member-id":"8e9e05c52164694d","forward-ticks":9,"forward-duration":"900ms","election-ticks":10,"election-timeout":"1s"}
{"level":"warn","ts":"2020-10-04T07:39:14.925+0800","caller":"etcdserver/metrics.go:212","msg":"failed to get file descriptor usage","error":"cannot get FDUsage on darwin"}
{"level":"info","ts":"2020-10-04T07:39:14.925+0800","caller":"raft/raft.go:1528","msg":"8e9e05c52164694d switched to configuration voters=(10276657743932975437)"}
{"level":"info","ts":"2020-10-04T07:39:14.925+0800","caller":"membership/cluster.go:385","msg":"added member","cluster-id":"cdf818194e3a8c32","local-member-id":"8e9e05c52164694d","added-peer-id":"8e9e05c52164694d","added-peer-peer-urls":["http://localhost:2380"]}
{"level":"info","ts":"2020-10-04T07:39:14.927+0800","caller":"embed/etcd.go:513","msg":"serving peer traffic","address":"127.0.0.1:2380"}
{"level":"info","ts":"2020-10-04T07:39:14.927+0800","caller":"embed/etcd.go:235","msg":"now serving peer/client/metrics","local-member-id":"8e9e05c52164694d","initial-advertise-peer-urls":["http://localhost:2380"],"listen-peer-urls":["http://localhost:2380"],"advertise-client-urls":["http://localhost:2379"],"listen-client-urls":["http://localhost:2379"],"listen-metrics-urls":[]}
{"level":"info","ts":"2020-10-04T07:39:15.866+0800","caller":"raft/raft.go:788","msg":"8e9e05c52164694d is starting a new election at term 1"}
{"level":"info","ts":"2020-10-04T07:39:15.866+0800","caller":"raft/raft.go:714","msg":"8e9e05c52164694d became candidate at term 2"}
{"level":"info","ts":"2020-10-04T07:39:15.866+0800","caller":"raft/raft.go:848","msg":"8e9e05c52164694d received MsgVoteResp from 8e9e05c52164694d at term 2"}
{"level":"info","ts":"2020-10-04T07:39:15.867+0800","caller":"raft/raft.go:766","msg":"8e9e05c52164694d became leader at term 2"}
{"level":"info","ts":"2020-10-04T07:39:15.867+0800","caller":"raft/node.go:327","msg":"raft.node: 8e9e05c52164694d elected leader 8e9e05c52164694d at term 2"}
{"level":"info","ts":"2020-10-04T07:39:15.868+0800","caller":"etcdserver/server.go:2285","msg":"setting up initial cluster version","cluster-version":"3.5"}
{"level":"info","ts":"2020-10-04T07:39:15.876+0800","caller":"membership/cluster.go:523","msg":"set initial cluster version","cluster-id":"cdf818194e3a8c32","local-member-id":"8e9e05c52164694d","cluster-version":"3.5"}
{"level":"info","ts":"2020-10-04T07:39:15.876+0800","caller":"embed/serve.go:97","msg":"ready to serve client requests"}
{"level":"info","ts":"2020-10-04T07:39:15.876+0800","caller":"api/capability.go:75","msg":"enabled capabilities for version","cluster-version":"3.5"}
{"level":"info","ts":"2020-10-04T07:39:15.876+0800","caller":"etcdserver/server.go:2305","msg":"cluster version is updated","cluster-version":"3.5"}
{"level":"info","ts":"2020-10-04T07:39:15.876+0800","caller":"etcdserver/server.go:1863","msg":"published local member to cluster through raft","local-member-id":"8e9e05c52164694d","local-member-attributes":"{Name:default ClientURLs:[http://localhost:2379]}","request-path":"/0/members/8e9e05c52164694d/attributes","cluster-id":"cdf818194e3a8c32","publish-timeout":"7s"}
{"level":"info","ts":"2020-10-04T07:39:15.876+0800","caller":"etcdmain/main.go:47","msg":"notifying init daemon"}
{"level":"info","ts":"2020-10-04T07:39:15.876+0800","caller":"etcdmain/main.go:53","msg":"successfully notified init daemon"}
{"level":"info","ts":"2020-10-04T07:39:15.877+0800","caller":"embed/serve.go:139","msg":"serving client traffic insecurely; this is strongly discouraged!","address":"127.0.0.1:2379"}

name 表示节点名称，默认为 default。
data-dir 表示 WAL 日志和 Snapshot 数据储存目录，默认为 ./default.etcd/ 目录。
使用 http://localhost:2380 和 etcd Cluster 中其他节点通信。
使用 http://localhost:2379 提供 HTTP API 服务，与客户端通信。
heartbeat 为 100ms，表示 Leader 多久发送一次心跳到所有 Followers。
election-timeout 为 1s，该参数的作用是重新投票的超时时间，如果 Follow 在该时间间隔内没有收到 Leader 发出的心跳包，就会触发重新投票。
snapshot-count 为 100000，该参数的作用是指定有多少次事务被提交后触发快照截取动作并持久化到磁盘。
cluster-id 为 cdf818194e3a8c32。
raft.node 为 8e9e05c52164694d。
启动的时候，会运行 Raft，选举出 Leader：elected leader 8e9e05c52164694d at term 2。
上述方法只是简单的启动了一个 etcd Server。当然，在生产环境中，通常使用 Systemd 来进行管理。

建立相关目录：
$ mkdir -p /var/lib/etcd/
$ mkdir -p /etc/etcd/config/
设定 etcd 配置文件：
$ cat <<EOF | sudo tee /etc/etcd/config/etcd.conf

# 节点名称
ETCD_NAME=$(hostname -s)
# 数据存放路径
ETCD_DATA_DIR=/var/lib/etcd

创建 systemd 配置文件：
$ cat <<EOF | sudo tee /etc/systemd/system/etcd.service

[Unit]
Description=Etcd Server
Documentation=https://github.com/coreos/etcd
After=network.target

[Service]
User=root
Type=notify
EnvironmentFile=-/opt/etcd/config/etcd.conf
ExecStart=~/workspace/etcd/bin
Restart=on-failure
RestartSec=10s
LimitNOFILE=40000

[Install]
WantedBy=multi-user.target

启动 etcd Server：
$ systemctl daemon-reload && systemctl enable etcd && systemctl start etcd

关键启动选型清单：

–name：指定 etcd Node 名称，可以使用 hostname。
–data-dir：指定 etcd Server 持久化数据存储目录路径。
–snapshot-count：指定有多少事务（transaction）被提交后，触发截取快照并持久化到磁盘。
–heartbeat-interval：指定 Leader 多久发送一次心跳到 Followers。
–eletion-timeout：指定重新投票的超时时间，如果 Follow 在该时间间隔没有收到 Leader 发出的心跳包，则会触发重新投票。
–listen-peer-urls：指定和 Cluster 其他 Node 通信的地址，比如：http://IP:2380，如果有多个，则使用逗号分隔。需要所有节点都能够访问，所以不要使用 localhost。
–listen-client-urls：指定对外提供服务的地址，比如：http://IP:2379,http://127.0.0.1:2379。
–advertise-client-urls：对外通告的该节点的客户端监听地址，会告诉集群中其他节点。
–initial-advertise-peer-urls：对外通告该节点的同伴（Peer）监听地址，这个值会告诉集群中其他节点。
–initial-cluster：指定集群中所有节点的信息，通常为 IP:Port 信息，格式为：node1=http://ip1:2380,node2=http://ip2:2380,…。注意，这里的 node1 就是 --name 指定的名字，ip1:2380 就是 --initial-advertise-peer-urls 指定的值。
–initial-cluster-state：新建集群时，这个值为 new；假如已经存在了集群，这个值为 existing。
–initial-cluster-token：创建集群的 token，这个值每个集群保持唯一。这样的话，如果你要重新创建集群，即使配置和之前一样，也会再次生成新的集群和节点 UUID；否则会导致多个集群之间的冲突，造成未知的错误。
容器部署
docker run on console
$ docker run --rm --name etcd -p 2379:2379 -e ETCDCTL_API=3 quay.io/coreos/etcd /usr/local/bin/etcd -advertise-client-urls http://0.0.0.0:2379 -listen-client-urls http://0.0.0.0:2379
$ docker ps -f name=etcd

etcdctl
$ docker exec -it etcd etcdctl version
etcdctl version: 3.3.8
API version: 3.3

$ docker exec -it etcd etcdctl endpoint health
127.0.0.1:2379 is healthy: successfully committed proposal: took = 571.723µs

集群部署
etcd Cluster 的部署，实际上就是多个主机上都部署 etcd Server，然后将它们加入到一个 Cluster 中。

注意：etcd Cluster 必须具有时间同步服务器，否则会导致 Cluster 频繁进行 Leader Change。在 OpenShift 的 etcd Cluster 中，会每隔 100ms 进行心跳检测。

在安装和启动 etcd 服务进程的时候，各个 Node 都需要知道 Cluster 中其他 Nodes 的信息，一般是 IP:Port 信息。根据用户是否提前知晓（规划）了每个 Node 的 IP 地址，有以下几种不同的集群部署方案：

静态配置：在启动 etcd Server 的时候，通过 --initial-cluster 参数配置好所有的节点信息。
注册到已有的 etcd Cluster：比如官方提供的 discovery.etcd.io。
使用 DNS 启动。
静态配置集群
如果 etcd Cluster 中的成员是已知的，且具有固定的 IP 地址，就可以静态的初始化一个集群。

每个 Node 都使用如下环境变量：

ETCD_INITIAL_CLUSTER="radon=http://10.0.2.1:2380,neon=http://10.0.3.1:2380"
ETCD_INITIAL_CLUSTER_STATE=new

或者使用如下指令行参数来指定集群成员：

--initial-cluster radon=http://10.0.2.1:2380,neon=http://10.0.3.1:2380
--initial-cluster-state new

初始化集群：

etcd --name radon --initial-advertise-peer-urls http://10.0.2.1:2380
  --listen-peer-urls http://10.0.2.1:2380
  --listen-client-urls http://10.0.2.1:2379,http://127.0.0.1:2379
  --advertise-client-urls http://10.0.2.1:2380
  --initial-cluster-token etcd.gmem.cc
  --initial-cluster radon=http://10.0.2.1:2380,neon=http://10.0.3.1:2380
  --initial-cluster-state new

注：所有以 --initial-cluster* 开头的选项，在第一次运行（Bootstrap）后都被忽略。

使用 TLS 加密，etcd 支持基于 TLS 加密的集群内部、集群外部（客户端与集群之间）的安全通信，每个集群节点都应该拥有被共享 CA 签名的证书：

# 密钥对、证书签名请求
openssl genrsa -out radon.key 2048
export SAN_CFG=$(printf "\n[SAN]\nsubjectAltName=IP:127.0.0.1,IP:10.0.2.1,DNS:radon.gmem.cc")
openssl req -new -sha256 -key radon.key -out radon.csr \
    -subj "/C=CN/ST=BeiJing/O=Gmem Studio/CN=Server Radon" \
    -reqexts SAN -config <(cat /etc/ssl/openssl.cnf <(echo $SAN_CFG))

# 执行签名
openssl x509 -req -sha256 -in radon.csr  -out radon.crt -CA ../ca.crt -CAkey ../ca.key -CAcreateserial -days 3650 \
     -extensions SAN -extfile <(echo "${SAN_CFG}")

初始化集群命令需要修改为：

etcd --name radon --initial-advertise-peer-urls https://10.0.2.1:2380
  --listen-peer-urls https://10.0.2.1:2380
  --listen-client-urls https://10.0.2.1:2379,https://127.0.0.1:2379
  --advertise-client-urls https://10.0.2.1:2380
  --initial-cluster-token etcd.gmem.cc
  --initial-cluster radon=https://10.0.2.1:2380,neon=https://10.0.3.1:2380      # 指定集群成员列表
  --initial-cluster-state new                                                                              # 初始化新集群时使用  
  --initial-cluster-state existing                                                                        # 加入已有集群时使用 

  # 客户端 TLS 相关参数
  --client-cert-auth 
  --trusted-ca-file=/usr/share/ca-certificates/GmemCA.crt
  --cert-file=/opt/etcd/cert/radon.crt
  --key-file=/opt/etcd/cert/radon.key

  # 集群内部 TLS 相关参数
  --peer-client-cert-auth
  --peer-trusted-ca-file=/usr/share/ca-certificates/GmemCA.crt
  --peer-cert-file=/opt/etcd/cert/radon.crt
  --peer-key-file=/opt/etcd/cert/radon.key

集群健康检测
ETCDCTL_API=3 /k8s/etcd/bin/etcdctl  endpoint health --write-out=table \
  --cacert=/k8s/kubernetes/ssl/ca.pem \
  --cert=/k8s/kubernetes/ssl/server.pem \
  --key=/k8s/kubernetes/ssl/server-key.pem \
  --endpoints=https://192.168.0.108:2379,https://192.168.0.109:2379,https://192.168.0.110:2379
------------------------------------------------
版权声明：本文为CSDN博主「范桂飓」的原创文章，遵循CC 4.0 BY-SA版权协议，转载请附上原文出处链接及本声明。
原文链接：https://blog.csdn.net/Jmilk/article/details/108914220
