---
title: "nginx知识梳理"
---



## 前言

本文已参与【[请查收｜你有一次免费申请掘金周边礼物的机会](https://juejin.cn/post/7000643252957216782)】活动。

参与评论，有机会获得掘金官方提供的 2 枚新版徽章，具体的抽奖细节请看文末。

对于初级开发、特别是是小白同学，希望小伙伴们认真看完，我相信你会有收获的。

创作不易，记得**点赞、关注、收藏**哟。

## Nginx概念

Nginx 是一个**高性能的 HTTP 和反向代理服务**。其特点是占有内存少，并发能力强，事实上nginx的并发能力在同类型的网页服务器中表现较好。

Nginx 专为性能优化而开发，性能是其最重要的考量指标，实现上非常注重效率，能经受住高负载的考验，有报告表明能支持高达50000个并发连接数。

在连接高并发的情况下，Nginx 是 Apache 服务不错的替代品：Nginx 在美国是做虚拟主机生意的老板们经常选择的软件平台之一。

## 反向代理

在说反向代理之前，先来说说什么是代理和正向代理。

### 代理

代理其实就是一个中介，A和B本来可以直连，中间插入一个C，C就是中介。刚开始的时候，代理多数是帮助**内网client**（局域网）访问**外网server**用的。 后来出现了反向代理，`反向`这个词在这儿的意思其实是指方向相反，即代理将来自外网客户端的请求转发到内网服务器，从外到内。

### 正向代理

>   正向代理即是客户端代理，代理客户端，服务端不知道实际发起请求的客户端。

正向代理类似一个跳板机，代理访问外部资源。

比如我们国内访问谷歌，直接访问访问不到，我们可以通过一个正向代理服务器，请求发到代理服服务上，代理服务器能够访问谷歌，这样由代理去访问谷歌取到返回数据，再返回给我们，这样我们就能访问谷歌了。

![image-20210904171928174](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/cd592a838773454694dd1e2f7adb103b~tplv-k3u1fbpfcp-watermark.awebp)

### 反向代理

>   反向代理即是服务端代理，代理服务端，客户端不知道实际提供服务的服务端。

客户端是感知不到代理服务器的存在。

是指以代理服务器来接受 Internet 上的连接请求，然后将请求转发给内部网络上的服务器，并将从服务器上得到的结果返回给 Internet 上请求连接的客户端，此时代理服务器对外就表现为一个反向代理服务器。

![image-20210904173138672](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/9ed3f205ac0146719d315215e685c4f2~tplv-k3u1fbpfcp-watermark.awebp)

## 负载均衡

关于负载均衡，先来举个例子：

>   地铁大家应该都坐过吧，我们一般在早高峰乘地铁时候，总有那么一个地铁口人最拥挤，这时候，一般会有个地铁工作人员A拿个大喇叭在喊“**着急的人员请走B口，B口人少车空**”。而这个地铁工作人员A就是负责负载均衡的。
>
>   为了提升网站的各方面能力，我们一般会把多台机器组成一个集群对外提供服务。然而，我们的网站对外提供的访问入口都是一个的，比如`www.taobao.com`。那么当用户在浏览器输入`www.taobao.com`的时候如何将用户的请求分发到集群中不同的机器上呢，这就是负载均衡在做的事情。

**负载均衡（Load Balance），意思是将负载（工作任务，访问请求）进行平衡、分摊到多个操作单元（服务器，组件）上进行执行。是解决高性能，单点故障（高可用），扩展性（水平伸缩）的终极解决方案。**

![image-20210904175555853](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/50cddae6f37542f5993a5fc500633755~tplv-k3u1fbpfcp-watermark.awebp)

>   Nginx提供的负载均衡主要有三种方式：轮询，加权轮询，Ip hash。

### 轮询

nginx默认就是轮询其权重都默认为1，服务器处理请求的顺序：ABCABCABCABC....

```java
upstream mysvr { 
    server 192.168.8.1:7070; 
    server 192.168.8.2:7071;
    server 192.168.8.3:7072;
}
复制代码
```

![img](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/d30e73ce66cf4c7c88a49f20073d253e~tplv-k3u1fbpfcp-watermark.awebp)

### 加权轮询

根据配置的权重的大小而分发给不同服务器不同数量的请求。如果不设置，则默认为1。下面服务器的请求顺序为：ABBCCCABBCCC....

```java
upstream mysvr { 
    server 192.168.8.1:7070 weight=1; 
    server 192.168.8.2:7071 weight=2;
    server 192.168.8.3:7072 weight=3;
}
复制代码
```

![img](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/eca7b9ca42a44e23a5cd70caf58e12fc~tplv-k3u1fbpfcp-watermark.awebp)

### ip_hash

iphash对客户端请求的ip进行hash操作，然后根据hash结果将同一个客户端ip的请求分发给同一台服务器进行处理，可以解决session不共享的问题。

```java
upstream mysvr { 
    server 192.168.8.1:7070; 
    server 192.168.8.2:7071;
    server 192.168.8.3:7072;
    ip_hash;
}
复制代码
```

![img](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/5a87faacfdb744c399ccb5a853c39a77~tplv-k3u1fbpfcp-watermark.awebp)

## 动静分离

### 动态与静态页面区别

-   静态资源： 当用户多次访问这个资源，资源的源代码永远不会改变的资源（如：HTML，JavaScript，CSS，img等文件）。
-   动态资源：当用户多次访问这个资源，资源的源代码可能会发送改变（如：.jsp、servlet 等）。

### 什么是动静分离

-   动静分离是让动态网站里的动态网页根据一定规则把不变的资源和经常变的资源区分开来，动静资源做好了拆分以后，我们就可以根据静态资源的特点将其做缓存操作，这就是网站静态化处理的核心思路。
-   动静分离简单的概括是：动态文件与静态文件的分离。

### 为什么要用动静分离

为了加快网站的解析速度，可以把动态资源和静态资源用不同的服务器来解析，加快解析速度。降低单个服务器的压力。

![image-20210904204757717](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/81802394179e462b8288ba6437c816f7~tplv-k3u1fbpfcp-watermark.awebp)

## Nginx安装

### windows下安装

**1、下载nginx**

[nginx.org/en/download…](https://link.juejin.cn/?target=http%3A%2F%2Fnginx.org%2Fen%2Fdownload.html) 下载稳定版本。以nginx/Windows-1.20.1为例，直接下载 nginx-1.20.1.zip。 下载后解压，解压后如下：

![image-20210905103735775](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/ccc0a0e77e7e407288ed14d10e692134~tplv-k3u1fbpfcp-watermark.awebp)

**2、启动nginx**

-   直接双击nginx.exe，双击后一个黑色的弹窗一闪而过
-   打开cmd命令窗口，切换到nginx解压目录下，输入命令 `nginx.exe` ，回车即可

**3、检查nginx是否启动成功**

直接在浏览器地址栏输入网址 [http://localhost:80](https://link.juejin.cn/?target=http%3A%2F%2Flocalhost%2F) 回车，出现以下页面说明启动成功！

![image-20210905103934702](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/eb9d673c19014693bfa910b35e5fe246~tplv-k3u1fbpfcp-watermark.awebp)

### Docker安装nginx

我之前的文章也讲过Linux下安装的步骤，我采用的是docker安装的，很简单。

相关链接如下：[Docker（三）：Docker部署Nginx和Tomcat](https://juejin.cn/post/6937814934646423560)

**1、查看所有本地的主机上的镜像**，使用命令`docker images`

![image-20210904232433522](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/eb55d059328d4522abea54ce04d33772~tplv-k3u1fbpfcp-watermark.awebp) **2、创建 nginx 容器 并启动容器**，使用命令`docker run -d --name nginx01 -p 3344:80 nginx`

![image-20210904233936797](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/b41fdf3ad1384caeb32b409f02afb8c6~tplv-k3u1fbpfcp-watermark.awebp)

**3、查看已启动的容器**，使用命令`docker ps`

![image-20210904234231750](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/cb62c7a5756b4f35b89469a994a21740~tplv-k3u1fbpfcp-watermark.awebp)

浏览器访问`服务器ip:3344`，如下，说明安装启动成功。

`注意`：如何连接不上，检查阿里云安全组是否开放端口，或者服务器防火墙是否开放端口！

![image-20210905104039595](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/d7c5106e35d84ff99ff4d88a593b13e3~tplv-k3u1fbpfcp-watermark.awebp)

### linux下安装

**1、安装gcc**

安装 nginx 需要先将官网下载的源码进行编译，编译依赖 gcc 环境，如果没有 gcc 环境，则需要安装：

```shell
yum install gcc-c++
复制代码
```

**2、PCRE pcre-devel 安装**

PCRE(Perl Compatible Regular Expressions) 是一个Perl库，包括 perl 兼容的正则表达式库。nginx 的 http 模块使用 pcre 来解析正则表达式，所以需要在 linux 上安装 pcre 库，pcre-devel 是使用 pcre 开发的一个二次开发库。nginx也需要此库。命令：

```shell
yum install -y pcre pcre-devel
复制代码
```

**3、zlib 安装**

zlib 库提供了很多种压缩和解压缩的方式， nginx 使用 zlib 对 http 包的内容进行 gzip ，所以需要在 Centos 上安装 zlib 库。

```shell
yum install -y zlib zlib-devel
复制代码
```

**4、OpenSSL 安装**

OpenSSL 是一个强大的安全套接字层密码库，囊括主要的密码算法、常用的密钥和证书封装管理功能及 SSL 协议，并提供丰富的应用程序供测试或其它目的使用。 nginx 不仅支持 http 协议，还支持 https（即在ssl协议上传输http），所以需要在 Centos 安装 OpenSSL 库。

```shell
yum install -y openssl openssl-devel
复制代码
```

**5、下载安装包**

手动下载.tar.gz安装包，地址：[nginx.org/en/download…](https://link.juejin.cn/?target=https%3A%2F%2Fnginx.org%2Fen%2Fdownload.html)

![image-20210905173049111](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/3deb20f65a8c474ca974e626624aeae6~tplv-k3u1fbpfcp-watermark.awebp)

下载完毕上传到服务器上 /root

**6、解压**

```shell
tar -zxvf nginx-1.20.1.tar.gz
cd nginx-1.20.1
复制代码
```

![image-20210905173212111](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/43d44bdd6f9640f7a6b929e9d13c01c0~tplv-k3u1fbpfcp-watermark.awebp)

**7、配置**

使用默认配置，在nginx根目录下执行

```shell
./configue
make
make install
复制代码
```

查找安装路径： `whereis nginx`

![image-20210905181408981](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/3c9269ca4dc44f4d89088d1418149d4b~tplv-k3u1fbpfcp-watermark.awebp)

8、启动 nginx

```shell
./nginx
复制代码
```

![image-20210905181510315](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/2a3c90205f3143d0b445ccd5ef22c8d2~tplv-k3u1fbpfcp-watermark.awebp)

启动成功，访问页面：ip:80

![image-20210905181740776](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/972e3378701b4248a2068876985c13e3~tplv-k3u1fbpfcp-watermark.awebp)

## Nginx常用命令

```
注意`：使用Nginx操作命令前提，必须进入到Nginx目录 `/usr/local/nginx/sbin
```

1、查看Nginx版本号：`./nginx -v`

![image-20210905203751070](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/e48f2e1f4c0341c586361719b374847d~tplv-k3u1fbpfcp-watermark.awebp)

2、启动 Nginx：`./nginx`

![image-20210905201929397](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/82363cd8b62441a092fa28b79f46499d~tplv-k3u1fbpfcp-watermark.awebp)

3、停止 Nginx：`./nginx -s stop` 或者`./nginx -s quit`

![image-20210905202644529](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/9ecaf7f50c8f4fcdaeb43dbda5585542~tplv-k3u1fbpfcp-watermark.awebp)

4、重新加载配置文件：`./nginx -s reload`

![image-20210905202753783](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/92630cf223c6409aaeb4ce9952d5fe76~tplv-k3u1fbpfcp-watermark.awebp)

5、查看nginx进程：`ps -ef|grep nginx`

![image-20210905202618893](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/591da4266e5b4ace98c573fa706477ea~tplv-k3u1fbpfcp-watermark.awebp)

## Nginx配置文件

Nginx配置文件的位置：`/usr/local/nginx/conf/nginx.conf`

![image-20210905204225730](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/58c6bb46f42a46d08afe95acd795c1f9~tplv-k3u1fbpfcp-watermark.awebp)

Nginx配置文件有3部分组成：


![image-20210906151628317](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/e97272cf5eb84ed49401fbbcb630e36b~tplv-k3u1fbpfcp-watermark.awebp)

**1、全局块**

从配置文件开始到 events 块之间的内容，主要会设置一些影响 nginx 服务器整体运行的配置指令，比如：`worker_processes 1`。

这是 Nginx 服务器并发处理服务的关键配置，worker_processes 值越大，可以支持的并发处理量也越多，但是会受到硬件、软件等设备的制约。一般设置值和CPU核心数一致。

**2、events块**

events 块涉及的指令主要影响 Nginx 服务器与用户的网络连接，比如：`worker_connections 1024`

表示每个 work process 支持的最大连接数为 1024，这部分的配置对 Nginx 的性能影响较大，在实际中应该灵活配置。

**3、http块**

```shell
http {
    include       mime.types;
    
    default_type  application/octet-stream;

    sendfile        on;

    keepalive_timeout  65;

    server {
        listen       80;#监听端口
        server_name  localhost;#域名

        location / {
            root   html;
            index  index.html index.htm;
        }

        error_page   500 502 503 504  /50x.html;
       
        location = /50x.html {
            root   html;
        }

    }

}
复制代码
```

这算是 Nginx 服务器配置中最频繁的部分。

## 演示示例

### 反向代理/负载均衡

![image-20210906150245880](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/3f002567c5574bf298a27b83685a1693~tplv-k3u1fbpfcp-watermark.awebp)

我们在windows下演示，首先我们创建两个springboot项目，端口是9001和9002，如下：

![image-20210906142059214](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/3eb33424d0fa4d61a03e18bdc6881513~tplv-k3u1fbpfcp-watermark.awebp)

我们要做的就是将`localhost:80`代理`localhost:9001`和`localhost:9002`这两个服务，并且让轮询访问这两个服务。

nginx配置如下：

```java
worker_processes  1;

events {
    worker_connections  1024;
}

http {
    include       mime.types;
    default_type  application/octet-stream;

    sendfile        on;
    keepalive_timeout  65;

	upstream jiangwang {
		server 127.0.0.1:9001 weight=1;//轮询其权重都默认为1
		server 127.0.0.1:9002 weight=1;
	}

    server {
        listen       80;
        server_name  localhost;

        #charset koi8-r;

        #access_log  logs/host.access.log  main;

        location / {
            root   html;
            index  index.html index.htm;
            proxy_pass http://jiangwang;
        }
    }

}
复制代码
```

我们先将项目打成jar包，然后命令行启动项目，然后在浏览器上访问`localhost`来访问这两个项目，我也在项目中打印了日志，操作一下来看看结果，是不是两个项目轮询被访问。

![image-20210906144202974](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/43259817ef54486a8d3e07b44e3abe38~tplv-k3u1fbpfcp-watermark.awebp)

![image-20210906144224135](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/23e2c01656814a399a08de8635e39865~tplv-k3u1fbpfcp-watermark.awebp)

可以看到，访问`localhost`，这两个项目轮询被访问。

接下来我们将权重改为如下设置：

```java
upstream jiangwang {
    server 127.0.0.1:9001 weight=1;
    server 127.0.0.1:9002 weight=3;
}
复制代码
```

重新加载一个nginx的配置文件：`nginx -s reload`

加载完毕，我们再访问其`localhost`，观察其访问的比例：

![image-20210906145104854](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/a917285dede645cfad80472b85774ee9~tplv-k3u1fbpfcp-watermark.awebp)

![image-20210906145132494](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/a642a970ee2746e5ba8487289f2afb3d~tplv-k3u1fbpfcp-watermark.awebp)

结果显示，9002端口的访问次数与9001访问的次数基本上是`3:1`。

### 动静分离

1、将静态资源放入本地新建的文件里面，例如：在D盘新建一个文件data，然后再data文件夹里面在新建两个文件夹，一个img文件夹，存放图片；一个html文件夹，存放html文件；如下图：

![image-20210906234145869](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/86dbdd39645549c3b581fa0b0b9b0d50~tplv-k3u1fbpfcp-watermark.awebp)

2、在html文件夹里面新建一个`a.html`文件，内容如下：

```html
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Html文件</title>
</head>
<body>
    <p>Hello World</p>
</body>
</html>
复制代码
```

3、在img文件夹里面放入一张照片，如下：

![img](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/e24c44c6467942e8a9fb8573d271d330~tplv-k3u1fbpfcp-watermark.awebp)

4、配置nginx中`nginx.conf`文件：

```java
location /html/ {
    root   D:/data/;
    index  index.html index.htm;
}

location /img/ {
    root   D:/data/;
    autoindex on;#表示列出当前文件夹中的所有内容
}
复制代码
```

5、启动nginx，访问其文件路径，在浏览器输入`http://localhost/html/a.html`，如下：

![image-20210906233944234](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/2788b426c72645e1a63daeb9be785eda~tplv-k3u1fbpfcp-watermark.awebp)

6、在浏览器输入`http://localhost/img/`

![image-20210906234039557](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/9fbb28181334479f8a1c8e65679c4731~tplv-k3u1fbpfcp-watermark.awebp)

## Nginx工作原理

### mater&worker

![image-20210907084101801](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/c3fe421794bb4bccb2164d63a002f886~tplv-k3u1fbpfcp-watermark.awebp)

master接收信号后将任务分配给worker进行执行，worker可有多个。

![image-20210906235920093](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/420fa25ad791410f844fa3329392cfd7~tplv-k3u1fbpfcp-watermark.awebp)

### worker如何工作

客户端发送一个请求到master后，worker获取任务的机制不是直接分配也不是轮询，而是一种争抢的机制，“抢”到任务后再执行任务，即选择目标服务器tomcat等，然后返回结果。

![image-20210907104204828](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/c403846424f042f5ad5b83b286f96acc~tplv-k3u1fbpfcp-watermark.awebp)

### worker_connection

发送请求占用了woker两个或四个连接数。

普通的静态访问最大并发数是：`worker_connections * worker_processes/ 2` ，若是 HTTP 作为反向代理来说，最大并发数量应该是 `worker_connections * worker_processes/ 4` 。当然了，worker数也不是越多越好，worker数和服务器的CPU数相等时最适宜的。

### 优点

可以使用 `nginx –s reload` 热部署，利用 nginx 进行热部署操作每个 woker 是独立的进程，若其中一个woker出现问题，其他继续进行争抢，实现请求过程，不会造成服务中断。

## 总结

关于 Nginx 的基本概念、安装教程、配置、使用实例以及工作原理，本文都做了详细阐述。希望本文对你有所帮助。

## 抽奖细则

1.  首先你要参与评论，我希望你能看完文章，不要只打个表情符号之类的。
2.  记得`点赞、关注`，动动你的小手指，谢谢！
3.  至于怎么公平，我也想了一下，我会创建一个抽奖群，我发红包，你们抢，金额最少的获得奖品。或者你们有更好的方式评论区告诉我。如果我觉得更好，我会采纳。
4.  如果我的评论区热度在Top 1-5 名，获得的 新版徽章1套 或 掘金新版IPT恤1件，我也会送给评论区的小伙伴，前提你关注我的我的公众号：微信搜【初念初恋】。