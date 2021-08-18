---
title: Git基础命令
---



1.   下载远程分支A到本地A

```
git checkout -b <BranchName> origin/<BranchName>
```

>   克隆一个远程分支本地，且本地分支不存在，如果分支存在报错，如果想要直接强制操作，可以报-b该为-f

