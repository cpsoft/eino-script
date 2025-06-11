# Eino-Script

## 🚀 概述

Eino-Script是一个基于eino开发的脚本驱动的AI工作流引擎。开发项目的目的，一是为了尝试Eino动态运行工作流的能力，二来是为了实现类似coze等在线编辑工作流的基础引擎。开源该引擎是为了方便将类似的引擎集成到私有化的项目中去。

### 🌟 核心特性

 - ⛰️ **多脚本格式支持**：通过viper库，支持toml、yaml、json等脚本格式
 - ⚡️ **支持多来源脚本**：支持文件脚本，也支持内存脚本（用于网络传递）

### 🥔 安装及使用

#### 由于eino的依赖关系问题，github.com/getkin/kin-openapi库 必须指定版本 v0.118.0
```bash
go get github.com/getkin/kin-openapi@v0.118.
```

使用：
#### 文件脚本方式：
```bash
eino-script -file [配置文件路径]
```

#### 服务器方式：
```bash
eino-script -server
```
采用gin作为服务器。客户端在
https://gitee.com/cpsoft13/ai_flow

需要 Go 1.18 或更高版本

## 🍉 开源协议

本项目采用 MIT 协议 - 详见 [LICENSE](LICENSE) 文件
