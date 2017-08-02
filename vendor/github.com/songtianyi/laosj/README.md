## laosj(老司机)

[![Build Status](https://travis-ci.org/songtianyi/laosj.svg?branch=master)](https://travis-ci.org/songtianyi/laosj)
[![Go Report Card](https://goreportcard.com/badge/github.com/songtianyi/laosj)](https://goreportcard.com/report/github.com/songtianyi/laosj)
[![codebeat badge](https://codebeat.co/badges/c05ec05d-e902-4091-b5e0-c1656f88ae3c)](https://codebeat.co/projects/github-com-songtianyi-laosj)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

[![logo](https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTaiDDQDv9P90h7lu9jARb1O8i6hmVMpgEuK9qY57l0CZjRVue2)](https://github.com/songtianyi/laosj)


基于goquery的轻量级爬虫, 图片下载器支持分布式, 可以自定义下载并发量。

## Releases
* mzitu 爬取mzitu.com/taiwan (图片下载依赖redis作为url队列)
* mzituzp 爬取mzitu.com/share (图片下载依赖redis作为url队列)
* aiss-darwin 根据接口下载爱丝图片(依赖redis作为url队列)
* jav 获取番号列表及根据番号获取种子文件 (无依赖)

## 代码上手
######下载
    go get -u -v github.com/songtianyi/laosj
######安装redis
	略
######golang.org/x依赖安装
```
mkdir $GOPATH/src/golang.org/x
cd $GOPATH/src/golang.org/x
git clone https://github.com/golang/net.git
```
######编译并运行样例
	go build mzituzp.go
	./mzituzp

## 截图

![laosj](http://i1.piimg.com/4851/a598ac03cd7ae15f.jpg)
