# 简介

本系统是集工单统计、任务钩子、权限管理、灵活配置流程与模版等等于一身的开源工单系统，当然也可以称之为工作流引擎。

致力于减少跨部门之间的沟通，自动任务的执行，提升工作效率与工作质量，减少不必要的工作量与人为出错率。

演示Demo: [http://fdevops.com:8001/](http://fdevops.com:8001/)

账号密码：admin/123456

Github: [https://github.com/lanyulei/ferry](https://github.com/lanyulei/ferry)

Gitee: [https://gitee.com/yllan/ferry](https://gitee.com/yllan/ferry)

文档：[https://www.fdevops.com/docs/ferry-tutorial-document/introduction](https://www.fdevops.com/docs/ferry-tutorial-document/introduction)

演示Demo上，将删除的功能全部隐藏了，因为之前发生过，有人恶意删除所有可删除的数据，包括流程数据和用户数据，因此，clone下来的代码是有删除之类的动作的。

## 功能

下面对本系统的功能做一个简单介绍。

工单系统相关功能：

*   工单提交申请
*   工单统计
*   多维度工单列表，包括（我创建的、我相关的、我待办的、所有工单）
*   自定义流程
*   自定义模版
*   任务钩子
*   任务管理
*   催办
*   转交
*   手动结单
*   加签
*   多维度处理人，包括（个人，变量(创建者、创建者负责人)）
*   排他网关，即根据条件判断进行工单跳转
*   并行网关，即多个节点同时进行审批处理
*   通知提醒（目前仅支持邮件）
*   流程分类管理

权限管理相关功能，使用casbin实现接口权限控制：

*   用户、角色、岗位的增删查改，批量删除，多条件搜索
*   角色、岗位数据导出Excel
*   重置用户密码
*   维护个人信息，上传管理头像，修改当前账户密码
*   部门的增删查改
*   菜单目录、跳转、按钮及API接口的增删查改
*   登陆日志管理
*   左菜单权限控制
*   页面按钮权限控制
*   API接口权限控制

目前大致上就是以上功能了，如果您觉得我有拉下的功能，还请留言提醒我，感谢。
