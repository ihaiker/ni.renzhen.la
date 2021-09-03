---
title: "国内node各种镜像源设置"
---

设置方式：

1、`shell`

```shell
npm set -g registry https://registry.npm.taobao.org
npm set -g disturl https://npm.taobao.org/dist
npm set -g puppeteer_download_host https://npm.taobao.org/mirrors
npm set -g electron_mirror https://npm.taobao.org/mirrors/electron/
npm set -g sharp_binary_host https://npm.taobao.org/mirrors/sharp
npm set -g sharp_libvips_binary_host https://npm.taobao.org/mirrors/sharp-libvips
```

2、在 `$HOME/.npmrc` 中设置

```properties
registry=https://registry.npm.taobao.org
disturl=https://npm.taobao.org/dist
puppeteer_download_host=https://npm.taobao.org/mirrors
electron_mirror=https://npm.taobao.org/mirrors/electron/
sharp_binary_host=https://npm.taobao.org/mirrors/sharp
sharp_libvips_binary_host=https://npm.taobao.org/mirrors/sharp-libvips
```

