---
title: mysql datetime和int分析对比
---

https://blog.csdn.net/riemann_/article/details/100679675



一、三者的区别
1、int

占用4个字节
建立索引之后，查询速度快
条件范围搜索可以使用使用between
不能使用mysql提供的时间函数
结论：适合需要进行大量时间范围查询的数据表

2、datetime

占用8个字节
允许为空值，可以自定义值，系统不会自动修改其值
实际格式储存（Just stores what you have stored and retrieves the same thing which you have stored.）
与时区无关（It has nothing to deal with the TIMEZONE and Conversion.）
可以在指定datetime字段的值的时候使用now()变量来自动插入系统的当前时间
datetime是可以设置默认值的, DEFAULT CURRENT_TIMESTAMP在mysql5.7中亲测可用
结论：datetime类型适合用来记录数据的原始的创建时间，因为无论你怎么更改记录中其他字段的值，datetime字段的值都不会改变，除非你手动更改它。

3、timestamp

占用4个字节
允许为空值，但是不可以自定义值，所以为空值时没有任何意义
TIMESTAMP值不能早于1970或晚于2037。这说明一个日期，例如’1968-01-01’，虽然对于DATETIME或DATE值是有效的，但对于TIMESTAMP值却无效，如果分配给这样一个对象将被转换为0。
值以UTC格式保存（ it stores the number of milliseconds）
时区转化 ，存储时对当前的时区进行转换，检索时再转换回当前的时区
默认值为CURRENT_TIMESTAMP()，其实也就是当前的系统时间
数据库会自动修改其值，所以在插入记录时不需要指定timestamp字段的名称和timestamp字段的值，你只需要在设计表的时候添加一个timestamp字段即可，插入后该字段的值会自动变为当前系统时间
默认情况下以后任何时间修改表中的记录时，对应记录的timestamp值会自动被更新为当前的系统时间
如果需要可以设置timestamp不自动更新。通过设置DEFAULT CURRENT_TIMESTAMP 可以实现
修改自动更新
field_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE
修改不自动更新
field_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
结论：timestamp类型适合用来记录数据的最后修改时间，因为只要你更改了记录中其他字段的值，timestamp字段的值都会被自动更新。（如果需要可以设置timestamp不自动更新）

二、MySQL中如何表示当前时间？
其实，表达方式还是蛮多的，汇总如下：

CURRENT_TIMESTAMP

CURRENT_TIMESTAMP()

NOW()

LOCALTIME

LOCALTIME()

LOCALTIMESTAMP

LOCALTIMESTAMP()

三、关于timestamp和datetime的比较
一个完整的日期格式如下：YYYY-MM-DD HH:MM:SS[.fraction]，它可分为两部分：date部分和time部分，其中，date部分对应格式中的“YYYY-MM-DD”，time部分对应格式中的“HH:MM:SS[.fraction]”。对于date字段来说，它只支持date部分，如果插入了time部分的内容，它会丢弃掉该部分的内容，并提示一个warning。

1、关于update_time设置成date类型


insert into user(update_time) values('20190909000000');
1




date后面的time部分都丢弃掉了。

2、timestamp和datetime的相同点

两者都可用来表示YYYY-MM-DD HH:MM:SS[.fraction]类型的日期。

3、timestamp和datetime的不同点

（1）、两者的存储方式不一样

对于TIMESTAMP，它把客户端插入的时间从当前时区转化为UTC（世界标准时间）进行存储。查询时，将其又转化为客户端当前时区进行返回。

而对于DATETIME，不做任何改变，基本上是原样输入和输出。

下面，我们来验证一下

首先创建两种测试表，一个使用timestamp格式，一个使用datetime格式。

user表 datetime类型

user2表 timestamp类型


insert into user(id, update_time) values(1, '20190910000800');
insert into user2(id, update_time) values(1, '20190910000800');

select * from user;
select * from user2;
1
2
3
4
5


两者输出是一样的。

其次修改当前会话的时区

show variables like '%time_zone%'; 
1


  set time_zone='+0:00';
1
select * from user;
select * from user2;
1
2



“CST”指的是MySQL所在主机的系统时间，是中国标准时间的缩写，China Standard Time UT+8:00，这里system_time_zone为空。

