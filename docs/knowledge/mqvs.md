---
title: "RocketMQ vs RabbitMQ vs Kafka"
---



# 三大主流MQ的组织结构

## 1、RabbitMQ

### RabbitMQ各组件的功能

![image.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/f43cadbf87e54ac8b821b1695fc3888f~tplv-k3u1fbpfcp-watermark.awebp)

-   **`Broker：`** **一个RabbitMQ实例**就是一个Broker
-   **`Virtual Host：`** 虚拟主机。**相当于Mysql的DataBase**, 一个Broker上可以存在多个vhost，vhost之间相互隔离。每个vhost都拥有自己的队列、交换机、绑定和权限机制。vhost必须在连接时指定，默认的vhost是 /。
-   **`Exchange：`** 交换机，用来接收生产者发送的消息并将这些消息路由给服务器中的队列。
-   **`Queue：`** 消息队列，用来保存消息直到发送给消费者。**它是消息的容器**。一个消息可投入一个或多个队列。
-   **`Banding：`** 绑定关系，用于**消息队列和交换机之间的关联**。通过路由键(**Routing Key**)将交换机和消息队列关联起来。
-   **`Channel：`** 管道，一条**双向数据流通道**。不管是发布消息、订阅队列还是接收消息，这些动作都是通过管道完成。因为对于操作系统来说，建立和销毁TCP都是非常昂贵的开销，所以引入了管道的概念，以复用一条TCP连接。
-   **`Connection：`** 生产者/消费者 与broker之间的TCP连接。
-   **`Publisher：`** 消息的生产者。
-   **`Consumer：`** 消息的消费者。
-   **`Message：`** 消息，它是由消息头和消息体组成。消息头则包括**Routing-Key**、**Priority**(优先级)等。

------

### RabbitMQ的多种交换机类型

`Exchange`分发消息给`Queue`时，`Exchange`的类型对应不同的分发策略，有3种类型的`Exchange`：**Direct**、**Fanout**、**Topic**。

-   **Direct**： 消息中的`Routing Key`如果和`Binding`中的`Routing Key`完全一致，`Exchange`就会将消息分发到对应的队列中。

![image.png](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/725c2fb9a52143aa9802394e3d5360ac~tplv-k3u1fbpfcp-watermark.awebp)

-   **Fanout**： 每个发到 Fanout 类型交换机的消息都会分发到所有绑定的队列上去。Fanout交换机没有`Routing Key`。**它在三种类型的交换机中转发消息是最快的**。

![image.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/5666e34fa88d410b88f0847ec69a79cc~tplv-k3u1fbpfcp-watermark.awebp)

-   **Topic**： Topic交换机通过模式匹配分配消息，将`Routing Key`和某个模式进行匹配。它只能识别两个**通配符**：**"#"和"\*"**。`#`匹配0个或多个单词，`*`匹配1个单词。

![image.png](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/1ebed3f0495b41c8b11c13a6115ad183~tplv-k3u1fbpfcp-watermark.awebp)

------

### TTL

TTL(Time To Live)：生存时间。RabbitMQ支持消息的过期时间，一共2种。

-   **在消息发送时进行指定**。通过配置消息体的`Properties`，可以指定当前消息的过期时间。
-   **在创建Exchange时指定**。从进入消息队列开始计算，只要超过了队列的超时时间配置，那么消息会自动清除。

------

### 生产者的消息确认机制

#### Confirm机制

-   消息的确认，是指生产者投递消息后，如果Broker收到消息，则会给我们生产者一个应答。
-   生产者进行接受应答，用来确认这条消息是否正常的发送到了Broker，这种方式也是**消息的可靠性投递的核心保障！**

>   如何实现Confirm确认消息？

![image.png](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/ad58d573fd0c47acac1e85d683ff4be1~tplv-k3u1fbpfcp-watermark.awebp)

1.  **在channel上开启确认模式**：`channel.confirmSelect()`
2.  **在channel上开启监听**：`addConfirmListener`，监听成功和失败的处理结果，根据具体的结果对消息进行重新发送或记录日志处理等后续操作。

#### Return消息机制

Return Listener**用于处理一些不可路由的消息**。

我们的消息生产者，通过指定一个Exchange和Routing，把消息送达到某一个队列中去，然后我们的消费者监听队列进行消息的消费处理操作。

但是在某些情况下，如果我们在发送消息的时候，当前的exchange不存在或者指定的路由key路由不到，这个时候我们需要监听这种不可达消息，就需要使用到Returrn Listener。

基础API中有个关键的配置项`Mandatory`：如果为true，监听器会收到路由不可达的消息，然后进行处理。如果为false，broker端会自动删除该消息。

同样，通过监听的方式，`chennel.addReturnListener(ReturnListener rl)`传入已经重写过handleReturn方法的ReturnListener。

------

### 消费端ACK与NACK

消费端进行消费的时候，如果由于业务异常可以进行日志的记录，然后进行补偿。但是对于服务器宕机等严重问题，我们需要**手动ACK**保障消费端消费成功。

