# FileService 简单文件服务器

    支持简单的路径挂载, 文件基本操作(CRUD). windows安装完成后会自动生成一个名字为FileService的服务, Linux目前为手动安装使用.

## 基本功能

| 模块          | 描述                                                                                   |
| ------------- | -------------------------------------------------------------------------------------- |
| Web服务       | 支持 http、https 协议                                                                  |
| 用户信息      | 支持简单的用户管理操作                                                                 |
| 文件权限      | 支持简单的文件夹授权操作, 暂时只支持: 可见、只读、可写三个维度的权限                   |
| 数据存取      | 支持基本的文件管理操作, 提供虚拟目录挂载、基础的文件在线预览等功能                          |

## 源码目录结构

| 项      | 描述                                                                    |
| ------- | ----------------------------------------------------------------------- |
| \_docs  | 文档存放目录, 如接口定义和测试信息                                         |
| \_setup | windows 打包物料和安装包目录, 打包好的安装包在`_setup/package-result`下  |
| server  | 后端接口源码, 包含 windows 和 Linux 编译好的可执行文件                  |
| web-ui  | 前端页面源码, 编译目标目录在`server/webapps`下                          |

## 程序初始化信息

| 项       | 值                      | 描述               |
| -------- | ----------------------- | ------------------ |
| 用户     | admin 密码空            | 进入系统后可改变 |
| 监听地址 | `http://127.0.0.1:8080` | 可通过修改配置改变 |

## 程序安装目录信息

| 目录     | 文件             | 描述                                       |
| -------- | ---------------- | ------------------------------------------ |
| /.conf    | fileservice.json | 程序配置文件, 如: 目录挂载信息、监听端口等 |
| /.conf    | key.pem,cert.pem | https 证书存放位置, 有配置则启用 https     |
| /webapps | \*               | web 页面                                   |

## 程序配置清单

| KEY                   | 默认值                              | 可选值 | 描述           |
| --------------------- | ----------------------------------- | ------ | -------------- |
| `listen.http.address` | 127.0.0.1:8080                      | `*`    | 服务对外端口   |
| `filedatas.mount./`   | `{"addr":"./datas","type":"LOCAL"}` | `*`    | 根目录挂载位置 |

    . 配置使用json格式存储, 格式示例:
        `{
            "listen": {
                "http": {
                    "address": "0.0.0.0:8080"
                }
            },
            "filedatas": {
                "mount": {
                    "/挂载的虚拟路径": {
                        "addr": "本地磁盘位置, 绝对位置或相对位置",
                        "type": "驱动类型, 默认 LOCAL"
                    }
                }
            }
        }`

## 目录挂载支持的文件系统

| 类型                   | 示例                              | 描述           |
| --------------------- | ----------------------------------- | -------------- |
| LOCAL | `{"addr":"./datas","type":"LOCAL"}` | 本地磁盘系统   |
| OSS | `{"addr": "OPENAPI@127.0.0.1:5051/mountdir","type": "OSS","passwd": ""},` | [对象存储](https://github.com/wup364/filestorage) |
