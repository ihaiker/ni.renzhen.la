# SQL语句不走索引时的排查

## 前言：

>   在索引优化时，经常会看到的一句话：如果索引字段出现隐式字符集转换的话，那么索引将失效，进而转为全表扫描，查询效率将大大降低，要避免出现隐式字符集转换；

## 在此我想问问同学们：

-   大家知道为什么隐式字符集转换会导致索引失效吗？
-   实际场景中有没有遇到过隐式字符集转换导致索引失效的场景，具体排查的过程；

## 本文主线

由上面的两个问题牵引出了本文的主线；

-   简单描述下隐式字符集转换导致索引失效的原因
-   然后模拟实际场景排查隐式字符集转换导致索引失效的过程

### 隐式字符集转换导致索引失效的原因

MySQL索引的数据结构是 B+Tree，想要走索引查询必须要满足其 **最左前缀原则** ，否则无法通过索引树进行查找，只能进行全表扫描；

例如：下面的这个SQL由于在 **索引字段** 上使用函数进行运算，导致索引失效

```sql
select * from t_user where SUBSTR(name, 1, 2) = '李彤'
```

上面的这个SQL怎么改造才能使索引生效呢？如下所示：

```sql
select * from t_user where name like '李彤%'
```

通过上面的小例子可以知道，如果在索引字段上使用函数运算，则会导致索引失效，而索引字段的 **隐式字符集转换** 由于MySQL会自动的在索引字段上加上 **转换函数** ，进而会导致索引失效；

那接下来我们就通过模拟的实际场景来具体看看是不是由于MySQL自动给加上了转换函数而导致索引失效的；

### 模拟场景 + 问题排查

>   由于导致索引失效的原因有很多，如果自己写的SQL怎么看都没问题，但是通过查看执行计划发现就是没有走索引查询，此时就会让很多人陷入困境，这到底是怎么导致的呢？

此时本文重点将要讲述的工具就要闪亮登场啦： **explain extended + show warnings** ；

使用这个工具可以将执行的SQL语句的一些扩展信息展示出来，这些扩展信息就包括：MySQL优化时可能会添加上字符集转换函数，使得字符集不匹配的SQL可以正确执行下去；

下面就来具体聊聊 **explain extended + show warnings** 的使用；

#### 模拟隐式字符集转换的场景：

**首先创建两个字符集不一样的表：**

```sql
CREATE TABLE `t_department` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `de_no` varchar(32) NOT NULL,
  `info` varchar(200) DEFAULT NULL,
  `de_name` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_de_no` (`de_no`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4;


CREATE TABLE `t_employees` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `em_no` varchar(32) NOT NULL,
  `de_no` varchar(32) NOT NULL,
  `age` int(11) DEFAULT NULL,
  `info` varchar(200) DEFAULT NULL,
  `em_name` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_em_no` (`de_no`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8;
```

**然后使用存储过程构造数据：**

```sql
# 如果存储过程存在则删除 
DROP PROCEDURE IF EXISTS proc_initData;

DELIMITER $
# 创建存储过程
CREATE PROCEDURE proc_initData()
BEGIN
    DECLARE i INT  DEFAULT 1;
    WHILE i<=30 DO
        # 新增数据
        INSERT INTO t_employees ( em_no, de_no, info, em_name , age) VALUES ( CONCAT('001', i), '003', 'test11', 'test2', i ); #执行的sql语句
        SET i = i+1;
    END WHILE;
END $

# 调用存储过程
CALL proc_initData();
```

注意：在构造数据时，记得将 t_employees 表中的 de_no 字段值构造的 **离散些** ，因为如果索引字段值的 **区分度很低** 的话，那么MyQSL优化器通过采样统计分析时，发现索引查询和全表扫描性能差不多，就会直接进行全表扫描了；

**索引失效的查询SQL语句：**

将表和数据构造完后，我们使用SQL语句进行查询下，然后再看看其执行计划；

```sql
explain 
select * from t_department a LEFT JOIN  t_employees b on a.de_no = b.de_no where a.id = 16
```

其执行计划如下：

![img](https://cdn.jsdelivr.net/gh/leishen6/ImgHosting/MuZiLei_blog_img/20210706092625.png)

发现 t_employees 表中的 de_no 字段有索引，但是没有走索引查询，type=ALL 走的全表扫描，但是通过查看SQL语句发现其没有问题呀，表面看上去都是满足走索引查询的条件呀，排查到这发现遇到了困境，苦恼啊！

![img](https://cdn.jsdelivr.net/gh/leishen6/ImgHosting/MuZiLei_blog_img/20210706093823.gif)

还好，通过在网络世界上遨游，最终发现了 **explain extended + show warnings** 利器，利用它快速发现了索引失效的根本原因，然后快速找到了解决方案；

下面就来聊聊这个利器的具体使用，开森！

![img](https://cdn.jsdelivr.net/gh/leishen6/ImgHosting/MuZiLei_blog_img/20210706094319.jpg)

#### 使用利器快速排查问题：

>   注意：explain 后面跟的关键字 EXTENDED（扩展信息） 在MySQL5.7及之后的版本中废弃了，但是该语法仍被识别为向后兼容，所以在5.7版本及后续版本中，可以不用在 explain 后面添加 EXTENDED 了；

![img](https://cdn.jsdelivr.net/gh/leishen6/ImgHosting/MuZiLei_blog_img/20210706095229.png)

EXTENDED关键字的具体查阅资料：[https://dev.mysql.com/doc/refman/5.7/en/explain-extended.html](https://www.oschina.net/action/GoToLink?url=https%3A%2F%2Fdev.mysql.com%2Fdoc%2Frefman%2F5.7%2Fen%2Fexplain-extended.html)

**具体使用方法如下：**

①、首先在MySQL的可视化工具中打开一个 **命令列介面** ：工具 --> 命令列介面

②、然后输入下面的SQL并按回车：

```sql
explain EXTENDED
select * from t_department a LEFT JOIN  t_employees b on a.de_no = b.de_no where a.id = 4019;
```

③、然后紧接着输入命令 **show warnings;** 并回车，会出现如下图所示内容：

![img](https://cdn.jsdelivr.net/gh/leishen6/ImgHosting/MuZiLei_blog_img/20210706100548.png)

通过展示出的执行SQL扩展信息，发现MySQL在字符集不一致时自动添加上字符集转换函数，因为是在 **索引字段 de_no** 上添加的转换函数，所以就导致了索引失效；

而如果我们没看扩展信息的话，那么可能直到我们查看表结构的时候才会发现是由于字符集不一致导致的，这样就会花费很多的时间；

#### 扩展：隐式类型转换

咱们聊完上面的隐式字符集转换导致索引失效的情况，再来简单聊聊另一种 **隐式类型转换** 导致索引失效的情况；

隐式类型转换：简单的说就是字段的类型与其赋值的类型不一致时会进行隐式的转换；

小例如下：

```sql
select * from t_employees where em_name = 123;
```

上面的SQL中 em_name 为索引字段，字段类型是 varchar，为其赋 int 类型的值时，会发现索引失效，这里也可以通过 **explain extended + show warnings** 查看，会发现如下图所示内容：

![img](https://cdn.jsdelivr.net/gh/leishen6/ImgHosting/MuZiLei_blog_img/20210706203110.png)

>   至此本文进入结尾，在此再说明下，上文中测试时使用的MySQL版本都是 **5.7** ；