```js
// deliveryTag：消息在mq中的唯一标识
// multiple：是否批量(和qos设置类似的参数)
// requeue：是否需要重回队列。或者丢弃或者重回队首再次消费。
public void basicNack(long deliveryTag, boolean multiple, boolean requeue) 
复制代码
```

如上代码，消息在**消费端重回队列**是为了对没有成功处理消息，把消息重新返回到Broker。一般来说，实际应用中都会关闭重回队列(**避免进入死循环**)，也就是设置为false。

------

### 死信队列DLX

**`死信队列(DLX Dead-Letter-Exchange)：`** 当消息在一个队列中变成死信之后，它会被重新推送到另一个队列，这个队列就是死信队列。

DLX也是一个正常的Exchange，和一般的Exchange没有区别，它能在任何的队列上被指定，实际上就是设置某个队列的属性。

当这个队列中有死信时，RabbitMQ就会自动的将这个消息重新发布到设置的Exchange上去，进而被路由到另一个队列。

## 2、RocketMQ

阿里巴巴双十一官方指定消息产品，支撑阿里巴巴集团所有的消息服务，历经十余年高可用与高可靠的严苛考验，是阿里巴巴交易链路的核心产品；

Rocket：火箭的意思。

![image.png](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/b7b1cb200a0f4530830c030dd42f2007~tplv-k3u1fbpfcp-watermark.awebp)

### RocketMQ的核心概念

他有以下核心概念：`Broker`、`Topic`、`Tag`、`MessageQueue`、`NameServer`、`Group`、`Offset`、`Producer`以及`Consumer`。

下面来详细介绍。

-   **Broker**：消息中转角色，负责**存储消息**，转发消息。
    -   -   **Broker**是具体提供业务的服务器，单个Broker节点与所有的NameServer节点保持长连接及心跳，并会定时将**Topic**信息注册到NameServer，顺带一提底层的通信和连接都是**基于Netty实现**的。
    -   -   **Broker**负责消息存储，以Topic为纬度支持轻量级的队列，单机可以支撑上万队列规模，支持消息推拉模型。
    -   -   官网上有数据显示：具有**上亿级消息堆积能力**，同时可**严格保证消息的有序性**。
-   **Topic**：主题！它是消息的第一级类型。比如一个电商系统可以分为：交易消息、物流消息等，一条消息必须有一个 Topic 。**Topic** 与生产者和消费者的关系非常松散，一个 Topic 可以有0个、1个、多个生产者向其发送消息，一个生产者也可以同时向不同的 Topic 发送消息。一个 Topic 也可以被 0个、1个、多个消费者订阅。
-   **Tag**：标签！可以看作子主题，它是消息的第二级类型，用于为用户提供额外的灵活性。使用标签，同一业务模块不同目的的消息就可以用相同 Topic 而不同的 **Tag** 来标识。比如交易消息又可以分为：交易创建消息、交易完成消息等，一条消息可以没有 **Tag** 。标签有助于保持您的代码干净和连贯，并且还可以为 **RocketMQ** 提供的查询系统提供帮助。
-   **MessageQueue**：一个 Topic 下可以设置多个消息队列，发送消息时执行该消息的 Topic ，RocketMQ 会轮询该 Topic 下的所有队列将消息发出去。消息的物理管理单位。一个Topic下可以有多个Queue，Queue的引入使得消息的存储可以分布式集群化，具有了水平扩展能力。
-   **NameServer**：类似Kafka中的Zookeeper, 但NameServer集群之间是**没有通信**的，相对ZK来说更加**轻量**。 它主要负责对于源数据的管理，包括了对于**Topic**和路由信息的管理。每个 Broker 在启动的时候会到 NameServer 注册，Producer 在发送消息前会根据 Topic 去 NameServer **获取对应 Broker 的路由信息**，Consumer 也会定时获取 Topic 的路由信息。
-   **Producer**： 生产者，支持三种方式发送消息：**同步、异步和单向**
    -   -   `单向发送`：消息发出去后，可以继续发送下一条消息或执行业务代码，不等待服务器回应，且**没有回调函数**。
    -   -   `异步发送`：消息发出去后，可以继续发送下一条消息或执行业务代码，不等待服务器回应，**有回调函数**。
    -   -   `同步发送`：消息发出去后，等待服务器响应成功或失败，才能继续后面的操作。
-   **Consumer**：消费者，支持`PUSH`和`PULL`两种消费模式，支持**集群消费**和**广播消费**
    -   -   `集群消费`: 该模式下一个消费者集群共同消费一个主题的多个队列，一个队列只会被一个消费者消费，如果某个消费者挂掉，分组内其它消费者会接替挂掉的消费者继续消费。
    -   -   `广播消费`: 会发给消费者组中的每一个消费者进行消费。相当于**RabbitMQ**的发布订阅模式。
