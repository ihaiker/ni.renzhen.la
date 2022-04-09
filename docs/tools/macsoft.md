---
title: mac软件修复
---



[“Mac应用”已损坏，打不开解决办法 for Mac - Mac毒 (macdo.cn)](https://www.macdo.cn/925.html)





### 2. macOS Catalina 10.15系统：

打开「终端.app」，输入以下命令并回车，输入开机密码回车

```
sudo xattr -rd com.apple.quarantine 空格 软件的路径
```

如Sketch.app

```
sudo xattr -rd com.apple.quarantine /Applications/Sketch.app
```

如CleanMyMac X.app

```
sudo xattr -rd com.apple.quarantine /Applications/CleanMyMac X.app
```

### 3. macOS Catalina 10.15.4 系统：

更新10.15.4系统后软件出现意外退出，可按照下面的方法给软件签名

**1.**打开「终端app」输入如下命令：

```
xcode-select --install
```

**2.给软件签名**

打开终端工具输入并执行如下命令：

```
sudo codesign --force --deep --sign - (应用路径)
```



**3.错误解决**

如出现以下错误提示：

```
/文件位置 : replacing existing signature
/文件位置 : resource fork,Finder information,or similar detritus not allowed
```

 

那么，先在终端执行：

```
xattr -cr /文件位置（直接将应用拖进去即可）
```

然后再次执行如下指令即可：

```
codesign --force --deep --sign - /文件位置（直接将应用拖进去即可）
```