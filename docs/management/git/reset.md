

```
reset命令有3种方式：
1：git reset –mixed：此为默认方式，不带任何参数的git reset，即时这种方式，它回退到某个版本，只保留源码，回退commit和index信息

2：git reset –soft：回退到某个版本，只回退了commit的信息，不会恢复到index file一级。如果还要提交，直接commit即可

3：git reset –hard：彻底回退到某个版本，本地的源码也会变为上一个版本的内容
```

git reset只是在本地仓库中回退版本，而远程仓库的版本不会变化。
以删除master分支为例

```
#新建一个备份的分支，数据无价
git branch old_master

#提交本地当前的文件到新建的分支
git push origin old_master:old_master

#本地可以彻底恢复到你想恢复到的版本了
git reset --hard 58093e1355716f0f861b64f1c3dfe59242be28f7

#在web端settings页面，修改默认分支为新建的分支，可以删除远程分支了
git push origin :master

#如果出现! [remote rejected] master (deletion of the current branch prohibited)，说明没有设置远程的默认分支，没有权限删除，请在web端settings页面，修改默认分支为新建的分支

#进行到这里，远程的master分支已经删除成功
#重新提交本地文件到master分支（此时会自动新建master分支）
git push origin master

#再体验一下删除分支
git push origin :old_master
```