-   **Group**：分组，一个组可以订阅多个Topic。分为 ProducerGroup，ConsumerGroup，代表某一类的生产者和消费者，一般来说同一个服务可以作为Group，同一个Group一般来说发送和消费的消息都是一样的
-   **Offset**：在**RocketMQ** 中，所有消息队列都是持久化，长度无限的数据结构，所谓长度无限是指队列中的每个存储单元都是定长，访问其中的存储单元使用Offset 来访问，Offset 为 java long 类型，64 位，理论上在 100年内不会溢出，所以认为是长度无限。也可以认为 Message Queue 是一个长度无限的数组，**Offset** 就是下标。

------

### 延时消息

开源版的RocketMQ不支持任意时间精度，仅支持特定的 level，例如定时 5s， 10s， 1min 等。其中，level=0 级表示不延时，level=1 表示 1 级延时，level=2 表示 2 级延时，以此类推。

延时等级如下：

```
messageDelayLevel=1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h
复制代码
```

### 顺序消息

消息有序指的是可以按照消息的发送顺序来消费（FIFO）。RocketMQ可以严格的保证消息有序，可以分为`分区有序`或者`全局有序`。

### 事务消息

![image.png](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/be59a61110fd4c2eb41c7db4a6b32021~tplv-k3u1fbpfcp-watermark.awebp) 消息队列 MQ 提供类似 X/Open XA 的分布式事务功能，通过消息队列 MQ 事务消息能达到分布式事务的最终一致。上图说明了事务消息的大致流程：正常事务消息的发送和提交、事务消息的补偿流程。

-   **事务消息发送及提交**：

1.  发送half消息
2.  服务端响应消息写入结果
3.  根据发送结果执行本地事务（如果写入失败，此时half消息对业务不可见，本地逻辑不执行）；
4.  根据本地事务状态执行Commit或Rollback（Commit操作生成消息索引，消息对消费者可见）。

-   **事务消息的补偿流程**：

1.  对没有Commit/Rollback的事务消息（pending状态的消息），从服务端发起一次“回查”；
2.  Producer收到回查消息，检查回查消息对应的本地事务的状态。
3.  根据本地事务状态，重新Commit或RollBack

其中，补偿阶段用于解决消息Commit或Rollback发生超时或者失败的情况。

-   **事务消息状态**：

事务消息共有三种状态：提交状态、回滚状态、中间状态：

1.  TransactionStatus.CommitTransaction：提交事务，它允许消费者消费此消息。
2.  TransactionStatus.RollbackTransaction：回滚事务，它代表该消息将被删除，不允许被消费。
3.  TransactionStatus.Unkonwn：中间状态，它代表需要检查消息队列来确定消息状态。

### RocketMQ的高可用机制

RocketMQ是天生支持分布式的，可以配置主从以及水平扩展

Master角色的Broker支持读和写，Slave角色的Broker仅支持读，也就是 Producer只能和Master角色的Broker连接写入消息；Consumer可以连接 Master角色的Broker，也可以连接Slave角色的Broker来读取消息。

#### 消息消费的高可用（主从）

在Consumer的配置文件中，并不需要设置是从Master读还是从Slave 读，当Master不可用或者繁忙的时候，Consumer会被自动切换到从Slave 读。有了自动切换Consumer这种机制，当一个Master角色的机器出现故障后，Consumer仍然可以从Slave读取消息，不影响Consumer程序。这就达到了消费端的高可用性。 **RocketMQ目前还不支持把Slave自动转成Master**，如果机器资源不足，需要把Slave转成Master，则要手动停止Slave角色的Broker，更改配置文件，用新的配置文件启动Broker。

#### 消息发送高可用（配置多个主节点）

在创建Topic的时候，把Topic的多个Message Queue创建在多个Broker组上（相同Broker名称，不同 brokerId的机器组成一个Broker组），这样当一个Broker组的Master不可用后，其他组的Master仍然可用，Producer仍然可以发送消息。

#### 主从复制

如果一个Broker组有Master和Slave，消息需要从Master复制到Slave 上，有同步和异步两种复制方式。

-   **同步复制**：同步复制方式是等Master和Slave均写成功后才反馈给客户端写成功状态。如果Master出故障， Slave上有全部的备份数据，容易恢复同步复制会增大数据写入延迟，降低系统吞吐量。
-   **异步复制**：异步复制方式是只要Master写成功 即可反馈给客户端写成功状态。在异步复制方式下，系统拥有较低的延迟和较高的吞吐量，但是如果Master出了故障，有些数据因为没有被写 入Slave，有可能会丢失

```
通常情况下，应该把Master和Save配置成同步刷盘方式，主从之间配置成异步的复制方式，这样即使有一台机器出故障，仍然能保证数据不丢，是个不错的选择。
```

------

### 负载均衡

####  Producer负载均衡

Producer端，每个实例在发消息的时候，默认会**轮询**所有的message queue发送，以达到让消息平均落在不同的queue上。而由于queue可以散落在不同的broker，所以消息就发送到不同的broker下，如下图：

![image.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/1ecb662186f24121bfe5ba1ccffaf750~tplv-k3u1fbpfcp-watermark.awebp)