通过结果可以看出，user2中返回的时间提前了8个小时，而user中时间则不变。这充分验证了两者的区别。

（2）、两者所能存储的时间范围不一样

timestamp所能存储的时间范围为：'1970-01-01 00:00:01.000000' 到 '2038-01-19 03:14:07.999999'。

datetime所能存储的时间范围为：'1000-01-01 00:00:00.000000' 到 '9999-12-31 23:59:59.999999'。

总结：timestamp和datetime除了存储范围和存储方式不一样，没有太大区别。当然，对于跨时区的业务，timestamp更为合适。

四、关于timestamp和datetime的自动初始化和更新
首先，我们来看下面的操作：

mysql> create table user3(id int, update_time timestam
Query OK, 0 rows affected (0.02 sec)

mysql> insert into user3(id) values(1);
Query OK, 1 row affected (0.00 sec)

mysql> select * from user3;
+------+---------------------+
| id   | update_time         |
+------+---------------------+
|    1 | 2019-09-10 00:37:57 |
+------+---------------------+
1 row in set (0.00 sec)

mysql> show create table user3\G
*************************** 1. row *******************
       Table: user3
Create Table: CREATE TABLE `user3` (
  `id` int(11) default NULL,
  `update_time` timestamp NOT NULL default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8
1 row in set (0.00 sec)

看起来是不是有点奇怪，我并没有对update_time字段进行插入操作，它的值自动修改为当前值，而且在创建表的时候，我也并没有定义“show create table test\G”结果中显示的“ default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP”。

其实，这个特性是自动初始化和自动更新（Automatic Initialization and Updating）。

自动初始化指的是如果对该字段（譬如上例中的update_time字段）没有显性赋值，则自动设置为当前系统时间。

自动更新指的是如果修改了其它字段，则该字段的值将自动更新为当前系统时间。

它与“explicit_defaults_for_timestamp”参数有关。

默认情况下，该参数的值为OFF，如下所示：

mysql> show variables like '%explicit_defaults_for_timestamp%';
+---------------------------------+-------+
| Variable_name                   | Value |
+---------------------------------+-------+
| explicit_defaults_for_timestamp | OFF   |
+---------------------------------+-------+
row in set (0.00 sec)

下面我们看看官档的说明：

By default, the first TIMESTAMP column has both DEFAULT CURRENT_TIMESTAMP and ON UPDATE CURRENT_TIMESTAMP if neither is specified explicitly。

很多时候，这并不是我们想要的，如何禁用呢？

将“explicit_defaults_for_timestamp”的值设置为ON。

“explicit_defaults_for_timestamp”的值依旧是OFF，也有两种方法可以禁用

<1> 用DEFAULT子句该该列指定一个默认值

<2> 为该列指定NULL属性。

如下所示：

mysql> create table user4(id int, update_time timestamp null);
Query OK, 0 rows affected (0.01 sec)

mysql> show create table user4\G
*************************** 1. row ***************************
       Table: user4
Create Table: CREATE TABLE `user4` (
  `id` int(11) DEFAULT NULL,
  `update_time` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf-8
row in set (0.00 sec)

mysql> create table user5(id int,update_time timestamp default 0);
Query OK, 0 rows affected (0.01 sec)

mysql> show create table user5\G
*************************** 1. row ***************************
       Table: user5
Create Table: CREATE TABLE `user5` (
  `id` int(11) DEFAULT NULL,
  `update_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00'
) ENGINE=InnoDB DEFAULT CHARSET=utf-8
row in set (0.00 sec)

在MySQL 5.6.5版本之前，Automatic Initialization and Updating只适用于TIMESTAMP，而且一张表中，最多允许一个TIMESTAMP字段采用该特性。从MySQL 5.6.5开始，Automatic Initialization and Updating同时适用于TIMESTAMP和DATETIME，且不限制数量。

------------------------------------------------
版权声明：本文为CSDN博主「老周聊架构」的原创文章，遵循CC 4.0 BY-SA版权协议，转载请附上原文出处链接及本声明。
原文链接：https://blog.csdn.net/riemann_/article/details/100679675