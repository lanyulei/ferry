# IDE开发

众多 IDE 里边，推荐使用 `goland IDE`进行调试

首先我们启动 `Goland` , 点击 `Open Project`，下图红框圈选部分；

![](https://www.fdevops.com/wp-content/uploads/2020/08/image.png)

选择 ferry 存放的路径，找到并打开；

# 配置 GOPORXY

然后选择`Goland` > `Preferences` ；

![](https://www.fdevops.com/wp-content/uploads/2020/08/image-1.png)

# 添加运行或调试配置

### 添加 init 配置

1\. 打开`Edit Configurations`

![](https://www.fdevops.com/wp-content/uploads/2020/08/image-2.png)

2\. 选择 `+` > `go build`

![](https://www.fdevops.com/wp-content/uploads/2020/08/image-3.png)

3\. 按照下图所示进行配置，注意：填写 `Program arguments` 为 `init -c=config/settings.dev.yml`，完成之后点击保存

![](https://www.fdevops.com/wp-content/uploads/2020/08/image-4.png)

4\. 修改数据库

![](https://www.fdevops.com/wp-content/uploads/2020/08/image-5.png)

5\. 初使化

![](https://www.fdevops.com/wp-content/uploads/2020/08/image-6.png)

### 添加 server 配置

1\. 打开`Edit Configurations`

![](https://www.fdevops.com/wp-content/uploads/2020/08/image-7.png)

2\. 选择 `+` > `go build`

![](https://www.fdevops.com/wp-content/uploads/2020/08/image-8.png)

3\. 按照下图所示进行配置，注意：填写 `Program arguments` 为`server -c=config/settings.dev.yml`，完成之后点击保存

![](https://www.fdevops.com/wp-content/uploads/2020/08/image-9.png)

4\. 启动 debug

![](https://www.fdevops.com/wp-content/uploads/2020/08/image-10.png)

转载自：[http://doc.zhangwj.com/go-admin-site/guide/ide.html#%E6%B7%BB%E5%8A%A0-server-%E9%85%8D%E7%BD%AE](http://doc.zhangwj.com/go-admin-site/guide/ide.html#%E6%B7%BB%E5%8A%A0-server-%E9%85%8D%E7%BD%AE)