#### Consumer负载均衡

>   如果consumer实例的数量比message queue的总数量还多的话，**多出来的consumer实例将无法分到queue**，也就无法消费到消息，也就无法起到分摊负载的作用了。所以需要控制让queue的总数量大于等于consumer的数量。

-   消费者的集群模式–启动多个消费者就可以保证消费者的负载均衡（均摊队列）
-   **默认使用的是均摊队列**：会按照queue的数量和实例的数量平均分配queue给每个实例，这样每个消费者可以均摊消费的队列，如下图所示6个队列和三个生产者。

![image.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/938eecc24e78468cac70a876278b1aba~tplv-k3u1fbpfcp-watermark.awebp)

-   另外一种平均的算法**环状轮流分queue**的形式，每个消费者，均摊不同主节点的一个消息队列，如下图所示：

![image.png](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/5ed73ddee6ef406481e5f811507aba33~tplv-k3u1fbpfcp-watermark.awebp)

>   对于广播模式并不是负载均衡的，要求一条消息需要投递到一个消费组下面所有的消费者实例，所以也就没有消息被分摊消费的说法。

------

### 死信队列

当一条消息消费失败，RocketMQ就会自动进行消息重试。而如果消息超过最大重试次数，RocketMQ就会认为这个消息有问题。但是此时，RocketMQ不会立刻将这个有问题的消息丢弃，而会将其发送到这个消费者组对应的一种特殊队列：死信队列。死信队列的名称是`%DLQ%+ConsumGroup`

**死信队列具有以下特性**：

1.  一个死信队列对应一个 Group ID， 而不是对应单个消费者实例。
2.  如果一个 Group ID 未产生死信消息，消息队列 RocketMQ 不会为其创建相应的死信队列。
3.  一个死信队列包含了对应 Group ID 产生的所有死信消息，不论该消息属于哪个 Topic。

## 3、Kafka

Kafka是一个分布式、支持分区的、多副本的，**基于Zookeeper**协调的分布式消息系统。

它最大的特性就是可以实时的处理大量数据以满足各种需求场景：比如基于hadoop的批处理系统、低延迟的实时系统、Storm/Spark流式处理引擎，Web/Nginx日志、访问日志，消息服务等等，用**Scala语言编写**。属于Apache基金会的顶级开源项目。

先看一下Kafka的架构图 ![image.png](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/9a32e13d0c3b4e87954e7a7abaf9eca2~tplv-k3u1fbpfcp-watermark.awebp)

### Kafka的核心概念

在Kafka中有几个核心概念：`Broker`、`Topic`、`Producer`、`Consumer`、`ConsumerGroup`、`Partition`、`Leader`、`Follower`、`Offset`。

-   **Broker** : 消息中间件处理节点，一个Kafka节点就是一个broker，一个或者多个Broker可以组成一个Kafka集群
-   **Topic** : Kafka根据topic对消息进行归类，发布到Kafka集群的每条消息都需要指定一个 topic
-   **Producer** : 消息生产者，向Broker发送消息的客户端
-   **Consumer** : 消息消费者，从Broker读取消息的客户端
-   **ConsumerGroup**: 每个Consumer属于一个特定的ConsumerGroup，一条消息可以被多个不同的ConsumerGroup消费，但是一个ConsumerGroup中只能有一个Consumer能够消费该消息
-   **Partition** : 物理上的概念，一个topic可以分为多个partition，每个partition内部消息是有序的
-   **Leader** : 每个partition有多个副本，其中有且仅有一个作为Leader，Leader是负责数据读写的partition。
-   **Follower** : Follower跟随Leader，所有写请求都通过Leader路由，数据变更会广播给所有Follower，Follower与Leader保持数据同步。如果Leader失效，则从Follower中选举出一个新的Leader。当Follower与Leader挂掉、卡住或者同步太慢，leader会把这个follower从`ISR列表`中删除，重新创建一个Follower。
-   **Offset** : 偏移量。kafka的存储文件都是按照offset.kafka来命名，用offset做名字的好处是方便查找。例如你想找位于2049的位置，只要找到2048.kafka的文件即可

**可以这么来理解Topic，Partition和Broker：**

一个Topic，代表逻辑上的一个业务数据集，比如订单相关操作消息放入订单Topic，用户相关操作消息放入用户Topic，对于大型网站来说，后端数据都是海量的，订单消息很可能是非常巨量的，比如有几百个G甚至达到TB级别，如果把这么多数据都放在一台机器上可定会有容量限制问题，那么就可以在Topic内部划分多个Partition来分片存储数据，不同的Partition可以位于不同的机器上，相当于**分布式存储**。每台机器上都运行一个Kafka的进程Broker。

------

### Kafka核心总控制器Controller

在Kafka集群中会有一个或者多个broker，其中有一个broker会被选举为控制器（Kafka Controller），可以理解为`Broker-Leader`, 它负责管理整个 集群中所有分区和副本的状态。

