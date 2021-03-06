---
title: 国内镜像加速
date: 2021-02-02T11:12:12+08:00
tags: docker
---

安装Kubernetes默认的image源是：k8s.gcr.io，中国区安装极为痛苦，不FQ就拉取不下来，类似的还有Red Head的quay.io，虽然能下载，但是速度非常的慢。



解决办法来了：

中国区有了以上的docker registry镜像源了：

中科大镜像 https://github.com/ustclug/mirrorrequest

Azure中国镜像 https://github.com/Azure/container-service-for-azure-china

先上总结，方便直接查阅：

## 总结

鉴于中科大的源经常的出现not fount的错误，所以这里使用Azure的镜像源。

### docker.io

#### docker官方维护

```
$ docker pull dockerhub.azk8s.cn/library/xxx:yyy
```

#### 非官方维护

```
$ docker pull dockerhub.azk8s.cn/uuu/xxx:yyy
```

### gcr.io 和 k8s.gcr.io

#### gcr.io

```
$ docker pull gcr.azk8s.cn/xxx/yyy:zzz
```

#### k8s.gcr.io

```
$ docker pull gcr.azk8s.cn/google-containers/xxx:yyy
```

### quay.io

```
$ docker pull quay.azk8s.cn/xxx/yyy:zzz
```

下面详细讲解：

## docker.io 镜像加速

