<!--
 * @Author: Z-Es-0 zes18642300628@qq.com
 * @Date: 2025-04-15 23:23:42
 * @LastEditors: Z-Es-0 zes18642300628@qq.com
 * @LastEditTime: 2025-04-15 23:57:55
 * @FilePath: \ZesOJ\Readme.md
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
-->
# go-dbg

> Golang implementation of the debugger

<!-- ![图片由Grok 3 生成](./pic/image.jpg ) -->


<img src="./pic/image.jpg" alt="图片由Grok 3 生成" width="300">

`go-dbg` 是一个用 Go 语言实现的调试器项目，目前支持 PE 文件解析、进程调试、内存操作和反汇编等简单功能。



### **1. 环境**
- 安装 [Go](https://golang.org/)（版本 1.20 或更高）。
- 确保你的系统含有必要的dll。

### **2. 本地编译**
```sh
git clone https://github.com/Z-Es-0/go-dbg.git
cd go-dbg
make build
```

### **3. 目前功能**

具有唯一的命令行参数 <被调试的进程>

```sh
debug.exe sever/test.exe
```


run, r: 继续运行程序

break < address >, b < address >: 设置断点

memory  < address > < size >, m  < address > < size >: 查看内存

context, c: 查看线程上下文(当前寄存器状态)


quit, q: 退出程序

list, l: 列出断点

help, h: 查看帮助