-   当某个`Partition-Leader`副本出现故障时，由控制器负责为该分区选举新的Leader副本。
-   当检测到某个分区的ISR集合发生变化时，由控制器负责通知所有broker更新其元数据信息。
-   当为某个topic增加分区数量时，同样还是由控制器负责让新分区被其他节点感知到。

------

### Controller选举机制

在kafka集群启动的时候，`选举的过程是集群中每个broker都会 尝试在zookeeper上创建一个 /controller 临时节点，zookeeper会保证有且仅有一个broker能创建成功`，这个broker 就会成为集群的总控器controller。

当这个controller角色的broker宕机了，此时zookeeper临时节点会消失，集群里其他broker会一直监听这个临时节 点，发现临时节点消失了，就竞争再次创建临时节点，就是我们上面说的选举机制，zookeeper又会保证有一个broker 成为新的controller。 具备控制器身份的broker需要比其他普通的broker多一份职责，具体细节如下：

1.  **监听broker相关的变化**。为Zookeeper中的/brokers/ids/节点添加BrokerChangeListener，用来处理broker 增减的变化。
2.  **监听topic相关的变化**。为Zookeeper中的/brokers/topics节点添加TopicChangeListener，用来处理topic增减 的变化；为Zookeeper中的/admin/delete_topics节点添加TopicDeletionListener，用来处理删除topic的动作。
3.  **从Zookeeper中读取获取当前所有与topic、partition以及broker有关的信息并进行相应的管理**。对于所有topic 所对应的Zookeeper中的/brokers/topics/节点添加PartitionModificationsListener，用来监听topic中的分区分配变化。
4.  **更新集群的元数据信息，同步到其他普通的broker节点中**

------

### Partition副本选举Leader机制

controller感知到分区Leader所在的broker挂了，controller会从 ISR列表(参数unclean.leader.election.enable=false的前提下)里挑第一个broker作为leader(第一个broker最先放进ISR 列表，可能是同步数据最多的副本)，如果参数unclean.leader.election.enable为true，代表在ISR列表里所有副本都挂 了的时候可以在ISR列表以外的副本中选leader，这种设置，可以提高可用性，但是选出的新leader有可能数据少很多。 副本进入ISR列表有两个条件：

1.  副本节点不能产生分区，必须能与zookeeper保持会话以及跟leader副本网络连通
2.  副本能复制leader上的所有写操作，并且不能落后太多。(与leader副本同步滞后的副本，是由 replica.lag.time.max.ms 配置决定的，超过这个时间都没有跟leader同步过的一次的副本会被移出ISR列表)

------

### 消费者消费消息的offset记录机制

每个consumer会定期将自己消费分区的offset提交给kafka内部topic：consumer_offsets，提交过去的时候，key是 consumerGroupId+topic+分区号，value就是当前offset的值，kafka会定期清理topic里的消息，最后就保留最新的 那条数据

因为__consumer_offsets可能会接收高并发的请求，kafka默认给其分配50个分区(可以通过 offsets.topic.num.partitions设置)，这样可以通过加机器的方式抗大并发。

------

### 消费者Rebalance机制

rebalance就是说**如果消费组里的消费者数量有变化或消费的分区数有变化，kafka会重新分配消费者与消费分区的关系**。 比如consumer group中某个消费者挂了，此时会自动把分配给他的分区交给其他的消费者，如果他又重启了，那么又会把一些分区重新交还给他。

**注意**：rebalance只针对subscribe这种不指定分区消费的情况，如果通过assign这种消费方式指定了分区，kafka不会进 行rebanlance。

如下情况可能会触发消费者rebalance:

1.  消费组里的consumer增加或减少了
2.  动态给topic增加了分区
3.  消费组订阅了更多的topic

rebalance过程中，消费者无法从kafka消费消息，这对kafka的TPS会有影响，如果kafka集群内节点较多，比如数百 个，那重平衡可能会耗时极多，所以应尽量避免在系统高峰期的重平衡发生。

### Rebalance过程如下

当有消费者加入消费组时，消费者、消费组及组协调器之间会经历以下几个阶段 ![image.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/95908b430e5f4d61b5fa9a055eb836aa~tplv-k3u1fbpfcp-watermark.awebp)

#### 第一阶段：选择组协调器

组协调器GroupCoordinator：每个consumer group都会选择一个broker作为自己的组协调器coordinator，负责监控 这个消费组里的所有消费者的心跳，以及判断是否宕机，然后开启消费者rebalance。 consumer group中的每个consumer启动时会向kafka集群中的某个节点发送 FindCoordinatorRequest 请求来查找对 应的组协调器GroupCoordinator，并跟其建立网络连接。 组协调器选择方式： 通过如下公式可以选出consumer消费的offset要提交到__consumer_offsets的哪个分区，这个分区leader对应的broker 就是这个consumer group的coordinator 公式：`hash(consumer group id) % 对应主题的分区数`

#### 第二阶段：加入消费组JOIN GROUP

