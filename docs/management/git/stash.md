---
title: stash 工作区保存 
---

-   **保存工作区内容**

    >   git stash save [message]
    >
    >   说明：
    >
    >   1.  将工作区未提交的修改封存，让工作区回到修改前的状态

-   **查看工作区列表**

    >   git stash list
    >
    >   说明：最新保存的工作区在最上面

-   **应用某个工作区**

    >   git stash apply [stash@{n}]

-   **删除工作区**

    >   git stash drop [stash@{n}] # 删除某一个工作区
    >
    >   git stash clear # 删除（清空）所有保存的工作区