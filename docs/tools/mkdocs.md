---
title: "Mkdocs"
date: 2021-05-31T17:31:16+08:00
draft: false
tags: ["项目管理"]
---


[MKDocs](https://www.mkdocs.org/) 是一个快速、简单、华丽的静态网站生成器，用于构建项目文档。
文档源文件是用Markdown编写的，并使用单个YAML配置文件进行配置。并构建的是完整静态的HTML站点，你可以托管静态页面到GtuhubPages、Amzone S3、AliCloud OSS、Nginx静态页。

## 安装

如果你有和使用包管理器(如apt-get,dnf,homebrew,yum,chocolatey等等),那么您可能想要搜索“MkDocs”包并安安装。例如：centos上使用
```shell
yum install mkdocs
```

### 手动安装前检查

如果你的机器上未安装包管理器，你仍然可以使用 `python`和`pip`安装.
为了可以手动安装`MkDocs`，你需要去检查`python`和包管理器`pip`是否已经安装成功。你可以使用一下命令检查是否安装。

```shell
$ python --version
Python 3.8.2
$ pip --version
pip 20.0.2 from /usr/local/lib/python3.8/site-packages/pip (python 3.8)
```

> MkDocs 支持 Python 版本为 3.5, 3.6, 3.7, 3.8, and pypy3.

### 手动安装

通过下面的命令安装：
```shell
pip install mkdocs
```

安装完成后使用：
```shell
$ mkdocs --version
mkdocs, version 0.15.3
```

## 入门指南

### 创建一个MKDocs项目
```shell
mkdocs new my-project
```

当执行成功后，你可以看到项目文件
```
my-project
	docs
		index.md
	mkdocs.yml
```

有一个名为mkdocs的配置文件。以及一个名为docs的文件夹，该文件夹将包含您的文档源文件， 现在docs文件夹只包含一个单独的文档页面，名为index.md。

### 开发实时预览

MkDocs带有内置的开发服务器，可以让你在工作时预览文档。确保您在与mkdocs相同的目录中。然后运行mkdocs serve命令启动服务器:

```shell
$ mkdocs serve
INFO    -  Building documentation...
INFO    -  Cleaning site directory
[I 160402 15:50:43 server:271] Serving on http://127.0.0.1:8000
[I 160402 15:50:43 handlers:58] Start watching changes
[I 160402 15:50:43 handlers:60] Start detecting changes
```
此时你就可以打开网页实时查看编写结果。

### mkdocs.yml 配置文件
```yaml
site_name: "网站名称"
nav:
    - Home: index.md
    - About: about.md
theme: readthedocs
```

site_name: 这里配置网站名称
theme： 配置文档生成主题
nav: 配置生成文档的导航条

在nav下面我们可以配置多个导航。


### 构建

在项目文件夹下面直接使用
```shell
mkdocs build 
```
执行命令后，我们将在当前文件夹下面看到生成的静态网站文件夹`site`.

### Github Pages 发布



## [更多MKDocs信息](https://www.mkdocs.org/)