在成功找到消费组所对应的 GroupCoordinator 之后就进入加入消费组的阶段，在此阶段的消费者会向 GroupCoordinator 发送 JoinGroupRequest 请求，并处理响应。然后GroupCoordinator 从一个consumer group中 选择第一个加入group的consumer作为leader(消费组协调器)，把consumer group情况发送给这个leader，接着这个 leader会负责制定分区方案。

#### 第三阶段（ SYNC GROUP)

consumer leader通过给GroupCoordinator发送SyncGroupRequest，接着GroupCoordinator就把分区方案下发给各 个consumer，他们会根据指定分区的leader broker进行网络连接以及消息消费。

------

### 消费者Rebalance分区分配策略

主要有三种rebalance的策略：`range`、`round-robin`、`sticky`。**默认情况为range分配策略**。

假设一个主题有10个分区(0-9)，现在有三个consumer消费：

**range策略**：`按照分区序号排序分配`，假设 n＝分区数／消费者数量 = 3， m＝分区数%消费者数量 = 1，那么前 m 个消 费者每个分配 n+1 个分区，后面的（消费者数量－m ）个消费者每个分配 n 个分区。 比如分区0~ 3给一个consumer，分区4~ 6给一个consumer，分区7~9给一个consumer。

**round-robin策略**：`轮询分配`，比如分区0、3、6、9给一个consumer，分区1、4、7给一个consumer，分区2、5、 8给一个consumer

**sticky策略**： **初始时分配策略与round-robin类似，但是在rebalance的时候，需要保证如下两个原则。**

1.  分区的分配要尽可能均匀 。
2.  分区的分配尽可能与上次分配的保持相同。

当两者发生冲突时，**第一个目标优先于第二个目**标 。这样可以最大程度维持原来的分区分配的策略。 比如对于第一种range情况的分配，如果第三个consumer挂了，那么重新用sticky策略分配的结果如下： consumer1除了原有的0~ 3，会再分配一个7 consumer2除了原有的4~ 6，会再分配8和9

------

### producer发布消息机制剖析

1.  写入方式

    producer 采用 push 模式将消息发布到 broker，每条消息都被 append 到 patition 中，属于顺序写磁盘（**顺序写磁盘 比 随机写 效率要高，保障 kafka 吞吐率**）。

2.  消息路由

    producer 发送消息到 broker 时，会根据分区算法选择将其存储到哪一个 partition。其路由机制为：

    1.  指定了 patition，则直接使用；
    2.  未指定 patition 但指定 key，通过`hash(key)%分区数`算出路由的patition, 如果patition 和 key 都未指定，使用轮询选出一个patition。

3.  写入流程

![image.png](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/9f8bf6e56004456c857a871e5d9f403b~tplv-k3u1fbpfcp-watermark.awebp)

1.  producer 先从 zookeeper 的 "/brokers/.../state" 节点找到该 partition 的 leader
2.  producer 将消息发送给该 leader
3.  leader 将消息写入本地 log
4.  followers 从 leader pull 消息，写入本地 log 后 向leader 发送 ACK
5.  leader 收到所有 ISR 中的 replica 的 ACK 后，增加 HW（high watermark，最后 commit 的 offset） 并向 producer 发送 ACK

------

### HW与LEO

`HW俗称高水位`，HighWatermark的缩写，取一个partition对应的ISR中最小的LEO(log-end-offset)作为HW， consumer最多只能消费到HW所在的位置。另外每个replica都有HW,leader和follower各自负责更新自己的HW的状 态。`对于leader新写入的消息，consumer不能立刻消费，leader会等待该消息被所有ISR中的replicas同步后更新HW， 此时消息才能被consumer消费`。**这样就保证了如果leader所在的broker失效，该消息仍然可以从新选举的leader中获取。对于来自内部broker的读取请求，没有HW的限制。**

------

### 日志分段存储

Kafka 一个分区的消息数据对应存储在一个文件夹下，以topic名称+分区号命名，消息在分区内是分段存储的， 每个段的消息都存储在不一样的log文件里，kafka规定了一个段位的 log 文 件最大为 1G，做这个限制目的是为了方便把 log 文件加载到内存去操作：

```js
1 # 部分消息的offset索引文件，kafka每次往分区发4K(可配置)消息就会记录一条当前消息的offset到index文件， 
2 # 如果要定位消息的offset会先在这个文件里快速定位，再去log文件里找具体消息 
3 00000000000000000000.index 
4 # 消息存储文件，主要存offset和消息体 
5 00000000000000000000.log 
6 # 消息的发送时间索引文件，kafka每次往分区发4K(可配置)消息就会记录一条当前消息的发送时间戳与对应的offset到timeindex文件， 
7 # 如果需要按照时间来定位消息的offset，会先在这个文件里查找 
8 00000000000000000000.timeindex 
9 
10 00000000000005367851.index 
11 00000000000005367851.log 
12 00000000000005367851.timeindex 
13 
14 00000000000009936472.index 
15 00000000000009936472.log 
16 00000000000009936472.timeindex
复制代码
```

