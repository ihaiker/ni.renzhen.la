---
tilte: mysql varchar  length 对性能影响分析
order: 1
---



1.MySQL建立索引时如果没有限制索引的大小，索引长度会默认采用的该字段的长度，也就是说varchar(20)和varchar(255)对应的索引长度分别为20*3(utf-8)(+2+1),255*3(utf-8)(+2+1)，其中"+2"用来存储长度信息，“+1”用来标记是否为空，加载索引信息时用varchar(255)类型会占用更多的内存； （备注：当字段定义为非空的时候，是否为空的标记将不占用字节）例如，测试sql(InnoDB引擎)如下:

```sql
CREATE DATABASE TestDataBase;
USE TestDataBase;

CREATE TABLE ABC (
	`id` int(11) DEFAULT null,
  `name` varchar(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT charset=utf8;

ALTER TABLE `ABC` ADD INDEX `nameIndex`(`name`);
explain select name from ABC;

alter table ABC modify name varchar(255);
explain select name from ABC;
```

varchr(10)变长字段且允许NULL:10*(Character Set：utf8=3,gbk=2,latin1=1)+1(NULL)+2(变长字段) 

varchr(10)变长字段且不允许NULL:10*(Character Set：utf8=3,gbk=2,latin1=1)+2(变长字段) 

char(10)固定字段且允许NULL:10*(Character Set：utf8=3,gbk=2,latin1=1)+1(NULL) char(10)固定

字段且允许NULL:10*(Character Set：utf8=3,gbk=2,latin1=1)根据这个值，就可以判断索引使用情况，特别是在组合索引的时候，判断所有的索引字段都被查询用到。

2.varchar(20)与varchar(255)都是保持可变的字符串，当使用ROW_FORMAT=FIXED创建MyISAM表时，会为每行使用固定的长度空间，这样设置不同的varchar长度值时，存储相同数据所占用的空间是不一样。

通常情况下使用varchar(20)和varchar(255)保持'hello'占用的空间都是一样的，但使用长度较短的列却有巨大的优势。较大的列使用更多的内存，因为MySQL通常会分配固定大小的内存块来保存值，这对排序或使用基于内存的临时表尤其不好。同样的事情也会发生在使用文件排序或者基于磁盘的临时表的时候。