<p align="center">
  <img src="https://www.fdevops.com/wp-content/uploads/2020/09/1599039924-ferry_log.png">
</p>


<p align="center">
  <a href="https://github.com/lanyulei/ferry">
    <img src="https://www.fdevops.com/wp-content/uploads/2020/07/1595067271-badge.png">
  </a>
  <a href="https://github.com/lanyulei/ferry">
    <img src="https://www.fdevops.com/wp-content/uploads/2020/07/1595067272-apistatus.png" alt="license">
  </a>
    <a href="https://github.com/lanyulei/ferry">
    <img src="https://www.fdevops.com/wp-content/uploads/2020/07/1595067269-donate.png" alt="donate">
  </a>
</p>

## 基于Gin + Vue + Element UI前后端分离的工单系统

**流程中心**

通过灵活的配置流程、模版等数据，非常快速方便的生成工单流程，通过对流程进行任务绑定，实现流程中的钩子操作，目前支持绑定邮件来通知处理，当然为兼容更多的通知方式，也可以自己写任务脚本来进行任务通知，可根据自己的需求定制。

兼容了多种处理情况，包括串行处理、并行处理以及根据条件判断进行节点跳转。

可通过变量设置处理人，例如：直接负责人、部门负责人、HRBP等变量数据。

**系统管理**

基于casbin的RBAC权限控制，借鉴了go-admin项目的前端权限管理，可以在页面对API、菜单、页面按钮等操作，进行灵活且简单的配置。

演示demo: [http://fdevops.com:8001/#/dashboard](http://fdevops.com:8001/#/dashboard)

```
账号：admin
密码：123456

演示demo登陆需要取消ldap验证，就是登陆页面取消ldap的打勾。
```

文档: [https://www.fdevops.com/docs/ferry](https://www.fdevops.com/docs/ferry-tutorial-document/introduction)

视频教程（由群内好友<稳定>提供，非常感谢。）：

* ferry工单系统需要的软件准备 https://www.bilibili.com/video/BV1sA411s7jE
* ferry源代码下载后第一次运行 https://www.bilibili.com/video/BV1oy4y1v7LR

官网：[http://ferry.fdevops.com](http://ferry.fdevops.com)

```
需注意，因有人恶意删除演示数据，将可删除的数据全都删除了，因此演示的Demo上已经将删除操作的隐藏了。

但是直接在Github或者Gitee下载下来的代码是完整的，请放心。

如果总是出现此类删除数据，关闭演示用户的情况的话，可能考虑不在维护demo，仅放置一些项目截图。

请大家一起监督。
```

## 功能介绍

<!-- wp:paragraph -->
<p>下面对本系统的功能做一个简单介绍。</p>
<!-- /wp:paragraph -->

<!-- wp:paragraph -->
<p>工单系统相关功能：</p>
<!-- /wp:paragraph -->

<!-- wp:list -->
<ul><li>工单提交申请</li><li>工单统计</li><li>多维度工单列表，包括（我创建的、我相关的、我待办的、所有工单）</li><li>自定义流程</li><li>自定义模版</li><li>任务钩子</li><li>任务管理</li><li>催办</li><li>转交</li><li>手动结单</li><li>加签</li><li>多维度处理人，包括（个人，变量(创建者、创建者负责人)）</li><li>排他网关，即根据条件判断进行工单跳转</li><li>并行网关，即多个节点同时进行审批处理</li><li>通知提醒（目前仅支持邮件）</li><li>流程分类管理</li></ul>
<!-- /wp:list -->

<!-- wp:paragraph -->
<p>权限管理相关功能，使用casbin实现接口权限控制：</p>
<!-- /wp:paragraph -->

<!-- wp:list -->
<ul><li>用户、角色、岗位的增删查改，批量删除，多条件搜索</li><li>角色、岗位数据导出Excel</li><li>重置用户密码</li><li>维护个人信息，上传管理头像，修改当前账户密码</li><li>部门的增删查改</li><li>菜单目录、跳转、按钮及API接口的增删查改</li><li>登陆日志管理</li><li>左菜单权限控制</li><li>页面按钮权限控制</li><li>API接口权限控制</li></ul>
<!-- /wp:list -->

快速安装部署:
```
bash build.sh install
```

启动服务：
```
bash build.sh start
```

## 交流

加群条件是需给项目一个star，不需要您费多大的功夫与力气，一个小小的star是作者能维护下去的动力。

如果您只是使用本项目的话，您可以在群内提出您使用中需要改进的地方，我会尽快修改。

如果您是想基于此项目二次开发的话，您可以在群里提出您在开发过程中的任何疑问，我会尽快答复并讲解。

群里只要不说骂人、侮辱人之类人身攻击的话，您就可以畅所欲言，有bug我及时修改，使用中有不懂的，我会及时回复，感谢。

<p>
  <img width="300" src="https://www.fdevops.com/wp-content/uploads/2022/10/1667027272-WechatIMG201.jpeg">
</p>

QQ群 4：493389089

[兰玉磊的技术博客](https://www.fdevops.com/)

### 个人微信，添加好友请描述地区、公司及名字，例如：北京-微软-xxx。

本项目现成立了付费的 VIP 交流群，赞助金额达到 299 以上，我将邀请进入付费交流群，本群能更快更优先的协助解决问题。

但是有一点需要跟大家明确好，因为本人需要上班，所以工作日的白天，小问题的话，我会尽快协助解决，大问题需晚上处理呢，大家请按需选择，是否进入，感谢支持。

目前正在使用最新技术重构最新版，最新本也将优先在本群发布。

本人微信：lanyulei666

## 战略合作

* [VForm 一款高效的Vue 2 / Vue3 的低代码表单可视化设计，一键生成源码，享受更多摸鱼时间](https://vform666.com/vform3.html)

## 特别感谢
* [go-amdin # 不错的后台开发框架](https://github.com/go-admin-team/go-admin)

* [vue-element-admin # 不错的前端模版框架](https://github.com/PanJiaChen/vue-element-admin)

* [vue-form-making # 表单设计器，开源版本比较简单，如果有能力的话可以自己进行二次开发 ](https://github.com/GavinZhuLei/vue-form-making.git)

* [wfd-vue # 流程设计器](https://github.com/guozhaolong/wfd-vue)

* [machinery # 任务队列](https://github.com/RichardKnop/machinery.git)

等等...

## 打赏

> 如果您觉得这个项目帮助到了您，您可以请作者喝一杯咖啡表示鼓励:

[打赏名人榜](https://www.fdevops.com/docs/ferry-tutorial-document/reward-celebrity-list)

<img class="no-margin" src="https://www.fdevops.com/wp-content/uploads/2020/07/1595075890-81595075871_.pic_hd.png"  height="200px" >

## 鸣谢

特别感谢 [JetBrains](https://www.jetbrains.com/?from=ferry) 为本开源项目提供免费的 [IntelliJ GoLand](https://www.jetbrains.com/go/?from=ferry) 授权

<p>
 <a href="https://www.jetbrains.com/?from=ferry">
   <img height="200" src="https://www.fdevops.com/wp-content/uploads/2020/09/1599213857-jetbrains-variant-4.png">
 </a>
</p>

## License

[MIT](https://github.com/lanyulei/ferry/blob/master/LICENSE)

Copyright (c) 2022 lanyulei