这个 9936472 之类的数字，就是代表了这个日志段文件里包含的起始 Offset，也就说明这个分区里至少都写入了接近 1000 万条数据了。 Kafka Broker 有一个参数，log.segment.bytes，限定了每个日志段文件的大小，最大就是 1GB。 一个日志段文件满了，就自动开一个新的日志段文件来写入，避免单个文件过大，影响文件的读写性能，这个过程叫做 log rolling，正在被写入的那个日志段文件，叫做 active log segment。

------

### 最后附一张zookeeper节点数据图

![image.png](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/8d6cdec325d34e77aace51314cd81560~tplv-k3u1fbpfcp-watermark.awebp)

# MQ带来的一些问题、及解决方案

## 1. 如何保证顺序消费？

1.  **RabbitMQ** :  一个Queue对应一个Consumer即可解决。

2.  RocketMQ

    -   全局有序：Topic里面只有一个MessageQueue即可。
    -   局部有序: 根据路由算法，比如`hash(key)%队列数`得到路由索引，使得需要保证有序的消息都路由到同一个MessageQueue。

3.  Kafka

    :

    -   全局有序：Topic里面只有一个Partition即可。
    -   局部有序: 根据路由算法，比如`hash(key)%分区数`得到路由索引，使得需要保证有序的消息都路由到同一个Partition。

## 2. 如何实现延迟消费？

1.  RabbitMQ

    : 两种方案

    -   死信队列 + TTL
    -   引入RabbitMQ的延迟插件

2.  **RocketMQ**：天生支持延时消息。

3.  Kafka

    : 步骤如下

    -   专门为要延迟的消息创建一个Topic
    -   新建一个消费者去消费这个Topic
    -   消息持久化
    -   再开一个线程定时去拉取持久化的消息，放入实际要消费的Topic
    -   实际消费的消费者从实际要消费的Topic拉取消息。

![image.png](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/1e53792d4f49475a96613d87c09fb940~tplv-k3u1fbpfcp-watermark.awebp)

## 3. 如何保证消息的可靠性投递

1.  **RabbitMQ**:

    -   Broker-->消费者 : 手动ACK
    -   生产者-->Broker: 两种方案
    -   -   `1. 数据库持久化`

    ```js
    1.将业务订单数据和生成的Message进行持久化操作（一般情况下插入数据库，这里如果分库的话可能涉及到分布式事务）
    2.将Message发送到Broker服务器中
    3.通过RabbitMQ的Confirm机制，在producer端，监听服务器是否ACK。
    4.如果ACK了，就将Message这条数据状态更新为已发送。如果失败，修改为失败状态。
    5.分布式定时任务查询数据库3分钟（这个具体时间应该根据的时效性来定）之前的发送失败的消息
    6.重新发送消息，记录发送次数
    7.如果发送次数过多仍然失败，那么就需要人工排查之类的操作。
    复制代码
    ```

    ![image.png](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/3609c1560670473282601e0eced567dc~tplv-k3u1fbpfcp-watermark.awebp)

    **优点**：能够保证消息百分百不丢失。

    **缺点**：第一步会涉及到分布式事务问题。

    -   -   `2. 消息的延迟投递`

    ```js
    流程图中，颜色不同的代表不同的message
    1.将业务订单持久化
    2.发送一条Message到broker(称之为主Message)，再发送相同的一条到不同的队列或者交换机(这条称为确认Message)中。
    3.主Message由实际业务处理端消费后，生成一条响应Message。之前的确认Message由Message Service应用处理入库。
    4~6.实际业务处理端发送的确认Message由Message Service接收后，将原Message状态修改。
    7.如果该条Message没有被确认，则通过rpc调用重新由producer进行全过程。
    复制代码
    ```

![image.png](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/a0a5c78c593d495488a249c6e79aeaf5~tplv-k3u1fbpfcp-watermark.awebp)

**优点**：相对于持久化方案来说响应速度有所提升

**缺点**：系统复杂性有点高，万一两条消息都失败了，消息存在丢失情况，仍需Confirm机制做补偿。

1.  **RocketMQ**：

-   **生产者弄丢数据**

    Producer在把Message发送Broker的过程中，因为网络问题等发生丢失，或者Message到了Broker，但是出了问题，没有保存下来。针对这个问题，RocketMQ对Producer发送消息设置了3种方式：

    1.  `同步发送`：天生保证了可靠性投递
    2.  `异步发送`：需要在回调函数中，根据broker响应的结果自定义实现。
    3.  `单向发送`：保证不了可靠性投递

-   **Broker弄丢数据**

　　Broker接收到Message暂存到内存，Consumer还没来得及消费，Broker挂掉了

　　可以通过`持久化`设置去解决：

　　1. 创建Queue的时候设置持久化，保证Broker持久化Queue的元数据，但是不会持久化Queue里面的消息

　　2. 将Message的deliveryMode设置为2，可以将消息持久化到磁盘，这样只有Message支持化到磁盘之后才会发送通知Producer ack

　　这两步过后，即使Broker挂了，Producer肯定收不到ack的，就可以进行重发