默认的docker registry是：[https://hub.docker.com](https://hub.docker.com/)

### 使用中科大镜像 docker.mirrors.ustc.edu.cn

默认的一种拉取方式：

```
$ docker pull xxx:yyy

这里没有指定用户，默认是docker官方维护的镜像
xxx为镜像的名称
yyy为镜像的tag版本
```

那么使用中科大镜像的方式：

```
$ docker pull docker.mirrors.ustc.edu.cn/library/xxx:yyy
```

例如：nginx:1.16-alpine

```
$ docker pull docker.mirrors.ustc.edu.cn/library/nginx:1.16-alpine
```

默认的另一种拉取方式：

```
$ docker pull uuu/xxx:yyy

uuu为具体的用户，有其负责维护该镜像
xxx为镜像的名称
yyy为镜像的tag版本
```

那么使用中科大镜像的方式：

```
$ docker pull docker.mirrors.ustc.edu.cn/uuu/xxx:yyy
```

例如：wanghkkk/busyboxplus:latest

```
$ docker pull docker.mirrors.ustc.edu.cn/wanghkkk/busyboxplus
```

默认中科大会提示：Error response from daemon: manifest for docker.mirrors.ustc.edu.cn/wanghkkk/busyboxplus:latest not found

因为默认中科大不会换成镜像的，需要多试几次，等中科大缓存完成后才能下载的。**所以这里不推荐用中科大的**。

### 使用Azure中国镜像 dockerhub.azk8s.cn

还是使用上面的两个例子：

nginx:1.16-alpine

wanghkkk/busyboxplus

对应的Azure中国镜像是：

```shell
$ docker pull dockerhub.azk8s.cn/library/nginx:1.16-alpine
$ docker pull dockerhub.azk8s.cn/wanghkkk/busyboxplus
```

这里推荐使用Azure，直接就拉取下来了，不像中科大需要他自己先下载缓存，然后才能下载，至少我试了好几次也没拉取下来。

## gcr.io 和 k8s.gcr.io 镜像加速

gcr.io 和 k8s.gcr.io 实际上都是Google的镜像，默认中国区是根本访问不到的。

### gcr.io

#### 使用中科大镜像 gcr.mirrors.ustc.edu.cn

默认gcr.io拉取方式：

```
$ docker pull gcr.io/xxx/yyy:zzz
```

那么更换成中科大的拉取方式：

```
$ docker pull gcr.mirrors.ustc.edu.cn/xxx/yyy:zzz
```

例如：

这里使用gcr.io/kubernetes-helm/tiller:v2.16.1为例：

```
$ docker pull gcr.mirrors.ustc.edu.cn/kubernetes-helm/tiller:v2.16.1
```

#### 使用Azure中国镜像 gcr.azk8s.cn

更换成Azure的拉取方式为：

```
$ docker pull gcr.azk8s.cn/kubernetes-helm/tiller:v2.16.1
```

### k8s.gcr.io

对于按照kubernetes时候，用到的就是k8s.gcr.io开头的镜像，其实k8s.gcr.io就是等价于gcr.io/google-containers

例如：k8s.gcr.io/kube-proxy:v1.15.5

所以呢，对于中科大的拉取方式为：

```
$ docker pull gcr.mirrors.ustc.edu.cn/google-containers/kube-proxy:v1.15.5
```

对于Azure的拉取方式为：

```
$ docker pull gcr.azk8s.cn/google-containers/kube-proxy:v1.15.5
```

对于中科大的拉取，经常的出现not found，而对于Azure来说非常的顺畅，很快的拉取下来了，所以非常的推荐使用Azure的镜像源。

## quay.io 镜像加速

例如拉取镜像：

quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.26.1

### 使用中科大镜像 quay.mirrors.ustc.edu.cn

```
$ docker pull quay.mirrors.ustc.edu.cn/kubernetes-ingress-controller/nginx-ingress-controller:0.26.1
```

### 使用Azure中国镜像 quay.azk8s.cn

```
$ docker pull quay.azk8s.cn/kubernetes-ingress-controller/nginx-ingress-controller:0.26.1
```

## 一个python脚本全搞定

docker-wrapper.py

```
#!/usr/bin/python
# coding=utf8

import os
import sys

# azure mirrors for gcr.io,k8s.gcr.io,quay.io in china
gcr_mirror = "gcr.azk8s.cn"
docker_mirror = "dockerhub.azk8s.cn"
quay_mirror = "quay.azk8s.cn"

k8s_namespace = "google_containers"

gcr_prefix = "gcr.io"
special_gcr_prefix = "k8s.gcr.io"
quay_prefix = "quay.io"


def execute_sys_cmd(cmd):
    result = os.system(cmd)
    if result != 0:
        print(cmd + " failed.")
        sys.exit(-1)


def usage():
    print("Usage: " + sys.argv[0] + " pull ")
    print("Examples:")
    print(sys.argv[0] + " pull k8s.gcr.io/kube-apiserver:v1.14.1")
    print(sys.argv[0] + " pull gcr.io/google_containers/kube-apiserver:v1.14.1")
    print(sys.argv[0] + " pull quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.26.1")


if __name__ == "__main__":
    if len(sys.argv) != 3:
        usage()
        sys.exit(-1)
    elif sys.argv[1] != 'pull':
        usage()
        sys.exit(-1)

    # image name like k8s.gcr.io/kube-apiserver:v1.14.1 or gcr.io/google_containers/kube-apiserver:v1.14.1
    image = sys.argv[2]
    imageArray = image.split("/")

    if imageArray[0] == gcr_prefix:
        imageArray[0] = gcr_mirror
    elif imageArray[0] == special_gcr_prefix:
        imageArray[0] = gcr_mirror
        imageArray.insert(1, k8s_namespace)
    elif imageArray[0] == quay_prefix:
        imageArray[0] = quay_mirror
    elif len(imageArray) == 1:
        imageArray.insert(0, docker_mirror)
        imageArray.insert(1, "library")
    elif len(imageArray) == 2:
        imageArray.insert(0, docker_mirror)

    temp_image = "/".join(imageArray)

    cmd = "docker pull {image}".format(image=temp_image)
    print("------Execute_cmd: %s" % cmd)
    execute_sys_cmd(cmd)

    cmd = "docker tag {newImage} {image}".format(newImage=temp_image, image=image)
    print("------Execute_cmd: %s" % cmd)
    execute_sys_cmd(cmd)

    cmd = "docker rmi {newImage}".format(newImage=temp_image)
    print("------Execute_cmd: %s" % cmd)
    execute_sys_cmd(cmd)

    print("------Pull %s done" % image)
    sys.exit(0)
```

使用方式：

```
$ chmod +x docker-wrapper.py
$ sudo mv docker-wrapper.py /usr/local/bin/
$ docker-wrapper.py pull xxx/yyy:zzz
```




## Docker Hub 镜像加速器列表

|镜像加速器|镜像加速器地址|专属加速器？|其它加速？|
|:---|---|---|---|
|Docker 中国官方镜像|https://registry.docker-cn.com|Docker Hub||
|DaoCloud 镜像站|http://f1361db2.m.daocloud.io|可登录，系统分配|Docker Hub|
|Azure 中国镜像|https://dockerhub.azk8s.cn|Docker Hub、GCR、Quay||
|科大镜像站|https://docker.mirrors.ustc.edu.cn|Docker Hub、GCR、Quay||
|阿里云|https://<your_code>.mirror.aliyuncs.com|需登录，系统分配|Docker Hub|
|七牛云|https://reg-mirror.qiniu.com||Docker Hub、GCR、Quay|
|网易云|https://hub-mirror.c.163.com||Docker Hub|
|腾讯云|https://mirror.ccs.tencentyun.com||Docker Hub|
