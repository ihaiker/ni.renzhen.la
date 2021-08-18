---
title: pstree 查看进程数
---



pstree -c -p -A $(pgrep dockerd)

The command below will list all the contents of /proc, and store the Redis PID for future use.

```
DBPID=$(pgrep redis-server) echo Redis is $DBPID ls /proc
```

Each process has it's own configuration and security settings defined within different files. `ls /proc/$DBPID`

For example, you can view and update the environment variables defined to that process. `cat /proc/$DBPID/environ`

```
docker exec -it db env
```