-   **消费者弄丢数据**

　　Consumer有消费到Message，但是内部出现问题，Message还没处理，Broker以为Consumer处理完了，只会把后续的消息发送。这时候，就要`关闭autoack，消息处理过后，进行手动ack`, 多次消费失败的消息，会进入`死信队列`，这时候需要人工干预。

1.  **Kafka**:

-   **生产者弄丢数据**

    设置了 `acks=all`，一定不会丢，要求是，你的 leader 接收到消息，所有的 follower 都同步到了消息之后，才认为本次写成功了。如果没满足这个条件，生产者会自动不断的重试，重试无限次。

-   **Broker弄丢数据**

Kafka 某个 broker 宕机，然后重新选举 partition 的 leader。大家想想，要是此时其他的 follower 刚好还有些数据没有同步，结果此时 leader 挂了，然后选举某个 follower 成 leader 之后，不就少了一些数据？这就丢了一些数据啊。

此时一般是要求起码设置如下 4 个参数：

-   -   给 topic 设置 `replication.factor` 参数：这个值必须大于 1，要求每个 partition 必须有至少 2 个副本。
-   -   在 Kafka 服务端设置 `min.insync.replicas` 参数：这个值必须大于 1，这个是要求一个 leader 至少感知到有至少一个 follower 还跟自己保持联系，没掉队，这样才能确保 leader 挂了还有一个 follower 吧。
-   -   在 producer 端设置 `acks=all`：这个是要求每条数据，必须是写入所有 replica 之后，才能认为是写成功了。
-   -   在 producer 端设置 `retries=MAX`（很大很大很大的一个值，无限次重试的意思）：这个是要求一旦写入失败，就无限重试，卡在这里了。

我们生产环境就是按照上述要求配置的，这样配置之后，至少在 Kafka broker 端就可以保证在 leader 所在 broker 发生故障，进行 leader 切换时，数据不会丢失。

-   **消费者弄丢数据**

你消费到了这个消息，然后消费者那边自动提交了 offset，让 Kafka 以为你已经消费好了这个消息，但其实你才刚准备处理这个消息，你还没处理，你自己就挂了，此时这条消息就丢咯。

这不是跟 RabbitMQ 差不多吗，大家都知道 Kafka 会自动提交 offset，那么只要**关闭自动提交 offset，在处理完之后自己手动提交 offset，就可以保证数据不会丢。** 但是此时确实还是可能会有重复消费，比如你刚处理完，还没提交 offset，结果自己挂了，此时肯定会重复消费一次，自己保证幂等性就好了。

## 4. 如何保证消息的幂等？

以 RocketMQ 为例，下面列出了消息重复的场景：

**1.发送时消息重复**

当一条消息已被成功发送到服务端并完成持久化，此时出现了网络闪断或者客户端宕机，导致服务端对客户端应答失败。如果此时生产者意识到消息发送失败并尝试再次发送消息，消费者后续会收到两条内容相同并且Message ID也相同的消息。

**2.投递时消息重复**

消息消费的场景下，消息已投递到消费者并完成业务处理，当客户端给服务端反馈应答的时候网络闪断。为了保证消息至少被消费一次，消息队列RocketMQ版的服务端将在网络恢复后再次尝试投递之前已被处理过的消息，消费者后续会收到两条内容相同并且Message ID也相同的消息。

3.**负载均衡时消息重复**（包括但不限于网络抖动、Broker重启以及消费者应用重启）

当消息队列RocketMQ版的Broker或客户端重启、扩容或缩容时，会触发Rebalance，此时消费者可能会收到重复消息。

那么，有什么解决方案呢？ 直接上图。

![image.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/cd2ff7a27a19411693a0f3e031af0411~tplv-k3u1fbpfcp-watermark.awebp)

## 5. 如何解决消息积压的问题？

关于这个问题，有几个点需要考虑：

**1. 如何快速让积压的消息被消费掉？**

临时写一个消息分发的消费者，把积压队列里的消息均匀分发到N个队列中，同时一个队列对应一个消费者，相当于消费速度提高了N倍。

修改前：

![image.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/4fc1111271aa42448aed1bac4d5c3994~tplv-k3u1fbpfcp-watermark.awebp)

修改后：

![image.png](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/645fc869b39742149b459ea1a0a2e594~tplv-k3u1fbpfcp-watermark.awebp)

**2. 积压时间太久，导致部分消息过期，怎么处理？**

批量重导。在业务不繁忙的时候，比如凌晨，提前准备好程序，把丢失的那批消息查出来，重新导入到MQ中。

**3. 消息大量积压，MQ磁盘被写满了，导致新消息进不来了，丢掉了大量消息，怎么处理？**

这个没办法。谁让【消息分发的消费者】写的太慢了，你临时写程序，接入数据来消费，消费一个丢弃一个，都不要了，快速消费掉所有的消息。然后走第二个方案，到了晚上再补数据吧。


作者：Boom
链接：https://juejin.cn/post/7006958043833303048
来源：掘金
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。