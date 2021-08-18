---
title: MySQL外键对性能影响
---

In MySQL, a foreign key requires an index. If an index already exists, the foreign key will use that index (even using the prefix of an existing multi-column index). If no index exists, defining the foreign key will build the index.

So the size increase and time to create a foreign key is about the same as to create an index on the same column(s).

The performance of a SELECT is not impacted significantly by the presence of a foreign key. Only a slight additional work for the optimizer, to consider the new indexes.

The performance of updating is more, because for each foreign key, an INSERT/UPDATE/DELETE has to check to see if the constraint is satisfied. That means a primary key lookup to the referenced tables. This impact is measurable, and it is greater if the referenced tables are not in the buffer pool.

Another impact is the locking. If I update a row in a child row that has a foreign key, InnoDB places a shared-lock on the referenced rows in the parent tables. That means no one can update those parent rows until I commit. If you have lots of threads updating child rows, then the parent rows may be locked most of the time, and this can make it hard to do concurrent work in the parent tables. Not so much a performance problem, but a concurrency problem.

As with all "how does that perform" questions, the answer really depends on your workload. If you don't have concurrent updates, for example, that issue may not effect you for all practical purposes. Testing it yourself with load testing is the only way to be sure. It's not something anyone on StackOverflow can answer precisely for you.

读后小结：

1.外键就想其他类型的索引一样，要说性能方面的影响，它影响的主要是写入操作(如，UPDATE/INSERT/DELETE)；

2.但它与单表索引不同的是，它会引用一张或多张父表，这样当对子表进行写入操作(UPDATE/INSERT)的时候，父表就会被加上“共享锁”，这样在对子表高并发进行写入操作的情况下，对父表的写入操作就会由于“共享锁”的存在，而会长时间不能得到更新！当然查询是可以的。

3.作者给出了一个观点：外键主要造成的问题，并不是影响性能多少的问题，而是一个并发访问的问题！

4.所以是否使用外键，关键是要看你对并发要求的高低了！